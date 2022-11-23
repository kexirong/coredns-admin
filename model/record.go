package model

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/miekg/dns"
)

//Record api data model
type Record struct {
	Type     Type   `json:"type"`
	TTL      uint32 `json:"ttl,omitempty"`
	Priority int    `json:"priority,omitempty"`
	Name     string `json:"name"`
	Content  string `json:"content"`
	Key      string `json:"key,omitempty"`
	//Path     string `json:"-"`
}

type Type uint16

// MarshalJSON marshals Type into json.
func (t Type) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, dns.Type(t).String())), nil
}

// UnmarshalJSON unmarshals Level from json.
func (t *Type) UnmarshalJSON(b []byte) error {
	t1, ok := dns.StringToType[string(bytes.Trim(b, `"`))]
	if !ok {
		return errors.New("invalid Type")
	}

	*t = Type(t1)
	return nil

}
