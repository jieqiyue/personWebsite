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
        return '<pre class="hljs"><code>' +
               hljs.highlight(str, { language: lang, ignoreIllegals: true }).value +
               '</code></pre>';
      } catch (__) {}
    }
    
    return '<pre class="hljs"><code>' + md.utils.escapeHtml(str) + '</code></pre>';
  }
})

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