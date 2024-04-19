# HTTP Cookie Encrypt

A simple reverse proxy designed specifically to encrypt cookies.

## Environment Variables

- `PORT` (_Default: 8000_) - On which port to listen
- `SECRET` (_Required, min 6 characters_) - Secret used to encrypt cookie values
- `BACKEND` (Example: http://google.com) - Destination backend including scheme
- `DOMAIN` (Example: mydomain.com) - Overwrite cookie domain field with custom value
- `HTTPONLY` (Default: false) - Toogle `HttpOnly` option on cookies
- `SECURE` (Default: false) - Toogle `Secure` option on cookies
- `IGNORE` (Example: uuid,load-balancer) - Comma separated list of Cookie names that should be ignored

## References

- [yowu/HttpProxy.go](https://gist.github.com/yowu/f7dc34bd4736a65ff28d)
- [JalfResi/revprox.go](https://gist.github.com/JalfResi/6287706)
- [Parsing the Cookie and Set-Cookie headers with Go](https://www.jvt.me/posts/2022/04/07/go-cookie-header/)
- [Learn Golang encryption and decryption](https://blog.logrocket.com/learn-golang-encryption-decryption/)
- [manishtpatel/main.go](https://gist.github.com/manishtpatel/8222606)
