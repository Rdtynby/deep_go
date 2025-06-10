package main

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type Group struct {
	cancel           context.CancelCauseFunc
	ctx              context.Context
	wg               sync.WaitGroup
	condition        *sync.Cond
	workersLimit     int
	workersNumber    int
	maxWorkersNumber int // Only for tests
}

func NewErrGroup(ctx context.Context) (*Group, context.Context) {
	newCtx, cancel := context.WithCancelCause(ctx)
	mutex := &sync.Mutex{}

	return &Group{
		ctx:       newCtx,
		cancel:    cancel,
		condition: sync.NewCond(mutex),
	}, newCtx
}

func (g *Group) Go(action func() error) {
	g.condition.L.Lock()
	defer g.condition.L.Unlock()

	if g.workersLimit > 0 && g.workersNumber >= g.workersLimit {
		g.condition.Wait()
	}

	g.workersNumber++
	g.wg.Add(1)

	go func() {
		defer func() {
			if g.workersNumber > g.maxWorkersNumber {
				g.maxWorkersNumber = g.workersNumber
			}

			g.workersNumber--
			g.condition.Signal()
			g.wg.Done()
		}()

		err := action()

		if err != nil {
			g.cancel(err)
		}
	}()
}

func (g *Group) Wait() error {
	g.wg.Wait()

	return context.Cause(g.ctx)
}

func (g *Group) SetLimit(limit int) {
	g.workersLimit = limit
}

func TestErrGroupWithoutError(t *testing.T) {
	var counter atomic.Int32
	group, _ := NewErrGroup(context.Background())

	for i := 0; i < 5; i++ {
		group.Go(func() error {
			time.Sleep(time.Second)
			counter.Add(1)
			return nil
		})
	}

	err := group.Wait()
	assert.Equal(t, int32(5), counter.Load())
	assert.NoError(t, err)
	assert.Equal(t, 5, group.maxWorkersNumber)
}

func TestErrGroupWithError(t *testing.T) {
	var counter atomic.Int32
	group, ctx := NewErrGroup(context.Background())

	for i := 0; i < 5; i++ {
		group.Go(func() error {
			timer := time.NewTimer(time.Second)
			defer timer.Stop()

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-timer.C:
				counter.Add(1)
				return nil
			}
		})
	}

	group.Go(func() error {
		return errors.New("error")
	})

	err := group.Wait()
	assert.Equal(t, int32(0), counter.Load())
	assert.Error(t, err)
	assert.Equal(t, 6, group.maxWorkersNumber)
}

func TestErrGroupWithoutErrorLimit(t *testing.T) {
	var counter atomic.Int32
	group, _ := NewErrGroup(context.Background())
	group.SetLimit(3)

	for i := 0; i < 5; i++ {
		group.Go(func() error {
			time.Sleep(time.Second)
			counter.Add(1)
			return nil
		})
	}

	err := group.Wait()
	assert.Equal(t, int32(5), counter.Load())
	assert.NoError(t, err)
	assert.Equal(t, 3, group.maxWorkersNumber)
}
