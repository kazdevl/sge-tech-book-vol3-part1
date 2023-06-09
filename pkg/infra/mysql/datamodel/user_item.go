// Code generated by SQLBoiler 4.14.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package datamodel

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// InsertAll inserts all rows with the specified column values, using an executor.
func (o UserItemSlice) InsertAll(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	ln := int64(len(o))
	if ln == 0 {
		return nil
	}
	var sql string
	vals := []interface{}{}
	for i, row := range o {

		if err := row.doBeforeInsertHooks(ctx, exec); err != nil {
			return err
		}

		nzDefaults := queries.NonZeroDefaultSet(userItemColumnsWithDefault, row)
		wl, _ := columns.InsertColumnSet(
			userItemAllColumns,
			userItemColumnsWithDefault,
			userItemColumnsWithoutDefault,
			nzDefaults,
		)
		if i == 0 {
			sql = "INSERT INTO `user_item` " + "(`" + strings.Join(wl, "`,`") + "`)" + " VALUES "
		}
		sql += strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), len(vals)+1, len(wl))
		if i != len(o)-1 {
			sql += ","
		}
		valMapping, err := queries.BindMapping(userItemType, userItemMapping, wl)
		if err != nil {
			return err
		}
		value := reflect.Indirect(reflect.ValueOf(row))
		vals = append(vals, queries.ValuesFromMapping(value, valMapping)...)
	}
	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, vals...)
	}

	_, err := exec.ExecContext(ctx, sql, vals...)
	if err != nil {
		return errors.Wrap(err, "datamodel: unable to insert into user_item")
	}

	return nil
}

// UserItem is an object representing the database table.
type UserItem struct { // ユーザID
	UserID int64 `db:"user_id" boil:"user_id" json:"user_id" toml:"user_id" yaml:"user_id"`
	// アイテムID
	ItemID int64 `db:"item_id" boil:"item_id" json:"item_id" toml:"item_id" yaml:"item_id"`
	// 数
	Count int64 `db:"count" boil:"count" json:"count" toml:"count" yaml:"count"`

	R *userItemR `db:"-" boil:"-" json:"-" toml:"-" yaml:"-"`
	L userItemL  `db:"-" boil:"-" json:"-" toml:"-" yaml:"-"`
}

var UserItemColumns = struct {
	UserID string
	ItemID string
	Count  string
}{
	UserID: "user_id",
	ItemID: "item_id",
	Count:  "count",
}

var UserItemTableColumns = struct {
	UserID string
	ItemID string
	Count  string
}{
	UserID: "user_item.user_id",
	ItemID: "user_item.item_id",
	Count:  "user_item.count",
}

// Generated where

var UserItemWhere = struct {
	UserID whereHelperint64
	ItemID whereHelperint64
	Count  whereHelperint64
}{
	UserID: whereHelperint64{field: "`user_item`.`user_id`"},
	ItemID: whereHelperint64{field: "`user_item`.`item_id`"},
	Count:  whereHelperint64{field: "`user_item`.`count`"},
}

// UserItemRels is where relationship names are stored.
var UserItemRels = struct {
}{}

// userItemR is where relationships are stored.
type userItemR struct {
}

// NewStruct creates a new relationship struct
func (*userItemR) NewStruct() *userItemR {
	return &userItemR{}
}

// userItemL is where Load methods for each relationship are stored.
type userItemL struct{}

var (
	userItemAllColumns            = []string{"user_id", "item_id", "count"}
	userItemColumnsWithoutDefault = []string{"user_id", "item_id", "count"}
	userItemColumnsWithDefault    = []string{}
	userItemPrimaryKeyColumns     = []string{"user_id", "item_id"}
	userItemGeneratedColumns      = []string{}
)

