package controller

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"time"
	"topos-backend-assignment/db"
	"topos-backend-assignment/models"
)

type TokenMessage struct {
	Token       string `json:"token"`
	TokenExpiry string `json:"tokenExpiry"`
}

var JwtSecretKey []byte

type Claims struct {
	Uname string `json:"uname"`
	jwt.StandardClaims
}

func SetJWTSecret(jkey []byte) {
	JwtSecretKey = jkey
}

func IsTokenValid(request *http.Request) (bool, []byte) {
	var err error
	var valid bool
	valid = false
	var js []byte
	var msg ResponseMessage
	var tkn *jwt.Token

	h := request.Header
	tknStr := h.Get("Authorization")

	if tknStr == "" {
		Logger.Println("Bad Request. Please provide token in the Authorization Header")
		msg = ResponseMessage{Status: http.StatusBadRequest, ErrorMsg: "No Token provided", Message: "Bad Request. Please provide token in the Authorization Header"}
		js, err = json.Marshal(msg)
	} else {
		// Initialize a new instance of `Claims`
		claims := &Claims{}

		// Parse the JWT string and store the result in `claims`.
		// Note that we are passing the key in this method as well. This method will return an error
		// if the token is invalid (if it has expired according to the expiry time we set on sign in),
		// or if the signature does not match
		tkn, err = jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return JwtSecretKey, nil
		})
		if tkn == nil {
			Logger.Println("Invalid Token Provided")
			msg = ResponseMessage{Status: http.StatusUnauthorized, ErrorMsg: "Invalid Token Provided", Message: "Invalid Token Provided"}
			js, err = json.Marshal(msg)
		} else {
			if !tkn.Valid {
				Logger.Println("Unauthorized User")
				msg = ResponseMessage{Status: http.StatusUnauthorized, ErrorMsg: err.Error(), Message: "Unauthorized User"}
				js, err = json.Marshal(msg)
			} else if err != nil {
				if err == jwt.ErrSignatureInvalid {
					Logger.Println("Unauthorized User")
					msg = ResponseMessage{Status: http.StatusUnauthorized, ErrorMsg: err.Error(), Message: "Unauthorized User"}
					js, err = json.Marshal(msg)
				} else {
					Logger.Println("Bad Request")
					msg = ResponseMessage{Status: http.StatusBadRequest, ErrorMsg: err.Error(), Message: "Bad Request"}
					js, err = json.Marshal(msg)
				}
			} else {
				valid = true
			}
		}
	}
	return valid, js
}

/**
API URL - http://<host>:<port>/buildingFootprints/signup
Method	- POST
Params	- None
Request Body - JSON Object containing username and password
This API signs up a new user
*/
func SignUp(writer http.ResponseWriter, request *http.Request) {
	var msg ResponseMessage
	var js []byte
	var err error
	var creds models.UserLogin

	defer request.Body.Close()

	// Parse and decode the request body into a new `Credentials` instance
	err = json.NewDecoder(request.Body).Decode(&creds)
	if err != nil {
		// If there is something wrong with the request body, return a 400 status
		Logger.Println("Bad Request", err)
		msg = ResponseMessage{Status: http.StatusBadRequest, ErrorMsg: err.Error(), Message: "Bad Request"}
		js, err = json.Marshal(msg)
	} else {
		// Salt and hash the password using the bcrypt algorithm
		// The second argument is the cost of hashing, which we arbitrarily set as 8 (this value can be more or less, depending on the computing power you wish to utilize)
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), 8)

		if err != nil {
			Logger.Println("Error signing up user", err)
			msg = ResponseMessage{Status: http.StatusInternalServerError, ErrorMsg: err.Error(), Message: "Error signing up user"}
			js, err = json.Marshal(msg)
		} else {
			creds.ID = bson.NewObjectId()
			creds.Password = string(hashedPassword)
			if err = creds.CreateUser(db.MgoSession, creds); err != nil {
				Logger.Println("Error creating user", err)
				msg = ResponseMessage{Status: http.StatusInternalServerError, ErrorMsg: err.Error(), Message: "Error Creating User"}
				js, err = json.Marshal(msg)
			} else {
				msg = ResponseMessage{Status: http.StatusCreated, ErrorMsg: "None", Message: "User signed up successfully!"}
				js, err = json.Marshal(msg)
			}
		}
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(js)
}

/**
Helper method to generate and return signed JWT token
*/
func generateToken(u string) (string, time.Time, error) {
	var token *jwt.Token
	var expirationTime time.Time

	// Declare the expiration time of the token
	// here, we have kept it as 60 minutes
	expirationTime = time.Now().Add(60 * time.Minute)

	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Uname: u,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Create the JWT string
	tokenString, err := token.SignedString(JwtSecretKey)

	return tokenString, expirationTime, err
}

/**
API URL - http://<host>:<port>/buildingFootprints/login
Method	- POST
Params	- None
Request Body - JSON Object containing username and password
This API verifies User and returns JWT for verified User
*/
func LoginUser(writer http.ResponseWriter, request *http.Request) {
	var user models.UserLogin
	var js []byte
	var err error
	var msg ResponseMessage
	var storedPass string
	var tok string
	var expTime time.Time

	err = json.NewDecoder(request.Body).Decode(&user)
	if err != nil {
		// If there is something wrong with the request body, return a 400 status
		Logger.Println("Bad Request", err)
		msg = ResponseMessage{Status: http.StatusBadRequest, ErrorMsg: err.Error(), Message: "Bad Request"}
		js, err = json.Marshal(msg)
	} else {
		storedPass, err = user.LoginUser(db.MgoSession, user.Username)
		if err != nil {
			Logger.Println("Error verifying user", err)
			msg = ResponseMessage{Status: http.StatusInternalServerError, ErrorMsg: err.Error(), Message: "Error verifying user"}
			js, err = json.Marshal(msg)
		} else {
			if storedPass != "" {
				if err = bcrypt.CompareHashAndPassword([]byte(storedPass), []byte(user.Password)); err != nil {
					// If the two passwords don't match, return a 401 status
					Logger.Println("Unauthorized User")
					msg = ResponseMessage{Status: http.StatusUnauthorized, ErrorMsg: err.Error(), Message: "Unauthorized User"}
					js, err = json.Marshal(msg)
				} else {
					// Valid User, generate and set token
					tok, expTime, err = generateToken(user.Username)
					msg := TokenMessage{Token: tok, TokenExpiry: expTime.String()}
					js, err = json.Marshal(msg)
				}
			} else {
				Logger.Println("Unauthorized User")
				msg = ResponseMessage{Status: http.StatusUnauthorized, ErrorMsg: "", Message: "No Such User exists!"}
				js, err = json.Marshal(msg)
			}
		}
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(js)
}
