package github

import (
	"github.com/google/go-github/v67/github"
)

type Release struct {
	Name        *string           `json:"name,omitempty"`
	TagName     *string           `json:"tag_name,omitempty"`
	Body        *string           `json:"body,omitempty"`
	Draft       *bool             `json:"draft,omitempty"`
	Prerelease  *bool             `json:"prerelease,omitempty"`
	CreatedAt   *github.Timestamp `json:"created_at,omitempty"`
	PublishedAt *github.Timestamp `json:"published_at,omitempty"`
	AuthorName  *string           `json:"author_name,omitempty"`
}
