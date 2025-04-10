package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type UserService struct {
	// not need to implement
	NotEmptyStruct bool
}
type MessageService struct {
	// not need to implement
	NotEmptyStruct bool
}

type Container struct {
	constructors          map[string]interface{}
	singletonConstructors map[string]interface{}
	singletons            map[string]interface{}
}

func NewContainer() *Container {
	return &Container{
		constructors:          make(map[string]interface{}),
		singletonConstructors: make(map[string]interface{}),
		singletons:            make(map[string]interface{}),
	}
}

func (c *Container) RegisterType(name string, constructor interface{}) {
	if c.constructors[name] == nil {
		c.constructors[name] = constructor
	}
}

func (c *Container) RegisterSingletonType(name string, constructor interface{}) {
	if c.singletonConstructors[name] == nil {
		c.singletonConstructors[name] = constructor
	}
}

func (c *Container) Resolve(name string) (interface{}, error) {
	constructor := c.constructors[name]

	if constructor == nil {
		return nil, errors.New("no constructor registered")
	}

	return constructor.(func() interface{})(), nil
}

func (c *Container) ResolveSingleton(name string) (interface{}, error) {
	singletonConstructors := c.singletonConstructors[name]

	if singletonConstructors == nil {
		return nil, errors.New("no singleton constructor registered")
	}

	if c.singletons[name] == nil {
		c.singletons[name] = singletonConstructors.(func() interface{})()
	}

	return c.singletons[name], nil
}

func TestDIContainer(t *testing.T) {
	container := NewContainer()
	container.RegisterType("UserService", func() interface{} {
		return &UserService{}
	})
	container.RegisterType("MessageService", func() interface{} {
		return &MessageService{}
	})

	userService1, err := container.Resolve("UserService")
	assert.NoError(t, err)
	userService2, err := container.Resolve("UserService")
	assert.NoError(t, err)

	u1 := userService1.(*UserService)
	u2 := userService2.(*UserService)
	assert.False(t, u1 == u2)

	messageService, err := container.Resolve("MessageService")
	assert.NoError(t, err)
	assert.NotNil(t, messageService)

	paymentService, err := container.Resolve("PaymentService")
	assert.Error(t, err)
	assert.Nil(t, paymentService)

	container.RegisterSingletonType("UserService", func() interface{} {
		return &UserService{}
	})
	container.RegisterSingletonType("MessageService", func() interface{} {
		return &MessageService{}
	})

	userSingletonService1, err := container.ResolveSingleton("UserService")
	assert.NoError(t, err)
	userSingletonService2, err := container.ResolveSingleton("UserService")
	assert.NoError(t, err)

	us1 := userSingletonService1.(*UserService)
	us2 := userSingletonService2.(*UserService)
	assert.True(t, us1 == us2)

	messageSingletonService, err := container.ResolveSingleton("MessageService")
	assert.NoError(t, err)
	assert.NotNil(t, messageSingletonService)

	paymentSingletonService, err := container.ResolveSingleton("PaymentService")
	assert.Error(t, err)
	assert.Nil(t, paymentSingletonService)
}
