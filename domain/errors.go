package domain

import "errors"

var (
	ErrMissingArgument    = errors.New("missing argument")
	ErrInvalidMessageType = errors.New("invalid message-type")
	ErrInvalidKey         = errors.New("invalid key")
	ErrInvalidGalleryName = errors.New("invalid gallery name")
	ErrInvalidFileName    = errors.New("invalid file name")
	ErrUrlProtocol        = errors.New("url protocol must be http or https")
	ErrUrlUnsupportedExt  = errors.New("url points to something that is not a media file")
	ErrNotFound           = errors.New("not found")
)
