package util

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image/jpeg"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/GirishBhutiya/gOpenfiatServer/app/config"
	"github.com/google/uuid"
)

const TimeFormat = "2006-01-02 15:04:05"

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
	resp, err := http.Get(fmt.Sprintf(config.OTPAPIURL, fmt.Sprintf("%s%d", "+91", phonenumber)))
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
func SaveProfilePic(imgname, imgBase64 string) (string, error) {
	unbased, err := base64.StdEncoding.DecodeString(imgBase64)
	if err != nil {
		return "", errors.New("cannot decode b64")
	}
	r := bytes.NewReader(unbased)

	im, err := jpeg.Decode(r)
	if err != nil {
		return "", errors.New("bad jpeg")
	}
	imgPath := fmt.Sprintf("profilepic/%s%s", imgname, ".jpeg")
	f, err := os.OpenFile(imgPath, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return "", errors.New("cannot open file")
	}

	return imgPath, jpeg.Encode(f, im, nil)
}
func ConvertUUIDToAny(ids []uuid.UUID) []any {
	var newIds []interface{}
	for _, id := range ids {
		newIds = append(newIds, id)
	}
	return newIds
}
func GenerateShortKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLength = 6

	rand.Seed(time.Now().UnixNano())
	shortKey := make([]byte, keyLength)
	for i := range shortKey {
		shortKey[i] = charset[rand.Intn(len(charset))]
	}
	return string(shortKey)
}

/* func ConvertOrderStringToOrder(order *model.OrderHandlerString) (model.OrderHandler, error) {
	var or model.OrderHandler
	parsed, err := time.Parse(TimeFormat, order.TimeLimit)
	if err != nil {
		return or, err
	}
	or.TimeLimit = parsed
	or.ID = order.ID
	or.FiatAmount = order.FiatAmount
	or.MinAmount = order.MinAmount
	or.Type = order.Type
	or.UserId = order.UserId
	or.Price = order.Price

	return or, nil
}
func ConvertTradeStringToTrade(trade *model.TradeHandlerUser) (model.TradeHandler, error) {
	var tr model.TradeHandler
	parsed, err := time.Parse(TimeFormat, trade.TradeTime)
	if err != nil {
		return tr, err
	}
	tr.ID = trade.ID
	tr.Orderid = trade.Orderid
	tr.TradeTime = parsed
	tr.Method = trade.Method

	return tr, nil
} */
