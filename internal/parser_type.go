package internal

import (
	"errors"
	"fmt"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/volatiletech/strmangle"
)

var ErrUnsupportedType = errors.New("unsupported type")

type TypeUnaryUnary struct {
	Type string `parser:"@( 'String' | 'Int8' | 'Int16' | 'Int32' | 'Int64' | 'Int128' | 'Int256' | 'UInt8' | 'UInt16' | 'UInt32' | 'UInt64' | 'UInt128' | 'UInt256' | 'Float32' | 'Float64' | 'Decimal' | 'Bool' | 'UUID' | 'Date' | 'Date32' | 'DateTime' | 'DateTime64' | 'IPv4' | 'IPv6' | 'Point' | 'Ring' | 'Polygon' | 'MultiPolygon' )"`
}

func (t TypeUnaryUnary) Field() (Field, error) {
	switch t.Type {
	case "String":
		return Field{Type: "string"}, nil
	case "Int8":
		return Field{Type: "int8"}, nil
	case "Int16":
		return Field{Type: "int16"}, nil
	case "Int32":
		return Field{Type: "int32"}, nil
	case "Int64":
		return Field{Type: "int64"}, nil
	case "Int128":
		return Field{Type: "big.Int", Imports: []string{"math/big"}}, nil
	case "Int256":
		return Field{Type: "big.Int", Imports: []string{"math/big"}}, nil
	case "UInt8":
		return Field{Type: "uint8"}, nil
	case "UInt16":
		return Field{Type: "uint16"}, nil
	case "UInt32":
		return Field{Type: "uint32"}, nil
	case "UInt64":
		return Field{Type: "uint64"}, nil
	case "UInt128":
		return Field{Type: "big.Int", Imports: []string{"math/big"}}, nil
	case "UInt256":
		return Field{Type: "big.Int", Imports: []string{"math/big"}}, nil
	case "Float32":
		return Field{Type: "float32"}, nil
	case "Float64":
		return Field{Type: "float64"}, nil
	case "Decimal":
		return Field{Type: "decimal.Decimal", Imports: []string{"github.com/shopspring/decimal"}}, nil
	case "Bool":
		return Field{Type: "bool"}, nil
	case "UUID":
		return Field{Type: "uuid.UUID", Imports: []string{"github.com/google/uuid"}}, nil
	case "Date":
		return Field{Type: "time.Time", Imports: []string{"time"}}, nil
	case "Date32":
		return Field{Type: "time.Time", Imports: []string{"time"}}, nil
	case "DateTime":
		return Field{Type: "time.Time", Imports: []string{"time"}}, nil
	case "DateTime64":
		return Field{Type: "time.Time", Imports: []string{"time"}}, nil
	case "IPv4":
		return Field{Type: "net.IP", Imports: []string{"net"}}, nil
	case "IPv6":
		return Field{Type: "net.IP", Imports: []string{"net"}}, nil
	case "Point":
		return Field{Type: "orb.Point", Imports: []string{"github.com/paulmach/orb"}}, nil
	case "Ring":
		return Field{Type: "orb.Ring", Imports: []string{"github.com/paulmach/orb"}}, nil
	case "Polygon":
		return Field{Type: "orb.Polygon", Imports: []string{"github.com/paulmach/orb"}}, nil
	case "MultiPolygon":
		return Field{Type: "orb.MultiPolygon", Imports: []string{"github.com/paulmach/orb"}}, nil
	default:
		return Field{}, fmt.Errorf("%w: %s", ErrUnsupportedType, t.Type)
	}
}

type TypeNullableUnary struct {
	Type TypeUnary `parser:"'Nullable' '(' @@ ')'"`
}

func (t TypeNullableUnary) Field() (Field, error) {
	f, err := t.Type.Field()
	if err != nil {
		return f, err
	}
	f.Type = "*" + f.Type
	return f, nil
}

type TypeLowCardinalityUnary struct {
	Type TypeUnary `parser:"'LowCardinality' '(' @@ ')'"`
}

func (t TypeLowCardinalityUnary) Field() (Field, error) {
	return t.Type.Field()
}

type TypeFixedString struct {
	Length int `parser:"'FixedString' '(' @Int ')'"`
}

func (t TypeFixedString) Field() (Field, error) {
	return Field{Type: "string"}, nil
}

type TypeUnary struct {
	Unary          *TypeUnaryUnary          `parser:"  @@"`
	Nullable       *TypeNullableUnary       `parser:"| @@"`
	LowCardinality *TypeLowCardinalityUnary `parser:"| @@"`
	FixedString    *TypeFixedString         `parser:"| @@"` // it's not exactly unary at db side, but resulting Go type is merely a string
}

func (t TypeUnary) Field() (f Field, err error) {
	switch {
	case t.Unary != nil:
		return t.Unary.Field()
	case t.Nullable != nil:
		return t.Nullable.Field()
	case t.LowCardinality != nil:
		return t.LowCardinality.Field()
	case t.FixedString != nil:
		return t.FixedString.Field()
	default:
		panic("must never happen")
	}
}

type TypeArray struct {
	Type TypeUnary `parser:"  'Array' '(' @@ ')'"`
}

func (t *TypeArray) Field() (Field, error) {
	f, err := t.Type.Field()
	if err != nil {
		return f, err
	}
	f.Type = "[]" + f.Type
	return f, nil
}

type TypeMap struct {
	KeyType   TypeUnary `parser:"'Map' '(' @@"`
	ValueType TypeUnary `parser:"            ',' @@ ')'"`
}

func (t TypeMap) Field() (f Field, err error) {
	kf, err := t.KeyType.Field()
	if err != nil {
		return f, err
	}
	if kf.Imports != nil {
		f.Imports = append(f.Imports, kf.Imports...)
	}
	vf, err := t.ValueType.Field()
	if err != nil {
		return f, err
	}
	if vf.Imports != nil {
		f.Imports = append(f.Imports, vf.Imports...)
	}
	f.Type = fmt.Sprintf("map[%s]%s", kf.Type, vf.Type)
	return f, nil
}

type TypeEnum struct {
	Elems []TypeEnumElem `parser:"('Enum8' | 'Enum16') '(' ( @@ ','? )+ ')'"`
}

func (t TypeEnum) Field() (Field, error) {
	f := Field{
		Type: "string",
	}
	f.Consts = make(map[string]string)
	for i := range t.Elems {
		v := strings.Trim(t.Elems[i].Label, "'")
		f.Consts[strmangle.TitleCase(v)] = v
	}
	return f, nil
}

type TypeEnumElem struct {
	Label string `parser:"     @(String|Char)"`
	Value int    `parser:" '=' @(Int)"`
}

type Type struct {
	Unary *TypeUnary `parser:"  @@"`
	Array *TypeArray `parser:"| @@"`
	Map   *TypeMap   `parser:"| @@"`
	Enum  *TypeEnum  `parser:"| @@"`
}

func (t Type) Field() (Field, error) {
	switch {
	case t.Unary != nil:
		return t.Unary.Field()
	case t.Array != nil:
		return t.Array.Field()
	case t.Map != nil:
		return t.Map.Field()
	case t.Enum != nil:
		return t.Enum.Field()
	default:
		panic("must never happen")
	}
}

var typeParser = participle.MustBuild(&Type{}, participle.Unquote())

func ParseType(exp string) (*Type, error) {
	var obj Type
	err := typeParser.ParseString("", exp, &obj)
	if err != nil {
		return nil, err
	}
	return &obj, nil
}
