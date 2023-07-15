package scipher

import (
	"crypto/cipher"
)

type CBCMode struct {
	iv	[]byte
}

func NewCBCMode(iv []byte) *CBCMode {
	return &CBCMode{
		iv: iv,
	}
}

func (c *CBCMode) GetEncrypter(b cipher.Block) (cipher.BlockMode, error) {
	if len(c.iv) != b.BlockSize() {
		return nil, ErrInvalidIVSize
	}

	return cipher.NewCBCEncrypter(b, c.iv), nil
}  
func (c *CBCMode) GetDecrypter(b cipher.Block) (cipher.BlockMode, error) {
	if len(c.iv) != b.BlockSize() {
		return nil, ErrInvalidIVSize
	}
	return cipher.NewCBCDecrypter(b, c.iv), nil
}
