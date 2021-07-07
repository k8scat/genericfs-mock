PUBKEY = /home/hsowan/workspace/genericfs/public-key.pem

run:
	go run -trimpath main.go -pubkey $(PUBKEY)

build:
	go build -trimpath -o genericfs main.go

test:
	go test -v -count 1 github.com/wanhuasong/genericfs/utils

upload:
	scp genericfs new-marsdev:/tmp/genericfs/genericfs
	scp public-key.pem new-marsdev:/tmp/genericfs/public-key.pem
	scp config.json new-marsdev:/tmp/genericfs/config.json
