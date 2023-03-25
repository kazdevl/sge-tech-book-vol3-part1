package cachedb

import (
	"context"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"sort"
)

//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=cachedb
type BulkExecutor interface {
	BulkInsert(ctx context.Context, tx *sqlx.Tx, contents []CacheContent) error
	BulkUpdate(ctx context.Context, tx *sqlx.Tx, contents []CacheContent) error
	BulkDelete(ctx context.Context, tx *sqlx.Tx, contents []CacheContent) error
}

type CacheDB struct {
	db                   *sqlx.DB
	tx                   *sqlx.Tx
	modelBulkExecutorMap map[string]BulkExecutor
}

func NewCacheDB(db *sqlx.DB, modelBulkExecutorMap map[string]BulkExecutor) *CacheDB {
	return &CacheDB{db: db, modelBulkExecutorMap: modelBulkExecutorMap}
}

func (cdb *CacheDB) Insert(ctx context.Context, content CacheContent) error {
	cacheManager, err := extractDBOperationCacheManager(ctx)
	if err != nil {
		return err
	}

	cacheManager.mutex.Lock()
	defer cacheManager.mutex.Unlock()

	cachedContent, err := cdb.extractCacheContent(ctx, content)
	if err != nil {
		return err
	}
	if cachedContent == nil {
		content.SetCacheStatus(Insert)
		cacheManager.dbOperationResult[content.Table()][content.UniqueKeyColumnValueStr()] = content.CreateCopy()
		return nil
	}

	switch cachedContent.GetCacheStatus() {
	case Insert, Select, Update:
		return errors.New("既に存在する内容をInsertしようとしています")
	case Delete:
		cachedContent.SetCacheStatus(Update)
		return cachedContent.Update(content)
	case None:
		cachedContent.SetCacheStatus(Insert)
		return cachedContent.Update(content)
	default:
		return errors.New("存在しないキャッシュのStatusです")
	}
}

// Update Cache済みのものに対してのみ使用される想定です
func (cdb *CacheDB) Update(ctx context.Context, content CacheContent) error {
	cacheManager, err := extractDBOperationCacheManager(ctx)
	if err != nil {
		return err
	}

	cacheManager.mutex.Lock()
	defer cacheManager.mutex.Unlock()

	cachedContent, err := cdb.extractCacheContent(ctx, content)
	if err != nil {
		return err
	}
	if cachedContent == nil {
		return errors.New("存在しないCacheについてUpdateしようとしています")
	}

	switch cachedContent.GetCacheStatus() {
	case Select, Update:
		cachedContent.SetCacheStatus(Update)
		if err = cachedContent.Update(content); err != nil {
			return err
		}
		return nil
	case Insert:
		if err = cachedContent.Update(content); err != nil {
			return err
		}
		return nil
	case Delete:
		return errors.New("削除済みのCacheに更新できません")
	default:
		return errors.New("存在しないCacheのStatusです")
	}
}

// Delete Cache済みのものに対してのみ使用される想定です
func (cdb *CacheDB) Delete(ctx context.Context, content CacheContent) error {
	cacheManager, err := extractDBOperationCacheManager(ctx)
	if err != nil {
		return err
	}

	cacheManager.mutex.Lock()
	defer cacheManager.mutex.Unlock()

	cachedContent, err := cdb.extractCacheContent(ctx, content)
	if err != nil {
		return err
	}
	if cachedContent == nil {
		return errors.New("存在しないCacheについてDeleteしようとしています")
	}

	switch cachedContent.GetCacheStatus() {
	case Select, Update:
		cachedContent.SetCacheStatus(Delete)
		return nil
	case Insert:
		cachedContent.SetCacheStatus(None)
		return nil
	case Delete:
		return nil
	default:
		return errors.New("存在しないCacheのStatusです")
	}
}

func (cdb *CacheDB) GetAndSet(ctx context.Context, content SelectCacheContent) error {
	cacheManager, err := extractDBOperationCacheManager(ctx)
	if err != nil {
		return err
	}

	cacheManager.mutex.Lock()
	defer cacheManager.mutex.Unlock()

	cachedContent, err := cdb.extractCacheContent(ctx, content)
	if err != nil {
		return err
	}
	if cachedContent == nil {
		// キャッシュに実際に存在しないので、データベースから直接データを取得してCacheに設定する
		query := fmt.Sprintf("SELECT * FROM %s WHERE %s", content.Table(), content.UniqueKeyKCondition())
		rows, err := cdb.db.QueryContext(ctx, query, content.UniqueKeyConditionValues()...)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			if err = content.Bind(rows); err != nil {
				return err
			}
		}
		if err = rows.Err(); err != nil {
			return err
		}

		content.SetCacheStatus(Select)
		cacheManager.dbOperationResult[content.Table()][content.UniqueKeyColumnValueStr()] = content.CreateCopy()
		return nil
	}
	cachedContent = cacheManager.dbOperationResult[content.Table()][content.UniqueKeyColumnValueStr()]
	if cachedContent.GetCacheStatus() == Delete || cachedContent.GetCacheStatus() == None {
		return nil
	}
	content.SetCacheStatus(cachedContent.GetCacheStatus())

	return content.Update(cachedContent)
}

