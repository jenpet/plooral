package api

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jenpet/plooral/boards"
	"github.com/jenpet/plooral/orgs"
	"net/http"
	"time"
)

const apiBasePath = "api/v1"

func Serve() {
	r := Server()
	_ = r.Run(":8079")
}

func Server() *gin.Engine {
	r := gin.New()
	r.Use(cors.Default())
	v1 := r.Group(apiBasePath)
	registerHandlers(v1)
	return r
}

func registerHandlers(r *gin.RouterGroup) {
	r.GET("/info", handleGetInfo)
	orgs.Bootstrap(r)
	boards.Bootstrap(r)
}

func handleGetInfo(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusOK)
	body := fmt.Sprintf(`{"status":"up", "ts": %d}`, time.Now().UTC().Unix())
	_, _ = c.Writer.WriteString(body)
}