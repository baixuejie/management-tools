# 个人工具网站 - 秘钥管理工具设计文档

**创建日期：** 2026-03-05
**版本：** 1.0

## 项目概述

构建一个个人工具网站，首个功能为秘钥管理工具。支持秘钥的批量上传、按规格分类、自定义模板复制、使用状态追踪等功能。系统采用移动端优先设计，同时兼容桌面端访问，支持长会话保持。

## 一、整体架构设计

### 1.1 容器组成

系统由4个Docker容器组成：

1. **MySQL容器**：存储用户、秘钥规格、秘钥数据
2. **Go后端容器**：提供RESTful API服务（端口8080）
3. **Vue前端容器**：Nginx托管静态文件（端口80）
4. **Nginx网关容器**：反向代理，统一入口（端口80）

### 1.2 网络架构

```
浏览器 → Nginx网关（80端口）
         ├─ 静态资源请求 → Vue前端容器
         └─ API请求（/api/*）→ Go后端容器
                              ↓
                          MySQL容器（3306端口）
```

### 1.3 技术栈

- **后端：** Go 1.25 + Gin + GORM + JWT
- **前端：** Vue 3 + Vant UI + Axios + Pinia
- **数据库：** MySQL 8.0
- **部署：** Docker Compose
- **安全：** JWT认证 + AES-256加密 + bcrypt密码哈希

## 二、数据库设计

### 2.1 秘钥规格表（key_specs）

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT | 主键，自增 |
| name | VARCHAR(100) | 规格名称（如"API密钥"、"License密钥"） |
| description | TEXT | 规格描述（可选） |
| created_at | DATETIME | 创建时间 |
| updated_at | DATETIME | 更新时间 |

### 2.2 秘钥表（keys）

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT | 主键，自增 |
| spec_id | BIGINT | 外键，关联key_specs.id |
| key_value | TEXT | 秘钥值（AES-256加密存储） |
| is_used | TINYINT | 使用状态（0=未使用，1=已使用） |
| used_at | DATETIME | 使用时间（NULL表示未使用） |
| created_at | DATETIME | 创建时间 |

**索引设计：**
- `idx_spec_used_time`：(spec_id, is_used, used_at) - 优化查询未使用秘钥和按时间排序
- `idx_spec_created`：(spec_id, created_at) - 优化按规格和时间查询

### 2.3 配置表（configs）

| 字段 | 类型 | 说明 |
|------|------|------|
| key | VARCHAR(100) | 配置键（如"copy_template"） |
| value | TEXT | 配置值（模板内容，如"API_KEY={{key}}"） |
| description | VARCHAR(255) | 配置说明 |
| updated_at | DATETIME | 更新时间 |

主键：key

## 三、Go后端API设计

### 3.1 技术组件

- **Gin框架**：轻量级HTTP框架
- **GORM**：ORM库，简化数据库操作
- **JWT**：用于登录认证的token
- **Viper**：配置文件管理

### 3.2 API路由

#### 认证相关
- `POST /api/auth/login` - 登录（验证用户名密码，返回JWT token）

#### 秘钥规格管理
- `GET /api/specs` - 获取所有规格列表
- `POST /api/specs` - 创建新规格
- `PUT /api/specs/:id` - 更新规格
- `DELETE /api/specs/:id` - 删除规格

#### 秘钥管理
- `GET /api/keys?spec_id=&is_used=&sort=used_at` - 查询秘钥（支持按规格、使用状态筛选，按时间排序）
- `POST /api/keys/batch` - 批量上传秘钥（接收多行文本）
- `POST /api/keys/:id/copy` - 复制秘钥（应用模板，标记为已使用，返回格式化后的文本）
- `DELETE /api/keys/:id` - 删除秘钥

#### 配置管理
- `GET /api/config/template` - 获取复制模板
- `PUT /api/config/template` - 更新复制模板

### 3.3 中间件

- **JWT认证中间件**：除登录接口外，所有API都需要验证token
- **CORS中间件**：允许前端跨域访问

### 3.4 项目结构

