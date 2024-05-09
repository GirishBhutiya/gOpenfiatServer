package db

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/GirishBhutiya/gOpenfiatServer/app/model"
	"github.com/GirishBhutiya/gOpenfiatServer/app/util"
	"github.com/datastax-ext/astra-go-sdk"
	"github.com/google/uuid"
)

// UserService represents a PostgreSQL implementation of myapp.UserService.
type DB struct {
	DB *astra.Client
}

type DatabaseService interface {
	//user
	InserLoginRegister(otp int, usr *model.UserLogin) error
	Login(phoneNumber int, otp int) (model.User, error)
	UpdateUser(profilepath string, usr *model.UserUpdate) (model.User, error)
	SubscribeGroupToUSer(userId uuid.UUID, group model.GroupUser) error
	UnsubscribeGroupToUSer(userId uuid.UUID, group model.GroupUser) error
	GetAllGroups(userId uuid.UUID) ([]model.GroupUser, error)
	CreateShortLink(group *model.GroupUser) (string, error)
	SubscribeGroupViaInvite(inviteKey string, userId uuid.UUID) error
	//groups
	CreateNewGroup(userId uuid.UUID, group model.CreateGroup) (model.GroupUser, error)
	UpdateGroup(group model.GroupUser) error
	DeleteGroup(userId uuid.UUID, group model.GroupUser) error

	//order
	CreateNewOrder(userId uuid.UUID, orderType string, order *model.OrderHandler) (model.OrderHandler, error)
	UpdateOrder(order *model.OrderHandler) error
	DeleteOrder(order *model.OrderUser) error
	GetAllOrders(userId uuid.UUID) ([]model.Order, error)

	//trade
	CreateTrade(userId uuid.UUID, trade *model.TradeHandlerUser) error
	ConfirmingTrade(trade *model.TradeUser) error
	ConfirmTrade(trade *model.TradeUser) error
	DisputedTrade(trade *model.TradeUser) error
	DeleteTrade(trade *model.TradeUser) error
	GetAllUserTrade(userId uuid.UUID) ([]model.UserTrades, error)
	GetOrderTrades(userId uuid.UUID, order model.OrderUser) ([]model.TradeHandler, error)
}

func (db *DB) InserLoginRegister(otp int, usr *model.UserLogin) error {
	//var usr model.User
	stmt := `INSERT INTO users (userid,phonenumber, verified,otp) values (?, ?, ?, ?) IF NOT EXISTS;`
	rows, err := db.DB.Query(stmt, uuid.New(), usr.PhoneNumber, false, otp).Exec()
	if err != nil {
		//log.Println("Error:", err)
		return err
	}
	if rows[0].Values()[0] == false {
		//fmt.Println("Girish:", rows[0].Values()[0])
		update := `UPDATE users SET otp=? WHERE phonenumber=? IF EXISTS;`
		row, err := db.DB.Query(update, otp, usr.PhoneNumber).Exec()
		if err != nil {
			//log.Println("in update Error:", err)
			return err
		}
		if row[0].Values()[0] == true {
			return nil
		}

	}
	/* for _, r := range rows {
		vals := r.Values()
		fmt.Println("Values:", vals)
	} */
	return nil
}

func (db *DB) Login(phoneNumber int, otp int) (model.User, error) {
	var user model.User
	var motp, phone int
	sel := `SELECT userid,phonenumber,verified,first_name,last_name,profile_pic,otp FROM users WHERE phonenumber=? ALLOW FILTERING;`
	rows, err := db.DB.Query(sel, phoneNumber).Exec()
	if err != nil {
		return user, err
	}
	if len(rows) == 0 {
		return user, errors.New("user not found")
	}
	for _, r := range rows {
		//vals := r.Values()
		r.Scan(&user.ID, &phone, &user.Verified, &user.FirstName, &user.LastName, &user.ProfilePic, &motp)
		/* fmt.Println("\nValues:", vals)
		fmt.Println("\nuser:", user)
		fmt.Println("opt:", otp, " motp:", motp) */
	}
	//log.Println(user, "\nopt", otp)
	if motp != otp {
		return user, errors.New("invalid otp")
	}

	stmt := `UPDATE users SET verified = true WHERE userid=? IF EXISTS;;`
	_, err = db.DB.Query(stmt, user.ID).Exec()
	if err != nil {
		return user, err
	}
	/* for _, r := range rows {
		vals := r.Values()
		fmt.Println("Values:", vals)
	} */
	user.Verified = true
	return user, nil
}

