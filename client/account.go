package client

import "encoding/json"

type Account struct {
	Data  Data  `json:"data,omitempty"`
	Links Links `json:"links,omitempty`
}

type Data struct {
	ID             string     `json:"id,omitempty"`
	OrganisationID string     `json:"organisation_id,omitempty"`
	Type           string     `json:"type,omitempty"`
	Version        int        `json:"version,omitempty"`
	CreatedOn      string     `json:"created_on,omitempty"`
	ModifiedOn     string     `json:"modified_on,omitempty"`
	Attributes     Attributes `json:"attributes,omitempty"`
}

type Attributes struct {
	Country                 string    `json:"country,omitempty"`
	BaseCurrency            string    `json:"base_currency,omitempty"`
	BankID                  string    `json:"bank_id,omitempty"`
	BankIDCode              string    `json:"bank_id_code,omitempty"`
	AccountNumber           string    `json:"account_number,omitempty"`
	BIC                     string    `json:"bic,omitempty"`
	IBAN                    string    `json:"iban,omitempty"`
	CustomerID              string    `json:"customer_id,omitempty"`
	Name                    [4]string `json:"name,omitempty"`
	AlternativeNames        [3]string `json:"alternative_names,omitempty"`
	AccountClassification   string    `json:"account_classification,omitempty"`
	JointAccount            bool      `json:"joint_account,omitempty"`
	AccountMatchingOptOut   bool      `json:"account_matching_opt_out,omitempty"`
	SecondaryIdentification string    `json:"secondary_identification,omitempty"`
	Switched                bool      `json:"switched,omitempty"`
	Status                  string    `json:"status,omitempty"`
}

type Links struct {
	Self  string `json:"self,omitempty"`
	First string `json:"first,omitempty"`
	Last  string `json:"last,omitempty"`
	Next  string `json:"next,omitempty"`
	Prev  string `json:"prev,omitempty"`
}

// Generates JSON string representation of an Account
func (a *Account) json() (string, error) {
	bytes, err := json.Marshal(a)
	if err != nil {
		return "", err
	}
	return string(bytes[:]), nil
}

// Builds Account from json data
func (a *Account) build(jsonData *[]byte) error {
	return json.Unmarshal(*jsonData, a)
}
