package database

import "database/sql"

type DBClient struct {
	*sql.DB
}
