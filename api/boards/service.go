package boards

import (
	"encoding/json"
	"errors"
)

var brds = []Board{
	{
		BoardInfo: BoardInfo{
			Slug: "test-board1",
			Name: "Test Board Alpha",
		},
		Content: nil,
	},
	{
		BoardInfo: BoardInfo{
			Slug: "test-board2",
			Name: "Test Board Beta",
		},
		Content: nil,
	},
}

type Boarder interface {
	getInfos(orgSlug string) ([]BoardInfo, error)
	getBoard(orgSlug string, slug string) (*Board, error)
	storeBoard(orgSlug string, slug string, content json.RawMessage) (*Board, error)
}

type service struct {}

func (s *service) getInfos(orgSlug string) ([]BoardInfo, error) {
	var infos []BoardInfo
	for _, b := range brds {
		infos = append(infos, b.BoardInfo)
	}
	return infos, nil
}

func (s *service) getBoard(orgSlug string, slug string) (*Board, error) {
	for _, b := range brds {
		if b.Slug == slug {
			return &b, nil
		}
	}
	return nil, errors.New("board not found")
}

func (s *service) storeBoard(orgSlug string, slug string, content json.RawMessage) (*Board, error) {
	for i := range brds {
		if brds[i].Slug == slug {
			brds[i].Content = content
			return &brds[i], nil
		}
	}
	return nil, errors.New("board not found")
}




