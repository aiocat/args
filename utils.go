package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Error struct to easily send json error messages
type Error struct {
	Message string `json:"error,omitempty"`
}

// For parse hcaptcha response
type HCaptchaResult struct {
	Success bool `json:"success"`
}

// Checks if key is valid
func HCaptchaChecker(resp string) bool {
	if resp == "" {
		return false
	}

	// Post given key
	request, err := http.NewRequest("POST", "https://hcaptcha.com/siteverify", bytes.NewBuffer([]byte("response="+resp+"&secret="+HCATPCHA_SECRET)))

	if err != nil {
		panic(err)
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Accept", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	result := new(HCaptchaResult)
	err = json.Unmarshal(body, &result)

	if err != nil {
		panic(err)
	}

	return result.Success
}

// Execute webhook to report argument
func ExecuteReportWebhook(argument *ArgumentReply) error {
	values := map[string]string{"username": "Report", "content": ("**Argument ID**: `" + argument.Id + "`\n**Argument**: ```" + argument.Argument + "```")}
	jsonToString, _ := json.Marshal(values)

	_, err := http.Post(REPORT_WEBHOOK_URL, "application/json", bytes.NewBuffer(jsonToString))
	return err
}
