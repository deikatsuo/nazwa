package dbquery

import (
	"nazwa/misc"

	"github.com/jmoiron/sqlx"
)

// DB variabel database global
var DB *sqlx.DB

var log = misc.Log
