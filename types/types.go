package types

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
