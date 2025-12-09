package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ExecuteDBQueryAll выполняет SELECT-запрос и возвращает все строки в виде [][]any
// с использованием контекста с таймаутом.
func ExecuteDBQueryAll(ctx context.Context, pool *pgxpool.Pool, queryTimeout time.Duration, query string, args ...any) ([][]any, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	rows, err := pool.Query(timeoutCtx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results [][]any
	for rows.Next() {
		vals, err := rows.Values()
		if err != nil {
			return nil, err
		}
		results = append(results, vals)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

// ExecuteDBQuery выполняет SELECT-запрос и возвращает только первую строку.
func ExecuteDBQuery(ctx context.Context, pool *pgxpool.Pool, queryTimeout time.Duration, query string, args ...any) ([]any, error) {
	rows, err := ExecuteDBQueryAll(ctx, pool, queryTimeout, query, args...)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, nil
	}
	return rows[0], nil
}

// ExecuteDBExec Fire and forget запрос к БД через SqlConnection.
func ExecuteDBExec(ctx context.Context, pool *pgxpool.Pool, queryTimeout time.Duration, query string, args ...any) (int64, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	ct, err := pool.Exec(timeoutCtx, query, args...)
	if err != nil {
		return 0, err
	}
	return ct.RowsAffected(), nil
}
