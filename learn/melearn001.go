package main

//
//type Cac struct {
//	max64     int64
//	nbytes    int64
//	ll        *list.List
//	cache     map[string]*list.Element
//	onEvicted func(key string, value Value)
//}
//
//type Value interface {
//	Len() int
//}
//type entry struct {
//	key   string
//	value Value
//}
//
//func New(max64 int64, onEnvicted func(key string, value Value)) *Cac {
//	return &Cac{
//		max64:     max64,
//		ll:        list.New(),
//		cache:     make(map[string]*list.Element),
//		onEvicted: onEnvicted,
//	}
//}
//
//func (c *Cac) Add(key string, value Value) {
//	if ele, ok := c.cache[key]; ok {
//		c.ll.MoveToFront(ele)
//		kv := ele.Value.(*entry)
//		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
//		kv.value = value
//	} else {
//		ele := c.ll.PushFront(&entry{key, value})
//		c.cache[key] = ele
//		c.nbytes += int64(len(key)) + int64(value.Len())
//	}
//
//	for c.max64 != 0 && c.max64 < c.nbytes {
//
//	}
//}
//
//func (c *Cac) RemoveOldest() {
//	ele := c.ll.Back()
//	if ele != nil {
//		c.ll.Remove(ele)
//		kv := ele.Value.(*entry)
//		delete(c.cache, kv.key)
//		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
//		if c.onEvicted != nil {
//			c.onEvicted(kv.key, kv.value)
//		}
//	}
//}
