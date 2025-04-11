package lru

import (
	"reflect"
	"testing" // Go的标准测试包，提供测试所需的工具
)

type String string

// 为 string 类型封装了一个实现了 Value 接口的类型
// 这个方法让String类型满足Value接口（在lru.go中定义的接口，要求实现Len() int方法）
func (d String) Len() int {
	return len(d)
}

// 测试缓存命中和缓存未命中的情况
/*
t *testing.T ----  Go测试函数的标准参数，用于报告测试结果
*/
func TestGet(t *testing.T) {
	lru := New(int64(0), nil)
	lru.Add("key1", String("1234"))
	// 第一个if语句：测试缓存命中的情况  如果命中失败则测试失败 || 类型断言 （接口转为String类型然后转为Go原生字符串比较）
	if v, ok := lru.Get("key1"); !ok || string(v.(String)) != "1234" {
		t.Fatalf("cache hit key1=1234 failed")
	}

	if _, ok := lru.Get("key2"); ok {
		t.Fatalf("cache miss key2 failed")
	}
}

/*
TestRemoveoldest - 测试淘汰机制

创建一个容量有限的缓存
添加多个键值对，触发淘汰机制
验证最旧的键值对被正确淘汰
*/
func TestRemoveoldest(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "k3"
	v1, v2, v3 := "value1", "value2", "v3"
	// 声明三组键值对
	//cap := len(k1 + k2 + v1 + v2)：计算缓存容量，刚好能容纳前两个键值对
	//创建有限容量的缓存
	cap := len(k1 + k2 + v1 + v2)
	lru := New(int64(cap), nil)
	lru.Add(k1, String(v1))
	lru.Add(k2, String(v2))
	lru.Add(k3, String(v3))

	if _, ok := lru.Get("key1"); ok || lru.Len() != 2 {
		t.Fatalf("Removeoldest key1 failed")
	}
}

/*
TestOnEvicted - 测试淘汰回调

设置一个回调函数，记录被淘汰的键
添加多个键值对，触发淘汰
验证回调函数被正确调用，并且淘汰的键符合预期
*/
func TestOnEvicted(t *testing.T) {
	// 创建一个空切片，用于记录被淘汰的键
	keys := make([]string, 0)
	// 定义一个回调函数，当项被淘汰时，将键添加到keys切片
	callback := func(key string, value Value) {
		keys = append(keys, key)
	}
	// 创建一个容量为10字节的缓存，并设置回调函数
	lru := New(int64(10), callback)
	lru.Add("key1", String("123456"))
	lru.Add("k2", String("k2"))
	lru.Add("k3", String("k3"))
	lru.Add("k4", String("k4"))

	expect := []string{"key1", "k2"}

	if !reflect.DeepEqual(expect, keys) {
		t.Fatalf("Call OnEvicted failed, expect keys equals to %s", expect)
	}
}

func TestAdd(t *testing.T) {
	lru := New(int64(0), nil)
	lru.Add("key", String("1"))
	lru.Add("key", String("111")) // 更新操作
	// 验证缓存占用的总字节数是否正确计算
	if lru.nbytes != int64(len("key")+len("111")) {
		t.Fatal("expected 6 but got", lru.nbytes)
	}
}
