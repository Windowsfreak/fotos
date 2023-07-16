package fotos

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"fotos/domain"
	"github.com/dvsekhvalnov/jose2go/base64url"
	"golang.org/x/crypto/sha3"
)

func MakeNonce() string {
	buff := make([]byte, 6)
	_, err := rand.Read(buff)
	if err != nil {
		println(fmt.Errorf("random number generation failed: %w", err).Error())
	}
	return base64.RawURLEncoding.EncodeToString(buff)
}

func hashFile(path string, nonce string) string {
	h := make([]byte, 32)
	d := sha3.NewShake128()
	d.Write([]byte(domain.Config.SecretKey))
	d.Write([]byte("zFSa4a0529q_rSzn9uly78CYvVCNQSHcViXFGul_oZ0"))
	d.Write([]byte(nonce))
	d.Write([]byte("#"))
	d.Write([]byte(path))
	d.Read(h)
	return base64url.Encode(h)
}

func splitHash(input string) string {
	if len(input) < 4 {
		return input
	}
	output := input[:1] + "/" + input[1:3] + "/" + input[3:]
	return output
}
