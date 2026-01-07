package service

import (
	"errors"
	"time"

	"github.com/aniqaqill/runners-list/internal/core/domain"
	"github.com/aniqaqill/runners-list/internal/port/mocks"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("EventService", func() {
	var (
		ctrl         *gomock.Controller
		mockRepo     *mocks.MockEventRepository
		eventService *EventService
	)

	BeforeEach(func() {
		// Initialize Gomock controller
		ctrl = gomock.NewController(GinkgoT())

		// Initialize mock repository
		mockRepo = mocks.NewMockEventRepository(ctrl)

		// Initialize EventService with the mock repository
		eventService = NewEventService(mockRepo)
	})

	AfterEach(func() {
		// Clean up Gomock controller
		ctrl.Finish()
	})

	Describe("CreateEvent", func() {
		It("should return an error if the event date is in the past", func() {
			pastEvent := &domain.Events{
				Name: "Past Event",
				Date: time.Now().AddDate(-1, 0, 0), // Past date
			}

			err := eventService.CreateEvent(pastEvent)
			Expect(err).To(MatchError(ErrEventDateInPast))
		})

		It("should return an error if the event name is not unique", func() {
			futureEvent := &domain.Events{
				Name: "Future Event",
				Date: time.Now().AddDate(1, 0, 0), // Future date
			}

			// Mock the EventNameExists method to return true (name exists)
			mockRepo.EXPECT().
				EventNameExists(futureEvent.Name).
				Return(true)

			err := eventService.CreateEvent(futureEvent)
			Expect(err).To(MatchError(ErrEventNameNotUnique))
		})

		It("should create the event if the date is in the future and the name is unique", func() {
			futureEvent := &domain.Events{
				Name: "Future Event",
				Date: time.Now().AddDate(1, 0, 0), // Future date
			}

			// Mock the EventNameExists method to return false (name is unique)
			mockRepo.EXPECT().
				EventNameExists(futureEvent.Name).
				Return(false)

			// Mock the Create method to return nil (success)
			mockRepo.EXPECT().
				Create(futureEvent).
				Return(nil)

			err := eventService.CreateEvent(futureEvent)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("ListEvents", func() {
		It("should return a list of events", func() {
			events := []domain.Events{
				{Name: "Event 1", Date: time.Now().AddDate(1, 0, 0)},
				{Name: "Event 2", Date: time.Now().AddDate(1, 1, 0)},
			}

			// Mock the FindAll method to return the events
			mockRepo.EXPECT().
				FindAll().
				Return(events, nil)

			result, err := eventService.ListEvents()
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(events))
		})

		It("should return an error if the repository fails", func() {
			// Mock the FindAll method to return an error
			mockRepo.EXPECT().
				FindAll().
				Return(nil, errors.New("repository error"))

			_, err := eventService.ListEvents()
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("DeleteEvent", func() {
		It("should delete the event if it exists", func() {
			eventID := uint(1)
			event := &domain.Events{Name: "Event 1", Date: time.Now().AddDate(1, 0, 0)}

			// Mock the FindByID method to return the event
			mockRepo.EXPECT().
				FindByID(eventID).
				Return(event, nil)

			// Mock the Delete method to return nil (success)
			mockRepo.EXPECT().
				Delete(event).
				Return(nil)

			err := eventService.DeleteEvent(eventID)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should return an error if the event does not exist", func() {
			eventID := uint(1)

			// Mock the FindByID method to return an error
			mockRepo.EXPECT().
				FindByID(eventID).
				Return(nil, errors.New("event not found"))

			err := eventService.DeleteEvent(eventID)
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("isEventDateInFuture", func() {
		It("should return true if the event date is in the future", func() {
			futureDate := time.Now().AddDate(1, 0, 0)
			Expect(isEventDateInFuture(futureDate)).To(BeTrue())
		})

		It("should return false if the event date is in the past", func() {
			pastDate := time.Now().AddDate(-1, 0, 0)
			Expect(isEventDateInFuture(pastDate)).To(BeFalse())
		})
	})
})
