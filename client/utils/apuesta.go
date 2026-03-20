package utils

type Fecha struct {
	Anio uint16
	Mes  uint8
	Dia  uint8
}

// Apuesta represents a bet made by a client. It contains the
// client's personal information and the bet number.
type Apuesta struct {
	Nombre     string
	Apellido   string
	Documento  uint32
	Nacimiento Fecha
	Numero     uint16
}
