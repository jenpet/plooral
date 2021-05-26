package orgs

import "github.com/jenpet/plooral/security"

type Organization struct {
	ID int `json:"id"`
	Slug string `json:"slug"`
	Name string `json:"name"`
	Description string `json:"description"`
	Hidden bool `json:"hidden"`
	Protected bool `json:"protected"`
	Tags []string `json:"tags"`
}

func (o *Organization) mergeWithPartial(u partialOrganization) {
	if u.Slug != nil {
		o.Slug = *u.Slug
	}

	if u.Name != nil {
		o.Name = *u.Name
	}

	if u.Description != nil {
		o.Description = *u.Description
	}

	if u.Hidden != nil {
		o.Hidden = *u.Hidden
	}

	if u.Protected != nil {
		o.Protected = *u.Protected
	}

	if u.Tags != nil {
		o.Tags = *u.Tags
	}
}

// partialOrganization is used for user input and updates
type partialOrganization struct {
	security.PartialPasswordSet
	Slug *string `json:"slug"`
	Name *string `json:"name"`
	Description *string `json:"description"`
	Hidden *bool `json:"hidden"`
	Protected *bool `json:"protected"`
	Tags *[]string `json:"tags"`
}

func (po *partialOrganization) setSlug(slug string) {
	po.Slug = &slug
}

func (po *partialOrganization) setName(name string) {
	po.Name = &name
}

func (po *partialOrganization) setDescription(description string) {
	po.Description = &description
}

func (po *partialOrganization) setHidden(hidden bool) {
	po.Hidden = &hidden
}

func (po *partialOrganization) setProtected(protected bool) {
	po.Protected = &protected
}

func (po *partialOrganization) setTags(tags []string) {
	po.Tags = &tags
}

func (po *partialOrganization) toOrganization() Organization {
	o := Organization{}
	o.mergeWithPartial(*po)
	return o
}