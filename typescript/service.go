package typescript

import (
	"fmt"
	"io"
	"reflect"
	"sort"
)

func New(registry map[string]any) *Service {
	return &Service{
		registry: registry,
	}
}

type Service struct {
	registry map[string]any
}

func (s *Service) Generate(writer io.Writer) error {
	mapping := map[string]string{
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
	}

	keys := []string{}
	for key, entry := range s.registry {
		// Keep track of the keys so they can be sorted and used later
		keys = append(keys, key)

		// This maps the give structs to what they should be converted to when encountered later
		mapping[getTypeIdentifier(reflect.ValueOf(entry).Type())] = key
	}

	sort.Strings(keys)

	tsItems := []typescriptGenerator{}

	for _, key := range keys {
		entry := s.registry[key]
		rv := reflect.ValueOf(entry)

		kind := rv.Kind()
		if x, exists := mapping[kind.String()]; exists {
			tsItems = append(tsItems, tsType{
				Name: key,
				Type: x,
			})
			continue
		}

		if kind == reflect.Struct {
			inter := tsInterface{
				Name:   key,
				Fields: []tsField{},
			}
			for i := 0; i < rv.NumField(); i++ {
				valueField := rv.Field(i)
				typeField := rv.Type().Field(i)
				actualType := valueField.Type()

				isSlice := typeField.Type.Kind() == reflect.Slice
				isPointer := typeField.Type.Kind() == reflect.Pointer

				if !typeField.IsExported() {
					continue
				}

				// Read through the pointer/slice
				if isPointer || isSlice {
					actualType = typeField.Type.Elem()
				}

				tag := parseJSONFieldTag(typeField.Tag.Get("json"))
				fieldName := typeField.Name
				if tag.NameOverride != "" {
					fieldName = tag.NameOverride
				}

				if tag.Ignored {
					continue
				}

				// Handle Standard Types
				if x, exists := mapping[getTypeIdentifier(actualType)]; exists {
					inter.Fields = append(inter.Fields, tsField{
						Name:     fieldName,
						Type:     x,
						Array:    isSlice,
						Nullable: isPointer || isSlice,
						Optional: tag.Omitempty,
					})
					continue
				}

				// Maps Handling
				if actualType.Kind() == reflect.Map {
					key := mapping[getTypeIdentifier(actualType.Key())]
					value := mapping[getTypeIdentifier(actualType.Elem())]
					if key != "" && value != "" {
						inter.Fields = append(inter.Fields, tsField{
							Name:     fieldName,
							Type:     fmt.Sprintf("Map<%s, %s>", key, value),
							Array:    isSlice,
							Nullable: isPointer || isSlice,
							Optional: tag.Omitempty,
						})

						continue
					}
				}
			}

			tsItems = append(tsItems, inter)
		}

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

func getTypeIdentifier(item reflect.Type) string {
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
