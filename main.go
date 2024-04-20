package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"unicode"
)

type Config struct {
	Secret, Backend, Domain, Ignore, Server string
	HttpOnly, Secure                        bool
}

var config Config

func main() {

	config.Server = os.Getenv("SERVER")
	if len(config.Server) == 0 {
		config.Server = "0.0.0.0:8000"
	}

	config.Secret = os.Getenv("SECRET")
	config.Backend = os.Getenv("BACKEND")

	if len(config.Backend) < 1 {
		panic("BACKEND environment variable is required")
	}

	if len(config.Secret) == 0 {
		panic("SECRET environment variable value length must be defined")
	}

	config.Domain = os.Getenv("DOMAIN")

	config.HttpOnly = len(os.Getenv("HTTPONLY")) > 0
	config.Secure = len(os.Getenv("SECURE")) > 0
	config.Ignore = os.Getenv("IGNORE")

	remote, err := url.Parse(config.Backend)
	if err != nil {
		panic(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.ModifyResponse = modifyResponse

	http.HandleFunc("/", handler(proxy))
	fmt.Printf("Server started on %s\n", config.Server)
	err = http.ListenAndServe(config.Server, nil)
	if err != nil {
		panic(err)
	}
}

func encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func decode(s string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return []byte{}, err
	}
	return data, nil
}

func encrypt(text, secret string) (string, error) {
	hash := sha256.Sum256([]byte(secret))
	block, err := aes.NewCipher(hash[:])
	if err != nil {
		return "", err
	}

	plainText := []byte(text)

	iv, err := generateIV()
	if err != nil {
		fmt.Print(err)
	}
	cipherText := append(iv, make([]byte, len(text))...)
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(cipherText[aes.BlockSize:], plainText)
	return encode(cipherText), nil
}

func decrypt(text, secret string) (string, error) {
	hash := sha256.Sum256([]byte(secret))
	block, err := aes.NewCipher(hash[:])
	if err != nil {
		return "", err
	}

	cipherText, err := decode(text)
	if err != nil {
		return "", err
	}

	if len(cipherText) < aes.BlockSize {
		return "", errors.New("cipherText too short")
	}

	iv := cipherText[:aes.BlockSize]
	cfb := cipher.NewCFBDecrypter(block, iv)

	cipherText = cipherText[aes.BlockSize:]
	cfb.XORKeyStream(cipherText, cipherText)

	return string(cipherText), nil
}

func joinCookies(cookies []*http.Cookie) string {
	var parts []string
	for _, cookie := range cookies {
		parts = append(parts, cookie.String())
	}
	return strings.Join(parts, "; ")
}

func generateIV() ([]byte, error) {
	iv := make([]byte, aes.BlockSize)

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	return iv, nil
}

func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}

func modifyResponse(resp *http.Response) (err error) {
	if resp.Cookies() != nil {

		var cookies []*http.Cookie

		cookies = []*http.Cookie{}

		for _, c := range resp.Cookies() {

			if !strings.Contains(config.Ignore, c.Name) {
				if len(config.Domain) > 0 {
					c.Domain = config.Domain
				}

				c.HttpOnly = config.HttpOnly
				c.Secure = config.Secure

				out, err := encrypt(c.Value, config.Secret)
				if err != nil {
					panic(err)
				}
				c.Value = out
			}

			cookies = append(cookies, c)
		}

		resp.Header.Del("Set-Cookie")

		for _, c := range cookies {
			resp.Header.Add("Set-Cookie", c.String())
		}
	}

	return nil
}

func handler(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	remote, err := url.Parse(config.Backend)
	if err != nil {
		panic(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var cookies []*http.Cookie
		cookies = []*http.Cookie{}

		for _, c := range r.Cookies() {
			if strings.Contains(config.Ignore, c.Name) {
				cookies = append(cookies, c)
				continue
			}

			out, err := decrypt(c.Value, config.Secret)

			if err != nil || !isASCII(out) {
				continue
			}

			c.Value = out
			cookies = append(cookies, c)
		}

		r.Header.Set("Cookie", joinCookies(cookies))
		r.Host = remote.Host

		p.ServeHTTP(w, r)

	}
}
