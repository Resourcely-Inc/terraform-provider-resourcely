package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type SystemService service

type SystemHealth struct {
	Status string `json:"status"`
}

func (s *SystemService) GetHealth(ctx context.Context) (*SystemHealth, *http.Response, error) {
	path := fmt.Sprintf("%s/system/health", s.Client.BasePath)
	body, resp, err := s.Client.Get(ctx, path, url.Values{}, new(SystemHealth))
	if err != nil {
		return nil, resp, err
	}
	return body.(*SystemHealth), resp, nil
}
