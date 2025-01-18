package service

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("UserService", func() {
	var (
		userService *UserService
	)

	BeforeEach(func() {
		// Initialize UserService
		userService = &UserService{}
	})

	Describe("CreateToken", func() {
		It("should create a valid JWT token with the correct claims", func() {
			// Set the JWT_SECRET environment variable for testing
			os.Setenv("JWT_SECRET", "test_secret")
			defer os.Unsetenv("JWT_SECRET")

			id := 1
			tokenString, err := userService.CreateToken(id)
			Expect(err).NotTo(HaveOccurred())

			// Parse the token to verify its claims
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("JWT_SECRET")), nil
			})
			Expect(err).NotTo(HaveOccurred())

			// Verify the claims
			claims, ok := token.Claims.(jwt.MapClaims)
			Expect(ok).To(BeTrue())
			Expect(claims["id"]).To(Equal(float64(id))) // JWT claims are float64
			Expect(claims["exp"]).To(BeNumerically(">", time.Now().Unix()))
		})

		It("should return an error if JWT_SECRET is not set", func() {
			// Ensure JWT_SECRET is not set
			os.Unsetenv("JWT_SECRET")

			id := 1
			_, err := userService.CreateToken(id)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("JWT_SECRET environment variable is not set"))
		})

		It("should return an error if ID is not set or negative", func() {
			// Set the JWT_SECRET environment variable for testing
			os.Setenv("JWT_SECRET", "test_secret")
			defer os.Unsetenv("JWT_SECRET")

			id := -1
			_, err := userService.CreateToken(id)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("ID cannot be negative or zero"))
		})
	})
})
