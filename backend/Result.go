package main

import "gopkg.in/validator.v2"

// Result — written result of the quiz
type Result struct {
  ID string `json:"id"`
  Name string `json:"name" validate:"min=0,max=50,regexp=^[a-zA-Z1-9]*$"`
  Answers []int8 `json:"answers" validate:"min=0"`
}

// Validate validates a result to be written
func (result Result) Validate() error {
	return validator.Validate(result);
}
