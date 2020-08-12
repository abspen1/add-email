package function

import (
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// Handle function
func Handle(w http.ResponseWriter, r *http.Request) {
	// var email []byte
	var name []byte

	if r.Body != nil {
		defer r.Body.Close()

		body, _ := ioutil.ReadAll(r.Body)
		name = body
		// body, _ = ioutil.ReadAll(r.Body)
		// email = body
	}

	// Check if email is valid
	// if e := string(email); !isEmailValid(e) {
	// 	w.Write([]byte("Invalid email"))
	// } else if string(name) == "" {
	// 	w.Write([]byte("Enter your name"))
	// } else {
	// 	secret, _ := getDBSecret("redis-password")
	// 	c, err := redis.Dial("tcp", "192.168.1.6:6379")
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	response, err := c.Do("AUTH", secret)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Println("Connected!", response)

	// 	c.Do("HSET", "emails", string(email), string(name))
	// 	defer c.Close()

	// 	w.WriteHeader(http.StatusOK)
	// }
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(name))

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
