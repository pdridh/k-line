package dining

import "github.com/pkg/errors"

var (
	ErrInvalidUUID        = errors.New("invalid uuid")
	ErrTableNoOpenSession = errors.New("table does not have an open session")
)