```
backend/
├── cmd/
│   └── main.go              # 入口文件
├── internal/
│   ├── config/              # 配置管理
│   ├── models/              # 数据模型（GORM）
│   ├── handlers/            # API处理器
│   ├── middleware/          # 中间件（JWT、CORS）
│   ├── services/            # 业务逻辑
│   └── utils/               # 工具函数（加密、JWT）
├── config.yaml              # 配置文件
├── Dockerfile
└── go.mod
```

## 四、Vue前端设计

### 4.1 技术栈

- **Vue 3 + Composition API**
- **Vant UI**：移动端优先组件库
- **Axios**：HTTP请求
- **Vue Router**：路由管理
- **Pinia**：状态管理

### 4.2 长会话保持方案

- JWT token有效期设置为7天（后端配置）
- 使用localStorage持久化存储token
- 应用启动时自动检查token有效性，有效则自动登录
- "记住我"选项：勾选后token有效期延长至30天

### 4.3 页面结构（移动端优先）

#### 1. 登录页（/login）
- 用户名、密码输入框
- "记住我"复选框（延长会话）
- 登录按钮（全宽）
- 自动登录逻辑：进入页面先检查token，有效则直接跳转主页

#### 2. 主页布局（/）
- 底部Tabbar导航：工具、设置
- 顶部NavBar：显示当前页面标题、退出按钮
- 内容区：显示当前选中的页面

#### 3. 秘钥管理页（/tools/keys）
- 顶部：规格选择下拉（Picker）
- 操作按钮区：批量上传、模板配置（固定在顶部）
- 秘钥列表（List组件）：
  - 每项显示：秘钥值（脱敏显示）、状态标签、使用时间
  - 右滑操作：复制、删除
  - 下拉刷新、上拉加载更多
  - 优先显示未使用秘钥

#### 4. 规格管理页（/settings/specs）
- 规格列表（List组件）
- 右滑操作：编辑、删除
- 底部添加按钮（FloatingButton）

### 4.4 桌面端适配

使用CSS媒体查询，屏幕宽度>768px时：
- 内容区最大宽度限制为600px，居中显示
- 增大字体和按钮尺寸
- Tabbar改为侧边栏

### 4.5 组件设计

- **KeyUploadDialog**：批量上传对话框
- **TemplateConfigDialog**：模板配置对话框
- **KeyTable**：秘钥列表表格组件

### 4.6 项目结构

```
frontend/
├── src/
│   ├── api/                 # API请求封装
│   ├── views/               # 页面组件
│   ├── components/          # 通用组件
│   ├── router/              # 路由配置
│   ├── stores/              # Pinia状态管理
│   ├── utils/               # 工具函数
│   └── App.vue
├── Dockerfile
└── package.json
```

## 五、Docker容器化配置

### 5.1 Docker Compose结构

```yaml
services:
  mysql:
    - MySQL 8.0镜像
    - 数据卷持久化（./data/mysql）
    - 初始化脚本自动创建数据库和表
    - 健康检查：mysqladmin ping

  backend:
    - 基于golang:1.25镜像构建
    - 多阶段构建：编译阶段 + 运行阶段（减小镜像体积）
    - 依赖mysql服务启动后再启动
    - 环境变量：数据库连接、JWT密钥、配置文件路径
    - 暴露8080端口

  frontend:
    - 基于node:18镜像构建Vue项目
    - 多阶段构建：npm build + nginx运行
    - 使用nginx:alpine镜像（轻量级）
    - 暴露80端口

  nginx-gateway:
    - nginx:alpine镜像
    - 反向代理配置：/ → frontend，/api → backend
    - 暴露宿主机80端口
    - 依赖frontend和backend启动
```

### 5.2 配置文件管理

- **config.yaml**：存储用户名密码哈希、JWT密钥、数据库连接、AES加密密钥等
- 通过volume挂载到backend容器
- **.env文件**：Docker Compose环境变量

### 5.3 一键启动

```bash
docker-compose up -d
```

## 六、安全性设计

### 6.1 认证安全

