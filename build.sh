rm -rf ./releases/*

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./releases/az-ops
