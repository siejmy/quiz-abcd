package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
)

var routeBase = os.Getenv("ROUTE_BASE")
func getRoute(subroute string) string {
        return fmt.Sprintf("/%s/%s", routeBase, subroute)
}
var staticRoute = getRoute("static")
var firebaseClient = initializeFirebase()
var firestoreClient = initializeFirestore()

func main() {
        log.Print("abcd: starting server...")
        log.Printf("routeBase='%s', staticRoute='%s'", routeBase, staticRoute)


        fs := http.FileServer(http.Dir("./static"))
        http.Handle(staticRoute, http.StripPrefix(staticRoute, fs))

        http.HandleFunc(getRoute(""), serveTemplate)
        http.HandleFunc(getRoute("demo"), serveFirestore)

        port := os.Getenv("PORT")
        if port == "" {
                port = "8080"
        }

        log.Printf("helloworld: listening on port %s", port)
        log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func initializeFirebase() *firebase.App {
        app, err := firebase.NewApp(context.Background(), nil)
        if err != nil {
                log.Fatalf("error initializing app: %v\n", err)
        }

        return app
}

func initializeFirestore() *firestore.Client {
        client, err := firebaseClient.Firestore(context.Background())
        if err != nil {
                log.Fatalf("error initializing firestore: %v\n", err)
        }
        return client
}

func serveTemplate(w http.ResponseWriter, r *http.Request) {
        lp := filepath.Join("templates", "layout.html")
        fp := filepath.Join("templates", "home.html")

        tmpl, err := template.ParseFiles(lp, fp)
        if err != nil {
                fmt.Fprintf(w, "Error: %v", err)
                return
        }
        templateData := make(map[string]interface{})
        templateData["routeBase"] = routeBase
        templateData["staticRoute"] = staticRoute
        tmpl.ExecuteTemplate(w, "layout", templateData)
}

func serveFirestore(w http.ResponseWriter, r *http.Request) {
        doc := firestoreClient.Doc("firestore_demo/demo")
        snapshot, err := doc.Get(context.Background())
        if err != nil {
                fmt.Fprintf(w, "Error: %v", err)
                return
        }
        if snapshot.Exists() != true {
                fmt.Fprintf(w, "Document doesnt exist")
                return
        }
        data := snapshot.Data()
        fmt.Fprintf(w, "Data: %+v", data)
}
