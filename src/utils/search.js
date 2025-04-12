/**
 * 文章搜索服务
 * 实现全文搜索功能
 */

// 搜索索引，存储所有文章的内容和元数据
let searchIndex = [];
let isIndexBuilt = false;

/**
 * 构建搜索索引
 * 从服务器获取所有文章内容并索引
 */
export const buildSearchIndex = async () => {
  if (isIndexBuilt) return;
  
  try {
    // 获取文章索引
    const indexResponse = await fetch('/markdown/articles/index.json');
    const { articles } = await indexResponse.json();
    
    // 临时存储所有文章数据的数组
    const indexData = [];
    
    // 获取每篇文章的内容
    for (const article of articles) {
      try {
        const contentResponse = await fetch(`/markdown/articles/${article.id}.md`);
        const content = await contentResponse.text();
        
        // 提取文章内容，移除YAML前置元数据
        let processedContent = content;
        if (content.trim().startsWith('---')) {
          const secondSeparatorIndex = content.indexOf('---', 3);
          if (secondSeparatorIndex !== -1) {
            processedContent = content.substring(secondSeparatorIndex + 3).trim();
          }
        }
        
        // 添加到索引
        indexData.push({
          ...article,
          content: processedContent,
          // 生成搜索文本（标题、摘要和内容的组合）
          searchText: `${article.title} ${article.excerpt || ''} ${processedContent}`.toLowerCase()
        });
      } catch (error) {
        console.error(`无法获取文章 ${article.id} 的内容:`, error);
      }
    }
    
    searchIndex = indexData;
    isIndexBuilt = true;
    console.log(`搜索索引构建完成，共索引 ${searchIndex.length} 篇文章`);
    
    return true;
  } catch (error) {
    console.error('构建搜索索引失败:', error);
    return false;
  }
};

/**
 * 搜索文章
 * @param {string} query 搜索查询
 * @param {Object} options 搜索选项
 * @param {Array} options.tags 按标签筛选
 * @param {number} options.limit 结果数量限制
 * @returns {Array} 搜索结果
 */
export const searchArticles = async (query, options = {}) => {
  // 确保索引已构建
  if (!isIndexBuilt) {
    await buildSearchIndex();
  }
  
  // 如果查询为空且没有标签筛选，返回空结果
  if ((!query || query.trim() === '') && (!options.tags || options.tags.length === 0)) {
    return [];
  }
  
  const { tags = [], limit = 20 } = options;
  const searchTerms = query ? query.toLowerCase().trim().split(/\s+/) : [];
  
  // 对每篇文章进行评分
  const results = searchIndex
    .map(article => {
      // 初始分数
      let score = 0;
      
      // 标签匹配加分
      if (tags.length > 0) {
        const tagMatches = article.tags ? article.tags.filter(tag => 
          tags.includes(tag.toLowerCase())
        ).length : 0;
        
        // 如果指定了标签但文章没有匹配，大幅降低分数
        if (tags.length > 0 && tagMatches === 0) {
          score -= 1000;
        } else {
          // 每个匹配标签加分
          score += tagMatches * 20;
        }
      }
      
      // 对每个搜索词进行评分
      for (const term of searchTerms) {
        if (term.length < 2) continue; // 忽略过短的搜索词
        
        // 标题匹配权重最高
        if (article.title.toLowerCase().includes(term)) {
          score += 50;
          
          // 标题开头匹配额外加分
          if (article.title.toLowerCase().startsWith(term)) {
            score += 20;
          }
        }
        
        // 摘要匹配次之
        if (article.excerpt && article.excerpt.toLowerCase().includes(term)) {
          score += 30;
        }
        
        // 内容匹配
        if (article.content.toLowerCase().includes(term)) {
          score += 10;
          
          // 计算匹配出现次数
          const matchCount = (article.content.toLowerCase().match(new RegExp(term, 'g')) || []).length;
          score += Math.min(matchCount, 10); // 最多额外加10分
        }
      }
      
      // 如果是空查询但有标签筛选，给所有匹配标签的文章一个基础分数
      if (searchTerms.length === 0 && tags.length > 0 && score > -1000) {
        score += 10;
      }
      
      // 日期新的文章略微加分
      const articleDate = new Date(article.date);
      const now = new Date();
      const ageInDays = (now - articleDate) / (1000 * 60 * 60 * 24);
      score += Math.max(0, 5 - Math.min(5, ageInDays / 30)); // 最近5个月内的文章加分
      
      return {
        ...article,
        score,
        // 生成摘要，突出显示匹配项
        searchExcerpt: generateSearchExcerpt(article.content, searchTerms)
      };
    })
    .filter(article => article.score > 0) // 只保留得分大于0的结果
    .sort((a, b) => b.score - a.score) // 按得分降序排序
    .slice(0, limit); // 只返回指定数量的结果
  
  return results;
};

/**
 * 生成搜索结果的摘要，突出显示匹配项
 * @param {string} content 文章内容
 * @param {Array} searchTerms 搜索词
 * @returns {string} 格式化的摘要
 */
function generateSearchExcerpt(content, searchTerms) {
  if (!content) return '';
  
  const maxLength = 200;
  let excerpt = content.substring(0, maxLength);
  
  // 如果没有搜索词，直接返回开头摘要
  if (searchTerms.length === 0 || searchTerms[0].length < 2) {
    return excerpt + (content.length > maxLength ? '...' : '');
  }
  
  // 尝试找到第一个匹配项附近的文本
  for (const term of searchTerms) {
    if (term.length < 2) continue;
    
    const termIndex = content.toLowerCase().indexOf(term);
    if (termIndex > 0) {
      // 确定摘要的起始位置，尽量使匹配词位于中间
      const start = Math.max(0, termIndex - Math.floor(maxLength / 2));
      excerpt = content.substring(start, start + maxLength);
      break;
    }
  }
  
  // 如果摘要被截断，添加省略号
  if (excerpt.length < content.length) {
    excerpt = excerpt + '...';
  }
  
  return excerpt;
}

/**
 * 获取搜索索引的统计信息
 * @returns {Object} 索引统计信息
 */
export const getSearchStats = () => {
  return {
    isIndexBuilt,
    articleCount: searchIndex.length,
  };
};

/**
 * 预热搜索索引(在应用启动时调用)
 */
export const preloadSearchIndex = () => {
  setTimeout(() => {
    buildSearchIndex();
  }, 2000); // 延迟2秒后构建索引，避免影响页面初始加载
}; 