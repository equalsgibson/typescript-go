package typescript

import (
	"fmt"
	"strings"
)

type Route struct {
	Path         string
	Method       string
	Params       map[string]any
	RequestBody  any
	ResponseBody any
}

type tsRoute struct {
	Name            string
	Path            string
	Method          string
	Params          []tsRouteParam
	RequestBodyType string
	ResponseType    string
}

type tsRouteParam struct {
	Name string
	Type string
}

func (ts tsRoute) GenerateTypeScript() string {
	arguments := []string{}
	for _, param := range ts.Params {
		arguments = append(
			arguments,
			fmt.Sprintf("%s: %s", param.Name, param.Type),
		)
	}

	if ts.RequestBodyType != "" {
		arguments = append(
			arguments,
			fmt.Sprintf("payload: %s", ts.RequestBodyType),
		)
	}

	output := fmt.Sprintf("export const %s = (%s) => {\n", ts.Name, strings.Join(arguments, ", "))

	if len(ts.Params) > 0 {
		output += "\tconst params = {\n"
		for _, param := range ts.Params {
			output += fmt.Sprintf("\t\t%s: %s,\n", param.Name, param.Name)
		}
		output += "\t}\n\n"

		output += "\tconst queryString = Object.keys(params).map((key) => {\n"
		output += "\t\treturn encodeURIComponent(key) + \"=\" + encodeURIComponent(params[key])\n"
		output += "\t}).join(\"&\")\n\n"

		output += fmt.Sprintf("\treturn fetch(`%s?${queryString}`, {\n", ts.Path)
	} else {
		output += fmt.Sprintf("\treturn fetch(\"%s\", {\n", ts.Path)
	}

	output += fmt.Sprintf("\t\tmethod: \"%s\",\n", ts.Method)

	if ts.RequestBodyType != "" {
		output += "\t\tbody: JSON.stringify(payload),\n"
	}

	output += fmt.Sprintf("\t}).then<%s>((response => response.json()))\n", ts.ResponseType)
	output += "}"

	return output
}
