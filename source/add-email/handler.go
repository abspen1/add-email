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

// Info struct
type Info struct {
	Name  string
	Email string
}

// Body struct
type Body struct {
	info Info
}

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// Handle function
func Handle(w http.ResponseWriter, r *http.Request) {
	var info []byte

	if r.Body != nil {
		defer r.Body.Close()
		body, _ := ioutil.ReadAll(r.Body)
		info = body
	}
	
	var requestBody Body
	
	requestBody = json.Unmarshall(info, &requestBody)

	// Creating the maps for JSON
	m := map[string]interface{}{}

	// Parsing/Unmarshalling JSON encoding/json
	_ = json.Unmarshal(info, &m)
	in.parseMap(m)

	// Check if email is valid
	if e := string(in.Email); !isEmailValid(e) {
		w.Write([]byte("Invalid email address"))
	} else if in.Name == "" {
		w.Write([]byte("Enter your name"))
	} else {
		in.success(w)
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

func (in *Info) parseMap(aMap map[string]interface{}) {
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

func (in Info) success(w http.ResponseWriter) {
	secret, _ := getDBSecret("redis-password")
	c, err := redis.Dial("tcp", "192.168.1.6:6379")
	if err != nil {
		log.Fatal(err)
	}

	_, err = c.Do("AUTH", secret)
	if err != nil {
		log.Fatal(err)
	}
	exists, _ := redis.Bool(c.Do("EXISTS", "emails", in.Email))
	if exists {
		w.Write([]byte("You're already signed up for emails :)"))
	} else {
		fmt.Println("Exists works with hashes")
		c.Do("HSET", "emails", in.Email, in.Name)
		message := fmt.Sprintf("Added %s: %s to database", in.Name, in.Email)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(message))
	}
	defer c.Close()
}
