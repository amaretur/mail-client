package eds

import (
	"bytes"
	"crypto/dsa"
	"encoding/asn1"
	"encoding/pem"
	"errors"
)

func ExportDSAPrivateKeyAsPemBytes(privKey *dsa.PrivateKey) ([]byte, error) {

	if privKey == nil {
		return nil, errors.New("nil pointer dereference")
	}

	buff := bytes.NewBuffer([]byte{})

	asn1Bytes, err := asn1.Marshal(*privKey)
	if err != nil {
		return nil, err
	}

	err = pem.Encode(
		buff,
		&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: asn1Bytes,
		},
	)

	return buff.Bytes(), err
}

func ParseDSAPrivateKeyFromPemBytes(privPem []byte) (*dsa.PrivateKey, error) {
	block, _ := pem.Decode(privPem)

	key := new(dsa.PrivateKey)
	_, err := asn1.Unmarshal(block.Bytes, key)
    
	return key, err
}

func ExportDSAPublicKeyAsPemBytes(pubKey *dsa.PublicKey) ([]byte, error) {

	if pubKey == nil {
		return nil, errors.New("nil pointer dereference")
	}

	buff := bytes.NewBuffer([]byte{})

	asn1Bytes, err := asn1.Marshal(*pubKey)
	if err != nil {
		return nil, err
	}

	err = pem.Encode(
		buff,
		&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: asn1Bytes,
		},
	)

	return buff.Bytes(), err
}

func ParseDSAPublicKeyFromPemBytes(pubPem []byte) (*dsa.PublicKey, error) {

	block, _ := pem.Decode(pubPem)
	if block == nil {
		return nil, errors.New("error parse key")
	}

	key := new(dsa.PublicKey)
	_, err := asn1.Unmarshal(block.Bytes, key)
    
	return key, err
}


