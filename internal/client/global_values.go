package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type GlobalValuesService service

type GlobalValue struct {
	Id       string `json:"id"`
	SeriesId string `json:"series_id"`
	Version  int64  `json:"version"`

	CommonGlobalValueFields

	Key  string `json:"key"`
	Type string `json:"type"`

	IsDeprecated bool `json:"is_deprecated"`
}

type NewGlobalValue struct {
	CommonGlobalValueFields

	Key  string `json:"key"`
	Type string `json:"type"`
}

type UpdatedGlobalValue struct {
	SeriesId string `json:"-"`

	CommonGlobalValueFields

	IsDeprecated bool `json:"is_deprecated"`
}

type CommonGlobalValueFields struct {
	Name        string              `json:"name"`
	Description string              `json:"description,omitempty"`
	Options     []GlobalValueOption `json:"options"`
}

type GlobalValueOption struct {
	Key         string      `json:"key"`
	Label       string      `json:"label"`
	Description string      `json:"description,omitempty"`
	Value       interface{} `json:"value"`
}

type GlobalValuesQueryResponse struct {
	Page       int `json:"page,omitempty"`
	PageSize   int `json:"page_size,omitempty"`
	TotalPages int `json:"total_pages,omitempty"`
	TotalItems int `json:"total_items"`

	PageItems []GlobalValue `json:"page_items"`
}

func (s *GlobalValuesService) GetGlobalValueBySeriesId(ctx context.Context, seriesId string) (*GlobalValue, *http.Response, error) {
	path := fmt.Sprintf("%s/presets/series/%s", s.Client.BasePath, seriesId)
	body, resp, err := s.Client.Get(ctx, path, url.Values{}, new(GlobalValue))
	if err != nil {
		return nil, resp, err
	}
	return body.(*GlobalValue), resp, nil
}

func (s *GlobalValuesService) GetGlobalValueByKey(ctx context.Context, key string) (*GlobalValue, *http.Response, error) {
	query := url.Values{}
	query.Set("key", key)
	query.Set("page_size", "2")
	query.Set("sort_field", "series_id")

	path := fmt.Sprintf("%s/presets", s.Client.BasePath)
	body, resp, err := s.Client.Get(ctx, path, query, new(GlobalValuesQueryResponse))
	if err != nil {
		return nil, resp, err
	}

	globalValues := body.(*GlobalValuesQueryResponse).PageItems
	switch len(globalValues) {
	case 0:
		return nil, resp, nil
	case 1:
		return &globalValues[0], resp, nil
	default:
		return &globalValues[0], resp, fmt.Errorf("Found multiple global values with the provided key. Expected just one.")
	}
}

func (s *GlobalValuesService) CreateGlobalValue(ctx context.Context, newGlobalValue *NewGlobalValue) (*GlobalValue, *http.Response, error) {
	path := fmt.Sprintf("%s/presets", s.Client.BasePath)
	body, resp, err := s.Client.Post(ctx, path, newGlobalValue, new(GlobalValue))
	if err != nil {
		return nil, resp, err
	}
	return body.(*GlobalValue), resp, nil
}

func (s *GlobalValuesService) UpdateGlobalValue(ctx context.Context, updatedGlobalValue *UpdatedGlobalValue) (*GlobalValue, *http.Response, error) {
	path := fmt.Sprintf("%s/presets/series/%s", s.Client.BasePath, updatedGlobalValue.SeriesId)
	body, resp, err := s.Client.Put(ctx, path, updatedGlobalValue, new(GlobalValue))
	if err != nil {
		return nil, resp, err
	}
	return body.(*GlobalValue), resp, nil
}
