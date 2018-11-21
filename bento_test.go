package bento

import (
	"fmt"
	"testing"
	"errors"
)

var SampleBusiness string = `{
  "businessId": 12345,
  "companyName": "My Company Inc",
  "nameOnCard": "My Company",
  "phone": "9998881234",
  "accountNumber": "820187766",
  "businessStructure": "LLC",
  "status": "APPROVED",
  "createdDate": 1495759408,
  "approvalDate": 1495759408,
  "approvalStatus": "Approved",
  "balance": 100.99,
  "timeZone": "America/Los_Angeles",
  "addresses": [
	{
	  "active": true,
	  "addressType": "BUSINESS_ADDRESS",
	  "city": "San Francisco",
	  "id": 12345,
	  "state": "CA",
	  "street": "123 Main Street",
	  "zipCode": "94123"
	}
  ]
}
`

var SampleCard string = `{
  "cardId": 12345,
  "type": "CategoryCard",
  "lifecycleStatus": "ACTIVATED",
  "status": "TURNED_ON",
  "expiration": "1221",
  "lastFour": "1234",
  "virtualCard": false,
  "alias": "My Card",
  "availableAmount": 123.45,
  "allowedDays": [
	"MONDAY"
  ],
  "allowedCategoriesActive": true,
  "allowedCategories": [
	{
	  "transactionCategoryId": 10
	}
  ],
  "createdOn": 1495759408,
  "updatedOn": 1495759408,
  "spendingLimit": {
	"active": true,
	"amount": 123.45,
	"period": "Day",
	"customStartDate": 1495759408,
	"customEndDate": 1495759408
  },
  "user": {
	"firstName": "John",
	"lastName": "Smith",
	"birthDate": 1495759408,
	"email": "me@myemail.com",
	"phone": "9998887654",
	"userId": 12345,
	"mobileAccess": true,
	"deleted": false,
	"created": 1495759408
  }
}`

type TestSession struct {
	Session
	method string
	endpoint string
	args interface{}
}

func testRequest(tbs *TestSession) func(session *Session, method, endpoint string, args interface{}) ([]byte, error) {
	return func(session *Session, method, endpoint string, args interface{}) ([]byte, error) {
		if tbs != nil {
			tbs.method = method
			tbs.endpoint = endpoint
			tbs.args = args
		}
		switch(endpoint) {
		case "/businesses/me":
			return []byte(SampleBusiness), nil
		case "/cards":
			if method == "GET" {
				return []byte(fmt.Sprintf("[%s,%s]", SampleCard, SampleCard)), nil
			} else if method == "POST" {
				return []byte(SampleCard), nil
			}
		case "/cards/12345":
			return []byte(SampleCard), nil
		}
		return nil, errors.New("No such testing endpoint.")
	}
}

func testRequestFailures(tbs *TestSession) func(session *Session, method, endpoint string, args interface{}) ([]byte, error) {
	return func(session *Session, method, endpoint string, args interface{}) ([]byte, error) {
		if tbs != nil {
			tbs.method = method
			tbs.endpoint = endpoint
			tbs.args = args
		}
		return nil, errors.New("<html>500 Error</html>")
	}
}

func TestGetBusiness(t *testing.T) {
	t.Log("TestGetBusiness")
	session := &Session{ requester: testRequest(nil) }

	business, err := session.GetBusiness()
	if err != nil {
		t.Errorf("Failed to unmarshal business object: %s", err)
		return
	}

	if business.BusinessId != 12345 {
		t.Error(`Expected BusinessId == 12345`)
	}
	if business.CompanyName != "My Company Inc" {
		t.Error(`Expected CompanyName == "My Company Inc"`)
	}
	if business.NameOnCard != "My Company" {
		t.Error(`Expected NameOnCard == "My Company"`)
	}
	if business.Phone != "9998881234" {
		t.Error(`Expected Phone == "9998881234"`)
	}
	if business.AccountNumber != "820187766" {
		t.Error(`Expected AccountNumber == "820187766"`)
	}
	if business.BusinessStructure != "LLC" {
		t.Error(`Expected BusinessStructure == "LLC"`)
	}
	if business.Status != "APPROVED" {
		t.Error(`Expected Status == "APPROVED"`)
	}
	if business.CreatedDate != 1495759408 {
		t.Error(`Expected CreatedDate == 1495759408`)
	}
	if business.ApprovalDate != 1495759408 {
		t.Error(`Expected ApprovalDate == 1495759408`)
	}
	if business.ApprovalStatus != "Approved" {
		t.Error(`Expected ApprovalStatus == "Approved"`)
	}
	if business.Balance != 100.99 {
		t.Error(`Expected Balance == 100.99`)
	}
	if business.TimeZone != "America/Los_Angeles" {
		t.Error(`Expected TimeZone == "America/Los_Angeles"`)
	}
	if !business.Addresses[0].Active {
		t.Error(`Expected address.Active == true`)
	}
	if business.Addresses[0].AddressType != "BUSINESS_ADDRESS" {
		t.Error(`Expected address.AddressType == "BUSINESS_ADDRESS"`)
	}
	if business.Addresses[0].City != "San Francisco" {
		t.Error(`Expected address.City == "San Francisco"`)
	}
	if business.Addresses[0].Id != 12345 {
		t.Error(`Expected address.Id == 12345`)
	}
	if business.Addresses[0].State != "CA" {
		t.Error(`Expected address.State == "CA"`)
	}
	if business.Addresses[0].Street != "123 Main Street" {
		t.Error(`Expected address.Street == "123 Main Street"`)
	}
	if business.Addresses[0].ZipCode != "94123" {
		t.Error(`Expected address.ZipCode == "94123"`)
	}
}

