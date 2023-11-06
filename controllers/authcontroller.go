package controllers

import (
	"context"
	"encoding/json"
	"go-server/constants"
	"go-server/db"
	"go-server/helpers"
	"go-server/middlewares"
	"go-server/models"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/register", registerProfile).Methods("POST")
	r.HandleFunc("/login", loginProfile).Methods("POST")
	r.HandleFunc("/getUserProfile", middlewares.AuthenticateMiddleware(getUserProfile)).Methods("GET")
}

func registerProfile(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is POST
	var response models.Response
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		response.Code = constants.NOT_FOUND_CODE
		response.Message = constants.INVALID_REQUEST_METHOD
		jsonResponse, _ := json.Marshal(response)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(jsonResponse)
		return
	}

	// Parse the request body into a User struct
	var newUser models.User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		response.Code = constants.ERROR_CODE
		response.Message = constants.FAILED_TO_PARSE_BODY
		jsonResponse, _ := json.Marshal(response)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonResponse)
		return
	}

	// Check for duplicate user validation by email
	if helpers.IsUserExists(newUser.Email) {
		response.Code = constants.ERROR_CODE
		response.Message = constants.USER_ALREADY_EXSISTS
		jsonResponse, _ := json.Marshal(response)
		w.WriteHeader(http.StatusConflict)
		w.Write(jsonResponse)
		return
	}
	if helpers.CheckEmpty(newUser.Password) {
		response.Code = constants.ERROR_CODE
		response.Message = constants.PASSWORD_NIL
		jsonResponse, _ := json.Marshal(response)
		w.WriteHeader(http.StatusConflict)
		w.Write(jsonResponse)
		return
	}

	// Hash the password before storing it
	hashedPassword, err := helpers.HashPassword(newUser.Password)
	if err != nil {
		response.Code = constants.ERROR_CODE
		response.Message = constants.GENERIC_ERROR_MESSAGE
		jsonResponse, _ := json.Marshal(response)
		w.WriteHeader(http.StatusConflict)
		w.Write(jsonResponse)
		return
	}
	newUser.Password = hashedPassword
	if newUser.UserRole == "" {
		newUser.UserRole = constants.USER
	}

	//generate UUid for the user
	newUser.UUid = helpers.GenerateUuid()

	// Insert the user data into the MongoDB collection
	insertResult, err := db.UsersEkart.InsertOne(context.Background(), newUser)

	if err != nil {
		response.Code = constants.ERROR_CODE
		response.Message = constants.GENERIC_ERROR_MESSAGE
		jsonResponse, _ := json.Marshal(response)
		w.WriteHeader(http.StatusConflict)
		w.Write(jsonResponse)
		return
	}

	// Respond with a success message

	response.Code = constants.OK
	response.Message = constants.USER_ADDED_SUCCESS
	response.Data = insertResult.InsertedID.(primitive.ObjectID).Hex()
	//response := map[string]string{"code": constants.success, "message": "User registered successfully", "userId": insertResult.InsertedID.(primitive.ObjectID).Hex()}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		response.Code = constants.ERROR_CODE
		response.Message = constants.GENERIC_ERROR_MESSAGE
		jsonResponse, _ := json.Marshal(response)
		w.WriteHeader(http.StatusConflict)
		w.Write(jsonResponse)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(jsonResponse)
}
func loginProfile(w http.ResponseWriter, r *http.Request) {
	var response models.Response
	w.Header().Set("Content-Type", "application/json")
	// Parse the request body into a User struct
	var User models.UserObject
	err := json.NewDecoder(r.Body).Decode(&User)
	if err != nil {
		response.Code = constants.ERROR_CODE
		response.Message = constants.FAILED_TO_PARSE_BODY
		jsonResponse, _ := json.Marshal(response)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonResponse)
		return
	}

	// Check for  user exsist validation by email
	if !helpers.IsUserAlreadyThere(User) {
		response.Code = constants.ERROR_CODE
		response.Message = constants.USER_DOEST_EXSISTS
		jsonResponse, _ := json.Marshal(response)
		w.WriteHeader(http.StatusConflict)
		w.Write(jsonResponse)
		return
	}
	//check for empty password
	if helpers.CheckEmpty(User.Password) {
		response.Code = constants.ERROR_CODE
		response.Message = constants.PASSWORD_NIL
		jsonResponse, _ := json.Marshal(response)
		w.WriteHeader(http.StatusConflict)
		w.Write(jsonResponse)
		return
	}

	user, authenticareError := helpers.AuthenticateUser(&User)

	if authenticareError != nil {
		response.Code = constants.ERROR_CODE
		response.Message = authenticareError["errorMessage"].(string)
		jsonResponse, _ := json.Marshal(response)
		w.WriteHeader(http.StatusConflict)
		w.Write(jsonResponse)
		return
	}

	user.Access_token = helpers.GenerateUuid()
	_, err = helpers.StoreUserSession(user)
	if err != nil {
		response.Code = constants.ERROR_CODE
		response.Message = constants.GENERIC_ERROR_MESSAGE
		jsonResponse, _ := json.Marshal(response)
		w.WriteHeader(http.StatusConflict)
		w.Write(jsonResponse)
		return
	}

	helpers.SetCookie(constants.USER_SESSID_ENUM, user.Access_token, w)

	response.Code = constants.OK
	response.Message = constants.LOGIN_SUCCESS
	formatterres, err := helpers.MarshalWithExclusion(user, []string{"password", "_id", "CreatedAt", "UpdatedAt"})
	response.Data = formatterres
	jsonResponse, _ := json.Marshal(response)
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonResponse)
	return
}
func getUserProfile(w http.ResponseWriter, r *http.Request) {
	var response models.Response

	user := r.Context().Value("user").(*models.User)

	// if !ok {
	// 	response.Code = constants.ERROR_CODE
	// 	response.Message = constants.UNAUTHORIZED
	// 	jsonResponse, _ := json.Marshal(response)
	// 	w.WriteHeader(http.StatusConflict)
	// 	w.Write(jsonResponse)
	// 	return
	// }

	response.Code = constants.OK
	response.Message = constants.DATA_FETCH_SUCCESS
	formatterres, _ := helpers.MarshalWithExclusion(user, []string{"password", "_id", "CreatedAt", "UpdatedAt", "access_token"})
	response.Data = formatterres
	jsonResponse, _ := json.Marshal(response)
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonResponse)
	return
}
