package model

import (
	"database/sql"
)

type (
	User struct {
		UserID     int    `pg:",pk"`
		UserName   string `pg:"type:varchar(50),unique"`
		FirstName  string
		LastName   string
		Email      string
		UserStatus string `pg:"type:varchar(1)"`
		Department sql.NullString
	}
)

const (
	Active     = "A"
	Inactive   = "I"
	Terminated = "T"
)
