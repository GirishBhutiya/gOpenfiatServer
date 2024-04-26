package db

import (
	"errors"

	"github.com/datastax-ext/astra-go-sdk"
	"github.com/google/uuid"
	"github.com/swisscdn/OpenfiatServer/types"
)

// UserService represents a PostgreSQL implementation of myapp.UserService.
type DB struct {
	DB *astra.Client
}

type DatabaseService interface {
	InserLoginRegister(otp int, usr *types.User) error
	VerifyOTP(phoneNumber int, otp int) (types.User, error)
	UpdateUser(usr *types.User) (types.User, error)
	DeleteUser(usr *types.User) error
	CreateNewOrder(order *types.Order) error
	UpdateOrderValue(order *types.Order) error
	ConfirmingOrder(order *types.Order) error
	ConfirmOrder(order *types.Order) error
	DisputedOrder(order *types.Order) error
	GetAllOrders(usr *types.User) ([]types.Order, error)
	DeleteOrder(order *types.Order) error
}

func (db *DB) InserLoginRegister(otp int, usr *types.User) error {
	//var usr types.User
	stmt := `INSERT INTO users (phonenumber, verified, roll_id,otp) values (?, ?, ?,?) IF NOT EXISTS;`
	rows, err := db.DB.Query(stmt, usr.PhoneNumber, false, 2, otp).Exec()
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

func (db *DB) VerifyOTP(phoneNumber int, otp int) (types.User, error) {
	var user types.User
	var motp int
	sel := `SELECT phonenumber,verified,first_name,last_name,roll_id,profile_pic,otp FROM users WHERE phonenumber=?;`
	rows, err := db.DB.Query(sel, phoneNumber).Exec()
	if err != nil {
		return user, err
	}
	if len(rows) == 0 {
		return user, errors.New("user not found")
	}
	for _, r := range rows {
		//vals := r.Values()
		r.Scan(&user.PhoneNumber, &user.Verified, &user.FirstName, &user.LastName, &user.RollId, &user.ProfilePic, &motp)
		/* fmt.Println("\nValues:", vals)
		fmt.Println("\nuser:", user)
		fmt.Println("opt:", otp, " motp:", motp) */
	}
	if motp != otp {
		return user, errors.New("invalid otp")
	}

	stmt := `UPDATE users SET verified = true WHERE phonenumber=? IF EXISTS;;`
	rows, err = db.DB.Query(stmt, phoneNumber).Exec()
	if err != nil {
		return user, err
	}
	if rows[0].Values()[0] == false {
		return user, errors.New("user not found")
	}
	/* for _, r := range rows {
		vals := r.Values()
		fmt.Println("Values:", vals)
	} */
	return user, nil
}

func (db *DB) UpdateUser(usr *types.User) (types.User, error) {

	stmt := `UPDATE users SET first_name=?, last_name=?, profile_pic=? WHERE phonenumber=?  IF EXISTS;`
	rows, err := db.DB.Query(stmt, usr.FirstName, usr.LastName, usr.ProfilePic, usr.PhoneNumber).Exec()
	if err != nil {
		return *usr, err
	}
	if rows[0].Values()[0] == false {
		return *usr, errors.New("user not found")
	}

	return *usr, nil
}
func (db *DB) DeleteUser(usr *types.User) error {

	stmt := `DELETE FROM users WHERE phonenumber=?  IF EXISTS;`
	rows, err := db.DB.Query(stmt, usr.PhoneNumber).Exec()
	if err != nil {
		return err
	}
	if rows[0].Values()[0] == false {
		return errors.New("user not found")
	}

	return nil
}
func (db *DB) CreateNewOrder(order *types.Order) error {
	stmt := `INSERT INTO orders (id, from_phonenumber, to_phonenumber, amount, status) values (?, ?, ?, ?, ?) IF NOT EXISTS;`
	rows, err := db.DB.Query(stmt, uuid.New(), order.FromPhone, order.ToPhone, order.Amount, types.ORDER_STATUS_PENDING).Exec()
	if err != nil {
		//log.Println("Error:", err)
		return err
	}
	if rows[0].Values()[0] == false {
		return errors.New("order not placed")
	}

	return nil
}
func (db *DB) UpdateOrderValue(order *types.Order) error {
	stmt := `UPDATE orders SET amount =? WHERE id=?  IF EXISTS;`
	rows, err := db.DB.Query(stmt, order.Amount, order.ID).Exec()
	if err != nil {
		//log.Println("Error:", err)
		return err
	}
	if rows[0].Values()[0] == false {
		return errors.New("order not found")
	}

	return nil
}
func (db *DB) ConfirmingOrder(order *types.Order) error {
	stmt := `UPDATE orders SET status =? WHERE id=?  IF EXISTS;`
	rows, err := db.DB.Query(stmt, types.ORDER_STATUS_CONFIRMING, order.ID).Exec()
	if err != nil {
		//log.Println("Error:", err)
		return err
	}
	if rows[0].Values()[0] == false {
		return errors.New("order not found")
	}

	return nil
}
func (db *DB) ConfirmOrder(order *types.Order) error {
	stmt := `UPDATE orders SET status =? WHERE id=?  IF EXISTS;`
	rows, err := db.DB.Query(stmt, types.ORDER_STATUS_CONFIRM, order.ID).Exec()
	if err != nil {
		//log.Println("Error:", err)
		return err
	}
	if rows[0].Values()[0] == false {
		return errors.New("order not found")
	}

	return nil
}
func (db *DB) DisputedOrder(order *types.Order) error {
	stmt := `UPDATE orders SET status =? WHERE id=?  IF EXISTS;`
	rows, err := db.DB.Query(stmt, types.ORDER_STATUS_DISPUTED, order.ID).Exec()
	if err != nil {
		//log.Println("Error:", err)
		return err
	}
	if rows[0].Values()[0] == false {
		return errors.New("order not found")
	}

	return nil
}
func (db *DB) GetAllOrders(usr *types.User) ([]types.Order, error) {
	var orders []types.Order
	sel := `SELECT id,from_phonenumber,to_phonenumber,amount,status FROM orders WHERE from_phonenumber=? OR to_phonenumber=? ALLOW FILTERING;`
	rows, err := db.DB.Query(sel, usr.PhoneNumber, usr.PhoneNumber).Exec()
	if err != nil {
		return orders, err
	}
	for _, r := range rows {
		var order types.Order
		//vals := r.Values()
		r.Scan(&order.ID, &order.FromPhone, &order.ToPhone, &order.Amount, &order.Status)
		//fmt.Println("\nValues:", vals)

		orders = append(orders, order)
	}

	return orders, nil
}
func (db *DB) DeleteOrder(order *types.Order) error {

	stmt := `DELETE FROM orders WHERE id=?  IF EXISTS;`
	rows, err := db.DB.Query(stmt, order.ID).Exec()
	if err != nil {
		return err
	}
	if rows[0].Values()[0] == false {
		return errors.New("order not found")
	}

	return nil
}
