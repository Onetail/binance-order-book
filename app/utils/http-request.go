package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
)

func closeReponseBody(body io.ReadCloser) {
	err := body.Close()
	if err != nil {
		logrus.Error(err.Error())
	}
}

func decodeResponseBody[T any](resp *http.Response) (*T, error) {
	var apiResponse map[string]interface{}
	defer closeReponseBody(resp.Body)
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		digest, _ := io.ReadAll(resp.Body)
		logrus.WithField("status code", resp.StatusCode).Error("unexpected response: ", string(digest))
		return nil, errors.New(string(digest))
	}

	err := json.NewDecoder(resp.Body).Decode(&apiResponse)
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}

	fmt.Println("\n\033[32m--- Debug ----")
	fmt.Printf("\033[36mapiResponse = %+v\n", apiResponse)
	fmt.Println("\033[32m\n---------------\033[0m")
	jsonString, _ := json.Marshal(apiResponse)
	var st T
	json.Unmarshal(jsonString, &st)

	return &st, nil
}

func CallAPI[T any](url string, method string, data any) (*T, error) {

	var bodyData bytes.Reader
	if data != nil {
		jsonValue, err := json.Marshal(data)
		if err != nil {
			logrus.Error(err.Error())
			return nil, err
		}
		bodyData = *bytes.NewReader([]byte(jsonValue))
	}

	req, httpErr := http.NewRequest(method, url, &bodyData)
	if httpErr != nil {
		logrus.Error(httpErr.Error())
		return nil, httpErr
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("accept", "*/*")
	req.Header.Set("Access-Control-Allow-Origin", "*")
	req.Header.Set("Access-Control-Allow-Credentials", "true")
	req.Header.Set("Access-Control-Max-Age", "86400")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}

	return decodeResponseBody[T](resp)
}
