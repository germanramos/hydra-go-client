package error

import (
	"errors"
)

var (
	IncorrectHydraServerResponse error = errors.New("Incorrect Hydra server response")
	InaccessibleHydraServer      error = errors.New("Inaccessible Hydra server")
)
