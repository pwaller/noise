package main

import (
	"runtime"
	"sync"
	"testing"
)

func BenchmarkMutex(b *testing.B) {
	var m sync.RWMutex
	for i := 0; i < b.N; i++ {
		m.Lock()
		m.Unlock()
	}
}

func ContendedMutex(b *testing.B, n int) {
	var m sync.RWMutex

	die := false

	for i := 0; i < n; i++ {
		go func() {
			for {
				m.Lock()
				if die {
					break
				}
				m.Unlock()
			}
		}()
	}

	for i := 0; i < b.N; i++ {
		m.Lock()
		m.Unlock()
	}

	m.Lock()
	die = true
	m.Unlock()
}

func BenchmarkMutexContended1(b *testing.B) {
	ContendedMutex(b, 1)
}

func BenchmarkMutexContended2(b *testing.B) {
	ContendedMutex(b, 2)
}

func BenchmarkMutexContended3(b *testing.B) {
	ContendedMutex(b, 3)
}

func BenchmarkMutexContended4(b *testing.B) {
	ContendedMutex(b, 4)
}

func BenchmarkChannelRead(b *testing.B) {
	c := make(chan struct{})

	for i := 0; i < b.N; i++ {
		select {
		case <-c:
		default:
		}
	}
}
func BenchmarkChannelReadSched(b *testing.B) {
	c := make(chan struct{})

	for i := 0; i < b.N; i++ {
		select {
		case <-c:
		default:
		}
		runtime.Gosched()
	}
}

func ChannelReadContended(b *testing.B, n int) {

	c := make(chan struct{})

	die := make(chan struct{})
	defer close(die)

	for i := 0; i < n; i++ {
		go func() {
			for {
				select {
				case <-c:
				case <-die:
					return
				default:
				}
				runtime.Gosched()
			}
		}()
	}

	for i := 0; i < b.N; i++ {
		select {
		case <-c:
		default:
		}
		runtime.Gosched()
	}
}

func BenchmarkGosched(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runtime.Gosched()
	}
}

func BenchmarkChannelReadContended1(b *testing.B) {
	ChannelReadContended(b, 1)
}

func BenchmarkChannelReadContended2(b *testing.B) {
	ChannelReadContended(b, 2)
}

func BenchmarkChannelReadContended3(b *testing.B) {
	ChannelReadContended(b, 3)
}

func BenchmarkChannelReadContended4(b *testing.B) {
	ChannelReadContended(b, 4)
}
