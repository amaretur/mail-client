package eds_test 

import (
	"testing"
	"crypto/dsa"

	"github.com/stretchr/testify/require"

	"github.com/amaretur/spg/internal/eds"
)

func TestDSAExportPublicKeyErr(t *testing.T) {

	cases := []*dsa.PublicKey{
		nil, 
		new(dsa.PublicKey),
	}

	for i, key := range cases {
		t.Logf("Test case #%d", i+1)

		_, err := eds.ExportDSAPublicKeyAsPemBytes(key)
		require.NotNil(t, err)
	}
}

func TestDSAExportPrivateKeyErr(t *testing.T) {

	cases := []*dsa.PrivateKey{
		nil, 
		new(dsa.PrivateKey),
	}

	for i, key := range cases {
		t.Logf("Test case #%d", i+1)

		_, err := eds.ExportDSAPrivateKeyAsPemBytes(key)
		require.NotNil(t, err)
	}
}

func TestDSAExportParsePublicKey(t *testing.T) {

	// key generation
	pubKey, _, err := eds.DSAKeyGen()
	if err != nil {
		t.Errorf("Incorrect result! err: %v", err)

	}

	// export to pem bytes 
	pubKeyBytes, err := eds.ExportDSAPublicKeyAsPemBytes(pubKey)
	if err != nil {
		t.Errorf("Incorrect result! err: %v", err)

	}

	// parse from pem bytes 
	pubKey1, err := eds.ParseDSAPublicKeyFromPemBytes(pubKeyBytes)
	if err != nil {
		t.Errorf("Incorrect result! err: %v", err)

	}

	// export to pem bytes
	pubKeyBytes1, err := eds.ExportDSAPublicKeyAsPemBytes(pubKey1)
	if err != nil {
		t.Errorf("Incorrect result! err: %v", err)
	}

	// validation
	if string(pubKeyBytes) != string(pubKeyBytes1) {
		t.Errorf("Incorrect result! err: %v", err)
	}
}

func TestDSAExportParsePrivateKey(t *testing.T) {

	// key generation
	_, privKey, err := eds.DSAKeyGen()
	if err != nil {
		t.Errorf("Incorrect result! err: %v", err)

	}

	// export to pem bytes 
	privKeyBytes, err := eds.ExportDSAPrivateKeyAsPemBytes(privKey)
	if err != nil {
		t.Errorf("Incorrect result! err: %v", err)

	}

	// parse from pem bytes 
	privKey1, err := eds.ParseDSAPrivateKeyFromPemBytes(privKeyBytes)
	if err != nil {
		t.Errorf("Incorrect result! err: %v", err)

	}

	// export to pem bytes
	privKeyBytes1, err := eds.ExportDSAPrivateKeyAsPemBytes(privKey1)
	if err != nil {
		t.Errorf("Incorrect result! err: %v", err)
	}

	// validation
	if string(privKeyBytes) != string(privKeyBytes1) {
		t.Errorf("Incorrect result! err: %v", err)
	}
}
