package helper

import (
	"fmt"
	"io"
	"net/http"
	"peanut/internal/repository"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
)

func ProxyDownload(repo repository.Repository, asset repository.Asset, w http.ResponseWriter, r *http.Request) {
	url := asset.ApiUrl.String
	client := &http.Client{
		Timeout: time.Second * 100,
	}

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/octet-stream")
	if len(repo.Token.String) > 0 {
		req.Header.Set("Authorization", "Bearer "+repo.Token.String)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msgf("Error downloading %s", url)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, asset.Name.String))
	w.Header().Set("Content-Length", strconv.FormatInt(resp.ContentLength, 10))
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		log.Error().Err(err).Msgf("Error copying %s/%s/%s", repo.Owner.String, repo.Name.String, asset.Name.String)
		return
	}
}
