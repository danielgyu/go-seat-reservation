# Seat-reservation backend in Go
For a detailed (in korean) description of this project, visit <https://bit.ly/3C3eyWC>

This is a simple project for me to learn usage of Go for building backend services. Use `go run cmd/reservation/main.go` to test the server. Database can be set up in `pkg/repository/data.sql`.

Here are some endpoints,
- GET /halls get all registered halls
- POST /halls register a hall (for admin)
- GET /reservation/:hallName reserve a seat in hall

The motivation for me to try out a go project was due to the fact that i had to implement a backend server specifically in Go at work. This project gave me a good foundation on the following.
- Project structure
- Database interaction for mysql and redis
- Depenency injection with services
- Utilizing channels and goroutines
- Overall understanding of the net/http package

Room for improvements exist, such as refactoring some code and better packaging them and some of them would probably be implemented in my next project that touches on gRPC with Go.