func (db *DB) UpdateUser(profilepath string, usr *model.UserUpdate) (model.User, error) {
	var user model.User
	stmt := `UPDATE users SET first_name=?, last_name=?, profile_pic=?, phonenumber=? WHERE userid=?  IF EXISTS;`
	rows, err := db.DB.Query(stmt, usr.FirstName, usr.LastName, profilepath, usr.PhoneNumber, usr.ID).Exec()
	if err != nil {
		return user, err
	}
	if rows[0].Values()[0] == false {
		return user, errors.New("user not found")
	}
	user.ID = usr.ID
	user.FirstName = usr.FirstName
	user.LastName = usr.LastName
	user.PhoneNumber = usr.PhoneNumber
	user.Verified = usr.Verified

	return user, nil
}

func (db *DB) SubscribeGroupToUSer(userId uuid.UUID, group model.GroupUser) error {
	/* stmt := fmt.Sprintf("UPDATE users SET groups = groups + { %s : '%s' } WHERE userid=? IF EXISTS;", group.ID, group.Name)
	//stmt := `UPDATE users SET groups = groups + { ? : '?' } WHERE userid=? IF EXISTS;`
	rows, err := db.DB.Query(stmt, userId).Exec()
	if err != nil {
		return err
	}
	if rows[0].Values()[0] == false {
		return errors.New("user not found")
	}

	return nil */

	stmt := `INSERT INTO groupuserrelation (groupid,userid) values (?, ?) IF NOT EXISTS;`
	rows, err := db.DB.Query(stmt, group.ID, userId).Exec()
	if err != nil {
		log.Println("Error:", err)
		return err
	}
	if rows[0].Values()[0] == false {
		return errors.New("user already subscribed")
	}
	return nil

}
func (db *DB) UnsubscribeGroupToUSer(userId uuid.UUID, group model.GroupUser) error {
	/* stmt := fmt.Sprintf("UPDATE users SET groups = groups - { %s } WHERE userid=? IF EXISTS;", group.ID)
	//stmt := `UPDATE users SET groups = groups - { ? } WHERE userid=? IF EXISTS;`
	rows, err := db.DB.Query(stmt, userId).Exec()
	log.Println(rows)
	if err != nil {
		return err
	}
	if rows[0].Values()[0] == false {
		return errors.New("user not found")
	}

	return nil */
	stmt := `DELETE FROM groupuserrelation WHERE groupid=? AND userid=? IF EXISTS;`
	rows, err := db.DB.Query(stmt, group.ID, userId).Exec()
	if err != nil {
		//log.Println("Error:", err)
		return err
	}
	if rows[0].Values()[0] == false {
		return errors.New("data not found")
	}
	return nil

}
func (db *DB) GetAllGroups(userId uuid.UUID) ([]model.GroupUser, error) {
	var usrGroups []model.GroupUser
	stmt := "SELECT groupid FROM groupuserrelation WHERE userid=? ALLOW FILTERING"
	rows, err := db.DB.Query(stmt, userId).Exec()
	if err != nil {
		return usrGroups, err
	}
	if len(rows) == 0 {
		return usrGroups, errors.New("no groups found")
	}
	var gids []uuid.UUID
	for _, r := range rows {
		var gid uuid.UUID
		//vals := r.Values()
		r.Scan(&gid)
		//fmt.Println("\n gIDS:", vals)
		gids = append(gids, gid)

	}
	var sb strings.Builder
	sb.WriteString("?")
	for range gids[1:] {
		sb.WriteString(", ?")
	}

	stmt = fmt.Sprintf("SELECT groupid,groupname FROM group WHERE groupid IN ( %s) ALLOW FILTERING", sb.String())
	/* stmt = "SELECT groupid,groupname FROM group WHERE groupid IN ( ?,? ) ALLOW FILTERING"
	log.Println("in else")
	var sb strings.Builder
	sb.WriteString(gids[0].String())
	for _, g := range gids[1:] {
		sb.WriteString(", ")
		sb.WriteString(g.String())
	}
	log.Println("in else", sb.String()) */
	rows, err = db.DB.Query(stmt, util.ConvertUUIDToAny(gids)...).Exec()

	if err != nil {
		return usrGroups, err
	}
	if len(rows) == 0 {
		return usrGroups, errors.New("no groups found")
	}
	for _, r := range rows {
		var g model.GroupUser
		//vals := r.Values()
		r.Scan(&g.ID, &g.Name)
		//fmt.Println("\n Groups:", vals)
		usrGroups = append(usrGroups, g)
	}
	return usrGroups, err
}

