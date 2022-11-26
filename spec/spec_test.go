package spec

import (
	"encoding/json"
	"testing"
)

// go test ./spec/ -v -run TestRef
func TestRef(t *testing.T) {
	t.Run("schema", func(t *testing.T) {
		if data, err := json.Marshal(SchemaR{}); err == nil {
			t.Log(string(data))
		} else {
			t.Error(err)
		}
		if data, err := json.Marshal(SchemaR{Ref: "http://some-where"}); err == nil {
			t.Log(string(data))
		} else {
			t.Error(err)
		}
		if data, err := json.Marshal(SchemaR{Schema: &Schema{Type: "string"}}); err == nil {
			t.Log(string(data))
		} else {
			t.Error(err)
		}
	})
}
