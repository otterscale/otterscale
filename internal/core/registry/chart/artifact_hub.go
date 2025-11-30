package chart

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	artifactHubSearchURL       = "https://artifacthub.io/api/v1/packages/search?kind=0&verified_publisher=true&official=true&sort=stars&limit=60&offset=0"
	artifactHubPackageTemplate = "https://artifacthub.io/api/v1/packages/helm/%s/%s"
)

type SearchResponse struct {
	Packages []Package `json:"packages"`
}

type Package struct {
	Name       string `json:"name"`
	Repository struct {
		Name string `json:"name"`
	} `json:"repository"`
}

type PackageContent struct {
	URL string `json:"content_url"`
}

func fetchPackages(ctx context.Context, url string) ([]Package, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to fetch package summaries: status %s, body: %s", resp.Status, string(bodyBytes))
	}

	var searchResponse SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResponse); err != nil {
		return nil, err
	}

	return searchResponse.Packages, nil
}

func fetchContentURL(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to fetch content url: status %s, body: %s", resp.Status, string(bodyBytes))
	}

	var content PackageContent
	if err := json.NewDecoder(resp.Body).Decode(&content); err != nil {
		return "", err
	}

	return content.URL, nil
}
