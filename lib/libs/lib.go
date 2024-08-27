package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/imroc/req/v3"
)

func GetRequestToken() (token string, err error) {
	username := os.Getenv("SYSTEM_USER")
	password := os.Getenv("SYSTEM_USER_PASSWORD")
	url := os.Getenv("AUTH_TOKEN_URL")
	client := req.C().DevMode()
	postBody, _ := json.Marshal(map[string]string{
		"username": username,
		"password": password,
	})
	requestBody := bytes.NewBuffer(postBody)
	r, err := client.R().
		SetErrorResult(&err).
		SetHeaders(map[string]string{
			"accept":       "*/*",
			"Content-Type": "application/json",
			"X-CSRFTOKEN":  "BMwdKEdnRpZxlDW4GGO3DUGYSlMyfgAZUlA1RgTOcFPJytopuAeQp9LXDzciHV2S",
		}).
		SetBody(requestBody).
		Post(url)
	if err != nil {
		panic(err.Error())
	}
	if !r.IsSuccessState() {
		fmt.Println("bad response status:", r.Status)
		return "", err
	}

	buf, _ := io.ReadAll(r.Body)
	response := map[string]string{}
	json.Unmarshal([]byte(buf), &response)
	token = response["access"]
	return token, nil
}
