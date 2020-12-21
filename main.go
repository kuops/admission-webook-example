package main

import (
	"flag"
	"github.com/gin-gonic/gin"
)

func main() {
	flag.Parse()
	r := gin.Default()
	r.POST("/mutate", mutationHandler)
	_ = r.RunTLS(":443","/etc/webhook/certs/cert.pem","/etc/webhook/certs/key.pem")
}
