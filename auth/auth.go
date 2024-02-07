package auth

import (
	"fmt"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
	"google.golang.org/grpc/metadata"
)

const (
	address string = "googleads.googleapis.com:443"
)

// NewCredentials creates an oauth2 credentials config from a valid client id and secret
func NewCredentials(clientID, clientSecret string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     google.Endpoint,
	}
}

// NewPartialToken creates a partial oauth2 token from a valid refresh token
func NewPartialToken(refreshToken string) *oauth2.Token {
	return &oauth2.Token{
		RefreshToken: refreshToken,
	}
}

// RefreshToken creates a new oauth2 token from an existing partial/expired token and valid credentials
func RefreshToken(creds *oauth2.Config, token *oauth2.Token, ctx ...context.Context) (*oauth2.Token, error) {
	defaultCtx := context.Background()
	if ctx != nil {
		defaultCtx = ctx[0]
	}
	tokenSource := creds.TokenSource(defaultCtx, token)
	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, err
	}
	return newToken, nil
}

// NewGrpcConnection creates a new authenticated grpc client and context with authentication headers
func NewGrpcConnection(token *oauth2.Token, developerToken string, loginCustomerID string) (*grpc.ClientConn, context.Context, error) {
	headers := createHeaders(token.AccessToken, developerToken, loginCustomerID)
	ctx := metadata.NewOutgoingContext(context.Background(), headers)

	creds := oauth.NewOauthAccess(token)
	transportCreds := grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, ""))
	conn, err := grpc.Dial(address, transportCreds, grpc.WithPerRPCCredentials(creds))
	if err != nil {
		return nil, nil, err
	}
	return conn, ctx, nil
}

// TODO: Add support for login-customer-id
// https://developers.google.com/google-ads/api/docs/concepts/call-structure#login-customer-id
func createHeaders(accessToken string, developerToken string, loginCustomerID string) metadata.MD {
	return metadata.Pairs(
		"Authorization", fmt.Sprintf("Bearer %s", accessToken),
		"developer-token", developerToken,
		"login-customer-id", loginCustomerID,
	)
}
