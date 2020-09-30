package model

import (
	"strings"
)

type Domain struct {
	Name      string
	SubDomain []*Domain
	idx       uint8
}

func NewDomainTree() *Domain {
	tree := new(Domain)
	tree.idx = 0
	return tree
}

func (d *Domain) AddSubDomain(qdn string) {
	idx := strings.Index(qdn, ".")
	if idx < 0 {
		idx = len(qdn)
	}
	if !d.hasSubDomain(qdn[:idx]) {
		d.SubDomain = append(d.SubDomain, &Domain{
			Name: qdn[:idx],
			idx:  d.idx + 1,
		})
	}
}

func (d *Domain) hasSubDomain(name string) bool {
	for _, i := range d.SubDomain {
		if i.Name == name {
			return true
		}
	}
	return false
}
