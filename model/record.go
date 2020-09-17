package model

import (
	"bytes"
	"errors"
	"fmt"
)

// Type is a DNS type.
type Type uint16

//supported types
const (
	TypeNone  = Type(0)
	TypeA     = Type(1)
	TypeNS    = Type(2)
	TypeCNAME = Type(5)
	TypePTR   = Type(12)
	TypeMX    = Type(15)
	TypeTXT   = Type(16)
	TypeAAAA  = Type(28)
	TypeSRV   = Type(33)
)

func (t Type) String() string {
	var enumVal string

	switch t {
	case TypeA:
		enumVal = "A"
	case TypeCNAME:
		enumVal = "CNAME"
	case TypePTR:
		enumVal = "PTR"
	case TypeMX:
		enumVal = "mx"
	case TypeTXT:
		enumVal = "txt"
	case TypeAAAA:
		enumVal = "aaa"
	case TypeSRV:
		enumVal = "srv"
	}

	return fmt.Sprintf(`"%s"`, enumVal)
}

// MarshalJSON marshals Type into json.
func (t Type) MarshalJSON() ([]byte, error) {
	return []byte(t.String()), nil
}

// UnmarshalJSON unmarshals Level from json.
func (t *Type) UnmarshalJSON(b []byte) error {

	switch string(bytes.Trim(b, `"`)) {
	case "A":
		*t = TypeA

	case "CAME":
		*t = TypeCNAME

	case "PTR":
		*t = TypePTR

	case "MX":
		*t = TypeMX

	case "TXT":
		*t = TypeTXT

	case "AAA":
		*t = TypeAAAA

	case "SRV":
		*t = TypeSRV

	default:
		return errors.New("invalid Type")
	}

	return nil
}

//Record api data model
type Record struct {
	Type     Type   `json:"type"`
	TTL      uint32 `json:"ttl,omitempty"`
	Priority int    `json:"priority,omitempty"`
	Name     string `json:"name"`
	Content  string `json:"content"`
	Path     string `json:"-"`
}
