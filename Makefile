automation: main.go
	env GOOS=linux GOARCH=arm GOARM=6 go build -o automation main.go
	scp automation  pi@192.168.2.109:automation
	rm automation

main: main.go
	go build main.go
	./main config.json
	rm main
	