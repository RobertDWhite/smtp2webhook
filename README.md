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
      - SMTP2WEBHOOK_DOMAIN=exampledomain.com
      - SMTP2WEBHOOK_CODE=webhook
      - SMTP2WEBHOOK_URL_IDENTIFIER=https://webhook.white.fm/IDENTIFIER
```

Will forward mail to `webhook+IDENTIFIER@exampledomain.com` to `https://webhook.white.fm/IDENTIFIER`.
