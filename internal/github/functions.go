package github

import (
	"context"
	"github.com/google/go-github/v67/github"
	"peanut/internal/config"
)

var (
	GHClient *github.Client
	GHCtx    = context.Background()
)

func Init() error {
	if len(config.Config.GHToken) > 0 {
		GHClient = github.NewClient(nil).WithAuthToken(config.Config.GHToken)
	} else {
		GHClient = github.NewClient(nil)
	}

	return nil
}
