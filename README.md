# smtp2webhook

SMTP to Webhook Relay.

```yml
version: "3"
services:
  smtp2webhook:
    restart: always
    image: registry.white.fm/smtp2webhook:latest
    ports:
      - "25:25"
    environment:
      - WEBHOOK_URL=https://EXAMPLE.COM
```

Will forward all mail to the specified Webhook.
