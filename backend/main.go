package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"github.com/google/uuid"
)

var routeBase = os.Getenv("ROUTE_BASE")
func getRoute(subroute string) string {
        return fmt.Sprintf("/%s/%s/", routeBase, subroute)
}
var staticRoute = getRoute("static")
var firebaseClient = initializeFirebase()
var firestoreClient = initializeFirestore()
var quiz = LoadQuiz()
var env = make(map[string]interface{})

func main() {
        log.Print("abcd: starting server...")
        log.Printf("routeBase='%s', staticRoute='%s'", routeBase, staticRoute)

        env["FACEBOOK_APP_ID"] = os.Getenv("FACEBOOK_APP_ID")
        env["RE_CAPTCHA_KEY"] = os.Getenv("RE_CAPTCHA_KEY")

        staticFileServer := http.FileServer(http.Dir("./static"))
        http.Handle(staticRoute, http.StripPrefix(staticRoute, staticFileServer))

        http.HandleFunc(getRoute("quiz"), templatedRouteFactory("quiz.html"))
        http.HandleFunc(getRoute("save"), handleSave)
        http.HandleFunc(getRoute("result"), templatedRouteFactory("result.html"))

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

func templatedRouteFactory(templateFile string) func(w http.ResponseWriter, r *http.Request) {
        return func (w http.ResponseWriter, r *http.Request) {
                lp := filepath.Join("templates", "layout.html")
                fp := filepath.Join("templates", templateFile)

                tmpl, err := template.ParseFiles(lp, fp)
                if err != nil {
                        fmt.Fprintf(w, "Error: %v", err)
                        return
                }
                templateData := make(map[string]interface{})
                templateData["routeBase"] = routeBase
                templateData["staticRoute"] = staticRoute
                templateData["quiz"] = quiz
                templateData["env"] = env
                templateData["quizJson"] = marshallToString(quiz)
                tmpl.ExecuteTemplate(w, "layout", templateData)
        }
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

func handleSave(w http.ResponseWriter, r *http.Request) {
    decoder := json.NewDecoder(r.Body)
    id := uuid.New().String()
    var result Result
    err := decoder.Decode(&result)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "Parse error: %v", err)
    }
    err = result.Validate()
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "Validation failed: %v", err)
    }
    response := make(map[string]interface{})
    response["id"] = id
    response["url"] = getRoute(fmt.Sprintf("result/%s", id))
    fmt.Fprintf(w, "%s", marshallToString(response))
}
