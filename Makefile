
mac:
	rm -f ./mac/pm
	go build -o ./mac/pm
	zip -r mac.zip ./mac/pm
dev:
	go mod tidy

windows:
	rm -f ./windows/pm.exe
	GOOS=windows GOARCH=amd64 go build -o ./windows/pm.exe
	zip -r windows.zip ./windows/pm.exe
linux:
	rm -f ./linux/pm
	GOOS=linux GOARCH=amd64 go build -o ./linux/pm
	zip -r linux.zip ./linux/pm

build:
	go build -o ./pm
	mv ./pm ~/pm_dir