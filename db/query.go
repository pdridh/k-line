package db

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

var PSQL = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

func GetCount(ctx context.Context, db *sqlx.DB, query squirrel.SelectBuilder) (int, error) {
	sql, args, err := query.ToSql()

	if err != nil {
		return 0, err
	}

	var count int
	if err := db.GetContext(ctx, &count, sql, args...); err != nil {
		return 0, err
	}

	return count, err
}
