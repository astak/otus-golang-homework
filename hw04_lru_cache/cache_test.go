package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		cache := NewCache(2)
		cache.Set("aaa", 1)
		cache.Set("bbb", 2)
		cache.Clear()
		_, ok := cache.Get("aaa")
		require.False(t, ok)
		_, ok = cache.Get("bbb")
		require.False(t, ok)
	})

	t.Run("capacity logic", func(t *testing.T) {
		c := NewCache(1)
		c.Set("aaa", 1)
		c.Set("bbb", 2)
		_, ok := c.Get("aaa")
		require.False(t, ok)

		c = NewCache(2)
		c.Set("aaa", 1)
		c.Set("bbb", 2)
		_, ok = c.Get("aaa")
		require.True(t, ok)
		c.Set("ccc", 3)
		c.Set("ddd", 4)
		_, ok = c.Get("aaa")
		require.False(t, ok)

		c = NewCache(3)
		c.Set("aaa", 1)
		c.Set("bbb", 2)
		c.Set("ccc", 3)
		_, ok = c.Get("aaa")
		require.True(t, ok)
		c.Set("ddd", 4)
		c.Set("eee", 5)
		c.Set("fff", 6)
		_, ok = c.Get("aaa")
		require.False(t, ok)
	})

	t.Run("mru logic", func(t *testing.T) {
		c := NewCache(2)
		c.Set("aaa", 1)
		c.Set("bbb", 2)
		c.Set("aaa", 3)
		c.Set("ccc", 4)
		_, ok := c.Get("aaa")
		require.True(t, ok)
		_, ok = c.Get("bbb")
		require.False(t, ok)

		c.Set("ddd", 5)
		_, ok = c.Get("aaa")
		require.True(t, ok)
		_, ok = c.Get("ccc")
		require.False(t, ok)

		c.Set("eee", 6)
		_, ok = c.Get("aaa")
		require.True(t, ok)
		_, ok = c.Get("ddd")
		require.False(t, ok)
	})

	t.Run("add logic", func(t *testing.T) {
		c := NewCache(2)
		c.Set("aaa", 1)
		c.Set("bbb", 2)
		c.Set("ccc", 3)
		_, ok := c.Get("bbb")
		require.True(t, ok)

		c = NewCache(3)
		c.Set("aaa", 1)
		c.Set("bbb", 2)
		c.Set("ccc", 3)
		c.Set("ddd", 4)
		c.Set("eee", 5)
		_, ok = c.Get("ccc")
		require.True(t, ok)

		c = NewCache(4)
		c.Set("aaa", 1)
		c.Set("bbb", 2)
		c.Set("ccc", 3)
		c.Set("ddd", 4)
		c.Set("eee", 5)
		c.Set("fff", 6)
		c.Set("ggg", 7)
		_, ok = c.Get("ddd")
		require.True(t, ok)
	})
}

func TestCacheMultithreading(_ *testing.T) {
	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
