package middlewares

import (
	"context"
	"encoding/json"
	"go-server/constants"
	"go-server/db"
	"go-server/helpers"
	"go-server/models"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

func AuthenticateMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var response models.Response

		// Get the session ID from the cookie
		sessionID, err := helpers.GetSessionIDFromCookie(r)
		if err != nil {
			response.Code = constants.ERROR_CODE
			response.Message = constants.UNAUTHORIZED
			jsonResponse, _ := json.Marshal(response)
			w.WriteHeader(http.StatusConflict)
			w.Write(jsonResponse)
			return
		}

		// Search for the user based on the session ID
		user, err := findUserBySessionID(sessionID)
		if err != nil {
			response.Code = constants.ERROR_CODE
			response.Message = constants.UNAUTHORIZED
			jsonResponse, _ := json.Marshal(response)
			w.WriteHeader(http.StatusConflict)
			w.Write(jsonResponse)
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)

		// Continue with the next handler if the user is authenticated
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func findUserBySessionID(sessionID string) (*models.User, error) {
	filter := bson.M{
		"access_token": sessionID,
	}
	var user *models.User
	err := db.UsersEkart.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
