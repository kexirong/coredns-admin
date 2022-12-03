package etcd

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/coredns/coredns/plugin/etcd/msg"
	"github.com/kexirong/coredns-admin/model"
	"github.com/kexirong/coredns-admin/service"
	"github.com/miekg/dns"
	"go.etcd.io/etcd/api/v3/mvccpb"
)

//Etcd is etcd's dns record item
type etcd struct {
	msg.Service
}

// HostType returns the DNS type of what is encoded in the Etcd Host field.
func (e etcd) HostType() uint16 {

	if e.Mail {
		return dns.TypeMX
	}
	if e.Port > 0 {
		return dns.TypeSRV
	}

	if strings.Contains(e.Key, "arpa/in-addr") {
		return dns.TypePTR
	}

	if strings.Contains(e.Key, "dns/ns") {
		return dns.TypeNS
	}
	ip := net.ParseIP(e.Host)

	switch {

	case ip == nil:
		if len(e.Text) == 0 {
			return dns.TypeCNAME
		}
		return dns.TypeTXT

	case ip.To4() != nil:
		return dns.TypeA

	case ip.To4() == nil:
		return dns.TypeAAAA
	}
	// This should never be reached.
	return dns.TypeNone
}

//ToRecord   Conversion to Record
func (e etcd) ToRecord(keyPrefix string) *model.Record {
	const sep = "/"
	r := &model.Record{
		TTL:      e.TTL,
		Priority: e.Priority,
		Key:      e.Key,
	}
	keyParts := strings.Split(strings.Trim(e.Key, sep), sep)
	if len(keyParts) < 2 {
		return nil
	}
	// r.Path = fmt.Sprintf("/%s/", keyParts[0])

	if keyParts[len(keyParts)-1][0] == '#' {
		keyParts = keyParts[:len(keyParts)-1]
	}

	tp := e.HostType()
	switch tp {
	case dns.TypeA, dns.TypeAAAA, dns.TypeCNAME, dns.TypeMX, dns.TypeTXT, dns.TypeSRV, dns.TypePTR:
		r.Type = model.Type(tp)
		n := len(strings.Split(strings.Trim(keyPrefix, sep), sep))
		if tp == dns.TypePTR {
			n += 2
		}
		for i, j := n, len(keyParts)-1; i < j; i, j = i+1, j-1 {
			keyParts[i], keyParts[j] = keyParts[j], keyParts[i]
		}

		r.Name = strings.Join(keyParts[n:], ".")

	//ns not suport
	case dns.TypeNS:
		return nil
	}
	switch tp {
	case dns.TypeA, dns.TypeAAAA, dns.TypeCNAME, dns.TypeMX, dns.TypePTR:
		r.Content = e.Host
	case dns.TypeTXT:
		r.Content = e.Text
	case dns.TypeSRV:
		r.Content = fmt.Sprintf("%d %d %s", e.Weight, e.Port, e.Host)
	}

	return r
}

var reSRV = regexp.MustCompile(`^(?P<weight>\d+) (?P<port>\d+) (?P<target>\S+)$`)

func EtcdFromRecord(r *model.Record, prefix string) (*etcd, error) {
	var e = new(etcd)
	if r.TTL != 0 {
		e.TTL = r.TTL
	}
	if r.Priority != 0 {
		e.Priority = r.Priority
	}
	switch uint16(r.Type) {
	case dns.TypeMX:
		e.Mail = true
		fallthrough
	case dns.TypeA, dns.TypeCNAME, dns.TypePTR, dns.TypeAAAA:
		e.Host = r.Content
	// case TypeNS:

	case dns.TypeTXT:
		e.Text = r.Content
	case dns.TypeSRV:

		matchs := reSRV.FindStringSubmatch(r.Content)
		if len(matchs) != 4 {
			return nil, errors.New("content field format is incorrect")
		}
		e.Weight, _ = strconv.Atoi(matchs[1])
		e.Port, _ = strconv.Atoi(matchs[2])
		e.Host = matchs[3]
	default:
		return nil, errors.New("type field invalid")
	}
	prefixPart := []string{prefix}
	// if !strings.HasSuffix(_path, "/") {
	// 	_path += "/"
	// }

	if uint16(r.Type) == dns.TypePTR {
		//_path += "arpa/in-addr/"
		prefixPart = append(prefixPart, "arpa", "in-addr")
	}
	keys := dns.SplitDomainName(r.Name)
	for i, j := 0, len(keys)-1; i < j; i, j = i+1, j-1 {
		keys[i], keys[j] = keys[j], keys[i]
	}
	e.Key = path.Join(append(prefixPart, keys...)...)
	return e, nil
}

//EtcdGetItems get  etcd items with path as prifix
func EtcdGetItems(path string) (ex []*etcd, err error) {
	kvs, err := service.EtcdGetKvs(path)
	if err != nil {
		return nil, err
	}

	for _, kv := range kvs {

		e := new(etcd)
		e.Key = string(kv.Key)

		if err := json.Unmarshal(kv.Value, e); err != nil {
			//log.Println(err.Error())
			continue
		}

		ex = append(ex, e)

	}
	return ex, nil
}

func EtcdPutItem(etcd *etcd) (err error) {

	value, err := json.Marshal(etcd)
	if err != nil {
		return err
	}

	return service.EtcdPutKv(etcd.Key, string(value))
}

func GetValueFromKVS(kvs []*mvccpb.KeyValue, key string) []byte {
	for _, kv := range kvs {
		if string(kv.Key) == key {
			return kv.Value
		}
	}
	return nil
}
