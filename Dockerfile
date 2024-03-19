FROM golang:1.21 AS build

RUN go install bou.ke/staticfiles@latest

FROM alpine:3.18

COPY --from=build /go/bin/staticfiles /bin/

ENTRYPOINT ["staticfiles"]
