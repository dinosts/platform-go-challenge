package audience

import "github.com/google/uuid"

type Audience struct {
	Id                 uuid.UUID
	Gender             string
	BirthCountry       string
	AgeGroup           string
	SocialMediaHours   float64
	PurchasesLastMonth int
}
