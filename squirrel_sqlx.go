package main

import (
	"context"

	sq "github.com/Masterminds/squirrel"
)

func (q *Queries) GetAuthorX(ctx context.Context, id int) (Author, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.Select(
		"id",
		"bio",
		"birth_year",
	).
		From("authors").
		Where(sq.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return Author{}, err
	}

	var a Author
	if err := q.dbx.GetContext(ctx, &a, sql, args...); err != nil {
		return Author{}, nil
	}

	return a, nil
}

func (q *Queries) ListAuthorsX(ctx context.Context) ([]Author, error) {
	var authors []Author

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.Select(
		"id",
		"bio",
		"birth_year",
	).
		From("authors").
		OrderBy("id")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	if err := q.dbx.SelectContext(ctx, &authors, sql, args...); err != nil {
		return nil, err
	}

	return authors, nil
}
