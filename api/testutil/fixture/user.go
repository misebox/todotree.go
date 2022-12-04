package fixture

import (
	"strconv"
	"time"
	"todotree/entity"
)

var cnt int = 10000

func NewUserForTest() *entity.User {
	cnt += 1
	test_digit := cnt
	// test_digit := rand.Int()
	test_prefix := strconv.Itoa(test_digit)[:5]
	result := &entity.User{
		ID:       entity.UserID(test_digit),
		Name:     "testuser" + test_prefix,
		Email:    "testuser" + test_prefix + "@example.com",
		Password: "testpassword",
		Role:     "admin",
		Created:  time.Now(),
		Modified: time.Now(),
	}
	return result
}
