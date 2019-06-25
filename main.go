package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
)

type Content struct {
	Nonce int    `json:"nonce"`
	Bytes string `json:"bytes"`
}

func main() {
	url := flag.String("url", "http://localhost:8080", "url")
	times := flag.Int("times", 3, "times")
	flag.Parse()

	for i := 0; i < *times; i++ {
		bytes, err := GenerateRandomString(32)
		if err != nil {
			fmt.Println(err)
			continue
		}
		content := Content{
			Nonce: i,
			Bytes: bytes,
		}
		contentBytes, err := json.Marshal(content)
		if err != nil {
			fmt.Println(err)
			continue
		}
		http_post(*url, contentBytes)
	}
}

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

func http_post(url string, jsonStr []byte) ([]byte, error) {
	count := 0
	for {
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
		if err == nil {
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err == nil {
				return body, nil
			} else {
				fmt.Println(err)
				count++
			}
		} else {
			fmt.Println(err)
			count++
		}
		if count == 2 {
			break
		}
	}
	return nil, nil
}
