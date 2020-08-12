package function

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"regexp"
	"strings"

	"github.com/gomodule/redigo/redis"
)

// Info struct
type Info struct {
	Name  string
	Email string
}

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// Handle function
func Handle(w http.ResponseWriter, r *http.Request) error {
	var err error

	var info []byte

	if r.Body != nil {
		defer r.Body.Close()

		body, _ := ioutil.ReadAll(r.Body)
		info = body
	}

	// Creating out struct
	var in Info
	// Creating the maps for JSON
	m := map[string]interface{}{}

	// Parsing/Unmarshalling JSON encoding/json
	eRR := json.Unmarshal(info, &m)

	if eRR != nil {
		log.Fatal(eRR)
	}
	in.parseMap(m)

	// Check if email is valid
	if e := string(in.Email); !isEmailValid(e) {
		w.Write([]byte("Enter a valid email address"))
		return errors.New("Invalid email")
	} else if in.Name == "" {
		w.Write([]byte("Enter your name"))
		return errors.New("Name field left empty")
	}

	secret, _ := getDBSecret("redis-password")
	c, err := redis.Dial("tcp", "192.168.1.6:6379")
	if err != nil {
		log.Fatal(err)
	}

	response, err := c.Do("AUTH", secret)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected!", response)

	c.Do("HSET", "emails", in.Email, in.Name)
	defer c.Close()

	message := fmt.Sprintf("Added %s: %s to database", in.Name, in.Email)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(message))
	return err
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

func (in *Info) parseMap(aMap map[string]interface{}) bool {
	for key, val := range aMap {
		switch concreteVal := val.(type) {
		case map[string]interface{}:
			in.parseMap(val.(map[string]interface{}))
		default:
			// fmt.Println(key, ":", concreteVal)
			if key == "name" {
				in.Name = concreteVal.(string)
			} else if key == "email" {
				in.Email = concreteVal.(string)
			}
		}
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