- 配置文件存储：用户名和密码的bcrypt哈希值存储在config.yaml
- JWT token机制：
  - 登录成功后签发JWT token（包含用户ID、过期时间）
  - 使用强随机密钥签名（配置文件配置）
  - "记住我"：30天有效期，否则7天
  - 每次API请求验证token有效性

### 6.2 秘钥数据加密

- 秘钥值使用AES-256加密存储到数据库
- 加密密钥存储在配置文件（与JWT密钥分离）
- 读取时解密，复制时应用模板后返回明文
- 防止数据库泄露导致秘钥直接暴露

### 6.3 传输安全

- 生产环境建议配置HTTPS（Nginx网关配置SSL证书）
- 敏感数据（密码、秘钥）通过HTTPS加密传输

### 6.4 其他安全措施

- 登录失败次数限制（防暴力破解）
- SQL注入防护（GORM参数化查询）
- XSS防护（前端输入验证和转义）
- CORS配置（限制允许的来源）

## 七、关键功能实现

### 7.1 批量上传秘钥

**前端：**
- 用户在文本框中粘贴多行文本，每行一个秘钥
- 选择秘钥规格
- 点击上传按钮

**后端：**
- 接收文本内容，按换行符分割
- 过滤空行和重复秘钥
- 批量加密并插入数据库
- 返回成功上传的数量

### 7.2 自定义模板复制

**模板格式：**
- 使用 `{{key}}` 作为秘钥值的占位符
- 示例：`API_KEY={{key}}`、`export TOKEN="{{key}}"`

**复制流程：**
1. 用户点击秘钥的复制按钮
2. 后端读取全局复制模板
3. 将秘钥值（解密后）替换模板中的 `{{key}}`
4. 标记秘钥为已使用，记录使用时间
5. 返回格式化后的文本
6. 前端复制到剪贴板，显示成功提示

### 7.3 秘钥列表展示

**排序规则：**
- 首次进入选择规格后，优先显示未使用的秘钥
- 未使用秘钥按创建时间倒序
- 已使用秘钥按使用时间倒序

**筛选功能：**
- 按规格筛选
- 按使用状态筛选（全部/未使用/已使用）

**移动端优化：**
- 下拉刷新
- 上拉加载更多（分页）
- 秘钥值脱敏显示（如：`abc***xyz`）

## 八、开发流程

### 8.1 初始化

1. 运行 `docker-compose up` 初始化数据库
2. 数据库自动创建表结构
3. 初始化默认配置（用户名密码、默认模板）

### 8.2 开发模式

- **后端开发**：本地运行 `go run cmd/main.go`，连接Docker中的MySQL
- **前端开发**：`npm run dev`，代理API请求到后端
- **联调测试**：`docker-compose up` 完整环境测试

### 8.3 生产部署

1. 构建镜像：`docker-compose build`
2. 启动服务：`docker-compose up -d`
3. 配置HTTPS（可选）
4. 配置域名解析（可选）

## 九、后续扩展能力

### 9.1 新工具添加

- 后端：添加新的API路由组
- 前端：添加新的页面和Tabbar菜单项
- 数据库：每个工具可以有独立的数据表

### 9.2 扩展示例

- 密码生成器
- 文本加密/解密工具
- 二维码生成器
- Base64编码/解码
- JSON格式化工具

### 9.3 架构优势

- 微服务架构天然支持功能模块化扩展
- 前后端分离，独立开发和部署
- Docker容器化，环境一致性好

## 十、预估工作量

- **后端开发**：2-3天
- **前端开发**：2-3天
- **Docker配置和联调**：1天
- **总计**：5-7天

## 十一、技术选型理由总结

| 技术 | 选择理由 |
|------|---------|
| Go + Gin | 轻量高效，适合API服务，部署简单 |
| Vue 3 | 现代化框架，生态丰富，学习曲线平缓 |
| Vant UI | 移动端优先，触摸交互体验好，轻量级 |
| MySQL | 成熟稳定，适合结构化数据存储 |
| Docker | 环境一致性，一键部署，易于扩展 |
| JWT | 无状态认证，适合前后端分离架构 |

---

**文档状态：** 已确认
**下一步：** 准备实施计划
