package cipher_test

import (
	"crypto/rsa"
	"bytes"
	"io"
	"testing"

	"github.com/amaretur/spg/internal/cipher"
)

func TestEcnrypt(t* testing.T) {

	pubKey, _, _ := cipher.RSAKeyGen()

	cases := []struct{
		pubKey	*rsa.PublicKey
		origin	io.Reader
		encrypt	io.Writer
		ok		bool
	}{
		{
			pubKey: pubKey,
			origin: bytes.NewBuffer([]byte("Hello world!")),
			encrypt: bytes.NewBuffer([]byte{}),
			ok: true,
		},		
		{
			pubKey: pubKey,
			origin: bytes.NewBuffer([]byte("")),
			encrypt: bytes.NewBuffer([]byte{}),
			ok: true,
		},
		{
			pubKey: nil,
			origin: bytes.NewBuffer([]byte("Hello world!")),
			encrypt: bytes.NewBuffer([]byte{}),
			ok: false,
		},
		{
			pubKey: new(rsa.PublicKey),
			origin: bytes.NewBuffer([]byte("")),
			encrypt: bytes.NewBuffer([]byte{}),
			ok: false,
		},
	}

	for i, test := range cases {

		t.Logf("Test case #%d", i+1)

		err := cipher.Encrypt(test.pubKey, test.origin, test.encrypt)
		if (err != nil && test.ok) || (err == nil && !test.ok)  {
			t.Errorf("Incorrect recult! err: %v", err)
		}
	}
}

func TestDecrypt(t* testing.T) {

	pubKey, privKey, _ := cipher.RSAKeyGen()

	encrypt := bytes.NewBuffer([]byte{})

	cipher.Encrypt(pubKey, bytes.NewBuffer([]byte("Hello world!")), encrypt)

	cases := []struct{
		privKey	*rsa.PrivateKey
		encrypt	io.Reader
		decrypt	io.Writer
		ok		bool
	}{
		{
			privKey: privKey,
			encrypt: encrypt,
			decrypt: bytes.NewBuffer([]byte{}),
			ok: true,
		},		
		{
			privKey: privKey,
			encrypt: bytes.NewBuffer([]byte("")),
			decrypt: bytes.NewBuffer([]byte{}),
			ok: false,
		},
		{
			privKey: privKey,
			encrypt: bytes.NewBuffer([]byte("Hello world!")),
			decrypt: bytes.NewBuffer([]byte{}),
			ok: false,
		},
		{
			privKey: nil,
			encrypt: encrypt,
			decrypt: bytes.NewBuffer([]byte{}),
			ok: false,
		},
		{
			privKey: new(rsa.PrivateKey),
			encrypt: bytes.NewBuffer([]byte("")),
			decrypt: bytes.NewBuffer([]byte{}),
			ok: false,
		},
	}

	for i, test := range cases {

		t.Logf("Test case #%d", i+1)

		err := cipher.Decrypt(test.privKey, test.encrypt, test.decrypt)
		if (err != nil && test.ok) || (err == nil && !test.ok)  {
			t.Errorf("Incorrect recult! err: %v", err)
		}
	}
}

func TestEncryptDecrypt(t *testing.T) {

	cases := []string{
		"Hello world!",
		"",
		"Hi!",
		"Привет мир!",
		"Input data 12345",
		`Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed 
		do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut 
		enim ad minim veniam, quis nostrud exercitation ullamco laboris 
		nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in 
		reprehenderit in voluptate velit esse cillum dolore eu fugiat 
		nulla pariatur. Excepteur sint occaecat cupidatat non proident, 
		sunt in culpa qui officia deserunt mollit anim id est laborum.`,
		`В своём стремлении улучшить пользовательский опыт мы 
		упускаем, что многие известные личности могут быть объединены в 
		целые кластеры себе подобных. И нет сомнений, что стремящиеся 
		вытеснить традиционное производство, нанотехнологии, вне 
		зависимости от их уровня, должны быть заблокированы в рамках 
		своих собственных рациональных ограничений. В целом, конечно, 
		реализация намеченных плановых заданий играет важную роль в 
		формировании укрепления моральных ценностей.`,
		`Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!Hello!!!`,
	}

	for i, data := range cases {

		t.Logf("Test case #%d", i+1)

		pubKey, privKey, _ := cipher.RSAKeyGen()

		origin := bytes.NewBuffer([]byte(data))
		encrypt := bytes.NewBuffer([]byte{})
		decrypt := bytes.NewBuffer([]byte{})

		cipher.Encrypt(pubKey, origin, encrypt)

		cipher.Decrypt(privKey, encrypt, decrypt)

		if  data != decrypt.String() {
			t.Errorf(
				"Incorrect result! \ndata: %s\ndecrypt: %s", 
				data, decrypt,
			)
		}
	}
}
