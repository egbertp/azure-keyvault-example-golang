package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/caarlos0/env"
)

var (
	CommitHash string
	Version    string
)

type Config struct {
	ClientID     string `env:"AZURE_CLIENT_ID,required"`
	ClientSecret string `env:"AZURE_CLIENT_SECRET,required"`
	TenantID     string `env:"AZURE_TENANT_ID,required"`
	BaseURL      string `env:"VAULT_BASE_URL,required"`
	SecretName   string `env:"SECRET_NAME,required"`
}

type TokenResponse struct {
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
}

type KeyvaultSecretResponse struct {
	Value string `json:"value"`
	Id    string `json:"id"`
	// attributes
}

type Client struct {
	Token string
	Cfg   Config
}

func main() {
	fmt.Printf("Starting keyvault-get-secret \t\tversion: %s, \t\tcommit hash: %s\n\n", Version, CommitHash)
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalln(err)
	}

	// fmt.Printf("%+v\n", cfg)

	authenticatedClient, err := authenticate(cfg)
	if err != nil {
		log.Fatalln(err)
	}

	secretResponse, err := authenticatedClient.getSecret(cfg.BaseURL, cfg.SecretName)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Secret value is: \t\t%v\n\n", secretResponse.Value)

}

func authenticate(cfg Config) (*Client, error) {
	const authenticationScope = "https://vault.azure.net/.default"

	var keyVaultransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
	}

	var keyVaultClient = &http.Client{
		Timeout:   time.Second * 5,
		Transport: keyVaultransport,
	}

	authenticationURL := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", cfg.TenantID)

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", cfg.ClientID)
	data.Set("client_secret", cfg.ClientSecret)
	data.Set("scope", authenticationScope)

	resp, err := keyVaultClient.Post(authenticationURL, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// log.Println(string(body))

	var tr TokenResponse
	errUnMarshal := json.Unmarshal(body, &tr)
	if errUnMarshal != nil {
		return nil, errUnMarshal
	}

	return &Client{
		Token: tr.AccessToken,
		Cfg:   cfg,
	}, nil
}

func (c *Client) getSecret(vaultBaseURL, secretName string) (*KeyvaultSecretResponse, error) {
	// Create a Bearer string by appending string access token
	var bearer = "Bearer " + c.Token
	log.Printf("Authorization: \t\t%s\n\n", bearer)

	url := fmt.Sprintf("%s/secrets/%s/?api-version=7.0", vaultBaseURL, secretName)

	// Create a new request using http
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// add authorization header to the req
	req.Header.Add("Authorization", bearer)

	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var sr KeyvaultSecretResponse
	errUnMarshal := json.Unmarshal([]byte(body), &sr)
	if errUnMarshal != nil {
		return nil, errUnMarshal
	}

	return &sr, nil
}
