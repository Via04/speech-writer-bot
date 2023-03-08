FROM golang:1.20.1-bullseye
ENV LANGUAGE="en"
RUN yes | apt install ffmpeg
WORKDIR /app
ARG CACHE_DATE=2023-03-01
RUN git clone https://github.com/Via04/speech-writer-bot
WORKDIR /app/speech-writer-bot
RUN go mod download
RUN go build
CMD ["./speech-writer"]