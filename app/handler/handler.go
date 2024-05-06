package handler

import (
	"encoding/json"
	"errors"

	"log"
	"net/http"
	"time"

	"github.com/GirishBhutiya/gOpenfiatServer/app/config"
	db "github.com/GirishBhutiya/gOpenfiatServer/app/database"
	"github.com/GirishBhutiya/gOpenfiatServer/app/model"
	"github.com/GirishBhutiya/gOpenfiatServer/app/token"
	util "github.com/GirishBhutiya/gOpenfiatServer/app/util"
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
type contextKey string

const UserIDKey contextKey = "userid"

var ser *Server

func InitServer(server *Server) {
	ser = server
}

type renewAccessTokenRequests struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type renewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
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
	User                  model.User `json:"user"`
	AccessToken           string     `json:"access_token"`
	AccessTokenExpiresAt  time.Time  `json:"access_token_expires_at"`
	RefreshToken          string     `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time  `json:"refresh_token_expires_at"`
	Authenticated         bool       `json:"authenticated"`
	ErrorMessage          string     `json:"message"`
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
// @Param user body model.UserLogin true "jsonResponse"
// @Success 200 {object} jsonResponse
// @Failure 401 {object} jsonResponse
// @Router /register [post]
func Register(w http.ResponseWriter, r *http.Request) {
	var user model.UserLogin

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
	if err != nil {
		log.Println(err)
		ErrorJSON(w, err)
		return
	}
	var payload jsonResponse
	payload.Error = false
	payload.Message = "Please check inbox to see  OTP"

	WriteJSON(w, http.StatusOK, payload)

}

// Login
// This API used to verify OTP which you get after register
// @Summary		Login
// @Description		Verify OTP which you get after register
// @Accept			json
// @Produce			json
// @Param			payload		body		AuthPayload	true	"LoginResponse"
// @Tags			Auth
// @Success			200		{object}	jsonResponse
// @Failure 401 {object} jsonResponse
// @Router /login [post]
func Login(w http.ResponseWriter, r *http.Request) {

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
	user, err := ser.Store.Login(payload.PhoneNumber, payload.OTP)
	if err != nil {
		log.Println(err)
		ErrorJSON(w, err)
		return
	}
	accessToken, accessPayload, err := ser.TokenMaker.CreateToken(user.ID, ser.Config.AccessTokenDuration)
	if err != nil {
		log.Println(err)
		ErrorJSON(w, err)
		return
	}
	refreshToken, refreshPayload, err := ser.TokenMaker.CreateToken(user.ID, ser.Config.RefreshTokenDuration)
	if err != nil {
		log.Println(err)
		ErrorJSON(w, err)
		return
	}
	var res LoginResponse
	res.Authenticated = true
	res.AccessToken = accessToken
	res.AccessTokenExpiresAt = accessPayload.ExpiredAt
	res.RefreshToken = refreshToken
	res.RefreshTokenExpiresAt = refreshPayload.ExpiredAt
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
// @Param user body model.User true "jsonResponse"
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
	userid, ok := r.Context().Value(UserIDKey).(uuid.UUID)
	if !ok {
		ErrorJSON(w, errors.New("can not get userid"))
		return
	}
	user.ID = userid
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

// SubscribeGroupToUSer : subscribe to group
// This API is used to subscribe to group
// @Summary Subscribe Group To USer
// @Description This API is used to subscribe to group
// @Tags User
// @Accept  json
// @Produce  json
// @Param user body model.GroupUser true "jsonResponse"
// @Success 200 {object} jsonResponse
// @Failure 401 {object} jsonResponse
// @Router /user/subscribe [post]
func SubscribeGroupToUSer(w http.ResponseWriter, r *http.Request) {
	var group model.GroupUser
	err := ReadJSON(w, r, &group)
	if err != nil {
		log.Println(err)
		ErrorJSON(w, err)
		return
	}
	userid, ok := r.Context().Value(UserIDKey).(uuid.UUID)
	if !ok {
		ErrorJSON(w, errors.New("can not get userid"))
		return
	}
	err = ser.Store.SubscribeGroupToUSer(userid, group)
	if err != nil {
		log.Println(err)
		ErrorJSON(w, err)
		return
	}
	var payload jsonResponse
	payload.Error = false
	payload.Message = "User Updated"

	WriteJSON(w, http.StatusOK, payload)
}

// UnsubscribeGroupToUSer : unsubscribe to group
// This API is used to unsubscribe to group
// @Summary Unsubscribe Group To USer
// @Description This API is used to unsubscribe to group
// @Tags User
// @Accept  json
// @Produce  json
// @Param user body model.GroupUser true "jsonResponse"
// @Success 200 {object} jsonResponse
// @Failure 401 {object} jsonResponse
// @Router /user/unsubscribe [post]
func UnsubscribeGroupToUSer(w http.ResponseWriter, r *http.Request) {
	var group model.GroupUser
	err := ReadJSON(w, r, &group)
	if err != nil {
		log.Println(err)
		ErrorJSON(w, err)
		return
	}
	userid, ok := r.Context().Value(UserIDKey).(uuid.UUID)
	if !ok {
		ErrorJSON(w, errors.New("can not get userid"))
		return
	}
	err = ser.Store.UnsubscribeGroupToUSer(userid, group)
	if err != nil {
		log.Println(err)
		ErrorJSON(w, err)
		return
	}
	var payload jsonResponse
	payload.Error = false
	payload.Message = "user unsubscribed to group"

	WriteJSON(w, http.StatusOK, payload)
}

// CreateGroup : create a new group
// This API is used to create a new group
// @Summary Create A New Group
// @Description This API is used to create a new group
// @Tags Group
// @Accept  json
// @Produce  json
// @Param user body model.CreateGroup true "model.GroupUser"
// @Success 200 {object} jsonResponse
// @Failure 401 {object} jsonResponse
// @Router /user/create-group [post]
func CreateGroup(w http.ResponseWriter, r *http.Request) {
	var group model.CreateGroup

	err := ReadJSON(w, r, &group)
	if err != nil {
		log.Println(err)
		ErrorJSON(w, err)
		return
	}

	userid, ok := r.Context().Value(UserIDKey).(uuid.UUID)
	if !ok {
		ErrorJSON(w, errors.New("can not get userid"))
		return
	}
	grp, err := ser.Store.CreateNewGroup(userid, group)
	if err != nil {
		log.Println(err)
		ErrorJSON(w, err)
		return
	}
	WriteJSON(w, http.StatusOK, grp)

}

// UpdateGroup : update group
// This API is used to update group
// @Summary Update Group
// @Description This API is used to update group
// @Tags Group
// @Accept  json
// @Produce  json
// @Param user body model.GroupUser true "jsonResponse"
// @Success 200 {object} jsonResponse
// @Failure 401 {object} jsonResponse
// @Router /user/update-group [post]
func UpdateGroup(w http.ResponseWriter, r *http.Request) {
	var group model.GroupUser

	err := ReadJSON(w, r, &group)
	if err != nil {
		log.Println(err)
		ErrorJSON(w, err)
		return
	}

	err = ser.Store.UpdateGroup(group)
	if err != nil {
		log.Println(err)
		ErrorJSON(w, err)
		return
	}
	var payload jsonResponse
	payload.Error = false
	payload.Message = "group updated"

	WriteJSON(w, http.StatusOK, payload)

}

// DeleteGroup : delete group
// This API is used to delete group
// @Summary Delete Group
// @Description This API is used to delete group
// @Tags Group
// @Accept  json
// @Produce  json
// @Param user body model.GroupUser true "jsonResponse"
// @Success 200 {object} jsonResponse
// @Failure 401 {object} jsonResponse
// @Router /user/delete-group [post]
func DeleteGroup(w http.ResponseWriter, r *http.Request) {
	var group model.GroupUser

	err := ReadJSON(w, r, &group)
	if err != nil {
		log.Println(err)
		ErrorJSON(w, err)
		return
	}
	userid, ok := r.Context().Value(UserIDKey).(uuid.UUID)
	if !ok {
		ErrorJSON(w, errors.New("can not get userid"))
		return
	}
	err = ser.Store.DeleteGroup(userid, group)
	if err != nil {
		log.Println(err)
		ErrorJSON(w, err)
		return
	}
	var payload jsonResponse
	payload.Error = false
	payload.Message = "group deleted"

	WriteJSON(w, http.StatusOK, payload)

}

// GetAllGroups : get all groups of user
// This API is used to get all groups which are related to user
// @Summary Get All Groups
// @Description This API is used to get all groups which are releted to user
// @Tags Group
// @Accept  json
// @Produce  json
// @Param user  true "model.UserGroups"
// @Success 200 {object} model.UserGroups
// @Failure 401 {object} jsonResponse
// @Router /user/getgroups [post]
func GetAllGroups(w http.ResponseWriter, r *http.Request) {
	userid, ok := r.Context().Value(UserIDKey).(uuid.UUID)
	if !ok {
		ErrorJSON(w, errors.New("can not get userid"))
		return
	}
	usrGroups, err := ser.Store.GetAllGroups(userid)
	if err != nil {
		log.Println(err)
		ErrorJSON(w, err)
		return
	}
	WriteJSON(w, http.StatusOK, usrGroups)
}

// CreateBuyOrder : Create a new buy order
// CreateBuyOrder : Create a new order with status pending and type buy
// @Summary Create a new buy order
// @Description Create a new order with status pending and type buy
// @Tags Order
// @Accept  json
// @Produce  json
// @Param order body model.OrderHandler true "model.OrderHandler"
// @Success 200 {object} model.OrderHandler
// @Failure 401 {object} jsonResponse
// @Router /user/create-buy-order [post]
func CreateBuyOrder(w http.ResponseWriter, r *http.Request) {
	var order model.OrderHandler

	err := ReadJSON(w, r, &order)
	if err != nil {
		log.Println(err)
		ErrorJSON(w, err)
		return
	}
	userid, ok := r.Context().Value(UserIDKey).(uuid.UUID)
	if !ok {
		ErrorJSON(w, errors.New("can not get userid"))
		return
	}
	/* or, err := util.ConvertOrderStringToOrder(&order)
	if err != nil {
		log.Println(err)
		ErrorJSON(w, err)
		return
	} */
	order, err = ser.Store.CreateNewOrder(userid, model.ORDER_TYPE_BUY, &order)
	if err != nil {
		log.Println(err)
		ErrorJSON(w, err)
		return
	}

	WriteJSON(w, http.StatusOK, order)
}

// CreateSellOrder : Create a new sell order
// CreateSellOrder : Create a new order with status pending and type sell
// @Summary Create a new sell order
// @Description Create a new order with status pending and type sell
// @Tags Order
// @Accept  json
// @Produce  json
// @Param order body model.OrderHandler true "model.OrderHandler"
// @Success 200 {object} model.OrderHandler
// @Failure 401 {object} jsonResponse
// @Router /user/create-sell-order [post]
func CreateSellOrder(w http.ResponseWriter, r *http.Request) {
	var order model.OrderHandler
	err := ReadJSON(w, r, &order)
	if err != nil {
		log.Println(err)
		ErrorJSON(w, err)
		return
	}
	userid, ok := r.Context().Value(UserIDKey).(uuid.UUID)
	if !ok {
		ErrorJSON(w, errors.New("can not get userid"))
		return
	}
	/* or, err := util.ConvertOrderStringToOrder(&order)
	if err != nil {
		log.Println(err)
		ErrorJSON(w, err)
		return
	} */
	order, err = ser.Store.CreateNewOrder(userid, model.ORDER_TYPE_SELL, &order)
	if err != nil {
		log.Println(err)
		ErrorJSON(w, err)
		return
	}

	WriteJSON(w, http.StatusOK, order)
}

// UpdateOrder godoc
// This Api is used to update the order
// @Summary Update Order Value
// @Description This Api is used to update the order
// @Tags Order
// @Accept  json
// @Produce  json
// @Param order body model.OrderHandler true "jsonResponse"
// @Success 200 {object} jsonResponse
// @Failure 401 {object} jsonResponse
// @Router /user/update-order [post]
func UpdateOrder(w http.ResponseWriter, r *http.Request) {
	var order model.OrderHandler

	err := ReadJSON(w, r, &order)
	if err != nil {
		log.Println(err)
		ErrorJSON(w, err)
		return
	}
	/* or, err := util.ConvertOrderStringToOrder(&order)
	if err != nil {
		log.Println(err)
		ErrorJSON(w, err)
		return
	} */
	err = ser.Store.UpdateOrder(&order)
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

// DeletedOrder godoc
// This API is used to delete an order
// @Summary Delete an order
// @Description Delete an order
// @Tags Order
// @Accept  json
// @Produce  json
// @Param order body model.OrderUser true "jsonResponse"
// @Success 200 {object} jsonResponse
// @Failure 401 {object} jsonResponse
// @Router /user/order-delete [post]
func DeleteOrder(w http.ResponseWriter, r *http.Request) {
	var order model.OrderUser

	err := ReadJSON(w, r, &order)
	if err != nil {
		log.Println(err)
		ErrorJSON(w, err)
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

}

// GetAllOrders : get all orders of user
// This API is used to get all orders which are related to user
// @Summary Get All Orders
// @Description This API is used to get all orders which are releted to user
// @Tags Order
// @Accept  json
// @Produce  json
// @Param user  true "[]model.Order"
// @Success 200 {object} []model.Order
// @Failure 401 {object} jsonResponse
// @Router /user/getorders [post]
func GetAllOrders(w http.ResponseWriter, r *http.Request) {
	userid, ok := r.Context().Value(UserIDKey).(uuid.UUID)
	if !ok {
		ErrorJSON(w, errors.New("can not get userid"))
		return
	}
	orders, err := ser.Store.GetAllOrders(userid)
	if err != nil {
		log.Println(err)
		ErrorJSON(w, err)
		return
	}
	WriteJSON(w, http.StatusOK, orders)
}

// CreateTrade godoc
// This Api is used to create new trade
// @Summary Create Trade
// @Description This Api is used create a new trade
// @Tags Trade
// @Accept  json
// @Produce  json
// @Param order body model.TradeHandlerUser true "jsonResponse"
// @Success 200 {object} jsonResponse
// @Failure 401 {object} jsonResponse
// @Router /user/create-trade [post]
func CreateTrade(w http.ResponseWriter, r *http.Request) {
	var trade model.TradeHandlerUser

	err := ReadJSON(w, r, &trade)
	if err != nil {
		log.Println(err)
		ErrorJSON(w, err)
		return
	}
	userid, ok := r.Context().Value(UserIDKey).(uuid.UUID)
	if !ok {
		ErrorJSON(w, errors.New("can not get userid"))
		return
	}
	/* tr, err := util.ConvertTradeStringToTrade(&trade)
	if err != nil {
		log.Println(err)
		ErrorJSON(w, err)
		return
	} */
	err = ser.Store.CreateTrade(userid, &trade)
	if err != nil {
		log.Println(err)
		ErrorJSON(w, err)
		return
	}
	var payload jsonResponse
	payload.Error = false
	payload.Message = "trade created"

	WriteJSON(w, http.StatusOK, payload)
}

// ConfirmingTrade godoc
// This Api is used to change the trade status to confirming
// @Summary ConfirmingTrade
// @Description This Apis is used to change the trade status to confirming
// @Tags Trade
// @Accept  json
// @Produce  json
// @Param trade body model.TradeUser true "jsonResponse"
// @Success 200 {object} jsonResponse
// @Failure 401 {object} jsonResponse
// @Router /user/trade-confirming [post]
func ConfirmingTrade(w http.ResponseWriter, r *http.Request) {
	var trade model.TradeUser

	err := ReadJSON(w, r, &trade)
	if err != nil {
		log.Println(err)
		ErrorJSON(w, err)
		return
	}

	err = ser.Store.ConfirmingTrade(&trade)
	if err != nil {
		log.Println(err)
		ErrorJSON(w, err)
		return
	}
	var payload jsonResponse
	payload.Error = false
	payload.Message = "trade status changed to confirming"

	WriteJSON(w, http.StatusOK, payload)

}

// ConfirmTrade godoc
// This Api is used to change the trade status to confirm
// @Summary ConfirmTrade
// @Description This Api is used to change the trade status to confirm
// @Tags Trade
// @Accept  json
// @Produce  json
// @Param trade body model.TradeUser true "jsonResponse"
// @Success 200 {object} jsonResponse
// @Failure 401 {object} jsonResponse
// @Router /user/trade-confirm [post]
func ConfirmTrade(w http.ResponseWriter, r *http.Request) {
	var trade model.TradeUser

	err := ReadJSON(w, r, &trade)
	if err != nil {
		log.Println(err)
		ErrorJSON(w, err)
		return
	}

	err = ser.Store.ConfirmTrade(&trade)
	if err != nil {
		log.Println(err)
		ErrorJSON(w, err)
		return
	}
	var payload jsonResponse
	payload.Error = false
	payload.Message = "trade status changed to confirm"

	WriteJSON(w, http.StatusOK, payload)

}

// DisputedTrade godoc
// This API is used to change the trade status to disputed
// @Summary DisputedTrade
// @Description This API is used to change the trade status to disputed
// @Tags Trade
// @Accept  json
// @Produce  json
// @Param trade body model.TradeUser  true "jsonResponse"
// @Success 200 {object} jsonResponse
// @Failure 401 {object} jsonResponse
// @Router /user/trade-disputed [post]
func DisputedTrade(w http.ResponseWriter, r *http.Request) {
	var trade model.TradeUser

	err := ReadJSON(w, r, &trade)
	if err != nil {
		log.Println(err)
		ErrorJSON(w, err)
		return
	}

	err = ser.Store.DisputedTrade(&trade)
	if err != nil {
		log.Println(err)
		ErrorJSON(w, err)
		return
	}
	var payload jsonResponse
	payload.Error = false
	payload.Message = "trade status changed to disputed"

	WriteJSON(w, http.StatusOK, payload)

}

// DeleteTrade godoc
// This API is used to delete an trade
// @Summary Delete an trade
// @Description Delete an trade
// @Tags Trade
// @Accept  json
// @Produce  json
// @Param trade body model.TradeUser true "jsonResponse"
// @Success 200 {object} jsonResponse
// @Failure 401 {object} jsonResponse
// @Router /user/trade-delete [post]
func DeleteTrade(w http.ResponseWriter, r *http.Request) {
	var trade model.TradeUser

	err := ReadJSON(w, r, &trade)
	if err != nil {
		log.Println(err)
		ErrorJSON(w, err)
		return
	}
	//Update to database
	err = ser.Store.DeleteTrade(&trade)
	if err != nil {
		log.Println(err)
		ErrorJSON(w, err)
		return
	}
	var payload jsonResponse
	payload.Error = false
	payload.Message = "trade Deleted"

	WriteJSON(w, http.StatusOK, payload)

	//WriteJSON(w, http.StatusAccepted, res)

}

// RenewAccessToken godoc
// This API is used to renew accesstoken using refreshtoken which you will get in verifyotp API.
// @Summary This API is used to renew accesstoken using refreshtoken which you will get in verifyotp API.
// @Description This API is used to renew accesstoken using refreshtoken which you will get in verifyotp API.
// @Tags User
// @Accept  json
// @Produce  json
// @Param user body renewAccessTokenRequests true "renewAccessTokenResponse"
// @Success 200 {object} jsonResponse
// @Failure 401 {object} jsonResponse
// @Router /renew-accesstoken [post]
func RenewAccessToken(w http.ResponseWriter, r *http.Request) {
	var req renewAccessTokenRequests
	err := ReadJSON(w, r, &req)
	if err != nil {
		log.Println(err)
		ErrorJSON(w, err)
		return
	}

	refreshPayload, err := ser.TokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		ErrorJSON(w, err, http.StatusUnauthorized)
		return
	}

	if time.Now().After(refreshPayload.ExpiredAt) {

		ErrorJSON(w, err, http.StatusUnauthorized)
		return
	}

	accessToken, accessPayload, err := ser.TokenMaker.CreateToken(
		refreshPayload.UserId,
		ser.Config.AccessTokenDuration,
	)
	if err != nil {
		ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	rsp := renewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
	}

	WriteJSON(w, http.StatusOK, rsp)
}
func GetAllUsersTrade(w http.ResponseWriter, r *http.Request) {
	userid, ok := r.Context().Value(UserIDKey).(uuid.UUID)
	if !ok {
		ErrorJSON(w, errors.New("can not get userid"))
		return
	}
	trades, err := ser.Store.GetAllUserTrade(userid)
	if err != nil {
		log.Println(err)
		ErrorJSON(w, err)
		return
	}
	WriteJSON(w, http.StatusOK, trades)
}

// GetOrderTrades : get all trades of order
// This API is used to get all trades which are related to order
// @Summary Get Order Trades
// @Description This API is used to get all trades which are releted to order
// @Tags Order
// @Accept  json
// @Produce  json
// @Param user  true "[]model.TradeHandler"
// @Success 200 {object} []model.TradeHandler
// @Failure 401 {object} jsonResponse
// @Router /user/getordertrades [post]
func GetOrderTrades(w http.ResponseWriter, r *http.Request) {
	var order model.OrderUser
	err := ReadJSON(w, r, &order)
	if err != nil {
		log.Println(err)
		ErrorJSON(w, err)
		return
	}
	userid, ok := r.Context().Value(UserIDKey).(uuid.UUID)
	if !ok {
		ErrorJSON(w, errors.New("can not get userid"))
		return
	}
	trades, err := ser.Store.GetOrderTrades(userid, order)
	if err != nil {
		log.Println(err)
		ErrorJSON(w, err)
		return
	}
	WriteJSON(w, http.StatusOK, trades)
}
