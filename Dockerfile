FROM golang:1.20.2-alpine
ENV LANGUAGE="en"
WORKDIR /app
RUN apk add  --no-cache ffmpeg
RUN apk add --no-cache git
ARG CACHE_DATE=2023-03-01
RUN git clone https://github.com/Via04/speech-writer-bot
WORKDIR /app/speech-writer-bot
RUN go mod download
RUN go build
CMD ["./speech-writer"]