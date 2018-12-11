package server_up_test

import (
	"github.com/cloudwan/gohan/db"
	"github.com/cloudwan/gohan/db/dbutil"
	"github.com/cloudwan/gohan/db/options"
	"github.com/cloudwan/gohan/schema"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"log"
	"os"
	"testing"
)

var (
	testDB    db.DB
	conn, dbType string
	whitelist = map[string]bool{
		"healthcheck": true,
	}
)

func TestServer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Server Up Suite")
}

func setupDB() {
	if os.Getenv("MYSQL_TEST") == "true" {
		conn = "root@/gohan_test"
		dbType = "mysql"
	} else {
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		conn = dir + "/test.db"
		dbType = "sqlite3"
	}
	var err error
	testDB, err = dbutil.ConnectDB(dbType, conn, db.DefaultMaxOpenConn, options.Default())
	Expect(err).ToNot(HaveOccurred(), "Failed to connect database.")
	if os.Getenv("MYSQL_TEST") == "true" {
		err = StartTestServer("../../server_test_mysql_config.yaml")
	} else {
		err = StartTestServer("../../server_test_config.yaml")
	}
	Expect(err).ToNot(HaveOccurred(), "Failed to start test server.")
}

func cleanupDB() {
	schema.ClearManager()
	os.Remove(conn)
}

var _ = Describe("Suit set up and tear down", func() {
	var _ = BeforeSuite(func() {
		setupDB()
	})

	var _ = AfterSuite(func() {
		cleanupDB()
	})
})
