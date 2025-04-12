# 在线聊天室与用户统计功能设计文档

## 1. 功能概述

本文档详细描述了为个人网站添加在线聊天室和实时在线人数统计两个功能的设计与实现方案。

### 1.1 在线聊天室功能

- 允许访问网站的用户进行实时文字交流
- 支持用户设置昵称
- 展示消息发送时间
- 支持简单的表情符号
- 提供历史消息加载功能

### 1.2 在线人数统计功能

- 实时显示当前浏览网站的用户数量
- 在导航栏或页脚位置展示
- 数据实时更新

## 2. 技术选型

系统采用前后端分离架构，前端使用Vue 3，后端使用Go语言，通过WebSocket进行实时通信。

### 2.1 后端技术

- **Go**：高性能、并发友好的编程语言
- **gorilla/websocket**：处理WebSocket连接
- **gin**：Web框架，提供REST API功能
- **Redis**：用于存储在线用户信息和聊天历史记录

### 2.2 前端技术

- **Vue 3**：现有前端框架
- **Pinia**：状态管理
- **date-fns**：日期格式化
- **Vue Router**：页面路由管理

## 3. 系统架构

### 3.1 整体架构

```
+------------------+       WebSocket       +------------------+
|                  |<--------------------->|                  |
|  Browser Client  |                       |  Go Server       |
|  (Vue.js App)    |<--------------------->|  (gin+websocket) |
|                  |       HTTP/API        |                  |
+------------------+                       +--------+---------+
                                                    |
                                                    | 数据存储
                                                    v
                                           +------------------+
                                           |                  |
                                           |  Redis Database  |
                                           |                  |
                                           +------------------+
```

### 3.2 功能模块划分

1. **前端模块**
   - 聊天室UI组件
   - 在线人数显示组件
   - 聊天状态管理模块
   - WebSocket连接管理模块

2. **后端模块**
   - WebSocket服务
   - 用户会话管理
   - 消息广播服务
   - Redis数据存储服务
   - HTTP API服务

## 4. 详细设计

### 4.1 数据结构设计

#### 4.1.1 消息数据结构

```go
// 后端Go结构体
type Message struct {
    ID        string    `json:"id"`
    Sender    User      `json:"sender"`
    Content   string    `json:"content"`
    Timestamp int64     `json:"timestamp"`
    Type      string    `json:"type"` // "text", "system", "emoji"
}
```

#### 4.1.2 用户数据结构

```go
// 后端Go结构体
type User struct {
    ID         string `json:"id"`         // WebSocket连接ID
    Nickname   string `json:"nickname"`   // 用户昵称
    JoinTime   int64  `json:"joinTime"`   // 加入时间
    LastActive int64  `json:"lastActive"` // 最后活跃时间
}
```

### 4.2 接口设计

#### 4.2.1 WebSocket事件

| 事件名称 | 方向 | 描述 |
|---------|------|------|
| connect | 客户端→服务器 | 用户连接 |
| disconnect | 服务器→客户端 | 用户断开连接 |
| chat:message | 双向 | 聊天消息 |
| user:join | 服务器→客户端 | 用户加入 |
| user:leave | 服务器→客户端 | 用户离开 |
| user:nickname | 客户端→服务器 | 用户修改昵称 |
| chat:welcome | 服务器→客户端 | 欢迎消息 |
| chat:history | 服务器→客户端 | 聊天历史 |
| online:count | 服务器→客户端 | 在线人数更新 |

#### 4.2.2 REST API端点

| 端点 | 方法 | 描述 |
|-----|------|------|
| /api/chat/history | GET | 获取聊天历史 |
| /api/chat/users | GET | 获取在线用户列表 |
| /api/chat/online-count | GET | 获取在线人数 |

### 4.3 Redis数据存储

- 使用Redis列表存储聊天历史：`chat:history`
- 使用Redis哈希表存储在线用户：`chat:users`
- 使用Redis计数器跟踪在线人数：`chat:onlineCount`

## 5. 安全与性能考虑

### 5.1 安全考虑

1. **XSS攻击防护**
   - 服务器端对聊天内容进行过滤和转义
   - 前端使用 v-html 时需小心处理

2. **消息频率限制**
   - 实现速率限制以防止垃圾消息
   - Go后端可以使用令牌桶算法进行限流

3. **用户身份**
   - 为每个会话生成唯一标识符
   - 考虑后期集成用户认证系统

### 5.2 性能考虑

1. **Go语言的并发优势**
   - 使用goroutine处理WebSocket连接，高效管理大量并发连接
   - 使用channel进行消息广播

2. **Redis优化**
   - 使用Redis管道减少网络往返
   - 定期清理过期数据
   - 对历史记录设置合理的存储上限

3. **前端优化**
   - 消息懒加载和分页
   - 减少不必要的DOM更新

## 6. 扩展功能（未来可考虑）

1. 私聊功能
2. 图片分享功能
3. 用户身份验证
4. 消息搜索功能
5. 消息通知
6. 多房间/频道支持

## 7. 测试策略

1. **单元测试**
   - 使用Go标准测试工具测试核心功能
   - Vue组件单元测试

2. **集成测试**
   - WebSocket连接测试
   - 消息流测试
   - Redis交互测试

3. **负载测试**
   - 模拟多用户并发连接
   - 评估系统在高负载下的性能 