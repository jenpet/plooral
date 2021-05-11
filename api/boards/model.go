package boards

import "encoding/json"

type Board struct {
	BoardInfo
	Content json.RawMessage `json:"content"`
}

type BoardInfo struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
}