type (
	// UserItemSlice is an alias for a slice of pointers to UserItem.
	// This should almost always be used instead of []UserItem.
	UserItemSlice []*UserItem
	// UserItemHook is the signature for custom UserItem hook methods
	UserItemHook func(context.Context, boil.ContextExecutor, *UserItem) error

	userItemQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	userItemType                 = reflect.TypeOf(&UserItem{})
	userItemMapping              = queries.MakeStructMapping(userItemType)
	userItemPrimaryKeyMapping, _ = queries.BindMapping(userItemType, userItemMapping, userItemPrimaryKeyColumns)
	userItemInsertCacheMut       sync.RWMutex
	userItemInsertCache          = make(map[string]insertCache)
	userItemUpdateCacheMut       sync.RWMutex
	userItemUpdateCache          = make(map[string]updateCache)
	userItemUpsertCacheMut       sync.RWMutex
	userItemUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var userItemAfterSelectHooks []UserItemHook

var userItemBeforeInsertHooks []UserItemHook
var userItemAfterInsertHooks []UserItemHook

var userItemBeforeUpdateHooks []UserItemHook
var userItemAfterUpdateHooks []UserItemHook

var userItemBeforeDeleteHooks []UserItemHook
var userItemAfterDeleteHooks []UserItemHook

var userItemBeforeUpsertHooks []UserItemHook
var userItemAfterUpsertHooks []UserItemHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *UserItem) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range userItemAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *UserItem) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range userItemBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *UserItem) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range userItemAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *UserItem) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range userItemBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *UserItem) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range userItemAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *UserItem) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range userItemBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *UserItem) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range userItemAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *UserItem) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range userItemBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *UserItem) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range userItemAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddUserItemHook registers your hook function for all future operations.
func AddUserItemHook(hookPoint boil.HookPoint, userItemHook UserItemHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		userItemAfterSelectHooks = append(userItemAfterSelectHooks, userItemHook)
	case boil.BeforeInsertHook:
		userItemBeforeInsertHooks = append(userItemBeforeInsertHooks, userItemHook)
	case boil.AfterInsertHook:
		userItemAfterInsertHooks = append(userItemAfterInsertHooks, userItemHook)
	case boil.BeforeUpdateHook:
		userItemBeforeUpdateHooks = append(userItemBeforeUpdateHooks, userItemHook)
	case boil.AfterUpdateHook:
		userItemAfterUpdateHooks = append(userItemAfterUpdateHooks, userItemHook)
	case boil.BeforeDeleteHook:
		userItemBeforeDeleteHooks = append(userItemBeforeDeleteHooks, userItemHook)
	case boil.AfterDeleteHook:
		userItemAfterDeleteHooks = append(userItemAfterDeleteHooks, userItemHook)
	case boil.BeforeUpsertHook:
		userItemBeforeUpsertHooks = append(userItemBeforeUpsertHooks, userItemHook)
	case boil.AfterUpsertHook:
		userItemAfterUpsertHooks = append(userItemAfterUpsertHooks, userItemHook)
	}
}

// One returns a single userItem record from the query.
func (q userItemQuery) One(ctx context.Context, exec boil.ContextExecutor) (*UserItem, error) {
	o := &UserItem{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "datamodel: failed to execute a one query for user_item")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all UserItem records from the query.
func (q userItemQuery) All(ctx context.Context, exec boil.ContextExecutor) (UserItemSlice, error) {
	var o []*UserItem

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "datamodel: failed to assign all query results to UserItem slice")
	}

	if len(userItemAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all UserItem records in the query.
func (q userItemQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "datamodel: failed to count user_item rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q userItemQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "datamodel: failed to check if user_item exists")
	}

	return count > 0, nil
}

// UserItems retrieves all the records using an executor.
func UserItems(mods ...qm.QueryMod) userItemQuery {
	mods = append(mods, qm.From("`user_item`"))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"`user_item`.*"})
	}

	return userItemQuery{q}
}

// FindUserItem retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindUserItem(ctx context.Context, exec boil.ContextExecutor, userID int64, itemID int64, selectCols ...string) (*UserItem, error) {
	userItemObj := &UserItem{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from `user_item` where `user_id`=? AND `item_id`=?", sel,
	)

	q := queries.Raw(query, userID, itemID)

	err := q.Bind(ctx, exec, userItemObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "datamodel: unable to select from user_item")
	}

	if err = userItemObj.doAfterSelectHooks(ctx, exec); err != nil {
		return userItemObj, err
	}

	return userItemObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *UserItem) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("datamodel: no user_item provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(userItemColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	userItemInsertCacheMut.RLock()
	cache, cached := userItemInsertCache[key]
	userItemInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			userItemAllColumns,
			userItemColumnsWithDefault,
			userItemColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(userItemType, userItemMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(userItemType, userItemMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO `user_item` (`%s`) %%sVALUES (%s)%%s", strings.Join(wl, "`,`"), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO `user_item` () VALUES ()%s%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			cache.retQuery = fmt.Sprintf("SELECT `%s` FROM `user_item` WHERE %s", strings.Join(returnColumns, "`,`"), strmangle.WhereClause("`", "`", 0, userItemPrimaryKeyColumns))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}
	_, err = exec.ExecContext(ctx, cache.query, vals...)

	if err != nil {
		return errors.Wrap(err, "datamodel: unable to insert into user_item")
	}

	var identifierCols []interface{}

	if len(cache.retMapping) == 0 {
		goto CacheNoHooks
	}

	identifierCols = []interface{}{
		o.UserID,
		o.ItemID,
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.retQuery)
		fmt.Fprintln(writer, identifierCols...)
	}
	err = exec.QueryRowContext(ctx, cache.retQuery, identifierCols...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	if err != nil {
		return errors.Wrap(err, "datamodel: unable to populate default values for user_item")
	}

CacheNoHooks:
	if !cached {
		userItemInsertCacheMut.Lock()
		userItemInsertCache[key] = cache
		userItemInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the UserItem.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *UserItem) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	userItemUpdateCacheMut.RLock()
	cache, cached := userItemUpdateCache[key]
	userItemUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			userItemAllColumns,
			userItemPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("datamodel: unable to update user_item, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE `user_item` SET %s WHERE %s",
			strmangle.SetParamNames("`", "`", 0, wl),
			strmangle.WhereClause("`", "`", 0, userItemPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(userItemType, userItemMapping, append(wl, userItemPrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, values)
	}
	var result sql.Result
	result, err = exec.ExecContext(ctx, cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "datamodel: unable to update user_item row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "datamodel: failed to get rows affected by update for user_item")
	}

	if !cached {
		userItemUpdateCacheMut.Lock()
		userItemUpdateCache[key] = cache
		userItemUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q userItemQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "datamodel: unable to update all for user_item")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "datamodel: unable to retrieve rows affected for user_item")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o UserItemSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("datamodel: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), userItemPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE `user_item` SET %s WHERE %s",
		strmangle.SetParamNames("`", "`", 0, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, userItemPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "datamodel: unable to update all in userItem slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "datamodel: unable to retrieve rows affected all in update all userItem")
	}
	return rowsAff, nil
}

// Delete deletes a single UserItem record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *UserItem) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("datamodel: no UserItem provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), userItemPrimaryKeyMapping)
	sql := "DELETE FROM `user_item` WHERE `user_id`=? AND `item_id`=?"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "datamodel: unable to delete from user_item")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "datamodel: failed to get rows affected by delete for user_item")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q userItemQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("datamodel: no userItemQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "datamodel: unable to delete all from user_item")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "datamodel: failed to get rows affected by deleteall for user_item")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o UserItemSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(userItemBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), userItemPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM `user_item` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, userItemPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "datamodel: unable to delete all from userItem slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "datamodel: failed to get rows affected by deleteall for user_item")
	}

	if len(userItemAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *UserItem) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindUserItem(ctx, exec, o.UserID, o.ItemID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *UserItemSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := UserItemSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), userItemPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT `user_item`.* FROM `user_item` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, userItemPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "datamodel: unable to reload all in UserItemSlice")
	}

	*o = slice

	return nil
}

