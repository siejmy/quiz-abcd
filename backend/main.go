package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

var routeBase = os.Getenv("ROUTE_BASE")
var domainName = os.Getenv("DOMAIN")
func getRoute(subroute string) string {
    return fmt.Sprintf("/%s/%s/", routeBase, subroute)
}
var staticRoute = getRoute("static")
var quiz = LoadQuiz()
var resultRepository = ResultRepositoryFirestore{}
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
    resultPercent := 100 * statsEntry.CorrectCount / statsEntry.TotalCount
    decileIndex := resultPercent / 10
    decileValue := getDecileValue(resultPercent, statsSummary)

    templateData := make(map[string]interface{})
    appendDefaultTemplateData(&templateData)
    templateData["id"] = id
    templateData["result"] = result
    templateData["resultMarshalled"] = marshallToString(result)
    templateData["statsEntry"] = statsEntry
    templateData["statsEntryMarshalled"] = marshallToString(statsEntry)
    templateData["statsSummary"] = statsSummary
    templateData["statsSummaryMarshalled"] = marshallToString(statsSummary)
    templateData["resultPercent"] = resultPercent
    templateData["decileIndex"] = decileIndex
    templateData["decileValue"] = decileValue
    templateData["title"] = getResultTitle(result, statsEntry)
    templateData["resultUrl"] = fmt.Sprintf("https://%s/%s/result/%s/", domainName, routeBase, id)
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

    err = WriteStats(quiz, result)
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

    templateWithFuncs := template.New("").Funcs(template.FuncMap{
        "seq": func(n int) []int { return make([]int, n) },
        "sub": func(a int, b int) int { return a - b },
        "urlencode": func(decoded string) string { return url.QueryEscape((decoded)) },
    })
    tmpl, err := templateWithFuncs.ParseFiles(lp, fp)
    if err != nil {
        fmt.Fprintf(w, "Error: %v", err)
        return
    }
    err = tmpl.ExecuteTemplate(w, "layout", templateData)
    if err != nil {
        log.Printf("Error while parsing template: %v", err)
    }
}

func appendDefaultTemplateData(templateData *map[string]interface{}) {
    (*templateData)["routeBase"] = routeBase
    (*templateData)["staticRoute"] = staticRoute
    (*templateData)["quiz"] = quiz
    (*templateData)["env"] = env
    (*templateData)["year"] = time.Now().UTC().Year()
    (*templateData)["quizJson"] = marshallToString(quiz)
    (*templateData)["quizUrl"] = fmt.Sprintf("https://%s/%s/", domainName, routeBase)
    (*templateData)["domainName"] = domainName
}

func getResultTitle(result Result, stats StatsEntry) string {
    resultPercent := 100 * stats.CorrectCount / stats.TotalCount
    if len(result.Name) > 0 {
        return fmt.Sprintf("%s uzyskał %d%% w quizie %s — Siejmy QUIZ", result.Name, resultPercent, quiz.Title)
    }
    return fmt.Sprintf("%d%% w quizie %s — Siejmy QUIZ", resultPercent, quiz.Title)
}

func getDecileValue(resultPercent int, stats StatsSummary) int {
    decileIndex := resultPercent / 10
    noOfWorse := 0
    for i := 0; i <= decileIndex;i++ {
        noOfWorse += stats.DecileHistogram[i]
    }
    return 100 * noOfWorse / stats.SampleCount
}
