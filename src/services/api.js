/**
 * API服务 - 集中管理所有后端请求
 */

// 获取基础URL
const getBaseUrl = () => {
  const protocol = window.location.protocol;
  const hostname = window.location.hostname;
  return `${protocol}//${hostname}`;
};

// WebSocket相关
export const getWebSocketUrl = (username, roomId) => {
  const wsProtocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
  return `${wsProtocol}//${window.location.hostname}/ws?username=${encodeURIComponent(username)}&roomId=${encodeURIComponent(roomId)}`;
};

// 房间相关API
export const roomApi = {
  // 获取所有房间列表
  getAllRooms: async () => {
    const response = await fetch(`${getBaseUrl()}/api/rooms`);
    return response.json();
  },
  
  // 获取特定房间的历史消息
  getRoomMessages: async (roomId, offset, limit) => {
    const response = await fetch(`${getBaseUrl()}/api/rooms/${roomId}/messages?offset=${offset}&limit=${limit}`);
    return response.json();
  }
};

// 版本信息API
export const getVersion = async () => {
  const response = await fetch(`${getBaseUrl()}/api/version`);
  return response.text();
};

// 可以根据需要继续添加其他API分类 