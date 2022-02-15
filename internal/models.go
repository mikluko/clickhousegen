package internal

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/volatiletech/strmangle"
)

type Table struct {
	Name    string
	Columns []Column
}

func (t Table) ToStruct() (Struct, error) {
	var s = Struct{
		Table:  t,
		Name:   strmangle.TitleCase(strmangle.Singular(t.Name)),
		Fields: make([]Field, len(t.Columns)),
	}
	var err error
	for i, c := range t.Columns {
		var colErr error
		s.Fields[i], colErr = c.ToField()
		if colErr != nil {
			err = multierror.Append(err, colErr)
		}
	}
	return s, err
}

type Column struct {
	Name string
	Type string
}

func (c Column) ToField() (f Field, err error) {
	typ, err := ParseType(c.Type)
	if err != nil {
		return
	}
	f, err = typ.Field()
	if err != nil {
		return
	}
	f.Name = strmangle.TitleCase(c.Name)
	f.Tag = fmt.Sprintf(`ch:"%s" json:"%s" yaml:"%s" toml:"%s" form:"%s"`, c.Name, c.Name, c.Name, c.Name, c.Name)
	f.Column = c
	return f, err
}

type Struct struct {
	Name   string
	Table  Table
	Fields []Field
}

type Field struct {
	Name    string
	Type    string
	Tag     string
	Imports []string
	Consts  map[string]string
	Column  Column
}
