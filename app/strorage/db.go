package storage

import (
	"SnipAndNeat/app/config"
	"context"
	"database/sql"
	"embed"
	"net/http"
	"os"
	"time"

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
				conn.AuthEnabled()
				conn.AuthUserAdd(db.cfg.DBConfig.User, db.cfg.DBConfig.Password, true)
				return conn.Ping(ctx)
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

	if err := db.Ping(ctx); err != nil {
		return errors.Wrap(err, "cannot ping database")
	}

	return db.migratingUp(ctx)
}

func (db *DB) Stop(context.Context) error {
	if db == nil || db.db == nil {
		return nil
	}
	return db.db.Close()
}

func (db *DB) Ping(context.Context) error {
	if db == nil || db.db == nil {
		return sql.ErrConnDone
	}
	return db.db.Ping()
}

func (db *DB) migratingUp(ctx context.Context) error {
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

// type OperFunc func(c context.Context, in any) (any, error)

// func (db *DB) WithTransaction(f OperFunc) OperFunc {
// 	return func(ctx context.Context, in any) (any, error) {
// 		tx := db.gorm.Begin()
// 		ctx = context.WithValue(ctx, "tx", tx)
// 		v, err := f(ctx, in)
// 		if err != nil {
// 			tx.Rollback()
// 			ctx = context.WithValue(ctx, "tx", nil)
// 			return nil, err
// 		} else {
// 			tx.Commit()
// 			ctx = context.WithValue(ctx, "tx", nil)
// 		}

// 		return v, nil
// 	}
// }
