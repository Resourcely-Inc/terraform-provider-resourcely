package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type BlueprintsService service

type Blueprint struct {
	Id       string `json:"id"`
	SeriesId string `json:"series_id"`
	Version  int64  `json:"version"`
	Scope    string `json:"scope"`
	CommonBlueprintFields
	Provider    string `json:"provider"`
	IsPublished bool   `json:"is_published"`
}

type NewBlueprint struct {
	CommonBlueprintFields
	Provider           string `json:"provider"`
	IsTerraformManaged bool   `json:"is_terraform_managed"`
	IsPublished        bool   `json:"is_published"`
}

type UpdatedBlueprint struct {
	SeriesId string `json:"-"`
	CommonBlueprintFields
}

type PatchedBlueprint struct {
	SeriesId    string `json:"-"`
	IsPublished bool   `json:"is_published"`
}

type CommonBlueprintFields struct {
	Name                          string   `json:"name"`
	Description                   string   `json:"description,omitempty"`
	Content                       string   `json:"content"`
	Categories                    []string `json:"categories"`
	Guidance                      string   `json:"guidance"`
	Labels                        []Label  `json:"labels"`
	ExcludedContextQuestionSeries []string `json:"excluded_context_question_series"`
}

func (s *BlueprintsService) GetBlueprintBySeriesId(ctx context.Context, seriesId string) (*Blueprint, *http.Response, error) {
	path := fmt.Sprintf("%s/blueprints/series/%s", s.Client.BasePath, seriesId)
	body, resp, err := s.Client.Get(ctx, path, url.Values{}, new(Blueprint))
	if err != nil {
		return nil, resp, err
	}
	return body.(*Blueprint), resp, nil
}

func (s *BlueprintsService) CreateBlueprint(ctx context.Context, newBlueprint *NewBlueprint) (*Blueprint, *http.Response, error) {
	path := fmt.Sprintf("%s/blueprints", s.Client.BasePath)
	body, resp, err := s.Client.Post(ctx, path, newBlueprint, new(Blueprint))
	if err != nil {
		return nil, resp, err
	}
	return body.(*Blueprint), resp, nil
}

func (s *BlueprintsService) UpdateBlueprint(ctx context.Context, updatedBlueprint *UpdatedBlueprint) (*Blueprint, *http.Response, error) {
	path := fmt.Sprintf("%s/blueprints/series/%s", s.Client.BasePath, updatedBlueprint.SeriesId)
	body, resp, err := s.Client.Put(ctx, path, updatedBlueprint, new(Blueprint))
	if err != nil {
		return nil, resp, err
	}
	return body.(*Blueprint), resp, nil
}

func (s *BlueprintsService) PatchBlueprint(ctx context.Context, patchedBlueprint *PatchedBlueprint) (*Blueprint, *http.Response, error) {
	path := fmt.Sprintf("%s/blueprints/series/%s", s.Client.BasePath, patchedBlueprint.SeriesId)
	body, resp, err := s.Client.Patch(ctx, path, patchedBlueprint, new(Blueprint))
	if err != nil {
		return nil, resp, err
	}
	return body.(*Blueprint), resp, nil
}

func (s *BlueprintsService) DeleteBlueprint(ctx context.Context, blueprintSeriesId string) (*http.Response, error) {
	path := fmt.Sprintf("%s/blueprints/series/%s", s.Client.BasePath, blueprintSeriesId)
	return s.Client.Delete(ctx, path)
}
