package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestSanitizer(t *testing.T) {
	bad1 := `Hello <STYLE>.XSS{background-image:url("javascript:alert('XSS')");}</STYLE><A CLASS=XSS></A>World`
	bad2 := `<a href="javascript:alert('XSS1')" onmouseover="alert('XSS2')">XSS<a>`
	markdown := `#Hello
	
	## This is cool
	
	*strong* and _emph_
	
	(A Link)[https://www.google.com]`

	expected1 := "Hello World"
	expected2 := "XSS"
	expectedMarkdown := markdown

	found1 := sanitize(bad1)
	found2 := sanitize(bad2)
	foundMD := sanitize(markdown)

	assert.Equal(t, expected1, found1)
	assert.Equal(t, expected2, found2)
	assert.Equal(t, expectedMarkdown, foundMD)
}

func TestTextGeneration(t *testing.T) {
	expectedBody := `Hello!
You have received the following contact 
Name: Kevin Eaton
Email: engineering@treelightsoftware.com
Subject: Unit Test

This is a test of the text generation

Phone: 555-555-5555`

	input := &MailerInput{
		Name:    "Kevin Eaton",
		Email:   "engineering@treelightsoftware.com",
		Subject: "Unit Test",
		Body:    "This is a test of the text generation",
		OtherData: map[string]string{
			"Phone": "555-555-5555",
		},
	}

	subject, body, err := generateText(input)
	assert.Nil(t, err)
	assert.True(t, strings.HasSuffix(subject, "Unit Test"))
	assert.Equal(t, expectedBody, body)

	badInput := &MailerInput{}
	subject, body, err = generateText(badInput)
	assert.NotNil(t, err)
	assert.Equal(t, "", subject)
	assert.Equal(t, "", body)
}

func TestSendErrorHelper(t *testing.T) {
	response, err := sendError("could not send that mail", http.StatusInternalServerError)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
	assert.Equal(t, "{\"message\":\"could not send that mail\"}", response.Body)
}

func TestSendSuccessHelper(t *testing.T) {
	input := map[string]string{
		"status": "success",
	}

	response, err := sendSuccess(input)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "{\"status\":\"success\"}", response.Body)
}

func TestEnvHelper(t *testing.T) {
	// if for some reason this is already set, we need to put it back after the test
	rand.Seed(time.Now().UnixNano())
	randID := rand.Int63()

	env := "TREEMAILER_TEST"
	original := os.Getenv(env)
	defer os.Setenv(env, original)
	os.Setenv(env, "testing")

	found := envHelper(env, "moo")
	assert.Equal(t, "testing", found)

	found = envHelper(fmt.Sprintf("TREEMAILER_%d", randID), "notfound")
	assert.Equal(t, "notfound", found)

}

func TestSendMail(t *testing.T) {
	// overwrite the values so we don't accidentally send emails in unit tests
	originalDomain := mgDomain
	originalKey := mgKey
	mgDomain = ""
	mgKey = ""
	response, id, err := sendMail("test@treelightsoftware.com", "testing@treelightsoftware.com",
		"cc@treelightsoftware.com", "Testing Send", "The message body")
	assert.NotNil(t, err)
	assert.Equal(t, "", response)
	assert.Equal(t, "", id)
	assert.Equal(t, "mailgun not configured", err.Error())

	// set them to dummy values
	mgDomain = "treelightsoftware.com"
	mgKey = "notarealkey"
	response, id, err = sendMail("test@treelightsoftware.com", "testing@treelightsoftware.com",
		"cc@treelightsoftware.com", "Testing Send", "The message body")
	assert.NotNil(t, err)

	mgDomain = originalDomain
	mgKey = originalKey
}

func TestHandler(t *testing.T) {
	// overwrite the values so we don't accidentally send emails in unit tests
	originalDomain := mgDomain
	originalKey := mgKey
	originalTo := to
	originalSiteName := siteName
	mgDomain = ""
	mgKey = ""
	to = ""
	siteName = ""

	request := events.APIGatewayProxyRequest{}

	response, err := Handler(request)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
	assert.Equal(t, "{\"message\":\"server is not configured properly\"}", response.Body)

	to = "testing@treelightsoftware.com"
	mgDomain = "treelightsoftware.com"
	mgKey = "notarealkey"
	siteName = "Treelight Software"

	// no json, so 400
	response, err = Handler(request)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)

	// send some json
	bytes, _ := json.Marshal(map[string]string{
		"from": "moo",
	})
	request.Body = string(bytes)
	response, err = Handler(request)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)

	// send some conformant json
	// this will fail since the domain is not setup correctly with Mailgun
	bytes, _ = json.Marshal(map[string]string{
		"name":    "Kevin Eaton",
		"email":   "testing@treelightsoftware.com",
		"subject": "Testing",
		"body":    "Testing a message",
	})
	request.Body = string(bytes)
	response, err = Handler(request)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusInternalServerError, response.StatusCode)

	mgDomain = originalDomain
	mgKey = originalKey
	to = originalTo
	siteName = originalSiteName
}
