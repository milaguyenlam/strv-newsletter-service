# strv-newsletter-service

## Endpoints
- POST /user/signup -> JWT Token
- POST /user/login -> JWT Token
- POST /newsletter/create (JWT auth)
- POST /newsletter/send (JWT auth)
- POST /newsletter/subscribe
- POST /newsletter/unsubscribe
- GET /newsletter/list (??)
  - Doesn't need to be implemented - set read rule to true to make Firestore public
  - Implement a simple client app that connects to Firestore and lists all subscribed newsletters for specified email address