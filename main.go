package main

import (
	"fmt"
	"net/http"

	"github.com/miekg/dns"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	v1 := r.Group("/v1")
	v2 := v1.Group("/v2")

	l := dns.Fqdn("www.voms.com")
	fmt.Println("IsDomainName:", l)
	v2.GET("/JSONP", func(c *gin.Context) {
		data := gin.H{
			"foo": "bar",
		}
		//callback is x
		// Will output  :   x({\"foo\":\"bar\"})
		c.JSONP(http.StatusOK, data)
	})
	r.GET("/someJSON", func(c *gin.Context) {
		names := []string{"lena", "austin", "foo"}

		// Will output  :   while(1);["lena","austin","foo"]
		c.SecureJSON(http.StatusOK, names)
	})
	// Listen and serve on 0.0.0.0:8080
	r.Run(":8088")

	// client
	// curl http://127.0.0.1:8080/JSONP?callback=x
}
