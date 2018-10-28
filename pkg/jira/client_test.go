package jira

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

const (
	testHostUrl  = "https://test.test"
	testUsername = "test"
	testPassword = "123456"
	testCookie   = "random value"

	testNotificationsCount = 3
)

var testClientNotifications = []Notification{
	{
		Title: "Testovací Notifikace",
	},
}

type httpClientMock struct{}

func (c httpClientMock) Do(r *http.Request) (*http.Response, error) {
	switch r.URL.String() {
	case restAuthUrl:
		return c.handleAuthRequest(r)
	case testHostUrl:
		return c.handleHostRequest(r)
	case testHostUrl + restNotifications:
		return c.handleNotificationsRequest(r)
	case testHostUrl + restNotificationCount:
		return c.handleNotificationsCountRequest(r)
	}

	res := new(http.Response)
	res.StatusCode = http.StatusNotFound
	return res, nil
}

func (c httpClientMock) handleAuthRequest(r *http.Request) (*http.Response, error) {
	decoder := json.NewDecoder(r.Body)
	var reqData map[string]string

	err := decoder.Decode(&reqData)
	if err != nil {
		return nil, err
	}

	resp := new(http.Response)

	if r.Method != http.MethodPost {
		resp.StatusCode = http.StatusNotFound
		return resp, nil
	}

	if reqData["username"] == testUsername && reqData["password"] == testPassword {
		resp.StatusCode = http.StatusOK
		resp.Header = make(http.Header)
		resp.Header.Set("Set-Cookie", testCookie)
	} else {
		resp.StatusCode = http.StatusForbidden
	}

	return resp, nil
}

func (c httpClientMock) handleHostRequest(r *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.StatusCode = http.StatusOK
	return resp, nil
}

func (c httpClientMock) handleNotificationsRequest(r *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	if r.Header.Get("Cookie") == testCookie {
		data := make(map[string]interface{})
		data["data"] = testClientNotifications
		notifications, _ := json.Marshal(data)
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBufferString(string(notifications)))
	} else {
		resp.StatusCode = http.StatusUnauthorized
	}
	return resp, nil
}

func (c httpClientMock) handleNotificationsCountRequest(r *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	if r.Header.Get("Cookie") == testCookie {
		count := Count{
			Count: testNotificationsCount,
		}
		data, _ := json.Marshal(count)
		resp.StatusCode = http.StatusOK
		resp.Body = ioutil.NopCloser(bytes.NewBufferString(string(data)))
	} else {
		resp.StatusCode = http.StatusUnauthorized
	}
	return resp, nil
}

func TestNewClient(t *testing.T) {
	expected := &client{
		testHostUrl,
		false,
		&http.Client{
			Timeout: time.Second * 5,
		},
		"",
	}
	actual := NewClient(testHostUrl)

	assert.Equal(t, expected, actual)
}

func TestClient_Login(t *testing.T) {
	client := &client{
		testHostUrl,
		false,
		&httpClientMock{},
		"",
	}

	err := client.Login(testUsername, testPassword)
	if assert.NoError(t, err) {
		assert.Equal(t, true, client.isLoggedIn)
	}
}

func TestClient_FetchNotificationCount(t *testing.T) {
	client := &client{
		testHostUrl,
		true,
		&httpClientMock{},
		testCookie,
	}

	count, err := client.FetchNotificationCount()
	if assert.NoError(t, err) {
		assert.Equal(t, testNotificationsCount, count)
	}
}

func TestClient_FetchNotifications(t *testing.T) {
	client := &client{
		testHostUrl,
		true,
		&httpClientMock{},
		testCookie,
	}

	notifications, err := client.FetchNotifications()
	if assert.NoError(t, err) {
		assert.Equal(t, testClientNotifications, notifications)
	}
}
