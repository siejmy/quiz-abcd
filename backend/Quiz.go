package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

// Quiz is configuration of the quiz itself
type Quiz struct {
	ID string `json:"id"`
  Type string `json:"type"`
  Title string `json:"title"`
  IntroHTML string `json:"introHtml"`
  IntroImageURL string `json:"introImageUrl"`
  Questions QuizABCDQuestion `json:"questions"`
}

// QuizABCDQuestion is a question type
type QuizABCDQuestion struct {
  Title string `json:"title"`
  ImageURL string `json:"imageUrl"`
  Distractors []string  `json:"distractors"`
  CorrectNo int8 `json:"correctNo"`
}

// LoadQuiz loads quiz from file
func LoadQuiz() Quiz {
	path := "./quiz.json"
	jsonFile, err := os.Open(path)
	if err != nil {
			log.Panicf("Cannot open quiz file %s: %v", path, err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result Quiz
	err = json.Unmarshal([]byte(byteValue), &result)
	if err != nil {
		log.Panicf("Cannot parse quiz file %s: %v", path, err)
	}

	return result
}
