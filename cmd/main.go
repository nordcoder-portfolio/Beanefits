package main

import (
	"Beanefits/internal/infra/jwtverifier"
	"context"
)

func main() {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlcyI6WyJDQVNISUVSIl0sImlzcyI6ImJlYW5lZml0cyIsInN1YiI6IjIiLCJleHAiOjE3NjY0ODEwMzYsImlhdCI6MTc2NjM5NDYzNn0.IKuOTkK1pKlIc3ap9dgYXJ5TGWw-fyQbWLw4v47y4hQ"
	v := jwtverifier.New("dev-secret-change-me", "beanefits")
	_, err := v.Verify(context.Background(), token)
	if err != nil {
		panic(err)
	}
}
