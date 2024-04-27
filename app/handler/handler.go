package handler

import (
	"encoding/json"
	"errors"

	"log"
	"net/http"
	"time"

	"github.com/GirishBhutiya/gOpenfiatServer/app/config"
	"github.com/GirishBhutiya/gOpenfiatServer/app/model"
	"github.com/GirishBhutiya/gOpenfiatServer/app/token"
	util "github.com/GirishBhutiya/gOpenfiatServer/app/util"
	"github.com/GirishBhutiya/gOpenfiatServer/db"
	"github.com/google/uuid"
)

type Server struct {
	Config config.Config
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

var ser *Server

func InitServer(server *Server) {
	log.Println("InitServer")
	ser = server
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
	User                 model.User `json:"user"`
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
// @Tags			Brocker
// @Accept			json
// @Produce		json
// @Required false
// @Success 200 {object} jsonResponse
// @Failure 401 {object} jsonResponse
// @Router / [get]
func Brocker(w http.ResponseWriter, r *http.Request) {

	payload := jsonResponse{
		Error:   false,
		Message: "Hit the brocker",
	}

	_ = WriteJSON(w, http.StatusOK, payload)

}

// Register godoc
// This API is used to register user with Phone number
// @Summary This API is used to register user with Phone number
// @Description This API is used to register user with Phone number
// @Tags			User
// @Accept			json
// @Produce			json
// @Param user body types.User true "User"
// @Success 200 {object} jsonResponse
// @Failure 401 {object} jsonResponse
// @Router /register [post]
func Register(w http.ResponseWriter, r *http.Request) {
	var user model.User

	err := ReadJSON(w, r, &user)
	if err != nil {
		log.Println(err)
		ErrorJSON(w, err)
		return
	}

	//fmt.Println("\nPhone is:", user.PhoneNumber)
	if util.LenLoop(user.PhoneNumber) < 10 {
		ErrorJSON(w, errors.New("please enter valid phonenumber"))
		return
	}
	//TODO: send OTP
	var otpRes OTPResponse
	if ser.Config.Production {
		data, err := util.SendOTP(user.PhoneNumber, &ser.Config)
		if err != nil {
			log.Println(err)
			ErrorJSON(w, err)
			return
		}

		err = json.Unmarshal(data, &otpRes)
		if err != nil {
			log.Println(err)
			ErrorJSON(w, err)
			return
		}
		//fmt.Println("OTP is:", otpRes.OTP)
		if otpRes.Status == "Error" {
			log.Println(err)
			ErrorJSON(w, errors.New(otpRes.Details))
			return
		}
	} else {
		otpRes.OTP = 123456
	}

	//Insert to database
	err = ser.Store.InserLoginRegister(otpRes.OTP, &user)
	//err = app.Store.InserLoginRegister(otpRes.OTP, &user)
	if err != nil {
		log.Println(err)
		ErrorJSON(w, err)
		return
	}
	var payload jsonResponse
	payload.Error = false
	payload.Message = "Please check inbox to see  OTP"

	WriteJSON(w, http.StatusOK, payload)

	//WriteJSON(w, http.StatusAccepted, res)

}

// VerifyOTP
// This API used to verify OTP which you get after register
// @Summary		Verify OTP
// @Description		Verify OTP which you get after register
// @Accept			json
// @Produce			json
// @Param			payload		body		AuthPayload	true	"payload"
// @Tags			Auth
// @Success			200		{object}	jsonResponse
// @Failure 401 {object} jsonResponse
// @Router /verifyotp [post]
func VerifyOTP(w http.ResponseWriter, r *http.Request) {

	var payload AuthPayload

	err := ReadJSON(w, r, &payload)
	if err != nil {
		log.Println(err)
		ErrorJSON(w, err)
		return
	}
	if util.LenLoop(payload.OTP) != 6 {
		ErrorJSON(w, errors.New("OTP must be in 6 digit"))
	}
	if util.LenLoop(payload.PhoneNumber) < 10 {
		ErrorJSON(w, errors.New("please enter valid phonenumber"))
		return
	}
	user, err := ser.Store.VerifyOTP(payload.PhoneNumber, payload.OTP)
	if err != nil {
		log.Println(err)
		ErrorJSON(w, err)
		return
	}
	accessToken, accessPayload, err := ser.TokenMaker.CreateToken(payload.PhoneNumber, ser.Config.AccessTokenDuration)
	if err != nil {
		log.Println(err)
		ErrorJSON(w, err)
		return
	}

	var res LoginResponse
	res.Authenticated = true
	res.AccessToken = accessToken
	res.AccessTokenExpiresAt = accessPayload.ExpiredAt
	res.User = user

	WriteJSON(w, http.StatusOK, res)

}

// UpdateUser : Update user
// This API is used to update user profile like First Name, Last Name etc
// @Summary Update User profile
// @Description This API is used to update user profile like First Name, Last Name etc
// @Tags User
// @Accept  json
// @Produce  json
// @Param user body types.User true "User"
// @Success 200 {object} jsonResponse
// @Failure 401 {object} jsonResponse
// @Router /user/update [post]
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user model.User

	err := ReadJSON(w, r, &user)
	if err != nil {
		log.Println(err)
		ErrorJSON(w, err)
		return
	}
	//fmt.Println("Phone is:", user.PhoneNumber)
	//TODO: send OTP
	if util.LenLoop(user.PhoneNumber) < 10 {
		ErrorJSON(w, errors.New("please enter valid phonenumber"))
		return
	}
	//Update to database
	_, err = ser.Store.UpdateUser(&user)
	if err != nil {
		log.Println(err)
		ErrorJSON(w, err)
		return
	}
	var payload jsonResponse
	payload.Error = false
	payload.Message = "User Updated"

	WriteJSON(w, http.StatusOK, payload)

	//WriteJSON(w, http.StatusAccepted, res)

}

// DeleteUser deletes a user
// This API is used to delete a user
// @Summary Delete User
// @Description Delete a user
// @Tags User
// @Accept  json
// @Produce  json
// @Param user body types.User true "User"
// @Success 200 {object} jsonResponse
// @Failure 401 {object} jsonResponse
// @Router /user/delete [post]
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	var user model.User

	err := ReadJSON(w, r, &user)
	if err != nil {
		log.Println(err)
		ErrorJSON(w, err)
		return
	}
	if util.LenLoop(user.PhoneNumber) < 10 {
		ErrorJSON(w, errors.New("please enter valid phonenumber"))
		return
	}
	//fmt.Println("Phone is:", user.PhoneNumber)
	//TODO: send OTP

	//Update to database
	err = ser.Store.DeleteUser(&user)
	if err != nil {
		log.Println(err)
		ErrorJSON(w, err)
		return
	}
	var payload jsonResponse
	payload.Error = false
	payload.Message = "Account Deleted"

	WriteJSON(w, http.StatusOK, payload)

	//WriteJSON(w, http.StatusAccepted, res)

}

