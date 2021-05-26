package orgs

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jenpet/plooral/errors"
	"github.com/jenpet/plooral/rest"
	"net/http"
)

func Bootstrap(rg *gin.RouterGroup, pw passwordService) {
	a := api{s: newDefaultService(pw)}
	orgs := rg.Group("/orgs")
	orgs.GET("", a.handleGetAll)
	orgs.GET("/:orgSlug", a.handleGetOrganization)
	orgs.POST("", a.handleCreateOrganization)
	orgs.PATCH("/:orgSlug", a.handleUpdateOrganization)
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

func (a *api) handleGetOrganization(c *gin.Context) {
	org, err := a.s.OrganizationBySlug(c.Param("orgSlug"))
	if err != nil {
		respondWithJSON(c, errors.ErrStatusCode(err), nil, err)
		return
	}
	respondWithJSON(c, http.StatusOK, org, nil)
}

func (a *api) handleCreateOrganization(c *gin.Context) {
	var body partialOrganization
	err := c.ShouldBindJSON(&body)
	if err != nil {
		respondWithJSON(c, http.StatusBadRequest, nil, err)
		return
	}
	org, err := a.s.CreateOrganization(body)
	if err != nil {
		respondWithJSON(c, errors.ErrStatusCode(err), nil, err)
		return
	}
	respondWithJSON(c, http.StatusCreated, org, nil)
}

func (a *api) handleUpdateOrganization(c *gin.Context) {
	var body partialOrganization
	err := c.ShouldBindJSON(&body)
	if err != nil {
		respondWithJSON(c, http.StatusBadRequest, nil, err)
		return
	}
	if body.Slug != nil && *body.Slug != c.Param("orgSlug") {
		respondWithJSON(c, http.StatusBadRequest, nil,
			errors.E("Path slug does not match body slug", rest.KUserInputInvalid))
		return
	}
	update, err := a.s.UpdateOrganization(body)
	if err != nil {
		respondWithJSON(c, http.StatusBadRequest, nil, err)
		return
	}
	respondWithJSON(c, http.StatusOK, update, nil)
}

func respondWithJSON(c *gin.Context, status int, o interface{}, err error) {
	if status <= 0 {
		status = http.StatusInternalServerError
	}
	c.Writer.WriteHeader(status)
	c.Header("Content-Type", "application/json")
	var errs []string
	if err != nil {
		errs = []string{err.Error()}
	}
	body := map[string]interface{}{"errors": errs, "data": o}
	b, _ := json.Marshal(body)
	_,_ = c.Writer.Write(b)
}