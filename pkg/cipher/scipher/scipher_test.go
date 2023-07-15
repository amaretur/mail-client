package scipher_test

import (
	"io"
	"testing"
	"crypto/des"
	"bytes"
	"errors"

	"github.com/amaretur/spg/pkg/scipher"
)

// valid unput data for use in test
const (
	validIv		= "12345678"
	validKey	= "123456789123456789123456"
)

// error value for test
var (
	ErrIO = errors.New("i/o error")
)

// a structure that simulates incorrect data input/output streams
type InvalidIO struct {
}

func (i InvalidIO) Read(dst []byte) (int, error) {
	return 0, ErrIO
}

func (i InvalidIO) Write(data []byte) (int, error) {
	return 0, ErrIO
}

// tests
func TestScipherBuffSize(t *testing.T) {

	cases := []struct{
		buffSize int
		err error
	}{
		{
			buffSize: -8,
			err: scipher.ErrInvalidBuffSize,
		},
		{
			buffSize: 0,
			err: scipher.ErrInvalidBuffSize,
		},
		{
			buffSize: 5,
			err: scipher.ErrInvalidBuffSize,
		},
		{
			buffSize: 12,
			err: scipher.ErrInvalidBuffSize,
		},
		{
			buffSize: 8,
			err: nil,
		},
		{
			buffSize: 1024,
			err: nil,
		},
	}

	for _, test := range cases {

		t.Logf("Test case: buffSize=%d, err=%v", test.buffSize, test.err)

		mode := scipher.NewCBCMode([]byte(validIv))
		block, _ := des.NewTripleDESCipher([]byte(validKey))

		_, err := scipher.New(block, mode, test.buffSize)

		if err != test.err {
			t.Errorf("Incorrect result! err: %v", err)
		}
	}
}

func TestEncrypt(t *testing.T) {

	validIO := bytes.NewBuffer([]byte{})
	invalidIO := InvalidIO{}

	cases := []struct{
		iv	[]byte
		r	io.Reader
		w	io.Writer
		err	error
	}{
		{
			iv: []byte(validIv),
			r: validIO,
			w: validIO,
			err: nil, 
		},
		{
			iv: []byte("1234"),
			r: validIO,
			w: validIO,
			err: scipher.ErrInvalidIVSize,
		},
		{
			iv: []byte("12345678910"),
			r: validIO,
			w: validIO,
			err: scipher.ErrInvalidIVSize,
		},
		{
			iv: []byte(validIv),
			r: invalidIO,
			w: validIO,
			err: ErrIO,
		},
		{
			iv: []byte(validIv),
			r: validIO,
			w: invalidIO,
			err: ErrIO,
		},
	}

	for i, test := range cases {

		mode := scipher.NewCBCMode([]byte(test.iv))
		block, _ := des.NewTripleDESCipher([]byte(validKey))
		cp, _ := scipher.New(block, mode, 16)

		t.Logf("Test case #%d", i+1)

		err := cp.Encrypt(test.r, test.w)
		if err != test.err {
			t.Errorf("Incorrect result! err: %v", err)
		}
	}
}

func TestDecrypt(t *testing.T) {

	validIO := bytes.NewBuffer([]byte{})
	invalidIO := InvalidIO{}

	cases := []struct{
		iv	[]byte
		r	io.Reader
		w	io.Writer
		err	error
	}{
		{
			iv: []byte(validIv),
			r: validIO,
			w: validIO,
			err: nil, 
		},
		{
			iv: []byte("1234"),
			r: validIO,
			w: validIO,
			err: scipher.ErrInvalidIVSize,
		},
		{
			iv: []byte("12345678910"),
			r: validIO,
			w: validIO,
			err: scipher.ErrInvalidIVSize,
		},
		{
			iv: []byte(validIv),
			r: invalidIO,
			w: validIO,
			err: ErrIO,
		},
		{
			iv: []byte(validIv),
			r: bytes.NewBuffer([]byte("00000000")),
			w: invalidIO,
			err: ErrIO,
		},
		{
			iv: []byte(validIv),
			r: bytes.NewBuffer([]byte("000")),
			w: validIO,
			err: scipher.ErrInvalidBlockSize,
		},
	}

	for i, test := range cases {

		mode := scipher.NewCBCMode([]byte(test.iv))
		block, _ := des.NewTripleDESCipher([]byte(validKey))
		cp, _ := scipher.New(block, mode, 16)

		t.Logf("Test case #%d", i+1)

		err := cp.Decrypt(test.r, test.w)
		if err != test.err {
			t.Errorf("Incorrect result! err: %v", err)
		}
	}
}

func TestEncryptAndDectypr(t *testing.T) {

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
		`LLorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehende`,
	}

	for i, text := range cases {

		mode := scipher.NewCBCMode([]byte(validIv))
		block, _ := des.NewTripleDESCipher([]byte(validKey))
		cp, _ := scipher.New(block, mode, 16)

		in := bytes.NewBuffer([]byte(text))
		encrypted := bytes.NewBuffer([]byte{})
		decrypted := bytes.NewBuffer([]byte{})

		cp.Encrypt(in, encrypted)
		cp.Decrypt(encrypted, decrypted)

		t.Logf("Test case #%d", i+1)

		if text != decrypted.String() {
			t.Errorf("Invalid result! decrypted: %v | %v", decrypted.Bytes(), []byte(text))
		}
	}
}

