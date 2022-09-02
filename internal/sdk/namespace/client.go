package namespace

import (
	"context"
	"github.com/google/uuid"
	"github.com/liukaho/terraform-provider-nacos/internal/sdk"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"net/url"
)

type NamespaceClient struct {
	host string
}

var (
	CREATE_NAMESPACE_ERROR error = errors.New("create namespace failed")
)

func NewNamespaceClient(host string) NamespaceClient {
	var namespaceClient NamespaceClient
	namespaceClient.host = host

	return namespaceClient
}

func (nc NamespaceClient) Create(ctx context.Context, namespaceName string) (string, error) {
	return nc.create(ctx, uuid.New().String(), namespaceName, "")
}

func (nc NamespaceClient) create(ctx context.Context, customNamespaceId, namespaceName, namespaceDesc string) (string, error) {
	createRequestForm := make(url.Values)
	createRequestForm.Add("customNamespaceId", customNamespaceId)
	createRequestForm.Add("namespaceName", namespaceName)
	createRequestForm.Add("namespaceDesc", namespaceDesc)
	url, err := url.Parse(nc.host + "/nacos/v1/console/namespaces")
	if err != nil {
		return "", errors.Wrap(err, "url parse failed at create namespace")
	}
	if accessToken := sdk.GetAccessToken(ctx); len(accessToken) > 0 {
		url.Query().Add(sdk.ACCESS_TOKEN, sdk.GetAccessToken(ctx))
	}

	resp, err := http.PostForm(url.String(), createRequestForm)
	if err != nil {
		return "", errors.Wrap(err, "create namespace failed")
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "read create namespace resp failed")
	}

	if !isSuccess(body) {
		return "", CREATE_NAMESPACE_ERROR
	}

	return customNamespaceId, nil
}

func isSuccess(body []byte) bool {
	return string(body) == "true"
}
