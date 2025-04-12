let cachedArticles = null;
let cachePromise = null;

/**
 * 获取文章列表数据，优先从缓存读取。
 * 如果缓存不存在，则发起网络请求并缓存结果。
 * 处理并发请求，确保只发起一次请求。
 * 
 * @returns {Promise<Array>} 包含文章对象的数组的Promise
 */
export const getArticles = async () => {
  // 如果已有缓存，直接返回
  if (cachedArticles) {
    console.log('从缓存加载文章列表');
    return Promise.resolve(cachedArticles);
  }

  // 如果正在请求中，返回该请求的Promise
  if (cachePromise) {
    console.log('等待正在进行的文章列表请求');
    return cachePromise;
  }

  // 发起新的请求
  console.log('首次请求文章列表');
  cachePromise = new Promise(async (resolve, reject) => {
    try {
      const response = await fetch('/markdown/articles/index.json');
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      const data = await response.json();
      
      // 缓存结果
      cachedArticles = data.articles;
      console.log(`文章列表缓存成功，共 ${cachedArticles.length} 篇`);
      resolve(cachedArticles);
    } catch (err) {
      console.error('加载或缓存文章列表失败:', err);
      // 清理Promise引用，允许重试
      cachePromise = null;
      reject('加载文章列表失败，请稍后重试'); 
    } 
  });

  return cachePromise;
}; 