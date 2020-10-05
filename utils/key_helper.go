package utils

import (
	"encoding/pem"
	"strings"
	ssh "golang.org/x/crypto/ssh"
	rsa "crypto/rsa"
)

func GeneratePublicKey(privateKey []byte) (string, error) {
	rsaPK, err := ssh.ParseRawPrivateKey(privateKey)
	if err != nil {
		return "", err
	}

	sshPK := rsaPK.(*rsa.PrivateKey)
	sshPubKey, err := ssh.NewPublicKey(&sshPK.PublicKey)
	if err != nil {
		return "", err
	}

	return MarshalPublickKey(sshPubKey), nil
}

func MarshalPublickKey(key ssh.PublicKey) string {
	pemBlock := &pem.Block{
		Type: "", // The type, taken from the preamble (i.e. "RSA PRIVATE KEY").
		// Headers: map[string]string{}, // Optional headers.
		Bytes: key.Marshal(), // The decoded bytes of the contents. Typically a DER encoded ASN.1 structure.
	}
	pemData := pem.EncodeToMemory(pemBlock)
	pemArray := strings.Split(string(pemData), "\n")

	return strings.Join(pemArray[1:len(pemArray)-2], "")
}
