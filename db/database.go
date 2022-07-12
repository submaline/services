package db

import "database/sql"

type DBClient struct {
	*sql.DB
}
