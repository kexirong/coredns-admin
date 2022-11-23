package etcd

import (
	"strings"
	"testing"
)

func TestGrowBasicPrefix(t *testing.T) {
	if growBasicPrefix("") != "#1" {
		t.Fatal()
	}
	if growBasicPrefix("#") != "#1" {
		t.Fatal()
	}
	if growBasicPrefix("#1") != "#2" {
		t.Fatal()
	}
	if growBasicPrefix("#19") != "#191" {
		t.Fatal()
	}
}

func TestReSrv(t *testing.T) {
	if reSRV.MatchString("1122 333222") {
		t.Fatal()
	}
	if reSRV.MatchString("1122 333222") {
		t.Fatal()
	}
	if !reSRV.MatchString("1122 333 wwww") {
		t.Fatal()
	}
	ss := reSRV.FindStringSubmatch(strings.Trim("1122 333 213.2.3. ", " ."))
	t.Fatal(len(ss), ss[3])
}
