package firebase

// import (
// 	"context"
// 	"os"

// 	"golang.org/x/oauth2/google"
// )

// // GetAccessToken retrieves a Firebase access token for FCM API calls
// func GetAccessToken() (string, error) {
// 	jsonPath := "firebase/service-account.json"
// 	data, err := os.ReadFile(jsonPath)
// 	if err != nil {
// 		return "", err
// 	}

// 	conf, err := google.JWTConfigFromJSON(data, "https://www.googleapis.com/auth/firebase.messaging")
// 	if err != nil {
// 		return "", err
// 	}

// 	token, err := conf.TokenSource(context.Background()).Token()
// 	if err != nil {
// 		return "", err
// 	}

// 	return token.AccessToken, nil
// }