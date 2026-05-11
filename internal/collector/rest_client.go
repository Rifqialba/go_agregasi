package collector

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"time"
)

type RestClient struct {
	client *http.Client
}

func NewRestClient() *RestClient {
	return &RestClient{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *RestClient) Fetch(
	ctx context.Context,
	url string,
	lastFetchedAt time.Time,
) ([]byte, int, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		return nil, 0, err
	}

	if !lastFetchedAt.IsZero() {
		req.Header.Set(
			"If-Modified-Since",
			lastFetchedAt.UTC().Format(http.TimeFormat),
		)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode

	switch {
	case statusCode == http.StatusNotModified:
		log.Println("no new data")

		return nil, statusCode, nil

	case statusCode >= 400 && statusCode < 500:
		return nil, statusCode, errors.New("client error")

	case statusCode >= 500:
		return nil, statusCode, errors.New("server error")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, statusCode, err
	}

	return body, statusCode, nil
}