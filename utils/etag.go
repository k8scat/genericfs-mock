package utils

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"io"
)

const (
	BLOCK_BITS = 22 // Indicate that the blocksize is 4M
	BLOCK_SIZE = 1 << BLOCK_BITS
)

func BlockCount(fsize int) int {
	return int((fsize + (BLOCK_SIZE - 1)) >> BLOCK_BITS)
}

func CalSha1(b []byte, r io.Reader) ([]byte, error) {
	h := sha1.New()
	_, err := io.Copy(h, r)
	if err != nil {
		return nil, err
	}
	return h.Sum(b), nil
}

func GetEtag(b []byte) (etag string, err error) {
	fsize := len(b)
	blockCnt := BlockCount(fsize)
	sha1Buf := make([]byte, 0, 21)

	r := bytes.NewReader(b)
	if blockCnt <= 1 { // file size <= 4M
		sha1Buf = append(sha1Buf, 0x16)
		sha1Buf, err = CalSha1(sha1Buf, r)
		if err != nil {
			return
		}
	} else { // file size > 4M
		sha1Buf = append(sha1Buf, 0x96)
		sha1BlockBuf := make([]byte, 0, blockCnt*20)
		for i := 0; i < blockCnt; i++ {
			body := io.LimitReader(r, BLOCK_SIZE)
			sha1BlockBuf, err = CalSha1(sha1BlockBuf, body)
			if err != nil {
				return
			}
		}
		sha1Buf, _ = CalSha1(sha1Buf, bytes.NewReader(sha1BlockBuf))
	}
	etag = base64.URLEncoding.EncodeToString(sha1Buf)
	return
}
