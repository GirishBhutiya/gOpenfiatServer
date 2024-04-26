package types

import "github.com/google/uuid"

var (
	// Constants for the order status
	ORDER_STATUS_PENDING    = "pending"
	ORDER_STATUS_CONFIRM    = "confirm"
	ORDER_STATUS_DISPUTED   = "disputed"
	ORDER_STATUS_CONFIRMING = "confirming"
)

// swagger:model User
type User struct {
	// phonenumber of the user
	// in: integer
	PhoneNumber int `json:"phonenumber"`
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

type Order struct {
	// id of the order
	// in: uuid
	ID uuid.UUID `json:"id"`
	// phonenumber of the user who placed the order
	// in: integer
	FromPhone int `json:"from_phonenumber"`
	// phonenumber of the user who is receiving the order
	// in: integer
	ToPhone int `json:"to_phonenumber"`
	// amount of the order
	// in: integer
	Amount float32 `json:"amount"`
	// status of the order
	// in: string
	Status string `json:"status"`
	/* // time of the order
	// in: string
	Time string `json:"time"` */
}
