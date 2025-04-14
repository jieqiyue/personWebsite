<template>
  <div class="chat-window" :class="{ 'is-open': isOpen }">
    <div class="chat-header">
      <h3>聊天室</h3>
      <div class="chat-controls">
        <button class="minimize-btn" @click="$emit('toggle-chat')">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <line x1="18" y1="6" x2="6" y2="18"></line>
            <line x1="6" y1="6" x2="18" y2="18"></line>
          </svg>
        </button>
      </div>
    </div>
    
    <div class="chat-body">
      <div class="connection-status" :class="{ 'connected': connected }">
        <span v-if="connected">已连接</span>
        <span v-else>未连接</span>
      </div>
      
      <div class="message-container" ref="messageContainer">
        <div v-if="messages.length === 0" class="no-messages">
          <p>欢迎来到聊天室，开始发送消息吧！</p>
        </div>
        
        <div 
          v-for="(message, index) in messages" 
          :key="index" 
          class="message" 
          :class="{ 
            'system-message': message.type === 'system' || message.type === 'join' || message.type === 'leave',
            'my-message': message.senderId === currentUserId,
            'other-message': message.senderId !== currentUserId && message.type !== 'system' && message.type !== 'join' && message.type !== 'leave'
          }"
        >
          <div v-if="message.type === 'system' || message.type === 'join' || message.type === 'leave'" class="system-content">
            {{ message.content }}
          </div>
          <template v-else>
            <div class="message-sender">{{ getSenderName(message) }}</div>
            <div class="message-content">{{ message.content }}</div>
            <div class="message-time">{{ formatTime(message.timestamp) }}</div>
          </template>
        </div>
      </div>
    </div>
    
    <div class="chat-footer">
      <select v-model="selectedRoom" class="room-selector">
        <option disabled value="">请选择聊天室</option>
        <option v-for="room in rooms" :key="room.id" :value="room.id">
          {{ room.name }} ({{ roomUserCounts[room.id] || 0 }}人在线)
        </option>
      </select>
      
      <div class="message-input-container">
        <input 
          type="text" 
          v-model="messageText" 
          @keyup.enter="sendMessage" 
          placeholder="输入消息..." 
          :disabled="!connected || !selectedRoom"
          class="message-input"
        />
        <button 
          @click="sendMessage" 
          class="send-button" 
          :disabled="!connected || !selectedRoom || !messageText.trim()"
        >
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <line x1="22" y1="2" x2="11" y2="13"></line>
            <polygon points="22 2 15 22 11 13 2 9 22 2"></polygon>
          </svg>
        </button>
      </div>
    </div>
    
    <div class="connect-controls">
      <div v-if="!connected" class="connect-form">
        <input type="text" v-model="username" placeholder="您的昵称" class="username-input" />
        <button @click="connect" class="connect-button" :disabled="!selectedRoom">连接</button>
      </div>
      <button v-else @click="disconnect" class="disconnect-button">断开连接</button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, computed, nextTick, watch } from 'vue';

// 定义属性和事件
const props = defineProps({
  isOpen: {
    type: Boolean,
    default: false
  }
});

const emit = defineEmits(['toggle-chat']);

// 状态变量
const connected = ref(false);
const messages = ref([]);
const messageText = ref('');
const selectedRoom = ref('');
const rooms = ref([]);
const username = ref(`游客${Math.floor(Math.random() * 10000)}`);
const currentUserId = ref('');
const ws = ref(null);
const messageContainer = ref(null);

// 添加在线人数计数器
const roomUserCounts = ref({}); // 房间ID -> 在线人数的映射

