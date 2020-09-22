package model

import (
	"encoding/json"
	"testing"
)

func TestRecordMarshalJSON(t *testing.T) {
	r := Record{TypeA, 10, 12, "test", "1.1.1.1", "/core.com.test", "core"}
	_, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
}
func TestRecordUnmarshalJSON(t *testing.T) {
	ss := `{"type":"A","ttl":10,"priority":12,"name":"test","content":"1.1.1.1"}`
	r := &Record{}
	err := json.Unmarshal([]byte(ss), r)
	if err != nil {
		t.Fatal(err)
	}
}
