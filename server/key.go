package main

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Apikey struct {
	Id     string
	Value  string
	Status string
}

func isKeyValid(req *http.Request) bool {

	key := req.FormValue("apiKey")

	if key != "" {
		var k Apikey
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