// 聊天功能
const connect = async () => {
  try {
    // 确保已经选择了房间
    if (!selectedRoom.value) {
      addMessage({
        type: 'system',
        content: '请先选择一个聊天室',
        timestamp: new Date()
      });
      return;
    }
    
    // 创建WebSocket连接
    const roomId = selectedRoom.value;
    const wsUrl = `ws://${window.location.hostname}:8080/ws?username=${encodeURIComponent(username.value)}&roomId=${encodeURIComponent(roomId)}`;
    
    ws.value = new WebSocket(wsUrl);
    
    ws.value.onopen = () => {
      connected.value = true;
      
      // 连接成功，自己已加入聊天室，但服务器会发送join消息，所以这里不需要增加计数
      addMessage({
        type: 'system',
        content: `已连接到聊天服务器，欢迎 ${username.value}`,
        timestamp: new Date()
      });
    };
    
    ws.value.onmessage = async (event) => {
      const msg = JSON.parse(event.data);
      
      // 调试输出所有接收到的消息
      console.log('收到消息:', msg);
      
      // 保存服务器分配的用户ID（无论消息类型）
      if (msg.metadata && msg.metadata.username === username.value) {
        currentUserId.value = msg.metadata.userId;
      }
      
      switch (msg.type) {
        case 'message':
        case 'text':
          // 处理普通消息
          // 检查是否是自己发送的消息，如果是则跳过（已经在本地添加过了）
          const isSelfMessage = msg.senderId === currentUserId.value;
          
          if (!isSelfMessage) {
            addMessage(msg);
          }
          break;
        
        case 'join':
          // 完整打印收到的加入消息对象，用于调试
          console.log('收到加入消息(完整对象):', JSON.stringify(msg));
          
          // 用户加入消息 - 更新所有房间人数
          if (msg.roomId && roomUserCounts.value[msg.roomId] !== undefined) {
            roomUserCounts.value[msg.roomId]++;
          }
          
          // 仅在当前选中的房间内显示加入消息
          if (msg.roomId === selectedRoom.value) {
            // 保持原始消息类型为'join'，并使用服务器发送的完整消息对象
            addMessage({
              ...msg,                                     // 保留所有原始字段
              type: 'join',                               // 确保类型为join
              content: msg.content || `用户加入了聊天室`,   // 提供默认内容
              timestamp: new Date(msg.timestamp || Date.now())
            });
          }
          break;
        
        case 'leave':
          // 完整打印收到的离开消息对象，用于调试
          console.log('收到离开消息(完整对象):', JSON.stringify(msg));
          
          // 用户离开消息 - 更新所有房间人数
          if (msg.roomId && roomUserCounts.value[msg.roomId] !== undefined && roomUserCounts.value[msg.roomId] > 0) {
            roomUserCounts.value[msg.roomId]--;
          }
          
          // 仅在当前选中的房间内显示离开消息
          if (msg.roomId === selectedRoom.value) {
            // 保持原始消息类型为'leave'，并使用服务器发送的完整消息对象
            addMessage({
              ...msg,                                     // 保留所有原始字段
              type: 'leave',                              // 确保类型为leave
              content: msg.content || `用户离开了聊天室`,   // 提供默认内容
              timestamp: new Date(msg.timestamp || Date.now())
            });
          }
          break;
          
        default:
          // 处理其他未知类型的消息
          console.log('收到未知类型消息:', msg);
          // 尝试显示未知类型消息的内容
          if (msg.content && msg.roomId === selectedRoom.value) {
            addMessage({
              type: msg.type || 'unknown',
              content: msg.content,
              senderId: msg.senderId,
              senderName: msg.senderName || (msg.metadata && msg.metadata.username) || '未知用户',
              timestamp: new Date(msg.timestamp || Date.now())
            });
          }
          break;
      }
    };
    
    ws.value.onclose = () => {
      connected.value = false;
      addMessage({
        type: 'system',
        content: '与服务器的连接已断开',
        timestamp: new Date()
      });
    };
    
    ws.value.onerror = (error) => {
      console.error('WebSocket错误:', error);
      addMessage({
        type: 'system',
        content: '连接发生错误',
        timestamp: new Date()
      });
    };
  } catch (error) {
    console.error('连接失败:', error);
    addMessage({
      type: 'system',
      content: `无法连接到聊天服务器: ${error.message}`,
      timestamp: new Date()
    });
  }
};

const disconnect = () => {
  if (ws.value) {
    // 断开前先减少当前房间的计数（只有在收不到服务器的leave消息时才有用）
    if (roomUserCounts.value[selectedRoom.value] !== undefined && roomUserCounts.value[selectedRoom.value] > 0) {
      roomUserCounts.value[selectedRoom.value]--;
    }
    
    ws.value.close();
    ws.value = null;
  }
};

