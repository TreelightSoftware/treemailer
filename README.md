# TreeMailer

TreeMailer is a simple Go-based AWS Lambda tool powered by Serverless to provide a simple backend for form contact requests, such as in portfolios. This application sets up an AWS
( or your choice of FaaS provider) to take in some JSON, generate an email, and send it with MailGun.

## Setup

As of this moment, you **must** hardcode your keys in the `mailer/main.go` file. You need to provide the following:

* `to`: Where is the contact going

* `mgDomain` The domain for Mailgun

* `mgKey`: The secret key for Mailgun

* `siteName`: A user-friendly sitename that is added to the subject

* `cc`: An additional `to` address to cc on all emails; useful for logging

### Serverless

This application assumes you are using Serverless for your FaaS framework. You will want to init a new serverless directory and then create these files in that directory.

## Deploying

If you have a serverless setup, you can run the following:

`make && serverless deploy`