#!/bin/bash

# 秘钥管理工具启动脚本

echo "正在启动秘钥管理工具..."

# 启动后端
echo "启动后端服务 (端口 8080)..."
cd backend
go run cmd/main.go &
BACKEND_PID=$!
cd ..

# 等待后端启动
sleep 3

# 启动前端
echo "启动前端服务 (端口 5173)..."
cd frontend
npm run dev &
FRONTEND_PID=$!
cd ..

echo ""
echo "✅ 服务启动成功！"
echo ""
echo "📱 前端地址: http://localhost:5173"
echo "🔧 后端API: http://localhost:8080/api"
echo ""
echo "🔑 默认账号: admin"
echo "🔒 默认密码: admin123"
echo ""
echo "按 Ctrl+C 停止服务"

# 等待用户中断
wait
