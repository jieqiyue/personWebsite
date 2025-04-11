import MarkdownIt from 'markdown-it'
import { full as emoji } from 'markdown-it-emoji'
import container from 'markdown-it-container'
import hljs from 'highlight.js'

// 创建基本的MarkdownIt实例
const md = new MarkdownIt({
  html: true,
  linkify: true,
  typographer: true,
  breaks: true,
  highlight: function (str, lang) {
    if (lang && hljs.getLanguage(lang)) {
      try {
        return '<pre class="hljs"><button class="code-copy-btn" title="复制代码"><svg viewBox="0 0 24 24" width="16" height="16" stroke="currentColor" stroke-width="2" fill="none"><path d="M8 17.929H6c-1.105 0-2-.912-2-2.036V5.036C4 3.912 4.895 3 6 3h8c1.105 0 2 .912 2 2.036v1.866m-6 .17h8c1.105 0 2 .91 2 2.035v10.857C20 21.09 19.105 22 18 22h-8c-1.105 0-2-.911-2-2.036V9.107c0-1.124.895-2.036 2-2.036z"/></svg></button><code>' +
               hljs.highlight(str, { language: lang, ignoreIllegals: true }).value +
               '</code></pre>';
      } catch (__) {}
    }
    
    return '<pre class="hljs"><button class="code-copy-btn" title="复制代码"><svg viewBox="0 0 24 24" width="16" height="16" stroke="currentColor" stroke-width="2" fill="none"><path d="M8 17.929H6c-1.105 0-2-.912-2-2.036V5.036C4 3.912 4.895 3 6 3h8c1.105 0 2 .912 2 2.036v1.866m-6 .17h8c1.105 0 2 .91 2 2.035v10.857C20 21.09 19.105 22 18 22h-8c-1.105 0-2-.911-2-2.036V9.107c0-1.124.895-2.036 2-2.036z"/></svg></button><code>' + md.utils.escapeHtml(str) + '</code></pre>';
  }
})

// 自定义图片渲染规则
const defaultRender = md.renderer.rules.image || function(tokens, idx, options, env, self) {
  return self.renderToken(tokens, idx, options);
};

md.renderer.rules.image = function (tokens, idx, options, env, self) {
  const token = tokens[idx];
  const srcIndex = token.attrIndex('src');
  const altIndex = token.attrIndex('alt');
  let src = token.attrs[srcIndex][1];
  const alt = altIndex >= 0 ? token.attrs[altIndex][1] : '';
  
  // 处理相对路径的图片链接，CSDN或其他图片服务器链接保持不变
  if (src && !src.startsWith('http') && !src.startsWith('/')) {
    // 如果是相对路径且不以/开头，转为绝对路径
    src = '/' + src;
    token.attrs[srcIndex][1] = src;
  }
  
  // 为图片添加加载错误处理
  const imgHtml = defaultRender(tokens, idx, options, env, self);
  return imgHtml.replace('<img', '<img loading="lazy" onload="this.classList.add(\'loaded\')" onerror="this.classList.add(\'error\'); this.setAttribute(\'data-error-src\', this.src); this.alt=\'图片加载失败: \' + this.alt"');
};

// 使用emoji插件，配置自定义表情符号快捷方式
md.use(emoji, {
  shortcuts: {
    // 添加一些常用表情符号的快捷方式
    'smile': [':)', ':-)', '😊'],
    'laughing': [':D', ':-D', '😄'],
    'wink': [';)', ';-)', '😉'],
    'frown': [':(', ':-(', '😞'],
    'stuck_out_tongue': [':P', ':-P', '😛'],
    'confused': [':/', ':-/', '😕'],
    'open_mouth': [':O', ':-O', '😮'],
    'heart': ['<3', '❤️']
  }
})

// 创建提示容器
md.use(container, 'tip', {
  validate: function(params) {
    return params.trim().match(/^tip\s+(.*)$/);
  },
  render: function (tokens, idx) {
    const m = tokens[idx].info.trim().match(/^tip\s+(.*)$/);
    if (tokens[idx].nesting === 1) {
      // 开始标签
      return `<div class="tip custom-block"><p class="custom-block-title">${md.utils.escapeHtml(m[1])}</p>\n`;
    } else {
      // 结束标签
      return '</div>\n';
    }
  }
});

// 创建警告容器
md.use(container, 'warning', {
  validate: function(params) {
    return params.trim().match(/^warning\s+(.*)$/);
  },
  render: function (tokens, idx) {
    const m = tokens[idx].info.trim().match(/^warning\s+(.*)$/);
    if (tokens[idx].nesting === 1) {
      // 开始标签
      return `<div class="warning custom-block"><p class="custom-block-title">${md.utils.escapeHtml(m[1])}</p>\n`;
    } else {
      // 结束标签
      return '</div>\n';
    }
  }
});

// 创建危险容器
md.use(container, 'danger', {
  validate: function(params) {
    return params.trim().match(/^danger\s+(.*)$/);
  },
  render: function (tokens, idx) {
    const m = tokens[idx].info.trim().match(/^danger\s+(.*)$/);
    if (tokens[idx].nesting === 1) {
      // 开始标签
      return `<div class="danger custom-block"><p class="custom-block-title">${md.utils.escapeHtml(m[1])}</p>\n`;
    } else {
      // 结束标签
      return '</div>\n';
    }
  }
});

// 创建成功容器
md.use(container, 'success', {
  validate: function(params) {
    return params.trim().match(/^success\s+(.*)$/);
  },
  render: function (tokens, idx) {
    const m = tokens[idx].info.trim().match(/^success\s+(.*)$/);
    if (tokens[idx].nesting === 1) {
      // 开始标签
      return `<div class="success custom-block"><p class="custom-block-title">${md.utils.escapeHtml(m[1])}</p>\n`;
    } else {
      // 结束标签
      return '</div>\n';
    }
  }
});

// 创建详情容器
md.use(container, 'details', {
  validate: function(params) {
    return params.trim().match(/^details\s+(.*)$/);
  },
  render: function (tokens, idx) {
    const m = tokens[idx].info.trim().match(/^details\s+(.*)$/);
    if (tokens[idx].nesting === 1) {
      // 开始标签
      return `<details class="custom-block"><summary>${md.utils.escapeHtml(m[1])}</summary>\n`;
    } else {
      // 结束标签
      return '</details>\n';
    }
  }
});

// 导出配置好的markdown-it实例
export default md; 