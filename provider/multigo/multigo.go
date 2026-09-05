package multigo

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"net/url"
	"strconv"
	"time"
	"uuid"
)

const (
	DefaultHost             = "multipass-api.golabs.vn"
	DefaultScheme           = "https"
	DefaultCredentialSuffix = ":Bus2019M@p_"

	PathGeolocation = "/v2/public/busmap/route_bus_gps"
)

var (
	DefaultURL = url.URL{
		Host:   DefaultHost,
		Scheme: DefaultScheme,
	}
)

type Client interface {
	ListGeolocations(context.Context, ListGeolocationParams) ([]Geolocation, error)
}

type client struct {
	httpClient *http.Client
}

type option func(c *client)

func NewClient(options ...option) Client {
	c := &client{
		httpClient: http.DefaultClient,
	}

	for _, option := range options {
		option(c)
	}

	return c
}

func WithHTTPClient(httpClient *http.Client) option {
	return func(c *client) {
		c.httpClient = httpClient
	}
}

func injectCredentials(header *http.Header) error {
	epoch := time.Now().UnixMilli()
	epochStr := strconv.FormatInt(epoch, 10)
	secret := epochStr + DefaultCredentialSuffix
	proof := md5.Sum([]byte(secret))

	header.Add("device-id", uuid.New().String())
	header.Add("epoch", epochStr)
	header.Add("proof", hex.EncodeToString(proof[:]))
	header.Add("client-version", "ios|26")

	return nil
}
