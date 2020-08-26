# goCache
go 语言-进程内缓存套件 demo；采用 
## 淘汰策略算法
### FIFO
采用 list.list 双向链表 和 map
### LFU
采用最小堆
### LRU
采用 list.list 双向链表 和 map

## BigCache 进程内缓存-优化学习
### 并发优化
1. 分片技术
- 先进行 N 个数组的分片
- 通过散列函数 hash(key)%N 获取分片，此过程不需要锁
- 获取具体分片锁进行读写

N 尽量取 2 的幂；因为 `x % N == x&(N-1)`；数学可证

2. hash 算法的选择
BigCache 默认采用 FNV64a Hash 算法

### GC优化
GO 1.5 之后的优化：当 map 中的键值对都是基本类型，GC 将会忽略扫描它们

**BigCache 中的解决方案**：采用一个 `map[uint64]uint32` ，key 为 Hash 值，value 为另一个 offset 索引；
通过把缓存对象序列化后放到一个预先分配的大的字节数组中，然后把数组中的 offset 作为 `map[uint64]uint32` 的 value

