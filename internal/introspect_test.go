//go:build acc

package internal

import (
	"context"
	"testing"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccIntrospect(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{"localhost:9000"},
	})
	require.NoError(t, err)

	err = conn.Ping(ctx)
	if err != nil {
		t.Fatal(err)
	}
	require.NoError(t, err)

	// language="ClickHouse"
	err = conn.Exec(ctx, `
		create table test_acc_introspect (
		    col1 String,
		    col2 Nullable(String)
		) engine Memory
	`)
	require.NoError(t, err)

	// language="ClickHouse"
	defer conn.Exec(ctx, `drop table if exists test_acc_introspect`)

	table, err := Introspect(ctx, conn, "default", "test_acc_introspect")
	require.NoError(t, err)

	assert.Len(t, table.Columns, 2)
	assert.Equal(t, Column{
		Name: "col1",
		Type: "String",
	}, table.Columns[0])
	assert.Equal(t, Column{
		Name: "col2",
		Type: "Nullable(String)",
	}, table.Columns[1])
}
