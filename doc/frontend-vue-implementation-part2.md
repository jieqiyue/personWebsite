# Vue前端实现设计（续）

## 3. 组件实现（续）

### 3.3 消息输入组件

```vue
<!-- src/components/chat/ChatInput.vue -->

<template>
  <div class="chat-input-container">
    <!-- 表情按钮 -->
    <button class="emoji-button" @click="toggleEmojiPicker" type="button">
      😊
    </button>
    
    <!-- 消息输入区域 -->
    <textarea
      ref="inputField"
      v-model="message"
      placeholder="输入消息..."
      @keydown.enter.prevent="sendMessage"
      @input="adjustHeight"
      rows="1"
      class="message-input"
    ></textarea>
    
    <!-- 发送按钮 -->
    <button
      class="send-button"
      @click="sendMessage"
      :disabled="!message.trim()"
      type="button"
    >
      发送
    </button>
    
    <!-- 表情选择器 -->
    <div v-if="showEmojiPicker" class="emoji-picker-container">
      <EmojiPicker @select="addEmoji" @close="showEmojiPicker = false" />
    </div>
  </div>
</template>

<script>
import { ref } from 'vue';
import EmojiPicker from './EmojiPicker.vue';

export default {
  name: 'ChatInput',
  components: {
    EmojiPicker
  },
  emits: ['send-message'],
  setup(props, { emit }) {
    const message = ref('');
    const inputField = ref(null);
    const showEmojiPicker = ref(false);
    
    // 发送消息
    const sendMessage = () => {
      if (message.value.trim()) {
        emit('send-message', message.value);
        message.value = '';
        
        // 重置输入框高度
        if (inputField.value) {
          inputField.value.style.height = 'auto';
        }
        
        // 聚焦输入框
        inputField.value.focus();
      }
    };
    
    // 自动调整输入框高度
    const adjustHeight = () => {
      const el = inputField.value;
      if (!el) return;
      
      // 重置高度以便得到正确的scrollHeight
      el.style.height = 'auto';
      
      // 使用scrollHeight调整高度，最大高度200px
      el.style.height = `${Math.min(el.scrollHeight, 200)}px`;
    };
    
    // 切换表情选择器
    const toggleEmojiPicker = () => {
      showEmojiPicker.value = !showEmojiPicker.value;
    };
    
    // 添加表情到消息
    const addEmoji = (emoji) => {
      message.value += emoji;
      showEmojiPicker.value = false;
      inputField.value.focus();
      
      // 调整高度
      adjustHeight();
    };
    
    return {
      message,
      inputField,
      showEmojiPicker,
      sendMessage,
      adjustHeight,
      toggleEmojiPicker,
      addEmoji
    };
  }
};
</script>

<style scoped>
.chat-input-container {
  display: flex;
  align-items: flex-end;
  padding: 10px 15px;
  background-color: var(--color-background-soft);
  border-top: 1px solid rgba(0, 0, 0, 0.1);
  position: relative;
}

.message-input {
  flex: 1;
  resize: none;
  border: 1px solid #ddd;
  border-radius: 20px;
  padding: 10px 15px;
  font-family: inherit;
  font-size: 1rem;
  outline: none;
  max-height: 200px;
  overflow-y: auto;
  line-height: 1.5;
  transition: border-color 0.2s;
}

.message-input:focus {
  border-color: var(--color-primary);
}

.emoji-button,
.send-button {
  background: none;
  border: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: transform 0.2s;
}

.emoji-button {
  font-size: 1.3rem;
  margin-right: 8px;
  padding: 5px;
  border-radius: 50%;
}

.emoji-button:hover {
  transform: scale(1.1);
  background-color: rgba(0, 0, 0, 0.05);
}

.send-button {
  background-color: var(--color-primary);
  color: white;
  border-radius: 50%;
  width: 36px;
  height: 36px;
  margin-left: 8px;
}

.send-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.send-button:hover:not(:disabled) {
  transform: scale(1.1);
}

.emoji-picker-container {
  position: absolute;
  bottom: 100%;
  left: 10px;
  z-index: 10;
  margin-bottom: 5px;
}

@media (max-width: 768px) {
  .message-input {
    font-size: 0.9rem;
    padding: 8px 12px;
  }
  
  .emoji-button {
    font-size: 1.1rem;
  }
  
  .send-button {
    width: 32px;
    height: 32px;
  }
}
</style>
```

### 3.4 表情选择器组件