// CreateOrder : Create a new order
// CreateOrder : Create a new order with status pending
// @Summary Create a new order
// @Description Create a new order with status pending
// @Tags Order
// @Accept  json
// @Produce  json
// @Param order body types.Order true "Order"
// @Success 200 {object} jsonResponse
// @Failure 401 {object} jsonResponse
// @Router /user/create-order [post]
func CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order model.Order

	err := ReadJSON(w, r, &order)
	if err != nil {
		log.Println(err)
		ErrorJSON(w, err)
		return
	}

	if order.Amount == 0 {
		log.Println(err)
		ErrorJSON(w, errors.New("amount value 0"))
		return
	}

	if order.FromPhone == 0 || util.LenLoop(order.FromPhone) < 10 {
		log.Println(err)
		ErrorJSON(w, errors.New("please check from phone number"))
		return
	}
	if order.ToPhone == 0 || util.LenLoop(order.ToPhone) < 10 {
		log.Println(err)
		ErrorJSON(w, errors.New("please check to phonenumber"))
		return
	}
	err = ser.Store.CreateNewOrder(&order)
	if err != nil {
		log.Println(err)
		ErrorJSON(w, err)
		return
	}
	var payload jsonResponse
	payload.Error = false
	payload.Message = "Order Created"

	WriteJSON(w, http.StatusOK, payload)

}

