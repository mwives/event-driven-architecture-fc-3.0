package events

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TestEvent struct {
	Name    string
	Payload interface{}
}

func (e *TestEvent) GetName() string {
	return e.Name
}

func (e *TestEvent) GetDateTime() time.Time {
	return time.Now()
}

func (e *TestEvent) GetPayload() interface{} {
	return e.Payload
}

func (e *TestEvent) SetPayload(payload interface{}) {
	e.Payload = payload
}

type TestEventHandler struct {
	ID int
}

// Function to implement the EventHandler interface
func (h *TestEventHandler) Handle(event EventInterface, wg *sync.WaitGroup) {}

// EventDispatcherTestSuite is the test suite structure used by testify/suite.
type EventDispatcherTestSuite struct {
	suite.Suite
	event           TestEvent        // First test event
	event2          TestEvent        // Second test event
	handler         TestEventHandler // First handler for event processing
	handler2        TestEventHandler // Second handler
	handler3        TestEventHandler // Third handler
	eventDispatcher *EventDispatcher // Event dispatcher instance
}

// SetupTest is a method that runs before each test to initialize the environment.
func (suite *EventDispatcherTestSuite) SetupTest() {
	// Create a new EventDispatcher instance for each test
	suite.eventDispatcher = NewEventDispatcher()
	// Initialize handlers and events
	suite.handler = TestEventHandler{ID: 1}
	suite.handler2 = TestEventHandler{ID: 2}
	suite.handler3 = TestEventHandler{ID: 3}
	suite.event = TestEvent{Name: "test", Payload: "test"}
	suite.event2 = TestEvent{Name: "test2", Payload: "test2"}
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Register() {
	err := suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)
	suite.Nil(err)
	suite.Equal(len(suite.eventDispatcher.handlers[suite.event.GetName()]), 1) // Check that one handler was registered

	err = suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler2)
	suite.Nil(err)
	suite.Equal(len(suite.eventDispatcher.handlers[suite.event.GetName()]), 2) // Now two handlers should be registered

	// Check that the correct handlers were registered in the correct order
	suite.Equal(&suite.handler, suite.eventDispatcher.handlers[suite.event.GetName()][0])
	suite.Equal(&suite.handler2, suite.eventDispatcher.handlers[suite.event.GetName()][1])
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Register_WithSameHandler() {
	err := suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)
	suite.Nil(err)
	suite.Equal(len(suite.eventDispatcher.handlers[suite.event.GetName()]), 1)

	err = suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)
	suite.Error(err, ErrHandlerAlreadyRegistered)                              // Assert that the appropriate error was returned
	suite.Equal(len(suite.eventDispatcher.handlers[suite.event.GetName()]), 1) // No duplicate handlers should be added
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Clear() {
	// Register handlers for both events
	err := suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)
	suite.Nil(err)
	suite.Equal(len(suite.eventDispatcher.handlers[suite.event.GetName()]), 1)

	err = suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler2)
	suite.Nil(err)
	suite.Equal(len(suite.eventDispatcher.handlers[suite.event.GetName()]), 2)

	err = suite.eventDispatcher.Register(suite.event2.GetName(), &suite.handler3)
	suite.Nil(err)
	suite.Equal(len(suite.eventDispatcher.handlers[suite.event2.GetName()]), 1)

	// Clear all registered handlers
	suite.eventDispatcher.Clear()
	suite.Equal(len(suite.eventDispatcher.handlers), 0) // No handlers should be left
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Has() {
	// Register two handlers for the event
	err := suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)
	suite.Nil(err)
	suite.Equal(len(suite.eventDispatcher.handlers[suite.event.GetName()]), 1)

	err = suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler2)
	suite.Nil(err)
	suite.Equal(len(suite.eventDispatcher.handlers[suite.event.GetName()]), 2)

	// Assert that the dispatcher has both registered handlers
	suite.True(suite.eventDispatcher.Has(suite.event.GetName(), &suite.handler))
	suite.True(suite.eventDispatcher.Has(suite.event.GetName(), &suite.handler2))
	// Assert that it does not have the third handler
	suite.False(suite.eventDispatcher.Has(suite.event.GetName(), &suite.handler3))
}

// MockHandler is a mock implementation of the EventHandler for testing.
// It uses testify's Mock functionality to simulate and track method calls.
type MockHandler struct {
	mock.Mock
}

// Handle is the mocked method that will simulate handling an event.
func (m *MockHandler) Handle(event EventInterface, wg *sync.WaitGroup) {
	m.Called(event) // Track that this method was called with a specific event
	wg.Done()
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Dispatch() {
	eh := &MockHandler{}
	eh.On("Handle", &suite.event) // Expect that the Handle method will be called with the event

	eh2 := &MockHandler{}
	eh2.On("Handle", &suite.event) // Expect that the Handle method will be called with the event

	// Register the mock handler to the event
	suite.eventDispatcher.Register(suite.event.GetName(), eh)
	suite.eventDispatcher.Register(suite.event.GetName(), eh2)
	// Dispatch the event, which should trigger the Handle method on the mock handler
	suite.eventDispatcher.Dispatch(&suite.event)

	// Assert that all expectations for the mock handlers were met (i.e., Handle was called)
	eh.AssertExpectations(suite.T())
	eh2.AssertExpectations(suite.T())
	// Assert that the Handle method was called exactly once
	eh.AssertNumberOfCalls(suite.T(), "Handle", 1)
	eh2.AssertNumberOfCalls(suite.T(), "Handle", 1)
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Remove() {
	// Register handlers for both events
	err := suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)
	suite.Nil(err)
	suite.Equal(len(suite.eventDispatcher.handlers[suite.event.GetName()]), 1)

	err = suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler2)
	suite.Nil(err)
	suite.Equal(len(suite.eventDispatcher.handlers[suite.event.GetName()]), 2)

	err = suite.eventDispatcher.Register(suite.event2.GetName(), &suite.handler3)
	suite.Nil(err)
	suite.Equal(len(suite.eventDispatcher.handlers[suite.event2.GetName()]), 1)

	// Remove the first handler for the first event
	suite.eventDispatcher.Remove(suite.event.GetName(), &suite.handler)
	suite.Equal(len(suite.eventDispatcher.handlers[suite.event.GetName()]), 1)
	// Assert that the remaining handler is handler2
	suite.Equal(&suite.handler2, suite.eventDispatcher.handlers[suite.event.GetName()][0])

	// Remove the second handler
	suite.eventDispatcher.Remove(suite.event.GetName(), &suite.handler2)
	suite.Equal(len(suite.eventDispatcher.handlers[suite.event.GetName()]), 0)

	// Remove the handler for the second event
	suite.eventDispatcher.Remove(suite.event2.GetName(), &suite.handler3)
	suite.Equal(len(suite.eventDispatcher.handlers[suite.event2.GetName()]), 0)
}

// TestSuite runs the test suite using testify's suite runner.
func TestSuite(t *testing.T) {
	suite.Run(t, new(EventDispatcherTestSuite))
}
