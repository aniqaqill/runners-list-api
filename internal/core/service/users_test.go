package service

import (
	"errors"

	"github.com/aniqaqill/runners-list/internal/core/domain"
	"github.com/aniqaqill/runners-list/internal/port/mocks"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"golang.org/x/crypto/bcrypt"
)

var _ = Describe("UserService", func() {
	var (
		ctrl        *gomock.Controller
		mockRepo    *mocks.MockUserRepository
		userService *UserService
	)

	BeforeEach(func() {
		// Initialize Gomock controller
		ctrl = gomock.NewController(GinkgoT())

		// Initialize mock repository
		mockRepo = mocks.NewMockUserRepository(ctrl)

		// Initialize UserService with the mock repository
		userService = NewUserService(mockRepo)
	})

	AfterEach(func() {
		// Clean up Gomock controller
		ctrl.Finish()
	})

	Describe("Register", func() {
		It("should return an error if the username already exists", func() {
			username := "existing_user"
			password := "password123"

			// Mock the FindByUsername method to return an existing user
			mockRepo.EXPECT().
				FindByUsername(username).
				Return(&domain.Users{Username: username}, nil)

			err := userService.Register(username, password)
			Expect(err).To(MatchError(ErrUsernameAlreadyExists))
		})

		It("should create the user if the username is unique", func() {
			username := "new_user"
			password := "password123"

			// Mock the FindByUsername method to return nil (user does not exist)
			mockRepo.EXPECT().
				FindByUsername(username).
				Return(nil, nil)

			// Mock the Create method to return nil (success)
			mockRepo.EXPECT().
				Create(gomock.Any()).
				Do(func(user *domain.Users) {
					// Verify that the password is hashed
					err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
					Expect(err).NotTo(HaveOccurred())
				}).
				Return(nil)

			err := userService.Register(username, password)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("GetUserByUsername", func() {
		It("should return the user if the username exists", func() {
			username := "existing_user"
			user := &domain.Users{Username: username}

			// Mock the FindByUsername method to return the user
			mockRepo.EXPECT().
				FindByUsername(username).
				Return(user, nil)

			result, err := userService.GetUserByUsername(username)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(user))
		})

		It("should return an error if the username does not exist", func() {
			username := "non_existing_user"

			// Mock the FindByUsername method to return an error
			mockRepo.EXPECT().
				FindByUsername(username).
				Return(nil, errors.New("user not found"))

			_, err := userService.GetUserByUsername(username)
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("Login", func() {
		It("should return an error if the username does not exist", func() {
			username := "non_existing_user"
			password := "password123"

			// Mock the FindByUsername method to return an error
			mockRepo.EXPECT().
				FindByUsername(username).
				Return(nil, errors.New("user not found"))

			_, err := userService.Login(username, password)
			Expect(err).To(MatchError(ErrInvalidCredentials))
		})

		It("should return an error if the password is incorrect", func() {
			username := "existing_user"
			password := "wrong_password"
			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correct_password"), 14)
			user := &domain.Users{
				Username: username,
				Password: string(hashedPassword),
			}

			// Mock the FindByUsername method to return the user
			mockRepo.EXPECT().
				FindByUsername(username).
				Return(user, nil)

			_, err := userService.Login(username, password)
			Expect(err).To(MatchError(ErrInvalidCredentials))
		})

		It("should return the user if the credentials are correct", func() {
			username := "existing_user"
			password := "correct_password"
			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
			user := &domain.Users{
				Username: username,
				Password: string(hashedPassword),
			}

			// Mock the FindByUsername method to return the user
			mockRepo.EXPECT().
				FindByUsername(username).
				Return(user, nil)

			result, err := userService.Login(username, password)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(user))
		})
	})

	Describe("ListUsers", func() {
		It("should return all users", func() {
			users := []domain.Users{
				{Username: "user1"},
				{Username: "user2"},
			}

			// Mock the FindAll method to return the list of users
			mockRepo.EXPECT().
				FindAll().
				Return(users, nil)

			result, err := userService.ListUsers()
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(users))
		})

		It("should return an error if there is an issue retrieving users", func() {
			// Mock the FindAll method to return an error
			mockRepo.EXPECT().
				FindAll().
				Return(nil, errors.New("failed to retrieve users"))

			_, err := userService.ListUsers()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("failed to retrieve users"))
		})
	})
})
