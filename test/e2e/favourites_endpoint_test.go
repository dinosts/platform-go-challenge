package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"platform-go-challenge/test"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFavourites(t *testing.T) {
	// Arrange
	server, token := test.StartServer()
	defer server.Close()

	client := server.Client()

	req, err := http.NewRequest(http.MethodGet, server.URL+"/v1/user/favourites", nil)
	req.Header.Add("Authorization", "bearer "+token)

	expected := map[string]any{
		"data": map[string]any{
			"charts": []any{
				map[string]any{
					"id":          "44444444-4444-4444-4444-444444444444",
					"description": "Main performance chart",
					"info": map[string]any{
						"id":           "11111111-1111-1111-1111-111111111111",
						"title":        "test chart",
						"x_axis_title": "commit number",
						"y_axis_title": "lines of code",
						"data": []any{
							map[string]any{"x": float64(1), "y": float64(100)},
							map[string]any{"x": float64(2), "y": float64(300)},
							map[string]any{"x": float64(3), "y": float64(500)},
						},
					},
				},
			},
			"insights": []any{
				map[string]any{
					"id":          "55555555-5555-5555-5555-555555555555",
					"description": "Great for Q2 presentation",
					"info": map[string]any{
						"Id":   "22222222-2222-2222-2222-222222222222",
						"Text": "40% of millennials spend more than 3 hours on social media daily",
					},
				},
			},
			"audiences": []any{
				map[string]any{
					"id":          "66666666-6666-6666-6666-666666666666",
					"description": "Target audience for campaign",
					"info": map[string]any{
						"id":                   "33333333-3333-3333-3333-333333333333",
						"gender":               "Male",
						"birth_country":        "United Kingdom",
						"age_group":            "25-34",
						"social_media_hours":   float64(3.5),
						"purchases_last_month": float64(7),
					},
				},
			},
		},
		"pagination": map[string]any{
			"page":     float64(0),
			"pageSize": float64(10),
			"maxPage":  float64(0),
		},
	}
	// Act
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Assert
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var result map[string]any
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestCreateFavourite(t *testing.T) {
	// Arrange
	server, token := test.StartServer()
	defer server.Close()

	requestBody := map[string]interface{}{
		"assetId":     "22222222-2222-2222-2222-222222222223",
		"description": "Good to know",
	}

	client := server.Client()

	bodyBytes, _ := json.Marshal(requestBody)
	req, err := http.NewRequest(http.MethodPost, server.URL+"/v1/user/favourites", bytes.NewReader(bodyBytes))
	req.Header.Add("Authorization", "bearer "+token)

	expected := map[string]interface{}{
		"user_id":     "a3973a1c-a77b-4a04-a296-ddec19034419",
		"asset_id":    "22222222-2222-2222-2222-222222222223",
		"asset_type":  "insight",
		"description": "Good to know",
	}

	// Act
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Assert
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	var result map[string]any
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)

	data := result["data"].(map[string]any)
	assert.Equal(t, data["user_id"], expected["user_id"])
	assert.Equal(t, data["asset_id"], expected["asset_id"])
	assert.Equal(t, data["asset_type"], expected["asset_type"])
	assert.Equal(t, data["description"], expected["description"])
}

func TestUpdateFavourite(t *testing.T) {
	// Arrange
	server, token := test.StartServer()
	defer server.Close()

	requestBody := map[string]interface{}{
		"description": "Good to know 2",
	}

	client := server.Client()

	bodyBytes, _ := json.Marshal(requestBody)
	req, err := http.NewRequest(http.MethodPatch, server.URL+"/v1/user/favourites/55555555-5555-5555-5555-555555555555", bytes.NewReader(bodyBytes))
	req.Header.Add("Authorization", "bearer "+token)

	expected := map[string]any{
		"id":          "55555555-5555-5555-5555-555555555555",
		"user_id":     "a3973a1c-a77b-4a04-a296-ddec19034419",
		"asset_id":    "22222222-2222-2222-2222-222222222222",
		"asset_type":  "insight",
		"description": "Good to know 2",
	}

	// Act
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Assert
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var result map[string]any
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)

	data := result["data"].(map[string]any)
	assert.Equal(t, data["user_id"], expected["user_id"])
	assert.Equal(t, data["asset_id"], expected["asset_id"])
	assert.Equal(t, data["asset_type"], expected["asset_type"])
	assert.Equal(t, data["description"], expected["description"])
}

func TestDeleteFavourite(t *testing.T) {
	// Arrange
	server, token := test.StartServer()
	defer server.Close()

	requestBody := map[string]interface{}{
		"description": "Good to know 2",
	}

	client := server.Client()

	bodyBytes, _ := json.Marshal(requestBody)
	req, err := http.NewRequest(http.MethodDelete, server.URL+"/v1/user/favourites/55555555-5555-5555-5555-555555555555", bytes.NewReader(bodyBytes))
	req.Header.Add("Authorization", "bearer "+token)
	expected := map[string]string{"message": "Favourite deleted"}
	// Act
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Assert
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var result map[string]string
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)

	assert.Equal(t, expected, result)
}
