package api

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/swisscdn/OpenfiatServer/db"
	"github.com/swisscdn/OpenfiatServer/token"
	"github.com/swisscdn/OpenfiatServer/types"
	util "github.com/swisscdn/OpenfiatServer/utils"
)

type Server struct {
	Config util.Config
	//AstraClient *astra.Client
	TokenMaker token.Maker
	Store      db.DatabaseService
}
type AuthPayload struct {
	// phonenumber of the user
	// in: integer
	PhoneNumber int `json:"phonenumber"`
	// otp of the user which get in phone
	// in: integer
	OTP int `json:"otp"`
}

/*
	 type UserResponse struct {
		PhoneNumber int    `json:"phonenumber"`
		FirstName   string `json:"first_name,omitempty"`
		LastName    string `json:"last_name,omitempty"`
		RollId      int    `json:"roll_id"`
		Verified    int    `json:"verified"`
		Roll        string `json:"roll"`
	}
*/
type LoginResponse struct {
	User                 types.User `json:"user"`
	AccessToken          string     `json:"access_token"`
	AccessTokenExpiresAt time.Time  `json:"access_token_expires_at"`
	Authenticated        bool       `json:"authenticated"`
	ErrorMessage         string     `json:"message"`
}
type OTPResponse struct {
	Status  string `json:"status"`
	Details string `json:"Details"`
	OTP     int    `json:"OTP"`
}

// Brocker godoc
// This API can be used as health check for this application.
// @Summary This API can be used as health check for this application.
// @Description This API can be used as health check for this application.
// @Tags			Groups
// @Accept			json
// @Produce		json
func (app *Server) Brocker(w http.ResponseWriter, r *http.Request) {

	payload := jsonResponse{
		Error:   false,
		Message: "Hit the brocker",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)

}
func (app *Server) Register(w http.ResponseWriter, r *http.Request) {
	var user types.User

	err := app.readJSON(w, r, &user)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}
	fmt.Println("\nPhone is:", user.PhoneNumber)
	//TODO: send OTP
	/* data, err := util.SendOTP(user.PhoneNumber, &app.Config)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}
	var otpRes OTPResponse
	err = json.Unmarshal(data, &otpRes)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}
	fmt.Println("OTP is:", otpRes.OTP)
	if otpRes.Status == "Error" {
		log.Println(err)
		app.errorJSON(w, errors.New(otpRes.Details))
		return
	} */
	//Insert to database
	err = app.Store.InserLoginRegister(123456, &user)
	//err = app.Store.InserLoginRegister(otpRes.OTP, &user)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}
	var payload jsonResponse
	payload.Error = false
	payload.Message = "Please check inbox to see  OTP"

	app.writeJSON(w, http.StatusAccepted, payload)

	//app.writeJSON(w, http.StatusAccepted, res)

}
func (app *Server) VerifyOTP(w http.ResponseWriter, r *http.Request) {

	var payload AuthPayload

	err := app.readJSON(w, r, &payload)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}
	if util.LenLoop(payload.OTP) != 6 {
		app.errorJSON(w, errors.New("OTP must be in 6 digit"))
	}
	user, err := app.Store.VerifyOTP(payload.PhoneNumber, payload.OTP)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}
	accessToken, accessPayload, err := app.TokenMaker.CreateToken(payload.PhoneNumber, app.Config.AccessTokenDuration)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}

	var res LoginResponse
	res.Authenticated = true
	res.AccessToken = accessToken
	res.AccessTokenExpiresAt = accessPayload.ExpiredAt
	res.User = user

	app.writeJSON(w, http.StatusAccepted, res)

}
func (app *Server) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user types.User

	err := app.readJSON(w, r, &user)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}
	fmt.Println("Phone is:", user.PhoneNumber)
	//TODO: send OTP

	//Update to database
	_, err = app.Store.UpdateUser(&user)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}
	var payload jsonResponse
	payload.Error = false
	payload.Message = "User Updated"

	app.writeJSON(w, http.StatusAccepted, payload)

	//app.writeJSON(w, http.StatusAccepted, res)

}
func (app *Server) DeleteUser(w http.ResponseWriter, r *http.Request) {
	var user types.User

	err := app.readJSON(w, r, &user)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}
	fmt.Println("Phone is:", user.PhoneNumber)
	//TODO: send OTP

	//Update to database
	err = app.Store.DeleteUser(&user)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}
	var payload jsonResponse
	payload.Error = false
	payload.Message = "Account Deleted"

	app.writeJSON(w, http.StatusAccepted, payload)

	//app.writeJSON(w, http.StatusAccepted, res)

}
