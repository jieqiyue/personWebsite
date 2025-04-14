#!/bin/bash

# 个人网站一键部署脚本

# 设置项目信息
GITHUB_REPO="git@github.com:jieqiyue/personWebsite.git"
PROJECT_DIR="/tmp/personWebsite"
NGINX_DIR="/usr/share/nginx/html"


echo "=== 开始部署个人网站 ==="

# 确保有sudo权限
if [ "$(id -u)" != "0" ]; then
   echo "此脚本需要sudo权限执行，请使用sudo运行" 
   exit 1
fi

# 1. 创建或清空临时项目目录
echo "=== 准备项目目录 ==="
if [ -d "$PROJECT_DIR" ]; then
    echo "清空现有目录: $PROJECT_DIR"
    rm -rf "$PROJECT_DIR"
fi
mkdir -p "$PROJECT_DIR"
echo "准备项目目录结束"

# 2. 克隆项目代码
echo "=== 从GitHub拉取代码 ==="
git clone "$GITHUB_REPO" "$PROJECT_DIR"
if [ $? -ne 0 ]; then
    echo "❌ 从GitHub拉取代码失败，请检查网络或仓库地址"
    exit 1
fi
echo "从GitHub拉取代码结束"

# 3. 进入项目目录
cd "$PROJECT_DIR"

# 4. 输出当前目录下面的文件
echo "=== 当前目录下面的文件 ==="
ls -la
echo "=== 当前目录下面的文件输出结束 ==="

# 5. 编译 + 打包 + 生成最终镜像
echo "=== 前端编译 + 打包 + 生成最终镜像 ==="
docker build -t my-vue-app .

# 停掉旧容器（如果有）
docker stop my-vue-app-container || true
docker rm my-vue-app-container || true

# 重新运行容器
docker run -d --name my-vue-app-container -p 80:80 my-vue-app
echo "前端编译 + 打包 + 生成最终镜像结束"

# 6. 删除悬空镜像
echo "=== 删除悬空镜像 ==="
docker image prune -f
echo "=== 删除悬空镜像结束 ==="

# 6. 输出当前docker状态
echo "=== 当前docker状态 ==="
docker ps 
echo "=== 当前docker状态输出结束 ==="

# 7. 开始进行后端部署
echo "=== 开始进行后端部署 ==="

echo "✅ 部署完成！"
echo "网站已成功部署到 $NGINX_DIR"
