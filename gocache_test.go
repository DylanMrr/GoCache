package gocache_test

import (
	"testing"
	"time"

	gocache "github.com/DylanMrr/GoCache"
)

func TestGetAdd(t *testing.T) {
	cycle := 100 * time.Millisecond
	defaultExpires := 5 * time.Millisecond

	c := gocache.New(cycle, defaultExpires)
	defer c.Close()

	c.Add("first", "now", gocache.Default, gocache.Expires{})
	c.Add("second", "future", gocache.Permanent, gocache.Expires{})
	c.Add("third", "next", gocache.Specific, gocache.Expires{ExpiresDuration: cycle / 4})

	now, found := c.Get("first")
	if !found {
		t.FailNow()
	}

	if now.(string) != "now" {
		t.FailNow()
	}

	time.Sleep(2 * cycle)

	_, found = c.Get("third")

	if found {
		t.FailNow()
	}

	time.Sleep(cycle * 2)

	_, found = c.Get("first")

	if found {
		t.FailNow()
	}

	_, found = c.Get("second")

	if !found {
		t.FailNow()
	}
}

func TestDelete(t *testing.T) {
	c := gocache.New(time.Minute, 1*time.Second)
	c.Add("hello", "Hello", gocache.Default, gocache.Expires{})
	_, found := c.Get("hello")

	if !found {
		t.FailNow()
	}

	c.Delete("hello")

	_, found = c.Get("hello")

	if found {
		t.FailNow()
	}
}

func TestCount(t *testing.T) {
	cycle := 100 * time.Millisecond
	defaultExpires := 1 * time.Second

	c := gocache.New(cycle, defaultExpires)
	defer c.Close()

	c.Add("first", "now", gocache.Default, gocache.Expires{})
	c.Add("second", "future", gocache.Permanent, gocache.Expires{})
	c.Add("third", "next", gocache.Specific, gocache.Expires{ExpiresDuration: cycle / 4})

	if c.Count() != 3 {
		t.FailNow()
	}
}
