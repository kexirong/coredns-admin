package redis

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"

	"github.com/kexirong/coredns-admin/model"
	"github.com/kexirong/coredns-admin/service"
	credis "github.com/kexirong/coredns-redis"
	"github.com/miekg/dns"
)

type redis struct {
	Key   string
	Items map[string]string
}

type redisItem struct {
	Field string
	Value string
}

func (i redisItem) ToRedis(key string) *redis {
	r := &redis{Key: key, Items: make(map[string]string)}
	r.Items[i.Field] = i.Value
	return r
}

func (r redis) ToRecords(keyPrefix string) []*model.Record {
	var rs []*model.Record
	const sep = ":"
	keyParts := strings.Split(strings.Trim(r.Key, sep), sep)
	if len(keyParts) < 2 {
		return nil
	}

	for i, j := 0, len(keyParts)-1; i < j; i, j = i+1, j-1 {
		keyParts[i], keyParts[j] = keyParts[j], keyParts[i]
	}
	kpl := len(keyParts) - len(strings.Split(strings.Trim(keyPrefix, sep), sep))

	for k, v := range r.Items {
		tp := dns.StringToType[k]
		switch tp {
		case dns.TypeA, dns.TypeAAAA:
			var ra credis.RecordA
			err := json.Unmarshal([]byte(v), &ra)
			if err != nil {
				continue
			}

			for _, a := range ra {
				r := &model.Record{
					Type:    model.Type(tp),
					TTL:     a.TTL,
					Content: a.IP.String(),
					Key:     r.Key,
					Name:    strings.Join(keyParts[:kpl], "."),
				}
				rs = append(rs, r)
			}

		case dns.TypeCNAME, dns.TypeNS, dns.TypePTR:
			var rc credis.RecordCNANE //[]ItemHost
			err := json.Unmarshal([]byte(v), &rc)
			if err != nil {
				continue
			}
			_n := kpl
			if tp == dns.TypePTR {
				_n -= 2
			}
			for _, c := range rc {
				r := &model.Record{
					Type:    model.Type(tp),
					TTL:     c.TTL,
					Content: c.Host,
					Key:     r.Key,
					Name:    strings.Join(keyParts[:_n], "."),
				}
				rs = append(rs, r)
			}

		case dns.TypeTXT:
			var rt credis.RecordTXT
			err := json.Unmarshal([]byte(v), &rt)
			if err != nil {
				continue
			}
			for _, t := range rt {
				r := &model.Record{
					Type:    model.Type(tp),
					TTL:     t.TTL,
					Content: t.Text,
					Key:     r.Key,
					Name:    strings.Join(keyParts[:kpl], "."),
				}
				rs = append(rs, r)
			}
		case dns.TypeMX:
			var rm credis.RecordMX
			err := json.Unmarshal([]byte(v), &rm)
			if err != nil {
				continue
			}
			for _, m := range rm {
				r := &model.Record{
					Priority: int(m.Preference),
					TTL:      m.TTL,
					Content:  m.Host,
					Key:      r.Key,
					Name:     strings.Join(keyParts[:kpl], "."),
				}
				rs = append(rs, r)
			}
		case dns.TypeSRV:
			var _rs credis.RecordSRV
			err := json.Unmarshal([]byte(v), &_rs)
			if err != nil {
				continue
			}
			for _, s := range _rs {
				r := &model.Record{
					Type:     model.Type(tp),
					Priority: int(s.Priority),
					TTL:      s.TTL,
					Content:  fmt.Sprintf("%d %d %s", s.Weight, s.Port, s.Target),
					Key:      r.Key,
					Name:     strings.Join(keyParts[:kpl], "."),
				}
				rs = append(rs, r)
			}
		case dns.TypeSOA:
			continue

		case dns.TypeCAA:
			var rc credis.RecordCAA
			err := json.Unmarshal([]byte(v), &rc)
			if err != nil {
				continue
			}
			for _, c := range rc {
				r := &model.Record{
					Type:    model.Type(tp),
					TTL:     c.TTL,
					Content: fmt.Sprintf("%d %s %s", c.Flag, c.Tag, c.Value),
					Key:     r.Key,
					Name:    strings.Join(keyParts[:kpl], "."),
				}
				rs = append(rs, r)
			}
		}

	}

	return rs
}
func RedisGetKeys(prefix string) ([]string, error) {
	return service.RedisKeys(prefix + "*")
}

func RedisGetItems(keys []string) (rx []*redis, err error) {
	for _, key := range keys {
		r, err := RedisGetItem(key)
		if err != nil {
			return nil, err
		}
		rx = append(rx, r)
	}
	return
}

func RedisGetItem(key string) (*redis, error) {

	items, err := service.RedisHGetAll(key)
	if err != nil {
		return nil, err
	}
	return &redis{Key: key, Items: items}, nil
}

func RedisGetValue(key, field string) (*redisItem, error) {

	value, err := service.RedisHGet(key, field)
	if err != nil {
		return nil, err
	}
	return &redisItem{Field: field, Value: value}, nil
}

