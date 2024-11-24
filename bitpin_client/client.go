package bitpin_client

import (
	"bitpin-tui/utils"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var (
	BaseUrl = "https://api.bitpin.market"
	Version = "v1"
)

type Client struct {
	HttpClient *http.Client

	AccessToken  string
	RefreshToken string

	ApiKey    string
	SecretKey string

	AutoRefresh bool
}

func NewClient(opts ClientOptions) (*Client, error) {
	if opts.BaseUrl != "" {
		BaseUrl = opts.BaseUrl
	}

	client := &Client{
		AutoRefresh: opts.AutoRefresh,
	}

	if opts.HttpClient != nil {
		client.HttpClient = opts.HttpClient
	} else {
		client.HttpClient = &http.Client{}
	}

	client.AccessToken = opts.AccessToken
	client.RefreshToken = opts.RefreshToken
	client.ApiKey = opts.ApiKey
	client.SecretKey = opts.SecretKey

	if err := client.handleAutoRefresh(); err != nil {
		return nil, fmt.Errorf("error handling auto refresh: %v", err)
	}

	if opts.ApiKey != "" && opts.SecretKey != "" {
		if _, err := client.Authenticate(opts.ApiKey, opts.SecretKey); err != nil {
			return nil, fmt.Errorf("error authenticating: %v", err)
		}
	}

	return client, nil
}

func assertAuth(client *Client) error {
	if client.AccessToken == "" {
		return errors.New("access token is empty")
	}
	if client.RefreshToken == "" {
		return errors.New("refresh token is empty")
	}
	return nil
}

func (c *Client) createApiURI(endpoint string, version string) string {
	if version == "" {
		version = Version
	}
	return fmt.Sprintf("%s/api/%s%s", BaseUrl, Version, endpoint)
}

func (c *Client) handleAutoRefresh() error {
	if c.AccessToken != "" {
		decoded, err := utils.DecodeJWT(c.AccessToken)
		if err != nil {
			return fmt.Errorf("error decoding access token: %v", err)
		}
		if decoded.IsExpired() {
			fmt.Println("Access token expired, refreshing...")
			err = c.RefreshAccessToken()
			if err != nil {
				return fmt.Errorf("error refreshing access token: %v", err)
			}
		}
	}

	if c.RefreshToken != "" {
		decoded, err := utils.DecodeJWT(c.RefreshToken)
		if err != nil {
			return fmt.Errorf("error decoding refresh token: %v", err)
		}

		if decoded.IsExpired() {
			fmt.Println("Refresh token expired, re-authenticating...")
			if c.ApiKey == "" || c.SecretKey == "" {
				return errors.New("API key and/or secret key are empty")
			}

			_, err = c.Authenticate(c.ApiKey, c.SecretKey)
			if err != nil {
				return fmt.Errorf("error re-authenticating: %v", err)
			}
		}
	}

	return nil
}

func (c *Client) Request(method string, url string, auth bool, body interface{}, result interface{}) error {
	var reqBody []byte
	var err error

	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return fmt.Errorf("error marshaling Request body: %v", err)
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("error creating Request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	if auth {
		if c.AutoRefresh {
			if err := c.handleAutoRefresh(); err != nil {
				return err
			}
		}

		if err := assertAuth(c); err != nil {
			return err
		}

		req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending Request: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			fmt.Println("error closing response body:", err)
		}
	}(resp.Body)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %v", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &APIError{
			StatusCode: resp.StatusCode,
			Message:    string(respBody),
		}
	}

	if result != nil {
		err = json.Unmarshal(respBody, result)
		if err != nil {
			return fmt.Errorf("error unmarshaling response: %v", err)
		}
	}

	return nil
}

func (c *Client) ApiRequest(method, endpoint string, version string, auth bool, body interface{}, result interface{}) error {
	url := c.createApiURI(endpoint, version)
	return c.Request(method, url, auth, body, result)
}

