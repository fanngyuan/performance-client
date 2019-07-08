package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
)

type Content struct {
	Nonce int    `json:"nonce"`
	Bytes string `json:"bytes"`
}

func main() {
	urls := flag.String("url", "http://localhost:8080", "url")
	times := flag.Int("times", 3, "times")
	flag.Parse()

	urls_slice := strings.Split(*urls, ",")

	var buf bytes.Buffer
	for i := 0; i < *times; i++ {
		url := urls_slice[rand.Intn(len(urls_slice))]
		bufLen := rand.Intn(512)
		bytes, bs, err := GenerateRandomString(bufLen)
		if err != nil {
			fmt.Println(err)
			continue
		}
		content := Content{
			Nonce: i,
			Bytes: bytes,
		}
		_, err = buf.Write(bs)
		if err != nil {
			fmt.Println(err)
			break
		}
		result := sha256.Sum256(buf.Bytes())
		buf.Reset()
		_, err = buf.Write(result[:])
		if err != nil {
			fmt.Println(err)
			break
		}
		contentBytes, err := json.Marshal(content)
		if err != nil {
			fmt.Println(err)
			continue
		}
		http_post(url, contentBytes)
	}
	fmt.Printf("hash hex is %s \r\n", hex.EncodeToString(buf.Bytes()))
}

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func GenerateRandomString(s int) (string, []byte, error) {
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), b, err
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
