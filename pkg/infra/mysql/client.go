package mysql

import (
	"fmt"
	"game-server-example/config"
	"github.com/jmoiron/sqlx"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type ApplicationDB struct {
	Db *sqlx.DB
	Tx *sqlx.Tx
}

func RetrieveSqlxDB(db *ApplicationDB) *sqlx.DB {
	return db.Db
}

func NewDB(mysqlConfig *config.Mysql) (*ApplicationDB, func(), error) {
	db, err := sqlx.Open("mysql", mysqlConfig.DNS())
	if err != nil {
		return nil, func() {}, err
	}
	i := 0
	for err = db.Ping(); err != nil; {
		if i > 3 {
			return nil, func() {}, err
		}
		time.Sleep(5 * time.Second)
		err = db.Ping()
		i++
	}
	db.SetMaxIdleConns(mysqlConfig.MaxIdleConns)
	db.SetMaxOpenConns(mysqlConfig.MaxOpenConns)
	db.SetConnMaxLifetime(mysqlConfig.ConnMaxLifetime)

	stopCh := make(chan struct{})
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("panic error: %+v", r)
			}
		}()

		ticker := time.NewTicker(5 * time.Second)
		for {
			select {
			case <-ticker.C:
				s := db.Stats()
				fmt.Println("************")
				fmt.Printf("stats: %+v\n", s)
			case <-stopCh:
				ticker.Stop()
				return
			}
		}
	}()

	return &ApplicationDB{Db: db}, func() {
		db.Close()
		close(stopCh)
	}, nil
}

func (a *ApplicationDB) Begin() {
	a.Tx = a.Db.MustBegin()
}

func (a *ApplicationDB) Commit() error {
	log.Println("commitの実行")
	return a.Tx.Commit()
}

func (a *ApplicationDB) RollBack() error {
	log.Println("rollbackの実行")
	return a.Tx.Rollback()
}
