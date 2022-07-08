package server

import (
	"errors"
)

var ErrAdminOnly = errors.New("admin only")