const sendMessage = () => {
  if (!connected.value || !selectedRoom.value || !messageText.value.trim()) return;
  
  const message = {
    type: 'text',
    content: messageText.value,
    senderId: currentUserId.value,
    roomId: selectedRoom.value,
    timestamp: new Date()
  };
  
  ws.value.send(JSON.stringify(message));
  
  // 清空输入框
  messageText.value = '';
  
  // 添加到本地消息列表
  addMessage({
    ...message,
    isMine: true // 标记为自己的消息
  });
};

const addMessage = (message) => {
  messages.value.push(message);
  
  // 限制消息数量，防止过多消息导致性能问题
  if (messages.value.length > 100) {
    messages.value = messages.value.slice(-100);
  }
  
  // 滚动到底部
  scrollToBottom();
};

const scrollToBottom = async () => {
  await nextTick();
  if (messageContainer.value) {
    messageContainer.value.scrollTop = messageContainer.value.scrollHeight;
  }
};

const getSenderName = (message) => {
  if (message.senderId === currentUserId.value) {
    return '我';
  }
  
  return message.senderName || message.metadata?.username || '其他用户';
};

const formatTime = (timestamp) => {
  const date = new Date(timestamp);
  const year = date.getFullYear();
  const month = (date.getMonth() + 1).toString().padStart(2, '0');
  const day = date.getDate().toString().padStart(2, '0');
  const hours = date.getHours().toString().padStart(2, '0');
  const minutes = date.getMinutes().toString().padStart(2, '0');
  const seconds = date.getSeconds().toString().padStart(2, '0');
  
  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
};

const loadRooms = async (force = false) => {
  try {
    const response = await fetch(`http://${window.location.hostname}:8080/api/rooms`);
    const data = await response.json();
    
    if (data.rooms && data.rooms.length > 0) {
      rooms.value = data.rooms;
      
      // 初始化在线人数计数器
      data.rooms.forEach(room => {
        roomUserCounts.value[room.id] = room.userCount || 0;
      });
      
      // 默认选择第一个房间
      if (!selectedRoom.value) {
        selectedRoom.value = data.rooms.find(room => room.isDefault)?.id || data.rooms[0].id;
      }
    } else {
      // 如果没有房间，使用默认房间
      rooms.value = [{ id: 'general', name: '常规聊天', userCount: 0, isDefault: true }];
      selectedRoom.value = 'general';
      roomUserCounts.value['general'] = 0; // 初始化默认房间计数器
    }
  } catch (error) {
    console.error('加载聊天室列表失败:', error);
    // 添加默认房间
    rooms.value = [{ id: 'general', name: '常规聊天', userCount: 0, isDefault: true }];
    selectedRoom.value = 'general';
    roomUserCounts.value['general'] = 0; // 初始化默认房间计数器
  }
};

// 由于现在可以从服务器得到所有房间的实时更新，可以移除或减少刷新间隔
// 例如，将刷新间隔从30秒改为2分钟作为备份机制
const startRefreshInterval = () => {
  // 清除之前的定时器
  if (refreshInterval.value) {
    clearInterval(refreshInterval.value);
  }
  
  // 设置新的定时器，每2分钟刷新一次，作为备份机制
  refreshInterval.value = setInterval(() => {
    loadRooms(true); // 强制刷新所有房间信息
  }, 120000); // 2分钟
};

// 监听聊天窗口的打开状态
watch(() => props.isOpen, (newVal) => {
  if (newVal) {
    // 窗口打开时，自动滚动到底部
    scrollToBottom();
  }
});

// 组件挂载和卸载时的处理
onMounted(async () => {
  // 加载房间列表
  await loadRooms();
});

onUnmounted(() => {
  // 确保在组件卸载时关闭连接
  disconnect();
});
</script>

