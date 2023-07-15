package scipher

import (
	"errors"
)

var (
	ErrInvalidBuffSize	= errors.New(
		"buffer size must be a multiple of the block size",
	)
	ErrInvalidIVSize	= errors.New("invalid size of init vector")
	ErrInvalidBlockSize	= errors.New("invalid size of block")
)
