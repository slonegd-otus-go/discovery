all:
	go build -o discovery ./example/
	go build -o ./update_signal/update_signal ./update_signal/main.go

deploy:
	scp discovery root@pm175.dev.sedmax.ru:/tmp
	scp ./update_signal/update_signal root@pm175.dev.sedmax.ru:/tmp

	