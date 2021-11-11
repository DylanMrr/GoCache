# GoCache 
Simple in-memory key-value cache with default or specific expiration time.

# Install
```
go get github.com/DylanMrr/GoCache
```
# Features
- Key-value storage
- Expiration time can be set as default, permanent or specific
- Native golang

# Usage
- Creating
```
//New cache
c := gocache.New(time.Second, 500 * time.Millisecond)
```
First parameter - clean time
Second parameter - default expiration time

- Adding
```
//Adding value "now" with key "first", default expiration time
c.Add("first", "now", gocache.Default, gocache.Expires{})
//Adding value "future" with key "second" without expiration time
c.Add("second", "future", gocache.Permanent, gocache.Expires{})
//Adding value "next" with key "third" and 4 seconds expiration time
c.Add("third", "next", gocache.Specific, gocache.Expires{ExpiresDuration: 4 * time.Second})
```

- Deleting
```
//Delete value with key "hello"
c.Delete("hello")
```

- Getting
```
//Get value with key "first"
now, found := c.Get("first")
```

- Range
```
//Range calls function sequentially for each key and value
c.Range(func(key, value interface{}) bool {
		count++
		return true
	})
```

- Count
```
//Count of elements in cache
length := c.Count()
```
