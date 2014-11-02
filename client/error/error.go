package error

type IncorrectHydraServerResponse struct {
	s string
}

func (i IncorrectHydraServerResponse) Error() string {
	return i.s
}

type InaccessibleHydraServer struct {
	s string
}

func (i InaccessibleHydraServer) Error() string {
	return i.s
}

type HydraNotAvailable struct {
	s string
}

func (h HydraNotAvailable) Error() string {
	return h.s
}

var (
	IncorrectHydraServerResponseError IncorrectHydraServerResponse = IncorrectHydraServerResponse{"Incorrect Hydra server response"}
	InaccessibleHydraServerError      InaccessibleHydraServer      = InaccessibleHydraServer{"Inaccessible Hydra server"}
	HydraNotAvailableError            HydraNotAvailable            = HydraNotAvailable{"Hydra is not available"}
)
