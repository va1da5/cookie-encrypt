FROM cgr.dev/chainguard/go AS builder
COPY . /app
RUN cd /app && go build -o cookie-encrypt .

FROM cgr.dev/chainguard/glibc-dynamic

COPY --from=builder /app/cookie-encrypt /usr/bin/
CMD ["/usr/bin/cookie-encrypt"]