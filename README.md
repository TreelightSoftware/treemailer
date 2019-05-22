# TreeMailer

TreeMailer is a simple Go-based AWS Lambda tool powered by Serverless to provide a simple backend for form contact requests, such as in portfolios. This application sets up an AWS
( or your choice of FaaS provider) to take in some JSON, generate an email, and send it with MailGun.

## Setup

As of this moment, you **must** hardcode your keys in the `mailer/main.go` file. You need to provide the following:

* `to`: Where is the contact going

* `mgDomain` The domain for Mailgun

* `mgKey`: The secret key for Mailgun

* `siteName`: A user-friendly sitename that is added to the subject