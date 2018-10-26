package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"
)

const (
	restNotifications = "/gateway/api/notification-log/api/2/notifications"
	restAuthUrl       = "https://id.atlassian.com/id/rest/login"
)

// Data
type (
	// notifications is the multiple of notification data object
	Notifications struct {
		Notifications []notification `json:"data"`
	}

	// notification is a single data object received by API endpoint
	notification struct {
		Title     string            `json:"title"`
		Users     map[string]string `json:"users"`
		Template  string            `json:"template"`
		Timestamp string            `json:"timestamp"`
		Metadata  metadata          `json:"metadata"`
	}

	metadata struct {
		User user `json:"user"`
	}

	user struct {
		AtlassianId string `json:"atlassianId"`
		Name        string `json:"name"`
	}

	// client is service for our data fetching
	client struct {
		host       string
		isLoggedIn bool
		client     httpClient
		cookie     string
	}
)

// public types
type (
	Client interface {
		FetchNewNotifications() *Notifications
		Login(username, password string) error
	}

	httpClient interface {
		Do(*http.Request) (*http.Response, error)
	}
)

// Factory to our fetcher
func New(host string) Client {
	if strings.HasSuffix(host, "/") {
		strings.TrimSuffix(host, "/")
	}

	jar, _ := cookiejar.New(nil)

	return &client{
		host,
		false,
		&http.Client{
			Timeout: time.Second * 2,
			Jar:     jar,
		},
		"",
	}
}

func (c *client) Login(username, password string) error {
	credentials := make(map[string]string)
	credentials["username"] = username
	credentials["password"] = password

	j, err := json.Marshal(credentials)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", restAuthUrl, bytes.NewBuffer(j))
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		return err
	}

	res, err := c.client.Do(req)

	c.cookie = res.Header.Get("Set-Cookie")
	switch res.StatusCode {
	case http.StatusNotFound:
		return errors.New("Nepodařilo se navázat spojení se serverem. Zkontrolujte zadanou adresu jiry")
	case http.StatusForbidden:
		return errors.New("Přihlašovací údaje byly zadány chybně")
	}

	c.isLoggedIn = true

	return nil
}

//
func (c client) FetchNewNotifications() *Notifications {
	if !c.isLoggedIn {
		log.Fatal(errors.New("Uživatel není přihlášen!"))
	}

	req, err := http.NewRequest(http.MethodGet, c.host+restNotifications, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Cookie", c.cookie)
	res, err := c.client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var data *Notifications
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatal(err)
	}

	return data
}
