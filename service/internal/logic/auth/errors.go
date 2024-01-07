package auth

import "errors"

var ErrTokenParse = errors.New("unable to parse token into claims")
