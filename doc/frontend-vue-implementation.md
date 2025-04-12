# Vue前端实现设计

本文档详细描述了聊天室和在线人数统计功能的Vue前端实现方案。

## 1. 项目结构

```
src/
├── components/
│   ├── chat/
│   │   ├── ChatRoom.vue          # 聊天室主容器
│   │   ├── ChatMessage.vue       # 单条消息组件
│   │   ├── ChatInput.vue         # 消息输入组件
│   │   ├── UserList.vue          # 用户列表组件
│   │   └── EmojiPicker.vue       # 表情选择器
│   ├── common/
│   │   └── OnlineCounter.vue     # 在线人数计数器
│   └── ...                       # 其他现有组件
├── views/
│   ├── ChatView.vue              # 聊天页面视图
│   └── ...                       # 其他现有视图
├── stores/
│   ├── chat.js                   # 聊天状态管理
│   └── ...                       # 其他现有状态管理
├── utils/
│   ├── websocket.js              # WebSocket工具
│   ├── date.js                   # 日期格式化工具
│   └── emoji.js                  # 表情工具函数
├── router/
│   └── index.js                  # 路由配置
└── assets/
    └── styles/
        └── chat.css              # 聊天相关样式
```

## 2. 状态管理

使用Pinia作为状态管理库，创建聊天状态存储：

