package bitop

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"io"
	"main/internal/app/config"
	"main/internal/app/model"
	"net/http"
	"net/url"
)

type Client struct {
	ctx    context.Context
	body   model.BitopBody
	client *http.Client
}

func New(ctx context.Context) *Client {
	cfg := config.FromContext(ctx)
	return &Client{
		ctx: ctx,
		body: model.BitopBody{
			Token: cfg.BITOPToken,
		},
		client: &http.Client{},
	}
}

func (c *Client) GetBranch(ctx context.Context, branch string) (*model.ResponseBody, error) {
	cfg := config.FromContext(ctx).BITOP

	url := url.URL{
		Scheme: cfg.Protocol,
		Host:   cfg.SiteAdress,
		Path:   cfg.Path,
	}

	log.Info("url created", url.String())

	reqBody, _ := json.Marshal(model.RequestBody{
		"",
		branch,
		"branch",
	})

	reqToApi, err := http.NewRequest("POST", url.String(), bytes.NewBuffer(reqBody))
	if err != nil {
		log.WithError(err).Error("cant create request")
		return nil, err
	}

	reqToApi.Header = http.Header{
		"x-bb-token": {c.body.Token},
	}

	rawResp, err := c.client.Do(reqToApi)
	if err != nil {
		log.WithError(err).Error("cant do request")
		return nil, err
	}

	body, err := io.ReadAll(rawResp.Body)
	if err != nil {
		log.WithError(err).Error("cant read response")
		return nil, err
	}

	var resp model.ResponseBody

	err = json.Unmarshal(body, &resp)
	if err != nil {
		log.WithError(err).Error("cant unmarshal response")
		return nil, err
	}

	return &resp, err
}

func (c *Client) GetFaculty(ctx context.Context, parentUUID uuid.UUID) (*model.ResponseBody, error) {
	cfg := config.FromContext(ctx).BITOP

	url := url.URL{
		Scheme: cfg.Protocol,
		Host:   cfg.SiteAdress,
		Path:   cfg.Path,
	}

	reqBody, _ := json.Marshal(model.RequestBody{
		parentUUID.String(),
		"",
		"faculty",
	})

	reqToApi, err := http.NewRequest("POST", url.String(), bytes.NewBuffer(reqBody))
	if err != nil {
		log.WithError(err).Error("cant create request")
	}

	reqToApi.Header = http.Header{
		"x-bb-token": {c.body.Token},
	}

	rawResp, err := c.client.Do(reqToApi)
	if err != nil {
		log.WithError(err).Error("cant do request")
		return nil, err
	}

	body, err := io.ReadAll(rawResp.Body)
	if err != nil {
		log.WithError(err).Error("cant read response")
		return nil, err
	}
	var resp model.ResponseBody

	err = json.Unmarshal(body, &resp)
	if err != nil {
		log.WithError(err).Error("cant unmarshal response")
		return nil, err
	}

	return &resp, err
}

func (c *Client) GetDepartment(ctx context.Context, parentUUID uuid.UUID) (*model.ResponseBody, error) {
	cfg := config.FromContext(ctx).BITOP

	url := url.URL{
		Scheme: cfg.Protocol,
		Host:   cfg.SiteAdress,
		Path:   cfg.Path,
	}

	reqBody, _ := json.Marshal(model.RequestBody{
		parentUUID.String(),
		"",
		"department",
	})

	reqToApi, err := http.NewRequest("POST", url.String(), bytes.NewBuffer(reqBody))
	if err != nil {
		log.WithError(err).Error("cant create request")
	}

	reqToApi.Header = http.Header{
		"x-bb-token": {c.body.Token},
	}

	rawResp, err := c.client.Do(reqToApi)
	if err != nil {
		log.WithError(err).Error("cant do request")
		return nil, err
	}

	body, err := io.ReadAll(rawResp.Body)
	if err != nil {
		log.WithError(err).Error("cant read response")
		return nil, err
	}
	var resp model.ResponseBody

	err = json.Unmarshal(body, &resp)
	if err != nil {
		log.WithError(err).Error("cant unmarshal response")
		return nil, err
	}

	return &resp, err
}

func (c *Client) GetGroup(ctx context.Context, parentUUID uuid.UUID) (*model.ResponseBody, error) {
	cfg := config.FromContext(ctx).BITOP

	url := url.URL{
		Scheme: cfg.Protocol,
		Host:   cfg.SiteAdress,
		Path:   cfg.Path,
	}

	reqBody, _ := json.Marshal(model.RequestBody{
		parentUUID.String(),
		"",
		"group",
	})

	reqToApi, err := http.NewRequest("POST", url.String(), bytes.NewBuffer(reqBody))
	if err != nil {
		log.WithError(err).Error("cant create request")
	}

	reqToApi.Header = http.Header{
		"x-bb-token": {c.body.Token},
	}

	rawResp, err := c.client.Do(reqToApi)
	if err != nil {
		log.WithError(err).Error("cant do request")
		return nil, err
	}

	body, err := io.ReadAll(rawResp.Body)
	if err != nil {
		log.WithError(err).Error("cant read response")
		return nil, err
	}
	var resp model.ResponseBody

	err = json.Unmarshal(body, &resp)
	if err != nil {
		log.WithError(err).Error("cant unmarshal response")
		return nil, err
	}

	return &resp, err
}
