all:
	go build -o discovery ./example/

deploy:
	scp discovery root@pm175.dev.sedmax.ru:/tmp