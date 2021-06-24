PUBKEY = /home/hsowan/workspace/genericfs/public-key.pem

run:
	go run -trimpath main.go -pubkey $(PUBKEY)

test:
	go test -v -count 1 github.com/wanhuasong/genericfs/utils