```vue
<!-- src/components/chat/EmojiPicker.vue -->

<template>
  <div class="emoji-picker" @click.stop="">
    <div class="emoji-header">
      <input
        type="text"
        v-model="searchQuery"
        placeholder="搜索表情..."
        class="emoji-search"
      />
      <button @click="$emit('close')" class="close-button">×</button>
    </div>
    
    <div class="emoji-categories">
      <button
        v-for="category in categories"
        :key="category.id"
        @click="currentCategory = category.id"
        class="category-button"
        :class="{ active: currentCategory === category.id }"
        :title="category.name"
      >
        {{ category.icon }}
      </button>
    </div>
    
    <div class="emoji-grid">
      <button
        v-for="emoji in filteredEmojis"
        :key="emoji.emoji"
        @click="selectEmoji(emoji.emoji)"
        class="emoji-item"
        :title="emoji.description"
      >
        {{ emoji.emoji }}
      </button>
    </div>
    
    <div class="emoji-footer">
      <div class="frequently-used">
        <button
          v-for="emoji in frequentlyUsed"
          :key="emoji"
          @click="selectEmoji(emoji)"
          class="emoji-item"
        >
          {{ emoji }}
        </button>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted, onBeforeUnmount } from 'vue';
import { emojiCategories, allEmojis } from '@/utils/emoji.js';

export default {
  name: 'EmojiPicker',
  emits: ['select', 'close'],
  setup(props, { emit }) {
    const searchQuery = ref('');
    const currentCategory = ref('smileys');
    const frequentlyUsed = ref(['😊', '👍', '❤️', '😂', '🎉', '🔥', '👋', '😎']);
    
    // 表情分类
    const categories = emojiCategories;
    
    // 根据搜索和分类过滤的表情
    const filteredEmojis = computed(() => {
      let emojis = allEmojis;
      
      // 如果有搜索查询，按描述过滤
      if (searchQuery.value) {
        const query = searchQuery.value.toLowerCase();
        return emojis.filter(emoji =>
          emoji.description.toLowerCase().includes(query)
        );
      }
      
      // 按分类过滤
      return emojis.filter(emoji => emoji.category === currentCategory.value);
    });
    
    // 选择表情
    const selectEmoji = (emoji) => {
      emit('select', emoji);
      
      // 更新常用表情
      if (!frequentlyUsed.value.includes(emoji)) {
        frequentlyUsed.value.pop(); // 移除最后一个
        frequentlyUsed.value.unshift(emoji); // 添加到开头
        
        // 存储到localStorage
        localStorage.setItem('frequentEmojis', JSON.stringify(frequentlyUsed.value));
      }
    };
    
    // 点击外部关闭选择器
    const handleClickOutside = (event) => {
      emit('close');
    };
    
    // 从localStorage加载常用表情
    const loadFrequentlyUsed = () => {
      const stored = localStorage.getItem('frequentEmojis');
      if (stored) {
        try {
          frequentlyUsed.value = JSON.parse(stored);
        } catch (e) {
          console.error('解析常用表情失败:', e);
        }
      }
    };
    
    onMounted(() => {
      loadFrequentlyUsed();
      document.addEventListener('click', handleClickOutside);
    });
    
    onBeforeUnmount(() => {
      document.removeEventListener('click', handleClickOutside);
    });
    
    return {
      searchQuery,
      currentCategory,
      categories,
      filteredEmojis,
      frequentlyUsed,
      selectEmoji
    };
  }
};
</script>

<style scoped>
.emoji-picker {
  background-color: var(--color-background);
  border-radius: 8px;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.2);
  width: 320px;
  max-height: 400px;
  display: flex;
  flex-direction: column;
}

.emoji-header {
  display: flex;
  padding: 10px;
  border-bottom: 1px solid var(--color-border);
}

.emoji-search {
  flex: 1;
  padding: 6px 10px;
  border: 1px solid var(--color-border);
  border-radius: 4px;
  font-size: 0.9rem;
}

.close-button {
  background: none;
  border: none;
  font-size: 1.2rem;
  cursor: pointer;
  margin-left: 8px;
  color: var(--color-text-light);
}

.emoji-categories {
  display: flex;
  overflow-x: auto;
  padding: 8px;
  border-bottom: 1px solid var(--color-border);
}

.category-button {
  background: none;
  border: none;
  font-size: 1.2rem;
  cursor: pointer;
  padding: 5px;
  border-radius: 4px;
  min-width: 32px;
}

.category-button.active {
  background-color: var(--color-background-mute);
}

.emoji-grid {
  display: grid;
  grid-template-columns: repeat(8, 1fr);
  gap: 5px;
  padding: 10px;
  overflow-y: auto;
  max-height: 250px;
}

.emoji-item {
  font-size: 1.5rem;
  background: none;
  border: none;
  padding: 5px;
  cursor: pointer;
  border-radius: 4px;
  transition: background-color 0.2s, transform 0.1s;
  display: flex;
  align-items: center;
  justify-content: center;
}

.emoji-item:hover {
  background-color: var(--color-background-soft);
  transform: scale(1.1);
}

.emoji-footer {
  padding: 8px;
  border-top: 1px solid var(--color-border);
}

.frequently-used {
  display: flex;
  overflow-x: auto;
  padding-bottom: 5px;
}

.frequently-used .emoji-item {
  font-size: 1.2rem;
}
</style>
```

