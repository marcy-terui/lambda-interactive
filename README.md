Run the bash shell interactively on Lambda Custom Runtime via ngrok
======

# Requirements
- Golang
- [ngrok.com](https://ngrok.com/)

# Setup

```sh
export SAM_PACKAGE_BUCKET=$YOUR_S3_BUCKET_NAME
export NGROK_AUTH_TOKEN=$YOUR_NGROK_AUTH_TOKEN
make init
```

# Deploy

```sh
make build
make deploy
```

# Start interaction

### 1. Start TCP Server on Lambda
```
make start
```

### 2. Check your ngrok endpoint and port on [ngrok dashboard](https://dashboard.ngrok.com/status)

### 3. Access the endpoint via telnet

```sh
telnet $YOUR_NGROK_ENDPOINT $YOUR_NGROK_PORT
```

# End interaction

```
exit
```
