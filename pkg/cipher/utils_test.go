package cipher_test 

import (
	"testing"
	"github.com/stretchr/testify/require"

	"crypto/rsa"

	"github.com/amaretur/spg/internal/cipher"
)

func TestRSAExportPublicKeyErr(t *testing.T) {

	cases := []*rsa.PublicKey{
		nil, 
		new(rsa.PublicKey),
	}

	for i, key := range cases {
		t.Logf("TestRSA case #%d", i+1)

		_, err := cipher.ExportRSAPublicKeyAsPemBytes(key)
		require.NotNil(t, err)
	}
}

func TestRSAExportPrivateKeyErr(t *testing.T) {

	cases := []*rsa.PrivateKey{
		nil, 
	}

	for i, key := range cases {
		t.Logf("TestRSA case #%d", i+1)

		_, err := cipher.ExportRSAPrivateKeyAsPemBytes(key)
		require.NotNil(t, err)
	}
}

func TestRSAParsePublicKeyErr(t *testing.T) {

	_, err := cipher.ParseRSAPublicKeyFromPemBytes([]byte{})
	require.NotNil(t, err)
}

func TestRSAParsePrivateKeyErr(t *testing.T) {

	_, err := cipher.ParseRSAPrivateKeyFromPemBytes([]byte{})
	require.NotNil(t, err)
}

func TestRSAExportParsePublicKey(t *testing.T) {

	// key generation
	pubKey, _, err := cipher.RSAKeyGen()
	if err != nil {
		t.Errorf("Incorrect result! err: %v", err)

	}

	// export to pem bytes 
	pubKeyBytes, err := cipher.ExportRSAPublicKeyAsPemBytes(pubKey)
	if err != nil {
		t.Errorf("Incorrect result! err: %v", err)

	}

	// parse from pem bytes 
	pubKey1, err := cipher.ParseRSAPublicKeyFromPemBytes(pubKeyBytes)
	if err != nil {
		t.Errorf("Incorrect result! err: %v", err)

	}

	// export to pem bytes
	pubKeyBytes1, err := cipher.ExportRSAPublicKeyAsPemBytes(pubKey1)
	if err != nil {
		t.Errorf("Incorrect result! err: %v", err)
	}

	// validation
	if string(pubKeyBytes) != string(pubKeyBytes1) {
		t.Errorf("Incorrect result! err: %v", err)
	}
}

func TestRSAExportParsePrivateKey(t *testing.T) {

	// key generation
	_, privKey, err := cipher.RSAKeyGen()
	if err != nil {
		t.Errorf("Incorrect result! err: %v", err)
	}

	// export to pem bytes 
	privKeyBytes, err := cipher.ExportRSAPrivateKeyAsPemBytes(privKey)
	if err != nil {
		t.Errorf("Incorrect result! err: %v", err)
	}

	// parse from pem bytes 
	privKey1, err := cipher.ParseRSAPrivateKeyFromPemBytes(privKeyBytes)
	if err != nil {
		t.Errorf("Incorrect result! err: %v", err)
	}

	// export to pem bytes
	privKeyBytes1, err := cipher.ExportRSAPrivateKeyAsPemBytes(privKey1)
	if err != nil {
		t.Errorf("Incorrect result! err: %v", err)
	}

	// validation
	if string(privKeyBytes) != string(privKeyBytes1) {
		t.Errorf("Incorrect result! err: %v", err)
	}
}

