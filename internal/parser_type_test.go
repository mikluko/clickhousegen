package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testParseTypeItems = []struct {
	In      string
	Type    string
	Imports []string
	Consts  map[string]string
	Err     error
}{
	{In: "String", Type: "string"},
	{In: "Nullable(String)", Type: "*string"},
	{In: "FixedString(2)", Type: "string"},
	{In: "Nullable(FixedString(2))", Type: "*string"},
	{In: "Int8", Type: "int8"},
	{In: "Nullable(Int8)", Type: "*int8"},
	{In: "Int16", Type: "int16"},
	{In: "Int32", Type: "int32"},
	{In: "Int64", Type: "int64"},
	{In: "Int128", Type: "big.Int", Imports: []string{"math/big"}},
	{In: "Int256", Type: "big.Int", Imports: []string{"math/big"}},
	{In: "UInt8", Type: "uint8"},
	{In: "UInt16", Type: "uint16"},
	{In: "UInt32", Type: "uint32"},
	{In: "UInt64", Type: "uint64"},
	{In: "UInt128", Type: "big.Int", Imports: []string{"math/big"}},
	{In: "UInt256", Type: "big.Int", Imports: []string{"math/big"}},
	{In: "Float32", Type: "float32"},
	{In: "Float64", Type: "float64"},
	{In: "Decimal", Type: "decimal.Decimal", Imports: []string{"github.com/shopspring/decimal"}},
	{In: "Bool", Type: "bool"}, {In: "UUID", Type: "uuid.UUID", Imports: []string{"github.com/google/uuid"}},
	{In: "Date", Type: "time.Time", Imports: []string{"time"}},
	{In: "Date32", Type: "time.Time", Imports: []string{"time"}},
	{In: "DateTime", Type: "time.Time", Imports: []string{"time"}},
	{In: "DateTime64", Type: "time.Time", Imports: []string{"time"}},
	{In: "IPv4", Type: "net.IP", Imports: []string{"net"}},
	{In: "IPv6", Type: "net.IP", Imports: []string{"net"}},
	{In: "Point", Type: "orb.Point", Imports: []string{"github.com/paulmach/orb"}},
	{In: "Ring", Type: "orb.Ring", Imports: []string{"github.com/paulmach/orb"}},
	{In: "Polygon", Type: "orb.Polygon", Imports: []string{"github.com/paulmach/orb"}},
	{In: "MultiPolygon", Type: "orb.MultiPolygon", Imports: []string{"github.com/paulmach/orb"}},
	{In: "Enum8( 'a' = 1, 'b' = 2 )", Type: "string", Consts: map[string]string{"A": "a", "B": "b"}},
}

func TestParseType(t *testing.T) {
	for _, item := range testParseTypeItems {
		t.Run(item.In, func(t *testing.T) {
			obj, err := ParseType(item.In)
			require.NoError(t, err)

			typ, err := obj.Field()
			if item.Err != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, item.Err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, item.Type, typ.Type)
				assert.Equal(t, item.Imports, typ.Imports)
				assert.Equal(t, item.Consts, typ.Consts)
			}
		})
	}
}
