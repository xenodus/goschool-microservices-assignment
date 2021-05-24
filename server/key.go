package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type ApiKey struct {
	Id     int
	Value  string
	Status string
}

func getKeyById(Id int) (*ApiKey, error) {
	var k ApiKey
	err := myDb.QueryRow("SELECT * FROM apikey WHERE Id = ? LIMIT 1", Id).Scan(&k.Id, &k.Value, &k.Status)
	switch {
	case err == sql.ErrNoRows:
		return nil, errKeyNotFound
	case err != nil:
		return nil, err
	default:
		return &k, nil
	}
}

func getKeyByValue(key string) (*ApiKey, error) {
	var k ApiKey
	err := myDb.QueryRow("SELECT * FROM apikey WHERE Value = ? LIMIT 1", key).Scan(&k.Id, &k.Value, &k.Status)
	switch {
	case err == sql.ErrNoRows:
		return nil, errKeyNotFound
	case err != nil:
		return nil, err
	default:
		return &k, nil
	}
}

func isKeyValid(req *http.Request) bool {

	key := req.FormValue("apiKey")

	if key != "" {
		var k ApiKey
		err := myDb.QueryRow("SELECT * FROM apikey WHERE Value = ? AND status = 'active' LIMIT 1", key).Scan(&k.Id, &k.Value, &k.Status)

		switch {
		case err != nil:
			return false
		default:
			return true
		}
	}

	return false
}

func handleKeyInvalid(res http.ResponseWriter) {
	printJSONResponse(res, JSONResponse{"error", http.StatusUnauthorized, errInvalidApiKey.Error()})
}

func generateKey() (string, error) {
	b := make([]byte, keyLength)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	key := hex.EncodeToString(b)

	results, dbErr := myDb.Exec("INSERT INTO apikey VALUES (NULL, ?,'active')", key)

	if dbErr != nil {
		return "", dbErr
	} else {
		rows, _ := results.RowsAffected()

		if rows > 0 {
			return key, nil
		}
	}

	return "", errors.New("unable to generate key")
}

func invalidateKey(key string) error {
	results, err := myDb.Exec("UPDATE apikey SET status='inactive' WHERE Value = ?", key)

	if err != nil {
		return err
	} else {
		rows, _ := results.RowsAffected()

		if rows > 0 {
			return nil
		}
	}

	return errors.New("unable to invalidate key")
}
