package imageresolver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

type registryTagLister struct {
	client *http.Client
}

func newRegistryTagLister() TagLister {
	return registryTagLister{
		client: &http.Client{Timeout: 15 * time.Second},
	}
}

func (l registryTagLister) ListTags(ctx context.Context, repository string) ([]string, error) {
	registryHost, repositoryPath, err := parseRepository(repository)
	if err != nil {
		return nil, err
	}

	tags, challenge, err := l.fetchTags(ctx, registryHost, repositoryPath, "")
	if err == nil {
		return tags, nil
	}
	if challenge == "" {
		return nil, err
	}

	token, tokenErr := l.fetchBearerToken(ctx, challenge, repositoryPath)
	if tokenErr != nil {
		return nil, fmt.Errorf("%w: %w", err, tokenErr)
	}

	tags, _, err = l.fetchTags(ctx, registryHost, repositoryPath, token)
	if err != nil {
		return nil, err
	}
	return tags, nil
}

func (l registryTagLister) fetchTags(ctx context.Context, registryHost string, repositoryPath string, token string) ([]string, string, error) {
	endpoint := fmt.Sprintf("https://%s/v2/%s/tags/list", registryHost, repositoryPath)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, "", fmt.Errorf("failed to build registry request: %w", err)
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := l.client.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("failed to query registry tags for %s: %w", repositoryPath, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		challenge := resp.Header.Get("Www-Authenticate")
		if challenge == "" {
			challenge = resp.Header.Get("WWW-Authenticate")
		}
		return nil, challenge, fmt.Errorf("registry authentication required for %s", repositoryPath)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("registry returned HTTP %d for %s", resp.StatusCode, repositoryPath)
	}

	var body struct {
		Tags []string `json:"tags"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, "", fmt.Errorf("failed to decode registry tag response: %w", err)
	}

	if len(body.Tags) == 0 {
		return []string{}, "", nil
	}
	return body.Tags, "", nil
}

func (l registryTagLister) fetchBearerToken(ctx context.Context, challenge string, repositoryPath string) (string, error) {
	realm, service, scope, err := parseBearerChallenge(challenge)
	if err != nil {
		return "", err
	}
	if scope == "" {
		scope = "repository:" + repositoryPath + ":pull"
	}

	tokenURL, err := url.Parse(realm)
	if err != nil {
		return "", fmt.Errorf("invalid auth realm: %w", err)
	}
	query := tokenURL.Query()
	if service != "" {
		query.Set("service", service)
	}
	query.Set("scope", scope)
	tokenURL.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, tokenURL.String(), nil)
	if err != nil {
		return "", fmt.Errorf("failed to build token request: %w", err)
	}
	resp, err := l.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch registry token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token endpoint returned HTTP %d", resp.StatusCode)
	}

	var body struct {
		Token       string `json:"token"`
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "", fmt.Errorf("failed to decode token response: %w", err)
	}

	if body.Token != "" {
		return body.Token, nil
	}
	if body.AccessToken != "" {
		return body.AccessToken, nil
	}
	return "", fmt.Errorf("token response missing token")
}

var bearerChallengeParamPattern = regexp.MustCompile(`([a-zA-Z_]+)="([^"]*)"`)

func parseBearerChallenge(challenge string) (string, string, string, error) {
	if challenge == "" {
		return "", "", "", fmt.Errorf("missing WWW-Authenticate challenge")
	}
	prefix, rest, found := strings.Cut(challenge, " ")
	if !found || !strings.EqualFold(prefix, "Bearer") {
		return "", "", "", fmt.Errorf("unsupported auth challenge")
	}

	matches := bearerChallengeParamPattern.FindAllStringSubmatch(rest, -1)
	values := make(map[string]string, len(matches))
	for _, m := range matches {
		if len(m) == 3 {
			values[strings.ToLower(m[1])] = m[2]
		}
	}

	realm := values["realm"]
	if realm == "" {
		return "", "", "", fmt.Errorf("auth challenge missing realm")
	}
	return realm, values["service"], values["scope"], nil
}

func parseRepository(repository string) (string, string, error) {
	if repository == "" {
		return "", "", fmt.Errorf("repository is required")
	}

	parts := strings.Split(repository, "/")
	if len(parts) == 0 {
		return "", "", fmt.Errorf("repository is required")
	}

	if isRegistryHost(parts[0]) {
		if len(parts) < 2 {
			return "", "", fmt.Errorf("repository path is required")
		}
		return parts[0], strings.Join(parts[1:], "/"), nil
	}

	repoPath := repository
	if !strings.Contains(repoPath, "/") {
		repoPath = "library/" + repoPath
	}
	return "registry-1.docker.io", repoPath, nil
}

func isRegistryHost(part string) bool {
	return part == "localhost" || strings.Contains(part, ".") || strings.Contains(part, ":")
}
