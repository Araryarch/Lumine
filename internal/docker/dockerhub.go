package docker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
)

type DockerHubResponse struct {
	Results []struct {
		Name string `json:"name"`
	} `json:"results"`
}

// FetchDockerHubTags fetches available tags for a Docker image from Docker Hub
func FetchDockerHubTags(imageName string) ([]string, error) {
	url := fmt.Sprintf("https://registry.hub.docker.com/v2/repositories/library/%s/tags?page_size=50", imageName)
	
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to fetch tags: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result DockerHubResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	var tags []string
	for _, tag := range result.Results {
		tags = append(tags, tag.Name)
	}

	// Sort tags
	sort.Strings(tags)

	return tags, nil
}
