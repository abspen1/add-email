package function

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"regexp"
	"strings"

	"github.com/gomodule/redigo/redis"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func Handle(w http.ResponseWriter, r *http.Request) {
	var email []byte
	var name []byte

	if r.Body != nil {
		defer r.Body.Close()

		body, _ := ioutil.ReadAll(r.Body.info.name)
		name = body
		body, _ = ioutil.ReadAll(r.Body.info.name)
		email = body
	}

	// Check if email is valid
	if e := string(email); !isEmailValid(e) {
		// http code 400 is for a bad request
		w.WriteHeader(http.Error(w, "Email not valid", 400))
	} 
	// Make sure name isn't empty
	else if string(name) == "" {
		w.WriteHeader(http.Error(w, "Enter your name", 400))
	} 
	else {
		secret, _ := getDBSecret("redis-password")
		c, err := redis.Dial("tcp", "192.168.1.6:6379")
		if err != nil {
			log.Fatal(err)
		}
		response, err := c.Do("AUTH", secret)
		if err != nil {
			log.Fatal(err)
		}

		c.Do("HSET", "emails", string(email), string(name))
		defer c.Close()

		w.WriteHeader(http.StatusOK)
	}
}

func isEmailValid(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	if !emailRegex.MatchString(e) {
		return false
	}
	parts := strings.Split(e, "@")
	mx, err := net.LookupMX(parts[1])
	if err != nil || len(mx) == 0 {
		return false
	}
	return true
}
