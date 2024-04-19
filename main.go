package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"unicode"
)

var Secret, Backend, Domain, Ignore, Port string
var HttpOnly = false
var Secure = false

func modifyResponse(resp *http.Response) (err error) {
	if resp.Cookies() != nil {

		var cookies []*http.Cookie

		cookies = []*http.Cookie{}

		for _, c := range resp.Cookies() {

			if !strings.Contains(Ignore, c.Name) {
				if len(Domain) > 0 {
					c.Domain = Domain
				}

				c.HttpOnly = HttpOnly
				c.Secure = Secure

				out, err := Encrypt(c.Value, Secret)
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

func main() {
	Port = os.Getenv("PORT")
	if len(Port) == 0 {
		Port = "8000"
	}

	Secret = fmt.Sprintf("%-*s", aes.BlockSize, os.Getenv("SECRET"))
	Backend = os.Getenv("BACKEND")

	if len(Backend) < 1 {
		panic("BACKEND environment variable is required")
	}

	if len(Secret) < 6 {
		panic("SECRET environment variable value length must be longer that 5 characters")
	}

	Domain = os.Getenv("DOMAIN")

	HttpOnly = len(os.Getenv("HTTPONLY")) > 0
	Secure = len(os.Getenv("SECURE")) > 0
	Ignore = os.Getenv("IGNORE")

	remote, err := url.Parse(Backend)
	if err != nil {
		panic(err)
	}

	handler := func(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			var cookies []*http.Cookie
			cookies = []*http.Cookie{}

			var error = false

			for _, c := range r.Cookies() {
				if strings.Contains(Ignore, c.Name) {
					cookies = append(cookies, c)
					continue
				}

				out, err := Decrypt(c.Value, Secret)

				if err != nil {
					fmt.Println(err)
					error = true
					break
				}
				if !isASCII(out) {
					fmt.Println("Cookie decoding failed")
					error = true
					break
				}

				c.Value = out
				cookies = append(cookies, c)
			}

			if !error {
				r.Header.Set("Cookie", joinCookies(cookies))
			}

			r.Host = remote.Host
			p.ServeHTTP(w, r)

		}
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.ModifyResponse = modifyResponse

	http.HandleFunc("/", handler(proxy))
	err = http.ListenAndServe(":"+Port, nil)
	if err != nil {
		panic(err)
	}
}

func Encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func Decode(s string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return []byte{}, err
	}
	return data, nil
}

func Encrypt(text, secret string) (string, error) {
	block, err := aes.NewCipher([]byte(secret))
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
	return Encode(cipherText), nil
}

func Decrypt(text, secret string) (string, error) {
	block, err := aes.NewCipher([]byte(secret))
	if err != nil {
		return "", err
	}
	cipherText, err := Decode(text)
	if err != nil {
		return "", err
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
