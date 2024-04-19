FROM cgr.dev/chainguard/go AS builder
COPY . /app
RUN cd /app && go build -o cookie-encrypt .

FROM cgr.dev/chainguard/glibc-dynamic

ENV PORT 8000

COPY --from=builder /app/cookie-encrypt /usr/bin/
CMD ["/usr/bin/cookie-encrypt"]