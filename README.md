## goLocalMemory 是什么?

goLocalMemory是一个基于Go语言本地内存开发的数据结构工具集，用来帮助Go业务系统更好的工作。

### 数据结构setGroup简介
setGroup是一个管理一组集合(set)的数据结构，适合集合共同管理相同的资源(member)的业务场景，而setGroup
提供了快捷的接口帮助使用者。它有下边的特点：

- 占用内存少，因为这组集合共享相同的资源，所以只存储一次
- 快速求组内多个集合的交、差、并，使用了bitMap和优化算法，时间复杂度 < O(n), n为集合组管理
元素的总个数
  
### setGroup接口方法
```bigquery
Add(setKey interface{}, member interface{}) 向setGroup中添加set元素
Remove(setKey interface{}, member interface{}) 移除setGroup中的元素
Intersect(setKeys ...interface{}) []interface{} 求setGroup中多个set的交集
Union(setKeys ...interface{}) []interface{} 求setGroup中多个set的并集
Different(setKeys ...interface{}) []interface{} 求setGroup中多个set的差集
```
  
### setGroup压测结果
```bigquery
goos: darwin
goarch: amd64
pkg: goLocalMemory
BenchmarkAddOneSetGroup-12            	 1000000	      1228 ns/op
BenchmarkAddMultiSetGroup-12          	  347254	     36716 ns/op
BenchmarkRemoveMultiSetGroup-12       	 2978799	       390 ns/op
BenchmarkMultiIntersectSetGroup-12    	 1493817	       799 ns/op
BenchmarkMultiUnionSetGroup-12        	 1413752	       849 ns/op
BenchmarkMultiDifferentSetGroup-12    	 1413374	       832 ns/op
PASS
ok  	goLocalMemory	22.555s
```
