package cipher 

import (
	"io"
	"crypto/rand"
	"crypto/des"
	"crypto/sha256" 
	"crypto/rsa"
	"errors"

	"github.com/amaretur/mail-client/pkg/cipher/scipher"
)

const (
	SKeySize	= 24 
	EncSKeySize	= 256
	BufferSize	= 1024
	RSAKeySize	= 2048
)

func Encrypt(pubKey *rsa.PublicKey, r io.Reader, w io.Writer) error {

	if pubKey == nil {
		return errors.New("nil pointer dereference")
	}

	sKey := scipher.SKeyGen(SKeySize)
	block, err := des.NewTripleDESCipher(sKey)
	if err != nil {
		return err
	}

	iv := sKey[:block.BlockSize()]

	// encrypt session key
	encSKey, err := rsa.EncryptOAEP(
		sha256.New(), rand.Reader, pubKey, sKey, []byte{},
	)
	if err != nil {
		return err
	}
	w.Write(encSKey)	// writing an encrypted key for a symmetric algorithm

	mode := scipher.NewCBCMode(iv)
	cp, err := scipher.New(block, mode, BufferSize)
	if err != nil {
		return err
	}

	return cp.Encrypt(r, w)
}

func Decrypt(privKey *rsa.PrivateKey, r io.Reader, w io.Writer) error {

	if privKey == nil {
		return errors.New("nil pointer dereference")
	}

	encSKey := make([]byte, EncSKeySize)

	// read encrypt key
	if _, err := r.Read(encSKey); err != nil {
		return err
	}

	// decrypt session key
	sKey, err := rsa.DecryptOAEP(
		sha256.New(), rand.Reader, privKey, encSKey, []byte{},
	)
	if err != nil {
		return err
	}

	block, err := des.NewTripleDESCipher(sKey)
	if err != nil {
		return err
	}

	iv := sKey[:block.BlockSize()]

	cp, err := scipher.New(block, scipher.NewCBCMode(iv), BufferSize)
	if err != nil {
		return err
	}

	return cp.Decrypt(r, w)
}

func RSAKeyGen() (*rsa.PublicKey, *rsa.PrivateKey, error) {

	privKey, err := rsa.GenerateKey(rand.Reader, RSAKeySize)  
	if err != nil {  
		return nil, nil, err
	}  

	return &privKey.PublicKey, privKey, nil
}
