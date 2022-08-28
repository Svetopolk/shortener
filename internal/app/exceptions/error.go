package exceptions

import "errors"

var ErrURLAlreadyExist = errors.New("url already exists")
var ErrHashAlreadyExist = errors.New("hash already exists")
