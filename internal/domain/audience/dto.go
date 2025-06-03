package audience

import "github.com/google/uuid"

type Audience struct {
	Id                 uuid.UUID `json:"id"`
	Gender             string    `json:"gender"`
	BirthCountry       string    `json:"birth_country"`
	AgeGroup           string    `json:"age_group"`
	SocialMediaHours   float64   `json:"social_media_hours"`
	PurchasesLastMonth int       `json:"purchases_last_month"`
}
