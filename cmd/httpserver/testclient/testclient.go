package testclient

import (
	"context"
	"github.com/google/uuid"
	"net"
	"net/http"
	"net/url"
	"testing"
)

//go:generate go tool oapi-codegen -config cfg.yaml ../../../api/api.yaml

type TClient struct {
	clientWithResponses *ClientWithResponses
}

func NewTestClient(t *testing.T, serverHost, serverPort string) *TClient {

	tmpUrl := url.URL{
		Scheme: "http",
		Host:   net.JoinHostPort(serverHost, serverPort),
	}

	c, err := NewClientWithResponses(tmpUrl.String())
	if err != nil {
		t.Fatal(err)
	}

	return &TClient{
		clientWithResponses: c,
	}
}

func (tc *TClient) WithBearer(JWT string) *TClient {

	server := tc.clientWithResponses.ClientInterface.(*Client).Server
	authString := "Bearer " + JWT

	ntc, err := NewClientWithResponses(server, WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
		req.Header.Add("Authorization", authString)
		return nil
	}))
	if err != nil {
		panic(err)
	}
	return &TClient{
		clientWithResponses: ntc,
	}
}

func (tc *TClient) DummyLogin(t *testing.T, body PostDummyLoginJSONBody) *PostDummyLoginResponse {
	t.Helper()
	res, err := tc.clientWithResponses.PostDummyLoginWithResponse(t.Context(), PostDummyLoginJSONRequestBody(body))
	if err != nil {
		t.Fatal(err)
	}
	return res
}

func (tc *TClient) CreatePVZ(t *testing.T, body PostPvzJSONRequestBody) *PostPvzResponse {
	t.Helper()
	res, err := tc.clientWithResponses.PostPvzWithResponse(t.Context(), body)
	if err != nil {
		t.Fatal(err)
	}
	return res
}

func (tc *TClient) CreateReception(t *testing.T, body PostReceptionsJSONBody) *PostReceptionsResponse {
	t.Helper()
	res, err := tc.clientWithResponses.PostReceptionsWithResponse(t.Context(), PostReceptionsJSONRequestBody(body))
	if err != nil {
		t.Fatal(t)
	}
	return res
}

func (tc *TClient) CreateProduct(t *testing.T, body PostProductsJSONBody) *PostProductsResponse {
	t.Helper()
	res, err := tc.clientWithResponses.PostProductsWithResponse(t.Context(), PostProductsJSONRequestBody(body))
	if err != nil {
		t.Fatal(t)
	}
	return res
}

func (tc *TClient) GetPVZ(t *testing.T, body GetPvzParams) *GetPvzResponse {
	t.Helper()
	res, err := tc.clientWithResponses.GetPvzWithResponse(t.Context(), &body)
	if err != nil {
		t.Fatal(err)
	}
	return res
}

func (tc *TClient) DeleteLastProduct(t *testing.T, pvzId uuid.UUID) *PostPvzPvzIdDeleteLastProductResponse {
	t.Helper()
	res, err := tc.clientWithResponses.PostPvzPvzIdDeleteLastProductWithResponse(t.Context(), pvzId)
	if err != nil {
		t.Fatal(err)
	}
	return res
}

func (tc *TClient) CloseLastReception(t *testing.T, pvzId uuid.UUID) *PostPvzPvzIdCloseLastReceptionResponse {
	t.Helper()
	res, err := tc.clientWithResponses.PostPvzPvzIdCloseLastReceptionWithResponse(t.Context(), pvzId)
	if err != nil {
		t.Fatal(err)
	}
	return res
}
