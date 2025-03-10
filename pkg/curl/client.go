package curl

import (
	"context"
	"crypto/tls"
	"net/http"

	"golang.org/x/oauth2"

	"github.com/marcofranssen/terraform-provider-curl/pkg/log"
	"github.com/marcofranssen/terraform-provider-curl/pkg/transport"
)

type HttpClientOptions struct {
	token      string
	disabletls bool
}

type HttpClient struct {
	httpClient *http.Client
}

func NewClient(ctx context.Context, opts HttpClientOptions) (*HttpClient, error) {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: opts.token})
	c := oauth2.NewClient(ctx, ts)

	if opts.disabletls {
		c.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	c.Transport = transport.TeeRoundTripper{
		RoundTripper: c.Transport,
		Writer:       log.NewTFLogger(ctx),
	}

	return &HttpClient{httpClient: c}, nil
}