func MergeRedisItem(itemA, itemB *redisItem) error {
	if itemA.Field != itemB.Field {
		return errors.New("itemA's Field is not equal to that of itemB")
	}
	var value []byte
	var err error
	switch itemA.Field {
	case "A", "AAAA":
		value, err = mergeJson[credis.ItemIP]([]byte(itemA.Value), []byte(itemB.Value))

	case "CNAME", "NS", "PTR":
		value, err = mergeJson[credis.ItemHost]([]byte(itemA.Value), []byte(itemB.Value))

	case "TXT":
		value, err = mergeJson[credis.ItemText]([]byte(itemA.Value), []byte(itemB.Value))

	case "MX":
		value, err = mergeJson[credis.ItemMX]([]byte(itemA.Value), []byte(itemB.Value))

	case "SRV":
		value, err = mergeJson[credis.ItemSRV]([]byte(itemA.Value), []byte(itemB.Value))

	case "CAA":
		value, err = mergeJson[credis.ItemCAA]([]byte(itemA.Value), []byte(itemB.Value))

	default:
		err = errors.New("class field invalid")
	}
	if err != nil {
		return err
	}
	itemA.Value = string(value)
	return nil
}

type item interface {
	credis.ItemIP | credis.ItemHost | credis.ItemText | credis.ItemMX | credis.ItemCAA | credis.ItemSRV
}

func mergeJson[T item](item1, item2 []byte) ([]byte, error) {
	var iis1 []T
	var iis2 []T
	err := json.Unmarshal(item1, &iis1)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(item2, &iis2)
	if err != nil {
		return nil, err
	}
	iis1 = append(iis1, iis2...)
	return json.Marshal(iis1)

}

var reSRV = regexp.MustCompile(`^(?P<weight>\d+) (?P<port>\d+) (?P<target>\S+)$`)
var reCAA = regexp.MustCompile(`^(?P<flag>\d+) (?P<tag>\w+) (?P<value>\S+)$`)

func RedisFromRecord(record *model.Record) (*redisItem, error) {
	var ri *redisItem

	switch uint16(record.Type) {
	case dns.TypeA, dns.TypeAAAA:
		var ii credis.ItemIP
		ii.IP = net.ParseIP(record.Content)
		if ii.IP == nil {
			return nil, errors.New("failed to parse content")
		}
		ii.TTL = record.TTL
		tp := record.Type.String()
		val, _ := json.Marshal(&credis.RecordA{ii})
		ri = &redisItem{Field: tp, Value: string(val)}

	case dns.TypeCNAME, dns.TypeNS, dns.TypePTR:
		var ih credis.ItemHost
		ih.Host = record.Content
		ih.TTL = record.TTL
		tp := record.Type.String()
		val, _ := json.Marshal(&credis.RecordCNANE{ih})
		ri = &redisItem{Field: tp, Value: string(val)}

	case dns.TypeTXT:
		var it credis.ItemText
		it.Text = record.Content
		it.TTL = record.TTL
		tp := record.Type.String()
		val, _ := json.Marshal(&credis.RecordTXT{it})
		ri = &redisItem{Field: tp, Value: string(val)}

	case dns.TypeMX:
		var im credis.ItemMX
		im.Host = record.Content
		im.TTL = record.TTL
		im.Preference = uint16(record.Priority)
		tp := record.Type.String()
		val, _ := json.Marshal(&credis.RecordMX{im})
		ri = &redisItem{Field: tp, Value: string(val)}

	case dns.TypeSRV:
		var im credis.ItemSRV
		matchs := reSRV.FindStringSubmatch(record.Content)
		if len(matchs) != 4 {
			return nil, errors.New("content field format is incorrect")
		}
		if v, err := strconv.Atoi(matchs[1]); err == nil {
			im.Weight = uint16(v)
		}
		if v, err := strconv.Atoi(matchs[2]); err == nil {
			im.Port = uint16(v)
		}
		im.Target = matchs[3]
		im.TTL = record.TTL
		im.Priority = uint16(record.Priority)
		tp := record.Type.String()
		val, _ := json.Marshal(&credis.RecordSRV{im})
		ri = &redisItem{Field: tp, Value: string(val)}

	case dns.TypeCAA:
		var ic credis.ItemCAA
		matchs := reCAA.FindStringSubmatch(record.Content)
		if len(matchs) != 4 {
			return nil, errors.New("content field format is incorrect")
		}
		if v, err := strconv.Atoi(matchs[1]); err == nil {
			ic.Flag = uint8(v)
		}
		ic.Tag = matchs[2]
		ic.Value = matchs[3]
		ic.TTL = record.TTL

		tp := record.Type.String()
		val, _ := json.Marshal(&credis.RecordCAA{ic})
		ri = &redisItem{Field: tp, Value: string(val)}

	default:
		return nil, errors.New("type field invalid")
	}

	return ri, nil
}

func RedisSet(r *redis) error {
	_, err := service.RedisHSet(r.Key, r.Items)
	return err
}
