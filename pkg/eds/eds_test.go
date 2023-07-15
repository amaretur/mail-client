package eds_test

import (
	"testing"
	"crypto/dsa"

	"github.com/amaretur/spg/internal/eds"
)

func TestSign(t *testing.T) {

	_, privKey, _ := eds.DSAKeyGen()

	cases := []struct{
		msg		[]byte
		privKey	*dsa.PrivateKey
		ok		bool
	} {
		{
			msg: []byte(""),
			privKey: privKey,
			ok: true,
		},
		{
			msg: []byte("Hello world"),
			privKey: nil,
			ok: false, 
		},
		{
			msg: nil,
			privKey: privKey,
			ok: false,
		},
	}

	for i, test := range cases {

		t.Logf("Test case #%d", i+1)
		_, err := eds.Sign(test.privKey, test.msg)
		if (err != nil && test.ok) || (err == nil && !test.ok) {
			t.Errorf("Incorrect result! err: %v", err)
		}
	}
}

func TestVarify(t *testing.T) {

	pubKey, privKey, _ := eds.DSAKeyGen()
	msg := []byte("Hello world!")
	sign, _ := eds.Sign(privKey, msg)

	cases := []struct{
		msg		[]byte
		pubKey	*dsa.PublicKey
		sign	[]byte
		ok		bool
	} {
		{
			msg: []byte(""),
			pubKey: pubKey,
			sign: sign,
			ok: true,
		},
		{
			msg: []byte("Hello world"),
			pubKey: nil,
			sign: sign,
			ok: false, 
		},
		{
			msg: nil,
			pubKey: pubKey,
			sign: sign,
			ok: false,
		},
		{
			msg: []byte("Hello world"),
			pubKey: pubKey,
			sign: nil,
			ok: false,
		},
		{
			msg: []byte("Hello world"),
			pubKey: pubKey,
			sign: make([]byte, eds.SignSize),
			ok: true,
		},		
		{
			msg: []byte("Hello world"),
			pubKey: pubKey,
			sign: make([]byte, 12),
			ok: false,
		},
	}

	for i, test := range cases {

		t.Logf("Test case #%d", i+1)
		_, err := eds.Verify(test.pubKey, test.sign, test.msg)
		if (err != nil && test.ok) || (err == nil && !test.ok) {
			t.Errorf("Incorrect result! err: %v", err)
		}
	}
}

func TestSignVerify(t* testing.T) {

	pubKey, privKey, _ := eds.DSAKeyGen()
	_, privKey2, _ := eds.DSAKeyGen()

	cases := []struct{
		msg		[]byte
		pubKey	*dsa.PublicKey
		privKey	*dsa.PrivateKey
		res		bool
	}{
		{
			msg: []byte(""),
			pubKey: pubKey,
			privKey: privKey,
			res: true,
		},
		{
			msg: []byte("Hello world! Hello world! Hello world! Hello world!"),
			pubKey: pubKey,
			privKey: privKey,
			res: true,
		},
		{
			msg: []byte("Привет мир! Привет мир! Привет мир! Привет мир!"),
			pubKey: pubKey,
			privKey: privKey,
			res: true,
		}, 
		{
			msg: []byte("Hello world! Привет мир! 1234567890=-!№;%:?*()+-*/"),
			pubKey: pubKey,
			privKey: privKey,
			res: true,
		},
		{
			msg: []byte("Hello world!"),
			pubKey: pubKey,
			privKey: privKey2,
			res: false,
		},
	}

	for i, test := range cases {

		t.Logf("Test case #%d", i+1)

		sign, _ := eds.Sign(test.privKey, test.msg)
		ok, err := eds.Verify(test.pubKey, sign, test.msg)

		if ok != test.res || err != nil {
			t.Errorf("Incorrect result! ok: %v, err: %v", ok, err)
		}
	}
}