```javascript
// src/stores/chat.js

import { defineStore } from 'pinia';
import { format } from 'date-fns';
import { zhCN } from 'date-fns/locale';

// WebSocket URL
const WS_URL = import.meta.env.VITE_WS_URL || 'ws://localhost:8080/ws';
// API URL
const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api';

export const useChatStore = defineStore('chat', {
  state: () => ({
    // WebSocket连接
    socket: null,
    connected: false,
    connecting: false,
    reconnectAttempts: 0,
    maxReconnectAttempts: 5,
    reconnectTimeout: null,
    
    // 用户数据
    currentUser: null,
    users: [],
    onlineCount: 0,
    
    // 消息数据
    messages: [],
    isLoadingHistory: false,
    historyPageSize: 20,
    historyPage: 0,
    hasMoreHistory: true,
    
    // 错误状态
    error: null
  }),
  
  getters: {
    isConnected: (state) => state.connected,
    
    // 格式化消息，添加时间显示
    formattedMessages: (state) => state.messages.map(msg => ({
      ...msg,
      formattedTime: format(new Date(msg.timestamp), 'HH:mm:ss', { locale: zhCN })
    })),
    
    // 排序后的用户列表
    sortedUsers: (state) => [...state.users].sort((a, b) => 
      a.nickname.localeCompare(b.nickname, 'zh-CN')
    )
  },
  
  actions: {
    // 连接WebSocket
    async connect() {
      if (this.socket && this.connected) return;
      
      this.connecting = true;
      this.error = null;
      
      try {
        // 创建WebSocket连接
        this.socket = new WebSocket(WS_URL);
        
        // 设置事件处理器
        this.socket.onopen = this.handleSocketOpen;
        this.socket.onmessage = this.handleSocketMessage;
        this.socket.onclose = this.handleSocketClose;
        this.socket.onerror = this.handleSocketError;
      } catch (error) {
        this.handleConnectionError(error);
      }
    },
    
    // 断开连接
    disconnect() {
      if (this.reconnectTimeout) {
        clearTimeout(this.reconnectTimeout);
        this.reconnectTimeout = null;
      }
      
      if (this.socket) {
        this.socket.close();
        this.socket = null;
        this.connected = false;
        this.connecting = false;
        this.currentUser = null;
      }
    },
    
    // 重新连接
    reconnect() {
      if (this.reconnectAttempts >= this.maxReconnectAttempts) {
        this.error = '无法连接到服务器，请稍后重试';
        this.connecting = false;
        return;
      }
      
      this.reconnectAttempts++;
      this.reconnectTimeout = setTimeout(() => {
        console.log(`尝试重新连接(${this.reconnectAttempts}/${this.maxReconnectAttempts})...`);
        this.connect();
      }, 2000 * Math.pow(2, this.reconnectAttempts - 1)); // 指数退避
    },
    
    // 发送消息
    sendMessage(content) {
      if (!this.socket || !this.connected || !content.trim()) return;
      
      const message = {
        type: 'chat:message',
        content: content.trim()
      };
      
      this.socket.send(JSON.stringify(message));
    },
    
    // 设置昵称
    setNickname(nickname) {
      if (!this.socket || !this.connected || !nickname.trim()) return;
      
      const message = {
        type: 'user:nickname',
        content: nickname.trim()
      };
      
      this.socket.send(JSON.stringify(message));
    },
    
    // 加载聊天历史
    async loadChatHistory(refresh = false) {
      if (this.isLoadingHistory || (!refresh && !this.hasMoreHistory)) return;
      
      const page = refresh ? 0 : this.historyPage;
      
      this.isLoadingHistory = true;
      
      try {
        const response = await fetch(
          `${API_URL}/chat/history?page=${page}&pageSize=${this.historyPageSize}`
        );
        
        if (!response.ok) {
          throw new Error('无法加载聊天历史');
        }
        
        const data = await response.json();
        
        if (refresh) {
          this.messages = data.messages;
          this.historyPage = 1;
        } else {
          this.messages = [...this.messages, ...data.messages];
          this.historyPage++;
        }
        
        this.hasMoreHistory = data.messages.length === this.historyPageSize;
      } catch (error) {
        console.error('加载聊天历史失败:', error);
        this.error = '加载聊天历史失败';
      } finally {
        this.isLoadingHistory = false;
      }
    },
    
    // 获取在线用户列表
    async fetchUsers() {
      try {
        const response = await fetch(`${API_URL}/chat/users`);
        
        if (!response.ok) {
          throw new Error('无法加载用户列表');
        }
        
        const data = await response.json();
        this.users = data.users;
      } catch (error) {
        console.error('获取用户列表失败:', error);
      }
    },
    
    // 获取在线人数
    async fetchOnlineCount() {
      try {
        const response = await fetch(`${API_URL}/chat/online-count`);
        
        if (!response.ok) {
          throw new Error('无法获取在线人数');
        }
        
        const data = await response.json();
        this.onlineCount = data.count;
      } catch (error) {
        console.error('获取在线人数失败:', error);
      }
    },
    
    // === WebSocket事件处理器 ===
    
    // 连接打开事件
    handleSocketOpen() {
      this.connected = true;
      this.connecting = false;
      this.reconnectAttempts = 0;
      this.error = null;
      console.log('已连接到聊天服务器');
      
      // 连接成功后加载历史消息
      this.loadChatHistory(true);
      this.fetchUsers();
    },
    
    // 收到消息事件
    handleSocketMessage(event) {
      try {
        const message = JSON.parse(event.data);
        
        switch (message.type) {
          case 'text': // 普通聊天消息
            this.messages.push(message);
            break;
            
          case 'system': // 系统消息
            this.messages.push(message);
            break;
            
          case 'user:join': // 用户加入
            // 添加到用户列表
            if (!this.users.some(u => u.id === message.sender.id)) {
              this.users.push(message.sender);
            }
            this.messages.push(message);
            break;
            
          case 'user:leave': // 用户离开
            // 从用户列表中移除
            this.users = this.users.filter(u => u.id !== message.sender.id);
            this.messages.push(message);
            break;
            
          case 'user:nickname': // 用户改名
            // 更新用户列表中的昵称
            const userIndex = this.users.findIndex(u => u.id === message.sender.id);
            if (userIndex >= 0) {
              this.users[userIndex].nickname = message.sender.nickname;
            }
            this.messages.push(message);
            break;
            
          case 'user:info': // 用户信息
            this.currentUser = message.sender;
            break;
            
          case 'online:count': // 在线人数更新
            this.onlineCount = message.data;
            break;
            
          case 'chat:history': // 聊天历史
            if (Array.isArray(message.data)) {
              this.messages = message.data;
            }
            break;
            
          case 'chat:error': // 错误消息
            this.error = message.content;
            break;
        }
      } catch (error) {
        console.error('处理消息失败:', error);
      }
    },
    
    // 连接关闭事件
    handleSocketClose(event) {
      if (this.connected) {
        console.log('与聊天服务器的连接已关闭');
      }
      
      this.connected = false;
      this.socket = null;
      
      // 如果不是主动关闭，尝试重新连接
      if (!event.wasClean && this.reconnectAttempts < this.maxReconnectAttempts) {
        this.reconnect();
      }
    },
    
    // 连接错误事件
    handleSocketError(error) {
      console.error('WebSocket连接错误:', error);
      this.handleConnectionError(error);
    },
    
    // 连接错误处理
    handleConnectionError(error) {
      this.connected = false;
      this.connecting = false;
      this.error = `连接错误: ${error.message || '未知错误'}`;
      console.error('连接错误:', error);
      
      // 尝试重新连接
      this.reconnect();
    },
    
    // 添加系统消息（本地）
    addSystemMessage(content) {
      this.messages.push({
        id: `system-${Date.now()}`,
        content,
        timestamp: Date.now(),
        type: 'system',
        sender: {
          id: 'system',
          nickname: '系统'
        }
      });
    }
  }
});
```