### 3.5 用户列表组件

```vue
<!-- src/components/chat/UserList.vue -->

<template>
  <div class="users-sidebar" :class="{ 'expanded': expanded }">
    <div class="users-header" @click="expanded = !expanded">
      <h3>在线用户 ({{ users.length }})</h3>
      <button class="toggle-button">
        {{ expanded ? '收起' : '展开' }}
      </button>
    </div>
    
    <ul v-if="expanded" class="users-list">
      <li
        v-for="user in users"
        :key="user.id"
        :class="{ 'current-user': user.id === currentUser?.id }"
      >
        <span class="user-status" :title="getActiveStatus(user)"></span>
        <span class="user-nickname">
          {{ user.nickname }}
          <span v-if="user.id === currentUser?.id">(你)</span>
        </span>
      </li>
    </ul>
  </div>
</template>

<script>
import { ref, computed } from 'vue';
import { formatDistanceToNow } from 'date-fns';
import { zhCN } from 'date-fns/locale';

export default {
  name: 'UserList',
  props: {
    users: {
      type: Array,
      required: true
    },
    currentUser: {
      type: Object,
      default: null
    }
  },
  setup(props) {
    const expanded = ref(true);
    
    // 获取用户活跃状态
    const getActiveStatus = (user) => {
      const lastActive = new Date(user.lastActive);
      return `最后活跃: ${formatDistanceToNow(lastActive, { 
        addSuffix: true,
        locale: zhCN
      })}`;
    };
    
    return {
      expanded,
      getActiveStatus
    };
  }
};
</script>

<style scoped>
.users-sidebar {
  width: 220px;
  background-color: var(--color-background-soft);
  border-left: 1px solid var(--color-border);
  display: flex;
  flex-direction: column;
  transition: width 0.3s;
}

.users-header {
  padding: 15px;
  border-bottom: 1px solid var(--color-border);
  display: flex;
  justify-content: space-between;
  align-items: center;
  cursor: pointer;
}

.users-header h3 {
  margin: 0;
  font-size: 1rem;
  color: var(--color-text);
}

.toggle-button {
  background: none;
  border: none;
  font-size: 0.8rem;
  color: var(--color-text-light);
  cursor: pointer;
}

.users-list {
  list-style: none;
  padding: 0;
  margin: 0;
  overflow-y: auto;
  flex: 1;
}

.users-list li {
  padding: 10px 15px;
  display: flex;
  align-items: center;
  transition: background-color 0.2s;
}

.users-list li:hover {
  background-color: rgba(0, 0, 0, 0.03);
}

.users-list .current-user {
  background-color: rgba(0, 0, 0, 0.05);
  font-weight: bold;
}

.user-status {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background-color: #4CAF50;
  margin-right: 10px;
  display: inline-block;
}

.user-nickname {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .users-sidebar {
    width: 100%;
    border-left: none;
    border-top: 1px solid var(--color-border);
  }
  
  .users-sidebar:not(.expanded) {
    height: 50px;
  }
  
  .expanded {
    height: 200px;
  }
}
</style>
```

### 3.6 在线人数计数器组件

