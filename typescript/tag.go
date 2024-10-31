package typescript

import "strings"

type jsonFieldTag struct {
	NameOverride string
	Omitempty    bool
	Ignored      bool
}

func parseJSONFieldTag(tagString string) jsonFieldTag {
	if tagString == "-" {
		return jsonFieldTag{
			Ignored: true,
		}
	}

	parts := strings.Split(tagString, ",")
	if len(parts) == 1 {

		return jsonFieldTag{
			NameOverride: overrideSpecialCharacters(parts[0]),
		}
	}

	tag := jsonFieldTag{}

	for i, part := range parts {
		if i == 0 {
			tag.NameOverride = overrideSpecialCharacters(part)

			continue
		}

		if strings.Contains(part, "omitempty") {
			tag.Omitempty = true
		}
	}

	return tag
}

// NOTE: This is required for json field names that have a "special" character in them.
// These need to be wrapped in " " and accessed via the [] notation.
//
// EX: Timestamp string `json:"@timestamp"` needs to be accessed via Object["@timestamp"]
func overrideSpecialCharacters(tagName string) string {
	runesToCheckFor := []rune{
		64, // @
		32, // ' '
	}

	for _, r := range runesToCheckFor {
		if strings.ContainsRune(tagName, r) {
			return "\"" + tagName + "\""
		}
	}
	return tagName
}
