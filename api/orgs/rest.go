package orgs

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Bootstrap(rg *gin.RouterGroup) {
	a := api{s: newDefaultService()}
	rg.GET("/orgs", a.handleGetAll)
}

type api struct {
	s *Service
}

func (a *api) handleGetAll(c *gin.Context) {
	orgs, err := a.s.AllOrganizations(false)
	if err != nil {
		respondWithJSON(c, http.StatusInternalServerError, nil, err)
		return
	}
	respondWithJSON(c, http.StatusOK, orgs, nil)
}

func respondWithJSON(c *gin.Context, status int, o interface{}, err error) {
	c.Writer.WriteHeader(status)
	var errs []string
	if err != nil {
		errs = []string{err.Error()}
	}
	body := map[string]interface{}{"errors": errs, "data": o}
	b, _ := json.Marshal(body)
	_,_ = c.Writer.Write(b)
}