func (c *Client) Authenticate(apiKey, secretKey string) (*AuthResponse, error) {
	if apiKey == "" || secretKey == "" {
		return nil, errors.New("API key and/or secret key are empty")
	}

	reqBody := map[string]string{
		"api_key":    apiKey,
		"secret_key": secretKey,
	}

	var authResponse AuthResponse
	err := c.ApiRequest("POST", "/usr/authenticate/", Version, false, reqBody, &authResponse)

	if err != nil {
		// Check for specific API errors here
		var apiErr *APIError
		if errors.As(err, &apiErr) {
			switch apiErr.StatusCode {
			case 401:
				return nil, fmt.Errorf("authentication failed: invalid API key or secret key")
			case 429:
				return nil, fmt.Errorf("authentication failed: rate limit exceeded")
			default:
				return nil, fmt.Errorf("authentication failed: %v", apiErr)
			}
		}
		return nil, fmt.Errorf("authentication failed: %v", err)
	}

	// Update the bitpin_client's tokens with the newly received ones
	c.AccessToken = authResponse.Access
	c.RefreshToken = authResponse.Refresh

	return &authResponse, nil
}

func (c *Client) RefreshAccessToken() error {
	reqBody := map[string]string{
		"refresh": c.RefreshToken,
	}

	var refreshResponse RefreshTokenResponse
	err := c.ApiRequest("POST", "/usr/refresh_token/", Version, false, reqBody, &refreshResponse)
	if err != nil {
		return err
	}

	// Update the bitpin_client's access token with the newly received one
	c.AccessToken = refreshResponse.Access

	return nil
}

func (c *Client) GetOrderBook(symbol string) (*OrderBook, error) {
	var orderBook *OrderBook
	err := c.ApiRequest("GET", fmt.Sprintf("/mth/orderbook/%s/", symbol), Version, false, nil, &orderBook)
	if err != nil {
		return nil, fmt.Errorf("error fetching order book: %v", err)
	}
	return orderBook, nil
}

func (c *Client) GetTickers() ([]*Ticker, error) {
	var tickers []*Ticker
	err := c.ApiRequest("GET", "/mkt/tickers/", Version, false, nil, &tickers)
	if err != nil {
		return nil, fmt.Errorf("error fetching tickers: %v", err)
	}
	return tickers, nil
}

func (c *Client) GetRecentTrades(symbol string) ([]*Trade, error) {
	var trades []*Trade
	err := c.ApiRequest("GET", fmt.Sprintf("/mth/matches/%s/", symbol), Version, false, nil, &trades)
	if err != nil {
		return nil, fmt.Errorf("error fetching recent trades: %v", err)
	}
	return trades, nil
}

func (c *Client) CancelOrder(orderId int) error {
	err := c.ApiRequest("DELETE", fmt.Sprintf("/odr/orders/%d/", orderId), Version, true, nil, nil)
	if err != nil {
		return fmt.Errorf("error canceling order: %v", err)
	}
	return nil
}

func (c *Client) CreateOrder(params CreateOrderParams) (*OrderStatus, error) {
	var orderStatus *OrderStatus
	err := c.ApiRequest("POST", "/odr/orders/", Version, true, params, &orderStatus)
	if err != nil {
		return nil, fmt.Errorf("error creating order: %v", err)
	}
	return orderStatus, nil
}

func (c *Client) GetOpenOrders(params GetOrdersParams) ([]*OrderStatus, error) {
	var orders []*OrderStatus
	err := c.ApiRequest("GET", "/odr/orders/?"+params.AsURLParams(), Version, true, nil, &orders)
	if err != nil {
		return nil, fmt.Errorf("error fetching open orders: %v", err)
	}
	return orders, nil
}

func (c *Client) GetOrderFills() (*Fills, error) {
	var fills *Fills
	err := c.ApiRequest("GET", "/odr/fills/", Version, true, nil, &fills)
	if err != nil {
		return nil, fmt.Errorf("error fetching order fills: %v", err)
	}
	return fills, nil
}

func (c *Client) GetBalances(params GetBalancesParams) ([]Balance, error) {
	var balances []Balance
	err := c.ApiRequest("GET", "/wlt/wallets/"+params.AsURLParams(), Version, true, nil, &balances)
	if err != nil {
		return nil, fmt.Errorf("error fetching balances: %v", err)
	}
	return balances, nil
}

func (c *Client) GetUserInfo() (*UserInfo, error) {
	var userInfo *UserInfo
	err := c.Request("GET", "https://api.bitpin.org/v3/usr/info/", true, nil, &userInfo)
	if err != nil {
		return &UserInfo{}, fmt.Errorf("error fetching user info: %v", err)
	}
	return userInfo, nil
}
