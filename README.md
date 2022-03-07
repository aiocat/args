<div align="center">

![Logo](/static/img/logo.png)

# Args

Anonymous argument sharing and advocacy platform

</div>

## Technologies

- **Programming Language**: Go
- **Database**: MongoDB
- **Server**: Gofiber
- **Captcha Service**: HCaptcha
- **Front-end**: Go's Official HTML Template Engine, CSS & JS
- **Hosting**: Heroku

## Hosting

Create a `.env` file and add these key-values:

- **MONGO_URI**: MongoDB database connection url.
- **WEBHOOK_URI**: Discord webhook url.
- **CAPTCHA_SECRET**: HCaptcha secret key.
- **PORT**: Port to serve.

Install Go programming language and run with: `go run .`

View demo here: https://args-app.herokuapp.com

## Routes

- `GET /`
- `GET /new`
- `GET /saves`
- `GET /arguments/:id`
- `POST /arguments`
- `POST /arguments/:id`
- `GET /delete`
- `DELETE /arguments/:secret`
- `GET /reports/:id`
- `POST /reports`

## License

Args is distributed under AGPLv3 license. for more information:

- https://raw.githubusercontent.com/aiocat/args/main/LICENSE
