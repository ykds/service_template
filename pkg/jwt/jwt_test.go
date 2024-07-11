package jwt

import (
	"fmt"
	"testing"
)

func TestJwt(t *testing.T) {
	i, err := ParseToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjAxNzMxNjIsIm5iZiI6MTcxOTU2ODM2MiwidXNlcl9pZCI6Mn0.UayA4SugX1Q-e1ur6PDWUNhNqKLTLdkneQ2JVrgQJ1M")
	if err != nil {
		panic(err)
	}
	fmt.Println(i)
}
