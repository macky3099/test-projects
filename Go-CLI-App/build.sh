#Windows:
GOOS=windows GOARCH=amd64 go build -o bin/bookingcli.exe main.go

#MacOS:
#64-bit
GOOS=darwin GOARCH=amd64 go build -o bin/bookingcli-macos-x86 main.go
#Arm based
GOOS=darwin GOARCH=arm64 go build -o bin/bookingcli-macos-arm main.go

#Linux:
#64-bit
GOOS=linux GOARCH=amd64 go build -o bin/bookingcli-linux-x86 main.go
#Arm based
GOOS=linux GOARCH=arm64 go build -o bin/bookingcli-linux-arm main.go
