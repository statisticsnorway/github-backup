package git

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/tomnomnom/linkheader"
)

type Repo struct {
	FullName string `json:"full_name"`
}

var client = http.Client{}

func ReposFor(org, authToken string) ([]Repo, error) {
	urlRaw := fmt.Sprintf("https://api.github.com/orgs/%s/repos", org)
	var allRepos []Repo
	for urlRaw != "" {
		fmt.Printf("Retrieving repos from %s\n", urlRaw)
		var reposPart []Repo
		res, err := getRequest(urlRaw, authToken)
		if res.StatusCode != 200 {
			return nil, fmt.Errorf("got a %d from GitHub", res.StatusCode)
		}
		if err != nil {
			return nil, err
		}
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(body, &reposPart)
		if err != nil {
			return nil, err
		}
		allRepos = append(allRepos, reposPart...)
		linkHeader := res.Header.Get("Link")
		urlRaw = nextUrl(linkHeader)
	}
	return allRepos, nil
}

func getRequest(rawUrl, authToken string) (*http.Response, error) {
	u, err := url.Parse(rawUrl)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header = http.Header{
		"Accept":        {"application/vnd.github.v3+json"},
		"User-Agent":    {"NAV IT Backup"},
		"Authorization": {fmt.Sprintf("Bearer %s", authToken)},
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func nextUrl(linkHeader string) string {
	links := linkheader.Parse(linkHeader)
	for _, link := range links {
		if link.Rel == "next" {
			return link.URL
		}
	}
	return ""
}
