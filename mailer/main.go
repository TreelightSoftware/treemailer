package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/mailgun/mailgun-go"
	"github.com/microcosm-cc/bluemonday"
)

// Set these as you see fit
// a good task would be to break these into environment variables
const (
	to       string = ""
	mgDomain string = ""
	mgKey    string = ""
	siteName string = ""
	cc       string = ""
)

var _sanitizer *bluemonday.Policy

// MailerInput is the input for the JSON of the form
type MailerInput struct {
	Name      string            `json:"name"`
	Email     string            `json:"email"`
	Subject   string            `json:"subject"`
	Body      string            `json:"body"`
	OtherData map[string]string `json:"otherData,omitempty"`
}

// MailResponse is a structure response
type MailResponse struct {
	Message string `json:"message"`
}

// main is the main entry that sets up the handler
func main() {
	lambda.Start(Handler)
}

// Handler handles the Lambda request. This is the bulk of the logic
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	input := &MailerInput{}

	// if the environment is not setup correctly, we return a 500
	if to == "" || mgDomain == "" || mgKey == "" || siteName == "" {
		return sendError("server is not configured properly", http.StatusInternalServerError)
	}

	// unmarshal the JSON
	err := json.Unmarshal([]byte(request.Body), input)
	if err != nil {
		return sendError(err.Error(), http.StatusBadRequest)
	}

	subject, body, err := generateText(input)
	if err != nil {
		return sendError(err.Error(), http.StatusBadRequest)
	}

	_, _, err = sendMail(to, input.Email, cc, subject, body)
	if err != nil {
		return sendError(err.Error(), http.StatusInternalServerError)
	}

	return sendSuccess(input)
}

// sendError sends an error to the API Gateway
func sendError(message string, code int) (events.APIGatewayProxyResponse, error) {
	response := MailResponse{}
	ret := []byte{}
	response.Message = message
	ret, _ = json.Marshal(response)
	return events.APIGatewayProxyResponse{Body: string(ret), StatusCode: code}, nil
}

// sendSuccess sends a success message to the gateway
func sendSuccess(retStruct interface{}) (events.APIGatewayProxyResponse, error) {
	ret, _ := json.Marshal(retStruct)
	return events.APIGatewayProxyResponse{Body: string(ret), StatusCode: 200}, nil
}

// sendMail sends the mail and returns information about the message from mailgun
func sendMail(to, from, cc, subject, body string) (string, string, error) {
	mg := mailgun.NewMailgun(mgDomain, mgKey)
	message := mg.NewMessage(
		from,
		subject,
		body)
	message.AddRecipient(to)
	if cc != "" {
		message.AddCC(cc)
	}
	message.SetReplyTo(from)

	return mg.Send(message)
}

// generateText generates the subject and body of the email based upon the input
func generateText(input *MailerInput) (subject, body string, err error) {

	// sanitize all the things
	input.Name = sanitize(input.Name)
	input.Email = sanitize(input.Email)
	input.Subject = sanitize(input.Subject)
	input.Body = sanitize(input.Body)

	// make sure the basic fields are there
	if input.Name == "" || input.Email == "" || input.Subject == "" || input.Body == "" {
		err = errors.New("name, email, subject, and body are all required")
		return
	}

	body = fmt.Sprintf("Hello!\nYou have received the following contact \nName: %s\n Email: %s\n Subject: %s\n %s", input.Name, input.Email, input.Subject, input.Body)
	if input.OtherData != nil {
		for k, v := range input.OtherData {
			body = fmt.Sprintf("%s\n%s: %v", body, sanitize(k), sanitize(v))
		}
	}

	subject = fmt.Sprintf("New Contact from %s: %s", siteName, input.Subject)
	return
}

// sanitize uses the sanitizer to make sure the text is clear of various bad things
func sanitize(input string) string {
	clean := _sanitizer.Sanitize(input)
	return clean
}
