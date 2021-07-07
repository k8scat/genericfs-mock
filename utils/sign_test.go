package utils

import (
	"fmt"
	"log"
	"testing"

	"github.com/wanhuasong/genericfs/config"
)

var (
	origin = "e=1625645423&t=1625641823&uuid=ButJmRBi"
	sig    string
)

func TestSign(t *testing.T) {
	Privkey = "/home/hsowan/workspace/genericfs/private-key.pem"

	// params := map[string]string{
	// 	"t": "1625564056",

	// 	"callback_url":  "https://devapi.myones.net/project/S1051/res/uploadcallback",
	// 	"callback_body": "hash=$(etag)&type=attachment&name=$(fname)&size=$(fsize)&mime=$(mimeType)&ext=$(ext)&exif=$(exif)&width=$(imageInfo.width)&height=$(imageInfo.height)&user=B2p77rtd&team=PuzvUjVK&resource=12Tet2mv&token=5b605155067aa4c1f998f06499234b76ce195d912a38fa0b29f090826d99ef73ed7c1d52a4127b97aa181de211fc373baaba845a8c89c472010c2b05ae923bc6013df6c1ba68affaa9a287a2c220312b88304ac610d6b30d699f477b2eb9e589b6478864dbedaa6a53716e1f2185ff4fc190f1056799d98e7ea4c446390a2176",

	// 	"uuid":           "12Tet2mv",
	// 	"reference_type": "8",
	// 	"reference_id":   "AdWX3aSz",
	// 	"team_uuid":      "PuzvUjVK",
	// 	"project_uuid":   "B2p77rtdRWU1q6A4",
	// 	"owner_uuid":     "B2p77rtd",
	// 	"modifier":       "A63euYZC",
	// 	"type":           "1",
	// 	"source":         "1",
	// 	"ext_id":         "FmL5XAJtJnP70JCPZECt74W2oAvr",
	// 	"name":           "Batrider.png",
	// 	"status":         "1",
	// 	"create_time":    "1559622205224752",
	// 	"description":    "",
	// 	"modify_time":    "1559622254814176",
	// }
	// origin = SortParams(params)
	log.Printf("Origin: %s", origin)

	var err error
	sig, err = Sign(origin)
	if err != nil {
		panic(err)
	}
	fmt.Printf("sig: %s\n", sig)
	fmt.Printf("sig len: %d\n", len(sig))
}

func TestVerifySig(t *testing.T) {
	origin = "e=1625469458&hash=Fk2wuo7kLZUYDWgJrY7FaGhHvWXJ&op=imageMogr2/auto-orient/thumbnail/x32&t=1625465858"
	sig = "5298797ba794b5ff6ad2a91b3011a67e127fca0c20a8a97f829729d8f272abdc599e8de6de404727864784a62f4493940a8d4f02cf6bc055fe26770974a75ad3e6e4c167c1c82db9526a430c2a0a0abf072fddf9ad8070c6285d2cdf84b42bc1e757e4999d67b76629192dbea4687ab42fc93f06b63c329dd35a28c84913c1a3"
	config.Cfg.PublicKey = "/home/hsowan/workspace/genericfs/public-key.pem"
	if err := VerifySig(origin, sig); err != nil {
		panic(err)
	}
}