// UpdateOrderValue godoc
// This Api is used to update the order amount
// @Summary Update Order Value
// @Description This Api is used to update the order amount
// @Tags Order
// @Accept  json
// @Produce  json
// @Param order body types.Order true "Order"
// @Success 200 {object} jsonResponse
// @Failure 401 {object} jsonResponse
// @Router /user/update-ordervalue [post]
func UpdateOrderValue(w http.ResponseWriter, r *http.Request) {
	var order model.Order

	err := ReadJSON(w, r, &order)
	if err != nil {
		log.Println(err)
		ErrorJSON(w, err)
		return
	}
	if order.ID == uuid.Nil {
		log.Println(err)
		ErrorJSON(w, errors.New("order id not found"))
		return
	}
	if order.Amount == 0 {
		log.Println(err)
		ErrorJSON(w, errors.New("amount value 0"))
		return
	}
	if order.FromPhone == 0 || util.LenLoop(order.FromPhone) < 10 {
		log.Println(err)
		ErrorJSON(w, errors.New("please check from phonenumber"))
		return
	}
	if order.ToPhone == 0 || util.LenLoop(order.ToPhone) < 10 {
		log.Println(err)
		ErrorJSON(w, errors.New("please check to phonenumber"))
		return
	}
	err = ser.Store.UpdateOrderValue(&order)
	if err != nil {
		log.Println(err)
		ErrorJSON(w, err)
		return
	}
	var payload jsonResponse
	payload.Error = false
	payload.Message = "Order Updates Sucessfully"

	WriteJSON(w, http.StatusOK, payload)

}

// ConfirmingOrder godoc
// This Api is used to change the order status to confirming
// @Summary ConfirmingOrder
// @Description This Apis is used to change the order status to confirming
// @Tags Order
// @Accept  json
// @Produce  json
// @Param order body types.Order true "Order"
// @Success 200 {object} jsonResponse
// @Failure 401 {object} jsonResponse
// @Router /user/order-confirming [post]
func ConfirmingOrder(w http.ResponseWriter, r *http.Request) {
	var order model.Order

	err := ReadJSON(w, r, &order)
	if err != nil {
		log.Println(err)
		ErrorJSON(w, err)
		return
	}
	if order.ID == uuid.Nil {
		log.Println(err)
		ErrorJSON(w, errors.New("order id not found"))
		return
	}
	if order.Amount == 0 {
		log.Println(err)
		ErrorJSON(w, errors.New("amount value 0"))
		return
	}
	if order.FromPhone == 0 || util.LenLoop(order.FromPhone) < 10 {
		log.Println(err)
		ErrorJSON(w, errors.New("please check from phonenumber"))
		return
	}
	if order.ToPhone == 0 || util.LenLoop(order.ToPhone) < 10 {
		log.Println(err)
		ErrorJSON(w, errors.New("please check to phonenumber"))
		return
	}
	err = ser.Store.ConfirmingOrder(&order)
	if err != nil {
		log.Println(err)
		ErrorJSON(w, err)
		return
	}
	var payload jsonResponse
	payload.Error = false
	payload.Message = "Order status changed to confirming"

	WriteJSON(w, http.StatusOK, payload)

}

// ConfirmOrder godoc
// This Api is used to change the order status to confirm
// @Summary ConfirmOrder
// @Description This Api is used to change the order status to confirm
// @Tags Order
// @Accept  json
// @Produce  json
// @Param order body types.Order true "Order"
// @Success 200 {object} jsonResponse
// @Failure 401 {object} jsonResponse
// @Router /user/order-confirm [post]
func ConfirmOrder(w http.ResponseWriter, r *http.Request) {
	var order model.Order

	err := ReadJSON(w, r, &order)
	if err != nil {
		log.Println(err)
		ErrorJSON(w, err)
		return
	}
	if order.ID == uuid.Nil {
		log.Println(err)
		ErrorJSON(w, errors.New("order id not found"))
		return
	}
	if order.Amount == 0 {
		log.Println(err)
		ErrorJSON(w, errors.New("amount value 0"))
		return
	}
	if order.FromPhone == 0 || util.LenLoop(order.FromPhone) < 10 {
		log.Println(err)
		ErrorJSON(w, errors.New("please check from phonenumber"))
		return
	}
	if order.ToPhone == 0 || util.LenLoop(order.ToPhone) < 10 {
		log.Println(err)
		ErrorJSON(w, errors.New("please check to phonenumber"))
		return
	}
	err = ser.Store.ConfirmOrder(&order)
	if err != nil {
		log.Println(err)
		ErrorJSON(w, err)
		return
	}
	var payload jsonResponse
	payload.Error = false
	payload.Message = "Order status changed to confirm"

	WriteJSON(w, http.StatusOK, payload)

}

