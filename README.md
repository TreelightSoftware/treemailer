# TreeMailer

TreeMailer is a simple Go-based AWS Lambda tool powered by Serverless to provide a simple backend for form contact requests, such as in portfolios. This application sets up an AWS
( or your choice of FaaS provider) to take in some JSON, generate an email, and send it with MailGun.

## Setup

Global variables are used to setup the configuration. You can provide these either hard coded or you can set them through the environment.

* `to`
  * Where is the contact message going
  * `TREEMAILER_TO` is the environment variable

* `mgDomain`
  * The domain for Mailgun
  * `TREEMAILER_MG_DOMAIN` is the environment variable

* `mgKey`
  * The secret key for Mailgun
  * `TREEMAILER_MG_KE` is the environment variable

* `siteName`
  * A user-friendly sitename that is added to the subject
  * `TREEMAILER_SITE_NAME` is the environment variable

* `cc`
  * An additional `to` address to cc on all emails; useful for logging
  * `TREEMAILER_CC` is the environment variable

### Serverless

This application assumes you are using Serverless for your FaaS framework. You will want to init a new serverless directory and then create these files in that directory.

## Deploying

If you have a serverless setup, you can run the following:

`make && serverless deploy`

## Testing

We use the [github.com/stretchr/testify/assert](Assert) library for testing. We try to keep our test coverage high, but there is always room for improvement.

To run the tests, you can use `make test` or `make cover`.

## Contributing

We love contributors, especially people new to open source and looking to help! Documentation, comments, examples, and tutorials are great ways to contribute, even if you are not familiar with Go. Simply checkout the `CONTRIBUTING.md` file and open a pull request!