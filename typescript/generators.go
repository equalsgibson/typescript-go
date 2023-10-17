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

	return fmt.Sprintf("\texport type %s = {\n%s\n\t}", ts.Name, strings.Join(fields, "\n"))
}

type tsField struct {
	Name     string
	Type     string
	Optional bool
}

func (ts tsField) GenerateTypeScript() string {
	o := ""
	if ts.Optional {
		o = "?"
	}

	return fmt.Sprintf("\t\t%s%s: %s", ts.Name, o, ts.Type)
}

type tsType struct {
	Name string
	Type string
}

func (ts tsType) GenerateTypeScript() string {
	return fmt.Sprintf("\texport type %s = %s", ts.Name, ts.Type)
}
