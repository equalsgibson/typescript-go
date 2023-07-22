package typescript

import (
	"fmt"
	"io"
	"reflect"
	"sort"
)

func New(
	registry map[string]any,
) *Service {
	return &Service{
		registry: registry,
		mapping: map[string]string{
			"string":       "string",
			"time.Time":    "string",
			"interface {}": "any",
			"bool":         "boolean",
			"int":          "number",
			"int8":         "number",
			"int16":        "number",
			"int32":        "number",
			"int64":        "number",
			"uint":         "number",
			"uint8":        "number",
			"uint16":       "number",
			"uint32":       "number",
			"uint64":       "number",
			"uintptr":      "number",
			"float32":      "number",
			"float64":      "number",
			"complex64":    "number",
			"complex128":   "number",
		},
	}
}

type Service struct {
	registry map[string]any
	mapping  map[string]string
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

		inter := tsInterface{
			Name:   key,
			Fields: []tsField{},
		}

		for i := 0; i < rv.NumField(); i++ {
			valueField := rv.Field(i)
			typeField := rv.Type().Field(i)
			actualType := valueField.Type()

			if !typeField.IsExported() {
				continue
			}

			tag := parseJSONFieldTag(typeField.Tag.Get("json"))
			fieldName := typeField.Name
			if tag.NameOverride != "" {
				fieldName = tag.NameOverride
			}

			if tag.Ignored {
				continue
			}

			tsType := s.convertGoTypeToTypeScriptType(actualType)

			inter.Fields = append(inter.Fields, tsField{
				Name:     fieldName,
				Type:     tsType,
				Optional: tag.Omitempty,
			})

		}

		tsItems = append(tsItems, inter)
	}

	for i, tsItem := range tsItems {
		s := tsItem.GenerateTypeScript()
		_, _ = writer.Write([]byte(s))
		_, _ = writer.Write([]byte("\n"))
		if i != len(tsItems)-1 {
			_, _ = writer.Write([]byte("\n"))
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

	// Read through the pointer/slice/map
	if isPointer || isSlice {
		item = item.Elem()
	}

	if isMap {
		return fmt.Sprintf(
			"Map<%s, %s> | null",
			s.convertGoTypeToTypeScriptType(item.Key()),
			s.convertGoTypeToTypeScriptType(item.Elem()),
		)
	}

	typeFromMapping, found := s.mapping[createStandardTypeIdentifier(item)]
	if !found {
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
