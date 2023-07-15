package scipher

import (
	"io"
	"bytes"
	"crypto/cipher"
)

type Mode interface {
	GetEncrypter(b cipher.Block) (cipher.BlockMode, error)
	GetDecrypter(b cipher.Block) (cipher.BlockMode, error)
}

type Scipher struct {
	block		cipher.Block
	mode		Mode
	buffSize	int
}

func New(block cipher.Block, mode Mode, buffSize int) (*Scipher, error) {

	if buffSize%block.BlockSize() != 0 || buffSize <= 0 {
		return nil, ErrInvalidBuffSize
	}

	return &Scipher {
		block: block,
		mode: mode,
		buffSize: buffSize,
	}, nil
}

func (s* Scipher) PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func (s* Scipher) PKCS7Unpadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])

	if length < unpadding {
		return origData
	}

	return origData[:(length - unpadding)]
}

func (s* Scipher) Encrypt(r io.Reader, w io.Writer) error {

	quit := false
	blockSize := s.block.BlockSize()

	buff := make([]byte, s.buffSize)

	encrypter, err := s.mode.GetEncrypter(s.block)
	if err != nil {
		return err
	}

	for n, err := r.Read(buff); !quit; n, err = r.Read(buff) {

		if err != nil && err != io.EOF {
			return err
		}

		if n < s.buffSize {
			buff = s.PKCS7Padding(buff[:n], blockSize)
			quit = true
		}

		out := make([]byte, len(buff))
		encrypter.CryptBlocks(out, buff)

		if _, err = w.Write(out); err != nil {
			return err
		}
	}

	return nil
}

func (s* Scipher) Decrypt(r io.Reader, w io.Writer) error {

	quit := false
	buff := make([]byte, s.buffSize)

	var old []byte

	blockSize := s.block.BlockSize()

	decrypter, err := s.mode.GetDecrypter(s.block)
	if err != nil {
		return err
	}

	for n, err := r.Read(buff); !quit; n, err = r.Read(buff) {

		if err != nil && err != io.EOF {
			return err
		}

		buff = buff[:n]
		if len(buff)%blockSize != 0 {
			return ErrInvalidBlockSize
		}

		out := make([]byte, n)
		decrypter.CryptBlocks(out, buff)

		if n < s.buffSize {

			if n == 0 {

				// if read only 0 bytes
				if old == nil {
					return nil
				}

				out = old
				old = nil
			}

			out = s.PKCS7Unpadding(out)
			quit = true
		}

		if old != nil {
			if _, err = w.Write(old); err != nil {
				return err
			}
		}

		old = out
	}

	_, err = w.Write(old)
	return err
}


