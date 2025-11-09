# Momon

A simple money tracker using LINE messaging platform.

## Features

- Receive and send messages via LINE chatbot
- Track expenses and income

## Setup

- Copy `.env.example` to `.env` and configure your LINE credentials
- Run the server:

```bash
go run main.go
```

- run ngrok:

```bash
NGROK_AUTHTOKEN=<YOUR_NGROK_AUTHTOKEN> HTTP_PORT=8080 docker-compose up
```

The server will start on port 8080.
