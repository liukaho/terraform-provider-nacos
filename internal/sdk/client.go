package sdk

import (
	"context"
	"github.com/liukaho/terraform-provider-nacos/internal/sdk/auth"
	"github.com/liukaho/terraform-provider-nacos/internal/sdk/namespace"
)

const ACCESS_TOKEN string = "accessToken"

type NacosClient struct {
	auth.AuthClient
	namespace.NamespaceClient
	username string
	password string
	host     string
	ctx      context.Context
}

func NewNacosClient(username, password, host string) (*NacosClient, error) {
	nacosClient := &NacosClient{
		username: username,
		password: password,
		host:     host,
	}

	nacosClient.initClient()
	if err := nacosClient.initAccessToken(); err != nil {
		return nil, err
	}
	return nacosClient, nil
}

func (nacosClient *NacosClient) initClient() {
	nacosClient.AuthClient = auth.NewAuthClient(nacosClient.host, nacosClient.username, nacosClient.password)
	nacosClient.NamespaceClient = namespace.NewNamespaceClient(nacosClient.host)
}

func (nacosClient *NacosClient) initAccessToken() error {
	if len(nacosClient.username) == 0 && len(nacosClient.password) == 0 {
		return nil
	}

	accessToken, err := nacosClient.Login()
	if err != nil {
		return err
	}
	nacosClient.ctx = context.WithValue(context.Background(), ACCESS_TOKEN, accessToken)

	return nil
}
func GetAccessToken(ctx context.Context) string {
	return ctx.Value(ACCESS_TOKEN).(string)
}
