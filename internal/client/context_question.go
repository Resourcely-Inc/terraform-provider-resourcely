package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type ContextQuestionsService service

type NewContextQuestion struct {
	CommonContextQuestionFields
}

type UpdatedContextQuestion struct {
	SeriesId string `json:"-"`
	CommonContextQuestionFields
}

type CommonContextQuestionFields struct {
	Label                   string                  `json:"label"`
	Prompt                  string                  `json:"prompt"`
	Qtype                   string                  `json:"qtype"`
	AnswerFormat            string                  `json:"answer_format,omit_empty"`
	Scope                   string                  `json:"scope"`
	ContextQuestionOptions  []ContextQuestionOption `json:"context_question_options"`
	BlueprintCategories     []string                `json:"blueprint_categories"`
	RegexPattern            string                  `json:"regex_pattern"`
	ExcludedBlueprintSeries []string                `json:"excluded_blueprint_series"`
}

type ContextQuestion struct {
	Id       string `json:"id"`
	SeriesId string `json:"series_id"`
	Version  int64  `json:"version"`

	CommonContextQuestionFields
}

func (s *ContextQuestionsService) GetContextQuestionBySeriesId(ctx context.Context, seriesId string) (*ContextQuestion, *http.Response, error) {
	path := fmt.Sprintf("%s/context-questions/series/%s", s.Client.BasePath, seriesId)
	body, resp, err := s.Client.Get(ctx, path, url.Values{}, new(ContextQuestion))
	if err != nil {
		return nil, resp, err
	}
	return body.(*ContextQuestion), resp, nil
}

func (s *ContextQuestionsService) CreateContextQuestion(ctx context.Context, newContextQuestion *NewContextQuestion) (*ContextQuestion, *http.Response, error) {
	path := fmt.Sprintf("%s/context-questions", s.Client.BasePath)
	body, resp, err := s.Client.Post(ctx, path, newContextQuestion, new(ContextQuestion))
	if err != nil {
		return nil, resp, err
	}
	return body.(*ContextQuestion), resp, nil
}

func (s *ContextQuestionsService) UpdateContextQuestion(ctx context.Context, updatedContextQuestion *UpdatedContextQuestion) (*ContextQuestion, *http.Response, error) {
	path := fmt.Sprintf("%s/context-questions/series/%s", s.Client.BasePath, updatedContextQuestion.SeriesId)
	body, resp, err := s.Client.Put(ctx, path, updatedContextQuestion, new(ContextQuestion))
	if err != nil {
		return nil, resp, err
	}
	return body.(*ContextQuestion), resp, nil
}

func (s *ContextQuestionsService) DeleteContextQuestion(ctx context.Context, ContextQuestionSeriesId string) (*http.Response, error) {
	path := fmt.Sprintf("%s/context-questions/series/%s", s.Client.BasePath, ContextQuestionSeriesId)
	return s.Client.Delete(ctx, path)
}
