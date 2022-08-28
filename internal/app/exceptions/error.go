package exceptions

import "errors"

var UrlAlreadyExist = errors.New("url already exists")
var HashAlreadyExist = errors.New("hash already exists")
