package orgs

import "github.com/jenpet/plooral/security"

type Organization struct {
	ID int                                `json:"-"`
	Slug string                           `json:"slug"`
	Name string                           `json:"name"`
	Description string                    `json:"description"`
	Hidden bool                           `json:"hidden"`
	Protected bool                        `json:"protected"`
	Tags []string                         `json:"tags"`
	UserSecurity *security.CredentialSet  `json:"user_credentials,omitempty"`
	OwnerSecurity *security.CredentialSet `json:"owner_credentials,omitempty"`
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

// clearCredentials sets user and owner credentials to nil so that they won't get marshalled into potential JSON responses
func (o *Organization) clearCredentials() {
	o.OwnerSecurity = nil
	o.UserSecurity = nil
}

// partialOrganization is used for user input and updates
type partialOrganization struct {
	*security.PartialCredentialSet
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

func (po *partialOrganization) isHidden() bool {
	return po.Hidden != nil && *po.Hidden != false
}

func (po *partialOrganization) setHidden(hidden bool) {
	po.Hidden = &hidden
}

func (po *partialOrganization) isProtected() bool {
	return po.Protected != nil && *po.Protected != false
}

func (po *partialOrganization) setProtected(protected bool) {
	po.Protected = &protected
}

func (po *partialOrganization) setTags(tags []string) {
	po.Tags = &tags
}

func (po *partialOrganization) setPassword(s string) {
	if po.PartialCredentialSet == nil {
		po.PartialCredentialSet = &security.PartialCredentialSet{}
	}
	po.Password = &s
}

func (po *partialOrganization) setPasswordConfirmation(s string) {
	if po.PartialCredentialSet == nil {
		po.PartialCredentialSet = &security.PartialCredentialSet{}
	}
	po.PasswordConfirmation = &s
}

func (po *partialOrganization) toOrganization() Organization {
	o := Organization{}
	o.mergeWithPartial(*po)
	return o
}