## 3. 组件实现

### 3.1 聊天室主容器组件

```vue
<!-- src/components/chat/ChatRoom.vue -->

<template>
  <div class="chat-container">
    <!-- 聊天室头部 -->
    <div class="chat-header">
      <h2>聊天室 <span class="online-count">({{ onlineCount }}人在线)</span></h2>
      <div v-if="!isConnected" class="connect-controls">
        <button @click="connect" :disabled="connecting" class="connect-btn">
          {{ connecting ? '连接中...' : '连接聊天室' }}
        </button>
      </div>
      <div v-else class="user-controls">
        <span class="user-greeting">你好, {{ currentUser?.nickname }}</span>
        <button @click="showNicknameModal = true" class="nickname-btn">
          修改昵称
        </button>
      </div>
    </div>
    
    <!-- 聊天内容区域 -->
    <div class="chat-content">
      <!-- 消息列表 -->
      <div 
        ref="messagesContainer" 
        class="messages-container"
        @scroll="handleScroll"
      >
        <!-- 加载历史按钮 -->
        <div v-if="hasMoreHistory && !isLoadingHistory" class="load-history">
          <button @click="loadMoreHistory" class="load-history-btn">
            加载更多历史消息
          </button>
        </div>
        
        <!-- 加载中提示 -->
        <div v-if="isLoadingHistory" class="loading-history">
          <div class="loading-spinner"></div>
          <span>加载历史消息中...</span>
        </div>
        
        <!-- 未连接提示 -->
        <div v-if="!isConnected" class="connect-message">
          <div v-if="error" class="error-message">{{ error }}</div>
          <div v-else>
            请点击"连接聊天室"按钮开始聊天
          </div>
        </div>
        
        <!-- 消息列表 -->
        <ChatMessage
          v-for="message in formattedMessages"
          :key="message.id"
          :message="message"
          :is-current-user="message.sender?.id === currentUser?.id"
        />
      </div>
      
      <!-- 用户列表侧边栏 -->
      <UserList
        :users="sortedUsers"
        :current-user="currentUser"
      />
    </div>
    
    <!-- 消息输入框 -->
    <ChatInput
      v-if="isConnected"
      @send-message="sendMessage"
    />
    
    <!-- 昵称修改弹窗 -->
    <div v-if="showNicknameModal" class="nickname-modal">
      <div class="nickname-modal-content">
        <h3>修改昵称</h3>
        <input
          v-model="newNickname"
          placeholder="输入新昵称"
          @keyup.enter="updateNickname"
          class="nickname-input"
        />
        <div class="modal-buttons">
          <button @click="showNicknameModal = false" class="modal-btn cancel-btn">
            取消
          </button>
          <button
            @click="updateNickname"
            :disabled="!newNickname.trim()"
            class="modal-btn confirm-btn"
          >
            确定
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted, onUnmounted, nextTick, watch } from 'vue';
import { useChatStore } from '@/stores/chat';
import ChatMessage from './ChatMessage.vue';
import ChatInput from './ChatInput.vue';
import UserList from './UserList.vue';

export default {
  name: 'ChatRoom',
  components: {
    ChatMessage,
    ChatInput,
    UserList
  },
  setup() {
    const chatStore = useChatStore();
    const messagesContainer = ref(null);
    const showNicknameModal = ref(false);
    const newNickname = ref('');
    const autoScrollEnabled = ref(true);
    
    // 计算属性
    const isConnected = computed(() => chatStore.isConnected);
    const connecting = computed(() => chatStore.connecting);
    const formattedMessages = computed(() => chatStore.formattedMessages);
    const currentUser = computed(() => chatStore.currentUser);
    const onlineCount = computed(() => chatStore.onlineCount);
    const sortedUsers = computed(() => chatStore.sortedUsers);
    const error = computed(() => chatStore.error);
    const isLoadingHistory = computed(() => chatStore.isLoadingHistory);
    const hasMoreHistory = computed(() => chatStore.hasMoreHistory);
    
    // 连接到聊天服务器
    const connect = () => {
      chatStore.connect();
    };
    
    // 发送消息
    const sendMessage = (content) => {
      chatStore.sendMessage(content);
      autoScrollEnabled.value = true; // 发送消息后启用自动滚动
    };
    
    // 更新昵称
    const updateNickname = () => {
      if (newNickname.value.trim()) {
        chatStore.setNickname(newNickname.value.trim());
        showNicknameModal.value = false;
        newNickname.value = '';
      }
    };
    
    // 加载更多历史消息
    const loadMoreHistory = () => {
      // 记录当前滚动位置
      const scrollContainer = messagesContainer.value;
      const oldScrollHeight = scrollContainer.scrollHeight;
      
      chatStore.loadChatHistory().then(() => {
        nextTick(() => {
          // 恢复滚动位置，保持相对位置不变
          const newScrollHeight = scrollContainer.scrollHeight;
          scrollContainer.scrollTop = newScrollHeight - oldScrollHeight;
        });
      });
    };
    
    // 处理消息容器滚动事件
    const handleScroll = () => {
      if (!messagesContainer.value) return;
      
      const { scrollTop, clientHeight, scrollHeight } = messagesContainer.value;
      
      // 如果用户向上滚动超过200px，禁用自动滚动
      autoScrollEnabled.value = (scrollHeight - scrollTop - clientHeight) < 200;
    };
    
    // 自动滚动到最新消息
    const scrollToBottom = () => {
      if (!messagesContainer.value || !autoScrollEnabled.value) return;
      
      nextTick(() => {
        const container = messagesContainer.value;
        container.scrollTop = container.scrollHeight;
      });
    };
    
    // 监听消息变化，自动滚动
    watch(() => chatStore.messages.length, scrollToBottom);
    
    // 生命周期钩子
    onMounted(() => {
      // 可选：自动连接
      // connect();
    });
    
    onUnmounted(() => {
      chatStore.disconnect();
    });
    
    return {
      isConnected,
      connecting,
      formattedMessages,
      currentUser,
      onlineCount,
      sortedUsers,
      error,
      isLoadingHistory,
      hasMoreHistory,
      messagesContainer,
      showNicknameModal,
      newNickname,
      connect,
      sendMessage,
      updateNickname,
      loadMoreHistory,
      handleScroll
    };
  }
};
</script>

<style scoped>
.chat-container {
  display: flex;
  flex-direction: column;
  height: 100%;
  max-width: 1200px;
  margin: 0 auto;
  border-radius: 8px;
  overflow: hidden;
  background-color: var(--color-background);
  box-shadow: 0 0 20px rgba(0, 0, 0, 0.1);
}

.chat-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px 20px;
  background-color: var(--color-primary);
  color: white;
}

.chat-header h2 {
  margin: 0;
  font-size: 1.5rem;
}

.online-count {
  font-size: 0.9rem;
  font-weight: normal;
}

.user-controls {
  display: flex;
  align-items: center;
  gap: 10px;
}

.user-greeting {
  font-weight: 500;
}

.nickname-btn {
  background-color: rgba(255, 255, 255, 0.2);
  border: none;
  color: white;
  padding: 5px 10px;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.nickname-btn:hover {
  background-color: rgba(255, 255, 255, 0.3);
}

.chat-content {
  display: flex;
  flex: 1;
  min-height: 0; /* 重要：防止Flex子元素溢出 */
}

.messages-container {
  flex: 1;
  padding: 20px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.connect-message {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
  color: var(--color-text-light);
  font-size: 1.2rem;
}

.error-message {
  color: var(--color-danger);
  text-align: center;
}

.load-history {
  display: flex;
  justify-content: center;
  margin-bottom: 15px;
}

.load-history-btn {
  background-color: var(--color-background-soft);
  border: none;
  padding: 8px 15px;
  border-radius: 20px;
  cursor: pointer;
  font-size: 0.9rem;
  color: var(--color-text);
  transition: background-color 0.2s;
}

.load-history-btn:hover {
  background-color: var(--color-background-mute);
}

.loading-history {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 10px;
  margin-bottom: 15px;
  color: var(--color-text-light);
  font-size: 0.9rem;
}

.loading-spinner {
  width: 20px;
  height: 20px;
  border: 2px solid transparent;
  border-top-color: var(--color-primary);
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.connect-btn {
  background-color: var(--color-primary);
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.connect-btn:hover:not(:disabled) {
  background-color: var(--color-primary-dark);
}

.connect-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.nickname-modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.nickname-modal-content {
  background-color: var(--color-background);
  padding: 20px;
  border-radius: 8px;
  width: 300px;
}

.nickname-modal-content h3 {
  margin-top: 0;
  margin-bottom: 15px;
  color: var(--color-text);
}

.nickname-input {
  width: 100%;
  padding: 10px;
  margin-bottom: 15px;
  border: 1px solid var(--color-border);
  border-radius: 4px;
  box-sizing: border-box;
  font-size: 1rem;
}

.modal-buttons {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

.modal-btn {
  padding: 8px 15px;
  border-radius: 4px;
  border: none;
  cursor: pointer;
  font-size: 0.9rem;
}

.cancel-btn {
  background-color: var(--color-text-light);
  color: white;
}

.confirm-btn {
  background-color: var(--color-primary);
  color: white;
}

.confirm-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .chat-content {
    flex-direction: column;
  }
  
  .user-greeting {
    display: none;
  }
}
</style>
```

