package typescript

import (
	"fmt"
	"strings"
)

type typescriptGenerator interface {
	GenerateTypeScript() string
}

type tsInterface struct {
	Name   string
	Fields []tsField
}

func (ts tsInterface) GenerateTypeScript() string {
	fields := []string{}

	for _, field := range ts.Fields {
		fields = append(fields, field.GenerateTypeScript())
	}

	return fmt.Sprintf("export interface %s {\n%s\n}", ts.Name, strings.Join(fields, "\n"))
}

type tsField struct {
	Name     string
	Type     string
	Array    bool
	Nullable bool
	Optional bool
}

func (ts tsField) GenerateTypeScript() string {
	t := ts.Type
	if ts.Array {
		t += "[]"
	}

	if ts.Nullable {
		t += " | null"
	}

	o := ""
	if ts.Optional {
		o = "?"
	}

	return fmt.Sprintf("\t%s%s: %s", ts.Name, o, t)
}

type tsType struct {
	Name string
	Type string
}

func (ts tsType) GenerateTypeScript() string {
	return fmt.Sprintf("export type %s = %s", ts.Name, ts.Type)
}
