package main

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type RWMutex struct {
	mutex          sync.Mutex
	readCond       *sync.Cond
	writeCond      *sync.Cond
	hasWriter      bool
	readers        int
	waitingWriters int
}

func NewRWMutex() *RWMutex {
	r := &RWMutex{}
	r.readCond = sync.NewCond(&r.mutex)
	r.writeCond = sync.NewCond(&r.mutex)

	return r
}

func (m *RWMutex) Lock() {
	m.mutex.Lock()
	m.waitingWriters++

	for m.hasWriter || m.readers > 0 {
		m.writeCond.Wait()
	}

	m.waitingWriters--
	m.hasWriter = true
	m.mutex.Unlock()
}

func (m *RWMutex) Unlock() {
	m.mutex.Lock()
	m.hasWriter = false

	if m.waitingWriters > 0 {
		m.writeCond.Signal()
	} else {
		m.readCond.Broadcast()
	}

	m.mutex.Unlock()
}

func (m *RWMutex) RLock() {
	m.mutex.Lock()

	for m.hasWriter || m.waitingWriters > 0 {
		m.readCond.Wait()
	}

	m.readers++
	m.mutex.Unlock()
}

func (m *RWMutex) RUnlock() {
	m.mutex.Lock()
	m.readers--

	if m.readers == 0 && m.waitingWriters > 0 {
		m.writeCond.Signal()
	}

	m.mutex.Unlock()
}

func (m *RWMutex) TryLock() bool {
	m.mutex.Lock()

	if m.hasWriter || m.readers > 0 {
		m.mutex.Unlock()

		return false
	}

	m.hasWriter = true
	m.mutex.Unlock()

	return true
}

func (m *RWMutex) TryRLock() bool {
	m.mutex.Lock()

	if m.hasWriter {
		m.mutex.Unlock()

		return false
	}

	m.readers++
	m.mutex.Unlock()

	return true
}

func TestRWMutexWithWriter(t *testing.T) {
	mutex := NewRWMutex()
	mutex.Lock() // writer

	var mutualExlusionWithWriter atomic.Bool
	mutualExlusionWithWriter.Store(true)
	var mutualExlusionWithReader atomic.Bool
	mutualExlusionWithReader.Store(true)

	go func() {
		mutex.Lock() // another writer
		mutualExlusionWithWriter.Store(false)
	}()

	go func() {
		mutex.RLock() // another reader
		mutualExlusionWithReader.Store(false)
	}()

	time.Sleep(time.Second)
	assert.True(t, mutualExlusionWithWriter.Load())
	assert.True(t, mutualExlusionWithReader.Load())
}

func TestRWMutexWithReaders(t *testing.T) {
	mutex := NewRWMutex()
	mutex.RLock() // reader

	var mutualExlusionWithWriter atomic.Bool
	mutualExlusionWithWriter.Store(true)

	go func() {
		mutex.Lock() // another writer
		mutualExlusionWithWriter.Store(false)
	}()

	time.Sleep(time.Second)
	assert.True(t, mutualExlusionWithWriter.Load())
}

func TestRWMutexMultipleReaders(t *testing.T) {
	mutex := NewRWMutex()
	mutex.RLock() // reader

	var readersCount atomic.Int32
	readersCount.Add(1)

	go func() {
		mutex.RLock() // another reader
		readersCount.Add(1)
	}()

	go func() {
		mutex.RLock() // another reader
		readersCount.Add(1)
	}()

	time.Sleep(time.Second)
	assert.Equal(t, int32(3), readersCount.Load())
}

func TestRWMutexWithWriterPriority(t *testing.T) {
	mutex := NewRWMutex()
	mutex.RLock() // reader

	var mutualExlusionWithWriter atomic.Bool
	mutualExlusionWithWriter.Store(true)
	var readersCount atomic.Int32
	readersCount.Add(1)

	go func() {
		mutex.Lock() // another writer is waiting for reader
		mutualExlusionWithWriter.Store(false)
	}()

	time.Sleep(time.Second)

	go func() {
		mutex.RLock() // another reader is waiting for a higher priority writer
		readersCount.Add(1)
	}()

	go func() {
		mutex.RLock() // another reader is waiting for a higher priority writer
		readersCount.Add(1)
	}()

	time.Sleep(time.Second)

	assert.True(t, mutualExlusionWithWriter.Load())
	assert.Equal(t, int32(1), readersCount.Load())
}

func TestTryRWMutexWithWriter(t *testing.T) {
	mutex := NewRWMutex()
	assert.True(t, mutex.TryLock())  // writer
	assert.False(t, mutex.TryLock()) // writer

	var mutualExlusionWithWriter atomic.Bool
	mutualExlusionWithWriter.Store(true)
	var mutualExlusionWithReader atomic.Bool
	mutualExlusionWithReader.Store(true)

	go func() {
		assert.False(t, mutex.TryLock()) // another writer
		mutualExlusionWithWriter.Store(false)
	}()

	go func() {
		assert.False(t, mutex.TryRLock()) // another reader
		mutualExlusionWithReader.Store(false)
	}()

	time.Sleep(time.Second)
	assert.False(t, mutualExlusionWithWriter.Load())
	assert.False(t, mutualExlusionWithReader.Load())
}

func TestTryRWMutexWithReaders(t *testing.T) {
	mutex := NewRWMutex()
	assert.True(t, mutex.TryRLock()) // reader
	assert.True(t, mutex.TryRLock()) // reader

	var mutualExlusionWithWriter atomic.Bool
	mutualExlusionWithWriter.Store(true)

	go func() {
		assert.False(t, mutex.TryLock()) // another writer
		mutualExlusionWithWriter.Store(false)
	}()

	time.Sleep(time.Second)
	assert.False(t, mutualExlusionWithWriter.Load())
}
