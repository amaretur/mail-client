package cipher

import (
    "crypto/rsa"
    "crypto/x509"
	"encoding/pem"
    "errors"
)

func ExportRSAPrivateKeyAsPemBytes(privKey *rsa.PrivateKey) ([]byte, error) {

	if privKey == nil {
		return nil, errors.New("nil pointer dereference")
	}

	return pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privKey),
		},
	), nil
}

func ParseRSAPrivateKeyFromPemBytes(privPem []byte) (*rsa.PrivateKey, error) {

	block, _ := pem.Decode(privPem)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func ExportRSAPublicKeyAsPemBytes(pubKey *rsa.PublicKey) ([]byte, error) {

	if pubKey == nil {
		return nil, errors.New("nil pointer dereference")
	}
	pubKeyPem, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		return nil, err
	}

	return pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: pubKeyPem,
		},
	), nil
}

func ParseRSAPublicKeyFromPemBytes(pubPem []byte) (*rsa.PublicKey, error) {

	block, _ := pem.Decode(pubPem)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch pub := pub.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		return nil, errors.New("Key type is not RSA")
	}
}
