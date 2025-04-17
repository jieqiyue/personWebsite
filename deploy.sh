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
echo "=== 当前目录 ==="
pwd
echo "=== 当前目录输出结束 ==="

echo "=== 当前目录下面的文件 ==="
ls -la
echo "=== 当前目录下面的文件输出结束 ==="

# 输出当前最近5次git提交记录
echo "=== 最近5次git提交记录 ==="
git log -2
echo "=== 最近5次git提交记录输出结束 ==="

sleep 5

# 5. 编译 + 打包 + 生成最终镜像
echo "=== 前端编译 + 打包 + 生成最终镜像 ==="
docker build -t my-vue-app .

# 停掉旧容器（如果有）
docker stop my-vue-app-container || true
docker rm my-vue-app-container || true

# 重新运行容器
#docker run -d --name my-vue-app-container -p 80:80 my-vue-app

docker run -d --name my-vue-app-container -p 80:80 -p 443:443 \
  -v /etc/letsencrypt/live/jieqyyy.top/fullchain.pem:/etc/letsencrypt/live/jieqyyy.top/fullchain.pem \
  -v /etc/letsencrypt/live/jieqyyy.top/privkey.pem:/etc/letsencrypt/live/jieqyyy.top/privkey.pem \
  -v /opt/nginxlog:/var/log/nginx \
  --add-host=host.docker.internal:host-gateway \
  my-vue-app
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

cd "$PROJECT_DIR/chat-server"

# 8. 输出当前目录下面的文件
echo "=== 当前目录 ==="
pwd
echo "=== 当前目录输出结束 ==="

echo "=== 当前目录下面的文件 ==="
ls -la
echo "=== 当前目录下面的文件输出结束 ==="   

# 9. 进入后端项目目录
cd "$PROJECT_DIR/chat-server/cmd/server"

# 10. 输出当前目录下面的文件
echo "=== 当前目录 ==="
pwd
echo "=== 当前目录输出结束 ==="

rm -rf ./app

echo "=== 当前目录下面的文件 ==="
ls -la
echo "=== 当前目录下面的文件输出结束 ==="

echo "🛑 尝试释放 8080 端口..."
lsof -i :8080 -t | xargs kill || true

sleep 10

# 输出当前8080端口占用情况
echo "=== 当前8080端口占用情况 ==="
lsof -i :8080
echo "=== 当前8080端口占用情况输出结束 ==="

sleep 10 
# 11. 开始进行后端编译
echo "=== 开始进行后端编译 ==="
go build -o app main.go
echo "=== 后端编译结束 ==="

# 12. 输出当前目录下面的文件

echo "=== 当前目录下面的文件 ==="
ls -la
echo "=== 当前目录下面的文件输出结束 ==="   

# 休息5秒
sleep 5

# 13. 开始启动后端服务
echo "=== 开始启动后端服务 ==="
rm -rf out.log
chmod +x ./app
nohup ./app > out.log 2>&1 &
echo "=== 后端服务启动结束 ==="

# 14. 判断是否启动成功
echo "=== 判断是否启动成功 ==="
if [ $? -ne 0 ]; then
        echo "❌ 后端服务启动失败，请检查日志"
    exit 1
fi

sleep 10 
# 15. 输出后端服务启动日志
echo "=== 后端服务启动日志 ==="
cat out.log
echo "=== 后端服务启动日志输出结束 ==="


# 写入当前版本信息
echo "📝 记录当前 Git 版本信息..."

# 拿到 commit id 和 message
commit_id=$(git rev-parse --short HEAD)
commit_msg=$(git log -1 --pretty=%B)

# 写入到 version 文件
echo "Commit ID: $commit_id" > version.txt
echo "Message: $commit_msg" >> version.txt

echo "✅ 写入 version.txt 完成"

echo "✅ 后端服务启动成功"
echo "✅ 部署完成！"

