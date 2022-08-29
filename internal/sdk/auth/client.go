package auth

import (
	"encoding/json"
	"github.com/liukaho/terraform-provider-nacos/internal/sdk"
	"io/ioutil"
	"net/http"
	"net/url"
)

type loginResp struct {
	AccessToken string `json:"accessToken"`
	TokenTtl    int    `json:"tokenTtl"`
	GlobalAdmin bool   `json:"globalAdmin"`
	Username    string `json:"username"`
}

type AuthClient struct {
	host     string
	username string
	password string
}

func NewAuthClient(host, username, password string) AuthClient {
	var authClient AuthClient
	authClient.host = host
	authClient.username = username
	authClient.password = password

	return authClient
}

func (auth AuthClient) Login() (string, error) {

	resp, err := http.PostForm(auth.host+sdk.LOGIN_PATH, auth.loginFormData())
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var loginResp loginResp
	if err = json.Unmarshal(respBody, &loginResp); err != nil {
		return "", err
	}

	return loginResp.AccessToken, nil
}

func (auth AuthClient) loginFormData() url.Values {
	formdata := make(url.Values)
	formdata.Add("username", auth.username)
	formdata.Add("password", auth.password)

	return formdata
}
