package helper

import (
	"fmt"
	"github.com/google/go-github/v67/github"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"peanut/internal/config"
	"strconv"
	"time"
)

func ProxyDownload(asset github.ReleaseAsset, token string, w http.ResponseWriter, r *http.Request) {
	url := *asset.URL
	client := &http.Client{
		Timeout: time.Second * 100,
	}

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/octet-stream")
	if len(config.Config.Github.Token) > 0 {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msgf("Error downloading %s", url)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, *asset.Name))
	w.Header().Set("Content-Length", strconv.FormatInt(resp.ContentLength, 10))
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		log.Error().Err(err).Msgf("Error copying %s", url)
		return
	}
}