// UserItemExists checks if the UserItem row exists.
func UserItemExists(ctx context.Context, exec boil.ContextExecutor, userID int64, itemID int64) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from `user_item` where `user_id`=? AND `item_id`=? limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, userID, itemID)
	}
	row := exec.QueryRowContext(ctx, sql, userID, itemID)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "datamodel: unable to check if user_item exists")
	}

	return exists, nil
}

// Exists checks if the UserItem row exists.
func (o *UserItem) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return UserItemExists(ctx, exec, o.UserID, o.ItemID)
}

var mySQLUserItemUniqueColumns = []string{}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *UserItem) Upsert(ctx context.Context, exec boil.ContextExecutor, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("datamodel: no user_item provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(userItemColumnsWithDefault, o)
	nzUniques := queries.NonZeroDefaultSet(mySQLUserItemUniqueColumns, o)

	if len(nzUniques) == 0 {
		return errors.New("cannot upsert with a table that cannot conflict on a unique column")
	}

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzUniques {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	userItemUpsertCacheMut.RLock()
	cache, cached := userItemUpsertCache[key]
	userItemUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			userItemAllColumns,
			userItemColumnsWithDefault,
			userItemColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			userItemAllColumns,
			userItemPrimaryKeyColumns,
		)

		if !updateColumns.IsNone() && len(update) == 0 {
			return errors.New("datamodel: unable to upsert user_item, could not build update column list")
		}

		ret = strmangle.SetComplement(ret, nzUniques)
		cache.query = buildUpsertQueryMySQL(dialect, "`user_item`", update, insert)
		cache.retQuery = fmt.Sprintf(
			"SELECT %s FROM `user_item` WHERE %s",
			strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, ret), ","),
			strmangle.WhereClause("`", "`", 0, nzUniques),
		)

		cache.valueMapping, err = queries.BindMapping(userItemType, userItemMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(userItemType, userItemMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}
	_, err = exec.ExecContext(ctx, cache.query, vals...)

	if err != nil {
		return errors.Wrap(err, "datamodel: unable to upsert for user_item")
	}

	var uniqueMap []uint64
	var nzUniqueCols []interface{}

	if len(cache.retMapping) == 0 {
		goto CacheNoHooks
	}

	uniqueMap, err = queries.BindMapping(userItemType, userItemMapping, nzUniques)
	if err != nil {
		return errors.Wrap(err, "datamodel: unable to retrieve unique values for user_item")
	}
	nzUniqueCols = queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), uniqueMap)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.retQuery)
		fmt.Fprintln(writer, nzUniqueCols...)
	}
	err = exec.QueryRowContext(ctx, cache.retQuery, nzUniqueCols...).Scan(returns...)
	if err != nil {
		return errors.Wrap(err, "datamodel: unable to populate default values for user_item")
	}

CacheNoHooks:
	if !cached {
		userItemUpsertCacheMut.Lock()
		userItemUpsertCache[key] = cache
		userItemUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}
