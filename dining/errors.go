package dining

import "github.com/pkg/errors"

var (
	ErrInvalidUUID       = errors.New("invalid uuid")
	ErrUnknownTable      = errors.New("table does not exist")
	ErrTableNotAvaliable = errors.New("table is not available")
	ErrUnknownOrder      = errors.New("order does not exist")
	ErrOrderNotOngoing   = errors.New("order is not ongoing")
	ErrUnkownOrderItem   = errors.New("order item does not exist")
)
