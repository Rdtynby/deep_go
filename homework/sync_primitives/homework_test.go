package main

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type RWMutex struct {
	mutex     sync.Mutex
	hasWriter atomic.Bool
	readers   atomic.Int32
}

func (m *RWMutex) Lock() {
	m.hasWriter.Store(true)
	m.mutex.Lock()
}

func (m *RWMutex) Unlock() {
	m.mutex.Unlock()
	m.hasWriter.Store(false)
}

func (m *RWMutex) RLock() {
	m.readers.Add(1)

	if m.readers.Load() == 1 {
		// It's the first reader, locking
		m.mutex.Lock()
	} else if m.hasWriter.Load() {
		// Waiting for writer
		m.mutex.Lock()
	}
}

func (m *RWMutex) RUnlock() {
	if m.readers.Load() == 1 {
		// It's the last reader
		m.mutex.Unlock()
	}

	m.readers.Add(-1)
}

func (m *RWMutex) TryLock() bool {
	m.hasWriter.Store(true)
	return m.mutex.TryLock()
}

func (m *RWMutex) TryRLock() bool {
	m.readers.Add(1)

	if m.readers.Load() == 1 {
		// It's the first reader, locking
		return m.mutex.TryLock()
	} else if m.hasWriter.Load() {
		// Waiting for writer
		return m.mutex.TryLock()
	}

	return false
}

func TestRWMutexWithWriter(t *testing.T) {
	var mutex RWMutex
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
	var mutex RWMutex
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
	var mutex RWMutex
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
	var mutex RWMutex
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
	var mutex RWMutex
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
	var mutex RWMutex
	assert.True(t, mutex.TryRLock())  // reader
	assert.False(t, mutex.TryRLock()) // reader

	var mutualExlusionWithWriter atomic.Bool
	mutualExlusionWithWriter.Store(true)

	go func() {
		assert.False(t, mutex.TryLock()) // another writer
		mutualExlusionWithWriter.Store(false)
	}()

	time.Sleep(time.Second)
	assert.False(t, mutualExlusionWithWriter.Load())
}
