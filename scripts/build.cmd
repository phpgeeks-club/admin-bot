set GOOS=linux
set GOARCH=amd64
go build -o geeksonator -trimpath -ldflags "-s -w"
