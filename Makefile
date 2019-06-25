main: main.go
	go build main.go
	./main config.json
	rm main

automation: main.go
	env GOOS=linux GOARCH=arm GOARM=6 go build -o automation main.go