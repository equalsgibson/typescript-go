package typescript

import (
	"encoding/json"
	"fmt"
)

type tsData struct {
	Name string
	Type string
	Data any
}

func (ts tsData) GenerateTypeScript() string {
	dataBytes, err := json.MarshalIndent(ts.Data, "\t", "\t")
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("\texport const %s: %s = %s", ts.Name, ts.Type, string(dataBytes))
}
