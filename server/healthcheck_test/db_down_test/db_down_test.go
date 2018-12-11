package server_up_test

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

var _ = Describe("Healthcheck test for inactive database", func() {
	It("should return 503 if db isn't working for members", func() {
		TestURL("GET", HealthcheckURL, MemberTokenID, nil, http.StatusServiceUnavailable)
	})

	It("should return 503 if db isn't working for non-members", func() {
		TestURL("GET", HealthcheckURL, "", nil, http.StatusServiceUnavailable)
	})
})

func StartTestServer(config string) error {
	var err error
	Server, err = srv.NewServer(config)
	if err != nil {
		return err
	}

	go func() {
		retry := 3
		for {
			err := Server.Start()
			retry--
			if err != nil && retry == 0{
				panic(err)
			}
			time.Sleep(50 * time.Millisecond)
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