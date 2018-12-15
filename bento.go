/*
Package bento provides an interface to the bentoforbusiness.com API.
Please see https://apidocs.bentoforbusiness.com/ for detailed info on
how the API works.

Basics

In order to begin interacting with Bento, you need to obtain a *Session
This can be done via the following functions:
For production:
	session, err := bento.GetProductionSession("myAccessKey", "mySecretKey")

For the sandbox:
	session, err := bento.GetTestSession("myTestAccessKey", "myTestSecretKey")

Once you have a session, you can begin doing things like getting and updating
cards and other resources associated with your bento account.

See the Session object's methods to see the types of objects you can interact
with. This is a good starting point from which you can begin to understand the
other types provided in this package.
*/
package bento
import (
	"fmt"
	"net/http"
	"bytes"
	"errors"
	"encoding/json"
	"io/ioutil"
	"log"
)

// Session provides the entry point to interact with the API.
// Created with GetProductionSession and GetTestSession.
type Session struct {
	apiUri string
	authorization string
	requester func(*Session, string, string, interface{}) ([]byte, error)
	logger *log.Logger
}

// AddressType can be "BUSINESS_ADDRESS" or "USER_ADDRESS"
type AddressType string

// Valid values for AddressType
const (
	BUSINESS_ADDRESS AddressType = "BUSINESS_ADDRESS"
	USER_ADDRESS     AddressType = "USER_ADDRESS"
)

type Address struct {
	Active bool             `json:"active"`
	AddressType AddressType `json:"addressType,omitempty"`
	City string             `json:"city,omitempty"`
	Id int64                `json:"id,omitempty"`
	State string            `json:"state,omitempty"`
	Street string           `json:"street,omitempty"`
	ZipCode string          `json:"zipCode,omitempty"`
}

type Business struct {
	BusinessId int64         `json:"businessId,omitempty"`
	CompanyName string       `json:"companyName,omitempty"`
	NameOnCard string        `json:"nameOnCard,omitempty"`
	Phone string             `json:"phone,omitempty"`
	AccountNumber string     `json:"accountNumber,omitempty"`
	BusinessStructure string `json:"businessStructure,omitempty"`
	Status string            `json:"status,omitempty"`
	CreatedDate int64        `json:"createdDate,omitempty"`
	ApprovalDate int64       `json:"approvalDate,omitempty"`
	ApprovalStatus string    `json:"approvalStatus,omitempty"`
	Balance float64          `json:"balance,omitempty"`
	TimeZone string          `json:"timeZone,omitempty"`
	Addresses []Address      `json:"addresses,omitempty"`
}

type ApiApplication struct {
	ApiApplicationId int64 `json:"apiApplicationId,omitempty"`
	Name string            `json:"name,omitempty"`
	AccessKey string       `json:"accessKey,omitempty"`
	Business Business      `json:"business,omitempty"`
}

// Period for a SpendingLimit
type Period string

// Valid periods for SpendingLimit
const (
	PERIOD_DAY    = "Day"
	PERIOD_WEEK   = "Week"
	PERIOD_MONTH  = "Month"
	PERIOD_CUSTOM = "Custom"
)

type SpendingLimit struct {
	Active bool           `json:"active"`
	Amount float64        `json:"amount,omitempty"`
	Period Period         `json:"period,omitempty"`
	CustomStartDate int64 `json:"customStartDate,omitempty"`
	CustomEndDate int64   `json:"customEndDate,omitempty"`
}

type User struct {
	FirstName string  `json:"firstName,omitempty"`
	LastName string   `json:"lastName,omitempty"`
	BirthDate int64   `json:"birthDate,omitempty"`
	Email string      `json:"email,omitempty"`
	Phone string      `json:"phone,omitempty"`
	UserId int64      `json:"userId,omitempty"`
	MobileAccess bool `json:"mobileAccess"`
	Deleted bool      `json:"deleted"`
	Created int64     `json:"created"`
	BentoType string  `json:"bentoType,omitempty"`
}

type Category struct {
	TransactionCategoryId int64 `json:"transactionCategoryId,omitempty"`
	Description string          `json:"description,omitempty"`
	Group string                `json:"group,omitempty"`
	Mccs []int64                `json:"mccs,omitempty"`
	Name string                 `json:"name,omitempty"`
	Type string                 `json:"type,omitempty"`
	BentoType string            `json:"bentoType,omitempty"`
}

