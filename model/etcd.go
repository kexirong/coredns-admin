package model

import (
	"net"
	"strings"
)

//Etcd is etcd's dns record item
type Etcd struct {
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
	Priority int    `json:"priority,omitempty"`
	Weight   int    `json:"weight,omitempty"`
	Text     string `json:"text,omitempty"`
	Mail     bool   `json:"mail,omitempty"` // Be an MX record. Priority becomes Preference.
	TTL      uint32 `json:"ttl,omitempty"`

	// When a SRV record with a "Host: IP-address" is added, we synthesize
	// a srv.Target domain name.  Normally we convert the full Key where
	// the record lives to a DNS name and use this as the srv.Target.  When
	// TargetStrip > 0 we strip the left most TargetStrip labels from the
	// DNS name.
	TargetStrip int `json:"targetstrip,omitempty"`

	// Group is used to group (or *not* to group) different services
	// together. Services with an identical Group are returned in the same
	// answer.
	Group string `json:"group,omitempty"`

	// Etcd key where we found this service and ignored from json un-/marshalling
	Key string `json:"-"`
}

// HostType returns the DNS type of what is encoded in the Etcd Host field.
func (e *Etcd) HostType() (what Type) {

	if e.Mail {
		return TypeMX
	}
	keyParts := strings.Split(e.Key, "/")
	if strings.Join(keyParts[1:3], "/") == "arpa/in-addr" {
		return TypePTR
	}
	if strings.Join(keyParts[len(keyParts)-2:len(keyParts)], "/") == "dns/ns" {
		return TypeNS
	}
	ip := net.ParseIP(e.Host)

	switch {

	case ip == nil:
		if len(e.Text) == 0 {

			return TypeCNAME
		}
		return TypeTXT

	case ip.To4() != nil:
		return TypeA

	case ip.To4() == nil:
		return TypeAAAA
	}
	// This should never be reached.
	return TypeNone
}
