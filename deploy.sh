#!/bin/bash

# 个人网站一键部署脚本

# 设置项目信息
GITHUB_REPO="https://github.com/yourusername/personWebsite.git"
PROJECT_DIR="/tmp/personWebsite"
NGINX_DIR="/usr/share/nginx/html"
BACKUP_DIR="/tmp/nginx_backup_$(date +%Y%m%d%H%M%S)"

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

# 2. 克隆项目代码
echo "=== 从GitHub拉取代码 ==="
git clone "$GITHUB_REPO" "$PROJECT_DIR"
if [ $? -ne 0 ]; then
    echo "❌ 从GitHub拉取代码失败，请检查网络或仓库地址"
    exit 1
fi

# 3. 进入项目目录并安装依赖
echo "=== 安装项目依赖 ==="
cd "$PROJECT_DIR"
npm install
if [ $? -ne 0 ]; then
    echo "❌ 安装依赖失败"
    exit 1
fi

# 4. 构建项目
echo "=== 构建项目 ==="
npm run build
if [ $? -ne 0 ]; then
    echo "❌ 项目构建失败"
    exit 1
fi

# 5. 备份Nginx现有文件
echo "=== 备份Nginx现有文件 ==="
if [ -d "$NGINX_DIR" ] && [ "$(ls -A $NGINX_DIR)" ]; then
    echo "备份现有Nginx文件到: $BACKUP_DIR"
    mkdir -p "$BACKUP_DIR"
    cp -r "$NGINX_DIR"/* "$BACKUP_DIR"
fi

# 6. 部署到Nginx目录
echo "=== 部署到Nginx目录 ==="
# 确保目标目录存在
mkdir -p "$NGINX_DIR"
# 清空目标目录，但保留目录本身
rm -rf "$NGINX_DIR"/*
# 复制构建结果到Nginx目录
cp -r dist/* "$NGINX_DIR"
if [ $? -ne 0 ]; then
    echo "❌ 复制文件到Nginx目录失败"
    echo "正在从备份恢复..."
    rm -rf "$NGINX_DIR"/*
    cp -r "$BACKUP_DIR"/* "$NGINX_DIR"
    exit 1
fi

# 7. 设置适当的权限
echo "=== 设置文件权限 ==="
chown -R nginx:nginx "$NGINX_DIR"
chmod -R 755 "$NGINX_DIR"

# 8. 重启Nginx
echo "=== 重启Nginx服务 ==="
systemctl restart nginx
if [ $? -ne 0 ]; then
    echo "❌ 重启Nginx失败，请手动检查Nginx状态"
    exit 1
fi

echo "✅ 部署完成！"
echo "网站已成功部署到 $NGINX_DIR"
echo "备份文件位置: $BACKUP_DIR" 