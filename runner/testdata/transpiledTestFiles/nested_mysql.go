package main

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

// Load___datagen_nested_mysql executes a single batch of records using the provided transaction.
func Load___datagen_nested_mysql(records []*__datagen_nested, tx *sql.Tx) error {
	if len(records) == 0 {
		return nil
	}

	ctx := context.Background()

	var b strings.Builder
	columns := []string{
		"`id`",
		"`user`",
	}
	b.WriteString("INSERT INTO nested (")
	b.WriteString(strings.Join(columns, ","))
	b.WriteString(") VALUES ")

	placeholderGroup := "(" + strings.Repeat("?,", 2)
	placeholderGroup = placeholderGroup[:len(placeholderGroup)-1] + ")"
	for i := range records {
		if i > 0 {
			b.WriteString(",")
		}
		b.WriteString(placeholderGroup)
	}
	sqlStmt := b.String()

	var args []interface{}
	for _, record := range records {
		args = append(args, record.id)
		args = append(args, record.user)

	}

	if _, err := tx.ExecContext(ctx, sqlStmt, args...); err != nil {
		return fmt.Errorf("insertion failed with error : %w", err)
	}

	return nil
}

// Truncate___datagen_nested_mysql() truncates the model's table using the shared connection.
func Truncate___datagen_nested_mysql(tx *sql.Tx) error {
	ctx := context.Background()
	if _, err := tx.ExecContext(ctx, "DELETE FROM nested;"); err != nil {
		return fmt.Errorf("delete failed with error : %w", err)
	}
	return nil
}
