package database

import (
	"github.com/google/uuid"
)

func IMpopulateStorageForDevEnv(
	db *IMDatabase,
	passwordHasher func(string) string,
) {
	// Constant UUIDs
	userId, _ := uuid.Parse("a3973a1c-a77b-4a04-a296-ddec19034419")
	chartId, _ := uuid.Parse("11111111-1111-1111-1111-111111111111")
	insightId, _ := uuid.Parse("22222222-2222-2222-2222-222222222222")
	insightId2, _ := uuid.Parse("22222222-2222-2222-2222-222222222223")
	audienceId, _ := uuid.Parse("33333333-3333-3333-3333-333333333333")
	favChartId, _ := uuid.Parse("44444444-4444-4444-4444-444444444444")
	favInsightId, _ := uuid.Parse("55555555-5555-5555-5555-555555555555")
	favAudienceId, _ := uuid.Parse("66666666-6666-6666-6666-666666666666")

	// User
	devUser := IMUserModel{
		Id:       userId,
		Email:    "test@test.com",
		Password: passwordHasher("pass"),
	}
	(db.UserStorage)[devUser.Id] = devUser

	// Chart
	chart := IMChartModel{
		Id:         chartId,
		Title:      "test chart",
		XAxisTitle: "commit number",
		YAxisTitle: "lines of code",
		Data: []map[string]float64{
			{"x": 1, "y": 100},
			{"x": 2, "y": 300},
			{"x": 3, "y": 500},
		},
	}
	db.ChartStorage[chart.Id] = chart

	// Insight
	insight := IMInsightModel{
		Id:   insightId,
		Text: "40% of millennials spend more than 3 hours on social media daily",
	}
	db.InsightStorage[insight.Id] = insight
	insight2 := IMInsightModel{
		Id:   insightId2,
		Text: "100% of zoomers spend more than 8 hours on watching memes",
	}
	db.InsightStorage[insight2.Id] = insight2

	// Audience
	audience := IMAudienceModel{
		Id:                 audienceId,
		Gender:             "Male",
		BirthCountry:       "United Kingdom",
		AgeGroup:           "25-34",
		SocialMediaHours:   3.5,
		PurchasesLastMonth: 7,
	}
	db.AudienceStorage[audience.Id] = audience

	// Favourites
	fav1 := IMFavouriteModel{
		Id:          favChartId,
		UserId:      devUser.Id,
		AssetId:     chart.Id,
		AssetType:   "chart",
		Description: "Main performance chart",
	}
	fav2 := IMFavouriteModel{
		Id:          favInsightId,
		UserId:      devUser.Id,
		AssetId:     insight.Id,
		AssetType:   "insight",
		Description: "Great for Q2 presentation",
	}
	fav3 := IMFavouriteModel{
		Id:          favAudienceId,
		UserId:      devUser.Id,
		AssetId:     audience.Id,
		AssetType:   "audience",
		Description: "Target audience for campaign",
	}
	db.FavouriteStorage[fav1.Id] = fav1
	db.FavouriteStorage[fav2.Id] = fav2
	db.FavouriteStorage[fav3.Id] = fav3
}
