package support

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/mailgun/mailgun-go/v4"
)

// Your available domain names can be found here:
// (https://app.mailgun.com/app/domains)
// var yourDomain string = "www.chantosweb.com" // e.g. mg.yourcompany.com
var yourDomain string = "https://api.mailgun.net/v3/sandbox98fb8bf69e90476faee0e0ccec15972d.mailgun.org" // e.g. mg.yourcompany.com

// You can find the Private API Key in your Account Menu, under "Settings":
// (https://app.mailgun.com/app/account/security)
var privateAPIKey string = "60969f97d9e5eac5f91ae07278d7350e-ea44b6dc-fc0b9bfe"


func email() {
    // Create an instance of the Mailgun Client
    mg := mailgun.NewMailgun(yourDomain, privateAPIKey)

    sender := "me@sandbox98fb8bf69e90476faee0e0ccec15972d.mailgun.org"
    subject := "Fancy subject!"
    body := "Hello from Mailgun Go!"
    recipient := "myrachanto@gmail.com"

    // The message object allows you to add attachments and Bcc recipients
    message := mg.NewMessage(sender, subject, body, recipient)

    ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
    defer cancel()

    // Send the message with a 10 second timeout
    resp, id, err := mg.Send(ctx, message)

    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("ID: %s Resp: %s\n", id, resp)
}