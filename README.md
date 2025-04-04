POST /users - create user
GET /user/<id> - get user
PATCH /user/<id> - edit user
 
type User struct {
  ID uuid
  Firstname string
  Lastname string
  Email string
  Age uint
  Created time.Time
}
