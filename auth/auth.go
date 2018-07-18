package auth

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type Auth struct {
	db              *sqlx.DB
	loginDelay      time.Duration
	listQuery       *sqlx.Stmt
	getQuery        *sqlx.Stmt
	getByLoginQuery *sqlx.Stmt
	createQuery     *sqlx.Stmt
	deleteQuery     *sqlx.Stmt
	updateQuery     *sqlx.Stmt
}

type UserAuth struct {
	UserId   int    `json:"user_id" db:"id"`
	UserRole string `json:"user_role" db:"-"`
}

type User struct {
	Id           int    `db:"id"`
	Login        string `db:"login"`
	PasswordHash string `db:"password_hash"`
}

func NewAuth(driver, dbConnect string, loginDelay time.Duration) (*Auth, error) {
	var err error

	d := &Auth{
		db:         sqlx.MustConnect(driver, dbConnect),
		loginDelay: loginDelay,
	}

	d.listQuery, err = d.db.Preparex("SELECT id, login FROM auth")
	if err != nil {
		return nil, err
	}

	d.getQuery, err = d.db.Preparex("SELECT id FROM auth WHERE login = ? AND password_hash = ?")
	if err != nil {
		return nil, err
	}

	d.getByLoginQuery, err = d.db.Preparex("SELECT id, password_hash FROM auth WHERE login = ?")
	if err != nil {
		return nil, err
	}

	d.createQuery, err = d.db.Preparex("INSERT INTO auth (login, password_hash) VALUES (?, ?)")
	if err != nil {
		return nil, err
	}

	d.updateQuery, err = d.db.Preparex("UPDATE auth SET password_hash = ? WHERE id = ?")
	if err != nil {
		return nil, err
	}

	d.deleteQuery, err = d.db.Preparex("DELETE FROM auth WHERE id = ?")
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (d *Auth) Authenticator(username string, password string, c *gin.Context) (interface{}, bool) {
	time.Sleep(d.loginDelay)
	return d.GetUser(username, password)
}

func (d *Auth) ListUsers() ([]User, error) {
	var users = make([]User, 0, 1024)
	rows, err := d.listQuery.Query()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var id int
		var login string
		err = rows.Scan(&id, &login)
		if err != nil {
			return nil, err
		}
		users = append(users, User{Id: id, Login: login})
	}

	return users, nil
}

func (d *Auth) GetUser(login string, passwordHash string) (*UserAuth, bool) {
	var user UserAuth
	err := d.getQuery.Get(&user, login, passwordHash)
	if err != nil {
		return nil, false
	}
	return &user, true
}

func (d *Auth) GetUserByLogin(login string) (*User, error) {
	var user User
	err := d.getByLoginQuery.Get(&user, login)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (d *Auth) CreateUser(login string, passwordHash string) error {
	res, err := d.createQuery.Exec(login, passwordHash)
	if err != nil {
		return err
	}
	num, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if num == 0 {
		errors.New("now rows affected")
	}
	return nil
}

func (d *Auth) UpdateUser(id int, passwordHash string) error {
	res, err := d.updateQuery.Exec(passwordHash, id)
	if err != nil {
		return err
	}
	num, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if num == 0 {
		errors.New("now rows affected")
	}
	return nil
}

func (d *Auth) DeleteUser(id int) error {
	_, err := d.deleteQuery.Exec(id)
	if err != nil {
		return err
	}
	return nil
}
