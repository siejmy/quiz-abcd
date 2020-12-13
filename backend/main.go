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
	"strings"
	"time"

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
var resultRepository = ResultRepositoryFirestore{
    firestoreClient,
}
var env = make(map[string]interface{})

func main() {
    log.Print("abcd: starting server...")
    log.Printf("routeBase='%s', staticRoute='%s'", routeBase, staticRoute)

    env["FACEBOOK_APP_ID"] = os.Getenv("FACEBOOK_APP_ID")
    env["RE_CAPTCHA_KEY"] = os.Getenv("RE_CAPTCHA_KEY")

    staticFileServer := http.FileServer(http.Dir("./static"))
    http.Handle(staticRoute, http.StripPrefix(staticRoute, staticFileServer))

    http.HandleFunc(fmt.Sprintf("/%s/", routeBase), handleQuiz)
    http.HandleFunc(getRoute("save"), handleSave)
    http.HandleFunc(getRoute("result"), handleResult)

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("abcd: listening on port %s", port)
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

func handleQuiz(w http.ResponseWriter, r *http.Request) {
    templateData := make(map[string]interface{})
    appendDefaultTemplateData(&templateData)
    respondWithTemplate(w, "quiz.html", templateData)
}

func handleResult(w http.ResponseWriter, r *http.Request) {
    pathParts := removeEmptyFromArray(strings.Split(r.URL.Path, "/"))
    if(len(pathParts) < 1) {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "Invalid url")
        return
    }
    id := pathParts[len(pathParts) - 1]
    result, err := resultRepository.GetByID(id)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "Cannot get result: %v", err)
        return
    }

    statsSummary, err := GetStatsSummary()
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "Cannot get statistics: %v", err)
        return
    }

    statsEntry := GetStatsEntryForResult(quiz, result)

    templateData := make(map[string]interface{})
    appendDefaultTemplateData(&templateData)
    templateData["id"] = id
    templateData["result"] = result
    templateData["resultMarshalled"] = marshallToString(result)
    templateData["statsEntry"] = statsEntry
    templateData["statsEntryMarshalled"] = marshallToString(statsEntry)
    templateData["statsSummary"] = statsSummary
    templateData["statsSummaryMarshalled"] = marshallToString(statsSummary)
    respondWithTemplate(w, "result.html", templateData)
}

func handleSave(w http.ResponseWriter, r *http.Request) {
    // throttle
    time.Sleep(200 * time.Millisecond)

    id := uuid.New().String()
    var result Result
    err := json.NewDecoder(r.Body).Decode(&result)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "POST parse error: %+v, parsed: %+v", err, result)
        return
    }
    err = result.Validate()
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "Validation failed: %v", err)
        return
    }

    err = resultRepository.Save(id, result)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "Cannot save to DB: %v", err)
        return
    }

    err = WriteStats(firestoreClient, quiz, result)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "Cannot save stats to DB: %v", err)
        return
    }

    response := make(map[string]interface{})
    response["id"] = id
    response["url"] = fmt.Sprintf("/%s/result/%s/", routeBase, id)
    fmt.Fprintf(w, "%s", marshallToString(response))
}

func respondWithTemplate(w http.ResponseWriter, templateFile string, templateData map[string]interface{}) {
    lp := filepath.Join("templates", "layout.html")
    fp := filepath.Join("templates", templateFile)

    tmpl, err := template.ParseFiles(lp, fp)
    if err != nil {
        fmt.Fprintf(w, "Error: %v", err)
        return
    }
    tmpl.ExecuteTemplate(w, "layout", templateData)
}

func appendDefaultTemplateData(templateData *map[string]interface{}) {
    (*templateData)["routeBase"] = routeBase
    (*templateData)["staticRoute"] = staticRoute
    (*templateData)["quiz"] = quiz
    (*templateData)["env"] = env
    (*templateData)["quizJson"] = marshallToString(quiz)
}
