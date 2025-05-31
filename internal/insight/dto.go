package insight

import "github.com/google/uuid"

type Insight struct {
	Id   uuid.UUID
	Text string
}
