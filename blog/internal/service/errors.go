package service

import "errors"

var (
	ErrorLinkTooShort        = errors.New("error link too short")
	ErrorInvalidSymbolInLink = errors.New("error invalid symbol in link")
	ErrorLinkAlreadyExists   = errors.New("error link already exists")
)
