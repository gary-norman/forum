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

func (dao *DAO[T]) All(ctx context.Context) ([]T, error) {
	var model T
	query := fmt.Sprintf("SELECT * FROM %s", model.TableName())

	rows, err := dao.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, _ := rows.Columns()
	var results []T

	for rows.Next() {
		var model T
		dest := make([]any, len(columns))
		val := reflect.ValueOf(&model).Elem()

		for i, col := range columns {
			field := val.FieldByNameFunc(func(name string) bool {
				field, _ := val.Type().FieldByName(name)
				return field.Tag.Get("db") == col
			})
			dest[i] = field.Addr().Interface()
		}

		if err := rows.Scan(dest...); err != nil {
			return nil, err
		}
		results = append(results, model)
	}

	return results, nil
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
