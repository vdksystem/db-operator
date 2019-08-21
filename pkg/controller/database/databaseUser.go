package database

import "github.com/sethvargo/go-password/password"

type user struct {
	username string
	password string
}

func genPassword() (string, error) {
	return password.Generate(12, 5, 3, false, false)
}
