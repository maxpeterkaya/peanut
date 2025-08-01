package helper

import (
	"github.com/google/go-github/v67/github"
	"github.com/rs/zerolog/log"
	"peanut/internal/cache"
	"strings"
)

func GetPlatformRelease(repo string, platform string) *github.ReleaseAsset {
	latestRelease, err := cache.FindLatest(repo)
	if err != nil {
		log.Error().Err(err).Str("repo", repo).Msg("Error finding latest release")
		return nil
	}
	assets := latestRelease.Assets

	for _, asset := range assets {
		if strings.HasSuffix(strings.ToLower(*asset.Name), strings.ToLower(platform)) {
			return asset
		}
	}

	return nil
}

type UserAgent struct {
	FirstPrediction  string
	SecondPrediction string
}

func ParseUserAgent(userAgent string) *UserAgent {
	ua := UserAgent{}

	if strings.Contains(userAgent, "win") || strings.Contains(userAgent, "windows") || strings.Contains(userAgent, "win32") {
		ua.FirstPrediction = "exe"
	} else if strings.Contains(userAgent, "debian") || strings.Contains(userAgent, "ubuntu") {
		ua.FirstPrediction = "deb"
		ua.SecondPrediction = "appimage"
	} else if strings.Contains(userAgent, "dmg") {
		ua.FirstPrediction = "dmg"
	} else if strings.Contains(userAgent, "fedora") {
		ua.FirstPrediction = "fedora"
	} else if strings.Contains(userAgent, "mac") || strings.Contains(userAgent, "macos") || strings.Contains(userAgent, "osx") || strings.Contains(userAgent, "macintosh") {
		ua.FirstPrediction = "darwin"
	} else {
		return nil
	}

	return &ua
}
