export GOROOT=/usr/local/go
export GOPATH=$HOME/go
export PATH=$GOPATH/bin:$GOROOT/bin:$PATH

export GIPHY_TOKEN=$GIPHY_TOKEN
export TELEGRAM_TOKEN=$TELEGRAM_TOKEN
export YANDEX_OAUTH=$YANDEX_OAUTH

docker compose up -d

go build cmd/app/app.go
