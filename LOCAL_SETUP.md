# 本地开发环境配置完成

## 服务状态

✅ **后端服务**: http://localhost:8080
- 已连接到本地MySQL (127.0.0.1:3306)
- 数据库: key_management
- 自动创建表结构

✅ **前端服务**: http://localhost:5173
- Vue 3 + Vant UI
- 移动端优先设计

## 登录信息

- **用户名**: admin
- **密码**: admin123

## 快速启动

下次启动服务，只需运行：

```bash
cd /c/baixuejie/code/go/awesomeProject/.worktrees/key-management-implementation
./start.sh
```

或者分别启动：

```bash
# 后端
cd backend
go run cmd/main.go

# 前端（新终端）
cd frontend
npm run dev
```

## 功能说明

### 1. 秘钥规格管理
- 创建不同类型的秘钥规格（如API密钥、License密钥）
- 每个规格可以有多个秘钥

### 2. 批量上传秘钥
- 在秘钥管理页面，点击"批量上传"
- 每行一个秘钥，支持批量导入
- 秘钥自动AES-256加密存储

### 3. 复制秘钥
- 点击"获取可用秘钥"优先获取未使用的秘钥
- 点击"复制"按钮，秘钥会按照模板格式复制到剪贴板
- 复制后自动标记为已使用

### 4. 自定义复制模板
- 在配置页面设置全局复制模板
- 使用 `{{key}}` 作为秘钥值占位符
- 例如: `export API_KEY={{key}}`

### 5. 使用追踪
- 查看秘钥使用状态（未使用/已使用）
- 按使用时间排序
- 过滤显示未使用秘钥

## 数据库配置

配置文件: `backend/config.yaml`

```yaml
database:
  host: 127.0.0.1
  port: 3306
  user: root
  password: "P@ssw0rd"
  dbname: key_management
```

## 技术栈

- **后端**: Go + Gin + GORM + JWT
- **前端**: Vue 3 + Vant UI + Pinia + Vue Router
- **数据库**: MySQL 8.0
- **安全**: AES-256加密 + JWT认证 + bcrypt密码哈希

## 项目结构

```
.
├── backend/
│   ├── cmd/main.go              # 入口文件
│   ├── internal/
│   │   ├── config/              # 配置管理
│   │   ├── database/            # 数据库连接
│   │   ├── models/              # 数据模型
│   │   ├── services/            # 业务逻辑
│   │   ├── handlers/            # API处理器
│   │   ├── middleware/          # 中间件
│   │   └── utils/               # 工具函数
│   └── config.yaml              # 配置文件
├── frontend/
│   └── src/
│       ├── views/               # 页面组件
│       ├── api/                 # API客户端
│       ├── stores/              # 状态管理
│       └── router/              # 路由配置
└── start.sh                     # 启动脚本
```

## 下一步

1. 访问 http://localhost:5173
2. 使用 admin/admin123 登录
3. 创建秘钥规格
4. 批量上传秘钥
5. 开始使用！
