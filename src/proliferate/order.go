package proliferate

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
)

// Orderer struct
type Orderer struct{}


//NewID generates UUID V4 ID
func NewID() string {
	id := uuid.Must(uuid.NewV4())
	return fmt.Sprintf("%s", id)
}
