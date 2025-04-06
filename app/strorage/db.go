package storage

import (
	"SnipAndNeat/app/config"
	"bufio"
	"context"
	"database/sql"
	"embed"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	models "SnipAndNeat/generated"

	migrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
	"github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"gorm.io/driver/sqlite" // Sqlite driver based on CGO
	// "github.com/glebarez/sqlite" // Pure go SQLite driver, checkout https://github.com/glebarez/sqlite for details
	"gorm.io/gorm"
)

const (
	msgMigrateUp = "migrate up"
	msgVersion   = "version"
)

//go:embed migrations
var migrations embed.FS

// DB data access layer PostgreSQL
type DB struct {
	log zerolog.Logger
	gdb *gorm.DB
	db  *sql.DB
	cfg *config.Config
}

// New конструктор DB
func New(cfg *config.Config) *DB {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	return &DB{
		log: zerolog.New(output),
		cfg: cfg,
	}
}

type MyLog func(format string, v ...any)

func regexFn(re, s string) (bool, error) {
	return regexp.MatchString(re, s)
}

func (db *DB) Start(ctx context.Context) error {
	newLogger := zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
			MessageKey:  "database",
			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalColorLevelEncoder,
		}),
		zapcore.Lock(os.Stdout),
		zapcore.DebugLevel))

	okCh, errCh := make(chan struct{}), make(chan error) // startConsistently
	var err error
	newLogger.Info("start db")

	go func() {
		sql.Register("sqlite3_with_extensions", &sqlite3.SQLiteDriver{
			ConnectHook: func(conn *sqlite3.SQLiteConn) error {
				if err := conn.RegisterFunc("regex", regexFn, true); err != nil {
					return fmt.Errorf("Error registering function regex: %s", err.Error())
				}
				return nil
			},
		})

		db.gdb, err = gorm.Open(sqlite.Open("viento.db"), &gorm.Config{})
		if err != nil {
			errCh <- errors.Wrap(err, "cannot open connection")
			return
		}

		db.db, err = db.gdb.DB()
		if err != nil {
			errCh <- errors.Wrap(err, "cannot get db connection")
			return
		}

		okCh <- struct{}{}
	}()

	select {
	case err = <-errCh:
		return err
	case <-time.After(db.cfg.StartTimeout):
		return nil
	case <-okCh:
	}

	db.db.SetMaxIdleConns(db.cfg.DBConfig.MaxIdleConns)
	db.db.SetMaxOpenConns(db.cfg.DBConfig.MaxOpenConns)
	db.db.SetConnMaxLifetime(db.cfg.DBConfig.ConnMaxLifetime)

	if err := db.Ping(); err != nil {
		return errors.Wrap(err, "cannot ping database")
	}

	if err := db.migratingUp(); err != nil {
		return errors.Wrap(err, "cannot migrate up")
	}

	if err := db.migrateVientoItems(); err != nil {
		return errors.Wrap(err, "cannot migrate viento items")
	}

	return nil
}

func (db *DB) Stop(ctx context.Context) error {
	if db == nil || db.db == nil {
		return nil
	}
	return db.db.Close()
}

func (db *DB) Ping() error {
	if db == nil || db.db == nil {
		return sql.ErrConnDone
	}
	return db.db.Ping()
}

func (db *DB) migratingUp() error {
	source, err := httpfs.New(http.FS(migrations), "migrations")
	if err != nil {
		return errors.Wrap(err, "failed to create fs source")
	}
	m, err := migrate.NewWithSourceInstance("httpfs", source, db.cfg.StringDB())
	if err != nil {
		return errors.Wrap(err, "migration init")
	}

	var mVersion uint
	mVersion, _, err = m.Version()
	db.log.Info().Err(err).Uint(msgVersion, mVersion).Msg("current migration version")
	err = m.Up()
	switch err {
	case nil:
		mVersion, _, err = m.Version()
		db.log.Info().Err(err).Uint(msgVersion, mVersion).Msg("new migration version")
		return nil
	case migrate.ErrNoChange:
		db.log.Info().Err(err).Msg(msgMigrateUp)
		return nil
	default:
		db.log.Info().Err(err).Msg(msgMigrateUp)
		return errors.Wrap(err, msgMigrateUp)
	}
}

func (db *DB) migrateVientoItems() error {
	f, err := os.Open("../../temp_scropt/prices.csv")
	if err != nil {
		return err
	}
	defer f.Close()

	db.log.Info().Msg("migrate viento items")
	scanner := bufio.NewScanner(f)
	var items []models.VientoItem
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")
		if len(parts) != 3 {
			db.log.Error().Msgf("invalid line format: %s", line)
			continue
		}
		name, priceStr, eanCodeStr := parts[0], parts[1], parts[2]

		eanCode, err := strconv.ParseInt(eanCodeStr, 10, 64)
		if err != nil {
			db.log.Error().Err(err).Msgf("invalid ean code: %s", eanCodeStr)
			continue
		}

		price, err := strconv.ParseFloat(priceStr, 32)
		if err != nil {
			db.log.Error().Err(err).Msgf("invalid price: %s", priceStr)
			continue
		}

		item := models.VientoItem{
			Name:  name,
			Price: price,
			Ean:   eanCode,
		}
		items = append(items, item)
	}

	if err := scanner.Err(); err != nil {
		db.log.Error().Err(err).Msg("failed to read file")
		return err
	}
	tx := db.gdb.Table("viento_items").Create(items)
	if tx.Error != nil {
		db.log.Error().Err(tx.Error).Msg("failed to migrate viento items")
		return tx.Error
	}

	db.log.Info().Msgf("migrate viento items %d finished", tx.RowsAffected)

	return nil
}

// func addPeriodQuery(q *gorm.DB, tabler schema.Tabler, field string, from, to time.Time) *gorm.DB {
// 	if tabler != nil {
// 		field = fmt.Sprintf("%s.%s", tabler.TableName(), field)
// 	}

// 	if !from.IsZero() {
// 		q = q.Where(fmt.Sprintf("%s >= ?", field), from.UTC())
// 	}
// 	if !to.IsZero() {
// 		q = q.Where(fmt.Sprintf("%s <= ?", field), to.UTC())
// 	}
// 	return q
// }

type OperFunc func(tx *gorm.DB, in any) (any, error)

func (db *DB) WithTransaction(ctx context.Context, f OperFunc, in any) (any, int64, error) {
	tx := db.gdb.Begin()
	v, err := f(tx, in)
	var res int64
	if err != nil {
		tx.Rollback()
		return nil, 0, err
	} else {
		res = tx.Commit().RowsAffected
	}

	return v, res, nil

}
