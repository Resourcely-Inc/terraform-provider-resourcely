package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type BlueprintTemplatesService service

type BlueprintTemplate struct {
	Id       string `json:"id"`
	SeriesId string `json:"series_id"`
	Version  int64  `json:"version"`
	Scope    string `json:"scope"`
	CommonBlueprintTemplateFields
	Provider string `json:"provider"`
}

type NewBlueprintTemplate struct {
	CommonBlueprintTemplateFields
	Provider string `json:"provider"`
}

type UpdatedBlueprintTemplate struct {
	SeriesId string `json:"-"`
	CommonBlueprintTemplateFields
}

type CommonBlueprintTemplateFields struct {
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Content     string   `json:"content"`
	Categories  []string `json:"categories"`
	Guidance    string   `json:"guidance"`
	Labels      []Label  `json:"labels"`
}

func (s *BlueprintTemplatesService) GetBlueprintTemplateBySeriesId(ctx context.Context, seriesId string) (*BlueprintTemplate, *http.Response, error) {
	path := fmt.Sprintf("%s/blueprint-templates/series/%s", s.Client.BasePath, seriesId)
	body, resp, err := s.Client.Get(ctx, path, url.Values{}, new(BlueprintTemplate))
	if err != nil {
		return nil, resp, err
	}
	return body.(*BlueprintTemplate), resp, nil
}

func (s *BlueprintTemplatesService) CreateBlueprintTemplate(ctx context.Context, newBlueprintTemplate *NewBlueprintTemplate) (*BlueprintTemplate, *http.Response, error) {
	path := fmt.Sprintf("%s/blueprint-templates", s.Client.BasePath)
	body, resp, err := s.Client.Post(ctx, path, newBlueprintTemplate, new(BlueprintTemplate))
	if err != nil {
		return nil, resp, err
	}
	return body.(*BlueprintTemplate), resp, nil
}

func (s *BlueprintTemplatesService) UpdateBlueprintTemplate(ctx context.Context, updatedBlueprintTemplate *UpdatedBlueprintTemplate) (*BlueprintTemplate, *http.Response, error) {
	path := fmt.Sprintf("%s/blueprint-templates/series/%s", s.Client.BasePath, updatedBlueprintTemplate.SeriesId)
	body, resp, err := s.Client.Put(ctx, path, updatedBlueprintTemplate, new(BlueprintTemplate))
	if err != nil {
		return nil, resp, err
	}
	return body.(*BlueprintTemplate), resp, nil
}

func (s *BlueprintTemplatesService) DeleteBlueprintTemplate(ctx context.Context, blueprintTemplateSeriesId string) (*http.Response, error) {
	path := fmt.Sprintf("%s/blueprint-templates/series/%s", s.Client.BasePath, blueprintTemplateSeriesId)
	return s.Client.Delete(ctx, path)
}
