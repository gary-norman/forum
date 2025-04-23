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

	for i := 0; i < typ.NumField(); i++ {
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
