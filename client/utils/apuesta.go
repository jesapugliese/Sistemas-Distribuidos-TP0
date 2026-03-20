package utils

// Apuesta represents a bet made by a client. It contains the
// client's personal information and the bet number.
type Apuesta struct {
	Nombre     string
	Apellido   string
	Documento  int
	Nacimiento string
	Numero     int
}
