# bento
--
    import "github.com/knusbaum/bento-go"

Package bento provides an interface to the bentoforbusiness.com API. Please see
https://apidocs.bentoforbusiness.com/ for detailed info on how the API works.


### Basics

In order to begin interacting with Bento, you need to obtain a *Session This can
be done via the following functions: For production:

    session, err := bento.GetProductionSession("myAccessKey", "mySecretKey")

For the sandbox:

    session, err := bento.GetTestSession("myTestAccessKey", "myTestSecretKey")

Once you have a session, you can begin doing things like getting and updating
cards and other resources associated with your bento account.

See the Session object's methods to see the types of objects you can interact
with. This is a good starting point from which you can begin to understand the
other types provided in this package.

## Usage

```go
const (
	PERIOD_DAY    = "Day"
	PERIOD_WEEK   = "Week"
	PERIOD_MONTH  = "Month"
	PERIOD_CUSTOM = "Custom"
)
```
Valid periods for SpendingLimit

```go
const (
	STATUS_CANCELED           = "CANCELED"
	STATUS_FRAUD_PREVENTION   = "FRAUD_PREVENTION"
	STATUS_TURNED_ON          = "TURNED_ON"
	STATUS_TURNED_OFF         = "TURNED_OFF"
	STATUS_WEEKLY_RESTRICTION = "WEEKLY_RESTRICTION"
)
```
Valid Card Statuses

#### type Address

```go
type Address struct {
	Active      bool        `json:"active"`
	AddressType AddressType `json:"addressType,omitempty"`
	City        string      `json:"city,omitempty"`
	Id          int64       `json:"id,omitempty"`
	State       string      `json:"state,omitempty"`
	Street      string      `json:"street,omitempty"`
	ZipCode     string      `json:"zipCode,omitempty"`
}
```


#### type AddressType

```go
type AddressType string
```

AddressType can be "BUSINESS_ADDRESS" or "USER_ADDRESS"

```go
const (
	BUSINESS_ADDRESS AddressType = "BUSINESS_ADDRESS"
	USER_ADDRESS     AddressType = "USER_ADDRESS"
)
```
Valid values for AddressType

#### type ApiApplication

```go
type ApiApplication struct {
	ApiApplicationId int64    `json:"apiApplicationId,omitempty"`
	Name             string   `json:"name,omitempty"`
	AccessKey        string   `json:"accessKey,omitempty"`
	Business         Business `json:"business,omitempty"`
}
```


#### type BentoError

```go
type BentoError struct {
	Message    string
	BentoError string `json:"error"`
}
```


#### func (BentoError) Error

```go
func (e BentoError) Error() string
```

#### type Business

```go
type Business struct {
	BusinessId        int64     `json:"businessId,omitempty"`
	CompanyName       string    `json:"companyName,omitempty"`
	NameOnCard        string    `json:"nameOnCard,omitempty"`
	Phone             string    `json:"phone,omitempty"`
	AccountNumber     string    `json:"accountNumber,omitempty"`
	BusinessStructure string    `json:"businessStructure,omitempty"`
	Status            string    `json:"status,omitempty"`
	CreatedDate       int64     `json:"createdDate,omitempty"`
	ApprovalDate      int64     `json:"approvalDate,omitempty"`
	ApprovalStatus    string    `json:"approvalStatus,omitempty"`
	Balance           float64   `json:"balance,omitempty"`
	TimeZone          string    `json:"timeZone,omitempty"`
	Addresses         []Address `json:"addresses,omitempty"`
}
```


#### type Card

```go
type Card struct {
	CardId                  int64           `json:"cardId,omitempty"`
	Type                    CardType        `json:"type,omitempty"`
	LifecycleStatus         string          `json:"lifecycleStatus,omitempty"`
	Status                  string          `json:"status,omitempty"`
	Expiration              string          `json:"expiration,omitempty"`
	LastFour                string          `json:"lastFour,omitempty"`
	VirtualCard             bool            `json:"virtualCard"`
	Alias                   string          `json:"alias,omitempty"`
	AvailableAmount         float64         `json:"availableAmount,omitempty"`
	AllowedDaysActive       bool            `json:"allowedDaysActive"`
	AllowedDays             []string        `json:"allowedDays,omitempty"`
	AllowedCategoriesActive bool            `json:"allowedCategoriesActive"`
	AllowedCategories       []Category      `json:"allowedCategories,omitempty"`
	TransactionCategoryId   int64           `json:"transactionCategoryId,omitempty"`
	CreatedOn               int64           `json:"createdOn,omitempty"`
	UpdatedOn               int64           `json:"updatedOn,omitempty"`
	SpendingLimit           SpendingLimit   `json:"spendingLimit,omitempty"`
	User                    User            `json:"user,omitempty"`
	Permissions             map[string]bool `json:"permissions,omitempty"`
	BentoType               string          `json:"bentoType,omitempty"`
}
```


#### func (*Card) Activate

```go
func (card *Card) Activate(lastFour string) (*Card, error)
```

#### func (*Card) Delete

```go
func (card *Card) Delete() (*Card, error)
```

#### func (*Card) GetBillingAddress

```go
func (card *Card) GetBillingAddress() (*Address, error)
```

