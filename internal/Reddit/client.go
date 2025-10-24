package reddit

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

type RedditClient struct {
	Username string
	Password string
	ID string
	Secret string
	UserAgent string
	Token string
	TokenExpiry time.Time
}
type TokenResp struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

func NewClient() *RedditClient {
	_ = godotenv.Load()
	return &RedditClient{
		Username: os.Getenv("REDDIT_USERNAME"),
		Password: os.Getenv("REDDIT_PASSWORD"),
		ID: os.Getenv("REDDIT_CLIENT_ID"),
		Secret: os.Getenv("REDDIT_CLIENT_SECRET"),
		UserAgent: os.Getenv("REDDIT_USER_AGENT"),
	}
}

func (client *RedditClient) getToken() error {
	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("username", client.Username)
	data.Set("password", client.Password)
	data.Set("scope", "read")

	req, _ := http.NewRequest("POST", "https://www.reddit.com/api/v1/access_token", strings.NewReader(data.Encode()))
	req.Header.Set("User-Agent", client.UserAgent)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(client.ID+":"+client.Secret)))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var tok TokenResp
	json.NewDecoder(resp.Body).Decode(&tok)
	
	client.Token = tok.AccessToken
	client.TokenExpiry = time.Now().Add(time.Duration(tok.ExpiresIn-300) * time.Second)

	return nil
}

func (client *RedditClient) ensureToken() error {
	if client.Token == "" || time.Now().After(client.TokenExpiry) {
		return client.getToken()
	}
	return nil
}

func (client *RedditClient) GetReddit(path string) (*RedditAPIResponse, error) {
	if err := client.ensureToken(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, "https://oauth.reddit.com"+path, nil)
	if err != nil {
		log.WithFields(log.Fields{
			"path": path,
		}).Error("Failed to fetch Reddit data")
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+client.Token)
	req.Header.Set("User-Agent", client.UserAgent)
	
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp RedditAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	return &apiResp, nil
}

