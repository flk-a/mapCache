package mapCache

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const TestSucceedSign = "\033[1;32m\u2713\033[0m"

func TestNewMapCache(t *testing.T) {
	t.Log("no ttl used")
	{
		mc := NewMapCache[uint, uint](10, 0)
		var i uint
		for ; i < 10; i++ {
			mc.Set(i, i)
		}
		require.Equal(t, 10, mc.curSize)
		require.Less(t, mc.front, 10)
		require.Less(t, mc.rear, 10)
		for i = 0; i < 10; i++ {
			v, found := mc.Get(i)
			require.Equal(t, i, v)
			require.True(t, found)
		}
		t.Logf("\t%s\t succeed", TestSucceedSign)

	}
	t.Log("concurrency")
	{
		mc := NewMapCache[uint, uint](10, time.Hour)
		var i uint
		var wg sync.WaitGroup
		wg.Add(10)
		for ; i < 10; i++ {
			go func(k uint) {
				defer wg.Done()
				mc.Set(k, k)
			}(i)
		}
		wg.Wait()
		require.Equal(t, 10, mc.curSize)
		require.Less(t, mc.front, 10)
		require.Less(t, mc.rear, 10)
		for i = 0; i < 10; i++ {
			v, found := mc.Get(i)
			require.Equal(t, i, v)
			require.True(t, found)
		}
		t.Logf("\t%s\t succeed", TestSucceedSign)
	}
	t.Log("inner pointers leak")
	{
		mc := NewMapCache[uint, uint](10, time.Hour)
		var i uint
		for ; i < 10000; i++ {
			mc.Set(i, i)
		}
		require.Equal(t, 10, mc.curSize)
		require.Less(t, mc.front, 10)
		require.Less(t, mc.rear, 10)
		t.Logf("\t%s\t succeed", TestSucceedSign)
	}

	t.Log("timeout exceed")
	{
		mc := NewMapCache[uint, uint](10, time.Millisecond)
		var i uint
		for ; i < 1000; i++ {
			mc.Set(i, i)
		}
		time.Sleep(10 * time.Millisecond)
		var v any
		v, found := mc.Get(1)
		require.Zero(t, v)
		require.False(t, found)
		require.True(t, mc.isEmpty())
		require.Empty(t, mc.m)
		t.Logf("\t%s\t succeed", TestSucceedSign)
	}
	t.Log("not totally filled")
	{
		mc := NewMapCache[uint, uint](10, time.Hour)
		var i uint
		for ; i < 9; i++ {
			mc.Set(i, i)
		}
		for i = 0; i < 9; i++ {
			v := mc.pop()
			delete(mc.m, v.val)
			require.Equal(t, i, v.val)
		}
		require.True(t, mc.isEmpty())
		require.Empty(t, mc.m)
		t.Logf("\t%s\t succeed", TestSucceedSign)
	}
	t.Log("lru check")
	{
		mc := NewMapCache[uint, uint](10, time.Hour)
		var i uint
		for ; i < 12; i++ {
			mc.Set(i, i)
		}
		for i = 2; i < 12; i++ {
			v, found := mc.Get(i)
			require.Equal(t, i, v)
			require.True(t, found)
		}
		t.Logf("\t%s\t succeed", TestSucceedSign)
	}
	t.Log("full cache")
	{
		mc := NewMapCache[uint, uint](10, time.Hour)
		var i uint
		for ; i < 10; i++ {
			mc.Set(i, i)
		}
		for i = 0; i < 10; i++ {
			v, found := mc.Get(i)
			require.Equal(t, i, v)
			require.True(t, found)
		}
		t.Logf("\t%s\t succeed", TestSucceedSign)
	}
}

func ExampleMapCache() {
	mc := NewMapCache[int64, uint](10000, time.Hour)
	// adding new element to cache
	mc.Set(12, 123456)
	// reading some element
	value, found := mc.Get(12)
	// use result
	fmt.Println(value, found)
}

func BenchmarkMapCacheRW(b *testing.B) {
	mc := NewMapCache[int, uint](10000, time.Hour)
	for i := 0; i < b.N; i++ {
		mc.Set(i, 123456789)
	}
	for i := 0; i < b.N; i++ {
		mc.Get(i)
	}
}

func BenchmarkMapCacheR(b *testing.B) {
	mc := NewMapCache[int, string](10000, time.Hour)
	for i := 0; i < 10000; i++ {
		mc.Set(i, "123456789")
	}
	for i := 0; i < b.N; i++ {
		mc.Get(i)
	}
}

func BenchmarkMapCacheW(b *testing.B) {
	mc := NewMapCache[int, uint](10000, time.Hour)
	for i := 0; i < b.N; i++ {
		mc.Set(i%10000, 123456789)
	}
}
