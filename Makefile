all:
	go build -o discovery ./example/address/
	go build -o http ./example/http

clean:
	rm -rf ./vendor

deploy:
	scp discovery root@pm175.dev.sedmax.ru:/tmp
	scp http root@pm175.dev.sedmax.ru:/tmp

vendor:
	go mod vendor