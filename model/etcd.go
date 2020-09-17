package model

import (
	"fmt"
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
func (e Etcd) HostType() Type {

	if e.Mail {
		return TypeMX
	}
	keyParts := strings.Split(strings.Trim(e.Key, "/"), "/")
	if strings.Join(keyParts[1:3], "/") == "arpa/in-addr" {
		return TypePTR
	}
	if strings.Join(keyParts[len(keyParts)-2:], "/") == "dns/ns" {
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

//ToRecord   Conversion to Record
func (e Etcd) ToRecord() *Record {
	r := &Record{
		TTL:      e.TTL,
		Priority: e.Priority,
	}
	keyParts := strings.Split(strings.Trim(e.Key, "/"), "/")
	r.Path = keyParts[0]
	tp := e.HostType()
	switch tp {
	case TypeA, TypeAAAA, TypeCNAME, TypeMX, TypeTXT, TypeSRV, TypePTR:
		r.Type = tp
		n := 1
		if tp == TypePTR {
			n = 3
		}
		for i, j := n, len(keyParts)-1; i < j; i, j = i+1, j-1 {
			keyParts[i], keyParts[j] = keyParts[j], keyParts[i]
		}
		r.Name = strings.Join(keyParts[n:], ".")

	//ns not suport
	case TypeNS:
		return nil

	}
	switch tp {
	case TypeA, TypeAAAA, TypeCNAME, TypeMX, TypePTR:
		r.Content = e.Host
	case TypeTXT:
		r.Content = e.Text
	case TypeSRV:
		r.Content = fmt.Sprintf("%d %d %d %s", e.Priority, e.Weight, e.Port, e.Host)
	}

	return r

}