func (db *DB) CreateNewGroup(userId uuid.UUID, group model.CreateGroup) (model.GroupUser, error) {
	stmt := `INSERT INTO group (groupid, groupname, creatorUserid, createtime) values (?, ?, ?, ?) IF NOT EXISTS;`
	id := uuid.New()
	var g model.GroupUser
	rows, err := db.DB.Query(stmt, id, group.Name, userId, time.Now()).Exec()
	if err != nil {
		//log.Println("Error:", err)
		return g, err
	}
	log.Println(rows)
	if rows[0].Values()[0] == false {
		return g, errors.New("group not created")
	}
	g.ID = id
	g.Name = group.Name

	err = db.SubscribeGroupToUSer(userId, g)
	if err != nil {
		//log.Println("Error:", err)
		return g, errors.New("group not subscribed to user")
	}
	return g, nil
}
func (db *DB) UpdateGroup(group model.GroupUser) error {

	stmt := `UPDATE group SET groupname=? WHERE groupid=?  IF EXISTS;`
	rows, err := db.DB.Query(stmt, group.Name, group.ID).Exec()
	if err != nil {
		//log.Println("Error:", err)
		return err
	}
	if rows[0].Values()[0] == false {
		return errors.New("group not created")
	}

	return nil
}
func (db *DB) DeleteGroup(userId uuid.UUID, group model.GroupUser) error {

	stmt := `DELETE FROM group WHERE groupid=? IF EXISTS;`
	rows, err := db.DB.Query(stmt, group.ID).Exec()
	if err != nil {
		return err
	}
	if rows[0].Values()[0] == false {
		return errors.New("group not found")
	}
	err = db.UnsubscribeGroupToUSer(userId, group)
	if err != nil {
		//log.Println("Error:", err)
		return errors.New("group not unsubscribed to user")
	}
	return nil
}

