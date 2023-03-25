package middleware

import (
	"errors"
	"game-server-example/pkg/handler/acontext"
	"game-server-example/pkg/handler/schema"
	"game-server-example/pkg/infra/mysql"
	"game-server-example/pkg/infra/mysql/cachedb"
	"game-server-example/pkg/infra/mysql/cachemodel"
	"github.com/labstack/echo/v4"
	"net/http"
)

type TransactionMiddleware struct {
	db                    *mysql.ApplicationDB
	cacheDB               *cachedb.CacheDB
	enableCacheRepository bool
}

func NewTransactionMiddleware(db *mysql.ApplicationDB, cacheDB *cachedb.CacheDB, enableCacheRepository bool) *TransactionMiddleware {
	return &TransactionMiddleware{db: db, cacheDB: cacheDB, enableCacheRepository: enableCacheRepository}
}

func (m *TransactionMiddleware) Intercept(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		var commonResp *schema.CommonResponse
		ctx := c.Request().Context()
		ctx = acontext.WithAPIResponse(ctx)
		ctx = cachedb.WithCacheContext(ctx)
		c.SetRequest(c.Request().WithContext(ctx))

		if m.enableCacheRepository {
			defer func() {
				if err != nil {
					m.cacheDB.Rollback(ctx)
				}
			}()

			if err = next(c); err != nil {
				return err
			}

			if err = m.cacheDB.Begin(); err != nil {
				return err
			}

			var cacheSyncContents map[string]map[cachedb.CacheStatus][]cachedb.CacheContent
			cacheSyncContents, err = m.cacheDB.SyncedToDB(ctx)
			if err != nil {
				return err
			}

			if err = m.cacheDB.Commit(ctx); err != nil {
				return err
			}

			if c.Request().URL.Path != "/user/register" {
				commonResp, err = m.toSchemaCommonResponse(cacheSyncContents)
				if err != nil {
					return err
				}
			}
		} else {
			defer func() {
				if err != nil {
					m.db.RollBack()
				}
			}()

			m.db.Begin()
			if err = next(c); err != nil {
				return err
			}

			if err = m.db.Commit(); err != nil {
				return err
			}
		}

		apiResponse, err := acontext.ExtractAPIResponse(ctx)
		if err != nil {
			return err
		}

		if c.Request().URL.Path == "/user/data" || commonResp == nil {
			return c.JSON(http.StatusOK, apiResponse)
		}

		apiResponse.Common = commonResp
		return c.JSON(http.StatusOK, apiResponse)
	}
}

func (m *TransactionMiddleware) toSchemaCommonResponse(cacheSyncContents map[string]map[cachedb.CacheStatus][]cachedb.CacheContent) (*schema.CommonResponse, error) {
	resp := &schema.CommonResponse{
		Delete: &struct {
			UserCoin    *[]schema.UserCoin    `json:"user_coin,omitempty"`
			UserItem    *[]schema.UserItem    `json:"user_item,omitempty"`
			UserMonster *[]schema.UserMonster `json:"user_monster,omitempty"`
		}{},
		Update: &struct {
			UserCoin    *[]schema.UserCoin    `json:"user_coin,omitempty"`
			UserItem    *[]schema.UserItem    `json:"user_item,omitempty"`
			UserMonster *[]schema.UserMonster `json:"user_monster,omitempty"`
		}{},
	}

	for table, contentsMap := range cacheSyncContents {
		for cacheStatus, contents := range contentsMap {
			switch cacheStatus {
			case cachedb.Insert:
				if err := m.setUpdateContents(resp, table, contents); err != nil {
					return nil, err
				}
			case cachedb.Update:
				if err := m.setUpdateContents(resp, table, contents); err != nil {
					return nil, err
				}
			case cachedb.Delete:
				if err := m.setDeleteContents(resp, table, contents); err != nil {
					return nil, err
				}
			}
		}
	}

	return resp, nil
}

func (m *TransactionMiddleware) setUpdateContents(dest *schema.CommonResponse, table string, contents []cachedb.CacheContent) error {
	switch table {
	case "user_coin":
		resp, err := m.toCoinTableResponse(contents)
		if err != nil {
			return err
		}
		dest.Update.UserCoin = resp
	case "user_item":
		resp, err := m.toItemTableResponse(contents)
		if err != nil {
			return err
		}
		dest.Update.UserItem = resp
	case "user_monster":
		resp, err := m.toMonsterTableResponse(contents)
		if err != nil {
			return err
		}
		dest.Update.UserMonster = resp
	}
	return nil
}

func (m *TransactionMiddleware) setDeleteContents(dest *schema.CommonResponse, table string, contents []cachedb.CacheContent) error {
	switch table {
	case "user_coin":
		resp, err := m.toCoinTableResponse(contents)
		if err != nil {
			return err
		}
		dest.Delete.UserCoin = resp
	case "user_item":
		resp, err := m.toItemTableResponse(contents)
		if err != nil {
			return err
		}
		dest.Delete.UserItem = resp
	case "user_monster":
		resp, err := m.toMonsterTableResponse(contents)
		if err != nil {
			return err
		}
		dest.Delete.UserMonster = resp
	}
	return nil
}

func (m *TransactionMiddleware) toCoinTableResponse(contents []cachedb.CacheContent) (*[]schema.UserCoin, error) {
	res := make([]schema.UserCoin, 0, len(contents))
	for _, content := range contents {
		cacheM, ok := content.(*cachemodel.UserCoinCacheModel)
		if !ok {
			return nil, errors.New("想定していない型です")
		}
		res = append(res, schema.UserCoin{
			Currency: &cacheM.Num,
			UserId:   &cacheM.UserID,
		})
	}
	return &res, nil
}

func (m *TransactionMiddleware) toItemTableResponse(contents []cachedb.CacheContent) (*[]schema.UserItem, error) {
	res := make([]schema.UserItem, 0, len(contents))
	for _, content := range contents {
		cacheM, ok := content.(*cachemodel.UserItemCacheModel)
		if !ok {
			return nil, errors.New("想定していない型です")
		}
		res = append(res, schema.UserItem{
			UserId: &cacheM.UserID,
			ItemId: &cacheM.ItemID,
			Count:  &cacheM.Count,
		})
	}
	return &res, nil
}

func (m *TransactionMiddleware) toMonsterTableResponse(contents []cachedb.CacheContent) (*[]schema.UserMonster, error) {
	res := make([]schema.UserMonster, 0, len(contents))
	for _, content := range contents {
		cacheM, ok := content.(*cachemodel.UserMonsterCacheModel)
		if !ok {
			return nil, errors.New("想定していない型です")
		}
		res = append(res, schema.UserMonster{
			UserId:    &cacheM.UserID,
			MonsterId: &cacheM.MonsterID,
			Exp:       &cacheM.Exp,
		})
	}
	return &res, nil
}
