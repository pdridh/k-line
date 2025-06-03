package dining

import "github.com/pkg/errors"

var (
	ErrTableNoOpenSession = errors.New("table does not have an open session")
)
