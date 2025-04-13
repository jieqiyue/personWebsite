<template>
  <div class="chat-button" @click="$emit('toggle-chat')">
    <span class="chat-icon">
      <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path>
      </svg>
    </span>
    <span class="chat-text">聊天</span>
  </div>
</template>

<script setup>
// 定义事件
defineEmits(['toggle-chat']);
</script>

<style scoped>
.chat-button {
  position: fixed;
  right: 20px;
  bottom: 20px;
  padding: 10px 20px;
  background-color: var(--background-alt, #f8f8f8);
  border: 1px solid var(--border, #e0e0e0);
  border-radius: 4px;
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  z-index: 1000;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: all 0.2s ease;
  user-select: none;
}

.chat-button:hover {
  background-color: var(--background, #fff);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.chat-icon {
  display: flex;
  align-items: center;
  justify-content: center;
}

.chat-text {
  font-size: 14px;
  font-weight: 500;
  color: var(--text, #333);
}

@media (max-width: 768px) {
  .chat-button {
    right: 10px;
    bottom: 10px;
    padding: 8px 16px;
  }
}
</style> 