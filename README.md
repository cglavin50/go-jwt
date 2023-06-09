# Go-JWT

Learning about JSON Web Tokens and using them in Golang using Fiber, Gorm, Bcrypt, and jwt-go. Run with `$ docker compose up`, then validate endpoints at localhost:3000.

- localhost:3000/signup
  - Encode Email and Password as two key/value pairs into body
- locahost:3000/login
  - Same as above, will return a cookie with the JWT to be used
- locahost:3000/validate
  - Using the cookie provided from /login, will parse and return a response accepting or denying your login credentials to acecss private information

## Functionality

Using Go to host the backend of a MPA, created endpoints that take in user information, store credentials in a postgres DB, then create and return JWT's as cookies for users to use to access client information. Deploymeny via docker, [repo here](https://hub.docker.com/r/cglavin50/go-jwt).

## Dangers of JWT (or tokens in general)

Usable until they expire, meaning if someone gets access to them, they can easily impersonate you (so if you fall to a MitM attack or cross-site scripting, your information is not secure). In general, not great to use as a session token unless extra measures (including keeping a list of expired tokens) are taken. For example, do not put user emails on JWTs, as it is PII and can easily be mined via MitM (never store at rest).

JWT's can be a powerful tool, however require explicit and fine-tuned parameters (such as sub (subject), and exp (experiation)) to avoid misuse. As they are public, no PII can be held on it, so make sure to use either other user ID's, or a system like opaque tokens to validate following API calls.

## Opaque (Phantom Tokens)

To avoid issues with expiration time or PII on JWTs, [phantom/opaque tokens](https://curity.io/resources/learn/phantom-token-pattern/) can be used instead. Essentially, there is a key/value store of personal (value) tokens/JWTs, and referential (key) tokens, allowing a separation of personal information from the client. The control flow is as follows:

1. User sends login request. Middleware receives, and sends values to authentication service.
2. Authentication service validates the information, and creates a JWT/Phantom token pair. The phantom token is then passed back to the middleware.
3. Middleware sends phantom token to the user to use in future requests.
4. Upon receiving a future request, middleware then uses the authentication service to perform a lookup, and then passes the JWT to backend services to handle the request.