#### func (*Card) GetPanAndCvv

```go
func (card *Card) GetPanAndCvv() (*PanAndCvv, error)
```

#### func (*Card) Put

```go
func (card *Card) Put() (*Card, error)
```

#### func (*Card) Reissue

```go
func (card *Card) Reissue() (*Card, error)
```

#### func (*Card) SetBillingAddress

```go
func (card *Card) SetBillingAddress(newAddress *Address) (*Address, error)
```

#### func (*Card) TurnOff

```go
func (card *Card) TurnOff() (*Card, error)
```

#### func (*Card) TurnOn

```go
func (card *Card) TurnOn() (*Card, error)
```

#### func (*Card) UpdateBillingAddress

```go
func (card *Card) UpdateBillingAddress(newAddress *Address) (*Address, error)
```

#### type CardType

```go
type CardType string
```


```go
const (
	BUSINESS_OWNER_CARD CardType = "BusinessOwnerCard"
	EMPLOYEE_CARD       CardType = "EmployeeCard"
	CATEGORY_CARD       CardType = "CategoryCard"
)
```
Valid Card Types

#### type Category

```go
type Category struct {
	TransactionCategoryId int64   `json:"transactionCategoryId,omitempty"`
	Description           string  `json:"description,omitempty"`
	Group                 string  `json:"group,omitempty"`
	Mccs                  []int64 `json:"mccs,omitempty"`
	Name                  string  `json:"name,omitempty"`
	Type                  string  `json:"type,omitempty"`
	BentoType             string  `json:"bentoType,omitempty"`
}
```


#### type PanAndCvv

```go
type PanAndCvv struct {
	Pan string `json:"pan,omitempty"`
	Cvv string `json:"cvv,omitempty"`
}
```


#### type Payee

```go
type Payee struct {
	Name    string `json:"name,omitempty"`
	City    string `json:"city,omitempty"`
	State   string `json:"state,omitempty"`
	Country string `json:"country,omitempty"`
	Zip     string `json:"zip,omitempty"`
}
```


#### type Period

```go
type Period string
```

Period for a SpendingLimit

#### type Session

```go
type Session struct {
}
```

Session provides the entry point to interact with the API. Created with
GetProductionSession and GetTestSession.

#### func  GetProductionSession

```go
func GetProductionSession(accessKey, secretKey string) (*Session, error)
```

#### func  GetTestSession

```go
func GetTestSession(accessKey, secretKey string) (*Session, error)
```

#### func (*Session) GetBusiness

```go
func (session *Session) GetBusiness() (*Business, error)
```

#### func (*Session) GetCard

```go
func (session *Session) GetCard(cardId int64) (*Card, error)
```

#### func (*Session) GetCards

```go
func (session *Session) GetCards() ([]Card, error)
```

#### func (*Session) GetTransactions

```go
func (session *Session) GetTransactions() (*Transactions, error)
```

#### func (*Session) NewCard

```go
func (session *Session) NewCard(cardType CardType, alias string) (*Card, error)
```

#### type SpendingLimit

```go
type SpendingLimit struct {
	Active          bool    `json:"active"`
	Amount          float64 `json:"amount,omitempty"`
	Period          Period  `json:"period,omitempty"`
	CustomStartDate int64   `json:"customStartDate,omitempty"`
	CustomEndDate   int64   `json:"customEndDate,omitempty"`
}
```


#### type Transaction

```go
type Transaction struct {
	CardTransactionId int64     `json:"cardTransactionId,omitempty"`
	Amount            float64   `json:"amount,omitempty"`
	ApprovalCode      string    `json:"approvalCode,omitempty"`
	AvailableBalance  float64   `json:"availableBalance,omitempty"`
	Card              *Card     `json:"card,omitempty"`
	Category          *Category `json:"category,omitempty"`
	Currency          string    `json:"currency,omitempty"`
	Deleted           bool      `json:"deleted,omitempty"`
	Fees              float64   `json:"fees,omitempty"`
	LedgerBalance     float64   `json:"ledgerBalance,omitempty"`
	Note              string    `json:"node,omitempty"`
	SettlementDate    int64     `json:"settlementDate,omitempty"`
	Status            string    `json:"status,omitempty"`
	Tags              []string  `json:"tags,omitempty"`
	TransactionDate   int64     `json:"transactionDate,omitempty"`
	Type              string    `json:"type,omitempty"`
	Payee             *Payee    `json:"payee,omitempty"`
}
```


#### type Transactions

```go
type Transactions struct {
	Amount           float64       `json:"amount,omitempty"`
	Size             int           `json:"size",omitempty"`
	CardTransactions []Transaction `json:"cardTransactions"`
}
```

Transactions

#### type User

```go
type User struct {
	FirstName    string `json:"firstName,omitempty"`
	LastName     string `json:"lastName,omitempty"`
	BirthDate    int64  `json:"birthDate,omitempty"`
	Email        string `json:"email,omitempty"`
	Phone        string `json:"phone,omitempty"`
	UserId       int64  `json:"userId,omitempty"`
	MobileAccess bool   `json:"mobileAccess"`
	Deleted      bool   `json:"deleted"`
	Created      int64  `json:"created"`
	BentoType    string `json:"bentoType,omitempty"`
}
```
