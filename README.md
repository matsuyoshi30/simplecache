# simplecache

## LRU (Least Recently Used)

```go
l := lru.NewLRU(3)

l.Put(1, 1)
l.Put(2, 2)
l.Put(3, 3)
l.Put(4, 4)

_, err := l.Get(1) // returns not found error
v, err := l.Get(2) 
```

## LFU (Least Frequently Used)

```go
l := lfu.NewLFU(3)

l.Put(1, 1)
l.Put(2, 2)
l.Put(3, 3)

v, err := l.Get(1) // update entry key 1

l.Put(4, 4)

__, err := l.Get(2) // returns not found error
```

# Author

[matsuyoshi30](https://twitter.com/matsuyoshi30)
