# 第一阶段：构建阶段
FROM node:18-alpine AS builder

WORKDIR /app
COPY package*.json ./
RUN npm install

COPY . .
RUN npm run build

# 第二阶段：部署阶段
FROM nginx:stable-alpine

# 拷贝构建好的前端静态文件到 nginx 的默认目录
COPY --from=builder /app/dist /usr/share/nginx/html

# 拷贝自定义 nginx 配置（可选）
COPY nginx.conf /etc/nginx/nginx.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]