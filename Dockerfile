FROM golang:1.11.1-alpine3.7

RUN apk add --update --no-cache git && go get github.com/bwmarrin/discordgo
ADD greg.go /greg.go
CMD go run /greg.go

