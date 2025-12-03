package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"
)

// Load___datagen_with_conditionals_postgres executes a single batch of records using the provided transaction.
func Load___datagen_with_conditionals_postgres(records []*__datagen_with_conditionals, tx *sql.Tx) error {
	if len(records) == 0 {
		slog.Warn(fmt.Sprintf("no records to insert for model %s", "with_conditionals"))
		return nil
	}

	ctx := context.Background()

	var b strings.Builder
	columns := []string{
		"\"id\"",
		"\"category\"",
		"\"value\"",
	}
	b.WriteString("INSERT INTO \"with_conditionals\" (")
	b.WriteString(strings.Join(columns, ","))
	b.WriteString(") VALUES ")

	// Build placeholders for Postgres ($1, $2, ... format)
	placeholderCount := 0
	for i := range records {
		if i > 0 {
			b.WriteString(",")
		}
		b.WriteString("(")
		for j := 0; j < 3; j++ {
			if j > 0 {
				b.WriteString(",")
			}
			placeholderCount++
			b.WriteString(fmt.Sprintf("$%d", placeholderCount))
		}
		b.WriteString(")")
	}
	sqlStmt := b.String()

	var args []interface{}
	for _, record := range records {
		args = append(args, record.id)
		args = append(args, record.category)
		args = append(args, record.value)

	}

	if _, err := tx.ExecContext(ctx, sqlStmt, args...); err != nil {
		return fmt.Errorf("insertion failed with error : %w", err)
	}

	return nil
}

// Truncate___datagen_with_conditionals_postgres() truncates the model's table using the shared connection.
func Truncate___datagen_with_conditionals_postgres(tx *sql.Tx) error {
	ctx := context.Background()
	if _, err := tx.ExecContext(ctx, "TRUNCATE TABLE \"with_conditionals\" RESTART IDENTITY CASCADE;"); err != nil {
		return fmt.Errorf("truncate failed with error : %w", err)
	}
	return nil
}
