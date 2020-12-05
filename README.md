# TreeMailer

[![CircleCI](https://circleci.com/gh/TreelightSoftware/treemailer/tree/master.svg?style=svg)](https://circleci.com/gh/TreelightSoftware/treemailer/tree/master)

[![codecov](https://codecov.io/gh/TreelightSoftware/treemailer/branch/master/graph/badge.svg)](https://codecov.io/gh/TreelightSoftware/treemailer)

TreeMailer is a simple Go-based AWS Lambda tool powered by Serverless to provide a simple backend for form contact requests, such as in portfolios. This application sets up an AWS
( or your choice of FaaS provider) to take in some JSON, generate an email, and send it with MailGun.

## Getting Started

You must have serverless >= 2 installed. See [https://www.serverless.com/framework/docs/getting-started/](Serverless - Getting Started).

You have two options.

1) If you have an existing serverless install, you can grab the `mailer` directory and drop it in there. Then, update the Makefile as appropriate.

2) You can checkout this directory and make your changes in here. This directory lacks the .serverless configuration, so you will need to configure serverless as you see fit.

You can specify the values in the `serverless.yml` file or inject the environment variables into the Lambda. If you use the `serverless.yml` file, you can also plug them in from another resources using [https://www.serverless.com/blog/serverless-v1.2.0](Serverless Variables). Ideally, you would use something like `git-crypt` to protect these as well.

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
  * `TREEMAILER_MG_KEY` is the environment variable

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

### Note on Security and SPAM

Currently, we do not require a CAPTCHA, although that feature is planned. Another alternative or addition is input field checking.

Many automated tools search the form for fields conveniently containing `name` or `id` or `email` or similar in the `input` tags. So we recommend you include a CSS-hidden form field with the name and id of email, with the real email address being stored in a field with a non-obvious name (probably best to not have the world email in there). When the button is clicked and the script prepares to send the JSON, check if there is a value in the hidden field. If there is, it's usually safe to assume it was filled in with a script (most users likely won't unhide a field in CSS and fill it in). Of course, other mitigation techniques are welcome, and if you have other best practices, feel free to open a PR with them in the README or another file.

## Testing

We use the [Assert](https://github.com/stretchr/testify) library for testing. We try to keep our test coverage high, but there is always room for improvement.

To run the tests, you can use `make test` or `make cover`.

## Contributing

We love contributors, especially people new to open source and looking to help! Documentation, comments, examples, and tutorials are great ways to contribute, even if you are not familiar with Go. Simply checkout the `CONTRIBUTING.md` file and open a pull request!