type CardType string

// Valid Card Types
const (
	BUSINESS_OWNER_CARD CardType = "BusinessOwnerCard"
	EMPLOYEE_CARD CardType       = "EmployeeCard"
	CATEGORY_CARD CardType       = "CategoryCard"
)

// Valid Card Statuses
const (
	STATUS_CANCELED           = "CANCELED"
	STATUS_FRAUD_PREVENTION   = "FRAUD_PREVENTION"
	STATUS_TURNED_ON          = "TURNED_ON"
	STATUS_TURNED_OFF         = "TURNED_OFF"
	STATUS_WEEKLY_RESTRICTION = "WEEKLY_RESTRICTION"
)

type Card struct {
	CardId int64                 `json:"cardId,omitempty"`
	Type CardType                `json:"type,omitempty"`
	LifecycleStatus string       `json:"lifecycleStatus,omitempty"`
	Status string                `json:"status,omitempty"`
	Expiration string            `json:"expiration,omitempty"`
	LastFour string              `json:"lastFour,omitempty"`
	VirtualCard bool             `json:"virtualCard"`
	Alias string                 `json:"alias,omitempty"`
	AvailableAmount float64      `json:"availableAmount,omitempty"`
	AllowedDaysActive bool       `json:"allowedDaysActive"`
	AllowedDays []string         `json:"allowedDays,omitempty"`
	AllowedCategoriesActive bool `json:"allowedCategoriesActive"`
	AllowedCategories []Category `json:"allowedCategories,omitempty"`
	TransactionCategoryId int64  `json:"transactionCategoryId,omitempty"`
	CreatedOn int64              `json:"createdOn,omitempty"`
	UpdatedOn int64              `json:"updatedOn,omitempty"`
	SpendingLimit SpendingLimit  `json:"spendingLimit,omitempty"`
	User User                    `json:"user,omitempty"`
	Permissions map[string]bool  `json:"permissions,omitempty"`
	BentoType string             `json:"bentoType,omitempty"`
	session *Session             `json:"-"`
}

type PanAndCvv struct {
	Pan string `json:"pan,omitempty"`
	Cvv string `json:"cvv,omitempty"`
}

type BentoError struct {
	Message string
	BentoError string `json:"error"`
}

func (e BentoError) Error() string {
	return fmt.Sprintf("Bento Error: [%s], [%s]", e.Message, e.BentoError)
}

func checkError(bs []byte) error {
	var bentoErr BentoError
	err := json.Unmarshal(bs, &bentoErr)
	if err != nil {
		// Even if we can't unmarshal, it may be another kind
		// of valid object.
		return nil
	}
	if bentoErr.Message != "" ||
		bentoErr.BentoError != "" {
		return bentoErr
	}
	return nil
}

var sandboxUri string = "https://sandbox-api.bentoforbusiness.com/api"
var productionUri string = "https://api.bentoforbusiness.com"

func GetProductionSession(accessKey, secretKey string) (*Session, error) {
	return getSession(productionUri, accessKey, secretKey)
}

func GetTestSession(accessKey, secretKey string) (*Session, error) {
	return getSession(sandboxUri, accessKey, secretKey)
}

func getSession(apiUri, accessKey, secretKey string) (*Session, error) {

	client := &http.Client{}

	bs, err := json.Marshal(
		map[string]string{
			"accessKey": accessKey,
			"secretKey": secretKey})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST",
		fmt.Sprintf("%s/sessions", apiUri),
		bytes.NewReader(bs))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application\\json")
	req.Header.Add("Accept", "*/*")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if !json.Valid(body) {
		return nil, errors.New(
			fmt.Sprintf("Invalid json response: [%s]", string(body)))
	}

	auth, ok := resp.Header["Authorization"]
	if !ok {
		return nil, errors.New("Server did not return an authorization token.")
	}

	app := &ApiApplication{}
	err = json.Unmarshal(body, app)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error unmarshalling: %s", err))
	}

	session := &Session{
		apiUri: apiUri,
		authorization: auth[0],
		requester: doRequest,
		logger: log.New(ioutil.Discard, "", 0),
	}
	return session, nil
}

func (session *Session) SetLogger(logger *log.Logger) {
	session.logger = logger
}

