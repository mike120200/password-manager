
build:
	go build -o ./mac/pm
dev:
	go mod tidy

windows:
	GOOS=windows GOARCH=amd64 go build -o ./windows/pm.exe
linux:
	GOOS=linux GOARCH=amd64 go build -o ./linux/pm