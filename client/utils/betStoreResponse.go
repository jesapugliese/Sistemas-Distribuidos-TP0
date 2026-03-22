package utils

// BetStoreResponse represents the response received from the server
// after sending a request for storage of a bet.
type BetStoreResponse struct {
	Success  bool
	Document uint32
	Number   uint16
}
