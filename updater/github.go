package updater

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type GithubRelease struct {
	TagName     string               `json:"tag_name"`
	Name        string               `json:"name"`
	Draft       bool                 `json:"draft"`
	Prerelease  bool                 `json:"prerelease"`
	CreatedAt   time.Time            `json:"created_at"`
	PublishedAt time.Time            `json:"published_at"`
	Assets      []GithubReleaseAsset `json:"assets"`
	TarballURL  string               `json:"tarball_url"`
	ZipballURL  string               `json:"zipball_url"`
	Body        string               `json:"body"`
}

type GithubReleaseAsset struct {
	Name               string    `json:"name"`
	ContentType        string    `json:"content_type"`
	State              string    `json:"state"`
	Size               int       `json:"size"`
	DownloadCount      int       `json:"download_count"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	BrowserDownloadURL string    `json:"browser_download_url"`
}

func GetReleases(repo string, binary string, draft bool, prerelease bool) ([]*ReleaseVersion, error) {

	// Trim github url prefix and generate github releases url.
	repo = strings.TrimPrefix(strings.TrimPrefix(repo, "https://"), "github.com/")
	url := "https://api.github.com/repos/" + repo + "/releases"

	// Get github releases.
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var ghReleases []GithubRelease
	err = json.Unmarshal(body, &ghReleases)
	if err != nil {
		return nil, err
	}

	// Filter github releases, only return releases with binary, hash and signature.
	var releases []*ReleaseVersion
	for _, ghRelease := range ghReleases {
		if ghRelease.Draft && !draft {
			continue
		}
		if ghRelease.Prerelease && !prerelease {
			continue
		}
		release := &ReleaseVersion{
			name: ghRelease.Name,
			tag:  ghRelease.TagName,
		}
		for _, asset := range ghRelease.Assets {
			if asset.Name == binary {
				release.binaryURL = asset.BrowserDownloadURL
			}
			if asset.Name == binary+".sha256" {
				release.checksumURL = asset.BrowserDownloadURL
			}
			if asset.Name == binary+".sig" {
				release.signatureURL = asset.BrowserDownloadURL
			}
		}
		if release.binaryURL != "" && release.checksumURL != "" && release.signatureURL != "" {
			releases = append(releases, release)
		}
	}
	return releases, nil
}
