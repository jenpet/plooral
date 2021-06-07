package orgs

import (
	"github.com/gin-gonic/gin"
	"github.com/jenpet/plooral/errors"
	"github.com/jenpet/plooral/rest"
	"net/http"
)

func Bootstrap(rg *gin.RouterGroup, pw credentialService) {
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
		rest.RespondWithJSONError(c.Writer, err)
		return
	}
	rest.RespondWithJSONData(c.Writer, http.StatusOK, orgs)
}

func (a *api) handleGetOrganization(c *gin.Context) {
	org, err := a.s.OrganizationBySlug(c, c.Param("orgSlug"))
	if err != nil {
		rest.RespondWithJSONError(c.Writer, err)
		return
	}
	rest.RespondWithJSONData(c.Writer, http.StatusOK, org)
}

func (a *api) handleCreateOrganization(c *gin.Context) {
	var body partialOrganization
	err := c.ShouldBindJSON(&body)
	if err != nil {
		rest.RespondWithJSONError(c.Writer, err)
		return
	}
	org, err := a.s.CreateOrganization(c, body)
	if err != nil {
		rest.RespondWithJSONError(c.Writer, err)
		return
	}
	rest.RespondWithJSONData(c.Writer, http.StatusCreated, org)
}

func (a *api) handleUpdateOrganization(c *gin.Context) {
	var body partialOrganization
	err := c.ShouldBindJSON(&body)
	if err != nil {
		rest.RespondWithJSONError(c.Writer, rest.WrapUserInputInvalidError(err))
		return
	}
	if body.Slug != nil && *body.Slug != c.Param("orgSlug") {
		err = rest.WrapUserInputInvalidError(errors.E("path slug does not match body slug"))
		rest.RespondWithJSONError(c.Writer, err)
		return
	}
	update, err := a.s.UpdateOrganization(c, body)
	if err != nil {
		rest.RespondWithJSONError(c.Writer, err)
		return
	}
	rest.RespondWithJSONData(c.Writer, http.StatusOK, update)
}
