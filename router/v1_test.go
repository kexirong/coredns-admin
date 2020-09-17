package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/kexirong/coredns-admin/config"
	"github.com/kexirong/coredns-admin/service"
)

func TestPingRoute(t *testing.T) {

	err := service.InitEtcdClient(*config.Get())
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/records", nil)
	router.ServeHTTP(w, req)
	t.Log("Code: ", w.Code)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}
