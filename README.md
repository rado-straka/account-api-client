## Account API Client
Solution to the [interview-accountapi](https://github.com/form3tech-oss/interview-accountapi/tree/HEAD@{2021-02-25}) exercise.

Run fake API and see client's tests running just by `docker-compose up`.

#### Usage of the client

```go
client := Client{URL: "<URL>"} // by default uses http.Client{}, override with Client: <HTTP_CLIENT>
account := Account{
	Data: Data{
		Type:           "accounts",
		ID:             "<UUID>",
		OrganisationID: "<UUID>",
		Attributes: Attributes{
			Country: "GB",
		},
	},
}
responseAccount, err := client.Create(account)
responseAccount, err := client.Fetch(account.Data.ID)
err := client.Delete(account.Data.ID, account.Data.Version)
```

#### Testing
- client_test.go: integration tests extensively covers functionality of not exported functions in client.go, therefore these functions does not have their own unit tests.
- Used t.Cleanup using client's own functionality to do clean up what is not the best do to. After some research of options to avoid that I would consider [DbCleaner](https://pkg.go.dev/gopkg.in/khaiql/dbcleaner.v2).
- Tested against required account fields only. To test against all fields it would be necessary to blank dynamic fields in response json to avoid unnecessary long(and errorprone) tests, what I considered above the scope of this exercise.

**Note**: attributes data.attributes.alternative_names, data.attributes.name (and some other I do not remember) are not returned by the fake API.