func (db *DB) CreateNewOrder(userId uuid.UUID, orderType string, order *model.OrderHandler) (model.OrderHandler, error) {
	var orderN model.OrderHandler
	id := uuid.New()
	stmt := `INSERT INTO orders (orderid, userid, fiatAmount, minAmount, price, timeLimit,type) values (?, ?, ?, ?, ?, ?, ?) IF NOT EXISTS;`
	rows, err := db.DB.Query(stmt, id, userId, order.FiatAmount, order.MinAmount, order.Price, order.TimeLimit, orderType).Exec()
	if err != nil {
		//log.Println("Error:", err)
		return orderN, err
	}
	if rows[0].Values()[0] == false {
		return orderN, errors.New("order not placed")
	}
	order.ID = id
	order.UserId = userId
	order.Type = orderType
	return *order, nil
}
func (db *DB) UpdateOrder(order *model.OrderHandler) error {
	stmt := `UPDATE orders SET fiatAmount=?, minAmount =?, price=?, timeLimit=? WHERE orderid=?  IF EXISTS;`
	rows, err := db.DB.Query(stmt, order.FiatAmount, order.MinAmount, order.Price, order.TimeLimit, order.ID).Exec()
	if err != nil {
		//log.Println("Error:", err)
		return err
	}
	if rows[0].Values()[0] == false {
		return errors.New("order not found")
	}

	return nil
}
func (db *DB) DeleteOrder(order *model.OrderUser) error {
	stmt := `DELETE FROM orders WHERE orderid=?  IF EXISTS;`
	rows, err := db.DB.Query(stmt, order.ID).Exec()
	if err != nil {
		return err
	}
	if rows[0].Values()[0] == false {
		return errors.New("order not found")
	}

	return nil
}
func (db *DB) GetAllOrders(userId uuid.UUID) ([]model.Order, error) {
	var orders []model.Order
	stmt := `SELECT orderid,fiatamount,minamount,price,timelimit,type FROM orders WHERE userid=? ALLOW FILTERING;`
	rows, err := db.DB.Query(stmt, userId).Exec()
	if err != nil {
		return orders, err
	}
	if len(rows) == 0 {
		return orders, errors.New("no urders found")
	}
	for _, r := range rows {
		var order model.Order
		//vals := r.Values()
		r.Scan(&order.ID, &order.FiatAmount, &order.MinAmount, &order.Price, &order.TimeLimit, &order.Type)

		//fmt.Println("\nValues:", vals)
		//fmt.Println("\norder:", order)
		orders = append(orders, order)

	}
	//fmt.Println("\norder:", orders)
	return orders, nil

}
func (db *DB) CreateTrade(userId uuid.UUID, trade *model.TradeHandlerUser) error {
	stmt := `INSERT INTO trade (tradeid, orderid, bidUserid, tradetime, status, method) values (?, ?, ?, ?, ?, ?) IF NOT EXISTS;`
	rows, err := db.DB.Query(stmt, uuid.New(), trade.Orderid, userId, trade.TradeTime, model.ORDER_STATUS_PENDING, trade.Method).Exec()
	if err != nil {
		//log.Println("Error:", err)
		return err
	}
	if rows[0].Values()[0] == false {
		return errors.New("trade not created")
	}

	return nil
}

func (db *DB) ConfirmingTrade(trade *model.TradeUser) error {
	stmt := `UPDATE trade SET status =? WHERE tradeid=?  IF EXISTS;`
	rows, err := db.DB.Query(stmt, model.ORDER_STATUS_CONFIRMING, trade.ID).Exec()
	if err != nil {
		//log.Println("Error:", err)
		return err
	}
	if rows[0].Values()[0] == false {
		return errors.New("trade not found")
	}

	return nil
}
func (db *DB) ConfirmTrade(trade *model.TradeUser) error {
	stmt := `UPDATE trade SET status =? WHERE tradeid=?  IF EXISTS;`
	rows, err := db.DB.Query(stmt, model.ORDER_STATUS_CONFIRM, trade.ID).Exec()
	if err != nil {
		//log.Println("Error:", err)
		return err
	}
	if rows[0].Values()[0] == false {
		return errors.New("trade not found")
	}

	return nil
}
func (db *DB) DisputedTrade(trade *model.TradeUser) error {
	stmt := `UPDATE trade SET status =? WHERE tradeid=?  IF EXISTS;`
	rows, err := db.DB.Query(stmt, model.ORDER_STATUS_DISPUTED, trade.ID).Exec()
	if err != nil {
		//log.Println("Error:", err)
		return err
	}
	if rows[0].Values()[0] == false {
		return errors.New("trade not found")
	}

	return nil
}
func (db *DB) DeleteTrade(trade *model.TradeUser) error {
	stmt := `DELETE FROM trade WHERE tradeid=?  IF EXISTS;`
	rows, err := db.DB.Query(stmt, trade.ID).Exec()
	if err != nil {
		return err
	}
	if rows[0].Values()[0] == false {
		return errors.New("trade not found")
	}

	return nil
}