### 3.2 单条消息组件

```vue
<!-- src/components/chat/ChatMessage.vue -->

<template>
  <div
    class="message-container"
    :class="{
      'system': message.type === 'system',
      'user-message': message.type === 'text',
      'current-user': isCurrentUser,
      'user-join': message.type === 'user:join',
      'user-leave': message.type === 'user:leave',
      'user-nickname': message.type === 'user:nickname'
    }"
  >
    <!-- 系统消息和用户事件消息 -->
    <div v-if="message.type !== 'text'" class="system-message">
      {{ message.content }}
      <span class="message-time">{{ message.formattedTime }}</span>
    </div>
    
    <!-- 聊天消息 -->
    <template v-else>
      <div class="message-header">
        <span class="message-sender">{{ message.sender.nickname }}</span>
        <span class="message-time">{{ message.formattedTime }}</span>
      </div>
      <div class="message-content" v-html="formatContent(message.content)"></div>
    </template>
  </div>
</template>

<script>
import { processEmojis } from '@/utils/emoji';

export default {
  name: 'ChatMessage',
  props: {
    message: {
      type: Object,
      required: true
    },
    isCurrentUser: {
      type: Boolean,
      default: false
    }
  },
  methods: {
    formatContent(content) {
      // 处理文本中的emoji
      return processEmojis(content);
    }
  }
};
</script>

<style scoped>
.message-container {
  margin-bottom: 5px;
  max-width: 85%;
  padding: 10px 12px;
  border-radius: 8px;
  position: relative;
  animation: fadeIn 0.3s ease-out;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

.user-message {
  background-color: var(--color-background-soft);
  align-self: flex-start;
}

.current-user {
  background-color: var(--color-primary-light);
  align-self: flex-end;
  color: var(--color-primary-contrast);
}

.system-message, .user-join, .user-leave, .user-nickname {
  background-color: rgba(0, 0, 0, 0.05);
  color: var(--color-text-light);
  align-self: center;
  font-style: italic;
  padding: 5px 10px;
  font-size: 0.9rem;
  border-radius: 15px;
  max-width: 70%;
  text-align: center;
}

.user-join {
  background-color: rgba(0, 128, 0, 0.1);
}

.user-leave {
  background-color: rgba(128, 0, 0, 0.1);
}

.user-nickname {
  background-color: rgba(0, 0, 128, 0.1);
}

.message-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 5px;
}

.message-sender {
  font-weight: bold;
  font-size: 0.9rem;
}

.message-time {
  font-size: 0.75rem;
  color: var(--color-text-light);
  margin-left: 10px;
}

.message-content {
  word-break: break-word;
  line-height: 1.4;
}

.system-message .message-time {
  margin-left: 5px;
  font-size: 0.75rem;
}

/* 表情符号样式 */
.emoji {
  font-size: 1.2em;
  vertical-align: middle;
}
</style>
``` 