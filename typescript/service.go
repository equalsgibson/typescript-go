package typescript

import (
	"fmt"
	"io"
	"reflect"
	"sort"
)

type ConfigFunc func(s *Service)

func WithRegistry(registry map[string]any) ConfigFunc {
	return func(s *Service) {
		s.registry = registry
	}
}

func WithData(data map[string]any) ConfigFunc {
	return func(s *Service) {
		s.data = data
	}
}

func WithRoutes(routes map[string]Route) ConfigFunc {
	return func(s *Service) {
		s.routes = routes
	}
}

func New(
	namespace string,
	configFuncs ...ConfigFunc,
) *Service {
	s := &Service{
		mapping: map[string]string{
			"string":        "string",
			"time.Time":     "string",
			"time.Duration": "number",
			"interface {}":  "any",
			"bool":          "boolean",
			"int":           "number",
			"int8":          "number",
			"int16":         "number",
			"int32":         "number",
			"int64":         "number",
			"uint":          "number",
			"uint8":         "number",
			"uint16":        "number",
			"uint32":        "number",
			"uint64":        "number",
			"uintptr":       "number",
			"float32":       "number",
			"float64":       "number",
			"complex64":     "number",
			"complex128":    "number",
		},
	}

	for _, configFunc := range configFuncs {
		configFunc(s)
	}

	return s
}

type Service struct {
	registry map[string]any
	mapping  map[string]string
	data     map[string]any
	routes   map[string]Route
}

func (s *Service) Generate(writer io.Writer) error {
	keys := []string{}
	for key, entry := range s.registry {
		// Keep track of the keys so they can be sorted and used later
		keys = append(keys, key)

		// This maps the give structs to what they should be converted to when encountered later
		s.mapping[createStandardTypeIdentifier(reflect.ValueOf(entry).Type())] = key
	}

	sort.Strings(keys)

	tsItems := []typescriptGenerator{}

	for _, key := range keys {
		entry := s.registry[key]
		rv := reflect.ValueOf(entry)

		kind := rv.Kind()

		// Look for standard "kinds" in the mapping and just return those as types
		if x, exists := s.mapping[kind.String()]; exists && x != key {
			tsItems = append(tsItems, tsType{
				Name: key,
				Type: x,
			})

			continue
		}

		if kind == reflect.Map {
			tsItems = append(tsItems, tsType{
				Name: key,
				Type: s.convertGoTypeToTypeScriptType(rv.Type()),
			})

			continue
		}

		fields := []tsField{}

		loopOverStructFields(rv, func(fieldDefinition reflect.StructField) {
			tag := parseJSONFieldTag(fieldDefinition.Tag.Get("json"))
			fieldName := fieldDefinition.Name

			if tag.NameOverride != "" {
				fieldName = tag.NameOverride
			}

			if tag.Ignored {
				return
			}

			tsType := s.convertGoTypeToTypeScriptType(fieldDefinition.Type)

			fields = append(fields, tsField{
				Name:     fieldName,
				Type:     tsType,
				Optional: tag.Omitempty,
			})
		})

		inter := tsInterface{
			Name:   key,
			Fields: fields,
		}

		tsItems = append(tsItems, inter)
	}

	// Add route endpoints
	routeNames := []string{}
	for routeName := range s.routes {
		routeNames = append(routeNames, routeName)
	}

	sort.Strings(routeNames)

	for _, routeName := range routeNames {
		route := s.routes[routeName]
		responseBodyType := s.convertGoTypeToTypeScriptType(reflect.ValueOf(route.ResponseBody).Type())
		requestBodyType := ""

		if route.RequestBody != nil {
			requestBodyType = s.convertGoTypeToTypeScriptType(reflect.ValueOf(route.RequestBody).Type())
		}

		params := []tsRouteParam{}

		paramKeys := []string{}
		for key := range route.Params {
			paramKeys = append(paramKeys, key)
		}

		sort.Strings(paramKeys)

		for _, key := range paramKeys {
			value := route.Params[key]
			params = append(params, tsRouteParam{
				Name: key,
				Type: s.convertGoTypeToTypeScriptType(reflect.ValueOf(value).Type()),
			})
		}

		tsItems = append(tsItems, tsRoute{
			Name:            routeName,
			Path:            route.Path,
			Method:          route.Method,
			Params:          params,
			RequestBodyType: requestBodyType,
			ResponseType:    responseBodyType,
		})
	}

	// Add data
	dataVarNames := []string{}
	for dataVarName := range s.data {
		dataVarNames = append(dataVarNames, dataVarName)
	}

	sort.Strings(dataVarNames)

	for _, dataVarName := range dataVarNames {
		data := s.data[dataVarName]
		tsItems = append(tsItems, tsData{
			Name: dataVarName,
			Type: s.convertGoTypeToTypeScriptType(reflect.ValueOf(data).Type()),
			Data: data,
		})
	}

	_, _ = writer.Write([]byte("export namespace GoGenerated {\n"))

	// Write all the items to the Writer
	for i, tsItem := range tsItems {
		s := tsItem.GenerateTypeScript()
		_, _ = writer.Write([]byte(s))
		_, _ = writer.Write([]byte("\n"))

		if i != len(tsItems)-1 {
			_, _ = writer.Write([]byte("\n"))
		} else {
			_, _ = writer.Write([]byte("}\n"))
		}
	}

	return nil
}

func createStandardTypeIdentifier(item reflect.Type) string {
	pkg := item.PkgPath()
	name := item.Name()

	if name != "" {
		if pkg != "" {
			x := pkg + "." + name

			return x
		}

		return name
	}

	return item.String()
}

func (s *Service) convertGoTypeToTypeScriptType(item reflect.Type) string {
	isSlice := item.Kind() == reflect.Slice
	isPointer := item.Kind() == reflect.Pointer
	isMap := item.Kind() == reflect.Map
	stringerTyp := reflect.TypeOf((*fmt.Stringer)(nil)).Elem()
	isStringer := item.Implements(stringerTyp)

	// Read through the pointer/slice/map
	if isPointer || isSlice {
		item = item.Elem()
	}

	if isMap {
		return fmt.Sprintf(
			"{ [key: %s]: %s } | null",
			s.convertGoTypeToTypeScriptType(item.Key()),
			s.convertGoTypeToTypeScriptType(item.Elem()),
		)
	}

	typeFromMapping, found := s.mapping[createStandardTypeIdentifier(item)]
	if !found {
		if isStringer {
			return "string"
		}

		return "unknown"
	}

	if isSlice {
		typeFromMapping += "[]"
	}

	if isSlice || isPointer {
		typeFromMapping += " | null"
	}

	return typeFromMapping
}

func loopOverStructFields(value reflect.Value, fieldHandler func(fieldDefinition reflect.StructField)) {
	if value.Kind() == reflect.Pointer {
		value = value.Elem()
	}

	for i := 0; i < value.NumField(); i++ {
		fieldValue := value.Field(i)
		fieldDefinition := value.Type().Field(i)

		if !fieldDefinition.IsExported() {
			continue
		}

		if fieldDefinition.Type.Kind() == reflect.Struct && fieldDefinition.Anonymous {
			loopOverStructFields(fieldValue, fieldHandler)

			continue
		}

		fieldHandler(fieldDefinition)
	}
}
