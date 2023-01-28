// Common tools and helper functions
package common

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/akhilrex/hammond/db"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

// A helper function to generate random string
func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// A helper to convert from litres to gallon
func LitreToGallon(litres float32) float32 {
	gallonConversionFactor :=  0.21997
	return litres * float32(gallonConversionFactor);
}

// A helper to convert from gallon to litres
func GallonToLitre(gallons float32) float32 {
	litreConversionFactor :=  3.785412
	return gallons * float32(litreConversionFactor);
}


// A helper to convert from km to miles
func KmToMiles(km float32) float32 {
	kmConversionFactor :=  0.62137119
	return km * float32(kmConversionFactor);
}

// A helper to convert from miles to km
func MilesToKm(miles float32) float32 {
	milesConversionFactor := 1.609344
	return miles * float32(milesConversionFactor);
}



// A Util function to generate jwt_token which can be used in the request header
func GenToken(id string, role db.Role) (string, string) {
	jwt_token := jwt.New(jwt.GetSigningMethod("HS256"))
	// Set some claims
	jwt_token.Claims = jwt.MapClaims{
		"id":   id,
		"role": role,
		"exp":  time.Now().Add(time.Hour * 24 * 3).Unix(),
	}
	// Sign and get the complete encoded token as a string
	token, _ := jwt_token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["id"] = id
	rtClaims["exp"] = time.Now().Add(time.Hour * 24 * 20).Unix()

	rt, _ := refreshToken.SignedString([]byte(os.Getenv("JWT_SECRET")))

	return token, rt
}

// My own Error type that will help return my customized Error info
//  {"database": {"hello":"no such table", error: "not_exists"}}
type CommonError struct {
	Errors map[string]interface{} `json:"errors"`
}

// To handle the error returned by c.Bind in gin framework
// https://github.com/go-playground/validator/blob/v9/_examples/translations/main.go
func NewValidatorError(err error) CommonError {
	res := CommonError{}
	res.Errors = make(map[string]interface{})
	errs := err.(validator.ValidationErrors)
	for _, v := range errs {
		// can translate each error one at a time.
		//fmt.Println("gg",v.NameNamespace)
		if v.Param() != "" {
			res.Errors[v.Field()] = fmt.Sprintf("{%v: %v}", v.Tag(), v.Param())
		} else {
			res.Errors[v.Field()] = fmt.Sprintf("{key: %v}", v.Tag())
		}

	}
	return res
}

// Warp the error info in a object
func NewError(key string, err error) CommonError {
	res := CommonError{}
	res.Errors = make(map[string]interface{})
	res.Errors[key] = err.Error()
	return res
}

// Changed the c.MustBindWith() ->  c.ShouldBindWith().
// I don't want to auto return 400 when error happened.
// origin function is here: https://github.com/gin-gonic/gin/blob/master/context.go
func Bind(c *gin.Context, obj interface{}) error {
	b := binding.Default(c.Request.Method, c.ContentType())
	return c.ShouldBindWith(obj, b)
}
