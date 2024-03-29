package pgx2txt

import (
	"database/sql"

	// registers "pgx"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewDB(conn string) (*sql.DB, error) { return sql.Open("pgx", conn) }
