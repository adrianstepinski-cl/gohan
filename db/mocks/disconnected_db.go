package mock_db

import (
	db "github.com/cloudwan/gohan/db"
	options "github.com/cloudwan/gohan/db/options"
	transaction "github.com/cloudwan/gohan/db/transaction"
	schema "github.com/cloudwan/gohan/schema"
	errors "github.com/pkg/errors"
)

type DBWithoutTransactions struct {
	OriginalDatabase db.DB
}

func (db DBWithoutTransactions) Connect(arg1 string, arg2 string, arg3 int) error {
	return db.OriginalDatabase.Connect(arg1, arg2, arg3)
}

func (db DBWithoutTransactions) Close() {
	db.OriginalDatabase.Close()
}

func (db DBWithoutTransactions) BeginTx(options ...transaction.Option) (transaction.Transaction, error) {
	return nil, errors.Errorf("mock DB.BeginTx")
}

func (db DBWithoutTransactions) RegisterTable(s *schema.Schema, cascade bool, migrate bool) error {
	return db.OriginalDatabase.RegisterTable(s, cascade, migrate)
}

func (db DBWithoutTransactions) DropTable(s *schema.Schema) error {
	return db.OriginalDatabase.DropTable(s)
}

func (db DBWithoutTransactions) Options() options.Options {
	return db.OriginalDatabase.Options()
}
