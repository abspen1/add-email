package function

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"regexp"
	"strings"

	"github.com/gomodule/redigo/redis"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// Handle function
func Handle(w http.ResponseWriter, r *http.Request) {
	var info []byte

	if r.Body != nil {
		defer r.Body.Close()
		info, _ = ioutil.ReadAll(r.Body)
	}

	// Get our JSON into nested map structure
	var data = map[string]map[string]string{}
	_ = json.Unmarshal(info, &data)
	name := (data["info"]["name"])
	email := (data["info"]["email"])

	// Check if email and name are valid, if so check database then add to database
	if !isEmailValid(email) {
		w.Write([]byte("Invalid email address"))
	} else if name == "" {
		w.Write([]byte("Enter your name"))
	} else {
		secret, _ := getDBSecret("redis-password")
		c, err := redis.Dial("tcp", "192.168.1.6:6379")
		if err != nil {
			log.Fatal(err)
		}
		defer c.Close()

		_, err = c.Do("AUTH", secret)
		if err != nil {
			log.Fatal(err)
		}

		var exists string
		exists, _ = redis.String(c.Do("HGET", "emails", email))
		if exists == "" {
			c.Do("HSET", "emails", email, name)
			message := fmt.Sprintf("Added %s: %s to database", name, email)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(message))
		} else {
			w.Write([]byte(fmt.Sprintf("%s, you're already signed up for emails :)", name)))
		}
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

func getDBSecret(secretName string) (secretBytes []byte, err error) {
	// read from the openfaas secrets folder
	secretBytes, err = ioutil.ReadFile("/var/openfaas/secrets/" + secretName)
	if err != nil {
		// read from the original location for backwards compatibility with openfaas <= 0.8.2
		secretBytes, err = ioutil.ReadFile("/run/secrets/" + secretName)
	}
	return secretBytes, err
}
