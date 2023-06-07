 

# Micro Service

## AFK Diagram

### Visualization

<img src="/home/panhong/go/src/github.com/panhongrainbow/designPatternsExplainedz/assets/image-20230607170756026.png" alt="image-20230607170756026" style="zoom:80%;" />  

### X-Axis Expansion

| Tool              | Description                                                  |
| ----------------- | ------------------------------------------------------------ |
| Goroutine         | Horizontally expand directly within applications and implements high concurrency |
| Redis Cluster     | Add Cluster nodes to increase load capacity                  |
| Service Discovery | List existing services in the service discovery list<br />(列表出服务而已) |
| K8S               | Known for horizontal service expansion<br />(就出名啊，一定有人认为微服务一定要用 K8S) |
| Distributed Locks | Without distributed locks, horizontal expansion will result in resource competetion<br />(发展水平扩展所需的元件) |

### Y-Axis Expansion

The Y-axis is special and can be used with or without K8S

(Y轴最特别，可进 K8S 也可不进)

| Tool                                                         | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ |
| Normalization Partitioning                                   | Partition database tables based on their characteristics     |
| Hot and Cold Data                                            | Separately partition high-frequency data                     |
| After service decomposition,<br />one service corresponds to one database | This makes expansion very convenient afterwards<br />(最理想，因全解隅了) |
| Redis as Middleware                                          | Stores high-frequency data in Redis for caching              |
| MongoDB's Data Center                                        | Put different data in different MongoDB data centers<br />(就不同地区的资料中心) |
| Sharding of TiDB                                             | Perform database sharding and table sharding<br />(全库分表，要不要进 K8S ，自己决定) |

### Z-Axis Expansion

| Tool                                 | Description                                                  |
| ------------------------------------ | ------------------------------------------------------------ |
| Mysql read-write Separation          | Classify the databases into read and write categories<br />(也有点像水平扩展，因可增加唯读 Slave 节点) |
| Service read-write Separation        | Services also need to be decomposed                          |
| Independent High Concurrent Services | Separate code that requires high concurrency independently<br />(就是 Kafka 解隅方式之一) |
| The Service Decomposition of Kafka   | Uses Kafka to decompose services                             |
| Data Stream of Kafka                 | Transforms the original data stream into another new servics<br />(就是 Kafka 解隅方式之一) |

### Compatible with the three axes

| Tool          | Description                                                  |
| ------------- | ------------------------------------------------------------ |
| openTelemetry | openTelemetry is a large project with many components.<br />It requires horizontal expansion, data partitioning and service decomposition at the same time.<br/>When openTelemetry is healthy, it provides support and services for the other three axes.<br />(openTelemetry 自己本身就需要三轴扩展，健状后又可以支援三轴的其他服务) |
