package orgs

type Organization struct {
	ID int `json:"id"`
	Slug string `json:"slug"`
	Name string `json:"name"`
	Description string `json:"description"`
	Hidden bool `json:"hidden"`
	Protected bool `json:"protected"`
	Tags []string `json:"tags"`
}
