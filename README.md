# Go-JWT

Learning about JSON Web Tokens and using them in Golang using Fiber, Gorm, Bcrypt, and jwt-go.

# DB
- Using a Docker postgres image to run a basic DB
- Starting with Docker compose to have this entire project be containerized and deployable

# Dangers of JWT (or tokens in general)

Usable until they expire, meaning if someone gets access to them, they can easily impersonate you (so if you fall to a MitM attack or cross-site scripting, your information is not secure). In general, not great to use as a session token unless extra measures (including keeping a list of expired tokens) are taken.
