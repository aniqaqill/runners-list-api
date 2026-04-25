package service

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const testJWTSecret = "test_secret"

var _ = Describe("UserService", func() {
	var (
		userService *UserService
	)

	BeforeEach(func() {
		// jwtSecret is now injected at construction — no os.Setenv needed
		userService = &UserService{jwtSecret: testJWTSecret}
	})

	Describe("CreateToken", func() {
		It("should create a valid JWT token with the correct claims", func() {
			id := 1
			tokenString, err := userService.CreateToken(id)
			Expect(err).NotTo(HaveOccurred())

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return []byte(testJWTSecret), nil
			})
			Expect(err).NotTo(HaveOccurred())

			claims, ok := token.Claims.(jwt.MapClaims)
			Expect(ok).To(BeTrue())
			Expect(claims["id"]).To(Equal(float64(id)))
			Expect(claims["exp"]).To(BeNumerically(">", time.Now().Unix()))
		})

		It("should succeed even when no environment variable is set (secret comes from service field)", func() {
			// Demonstrates the improvement: secret is injected, not read from env
			svc := &UserService{jwtSecret: "any_injected_secret"}
			tokenString, err := svc.CreateToken(1)
			Expect(err).NotTo(HaveOccurred())
			Expect(tokenString).NotTo(BeEmpty())
		})

		It("should return an error if ID is not set or negative", func() {
			id := -1
			_, err := userService.CreateToken(id)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("ID cannot be negative or zero"))
		})
	})
})
