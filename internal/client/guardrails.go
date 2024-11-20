package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type GuardrailsService service

type Guardrail struct {
	Id       string `json:"id"`
	SeriesId string `json:"series_id"`
	Version  int64  `json:"version"`
	Scope    string `json:"scope"`

	CommonGuardrailFields

	GuardrailTemplate struct {
		SeriesId string `json:"series_id"`
	} `json:"guardrail_template"`
	GuardrailTemplateInputs interface{} `json:"guardrail_template_inputs"`
}

type NewGuardrail struct {
	CommonGuardrailFields

	GuardrailTemplateSeriesId string      `json:"guardrail_template_series_id"`
	GuardrailTemplateInputs   interface{} `json:"guardrail_template_inputs"`

	IsTerraformManaged bool `json:"is_terraform_managed"`
}

type UpdatedGuardrail struct {
	SeriesId string `json:"-"`

	CommonGuardrailFields

	GuardrailTemplateSeriesId string      `json:"guardrail_template_series_id"`
	GuardrailTemplateInputs   interface{} `json:"guardrail_template_inputs"`
}

type CommonGuardrailFields struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Provider    string `json:"provider"`
	Category    string `json:"category"`
	State       string `json:"state"`
	Content     string `json:"content"`
}

func (s *GuardrailsService) GetGuardrailBySeriesId(ctx context.Context, seriesId string) (*Guardrail, *http.Response, error) {
	path := fmt.Sprintf("%s/guardrails/series/%s", s.Client.BasePath, seriesId)
	body, resp, err := s.Client.Get(ctx, path, url.Values{}, new(Guardrail))
	if err != nil {
		return nil, resp, err
	}
	return body.(*Guardrail), resp, nil
}

func (s *GuardrailsService) CreateGuardrail(ctx context.Context, newGuardrail *NewGuardrail) (*Guardrail, *http.Response, error) {
	path := fmt.Sprintf("%s/guardrails", s.Client.BasePath)
	body, resp, err := s.Client.Post(ctx, path, newGuardrail, new(Guardrail))
	if err != nil {
		return nil, resp, err
	}
	return body.(*Guardrail), resp, nil
}

func (s *GuardrailsService) UpdateGuardrail(ctx context.Context, updatedGuardrail *UpdatedGuardrail) (*Guardrail, *http.Response, error) {
	path := fmt.Sprintf("%s/guardrails/series/%s", s.Client.BasePath, updatedGuardrail.SeriesId)
	body, resp, err := s.Client.Put(ctx, path, updatedGuardrail, new(Guardrail))
	if err != nil {
		return nil, resp, err
	}
	return body.(*Guardrail), resp, nil
}

func (s *GuardrailsService) DeleteGuardrail(ctx context.Context, guardrailSeriesId string) (*http.Response, error) {
	path := fmt.Sprintf("%s/guardrails/series/%s", s.Client.BasePath, guardrailSeriesId)
	return s.Client.Delete(ctx, path)
}
