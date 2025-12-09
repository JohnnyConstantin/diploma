package db

import (
	"context"
	"database/sql"
	"diploma/pkg/config"
	"diploma/pkg/logger"
	"diploma/pkg/logger/message"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

// SqlConnection объект подключения к БД
type SqlConnection struct {
	PgSql         *pgxpool.Pool
	SqlDB         *sql.DB
	Timeout       time.Duration
	DBName        string
	MigrationsDir string
}

// InitSQL создаёт и инициализирует глобальное подключение к БД Postgres.
func InitSQL(cfg *config.DBConfig) (*SqlConnection, error) {
	var inst *SqlConnection
	var initErr error
	var once sync.Once
	timeout := time.Duration(cfg.DBTimeout) * time.Millisecond

	once.Do(func() {
		pgCfg, err := pgxpool.ParseConfig(cfg.DatabaseUri)
		if err != nil {
			initErr = err
			return
		}

		// Ненастраиваемые параметры. Нехорошо, но пользователю их явно показывать нельзя
		pgCfg.MaxConns = int32(cfg.PoolSize)
		pgCfg.MaxConnLifetime = 30 * time.Minute
		pgCfg.MaxConnIdleTime = 5 * time.Minute
		pgCfg.HealthCheckPeriod = time.Minute

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		dbPool, err := pgxpool.NewWithConfig(ctx, pgCfg)
		if err != nil {
			initErr = fmt.Errorf("pgxpool.NewWithConfig: %w", err)
			return
		}

		if err := dbPool.Ping(ctx); err != nil {
			dbPool.Close()
			initErr = fmt.Errorf("failed ping: %w", err)
			return
		}

		sqlDB := stdlib.OpenDBFromPool(dbPool)

		inst = &SqlConnection{
			PgSql:         dbPool,
			SqlDB:         sqlDB,
			Timeout:       timeout,
			DBName:        cfg.DBName,
			MigrationsDir: cfg.MigrationsDir,
		}
	})

	if initErr != nil || inst == nil {
		return nil, initErr
	}
	return inst, nil
}

// CloseSqlInstance закрывает пул соединений и связанные ресурсы БД.
func (s *SqlConnection) CloseSqlInstance() {
	// Значит объекта соединения с БД нет
	if s == nil {
		return
	}

	// Закрываем адаптер соединения с БД
	if s.SqlDB != nil {
		_ = s.SqlDB.Close()
	}

	// Закрываем пул Pgx
	if s.PgSql != nil {
		logger.Log.Info(&message.LogMessage{Message: "Closing database connection pool..."})
		s.PgSql.Close()
		logger.Log.Info(&message.LogMessage{Message: "Database connection pool closed."})
	}
}
