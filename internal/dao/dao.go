package dao

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/gary-norman/forum/internal/models"
)

type DAO[T models.DBModel] struct {
	DB *sql.DB
}

func (dao *DAO[T]) All(ctx context.Context) ([]T, error) {
	var model T
	query := fmt.Sprintf("SELECT * FROM %s", model.TableName())

	rows, err := dao.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	return scanRowsToStructs[T](rows)
}

func (dao *DAO[T]) GetByID(ctx context.Context, id int64) (*T, error) {
	var model T
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = ?", model.TableName())

	rows, err := dao.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results, err := scanRowsToStructs[T](rows)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, sql.ErrNoRows
	}
	return &results[0], nil
}

func (dao *DAO[T]) Insert(ctx context.Context, model T) error {
	val := reflect.ValueOf(model)
	typ := reflect.TypeOf(model)

	var cols []string
	var placeholders []string
	var args []any

	for i := range typ.NumField() {
		field := typ.Field(i)
		tag := field.Tag.Get("db")
		if tag == "" || tag == "-" {
			continue
		}
		cols = append(cols, tag)
		placeholders = append(placeholders, "?")
		args = append(args, val.Field(i).Interface())
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		model.TableName(),
		strings.Join(cols, ", "),
		strings.Join(placeholders, ", "))

	_, err := dao.DB.ExecContext(ctx, query, args...)
	return err
}

func (dao *DAO[T]) Delete(ctx context.Context, id int64) (*T, error) {
	var model T
	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", model.TableName())

	rows, err := dao.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result, err := scanRowsToStructs[T](rows)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("no rows returned after delete")
	}
	return &result[0], nil
}

func scanRowsToStructs[T any](rows *sql.Rows) ([]T, error) {
	defer rows.Close()

	var results []T
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var t T
		val := reflect.ValueOf(&t).Elem()
		dest := make([]any, len(columns))

		for i, col := range columns {
			field := val.FieldByNameFunc(func(name string) bool {
				sf, ok := val.Type().FieldByName(name)
				return ok && sf.Tag.Get("db") == col
			})

			if field.IsValid() && field.CanSet() {
				dest[i] = field.Addr().Interface()
			} else {
				var discard any
				dest[i] = &discard
			}
		}

		if err := rows.Scan(dest...); err != nil {
			return nil, err
		}

		results = append(results, t)
	}

	return results, rows.Err()
}
