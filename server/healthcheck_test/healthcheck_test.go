package healthcheck_test

import (
	"errors"
	"net"
	"net/http"
	"time"

	srv "github.com/cloudwan/gohan/server"
	. "github.com/cloudwan/gohan/server/test_common_part"
	. "github.com/onsi/ginkgo"

)

var (
	Server *srv.Server
)

var _ = Describe("Healthcheck test for active database", func() {
	It("should return 200 if db works for members", func() {
		TestURL("GET", HealthcheckURL, MemberTokenID, nil, http.StatusOK)
	})

	It("should return 200 if db works for non-members", func() {
		TestURL("GET", HealthcheckURL, "", nil, http.StatusOK)
	})
})

/*var _ = Describe("Healthcheck test for inactive database", func() {
	BeforeEach(func() {
		Server.Stop()
	})

	AfterEach(func() {
		Server.Start()
	})

	It("should return 503 if db isn't working for members", func() {
		TestURL("GET", HealthcheckURL, MemberTokenID, nil, http.StatusServiceUnavailable)
	})

	It("should return 503 if db isn't working for non-members", func() {
		TestURL("GET", HealthcheckURL, "", nil, http.StatusServiceUnavailable)
	})
})*/

func StartTestServer(config string) error {
	var err error
	Server, err = srv.NewServer(config)
	if err != nil {
		return err
	}

	go func() {
		err := Server.Start()
		if err != nil {
			panic(err)
		}
	}()

	retry := 3
	for {
		conn, err := net.Dial("tcp", Server.Address())
		if err == nil {
			conn.Close()
			break
		}
		retry--
		if retry == 0 {
			return errors.New("server not started")
		}
		time.Sleep(50 * time.Millisecond)
	}
	Server.SetRunning(true)

	return nil
}