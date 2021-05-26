package api

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jenpet/plooral/boards"
	"github.com/jenpet/plooral/database"
	"github.com/jenpet/plooral/orgs"
	"github.com/jenpet/plooral/security"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const apiBasePath = "api/v1"

func Serve() {
	r := Server()
	_ = r.Run(":8079")
}

func Server() *gin.Engine {
	// apply database migrations
	database.DefaultMigrate(logrus.StandardLogger())

	r := gin.New()
	r.Use(cors.Default())
	v1 := r.Group(apiBasePath)
	setupDependencies(v1)
	return r
}

func setupDependencies(r *gin.RouterGroup) {
	r.GET("/info", handleGetInfo)
	secSvc := security.Bootstrap()
	orgs.Bootstrap(r, secSvc)
	boards.Bootstrap(r)
}

func handleGetInfo(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusOK)
	body := fmt.Sprintf(`{"status":"up", "ts": %d}`, time.Now().UTC().Unix())
	_, _ = c.Writer.WriteString(body)
}