package orgs

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jenpet/plooral/errors"
	"net/http"
)

func Bootstrap(rg *gin.RouterGroup) {
	a := api{s: newDefaultService()}
	orgs := rg.Group("/orgs")
	orgs.GET("", a.handleGetAll)
	orgs.GET("/:orgSlug", a.handleGetOrganization)
	orgs.POST("", a.handleCreateOrganization)
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
	var body createOrganizationBody
	err := c.ShouldBindJSON(&body)
	if err != nil {
		respondWithJSON(c, http.StatusBadRequest, nil, err)
		return
	}
	org, err := a.s.UpsertOrganization(body.toDomain())
	if err != nil {
		respondWithJSON(c, errors.ErrStatusCode(err), nil, err)
		return
	}
	respondWithJSON(c, http.StatusCreated, org, nil)
}

func respondWithJSON(c *gin.Context, status int, o interface{}, err error) {
	if status <= 0 {
		status = http.StatusInternalServerError
	}
	c.Writer.WriteHeader(status)
	var errs []string
	if err != nil {
		errs = []string{err.Error()}
	}
	body := map[string]interface{}{"errors": errs, "data": o}
	b, _ := json.Marshal(body)
	_,_ = c.Writer.Write(b)
}

type createOrganizationBody struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
	Description string `json:"description"`
	Hidden bool `json:"hidden"`
	Protected bool `json:"protected"`
	Tags []string `json:"tags"`
}

func (cob createOrganizationBody) toDomain() Organization {
	return Organization{
		Slug: cob.Slug,
		Name: cob.Name,
		Description: cob.Description,
		Hidden: cob.Hidden,
		Protected: cob.Protected,
		Tags: cob.Tags,
	}
}