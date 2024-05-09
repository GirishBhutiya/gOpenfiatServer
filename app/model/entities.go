package model

import (
	"time"

	"github.com/google/uuid"
)

var (
	// Constants for the order status
	ORDER_STATUS_PENDING    = "pending"
	ORDER_STATUS_CONFIRM    = "confirm"
	ORDER_STATUS_DISPUTED   = "disputed"
	ORDER_STATUS_CONFIRMING = "confirming"
	ORDER_TYPE_BUY          = "buy"
	ORDER_TYPE_SELL         = "sell"
)

// swagger:model UserLogin
type UserLogin struct {
	// phonenumber of the user
	// in: integer
	PhoneNumber int `json:"phonenumber,omitempty"`
}

// swagger:model User
type User struct {
	// id of the user
	// in: integer
	ID uuid.UUID `json:"userid,omitempty"`
	// phonenumber of the user
	// in: integer
	PhoneNumber int `json:"phonenumber,omitempty"`
	// first name of the user
	// in: string
	FirstName string `json:"first_name,omitempty"`
	// last name of the user
	// in: string
	LastName string `json:"last_name,omitempty"`
	//Is users phone number verified
	// in: boolean
	Verified bool `json:"verified"`
	// profile pic link of the user
	// in: string
	ProfilePic string `json:"profile_pic"`
}

// swagger:model User
type UserUpdate struct {
	// id of the user
	// in: integer
	ID uuid.UUID `json:"userid,omitempty"`
	// phonenumber of the user
	// in: integer
	PhoneNumber int `json:"phonenumber,omitempty"`
	// first name of the user
	// in: string
	FirstName string `json:"first_name,omitempty"`
	// last name of the user
	// in: string
	LastName string `json:"last_name,omitempty"`
	//Is users phone number verified
	// in: boolean
	Verified bool `json:"verified"`

	// profile pic link of the user
	// in: string
	Base64JPGIMG string `json:"base64jpgimg"`
}
type UserGroups struct {
	Groups map[uuid.UUID]string `json:"groups"`
}

// swagger:model User
type UserHandler1 struct {
	// phonenumber of the user
	// in: integer
	PhoneNumber int `json:"phonenumber,omitempty"`
	// first name of the user
	// in: string
	FirstName string `json:"first_name,omitempty"`
	// last name of the user
	// in: string
	LastName string `json:"last_name,omitempty"`
	//Is users phone number verified
	// in: boolean
	Verified bool `json:"verified"`
	// roll of the user admin or user
	// in: integer
	RollId int `json:"roll_id"`
	// profile pic link of the user
	// in: string
	ProfilePic string `json:"profile_pic"`
}

type GroupUser struct {
	// id of the group
	// in: uuid
	ID uuid.UUID `json:"groupid"`
	// name of the group
	// in: string
	Name string `json:"groupname"`
}

// swagger:model CreateGroup
type CreateGroup struct {
	// name of the group
	// in: string
	Name string `json:"groupname"`
}
type GroupHandler struct {
	ID            uuid.UUID `json:"groupid"`
	Name          string    `json:"name"`
	CreatorUserid uuid.UUID `json:"createruserid"`
	CreateTime    time.Time `json:"createtime"`
}
type OrderHandler struct {
	ID         uuid.UUID `json:"orderid"`
	UserId     uuid.UUID `json:"userid"`
	FiatAmount float32   `json:"fiatAmount,omitempty"`
	MinAmount  float32   `json:"minAmount,omitempty"`
	Price      float32   `json:"price,omitempty"`
	TimeLimit  int64     `json:"timeLimit,omitempty"`
	Type       string    `json:"type"`
}
type OrderHandlerString struct {
	ID         uuid.UUID `json:"orderid"`
	UserId     uuid.UUID `json:"userid"`
	FiatAmount float32   `json:"fiatAmount,omitempty"`
	MinAmount  float32   `json:"minAmount,omitempty"`
	Price      float32   `json:"price,omitempty"`
	TimeLimit  string    `json:"timeLimit,omitempty"`
	Type       string    `json:"type"`
}
type Order struct {
	ID         uuid.UUID `json:"orderid"`
	FiatAmount float32   `json:"fiatAmount,omitempty"`
	MinAmount  float32   `json:"minAmount,omitempty"`
	Price      float32   `json:"price,omitempty"`
	TimeLimit  string    `json:"timeLimit,omitempty"`
	Type       string    `json:"type"`
}

// swagger:model OrderUser
type OrderUser struct {
	// id of the order
	// in: uuid
	ID uuid.UUID `json:"orderid"`
}

/* type OrderHandler struct {
	// id of the order
	// in: uuid
	ID uuid.UUID `json:"id"`
	// phonenumber of the user who placed the order
	// in: integer
	Phonenumber int `json:"phonenumber,omitempty"`
	// flat amount of the order
	// in: integer
	FlatAmount float32 `json:"flat_amount,omitempty"`
	// minimum amount of the order
	// in: integer
	MinAmount float32 `json:"min_amount,omitempty"`
	// price of the order
	// in: integer
	Price float32 `json:"price,omitempty"`
	// status of the order
	// in: string
	Status string `json:"status"`
	// time limit of the order
	// in: string
	TimeLimit string `json:"time_limit"`
	//type of the order buy/sell
	// in: string
	OrderType string `json:"order_type"`
}
type Order struct {
	// id of the order
	// in: uuid
	ID uuid.UUID `json:"id"`
	// amount of the order
	// in: integer
	MinAmount float32 `json:"min_amount,omitempty"`
	// time limit of the order
	// in: string
	TimeLimit string `json:"time_limit"`
}
*/

type TradeHandler struct {
	ID        uuid.UUID `json:"tradeid"`
	Orderid   uuid.UUID `json:"orderid"`
	BidUserid uuid.UUID `json:"bidUserid"`
	TradeTime int64     `json:"tradetime"`
	Status    string    `json:"status"`
	Method    string    `json:"method"`
}
type TradeHandlerUser struct {
	ID        uuid.UUID `json:"tradeid"`
	Orderid   uuid.UUID `json:"orderid"`
	TradeTime int64     `json:"tradetime"`
	Method    string    `json:"method"`
}
type UserTrades struct {
	ID        uuid.UUID `json:"tradeid"`
	Orderid   uuid.UUID `json:"orderid"`
	TradeTime int64     `json:"tradetime"`
	Status    string    `json:"status"`
	Method    string    `json:"method"`
}

// swagger:model CreateGroup
type TradeUser struct {
	// id of the trade
	// in: uuid
	ID uuid.UUID `json:"tradeid"`
}

// swagger:model InviteLink
type InviteLink struct {
	InviteLink string `json:"invitelink"`
}
