package main

import (
	"github.com/ErikBjare/Futarchio/src/api"
	"github.com/emicklei/go-restful"
	"github.com/golang/oauth2"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func oauth_test() {
	file, err := os.Open("secrets/key.pem")
	if err != nil {
		panic(err)
	}
	key := []byte{}
	file.Read(key)

	conf, err := oauth2.NewJWTConfig(&oauth2.JWTOptions{
		Email: "643992545442-u8dubmhq38dor5bvltjb2o98tv3musqq@developer.gserviceaccount.com",
		// The contents of your RSA private key or your PEM file
		// that contains a private key.
		// If you have a p12 file instead, you
		// can use `openssl` to export the private key into a pem file.
		//
		//    $ openssl pkcs12 -in key.p12 -out key.pem -nodes
		//
		// It only supports PEM containers with no passphrase.
		PrivateKey: key,
		Scopes:     []string{"profile"},
	},
		"https://provider.com/o/oauth2/token")
	if err != nil {
		log.Fatal(err)
	}

	// Initiate an http.Client, the following GET request will be
	// authorized and authenticated on the behalf of
	// xxx@developer.gserviceaccount.com.
	client := http.Client{Transport: conf.NewTransport()}
	client.Get("...")

	// If you would like to impersonate a user, you can
	// create a transport with a subject. The following GET
	// request will be made on the behalf of user@example.com.
	client = http.Client{Transport: conf.NewTransportWithUser("user@example.com")}
	client.Get("...")
}

func serve(wsContainer *restful.Container) {
	mux := http.NewServeMux()
	mux.Handle("/api/0/", wsContainer)
	mux.Handle("/", http.FileServer(http.Dir("site/dist")))
	server := &http.Server{Addr: ":8080", Handler: mux}

	log.Println("Frontend is serving on: http://localhost:8080")
	log.Println("API is serving on: http://localhost:8080/api/")
	server.ListenAndServe()
}

func main() {
	log.Println("Started...")
	rand.Seed(time.Now().Unix())

	wsContainer := restful.NewContainer()

	api.Users.Register(wsContainer)
	api.Users.Init()

	go serve(wsContainer)

	queue := make(chan error)
	for {
		err := <-queue
		log.Println(err)
	}

	log.Println("Quitting")
}
