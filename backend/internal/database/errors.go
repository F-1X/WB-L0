package database

import "errors"

var ErrNotExist = errors.New("database: not exist")
var ErrBeginTransaction = errors.New("database: failed begin transaction")
var ErrCommitTransaction = errors.New("database: failed to commit trasaction")
var ErrRollbackTransaction = errors.New("database: failed to rollback transaction")
var ErrOrderExists = errors.New("database: order already exist")
