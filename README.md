# Mercurity

A simple email server that allows you to create inboxes and rules to filter emails in go.

A project for a PI planning.



## Running the server

```shell
go run main.go
```

## Testing API


Create an Account

```shell
curl -X POST http://localhost:8080/accounts -H "Content-Type: application/json" -d '{"name": "Test Account"}'
```

Create an Inbox

```shell
curl -X POST http://localhost:8080/accounts/1/inboxes -H "Content-Type: application/json" -d '{"email": "inbox@example.com"}'
```

Create Rules

```shell
curl -X POST http://localhost:8080/accounts/1/inboxes/1/rules -H "Content-Type: application/json" -d '{"sender": "sender@example.com", "receiver": "inbox@example.com", "subject": "Test Subject"}'
```

Send an email

```shell
swaks --to inbox@example.com --from sender@example.com --server localhost:1025 --header "Subject: Test Subject" --body "This is a test email."
```

## Architecture

- https://archiveopteryx.org/db/
