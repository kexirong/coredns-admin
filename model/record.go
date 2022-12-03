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
}

type Type uint16

// MarshalJSON marshals Type into json.
func (t Type) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, t.String())), nil
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

func (t Type) String() string {
	return dns.Type(t).String()
}

func (r *Record) Signature() string {
	sum := hashNew()

	sum = hashAdd(sum, r.Name)
	sum = hashAdd(sum, r.Type.String())
	sum = hashAdd(sum, r.Content)
	return fmt.Sprintf("%016x", sum)
}

// Inline and byte-free variant of hash/fnv's fnv64a.

const (
	offset64 = 14695981039346656037
	prime64  = 1099511628211
)

// hashNew initializes a new fnv64a hash value.
func hashNew() uint64 {
	return offset64
}

// hashAdd adds a string to a fnv64a hash value, returning the updated hash.
func hashAdd(h uint64, s string) uint64 {
	for i := 0; i < len(s); i++ {
		h ^= uint64(s[i])
		h *= prime64
	}
	return h
}
