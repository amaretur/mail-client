package eds

import (
	"math/big"
	"crypto/dsa"
	"crypto/rand"
	"crypto/md5"
	"errors"
)

const (
	SignSize = 40
)

func hash(data []byte) ([]byte, error) {

	if data == nil {
		return nil, errors.New("nil pointer dereference")
	}

	h := md5.New()
	h.Write(data)

	return h.Sum(nil), nil
}

func Sign(privKey *dsa.PrivateKey, data []byte) ([]byte, error) {

	if privKey == nil {
		return nil, errors.New("nil pointer dereference")
	}

	summ, err := hash(data)
	if err != nil {
		return nil, err
	}

	r, s, err := dsa.Sign(rand.Reader, privKey, summ)
	if err != nil {
		return nil, err
	}

	sign := r.Bytes()
	sign = append(sign, s.Bytes()...)

	return sign, nil
}

func Verify(pubKey *dsa.PublicKey, sign, data []byte) (bool, error) {

	if pubKey == nil || sign == nil {
		return false, errors.New("nil pointer dereference")
	}

	if len(sign) != SignSize {
		return false, errors.New("invalid sign")
	}

	r := new(big.Int).SetBytes(sign[0:SignSize/2])
	s := new(big.Int).SetBytes(sign[SignSize/2:SignSize])

	summ, err := hash(data)
	if err != nil {
		return false, err
	}

	return dsa.Verify(pubKey, summ, r, s), nil
}

func DSAKeyGen() (*dsa.PublicKey, *dsa.PrivateKey, error) {

	params := new(dsa.Parameters)

	err := dsa.GenerateParameters(params, rand.Reader, dsa.L1024N160)
	if err != nil {
		return nil, nil, err
	}

	privKey := new(dsa.PrivateKey)
	privKey.PublicKey.Parameters = *params

	dsa.GenerateKey(privKey, rand.Reader)

	return &privKey.PublicKey, privKey, nil
}