```vue
<!-- src/components/common/OnlineCounter.vue -->

<template>
  <div 
    class="online-counter" 
    :class="{ 'connected': connected }"
    @click="onClick"
  >
    <span class="counter-icon">👥</span>
    <span class="counter-number">{{ count }}</span>
  </div>
</template>

<script>
import { computed, onMounted, onUnmounted } from 'vue';
import { useChatStore } from '@/stores/chat';
import { useRouter } from 'vue-router';

export default {
  name: 'OnlineCounter',
  props: {
    autoConnect: {
      type: Boolean,
      default: false
    }
  },
  setup(props) {
    const chatStore = useChatStore();
    const router = useRouter();
    
    const count = computed(() => chatStore.onlineCount);
    const connected = computed(() => chatStore.connected);
    
    // 点击计数器导航到聊天页面
    const onClick = () => {
      router.push('/chat');
    };
    
    onMounted(() => {
      // 如果设置了自动连接，则连接WebSocket
      if (props.autoConnect) {
        chatStore.connect();
      } else if (!chatStore.socket) {
        // 否则只获取在线人数
        chatStore.fetchOnlineCount();
      }
    });
    
    // 选择是否在组件卸载时断开连接
    onUnmounted(() => {
      // 如果是自动连接模式，则断开连接
      // 否则保持连接状态，让其他组件决定何时断开
      if (props.autoConnect) {
        chatStore.disconnect();
      }
    });
    
    return {
      count,
      connected,
      onClick
    };
  }
};
</script>

<style scoped>
.online-counter {
  display: flex;
  align-items: center;
  font-size: 0.9rem;
  color: var(--color-text-light);
  gap: 5px;
  cursor: pointer;
  transition: color 0.3s;
  padding: 4px 8px;
  border-radius: 12px;
}

.online-counter:hover {
  background-color: rgba(0, 0, 0, 0.05);
}

.online-counter.connected {
  color: var(--color-primary);
}

.counter-icon {
  font-size: 1.1rem;
}

.counter-number {
  font-weight: 500;
}
</style>
```

## 4. 视图和路由集成

### 4.1 聊天页面视图

```vue
<!-- src/views/ChatView.vue -->

<template>
  <div class="chat-view">
    <ChatRoom />
  </div>
</template>

<script>
import ChatRoom from '@/components/chat/ChatRoom.vue';

export default {
  name: 'ChatView',
  components: {
    ChatRoom
  }
};
</script>

<style scoped>
.chat-view {
  height: 100%;
  padding: 20px;
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
}

@media (max-width: 768px) {
  .chat-view {
    padding: 10px;
  }
}
</style>
```

### 4.2 路由配置

更新路由配置以添加聊天页面：

```javascript
// src/router/index.js

import { createRouter, createWebHistory } from 'vue-router';
import HomeView from '../views/HomeView.vue';
import ChatView from '../views/ChatView.vue';

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView
    },
    {
      path: '/chat',
      name: 'chat',
      component: ChatView
    },
    // 其他现有路由...
  ]
});

export default router;
```

### 4.3 导航栏集成

在应用导航栏中添加聊天入口和在线人数：

```vue
<!-- 在主导航组件中添加 -->
<template>
  <nav class="navbar">
    <!-- 其他导航项 -->
    <div class="nav-links">
      <router-link to="/">首页</router-link>
      <router-link to="/articles">文章</router-link>
      <router-link to="/chat" class="chat-link">
        聊天室
        <OnlineCounter />
      </router-link>
      <!-- 其他链接 -->
    </div>
  </nav>
</template>

<script>
import OnlineCounter from '@/components/common/OnlineCounter.vue';

export default {
  components: {
    OnlineCounter
  }
};
</script>

<style scoped>
.chat-link {
  display: flex;
  align-items: center;
  gap: 5px;
}
</style>
```

## 5. 工具函数

### 5.1 表情工具函数

```javascript
// src/utils/emoji.js

// 表情分类
export const emojiCategories = [
  { id: 'smileys', name: '表情笑脸', icon: '😀' },
  { id: 'people', name: '人物', icon: '👨' },
  { id: 'animals', name: '动物', icon: '🐶' },
  { id: 'food', name: '食物', icon: '🍎' },
  { id: 'activities', name: '活动', icon: '⚽' },
  { id: 'travel', name: '旅行', icon: '🚗' },
  { id: 'objects', name: '物品', icon: '💡' },
  { id: 'symbols', name: '符号', icon: '❤️' },
  { id: 'flags', name: '旗帜', icon: '🏁' }
];

// 表情列表（简化版）
export const allEmojis = [
  // 笑脸表情
  { emoji: '😀', description: '笑脸', category: 'smileys' },
  { emoji: '😃', description: '大笑', category: 'smileys' },
  { emoji: '😄', description: '笑眼', category: 'smileys' },
  { emoji: '😁', description: '咧嘴笑', category: 'smileys' },
  // ...更多表情
];

// 常见表情符号的正则表达式
const emojiRegex = /[\u{1F300}-\u{1F6FF}\u{2600}-\u{26FF}]/gu;

// 处理文本中的表情符号
export function processEmojis(text) {
  // 将文本中的表情符号包装在span中以便样式化
  return text.replace(emojiRegex, match => 
    `<span class="emoji">${match}</span>`
  );
}
```

