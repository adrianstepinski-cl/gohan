package healthcheck_test

import (
	"os"
	"testing"

	"github.com/cloudwan/gohan/db"
	"github.com/cloudwan/gohan/db/dbutil"
	"github.com/cloudwan/gohan/db/options"
	"github.com/cloudwan/gohan/schema"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	testDB    db.DB
	conn, dbType string
	whitelist = map[string]bool{
		"schema":    true,
		"policy":    true,
		"extension": true,
		"namespace": true,
		"version":   true,
		"healthcheck": true,
	}
)

func TestServer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Server Suite")
}

func SetupDB() {
	if os.Getenv("MYSQL_TEST") == "true" {
		conn = "root@/gohan_test"
		dbType = "mysql"
	} else {
		conn = "./test.db"
		dbType = "sqlite3"
	}
	var err error
	testDB, err = dbutil.ConnectDB(dbType, conn, db.DefaultMaxOpenConn, options.Default())
	Expect(err).ToNot(HaveOccurred(), "Failed to connect database.")
	if os.Getenv("MYSQL_TEST") == "true" {
		err = StartTestServer("/home/adrian/work/src/github.com/cloudwan/gohan/server/server_test_mysql_config.yaml")
	} else {
		err = StartTestServer("/home/adrian/work/src/github.com/cloudwan/gohan/server/server_test_config.yaml")
	}
	Expect(err).ToNot(HaveOccurred(), "Failed to start test server.")
}

func CleanupDB() {
	schema.ClearManager()
	os.Remove(conn)
}

var _ = Describe("Suit set up and tear down", func() {
	var _ = BeforeSuite(func() {
		SetupDB()
	})

	var _ = AfterSuite(func() {
		CleanupDB()
	})
})

