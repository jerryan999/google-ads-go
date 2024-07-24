package ads

import (
	"encoding/json"
	"io/ioutil"

	"github.com/jerryan999/google-ads-go/v17/auth"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
)

const (
	defaultVersion string = "v0"
)

type GoogleAdsClient struct {
	credentials     *oauth2.Config
	token           *oauth2.Token
	conn            *grpc.ClientConn
	developerToken  string
	loginCustomerID string
	ctx             context.Context
}

type GoogleAdsClientParams struct {
	ClientID        string
	ClientSecret    string
	DeveloperToken  string
	RefreshToken    string
	LoginCustomerID string
}

type googleAdsStorageParams struct {
	ClientID        string `json:"client_id"`
	ClientSecret    string `json:"client_secret"`
	RefreshToken    string `json:"refresh_token"`
	DeveloperToken  string `json:"developer_token"`
	LoginCustomerID string `json:"login_customer_id",omitempty`
}

// NewClient creates a new client with specified credential params
func NewClient(params *GoogleAdsClientParams) (*GoogleAdsClient, error) {
	credentials := auth.NewCredentials(params.ClientID, params.ClientSecret)
	initialToken := auth.NewPartialToken(params.RefreshToken)

	c := &GoogleAdsClient{
		credentials:     credentials,
		token:           initialToken,
		developerToken:  params.DeveloperToken,
		loginCustomerID: params.LoginCustomerID,
	}

	newToken, err := auth.RefreshToken(c.credentials, c.token)
	if err != nil {
		return nil, err
	}
	c.token = newToken

	conn, ctx, err := auth.NewGrpcConnection(c.token, c.developerToken, c.loginCustomerID)
	if err != nil {
		return nil, err
	}
	c.conn = conn
	c.ctx = ctx

	return c, nil
}

// NewClientFromStorage creates a new client instance from specified "google-ads.json" file
func NewClientFromStorage(filepath string) (*GoogleAdsClient, error) {
	params, err := ReadCredentialsFile(filepath)
	if err != nil {
		return nil, err
	}
	client, err := NewClient(params)
	return client, err
}

// ReadCredentialsFile reads a credentials JSON file and returns the exported config
func ReadCredentialsFile(filepath string) (*GoogleAdsClientParams, error) {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	var g googleAdsStorageParams
	if err = json.Unmarshal(file, &g); err != nil {
		return nil, err
	}

	return &GoogleAdsClientParams{
		ClientID:        g.ClientID,
		ClientSecret:    g.ClientSecret,
		RefreshToken:    g.RefreshToken,
		DeveloperToken:  g.DeveloperToken,
		LoginCustomerID: g.LoginCustomerID,
	}, nil
}

// RefreshTokenIfNeeded checks and refreshes the access token if it is not valid
func (g *GoogleAdsClient) RefreshTokenIfNeeded() error {
	if !g.token.Valid() {
		newToken, err := auth.RefreshToken(g.credentials, g.token)
		if err != nil {
			return err
		}
		g.token = newToken

		// Refresh gRPC connection with new token
		conn, ctx, err := auth.NewGrpcConnection(g.token, g.developerToken, g.loginCustomerID)
		if err != nil {
			return err
		}
		g.conn = conn
		g.ctx = ctx
	}
	return nil
}

// Conn returns a pointer to the clients gRPC connection
func (g *GoogleAdsClient) Conn() *grpc.ClientConn {
	return g.conn
}

// Context returns the context of the client
func (g *GoogleAdsClient) Context() context.Context {
	return g.ctx
}

// TokenIsValid returns a bool indicating if the generated access token is valid
func (g *GoogleAdsClient) TokenIsValid() bool {
	return g.token.Valid()
}