### 5.2 WebSocket重连工具

```javascript
// src/utils/websocket.js

/**
 * 创建一个带有自动重连功能的WebSocket
 * @param {string} url WebSocket连接URL
 * @param {Object} options 配置项
 * @returns {Object} WebSocket实例和控制方法
 */
export function createReconnectingWebSocket(url, options = {}) {
  const {
    maxReconnectAttempts = 5,
    reconnectInterval = 1000,
    reconnectDecay = 1.5,
    maxReconnectInterval = 30000,
    onOpen,
    onMessage,
    onClose,
    onError,
    onReconnect,
    onMaxReconnectsExceeded
  } = options;
  
  let socket = null;
  let reconnectAttempts = 0;
  let reconnectTimeout = null;
  let forceClosed = false;
  
  // 创建WebSocket连接
  const connect = () => {
    // 清除之前的重连定时器
    if (reconnectTimeout) {
      clearTimeout(reconnectTimeout);
      reconnectTimeout = null;
    }
    
    // 创建新的WebSocket
    socket = new WebSocket(url);
    
    // 连接成功
    socket.onopen = (event) => {
      reconnectAttempts = 0;
      if (onOpen) onOpen(event);
    };
    
    // 接收消息
    socket.onmessage = (event) => {
      if (onMessage) onMessage(event);
    };
    
    // 连接关闭
    socket.onclose = (event) => {
      if (onClose) onClose(event);
      
      // 如果不是强制关闭，尝试重连
      if (!forceClosed && reconnectAttempts < maxReconnectAttempts) {
        attemptReconnect();
      } else if (!forceClosed) {
        // 达到最大重连次数
        if (onMaxReconnectsExceeded) onMaxReconnectsExceeded();
      }
    };
    
    // 连接错误
    socket.onerror = (error) => {
      if (onError) onError(error);
    };
  };
  
  // 尝试重连
  const attemptReconnect = () => {
    reconnectAttempts++;
    if (onReconnect) onReconnect(reconnectAttempts);
    
    // 计算下次重连间隔（指数退避）
    const delay = Math.min(
      reconnectInterval * Math.pow(reconnectDecay, reconnectAttempts - 1),
      maxReconnectInterval
    );
    
    // 设置重连定时器
    reconnectTimeout = setTimeout(connect, delay);
  };
  
  // 关闭连接
  const close = () => {
    forceClosed = true;
    if (reconnectTimeout) {
      clearTimeout(reconnectTimeout);
      reconnectTimeout = null;
    }
    
    if (socket) {
      socket.close();
    }
  };
  
  // 发送消息
  const send = (data) => {
    if (socket && socket.readyState === WebSocket.OPEN) {
      if (typeof data !== 'string') {
        data = JSON.stringify(data);
      }
      socket.send(data);
      return true;
    }
    return false;
  };
  
  // 立即连接
  connect();
  
  // 返回WebSocket实例和控制方法
  return {
    get socket() { return socket; },
    get reconnectAttempts() { return reconnectAttempts; },
    send,
    close,
    connect
  };
}
```

## 6. 总结和部署注意事项

### 6.1 总结

本文档详细描述了聊天室和在线人数统计功能的Vue前端实现方案，包括：

1. Pinia状态管理
2. WebSocket通信
3. 聊天界面组件
4. 表情选择器
5. 在线人数统计
6. 网络和安全处理

### 6.2 部署注意事项

1. 确保设置正确的环境变量：

```
# .env.production
VITE_WS_URL=wss://chat-api.yourdomain.com/ws
VITE_API_URL=https://chat-api.yourdomain.com/api
```

2. 确保WebSocket服务器支持SSL (wss://)

3. 确保CORS配置正确，允许前端域名访问

4. 压缩静态资源以提高加载性能

5. 考虑使用服务工作线程(Service Worker)实现离线功能 