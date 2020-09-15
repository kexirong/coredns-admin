package model

import "errors"

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
		enumVal = "a"
	case TypeCNAME:
		enumVal = "cname"
	case TypePTR:
		enumVal = "ptr"
	case TypeMX:
		enumVal = "mx"
	case TypeTXT:
		enumVal = "txt"
	case TypeAAAA:
		enumVal = "aaa"
	case TypeSRV:
		enumVal = "srv"
	}

	return enumVal
}

// MarshalJSON marshals Type into json.
func (t *Type) MarshalJSON() ([]byte, error) {
	return []byte(t.String()), nil
}

// UnmarshalJSON unmarshals Level from json.
func (t *Type) UnmarshalJSON(b []byte) error {
	switch string(b) {
	case "a":
		*t = TypeA

	case "cname":
		*t = TypeCNAME

	case "ptr":
		*t = TypePTR

	case "mx":
		*t = TypeMX

	case "txt":
		*t = TypeTXT

	case "aaa":
		*t = TypeAAAA

	case "srv":
		*t = TypeSRV

	default:
		return errors.New("invalid Type")
	}

	return nil
}

//Record   api data model
type Record struct {
	Type     Type   `json:"type"`
	TTL      uint32 `json:"ttl,omitempty"`
	Priority int    `json:"priority,omitempty"`
	Name     string `json:"name"`
	Content  string `json:"content,omitempty"`
}
