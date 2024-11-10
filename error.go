package admob

import "errors"

var (
	ErrLoadPublicKey        = errors.New("load Google Public Key Failed")
	ErrCannotFoundPublicKey = errors.New("cannot found public key by key id")
	ErrInvalidSignature     = errors.New("invalid signature")
)
