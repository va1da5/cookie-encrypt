# 🔐 HTTP Cookie Encrypt

A simple reverse proxy designed specifically to encrypt cookies.

## Container Image

The prebuilt container image can be found in [DokerHub.io](https://hub.docker.com/r/va1da5/cookie-encrypt)

## Environment Variables

- `SERVER` (_Default: 0.0.0.0:8000_) - On which port to listen
- `SECRET` (_Required_) - Secret used to encrypt cookie values
- `BACKEND` (_Example: http://google.com_) - Destination backend including scheme
- `DOMAIN` (_Example: mydomain.com_) - Overwrite cookie domain field with custom value
- `HTTPONLY` (_Default: false_) - Toogle `HttpOnly` option on cookies
- `SECURE` (_Default: false_) - Toogle `Secure` option on cookies
- `IGNORE` (_Example: uuid,PHPSESSID_) - Comma separated list of Cookie names that should be ignored

## Development

```
cp sample.env .env
export $(cat .env |xargs)

# start local server
go run main.go
```

## References

- [yowu/HttpProxy.go](https://gist.github.com/yowu/f7dc34bd4736a65ff28d)
- [JalfResi/revprox.go](https://gist.github.com/JalfResi/6287706)
- [Parsing the Cookie and Set-Cookie headers with Go](https://www.jvt.me/posts/2022/04/07/go-cookie-header/)
- [Learn Golang encryption and decryption](https://blog.logrocket.com/learn-golang-encryption-decryption/)
- [manishtpatel/main.go](https://gist.github.com/manishtpatel/8222606)
