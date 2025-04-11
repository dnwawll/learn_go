package lru

import "container/list"

// Cache is a LRU cache. It is not safe for concurrent access.
type Cache struct {
	maxBytes  int64                         // 缓存的最大容量（字节为单位）
	nbytes    int64                         // 当前已使用的字节数
	ll        *list.List                    // 双向链表，用于维护数据项的访问顺序
	cache     map[string]*list.Element      // 哈希表，用于快速查找数据   值 (value) 是 *list.Element 类型，即指向 list.Element 的指针
	OnEvicted func(key string, value Value) // 可选的回调函数，在数据被淘汰时执行
}

// 存储在链表中的数据项
type entry struct {
	key   string
	value Value
}

// Value use Len to count how many bytes it takes
type Value interface {
	Len() int // 这个接口只定义了一个方法 Len()，该方法没有参数并返回一个 int。任何类型只要实现了 Len() int 方法，就自动满足了 Value 接口。
}

// New is the Constructor of Cache 初始化缓存结构，设置最大容量和淘汰回调函数 || 返回值类型是 *Cache（指向 Cache 的指针）
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// Add adds a value to the cache.
/*(c *Cache) 是接收者，表示这个方法属于 *Cache 类型 接受两个参数：key (字符串) 和 value (实现了 Value 接口的任何类型)
对value来说 -->
interface{} 是一个特殊的接口类型，它不包含任何方法定义，因此任何类型都满足这个接口。它类似于 Java 中的 Object 或 C++ 中的 void*，可以存储任何类型的值。
ele.Value 的类型是 interface{}
(*entry) 表示我们断言这个接口值实际存储的是一个 *entry 类型的值
如果断言成功，kv 将获得这个 *entry 值
如果断言失败（实际存储的不是 *entry 类型），程序会产生 panic
-----
在 LRU 缓存实现中，链表 list.List 存储的是 *entry 类型的值，
但是由于 list.Element.Value 的类型是 interface{}，所以我们存入后又取出时，需要通过类型断言转回原来的 *entry 类型。
*/
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok { // 在 Go 中，map 查找会返回两个值：查找到的元素和一个表示是否找到的布尔值
		c.ll.MoveToFront(ele) // 链表元素移动到最前面表示最近使用
		kv := ele.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len()) // 更新已用字节数，计算新值和旧值的大小差异
		kv.value = value
	} else {
		// entry{key, value} 创建一个 entry 结构体  & 获取这个结构体的指针
		ele := c.ll.PushFront(&entry{key, value}) // 将新条目添加到链表前端，返回链表元素
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

// Get look ups a key's value
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// RemoveOldest removes the oldest item
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// Len the number of cache entries
func (c *Cache) Len() int {
	return c.ll.Len()
}
