package my_fb

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func CreateFirebaseApp() (*firebase.App, error){
	opt := option.WithCredentialsFile("./credentials/firebase_admin_credentials.json")
	
	ctx := context.Background()
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// firestoreCli, err := app.Firestore(ctx)
	// if err != nil {
	// log.Fatalln(err)
	// }
	// defer firestoreCli.Close()

	return app, nil
}