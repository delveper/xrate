// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mock

import (
	"context"
	"github.com/GenesisEducationKyiv/main-project-delveper/internal/subs"
	"sync"
)

// Ensure, that SubscriberRepositoryMock does implement subs.SubscriberRepository.
// If this is not the case, regenerate this file with moq.
var _ subs.SubscriberRepository = &SubscriberRepositoryMock{}

// SubscriberRepositoryMock is a mock implementation of subs.SubscriberRepository.
//
//	func TestSomethingThatUsesSubscriberRepository(t *testing.T) {
//
//		// make and configure a mocked subs.SubscriberRepository
//		mockedSubscriberRepository := &SubscriberRepositoryMock{
//			AddFunc: func(contextMoqParam context.Context, subscription subs.Subscription) error {
//				panic("mock out the Add method")
//			},
//			ListFunc: func(contextMoqParam context.Context) ([]subs.Subscription, error) {
//				panic("mock out the List method")
//			},
//		}
//
//		// use mockedSubscriberRepository in code that requires subs.SubscriberRepository
//		// and then make assertions.
//
//	}
type SubscriberRepositoryMock struct {
	// AddFunc mocks the Add method.
	AddFunc func(contextMoqParam context.Context, subscription subs.Subscription) error

	// ListFunc mocks the List method.
	ListFunc func(contextMoqParam context.Context) ([]subs.Subscription, error)

	// calls tracks calls to the methods.
	calls struct {
		// Add holds details about calls to the Add method.
		Add []struct {
			// ContextMoqParam is the contextMoqParam argument value.
			ContextMoqParam context.Context
			// Subscription is the subscription argument value.
			Subscription subs.Subscription
		}
		// List holds details about calls to the List method.
		List []struct {
			// ContextMoqParam is the contextMoqParam argument value.
			ContextMoqParam context.Context
		}
	}
	lockAdd  sync.RWMutex
	lockList sync.RWMutex
}

// Add calls AddFunc.
func (mock *SubscriberRepositoryMock) Add(contextMoqParam context.Context, subscription subs.Subscription) error {
	if mock.AddFunc == nil {
		panic("SubscriberRepositoryMock.AddFunc: method is nil but SubscriberRepository.Add was just called")
	}
	callInfo := struct {
		ContextMoqParam context.Context
		Subscription    subs.Subscription
	}{
		ContextMoqParam: contextMoqParam,
		Subscription:    subscription,
	}
	mock.lockAdd.Lock()
	mock.calls.Add = append(mock.calls.Add, callInfo)
	mock.lockAdd.Unlock()
	return mock.AddFunc(contextMoqParam, subscription)
}

// AddCalls gets all the calls that were made to Add.
// Check the length with:
//
//	len(mockedSubscriberRepository.AddCalls())
func (mock *SubscriberRepositoryMock) AddCalls() []struct {
	ContextMoqParam context.Context
	Subscription    subs.Subscription
} {
	var calls []struct {
		ContextMoqParam context.Context
		Subscription    subs.Subscription
	}
	mock.lockAdd.RLock()
	calls = mock.calls.Add
	mock.lockAdd.RUnlock()
	return calls
}

// List calls ListFunc.
func (mock *SubscriberRepositoryMock) List(contextMoqParam context.Context) ([]subs.Subscription, error) {
	if mock.ListFunc == nil {
		panic("SubscriberRepositoryMock.ListFunc: method is nil but SubscriberRepository.List was just called")
	}
	callInfo := struct {
		ContextMoqParam context.Context
	}{
		ContextMoqParam: contextMoqParam,
	}
	mock.lockList.Lock()
	mock.calls.List = append(mock.calls.List, callInfo)
	mock.lockList.Unlock()
	return mock.ListFunc(contextMoqParam)
}

// ListCalls gets all the calls that were made to List.
// Check the length with:
//
//	len(mockedSubscriberRepository.ListCalls())
func (mock *SubscriberRepositoryMock) ListCalls() []struct {
	ContextMoqParam context.Context
} {
	var calls []struct {
		ContextMoqParam context.Context
	}
	mock.lockList.RLock()
	calls = mock.calls.List
	mock.lockList.RUnlock()
	return calls
}