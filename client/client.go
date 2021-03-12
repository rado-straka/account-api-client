package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Client is an entrypoint for the use of its operations
type Client struct {
	// Represents base url for the endpoints
	URL    string
	Client http.Client
}

var accountsEndpoint = "/v1/organisation/accounts"

func (c *Client) Create(account Account) (*Account, error) {
	// build POST request from an account
	accountJSON, err := account.json()
	if err != nil {
		return nil, err
	}
	body := strings.NewReader(accountJSON)
	url := c.URL + accountsEndpoint
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/vnd.api+json")

	// execute request and read body
	res, resBody, err := c.execRequest(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// build and return error message if POST does not create an Account record
	err = checkResponseStatus(res, 201, resBody)
	if err != nil {
		return nil, err
	}

	// build and return created Account from response body
	resAccount := Account{}
	err = resAccount.build(resBody)
	if err != nil {
		return nil, err
	}
	return &resAccount, nil
}

// Executes request then reads response body
func (c *Client) execRequest(req *http.Request) (*http.Response, *[]byte, error) {
	res, err := c.Client.Do(req)
	if err != nil {
		return nil, nil, err
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return res, nil, err
	}
	return res, &resBody, err
}

func checkResponseStatus(res *http.Response, expectedStatusCode int,
	resBody *[]byte) error {
	if res.StatusCode != expectedStatusCode {
		errorMessage, err := errorMessage(res.Status, resBody)
		if err != nil {
			return err
		}
		return errors.New(errorMessage)
	}
	return nil
}

func errorMessage(resStatus string, resBody *[]byte) (string, error) {
	// unmarshal response body json into a map to access
	// top level string attributes
	var resJSON map[string]string
	err := json.Unmarshal(*resBody, &resJSON)
	if err != nil {
		return "", err
	}

	errorMessage, ok := resJSON["error_message"]
	if ok {
		return fmt.Sprintf("%s: %s", resStatus, errorMessage), nil
	}
	return resStatus, nil
}

func (c *Client) Fetch(accountID string) (*Account, error) {
	// check not empty id to not call List function of the API
	if accountID == "" {
		return nil, errors.New("Account.Data.ID is empty")
	}

	// build GET request with account id
	url := c.URL + accountsEndpoint + "/" + accountID
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// execute request and read body
	res, resBody, err := c.execRequest(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// build and return error message if GET does not return an Account record
	err = checkResponseStatus(res, 200, resBody)
	if err != nil {
		return nil, err
	}

	// build and return created Account from response body
	resAccount := Account{}
	err = resAccount.build(resBody)
	if err != nil {
		return nil, err
	}
	return &resAccount, nil
}

func (c *Client) Delete(accountID string, accountVersion int) error {
	// build DELETE request with an account id and an account version
	url := fmt.Sprintf("%s%s/%s?version=%d",
		c.URL, accountsEndpoint, accountID, accountVersion)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	// execute request and read body
	res, resBody, err := c.execRequest(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// build and return error message if DELETE fails
	err = checkResponseStatus(res, 204, resBody)
	if err != nil {
		return err
	}

	return nil
}
