package helpers

import (
	"context"
	"encoding/json"
	"go-server/constants"
	"go-server/db"
	"go-server/models"
	"log"
	"net/http"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func RequestBodyParser(req *http.Request, v interface{}) error {
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(v); err != nil {
		log.Fatal("body error")
		return err
	}
	return nil
}

func IsUserExists(email string) bool {
	filter := bson.M{"email": email}
	var existingUser models.User
	err := db.UsersEkart.FindOne(context.Background(), filter).Decode(&existingUser)
	return err == nil
}

func IsUserAlreadyThere(user models.UserObject) bool {
	filter := bson.M{
		"$or": []bson.M{
			{"email": user.Email},
			{"email": user.UsernameOrEmail},
			{"username": user.UsernameOrEmail},
			{"username": user.Username},
		},
	}

	var existingUser models.User
	err := db.UsersEkart.FindOne(context.Background(), filter).Decode(&existingUser)
	return err == nil
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func AuthenticateUser(user *models.UserObject)(*models.User, map[string]interface{}){
	filter := bson.M{
			"email": user.Email,
	}
	dynamicObject := make(map[string]interface{})
	var findeduser models.User
	err := db.UsersEkart.FindOne(context.Background(), filter).Decode(&findeduser)
	if err != nil {

		dynamicObject["errorMessage"]=constants.USER_DOEST_EXSISTS
		return nil, dynamicObject
		
	}
	err=bcrypt.CompareHashAndPassword([]byte(findeduser.Password),[]byte(user.Password))
	if err != nil {
		dynamicObject["errorMessage"]=constants.INVALID_CREDENTIALS

		return nil, dynamicObject
	}
	return &findeduser, nil
	//return string(hashedPassword), nil
}

func StoreUserSession(user *models.User)(*models.User,error){
	filter := bson.M{"email": user.Email}
	update := bson.M{"$set": bson.M{"access_token": user.Access_token}}
	_, err := db.UsersEkart.UpdateOne(context.Background(), filter, update)
	return nil,err
}

func CheckEmpty(s any) bool {
	return s == ""
}

func GenerateUuid() string {
	return uuid.New().String()
}

func SetCookie(cookieName string, cookieValue string,responseObject http.ResponseWriter){
	cookie := &http.Cookie{
		Name:     cookieName,
		Value:    cookieValue,
		Path:     "/",
		HttpOnly: true, // Make the cookie accessible only through HTTP
		SameSite: http.SameSiteLaxMode, // Set SameSite to "Lax"
	}

	http.SetCookie(responseObject, cookie)
}

func GetSessionIDFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie(constants.USER_SESSID_ENUM)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func MarshalWithExclusion(data interface{}, excludedField []string) (map[string]interface{}, error) {
    // Create a map to store the data to be marshaled
    dataMap := make(map[string]interface{})

    // Marshal the input data into a map
    dataJSON, err := json.Marshal(data)
    if err != nil {
        return nil, err
    }

    // Unmarshal the JSON into the map
    if err := json.Unmarshal(dataJSON, &dataMap); err != nil {
        return nil, err
    }

    for _, field := range excludedField {
		delete(dataMap, field)
	}

    // Marshal the modified map back into JSON
	by,err :=json.Marshal(dataMap)
	if err != nil{
		return nil,err
	}
	stringifyJson:=string(by)
	
	var parsedData map[string]interface{}
	if err := json.Unmarshal([]byte(stringifyJson), &parsedData); err != nil {
		
		return nil, err
	}
    return parsedData,nil
}