package httpx

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func QueryToCSV(ctx context.Context, pool *pgxpool.Pool, sql string, args ...any) ([]byte, error) {
	rows, err := pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	fields := rows.FieldDescriptions()
	headers := make([]string, len(fields))
	for i, f := range fields {
		headers[i] = string(f.Name)
	}

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)

	if err := w.Write(headers); err != nil {
		return nil, err
	}

	for rows.Next() {
		vals, err := rows.Values()
		if err != nil {
			return nil, err
		}
		rec := make([]string, len(vals))
		for i, v := range vals {
			if v == nil {
				rec[i] = ""
				continue
			}
			rec[i] = fmt.Sprint(v)
		}
		if err := w.Write(rec); err != nil {
			return nil, err
		}
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