func doRequest(session *Session, method, endpoint string, args interface{}) ([]byte, error) {
	client := &http.Client{}

	var err error
	var req *http.Request
	if args != nil {
		bs, err := json.Marshal(args)
		if err != nil {
			return nil, err
		}

		req, err = http.NewRequest(method,
			fmt.Sprintf("%s%s", session.apiUri, endpoint),
			bytes.NewReader(bs))
		if err != nil {
			return nil, err
		}

		req.Header.Add("Content-Type", "application\\json")
		req.Header.Add("Accept", "*/*")
		req.Header.Add("Authorization", session.authorization)
		session.logger.Printf("Sending request: [method: %s] [uri: %s] body: %s",
			method, fmt.Sprintf("%s%s", session.apiUri, endpoint), string(bs))
	} else {
		req, err = http.NewRequest(method,
			fmt.Sprintf("%s%s", session.apiUri, endpoint),
			nil)
		if err != nil {
			return nil, err
		}

		req.Header.Add("Accept", "*/*")
		req.Header.Add("Authorization", session.authorization)
		session.logger.Printf("Sending request: [method: %s] [uri: %s]",
			method, fmt.Sprintf("%s%s", session.apiUri, endpoint))
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	session.logger.Printf("Received response: %s", string(body))
	if err != nil {
		return nil, err
	}

	if !json.Valid(body) {
		return nil, errors.New(
			fmt.Sprintf("Server returned non-json value: [%s]", string(body)))
	}

	err = checkError(body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (session *Session) GetBusiness() (*Business, error) {
	bs, err := session.requester(session, "GET", "/businesses/me", nil)
	if err != nil {
		return nil, err
	}

	business := &Business{}
	err = json.Unmarshal(bs, business)
	if err != nil {
		return nil, err
	}

	return business, nil
}

func (session *Session) GetCards() ([]Card, error) {
	bs, err := session.requester(session, "GET", "/cards", nil)
	if err != nil {
		return nil, err
	}

	var cards []Card
	err = json.Unmarshal(bs, &cards)
	if err != nil {
		return nil, err
	}

	for i := range cards {
		cards[i].session = session
	}

	return cards, nil
}

func (session *Session) GetCard(cardId int64) (*Card, error) {
	bs, err := session.requester(session, "GET", fmt.Sprintf("/cards/%d", cardId), nil)
	if err != nil {
		return nil, err
	}

	var card Card
	err = json.Unmarshal(bs, &card)
	if err != nil {
		return nil, err
	}

	card.session = session
	return &card, nil
}

func (session *Session) NewCard(cardType CardType, alias string) (*Card, error) {
	bs, err := session.requester(session, "POST", "/cards",
		map[string]interface{}{
			"type": cardType,
			"alias": alias,
			"virtualCard": false,
			"lastFour": 3215,
		})
	if err != nil {
		return nil, err
	}

	var cardResp Card
	err = json.Unmarshal(bs, &cardResp)
	if err != nil {
		return nil, err
	}

	cardResp.session = session
	return &cardResp, nil
}

func (card *Card) Put() (*Card, error) {
	bs, err := card.session.requester(card.session, "PUT", fmt.Sprintf("/cards/%d", card.CardId), card)
	if err != nil {
		return nil, err
	}

	var cardResp Card
	err = json.Unmarshal(bs, &cardResp)
	if err != nil {
		return nil, err
	}

	cardResp.session = card.session
	return &cardResp, nil
}

func (card *Card) Delete() (*Card, error) {
	bs, err := card.session.requester(card.session, "DELETE", fmt.Sprintf("/cards/%d", card.CardId), nil)
	if err != nil {
		return nil, err
	}

	var cardResp Card
	err = json.Unmarshal(bs, &cardResp)
	if err != nil {
		return nil, err
	}

	cardResp.session = card.session
	return &cardResp, nil
}

func (card *Card) Activate(lastFour string) (*Card, error) {
	card.LastFour = lastFour
	bs, err := card.session.requester(card.session,
		"POST",
		fmt.Sprintf("/cards/%d/activation", card.CardId),
		card)
	if err != nil {
		return nil, err
	}

	var cardResp Card
	err = json.Unmarshal(bs, &cardResp)
	if err != nil {
		return nil, err
	}

	cardResp.session = card.session
	return &cardResp, nil
}

func (card *Card) TurnOn() (*Card, error) {
	card.Status = STATUS_TURNED_ON
	card, err := card.Put()
	if err != nil {
		return nil, err
	}
	if card.Status != STATUS_TURNED_ON {
		return nil, errors.New(fmt.Sprintf("Bento returned success for turn on, but card's status is: %s", card.Status))
	}
	return card, nil
}

func (card *Card) TurnOff() (*Card, error) {
	card.Status = STATUS_TURNED_OFF
	card, err := card.Put()
	if err != nil {
		return nil, err
	}
	if card.Status != STATUS_TURNED_OFF {
		return nil, errors.New(fmt.Sprintf("Bento returned success for turn off, but card's status is: %s", card.Status))
	}
	return card, nil
}

func (card *Card) Reissue() (*Card, error) {
	bs, err := card.session.requester(card.session,
		"POST",
		fmt.Sprintf("/cards/%d/reissue", card.CardId),
		nil)
	if err != nil {
		return nil, err
	}

	var cardResp Card
	err = json.Unmarshal(bs, &cardResp)
	if err != nil {
		return nil, err
	}

	cardResp.session = card.session
	return &cardResp, nil
}

func (card *Card) GetPanAndCvv() (*PanAndCvv, error) {
	bs, err := card.session.requester(card.session,
		"GET",
		fmt.Sprintf("/cards/%d/pan", card.CardId),
		nil)
	if err != nil {
		return nil, err
	}

	var panAndCvv PanAndCvv
	err = json.Unmarshal(bs, &panAndCvv)
	if err != nil {
		return nil, err
	}

	return &panAndCvv, nil
}

func (card *Card) GetBillingAddress() (*Address, error) {
	bs, err := card.session.requester(card.session,
		"GET",
		fmt.Sprintf("/cards/%d/billingAddress", card.CardId),
		nil)
	if err != nil {
		return nil, err
	}

	var address Address
	err = json.Unmarshal(bs, &address)
	if err != nil {
		return nil, err
	}

	return &address, nil
}

func (card *Card) SetBillingAddress(newAddress *Address) (*Address, error) {
	bs, err := card.session.requester(card.session,
		"POST",
		fmt.Sprintf("/cards/%d/billingAddress", card.CardId),
		newAddress)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return nil, err
	}

	var address Address
	err = json.Unmarshal(bs, &address)
	if err != nil {
		return nil, err
	}

	return &address, nil
}

func (card *Card) UpdateBillingAddress(newAddress *Address) (*Address, error) {
	bs, err := card.session.requester(card.session,
		"PUT",
		fmt.Sprintf("/cards/%d/billingAddress", card.CardId),
		newAddress)
	if err != nil {
		return nil, err
	}

	var address Address
	err = json.Unmarshal(bs, &address)
	if err != nil {
		return nil, err
	}

	return &address, nil
}

// Transactions
type Transactions struct {
	Amount float64                 `json:"amount,omitempty"`
	Size int                       `json:"size",omitempty"`
	CardTransactions []Transaction `json:"cardTransactions"`
}

type Payee struct {
	Name string    `json:"name,omitempty"`
	City string    `json:"city,omitempty"`
	State string   `json:"state,omitempty"`
	Country string `json:"country,omitempty"`
	Zip string     `json:"zip,omitempty"`
}

type Transaction struct {
	CardTransactionId int64  `json:"cardTransactionId,omitempty"`
	Amount float64           `json:"amount,omitempty"`
	ApprovalCode string      `json:"approvalCode,omitempty"`
	AvailableBalance float64 `json:"availableBalance,omitempty"`
	Card *Card               `json:"card,omitempty"`
	Category *Category       `json:"category,omitempty"`
	Currency string          `json:"currency,omitempty"`
	Deleted bool             `json:"deleted,omitempty"`
	Fees float64             `json:"fees,omitempty"`
	LedgerBalance float64    `json:"ledgerBalance,omitempty"`
	Note string              `json:"node,omitempty"`
	SettlementDate int64     `json:"settlementDate,omitempty"`
	Status string            `json:"status,omitempty"`
	Tags []string            `json:"tags,omitempty"`
	TransactionDate int64    `json:"transactionDate,omitempty"`
	Type string              `json:"type,omitempty"`
	Payee *Payee             `json:"payee,omitempty"`
}


func (session *Session) GetTransactions() (*Transactions, error) {
	bs, err := session.requester(session, "GET", "/transactions", nil)
	if err != nil {
		return nil, err
	}

	var transaction Transactions
	err = json.Unmarshal(bs, &transaction)
	if err != nil {
		return nil, err
	}

	return &transaction, nil
}
