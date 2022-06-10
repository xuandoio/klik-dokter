package model

import (
	"fmt"
	"log"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/xuandoio/klik-dokter/internal/config"
	"golang.org/x/crypto/bcrypt"
)

// Transaction is an interface that models the standard transaction in `pg.DB`
// To ensure `TxFn` func cannot commit or rollback a transaction (which is
// handled by `WithTransaction`), those methods are not included here.
type Transaction interface {
	Exec(query interface{}, params ...interface{}) (pg.Result, error)
	Prepare(q string) (*pg.Stmt, error)
	Query(model interface{}, query interface{}, params ...interface{}) (pg.Result, error)
	QueryOne(model interface{}, query interface{}, params ...interface{}) (pg.Result, error)
}

// A TxFn is a function that will be called with an initialized `Transaction` object
// that can be used for executing statements and queries against a database.
type TxFn func(db orm.DB) error

// WithTransaction creates a new transaction and handles rollback/commit based on the
// error object returned by the `TxFn`
func WithTransaction(db *pg.DB, fn TxFn) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Close()

	defer func() {
		if p := recover(); p != nil {
			// a panic occurred, rollback and re panic
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			// something went wrong, rollback
			_ = tx.Rollback()
		} else {
			// all good, commit
			err = tx.Commit()
		}
	}()

	err = fn(tx)
	return err
}

// NewConnection /**
func NewConnection(c *config.Config) (*pg.DB, error) {
	sqlConnectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.Database.Username,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Name,
	)
	opt, err := pg.ParseURL(sqlConnectionString)
	if err != nil {
		log.Fatalf("Error connecting to SQL database:%v", err)
	}

	return pg.Connect(opt), nil
}

// HashPassword /**
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

// CheckPasswordHash /**
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