func (cdb *CacheDB) FindAndSetByConditions(ctx context.Context, contents SelectCacheContents, conditions []*ConditionValue) error {
	// 条件の順番の設定
	sort.Slice(conditions, func(i, j int) bool {
		return conditions[i].TargetColumn < conditions[j].TargetColumn
	})
	sort.Slice(conditions, func(i, j int) bool {
		return conditions[i].ConditionType < conditions[j].ConditionType
	})
	for i, condition := range conditions {
		condition.setOrder(i)
	}

	cacheManager, err := extractDBOperationCacheManager(ctx)
	if err != nil {
		return err
	}

	cacheManager.mutex.Lock()
	defer cacheManager.mutex.Unlock()

	if _, ok := cacheManager.selectQueryCondition[contents.Table()]; !ok {
		cacheManager.selectQueryCondition[contents.Table()] = make(map[string]ConditionMaps)
	}

	shouldFetchFromDB := false
	conditionMap := toConditionMap(conditions)

	key, query, args := conditionMap.CreateQuery(contents.Table())
	if _, ok := cacheManager.selectQueryCondition[contents.Table()][key]; !ok {
		// 一度も引数の条件式でSelect文を実行したことがない
		shouldFetchFromDB = true
		cacheManager.selectQueryCondition[contents.Table()][key] = ConditionMaps{conditionMap}
	}

	cachedConditionMaps, _ := cacheManager.selectQueryCondition[contents.Table()][key]
	if !cachedConditionMaps.IsOnceCalled(conditions) {
		// 引数の条件式でSelect文を実行したことはあるが条件値が異なるか漏れがある
		shouldFetchFromDB = true
		if !cachedConditionMaps.UnionConditionValues(conditions) {
			cachedConditionMaps = append(cachedConditionMaps, conditionMap)
		}
	}

	// データベースからデータを取得して、そのデータをキャッシュする
	if shouldFetchFromDB {
		rows, err := cdb.db.QueryContext(ctx, query, args...)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			if err = contents.BindAndAddContent(rows); err != nil {
				return err
			}
		}
		if err = rows.Err(); err != nil {
			return err
		}

		contents.ForEach(func(content CacheContent) {
			if _, ok := cacheManager.dbOperationResult[content.Table()]; !ok {
				cacheManager.dbOperationResult[content.Table()] = make(map[string]CacheContent)
			}

			if _, ok := cacheManager.dbOperationResult[content.Table()][content.UniqueKeyColumnValueStr()]; !ok {
				content.SetCacheStatus(Select)
				cacheManager.dbOperationResult[content.Table()][content.UniqueKeyColumnValueStr()] = content.CreateCopy()
			}
		})
	}

	// 条件に合致するキャッシュを引数のデータに追加する
	contents.Reset()
	for _, cachedContent := range cacheManager.dbOperationResult[contents.Table()] {
		if cachedContent.GetCacheStatus() == Delete || cachedContent.GetCacheStatus() == None {
			continue
		}
		for _, condition := range conditions {
			switch condition.ConditionType {
			case Eq:
				if !cachedContent.IsSame(condition.TargetColumn, condition.Values[0]) {
					continue
				}
			case In:
				if !cachedContent.IsInclude(condition.TargetColumn, condition.Values) {
					continue
				}
			}
		}
		contents.Add(cachedContent.CreateCopy())
	}
	return nil
}

func (cdb *CacheDB) extractCacheContent(ctx context.Context, content CacheContent) (CacheContent, error) {
	cacheManager, err := extractDBOperationCacheManager(ctx)
	if err != nil {
		return nil, err
	}

	_, ok := cacheManager.dbOperationResult[content.Table()]
	if !ok {
		cacheManager.dbOperationResult[content.Table()] = make(map[string]CacheContent)
	}
	cachedContent, ok := cacheManager.dbOperationResult[content.Table()][content.UniqueKeyColumnValueStr()]
	if !ok {
		return nil, nil
	}
	return cachedContent, nil
}
