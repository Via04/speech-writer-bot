FROM golang:1.20.1-bullseye
ENV LANGUAGE="en"
WORKDIR /app
RUN git clone https://github.com/Via04/speech-writer-bot
WORKDIR /app/speech-writer-bot
RUN go mod download
RUN go build
CMD ["./speech-writer"]