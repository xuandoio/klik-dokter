package migration

import (
	"fmt"
	"log"

	"github.com/go-pg/pg/v10"
	"github.com/xuandoio/klik-dokter/internal/app/model"
	"github.com/xuandoio/klik-dokter/internal/config"
)

type Migrator interface {
	Up() (err error)
	Down() (err error)
	GetVersion() string
	SetDB(db *pg.DB)
	GetDB() *pg.DB
}

var Migrators = []Migrator{
	&InitSchema{},
}

// Migration /**
type Migration struct {
	tableName struct{} `pg:"migrations,alias:migration"`
	Id        int      `pg:"id,pk"`
	Migration string   `pg:"migration"`
	Batch     int      `pg:"batch"`
	DB        *pg.DB   `pg:"-"`
}

// SetDB /**
func (m *Migration) SetDB(db *pg.DB) {
	m.DB = db
}

// GetDB /**
func (m *Migration) GetDB() *pg.DB {
	return m.DB
}

// Engine /**
type Engine struct {
	DB *pg.DB
}

// NewEngine /**
func NewEngine(c *config.Config) *Engine {
	const migration = `CREATE TABLE IF NOT EXISTS migrations(
		id SERIAL PRIMARY KEY,
		migration VARCHAR(255) NOT NULL,
		batch INT NOT NULL)`

	db, err := model.NewConnection(c)
	if err != nil {
		log.Fatalln("Error connecting to database...")
	}

	// init the "migration" table
	_, err = db.ExecOne(migration)

	if err != nil {
		log.Fatalf("Error while creating migrations table:%v\n", err)
	}

	return &Engine{
		DB: db,
	}
}

// Migrate /**
func (engine *Engine) Migrate() (err error) {
	lastBatch, err := engine.getLastBatch()

	if err != nil {
		return err
	}

	successMigration := 0

	fmt.Println("Begin migrating...")

	for _, migrator := range Migrators {
		migrator.SetDB(engine.DB)

		version := migrator.GetVersion()

		migration, err := engine.getMigrationByVersion(version)
		if err != nil { // not migrated yet
			migration = Migration{
				Migration: version,
				Batch:     lastBatch + 1,
				DB:        engine.DB,
			}
		} else { // already migrated
			continue
		}

		// "up" the migration changes
		fmt.Println("Migrating ", migrator.GetVersion())
		err = migrator.Up()

		if err != nil {
			fmt.Printf("Error when migrating %s: %v\n", migrator.GetVersion(), err)
			panic(err)
		}

		// insert a new record to migration table
		err = migration.Create()
		if err != nil {
			panic(err)
		}
		successMigration++
		fmt.Println("Migrated ", migrator.GetVersion())
	}

	if successMigration > 0 {
		fmt.Println("Migrate Done.")
	} else {
		fmt.Println("Nothing to migrate.")
	}

	return err
}

// Rollback /**
func (engine *Engine) Rollback() (err error) {
	lastBatch, err := engine.getLastBatch()
	migrations, err := engine.getMigrationByBatch(lastBatch)

	successRollback := 0

	fmt.Println("Begin rolling back...")
	for _, m := range migrations {
		m.SetDB(engine.DB)
		fmt.Println("Rolling back ", m.Migration)
		err = m.Rollback()
		if err != nil {
			return err
		}
		err = m.Delete()
		if err != nil {
			return err
		}
		fmt.Println("Rolled back ", m.Migration)
		successRollback++
	}

	if successRollback > 0 {
		fmt.Println("Rollback Done.")
	} else {
		fmt.Println("Nothing to rollback.")
	}
	return
}

/**
Get the latest batch on migration table.
*/
func (engine *Engine) getLastBatch() (maxBatch int, err error) {
	statement := "SELECT COALESCE(MAX(batch), 0) AS b FROM migrations"
	_, err = engine.DB.QueryOne(pg.Scan(&maxBatch), statement)
	return maxBatch, err
}

// Create /**
func (m *Migration) Create() (err error) {
	statement := "INSERT INTO migrations(migration, batch) VALUES(?, ?) RETURNING id"

	_, err = m.DB.ExecOne(statement, m.Migration, m.Batch)
	return err
}

// Delete /**
func (m *Migration) Delete() (err error) {
	_, err = m.DB.ExecOne("DELETE FROM migrations WHERE id=?", m.Id)
	return err
}

// Rollback /**
func (m *Migration) Rollback() (err error) {
	for _, migrator := range Migrators {
		migrator.SetDB(m.DB)
		if m.Migration == migrator.GetVersion() {
			err = migrator.Down()
		}
	}
	return err
}

/**
Get migrations by batch number
*/
func (engine *Engine) getMigrationByBatch(batch int) (migrations []Migration, err error) {
	_, err = engine.DB.Query(&migrations, "SELECT id, migration, batch FROM migrations WHERE batch=? ORDER BY id DESC", batch)
	return migrations, err
}

/**
Get specific migration by version
*/
func (engine *Engine) getMigrationByVersion(version string) (migration Migration, err error) {
	statement := "SELECT id, migration, batch FROM migrations WHERE migration=?"
	_, err = engine.DB.QueryOne(&migration, statement, version)
	migration.SetDB(engine.DB)
	return migration, err
}

// Reset /**
func (engine *Engine) Reset() (err error) {
	var dropTableStatements []string
	_, err = engine.DB.Query(&dropTableStatements, `SELECT 'drop table if exists ' || tablename || ' cascade;' 
		FROM pg_tables 
		WHERE schemaname = 'public' 
		AND tablename NOT IN (?)`, pg.In([]string{"migrations"}))
	if err != nil {
		return err
	}

	fmt.Println("Begin drop all tables...")
	for _, statement := range dropTableStatements {
		_, err = engine.DB.ExecOne(statement)
	}
	fmt.Println("Drop all tables done.")
	// reset the "migrations" table
	_, err = engine.DB.Exec("TRUNCATE ONLY migrations RESTART IDENTITY;")
	if err != nil {
		return err
	}

	err = engine.Migrate()
	return err
}
