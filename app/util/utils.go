package util

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"

	"github.com/GirishBhutiya/gOpenfiatServer/app/config"
)

func CheckPasswordValidity(password string) bool {
	if regexp.MustCompile(`\s`).MatchString(password) {
		log.Println("in regex")
		return false
	}
	if len(password) < 8 {
		log.Println("in length")
		return false

	}
	return true
}
func SendOTP(phonenumber int, config *config.Config) ([]byte, error) {
	resp, err := http.Get(fmt.Sprintf(config.OTPAPIURL, fmt.Sprintf("%s%d", "+", phonenumber)))
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	//We Read the response body on the line below.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	log.Println(string(body[:]))

	return body, nil
}
func LenLoop(i int) int {
	if i == 0 {
		return 1
	}
	count := 0
	for i != 0 {
		i /= 10
		count++
	}
	return count
}
