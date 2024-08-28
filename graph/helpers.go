package graph

import (
	"encoding/json"
	"fmt"
	"io"
)

// JSON is a custom scalar type to handle arbitrary JSON objects
type JSON map[string]interface{}

// MarshalGQL marshals the JSON value to a GraphQL string
func (j JSON) MarshalGQL(w io.Writer) {
	b, err := json.Marshal(j)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`"error": "%v"`, err)))
		return
	}
	w.Write(b)
}

// UnmarshalGQL unmarshals a GraphQL value to a JSON value
func (j *JSON) UnmarshalGQL(v interface{}) error {
	switch v := v.(type) {
	case string:
		return json.Unmarshal([]byte(v), j)
	case map[string]interface{}:
		*j = v
		return nil
	default:
		return fmt.Errorf("json: cannot unmarshal %T", v)
	}
}