// DisputedOrder godoc
// This API is used to change the order status to disputed
// @Summary DisputedOrder
// @Description This API is used to change the order status to disputed
// @Tags Order
// @Accept  json
// @Produce  json
// @Param order body types.Order true "Order"
// @Success 200 {object} jsonResponse
// @Failure 401 {object} jsonResponse
// @Router /user/order-disputed [post]
func DisputedOrder(w http.ResponseWriter, r *http.Request) {
	var order model.Order

	err := ReadJSON(w, r, &order)
	if err != nil {
		log.Println(err)
		ErrorJSON(w, err)
		return
	}
	if order.ID == uuid.Nil {
		log.Println(err)
		ErrorJSON(w, errors.New("order id not found"))
		return
	}
	if order.Amount == 0 {
		log.Println(err)
		ErrorJSON(w, errors.New("amount value 0"))
		return
	}
	if order.FromPhone == 0 || util.LenLoop(order.FromPhone) < 10 {
		log.Println(err)
		ErrorJSON(w, errors.New("please check from phonenumber"))
		return
	}
	if order.ToPhone == 0 || util.LenLoop(order.ToPhone) < 10 {
		log.Println(err)
		ErrorJSON(w, errors.New("please check to phonenumber"))
		return
	}
	err = ser.Store.DisputedOrder(&order)
	if err != nil {
		log.Println(err)
		ErrorJSON(w, err)
		return
	}
	var payload jsonResponse
	payload.Error = false
	payload.Message = "Order status changed to disputed"

	WriteJSON(w, http.StatusOK, payload)

}

// DeletedOrder godoc
// This API is used to delete an order
// @Summary Delete an order
// @Description Delete an order
// @Tags Order
// @Accept  json
// @Produce  json
// @Param order body types.Order true "Order"
// @Success 200 {object} jsonResponse
// @Failure 401 {object} jsonResponse
// @Router /user/order-delete [post]
func DeleteOrder(w http.ResponseWriter, r *http.Request) {
	var order model.Order

	err := ReadJSON(w, r, &order)
	if err != nil {
		log.Println(err)
		ErrorJSON(w, err)
		return
	}
	if order.ID == uuid.Nil {
		log.Println(err)
		ErrorJSON(w, errors.New("order id not found"))
		return
	}
	if order.Amount == 0 {
		log.Println(err)
		ErrorJSON(w, errors.New("amount value 0"))
		return
	}
	if order.FromPhone == 0 || util.LenLoop(order.FromPhone) < 10 {
		log.Println(err)
		ErrorJSON(w, errors.New("please check from phonenumber"))
		return
	}
	if order.ToPhone == 0 || util.LenLoop(order.ToPhone) < 10 {
		log.Println(err)
		ErrorJSON(w, errors.New("please check to phonenumber"))
		return
	}
	//Update to database
	err = ser.Store.DeleteOrder(&order)
	if err != nil {
		log.Println(err)
		ErrorJSON(w, err)
		return
	}
	var payload jsonResponse
	payload.Error = false
	payload.Message = "Order Deleted"

	WriteJSON(w, http.StatusOK, payload)

	//WriteJSON(w, http.StatusAccepted, res)

}

// GetUserAllOrders godoc
// This API is used to get all orders of a user
// @Summary Get all orders of a user
// @Description This API is used to get all orders of a user
// @Tags Order
// @Accept  json
// @Produce  json
// @Param user body types.User true "User"
// @Success 200 {object} jsonResponse
// @Failure 401 {object} jsonResponse
// @Router /user/allorders [post]
func GetUserAllOrders(w http.ResponseWriter, r *http.Request) {
	var user model.User

	err := ReadJSON(w, r, &user)
	if err != nil {
		log.Println(err)
		ErrorJSON(w, err)
		return
	}
	if util.LenLoop(user.PhoneNumber) < 10 {
		ErrorJSON(w, errors.New("please enter valid phonenumber"))
		return
	}
	orders, err := ser.Store.GetAllOrders(&user)
	if err != nil {
		log.Println(err)
		ErrorJSON(w, err)
		return
	}
	var payload jsonResponse
	payload.Error = false
	payload.Message = "Order Deleted"

	WriteJSON(w, http.StatusOK, orders)
}
