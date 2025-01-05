package middleware

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"os"
	"time"

	"github.com/aniqaqill/runners-list/internal/core/domain"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Middleware", func() {
	var app *fiber.App
	var ctrl *gomock.Controller

	BeforeEach(func() {
		app = fiber.New()
		ctrl = gomock.NewController(GinkgoT())
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("ValidateRegisterInput", func() {
		It("should return 400 for invalid JSON input", func() {
			app.Post("/register", ValidateRegisterInput, func(c *fiber.Ctx) error {
				return c.SendStatus(fiber.StatusOK)
			})

			req := httptest.NewRequest("POST", "/register", bytes.NewReader([]byte(`{"invalid_json`)))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(fiber.StatusBadRequest))
		})
		It("should return 400 for missing required fields", func() {
			app.Post("/register", ValidateRegisterInput, func(c *fiber.Ctx) error {
				return c.SendStatus(fiber.StatusOK)
			})

			req := httptest.NewRequest("POST", "/register", bytes.NewReader([]byte(`{"username": "testuser"}`)))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(fiber.StatusBadRequest))

			// Parse the response body to check the error message
			var response map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&response)
			Expect(err).NotTo(HaveOccurred())
			Expect(response["error"]).To(BeTrue())
			Expect(response["message"]).To(Equal("Validation failed"))
			// Expect(response["details"]).To(ContainSubstring("Password is a required field"))
		})

		It("should pass validation and store input in context", func() {
			user := domain.Users{
				Username: "testuser",
				Password: "testpassword",
			}
			body, _ := json.Marshal(user)

			app.Post("/register", ValidateRegisterInput, func(c *fiber.Ctx) error {
				input := c.Locals("registerInput").(domain.Users)
				Expect(input.Username).To(Equal("testuser"))
				Expect(input.Password).To(Equal("testpassword"))
				return c.SendStatus(fiber.StatusOK)
			})

			req := httptest.NewRequest("POST", "/register", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(fiber.StatusOK))
		})
	})

	Describe("ValidateLoginInput", func() {
		It("should return 400 for invalid JSON input", func() {
			app.Post("/login", ValidateLoginInput, func(c *fiber.Ctx) error {
				return c.SendStatus(fiber.StatusOK)
			})

			req := httptest.NewRequest("POST", "/login", bytes.NewReader([]byte(`{"invalid_json`)))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(fiber.StatusBadRequest))
		})

		It("should return 400 for missing required fields", func() {
			app.Post("/login", ValidateLoginInput, func(c *fiber.Ctx) error {
				return c.SendStatus(fiber.StatusOK)
			})

			req := httptest.NewRequest("POST", "/login", bytes.NewReader([]byte(`{"username": "testuser"}`)))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(fiber.StatusBadRequest))
		})

		It("should pass validation and store input in context", func() {
			user := domain.Users{
				Username: "testuser",
				Password: "testpassword",
			}
			body, _ := json.Marshal(user)

			app.Post("/login", ValidateLoginInput, func(c *fiber.Ctx) error {
				input := c.Locals("loginInput").(domain.Users)
				Expect(input.Username).To(Equal("testuser"))
				Expect(input.Password).To(Equal("testpassword"))
				return c.SendStatus(fiber.StatusOK)
			})

			req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(fiber.StatusOK))
		})
	})

	Describe("JWTProtected", func() {
		It("should return 401 for missing Authorization header", func() {
			app.Get("/protected", JWTProtected(), func(c *fiber.Ctx) error {
				return c.SendStatus(fiber.StatusOK)
			})

			req := httptest.NewRequest("GET", "/protected", nil)
			resp, err := app.Test(req)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(fiber.StatusUnauthorized))
		})

		It("should return 401 for invalid token", func() {
			app.Get("/protected", JWTProtected(), func(c *fiber.Ctx) error {
				return c.SendStatus(fiber.StatusOK)
			})

			req := httptest.NewRequest("GET", "/protected", nil)
			req.Header.Set("Authorization", "Bearer invalidtoken")
			resp, err := app.Test(req)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(fiber.StatusUnauthorized))
		})

		It("should return 401 for expired token", func() {
			os.Setenv("JWT_SECRET", "supersecretkey")
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"sub": "1234567890",
				"exp": time.Now().Add(-time.Hour).Unix(), // Expired token
			})
			tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
			Expect(err).NotTo(HaveOccurred())

			app.Get("/protected", JWTProtected(), func(c *fiber.Ctx) error {
				return c.SendStatus(fiber.StatusOK)
			})

			req := httptest.NewRequest("GET", "/protected", nil)
			req.Header.Set("Authorization", "Bearer "+tokenString)
			resp, err := app.Test(req)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(fiber.StatusUnauthorized))
		})

		It("should allow access with a valid token", func() {
			os.Setenv("JWT_SECRET", "supersecretkey")
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"sub": "1234567890",
				"exp": time.Now().Add(time.Hour * 1).Unix(),
			})
			tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
			Expect(err).NotTo(HaveOccurred())

			app.Get("/protected", JWTProtected(), func(c *fiber.Ctx) error {
				return c.SendStatus(fiber.StatusOK)
			})

			req := httptest.NewRequest("GET", "/protected", nil)
			req.Header.Set("Authorization", "Bearer "+tokenString)
			resp, err := app.Test(req)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(fiber.StatusOK))
		})

	})
})