func TestGetBusinessFailure(t *testing.T) {
	t.Log("TestGetBusinessFailure")
	session := &Session{ requester: testRequestFailures(nil) }

	_, err := session.GetBusiness()
	if err == nil {
		t.Error("Expected failure when calling GetBusiness")
	}
}

func TestGetCards(t *testing.T) {
	t.Log("TestGetCards")
	session := &Session{ requester: testRequest(nil) }

	cards, err := session.GetCards()
	if err != nil {
		t.Errorf("Failed to unmarshal cards: %s", err)
		return
	}

	if len(cards) != 2 {
		t.Error("Expected len(cards) == 2.")
	}
}

func TestGetCardsFail(t *testing.T) {
	t.Log("TestGetCardsFail")
	session := &Session{ requester: testRequestFailures(nil) }

	_, err := session.GetCards()
	if err == nil {
		t.Error("Expected failure when calling GetCards.")
	}
}

func TestGetCard(t *testing.T) {
	t.Log("TestGetCard")

	session := &Session{ requester: testRequest(nil) }
	card, err := session.GetCard(12345)
	if err != nil {
		t.Errorf("Failed to unmarshal card: %s", err)
		return
	}

	if card.CardId != 12345 {
		t.Error(`Expected CardId == 12345`)
	}
	if card.Type != "CategoryCard" {
		t.Error(`Expected Type == "CategoryCard"`)
	}
	if card.LifecycleStatus != "ACTIVATED" {
		t.Error(`Expected LifecycleStatus == "ACTIVATED"`)
	}
	if card.Status != "TURNED_ON" {
		t.Error(`Expected Status == "TURNED_ON"`)
	}
	if card.Expiration != "1221" {
		t.Error(`Expected Expiration == "1221"`)
	}
	if card.LastFour != "1234" {
		t.Error(`Expected LastFour == "1234"`)
	}
	if card.VirtualCard != false {
		t.Error(`Expected VirtualCard == false`)
	}
	if card.Alias != "My Card" {
		t.Error(`Expected Alias == "My Card"`)
	}
	if card.AvailableAmount != 123.45 {
		t.Error(`Expected AvailableAmount == 123.45`)
	}
	if len(card.AllowedDays) != 1 || card.AllowedDays[0] != "MONDAY" {
		t.Error(`Expected AllowedDays == "MONDAY"`)
	}
	if card.AllowedCategoriesActive != true {
		t.Error(`Expected AllowedCategoriesActive == true`)
	}
	if len(card.AllowedCategories) != 1 ||
		card.AllowedCategories[0].TransactionCategoryId != 10 {
		t.Error(`Expected AllowedCategories == [bento.Category{TransactionCategoryId: 10}]`)
	}
	if card.CreatedOn != 1495759408 {
		t.Error(`Expected CreatedOn == 1495759408`)
	}
	if card.UpdatedOn != 1495759408 {
		t.Error(`Expected UpdatedOn == 1495759408`)
	}
	if card.SpendingLimit.Active != true {
		t.Error(`Expected SpendingLimit.Active == true`)
	}
	if card.SpendingLimit.Amount != 123.45 {
		t.Error(`Expected SpendingLimit.Amount == 123.45`)
	}
	if card.SpendingLimit.Period != "Day" {
		t.Error(`Expected SpendingLimit.Period == "Day"`)
	}
	if card.SpendingLimit.CustomStartDate != 1495759408 {
		t.Error(`Expected SpendingLimit.CustomStartDate == 1495759408`)
	}
	if card.SpendingLimit.CustomEndDate != 1495759408 {
		t.Error(`Expected SpendingLimit.CustomEndDate == 1495759408`)
	}

	if card.User.FirstName != "John" {
		t.Error(`Expected card.User.FirstName == "John"`)
	}
	if card.User.LastName != "Smith" {
		t.Error(`Expected card.User.LastName == "Smith"`)
	}
	if card.User.BirthDate != 1495759408 {
		t.Error(`Expected card.User.BirthDate == 1495759408`)
	}
	if card.User.Email != "me@myemail.com" {
		t.Error(`Expected card.User.Email == "me@myemail.com"`)
	}
	if card.User.Phone != "9998887654" {
		t.Error(`Expected card.User.Phone == "9998887654"`)
	}
	if card.User.UserId != 12345 {
		t.Error(`Expected card.User.UserId == 12345`)
	}
	if card.User.MobileAccess != true {
		t.Error(`Expected card.User.MobileAccess == true`)
	}
	if card.User.Deleted != false {
		t.Error(`Expected card.User.Deleted == false`)
	}
	if card.User.Created != 1495759408 {
		t.Error(`Expected card.User.Created == 1495759408`)
	}
}

func TestNewCard(t *testing.T) {
	t.Log("TestNewCard")
	
	session := &TestSession{}
	session.requester = testRequest(session)
	_, err := session.NewCard(EMPLOYEE_CARD, "Testing Card")
	if err != nil {
		t.Error("Unable to unmarshal card:", err)
		return
	}

	if session.method != "POST" {
		t.Error(`Expected POST request.`)
	}
	if session.endpoint != "/cards" {
		t.Error(`Expected endpoint "/cards".`)
	}
}
