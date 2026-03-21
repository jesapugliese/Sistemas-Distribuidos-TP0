package utils

type Date struct {
	Year  uint16
	Month uint8
	Day   uint8
}

// Bet represents a bet made by a client. It contains the
// client's personal information and the bet number.
type Bet struct {
	AgencyID  uint8
	FirstName string
	LastName  string
	Document  uint32
	Birthdate Date
	Number    uint16
}
