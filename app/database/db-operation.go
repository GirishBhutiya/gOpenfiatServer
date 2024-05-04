package db

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/GirishBhutiya/gOpenfiatServer/app/model"
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
	UpdateUser(usr *model.User) (model.User, error)
	SubscribeGroupToUSer(userId uuid.UUID, group model.GroupUser) error
	UnsubscribeGroupToUSer(userId uuid.UUID, group model.GroupUser) error

	//groups
	CreateNewGroup(userId uuid.UUID, group model.CreateGroup) error
	UpdateGroup(group model.GroupUser) error
	DeleteGroup(group model.GroupUser) error

	//order
	CreateNewOrder(userId uuid.UUID, orderType string, order *model.OrderHandler) error
	UpdateOrder(order *model.OrderHandler) error
	DeleteOrder(order *model.OrderUser) error

	//trade
	CreateTrade(userId uuid.UUID, trade *model.TradeHandler) error
	ConfirmingTrade(trade *model.TradeUser) error
	ConfirmTrade(trade *model.TradeUser) error
	DisputedTrade(trade *model.TradeUser) error
	DeleteTrade(trade *model.TradeUser) error
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
	sel := `SELECT userid,phonenumber,verified,first_name,last_name,groups,profile_pic,otp FROM users WHERE phonenumber=? ALLOW FILTERING;`
	rows, err := db.DB.Query(sel, phoneNumber).Exec()
	if err != nil {
		return user, err
	}
	if len(rows) == 0 {
		return user, errors.New("user not found")
	}
	for _, r := range rows {
		//vals := r.Values()
		r.Scan(&user.ID, &phone, &user.Verified, &user.FirstName, &user.LastName, &user.Groups, &user.ProfilePic, &motp)
		/* fmt.Println("\nValues:", vals)
		fmt.Println("\nuser:", user)
		fmt.Println("opt:", otp, " motp:", motp) */
	}
	log.Println(user, "\nopt", otp)
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
	return user, nil
}

func (db *DB) UpdateUser(usr *model.User) (model.User, error) {

	stmt := `UPDATE users SET first_name=?, last_name=?, profile_pic=?, phonenumber=? WHERE userid=?  IF EXISTS;`
	rows, err := db.DB.Query(stmt, usr.FirstName, usr.LastName, usr.ProfilePic, usr.PhoneNumber, usr.ID).Exec()
	if err != nil {
		return *usr, err
	}
	if rows[0].Values()[0] == false {
		return *usr, errors.New("user not found")
	}

	return *usr, nil
}

func (db *DB) SubscribeGroupToUSer(userId uuid.UUID, group model.GroupUser) error {
	stmt := fmt.Sprintf("UPDATE users SET groups = groups + { %s : '%s' } WHERE userid=? IF EXISTS;", group.ID, group.Name)
	//stmt := `UPDATE users SET groups = groups + { ? : '?' } WHERE userid=? IF EXISTS;`
	rows, err := db.DB.Query(stmt, userId).Exec()
	if err != nil {
		return err
	}
	if rows[0].Values()[0] == false {
		return errors.New("user not found")
	}

	return nil
}
func (db *DB) UnsubscribeGroupToUSer(userId uuid.UUID, group model.GroupUser) error {
	stmt := fmt.Sprintf("UPDATE users SET groups = groups - { %s } WHERE userid=? IF EXISTS;", group.ID)
	//stmt := `UPDATE users SET groups = groups - { ? } WHERE userid=? IF EXISTS;`
	rows, err := db.DB.Query(stmt, userId).Exec()
	log.Println(rows)
	if err != nil {
		return err
	}
	if rows[0].Values()[0] == false {
		return errors.New("user not found")
	}

	return nil
}

func (db *DB) CreateNewGroup(userId uuid.UUID, group model.CreateGroup) error {
	stmt := `INSERT INTO group (groupid, groupname, creatorUserid, createtime) values (?, ?, ?, ?) IF NOT EXISTS;`
	rows, err := db.DB.Query(stmt, uuid.New(), group.Name, userId, time.Now()).Exec()
	if err != nil {
		//log.Println("Error:", err)
		return err
	}
	if rows[0].Values()[0] == false {
		return errors.New("group not created")
	}

	return nil
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
func (db *DB) DeleteGroup(group model.GroupUser) error {

	stmt := `DELETE FROM group WHERE groupid=?  IF EXISTS;`
	rows, err := db.DB.Query(stmt, group.ID).Exec()
	if err != nil {
		return err
	}
	if rows[0].Values()[0] == false {
		return errors.New("group not found")
	}

	return nil
}

func (db *DB) CreateNewOrder(userId uuid.UUID, orderType string, order *model.OrderHandler) error {
	stmt := `INSERT INTO orders (orderid, userid, fiatAmount, minAmount, price, timeLimit,type) values (?, ?, ?, ?, ?, ?, ?) IF NOT EXISTS;`
	rows, err := db.DB.Query(stmt, uuid.New(), userId, order.FiatAmount, order.MinAmount, order.Price, order.TimeLimit, orderType).Exec()
	if err != nil {
		//log.Println("Error:", err)
		return err
	}
	if rows[0].Values()[0] == false {
		return errors.New("order not placed")
	}

	return nil
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

func (db *DB) CreateTrade(userId uuid.UUID, trade *model.TradeHandler) error {
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
