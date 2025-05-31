package config

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

var FirebaseApp *firebase.App

func InitFirebase() {
	opt := option.WithCredentialsFile("firebase/firebase-service-account.json")

	app, err := firebase.NewApp(context.Background(), &firebase.Config{
		ProjectID: "", // <-- Isi dengan project ID Firebase kamu
	}, opt)
	if err != nil {
		log.Fatalf("Firebase init error: %v\n", err)
	}
	FirebaseApp = app
}
