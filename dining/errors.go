package dining

import "github.com/pkg/errors"

var (
	ErrInvalidUUID       = errors.New("invalid uuid")
	ErrUnknownTable      = errors.New("table does not exist")
	ErrTableNotAvaliable = errors.New("table is not available")
)
