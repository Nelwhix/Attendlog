package responses

import "github.com/Nelwhix/Attendlog/models"

type Link struct {
	ID        string
	Title     string
	CreatedAt int64
}

func NewLink(link models.Link) Link {
	return Link{
		ID:        link.ID,
		Title:     link.Title,
		CreatedAt: link.CreatedAt.Unix(),
	}
}
