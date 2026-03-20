package communication

import (
	"encoding/binary"
	"fmt"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/utils"
)

type ApuestaSerializer struct {
}

func NewApuestaSerializer() ApuestaSerializer {
	return ApuestaSerializer{}
}

// Serialize converts an Apuesta struct into a byte slice that can be sent
// over the network.
// Serialization format:
//   - Nombre: variable length string (preceded by its 1-byte length)
//   - Apellido: variable length string (preceded by its 1-byte length)
//   - Documento: 4 bytes (integer)
//   - Nacimiento:
//     > Año: 2 bytes (integer)
//     > Mes: 1 byte (integer)
//     > Dia: 1 byte (integer)
//   - Numero: 2 bytes (integer)
func (as ApuestaSerializer) Serialize(apuesta utils.Apuesta) ([]byte, error) {
	var serialized []byte

	// Serializar Nombre
	nombreBytes := []byte(apuesta.Nombre)
	if len(nombreBytes) > 255 {
		return nil, fmt.Errorf("Nombre demasiado largo")
	}
	serialized = append(serialized, byte(len(nombreBytes)))
	serialized = append(serialized, nombreBytes...)

	// Serializar Apellido
	apellidoBytes := []byte(apuesta.Apellido)
	if len(apellidoBytes) > 255 {
		return nil, fmt.Errorf("Apellido demasiado largo")
	}
	serialized = append(serialized, byte(len(apellidoBytes)))
	serialized = append(serialized, apellidoBytes...)

	// Serializar Documento
	var tmp4 [4]byte
	binary.BigEndian.PutUint32(tmp4[:], apuesta.Documento)
	serialized = append(serialized, tmp4[:]...)

	// Serializar Nacimiento
	binary.BigEndian.PutUint16(tmp4[0:2], apuesta.Nacimiento.Anio)
	tmp4[2] = apuesta.Nacimiento.Mes
	tmp4[3] = apuesta.Nacimiento.Dia
	serialized = append(serialized, tmp4[:]...)

	// Serializar Numero
	var tmp2 [2]byte
	binary.BigEndian.PutUint16(tmp2[:], uint16(apuesta.Numero))
	serialized = append(serialized, tmp2[:]...)

	return serialized, nil
}
