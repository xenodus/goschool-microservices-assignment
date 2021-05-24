package main

import (
	"database/sql"
	"errors"
	"strings"
	"unicode/utf8"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Id       string
	Email    string
	Password string
	ApiKeyId int
	Admin    int
}

type UserAuth struct {
	Email    string
	Password string
}

func (u *User) register() error {
	results, err := myDb.Exec("INSERT INTO user VALUES (?,?,?,?,?)", u.Id, u.Email, u.Password, 0, 0)

	if err != nil {
		return err
	} else {
		rows, _ := results.RowsAffected()

		if rows > 0 {
			u.setKey()
			return nil
		}
	}

	return errors.New("unable to create user")
}

func (u *User) setKey() error {
	key, err := generateKey()

	if err != nil {
		return err
	} else {
		newKey, newKeyErr := getKeyByValue(key)

		if newKeyErr != nil {
			return newKeyErr
		}

		u.ApiKeyId = newKey.Id
		updateErr := u.update()

		if updateErr != nil {
			return updateErr
		}
	}

	return nil
}

func (u *User) getKey() (*ApiKey, error) {
	key, err := getKeyById(u.ApiKeyId)

	if err != nil {
		return nil, err
	}

	return key, nil
}

func (u *User) update() error {
	results, err := myDb.Exec("UPDATE user SET Email=?, Password=?, ApiKeyId=?, Admin=? WHERE Id=?", u.Email, u.Password, u.ApiKeyId, u.Admin, u.Id)

	if err != nil {
		return err
	} else {
		rows, _ := results.RowsAffected()

		if rows > 0 {
			return nil
		}
	}

	return errors.New("unable to update user")
}

func getUserByEmail(email string) (*User, error) {
	var u User
	err := myDb.QueryRow("SELECT * FROM user WHERE email = ? LIMIT 1", email).Scan(&u.Id, &u.Email, &u.Password, &u.ApiKeyId, &u.Admin)
	switch {
	case err == sql.ErrNoRows:
		return nil, errUserNotFound
	case err != nil:
		return nil, err
	default:
		return &u, nil
	}
}

func (u *UserAuth) validateFields() error {

	fErr := u.fieldsCheck()

	if fErr != nil {
		return fErr
	}

	sErr := u.dbSchemaCheck()

	if sErr != nil {
		return sErr
	}

	return nil
}

func (u *UserAuth) fieldsCheck() error {

	if strings.TrimSpace(u.Email) == "" || strings.TrimSpace(u.Password) == "" {
		return errInvalidUserInfo
	}

	return nil
}

func (u *UserAuth) dbSchemaCheck() error {

	if utf8.RuneCountInString(u.Email) > 128 {
		return errInvalidEmailLength
	}

	return nil
}
