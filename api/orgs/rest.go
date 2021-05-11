package orgs

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
)

var orgs = []Organization{
	{
		"test-org",
		"Test Organization",
	},
	{
		"private-org",
		"Private Organization",
	},
}

func Bootstrap(rg *gin.RouterGroup) {
	rg.GET("/orgs", handleGetAll)
}

func handleGetAll(c *gin.Context) {
	respondWithJSON(c, 200, orgs)
}

func respondWithJSON(c *gin.Context, status int, o interface{}) {
	c.Writer.WriteHeader(status)
	body := map[string]interface{}{"errors": []string{}, "data": o}
	b, _ := json.Marshal(body)
	_,_ = c.Writer.Write(b)
}