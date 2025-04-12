/**
 * 文章相关工具函数
 */

/**
 * 估算文章阅读时间
 * @param {string} content 文章内容
 * @param {number} wordsPerMinute 平均阅读速度(每分钟字数)
 * @returns {number} 预计阅读时间(分钟)
 */
export const estimateReadingTime = (content, wordsPerMinute = 300) => {
  if (!content) return 1;
  
  // 移除Markdown语法和HTML标签
  const cleanText = content
    .replace(/!\[.*?\]\(.*?\)/g, '') // 图片
    .replace(/\[.*?\]\(.*?\)/g, '$1') // 链接
    .replace(/#{1,6}\s+/g, '') // 标题
    .replace(/```[\s\S]*?```/g, '') // 代码块
    .replace(/`.*?`/g, '') // 行内代码
    .replace(/<[^>]*>/g, ''); // HTML标签
  
  // 汉字计数(一个汉字按两个单词计算)
  const chineseCharCount = (cleanText.match(/[\u4e00-\u9fa5]/g) || []).length;
  
  // 非汉字单词计数
  const wordCount = cleanText
    .replace(/[\u4e00-\u9fa5]/g, '') // 移除汉字
    .trim()
    .split(/\s+/)
    .filter(Boolean).length;
  
  // 总计词数(汉字计为2个词)
  const totalWordCount = wordCount + (chineseCharCount * 2);
  
  // 计算阅读时间(分钟)
  const minutes = totalWordCount / wordsPerMinute;
  
  // 至少返回1分钟
  return Math.max(1, Math.round(minutes));
};

/**
 * 格式化阅读时间
 * @param {number} minutes 阅读时间(分钟)
 * @returns {string} 格式化后的阅读时间
 */
export const formatReadingTime = (minutes) => {
  if (minutes < 1) return '不到1分钟';
  if (minutes === 1) return '1分钟';
  if (minutes < 60) return `${minutes}分钟`;
  
  const hours = Math.floor(minutes / 60);
  const remainingMinutes = Math.round(minutes % 60);
  
  if (remainingMinutes === 0) {
    return `${hours}小时`;
  }
  
  return `${hours}小时${remainingMinutes}分钟`;
}; 