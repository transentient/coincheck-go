package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
	"strings"
)

type CoinCheck struct {
	accessKey string
	secretKey string
	account Account
	bank_account BankAccount
}

func (c CoinCheck) NewClient(accessKey string, secretKey string) CoinCheck{
	c.accessKey = accessKey
	c.secretKey = secretKey
	c.account = Account{&c}
	c.bank_account = BankAccount{&c}
	return c
}

func (c CoinCheck) Request(method string, path string, param string) string {
	url := "https://coincheck.jp/" + path
	nonce := strconv.FormatInt(CreateNonce(), 10)
	message := nonce + url + param
	req := &http.Request{}
	if method == "POST" {
		payload := strings.NewReader(param)
		req, _ = http.NewRequest(method, url, payload)
	} else {
		req, _ = http.NewRequest(method, url, nil)
	}
	signature := ComputeHmac256(message, c.secretKey)
	req.Header.Add("access-key", c.accessKey)
	req.Header.Add("access-nonce", nonce)
	req.Header.Add("access-signature", signature)
	req.Header.Add("content-type", "application/json")
	req.Header.Add("cache-control", "no-cache")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return ""
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))

	return string(body)
}

//create nonce by milliseconds
func CreateNonce() int64 {
	return time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}

//create signature
func ComputeHmac256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}


