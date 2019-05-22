package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/mailgun/mailgun-go"
	"github.com/microcosm-cc/bluemonday"
)

// Set these as you see fit
const (
	to       string = ""
	mgDomain string = ""
	mgKey    string = ""
	siteName string = ""
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

// Handler handles the Lambda request
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	input := &MailerInput{}
	response := MailResponse{}
	ret := []byte{}

	// if the environemtn is not setup correctly, we return a 500
	if to == "" || mgDomain == "" || mgKey == "" || siteName == "" {
		response.Message = "server is not configured properly"
		ret, _ = json.Marshal(response)
		return events.APIGatewayProxyResponse{Body: string(ret), StatusCode: 500}, nil
	}

	// unmarshal the JSON
	err := json.Unmarshal([]byte(request.Body), input)
	if err != nil {
		response.Message = err.Error()
		ret, _ = json.Marshal(response)
		return events.APIGatewayProxyResponse{Body: string(ret), StatusCode: 400}, nil
	}

	// sanitize all the things
	input.Name = sanitize(input.Name)
	input.Email = sanitize(input.Email)
	input.Subject = sanitize(input.Subject)
	input.Body = sanitize(input.Body)

	// make sure the basic fields are there
	if input.Name == "" || input.Email == "" || input.Subject == "" || input.Body == "" {
		response.Message = "name, email, subject, and body are all required"
		ret, _ = json.Marshal(response)
		return events.APIGatewayProxyResponse{Body: string(ret), StatusCode: 400}, nil
	}

	// now we build the email
	mg := mailgun.NewMailgun(mgDomain, mgKey)

	body := fmt.Sprintf("Hello!\nYou have received the following contact \nName: %s\n Email: %s\n Subject: %s\n %s", input.Name, input.Email, input.Subject, input.Body)
	if input.OtherData != nil {
		for k, v := range input.OtherData {
			body = fmt.Sprintf("%s\n%s: %v", body, sanitize(k), sanitize(v))
		}
	}

	subject := fmt.Sprintf("New Contact from %s: %s", siteName, input.Subject)

	m := mg.NewMessage(
		input.Email,
		subject,
		body)
	m.AddRecipient(to)
	m.AddCC("kevin.eaton@kvsstechnologies.com")
	m.SetReplyTo(input.Email)

	_, _, err = mg.Send(m)

	ret, _ = json.Marshal(input)
	return events.APIGatewayProxyResponse{Body: string(ret), StatusCode: 200}, nil
}

func main() {
	_sanitizer = bluemonday.StrictPolicy()
	lambda.Start(Handler)
}

func sanitize(input string) string {
	clean := _sanitizer.Sanitize(input)
	return clean
}
