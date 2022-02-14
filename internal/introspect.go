package internal

import (
	"context"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/hashicorp/go-multierror"
)

func Introspect(ctx context.Context, conn clickhouse.Conn, dbName, tableName string) (t Table, err error) {
	// language="ClickHouse"
	const query = `SELECT name, type FROM system.columns WHERE database=$1 AND table=$2`
	rows, err := conn.Query(ctx, query, dbName, tableName)
	if err != nil {
		return
	}
	for rows.Next() {
		var col Column
		rowErr := rows.Scan(&col.Name, &col.Type)
		if rowErr != nil {
			err = multierror.Append(err, rowErr)
		}
		t.Columns = append(t.Columns, col)
	}
	t.Name = tableName
	return
}
