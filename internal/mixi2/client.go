package mixi2

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/koutya0akari/arxiv-plug-in-for-mixi2/internal/config"
	"github.com/mixigroup/mixi2-application-sdk-go/auth"
	applicationapiv1 "github.com/mixigroup/mixi2-application-sdk-go/gen/go/social/mixi/application/service/application_api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Poster interface {
	Post(context.Context, string) error
	Close() error
}

type Client struct {
	authenticator auth.Authenticator
	conn          *grpc.ClientConn
	client        applicationapiv1.ApplicationServiceClient
	communityID   string
}

func New(creds config.Credentials) (*Client, error) {
	authenticator, err := auth.NewAuthenticator(creds.ClientID, creds.ClientSecret, creds.TokenURL)
	if err != nil {
		return nil, fmt.Errorf("create authenticator: %w", err)
	}
	conn, err := grpc.NewClient(
		creds.APIAddress,
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{MinVersion: tls.VersionTLS12})),
	)
	if err != nil {
		return nil, fmt.Errorf("connect api: %w", err)
	}
	return &Client{
		authenticator: authenticator,
		conn:          conn,
		client:        applicationapiv1.NewApplicationServiceClient(conn),
		communityID:   creds.CommunityID,
	}, nil
}

func (c *Client) Post(ctx context.Context, text string) error {
	authCtx, err := c.authenticator.AuthorizedContext(ctx)
	if err != nil {
		return fmt.Errorf("authorize: %w", err)
	}
	_, err = c.client.CreatePost(authCtx, &applicationapiv1.CreatePostRequest{
		Text:        text,
		CommunityId: &c.communityID,
	})
	if err != nil {
		return fmt.Errorf("create post: %w", err)
	}
	return nil
}

func (c *Client) Close() error {
	if c.conn == nil {
		return nil
	}
	return c.conn.Close()
}
