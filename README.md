## goLocalMemory 是什么?

goLocalMemory是一个基于Go语言本地内存开发的数据结构工具集，用来帮助Go业务系统更好的工作。

### setGroup简介
setGroup是一个管理一组集合(set)的数据结构，这组集合共同管理相同的资源(member)，而setGroup
提供了快捷的接口帮助使用者。它有下边的特点：

- 占用内存少，因为这组集合共享相同的资源，所以只存储一次
- 快速求组内多个集合的交、差、并，使用了bitMap和优化算法，时间复杂度 < Olog(n), n为集合组管理
元素的总和
  
//todo:项目处于初始阶段，后续会添加更多有用的数据结构和使用文档、压测结果等。