func (db *DB) GetAllUserTrade(userId uuid.UUID) ([]model.UserTrades, error) {
	var trades []model.UserTrades
	stmt := `SELECT tradeid,method,orderid,status,tradetime FROM trade WHERE biduserid=? ALLOW FILTERING;`
	rows, err := db.DB.Query(stmt, userId).Exec()
	if err != nil {
		return trades, err
	}
	if len(rows) == 0 {
		return trades, errors.New("no trade found")
	}
	for _, r := range rows {
		var trade model.UserTrades
		//vals := r.Values()
		r.Scan(&trade.ID, &trade.Method, &trade.Orderid, &trade.Status, &trade.TradeTime)

		//fmt.Println("\nValues:", vals)
		//fmt.Println("\norder:", order)
		trades = append(trades, trade)

	}
	//fmt.Println("\norder:", trades)
	return trades, nil

}
func (db *DB) GetOrderTrades(userId uuid.UUID, order model.OrderUser) ([]model.TradeHandler, error) {
	var trades []model.TradeHandler
	stmt := `SELECT tradeid,biduserid,method,orderid,status,tradetime FROM trade WHERE orderid=? ALLOW FILTERING;`
	rows, err := db.DB.Query(stmt, order.ID).Exec()
	if err != nil {
		return trades, err
	}
	if len(rows) == 0 {
		return trades, errors.New("no trade found")
	}
	for _, r := range rows {
		var trade model.TradeHandler
		//vals := r.Values()
		r.Scan(&trade.ID, &trade.BidUserid, &trade.Method, &trade.Orderid, &trade.Status, &trade.TradeTime)

		//fmt.Println("\nValues:", vals)
		//fmt.Println("\norder:", trade)
		trades = append(trades, trade)

	}
	//fmt.Println("\norder:", trades)
	return trades, nil

}
func (db *DB) CreateShortLink(group *model.GroupUser) (string, error) {
	inviteKey := util.GenerateShortKey()
	stmt := `INSERT INTO groupinvite (groupid,invitekey) values (?, ?) IF NOT EXISTS;`
	rows, err := db.DB.Query(stmt, group.ID, inviteKey).Exec()
	if err != nil {
		//log.Println("Error:", err)
		return "", err
	}
	if rows[0].Values()[0] == false {
		//fmt.Println("Girish:", rows[0].Values()[0])
		sect := `SELECT invitekey FROM groupinvite WHERE groupid=? ALLOW FILTERING;`
		rows, err := db.DB.Query(sect, group.ID).Exec()
		if err != nil {
			return "", err
		}

		for _, r := range rows {
			//vals := r.Values()
			r.Scan(&inviteKey)

			//fmt.Println("\nValues:", vals)
			//fmt.Println("\norder:", trade)

		}

	}

	return inviteKey, nil
}
func (db *DB) SubscribeGroupViaInvite(inviteKey string, userId uuid.UUID) error {
	var group model.GroupUser
	sect := `SELECT groupid FROM groupinvite WHERE invitekey=? ALLOW FILTERING;`
	rows, err := db.DB.Query(sect, inviteKey).Exec()
	if err != nil {
		return err
	}

	for _, r := range rows {
		//vals := r.Values()
		r.Scan(&group.ID)
		//fmt.Println("\nValues:", vals)
		//fmt.Println("\norder:", trade)
	}
	log.Println("groupid", group.ID, " userid", userId)
	err = db.SubscribeGroupToUSer(userId, group)
	if err != nil {
		log.Println("Error:", err)
		return errors.New("group not subscribed to user")
	}
	return nil
}
