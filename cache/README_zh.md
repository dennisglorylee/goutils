**其他语言版本: [English](README.md), [中文](README_zh.md).**



<!-- toc -->

- [缓存模块](#%E7%BC%93%E5%AD%98%E6%A8%A1%E5%9D%97)
  * [func_cache_test.go](#func_cache_testgo)

<!-- tocstop -->

# 缓存模块
## func_cache_test.go
### TestRdscCacheFunc
```go

ctx := context.Background()
cacher := mockGetCacher()

err := DeleteCache(ctx, cacher, "UTKey")
if err != nil {
	t.Error(err)
}

value1, err := CacheFunc0(ctx, cacher, "UTKey", 60*time.Second, rawGetFunc0)
if err != nil {
	t.Error(err)
}

value2, err := CacheFunc0(ctx, cacher, "UTKey", 60*time.Second, rawGetFunc0)
if err != nil {
	t.Error(err)
}

if value1 != value2 {
	t.Error(value1, value2)
}

err = DeleteCache(ctx, cacher, fmt.Sprintf("UT:%v", "p1"))
if err != nil {
	t.Error(err)
}

value1, err = CacheFunc1(ctx, cacher, fmt.Sprintf("UT:%v", "p1"), 60*time.Second, rawGetFunc1, "p1")
if err != nil {
	t.Error(err)
}

value2, err = CacheFunc1(ctx, cacher, fmt.Sprintf("UT:%v", "p1"), 60*time.Second, rawGetFunc1, "p1")
if err != nil {
	t.Error(err)
}

if value1 != value2 {
	t.Error(value1, value2)
}

err = DeleteCache(ctx, cacher, fmt.Sprintf("UT:%v:%v", "p1", "p2"))
if err != nil {
	t.Error(err)
}

value1, err = CacheFunc2(ctx, cacher, fmt.Sprintf("UT:%v:%v", "p1", "p2"), 60*time.Second, rawGetFunc2, "p1", "p2")
if err != nil {
	t.Error(err)
}

value2, err = CacheFunc2(ctx, cacher, fmt.Sprintf("UT:%v:%v", "p1", "p2"), 60*time.Second, rawGetFunc2, "p1", "p2")
if err != nil {
	t.Error(err)
}

if value1 != value2 {
	t.Error(value1, value2)
}

param3 := &testCacheParam{Param1: "p3"}
err = DeleteCache(ctx, cacher, fmt.Sprintf("UT:%v:%v:%v", "p1", "p2", param3.Param1))
if err != nil {
	t.Error(err)
}

value1, err = CacheFunc3(ctx, cacher, fmt.Sprintf("UT:%v:%v:%v", "p1", "p2", param3.Param1), 60*time.Second, rawGetFunc3, "p1", "p2", param3)
if err != nil {
	t.Error(err)
}

value2, err = CacheFunc3(ctx, cacher, fmt.Sprintf("UT:%v:%v:%v", "p1", "p2", param3.Param1), 60*time.Second, rawGetFunc3, "p1", "p2", param3)
if err != nil {
	t.Error(err)
}

if value1 != value2 {
	t.Error(value1, value2)
}
```
### TestRdscCacheFuncErr
```go

ctx := context.Background()
cacher := mockGetCacher()

err := DeleteCache(ctx, cacher, "UTKey")
if err != nil {
	t.Error(err)
}

value1, _ := CacheFunc0(ctx, cacher, "UTKey", 60*time.Second, rawGetFunc0Err)
value2, _ := CacheFunc0(ctx, cacher, "UTKey", 60*time.Second, rawGetFunc0Err)

if value1 != value2 {
	t.Error(value1, value2)
}

err = DeleteCache(ctx, cacher, fmt.Sprintf("UT:%v", "error"))
if err != nil {
	t.Error(err)
}

value1, _ = CacheFunc1(ctx, cacher, fmt.Sprintf("UT:%v", "error"), 60*time.Second, rawGetFunc1, "error")
value2, _ = CacheFunc1(ctx, cacher, fmt.Sprintf("UT:%v", "error"), 60*time.Second, rawGetFunc1, "error")

if value1 != value2 {
	t.Error(value1, value2)
}

err = DeleteCache(ctx, cacher, fmt.Sprintf("UT:%v:%v", "error", "p2"))
if err != nil {
	t.Error(err)
}

value1, _ = CacheFunc2(ctx, cacher, fmt.Sprintf("UT:%v:%v", "error", "p2"), 60*time.Second, rawGetFunc2, "error", "p2")
value2, _ = CacheFunc2(ctx, cacher, fmt.Sprintf("UT:%v:%v", "error", "p2"), 60*time.Second, rawGetFunc2, "error", "p2")

if value1 != value2 {
	t.Error(value1, value2)
}

param3 := &testCacheParam{Param1: "p3"}
err = DeleteCache(ctx, cacher, fmt.Sprintf("UT:%v:%v:%v", "error", "p2", param3.Param1))
if err != nil {
	t.Error(err)
}

value1, _ = CacheFunc3(ctx, cacher, fmt.Sprintf("UT:%v:%v:%v", "error", "p2", param3.Param1), 60*time.Second, rawGetFunc3, "error", "p2", param3)
value2, _ = CacheFunc3(ctx, cacher, fmt.Sprintf("UT:%v:%v:%v", "error", "p2", param3.Param1), 60*time.Second, rawGetFunc3, "error", "p2", param3)
if value1 != value2 {
	t.Error(value1, value2)
}
```
### TestRdsDeleteCacheTestMore
```go

ctx := context.Background()
cacher := mockGetCacher()

var err error

err = DeleteCache(ctx, cacher, fmt.Sprintf("GUT:%v:%v", "p1", "p2"))
if err != nil {
	t.Error(err)
}

result1, err := CacheFunc2(ctx, cacher, fmt.Sprintf("GUT:%v:%v", "p1", "p2"), 60*time.Second, rawGetFunc4, "p1", "p2")
if err != nil {
	t.Error(err)
}
log.Info(ctx, "%v %v %v", result1, err, mockPrintKind(result1))

result2, err := CacheFunc2(ctx, cacher, fmt.Sprintf("GUT:%v:%v", "p1", "p2"), 60*time.Second, rawGetFunc4, "p1", "p2")
if err != nil {
	t.Error(err)
}
log.Info(ctx, "%v %v %v", result2, err, mockPrintKind(result2))

if !reflect.DeepEqual(result1, result2) {
	t.Error(result1, result2)
}

err = DeleteCache(ctx, cacher, fmt.Sprintf("GUT:%v:%v", "p1", "p2"))
if err != nil {
	t.Error(err)
}

result3, err := CacheFunc2(ctx, cacher, fmt.Sprintf("GUT:%v:%v", "p1", "p2"), 60*time.Second, rawGetFunc5, "p1", "p2")
if err != nil {
	t.Error(err)
}
log.Info(ctx, "%v %v %v", result3, err, mockPrintKind(result3))

result4, err := CacheFunc2(ctx, cacher, fmt.Sprintf("GUT:%v:%v", "p1", "p2"), 60*time.Second, rawGetFunc5, "p1", "p2")
if err != nil {
	t.Error(err)
}
log.Info(ctx, "%v %v %v", result4, err, mockPrintKind(result4))

if !reflect.DeepEqual(result3, result4) {
	t.Error(result3, result4)
}
```
