# 聊天服务器

这是一个基于WebSocket的实时聊天服务器实现，使用Go语言开发。

## 功能特性

* 基于WebSocket的实时消息传递
* 多房间支持，用户可以在不同房间间切换
* 消息历史记录，使用Redis存储
* 用户在线状态管理
* RESTful API接口
* 跨域支持

## 系统要求

* Go 1.18+
* Redis 6.0+

## 配置和启动

1. 编辑 `configs/config.yaml` 文件配置服务器参数

2. 启动Redis:
```
# 在Windows上
redis-server

# 在Linux/Mac上
redis-server
```

3. 构建并运行服务器:
```
cd chat-server
go build -o chat-server ./cmd/server
./chat-server  # 在Windows上使用 chat-server.exe
```

## WebSocket连接参数

连接URL: `ws://localhost:8080/ws?userId=xxx&username=xxx&roomId=general`

参数说明:
* userId: 用户唯一ID
* username: 用户名
* roomId: 房间ID（可选，默认为"general"）

## API接口

### 获取聊天室列表
```
GET /api/rooms
```

### 获取聊天室信息
```
GET /api/rooms/:id
```

### 获取聊天室消息历史
```
GET /api/rooms/:id/messages
```

### 获取用户信息
```
GET /api/users/:id
```

### 用户登录
```
POST /api/users/login
```

### 用户注册
```
POST /api/users/register
```

## 项目结构说明

```
chat-server/
├── cmd/                  # 命令行入口
│   └── server/           # 服务器主程序
├── configs/              # 配置文件目录
├── internal/             # 内部包
│   ├── api/              # API处理
│   ├── config/           # 配置加载
│   ├── models/           # 数据模型
│   ├── redis/            # Redis连接管理
│   ├── websocket/        # WebSocket处理
│   └── utils/            # 工具函数
├── go.mod                # Go模块文件
└── go.sum                # 依赖校验和
```

## 消息格式

```json
{
  "id": "消息ID",
  "roomId": "房间ID",
  "senderId": "发送者ID",
  "content": "消息内容",
  "type": "text/image/system/join/leave/command",
  "timestamp": "2023-01-01T00:00:00Z",
  "metadata": {}  // 可选的额外数据
}
```

## 许可证

MIT License 