<style scoped>
.chat-window {
  position: fixed;
  right: 20px;
  bottom: 20px;
  width: 450px;
  height: 600px;
  background-color: var(--background, #fff);
  border: 1px solid var(--border, #e0e0e0);
  border-radius: 8px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
  display: flex;
  flex-direction: column;
  z-index: 1001;
  transform: translateY(calc(100% + 20px));
  opacity: 0;
  visibility: hidden;
  transition: all 0.3s ease;
  overflow: hidden;
}

.chat-window.is-open {
  transform: translateY(0);
  opacity: 1;
  visibility: visible;
}

.chat-header {
  padding: 12px 16px;
  background-color: var(--background-alt, #f8f8f8);
  border-bottom: 1px solid var(--border, #e0e0e0);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.chat-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: var(--primary, #333);
}

.chat-controls button {
  background: none;
  border: none;
  cursor: pointer;
  padding: 4px;
  color: var(--text-light, #888);
  display: flex;
  align-items: center;
  justify-content: center;
  transition: color 0.2s ease;
}

.chat-controls button:hover {
  color: var(--primary, #333);
}

.chat-body {
  flex: 1;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.connection-status {
  padding: 6px 12px;
  text-align: center;
  font-size: 12px;
  background-color: #ffe6e6;
  color: #cc0000;
}

.connection-status.connected {
  background-color: #e6ffe6;
  color: #007700;
}

.message-container {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.no-messages {
  text-align: center;
  color: var(--text-light, #888);
  margin: auto 0;
}

.message {
  max-width: 80%;
  padding: 8px 12px;
  border-radius: 12px;
  position: relative;
  word-break: break-word;
}

.system-message {
  width: 100%;
  max-width: 100%;
  text-align: center;
  padding: 6px 12px;
  background-color: var(--background-alt, #f5f5f5);
  color: var(--text-light, #888);
  font-size: 12px;
  border-radius: 16px;
  margin: 6px auto;
}

.my-message {
  align-self: flex-end;
  background-color: #e6f7ff;
  color: #004977;
  border-bottom-right-radius: 4px;
}

.other-message {
  align-self: flex-start;
  background-color: var(--background-alt, #f5f5f5);
  color: var(--text, #333);
  border-bottom-left-radius: 4px;
}

.message-sender {
  font-size: 14px;
  font-weight: 600;
  margin-bottom: 4px;
  color: var(--text-light, #888);
}

.message-content {
  font-size: 16px;
  line-height: 1.5;
}

.message-time {
  font-size: 12px;
  color: var(--text-lighter, #aaa);
  margin-top: 5px;
  text-align: right;
}

.chat-footer {
  padding: 12px;
  border-top: 1px solid var(--border, #e0e0e0);
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.room-selector {
  padding: 8px;
  border: 1px solid var(--border, #e0e0e0);
  border-radius: 4px;
  background-color: var(--background, #fff);
  color: var(--text, #333);
  font-size: 14px;
  width: 100%;
}

.message-input-container {
  display: flex;
  gap: 8px;
}

.message-input {
  flex: 1;
  padding: 10px 14px;
  border: 1px solid var(--border, #e0e0e0);
  border-radius: 4px;
  font-size: 16px;
  outline: none;
  transition: border-color 0.2s ease;
}

.message-input:focus {
  border-color: var(--primary, #007BFF);
}

.send-button {
  padding: 10px 14px;
  background-color: var(--primary, #007BFF);
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background-color 0.2s ease;
}

.send-button:hover {
  background-color: var(--primary-dark, #0069d9);
}

.send-button:disabled {
  background-color: var(--text-lighter, #ccc);
  cursor: not-allowed;
}

.connect-controls {
  padding: 12px;
  border-top: 1px solid var(--border, #e0e0e0);
}

.connect-form {
  display: flex;
  gap: 8px;
}

.username-input {
  flex: 1;
  padding: 8px 12px;
  border: 1px solid var(--border, #e0e0e0);
  border-radius: 4px;
  font-size: 14px;
}

.connect-button, .disconnect-button {
  padding: 8px 12px;
  border-radius: 4px;
  font-size: 14px;
  cursor: pointer;
}

.connect-button {
  background-color: var(--primary, #007BFF);
  color: white;
  border: none;
}

.disconnect-button {
  width: 100%;
  background-color: var(--background-alt, #f8f8f8);
  color: var(--text, #333);
  border: 1px solid var(--border, #e0e0e0);
}

.connect-button:hover {
  background-color: var(--primary-dark, #0069d9);
}

.disconnect-button:hover {
  background-color: #ffe6e6;
  border-color: #ffcccc;
  color: #cc0000;
}

@media (max-width: 768px) {
  .chat-window {
    right: 5px;
    bottom: 5px;
    width: calc(100% - 10px);
    height: 80vh;
    max-height: 600px;
  }
}
</style> 