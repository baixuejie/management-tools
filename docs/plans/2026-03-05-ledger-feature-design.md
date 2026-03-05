# 账本功能设计文档

**创建日期：** 2026-03-05
**版本：** 1.0

## 项目概述

在现有秘钥管理工具基础上，新增账本功能模块。支持成本记录、交易记账（区分新客/续费）、购买人管理、渠道手续费计算、统计分析等功能。同时优化UI样式和用户体验，确保最常用功能（记账、复制秘钥）的操作路径最短。

## 一、核心需求

### 1.1 功能需求

**成本管理**
- 记录每笔投入成本的金额和备注
- 查看成本明细列表
- 支持删除成本记录

**交易记账**
- 区分新客和续费两种类型
- 新客：输入购买人姓名，自动创建购买人记录
- 续费：从已有购买人中选择，支持名称搜索
- 记录交易金额、渠道（闲鱼/微信）、记录者
- 自动计算渠道手续费（闲鱼7‰，微信0）

**购买人管理**
- 存储购买人姓名
- 支持修改购买人名称
- 支持按名称搜索

**统计分析**
- 总成本、总收入、净利润
- 新客数量、续费数量
- 总手续费

### 1.2 用户体验需求

- 保持H5页面操作便捷性
- 记账和复制秘钥功能操作路径最短（1-2步完成）
- 优化页面样式，提升视觉体验
- 支持多用户独立登录

## 二、数据库设计

### 2.1 用户表（users）

| 字段 | 类型 | 约束 | 说明 |
|------|------|----|
| id | BIGINT | PRIMARY KEY, AUTO_INCREMENT | 用户ID |
| username | VARCHAR(50) | UNIQUE, NOT NULL | 登录用户名 |
| password_hash | VARCHAR(255) | NOT NULL | bcrypt密码哈希 |
| display_name | VARCHAR(100) | NOT NULL | 显示名称 |
| created_at | DATETIME | NOT NULL | 创建时间 |
| updated_at | DATETIME | NOT NULL | 更新时间 |

**初始数据：**
```sql
INSERT INTO users (username, password_hash, display_name) VALUES
('admin', '$2a$10$...', '白了个白'),
('fanchen', '$2a$10$...', '凡尘');
```

### 2.2 成本记录表（cost_records）

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGINT | PRIMARY KEY, AUTO_INCREMENT | 记录ID |
| amount | DECIMAL(10,2) | NOT NULL | 投入金额 |
| note | TEXT | | 备注说明 |
| recorded_by | BIGINT | NOT NULL, FOREIGN KEY | 记录者ID |
| created_at | DATETIME | NOT NULL | 记录时间 |

**索引：**
- `idx_created_at` (created_at DESC)

### 2.3 购买人表（customers）

| 字段 | 类型 | 约束 | 说明 |
|----|------|---|------|
| id | BIGINT | PRIMARY KEY, AUTO_INCREMENT | 购买人ID |
| name | VARCHAR(100) | NOT NULL | 购买人名称 |
| created_by | BIGINT | NOT NULL, FOREIGN KEY | 创建者ID |
| created_at | DATETIME | NOT NULL | 创建时间 |
| updated_at | DATETIME | NOT NULL | 更新时间 |

**索引：**
- `idx_name` (name)

### 2.4 交易记录表（transactions）

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | BIGINT | PRIMARY KEY, AUTO_INCREMENT | 交易ID |
| customer_id | BIGINT | NOT NULL, FOREIGN KEY | 购买人ID |
| amount | DECIMAL(10,2) | NOT NULL | 交易金额 |
| channel | VARCHAR(20) | NOT NULL | 渠道（xianyu/wechat） |
| commission_rate | DECIMAL(5,4) | NOT NULL | 手续费率 |
| commission_amount | DECIMAL(10,2) | NOT NULL | 手续费金额 |
| is_new_customer | TINYINT | NOT NULL | 是否新客（1=新客，0=续费） |
| recorded_by | BIGINT | NOT NULL, FOREIGN KEY | 记录者ID |
| created_at | DATETIME | NOT NULL | 交易时间 |

**索引：**
- `idx_customer_time` (customer_id, created_at DESC)
- `idx_recorder_time` (recorded_by, created_at DESC)
- `idx_created_at` (created_at DESC)

**手续费规则：**
- 闲鱼（xianyu）：commission_rate = 0.0070（7‰）
- 微信（wechat）：commission_rate = 0.0000（0）
- commission_amount = amount * commission_rate

## 三、后端API设计

### 3.1 认证系统重构

#### POST /api/auth/login
登录接口（重构为数据库验证）

**请求：**
```json
{
  "username": "admin",
  "password": "admin123"
}
```

**响应：**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": 1,
    "username": "admin",
    "display_name": "白了个白"
  }
}
```

#### GET /api/auth/me
获取当前登录用户信息

**响应：**
```json
{
  "id": 1,
  "username": "admin",
  "display_name": "白了个白"
}
```

### 3.2 成本管理API

#### GET /api/ledger/costs
查询成本记录列表

**查询参数：**
- `page`: 页码（默认1）
- `limit`: 每页数量（默认20）

**响应：**
```json
{
  "total": 100,
  "items": [
    {
      "id": 1,
      "amount": 500.00,
      "note": "服务器续费",
      "recorded_by": 1,
      "recorder_name": "白了个白",
      "created_at": "2026-03-05T10:30:00Z"
    }
  ]
}
```

#### POST /api/ledger/costs
添加成本记录

**请求：**
```json
{
  "amount": 500.00,
  "note": "服务器续费"
}
```

**说明：** recorded_by自动取当前登录用户ID

#### DELETE /api/ledger/costs/:id
删除成本记录

### 3.3 购买人管理API

#### GET /api/ledger/customers
查询购买人列表

**查询参数：**
- `search`: 搜索关键词（可选）

**响应：**
```json
{
  "items": [
    {
    "id": 1,
      "name": "张三",
      "created_by": 1,
      "created_at": "2026-03-05T10:00:00Z"
    }
  ]
}
```

#### POST /api/ledger/customers
创建购买人（新客记账时自动调用）

**请求：**
```json
{
  "name": "张三"
}
```

#### PUT /api/ledger/customers/:id
修改购买人名称

**请求：**
```json
{
  "name": "张三（已改名）"
}
```

### 3.4 交易记录API

#### GET /api/ledger/transactions
查询交易记录列表
**查询参数：**
- `customer_id`: 购买人ID（可选）
- `is_new_customer`: 是否新客（可选，0/1）
- `page`: 页码（默认1）
- `limit`: 每页数量（默认20）

**响应：**
```json
{
  "total": 200,
  "items": [
    {
      "id": 1,
      "customer_id": 1,
      "customer_name": "张三",
      "amount": 100.00,
      "channel": "xianyu",
      "commission_rate": 0.0070,
      "commission_amount": 0.70,
      "is_new_customer": true,
      "recorded_by": 1,
      "recorder_name": "白了个白",
      "created_at": "2026-03-05T11:00:00Z"
    }
  ]
}
```

#### POST /api/ledger/transactions
创建交易记录

**请求（新客）：**
```json
{
  "customer_name": "张三",
  "amount": 100.00,
  "channel": "xianyu",
  "is_new_customer": true
}
```

**请求（续费）：**
```json
{
  "customer_id": 1,
  "amount": 100.00,
  "channel": "wechat",
  "is_new_customer": false
}
```

**处理逻辑：**
1. 如果is_new_customer=true，先创建customer记录
2. 根据channel计算commission_rate和commission_amount
3. recorded_by自动取当前用户ID
4. 插入transaction记录

### 3.5 统计分析API

#### GET /api/ledger/statistics
获取账本统计数据

**响应：**
```json
{
  "total_cost": 5000.00,
  "total_revenue": 8000.00,
  "total_commission": 56.00,
  "net_profit": 2944.00,
  "new_customers": 50,
  "renewal_customers": 30
}
```

**计算规则：**
- total_cost: SUM(cost_records.amount)
- total_revenue: SUM(transactions.amount)
- total_commission: SUM(transactions.commission_amount)
- net_profit: total_revenue - total_cost - total_commission
- new_customers: COUNT(DISTINCT transactions WHERE is_new_customer=1)
- renewal_customers: COUNT(DISTINCT transactions WHERE is_new_customer=0)

## 四、前端页面设计

### 4.1 路由结构

```
/login              - 登录页
/ledger             - 记账首页（默认页）
/keys               - 秘钥管理页
/more               - 更多功能页
  /more/costs       - 成本记录页
  /more/customers   - 购买人管理页
  /more/specs       - 规格管理页
  /more/template    - 模板配置页
  /more/profile     - 个人信息页
```

### 4.2 底部Tabbar导航

```
┌─────────┬─────┬─────────┐
│  记账   │  秘钥   │  更多   │
└─────────┴─────────┴─────────┘
```

- 记账：/ledger（默认选中）
- 秘钥：/keys
- 更多：/more

### 4.3 登录页（/login）

**组件：**
- 用户名下拉选择（van-dropdown-menu）
  - 选项：admin（白了个白）、fanchen（凡尘）
- 密码输入框（van-field，type=password）
- 登录按钮（van-button，type=primary，block）

**逻辑：**
- 页面加载时检查localStorage中的token
- 如果token有效，自动跳转到/ledger
- 登录成功后存储token和用户信息

### 4.4 记账首页（/ledger）

**布局结构：**

```
┌────────────────────┐
│  统计卡片（可折叠）          │
│  总成本 | 总收入 | 净利润  │
│  新客数 | 续费数 | 手续费    │
└─────────────────┘
┌───────────────┐
│  快速记账区域                │
│  [新客] [续费]           │
│  购买人：___________         │
│  金额：___________           │
│  渠道：[闲鱼] [微信]         │
│  [提交]         │
└─────────────────────────┘
┌───────────────┐
│  交易记录列表                │
│  ┌────────────┐    │
│  │ 张三 [新客]         │    │
│  │ ¥100.00  闲鱼 -0.70 │    │
│  │ 白了个白 2026-03-05 │    │
│  └──────────┘    │
│  ...                   │
└──────────────┘
[+成本] ← 右下角浮动按钮
```

**统计卡片（van-cell-group）**
- 渐变背景（#4facfe → #00f2fe）
- 白色文字
- 数字大字体（20px），带滚动动画
- 点击展开/收起

**快速记账区域（van-form）**
- 客户类型Tab（van-tabs）：新客/续费
- 新客模式：
  - 购买人姓名输入框（van-field）
  - 金额输入框（van-field，type=digit）
  - 渠道选择（自定义按钮组）
  - 提交按钮（van-button，渐变背景）
- 续费模式：
  - 购买人选择（van-search + van-list）
  - 金额输入框
  - 渠道选择
  - 提交按钮

**交易记录列表（van-list）**
- 下拉刷新（van-pull-refresh）
- 上拉加载更多
- 每条记录：
  - 购买人名称 + 新客/续费标签（van-tag）
  - 金额（绿色，大字）+ 渠道图标 + 手续费（灰色小字）
  - 记录者 + 时间（灰色小字）

**浮动按钮（van-floating-bubble）**
- 图标：加号
- 点击弹出成本添加对话框

### 4.5 秘钥管理页（/keys）

**优化后布局：**

```
┌───────────────┐
│  规格选择：[下拉]  [批量上传]│
└──────────────────────┘
┌──────────────┐
│  快速复制区域                │
│  [获取可用秘钥]              │
└────────────────┘
┌─────────────────────┐
│  [查看全部秘钥 ▼]            │
│  （点击展开秘钥列表）         │
└───────────┘
```

**快速复制区域**
- 大按钮（van-button，渐变背景，图标+文字）
- 点击后：
  1. 调用API获取未使用秘钥
  2. 弹出对话框显示秘钥值（大字体）
  3. 自动复制到剪贴板
  4. 显示Toast提示"已复制！"
  5. 标记秘钥为已使用

**秘钥列表**
- 默认折叠
- 点击"查看全部秘钥"展开
- 保持原有列表功能

### 4.6 更多功能页（/more）

**功能卡片列表（van-cell-group）**

**账本管理**
- 成本记录（/more/costs）
- 购买人管理（/more/customers）

**秘钥管理**
- 规格管理（/more/specs）
- 复制模板配置（/more/template）

**系统设置**
- 个人信息（/more/profile）
- 退出登录

## 五、UI样式设计

### 5.1 色彩方案

**主色调**
- 渐变蓝绿：#4facfe → #00f2fe
- 用于：主按钮、统计卡片背景、Tab选中状态

**功能色**
- 成功色：#07c160（微信绿）- 金额、成功提示
- 警告色：#ff976a（闲鱼橙）- 闲鱼渠道、手续费
- 信息色：#1989fa（蓝色）- 新客标签
- 危险色：#ee0a24（红色）- 删除操作

**中性色**
- 背景色：#f7f8fa（浅灰）
- 卡片背景：#ffffff（白色）
- 文字主色：#323233（深灰）
- 文字辅助色：#969799（灰色）
- 边框色：#ebedf0（浅灰）

### 5.2 字体层级

- 大标题：18px，粗体
- 标题：16px，粗体
- 金额数字：20px，粗体
- 正文：14px，常规
- 辅助信息：12px，常规

### 5.3 组件样式

**卡片样式**
- 背景：白色
- 圆角：12px
- 阴影：0 2px 8px rgba(0,0,0,0.08)
- 间距：12px

**按钮样式**
- 主按钮：渐变背景，白色文字，圆角8px，高度44px
- 次要按钮：白色背景，主色边框，主色文字
- 渠道按钮：选中有背景色（闲鱼橙/微信绿）

**输入框样式**
- 边框：1px solid #ebedf0
- 聚焦：1px solid #4facfe
- 圆角：8px
- 高度：44px
- 字体：14px

**标签样式**
- 新客：蓝色背景，白色文字，圆角4px
- 续费：橙色背景，白色文字，圆角4px

### 5.4 交互动画

**按钮点击**
```css
.button:active {
  transform: scale(0.95);
  transition: transform 0.1s;
}
```
**页面切换**
- 左右滑动过渡（300ms）

**列表加载**
- 骨架屏占位（van-skeleton）

**数据提交**
- 按钮loading状态（van-loading）

**成功提示**
- Toast提示 + 震动反馈（navigator.vibrate）

### 5.5 响应式设计

**移动端（<768px）**
- 全屏布局
- 底部Tabbar固定
- 输入框高度44px
- 按钮高度44px
- 字体14px

**桌面端（>768px）**
- 内容区最大宽度600px，居中
- Tabbar改为左侧边栏（宽度200px）
- 输入框高度48px
- 按钮高度48px
- 字体16px
- 鼠标悬停效果

## 六、技术实现要点

### 6.1 后端重构

**认证系统迁移**
1. 创建users表，插入初始用户数据
2. 修改登录逻辑：从数据库查询用户，验证密码
3. JWT payload包含user_id
4. 中间件从token解析user_id，注入到context

**项目结构调整**
```
backend/
├── internal/
│   ├── models/
│   │   ├── user.go          # 新增
│   │   ├── cost_record.go   # 新增
│   │   ├── customer.go      # 新增
│   │   ├── transaction.go   # 新增
│   ├── services/
│   │   ├── auth_service.go  # 重构
│   │   ├── ledger_service.go # 新增
│   ├── handlers/
│   │   ├── auth.go          # 重构
│   │   ├── ledger.go        # 新增
```

### 6.2 前端重构

**状态管理（Pinia）**
```javascript
// stores/auth.js - 重构
export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: localStorage.getItem('token'),
    user: JSON.parse(localStorage.getItem('user') || 'null')
  }),
  actions: {
    async login(username, password) { ... },
    logout() { ... }
  }
})

// stores/ledger.js - 新增
export const useLedgerStore = defineStore('ledger', {
  state: () => ({
    statistics: {},
    transactions: [],
    customers: []
  }),
  actions: {
    async fetchStatistics() { ... },
    async createTransaction(data) { ... }
  }
})
```

**路由守卫**
```javascript
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  if (to.path !== '/login' && !authStore.token) {
    next('/login')
  } else {
    next()
  }
})
```

### 6.3 数据库迁移

**迁移脚本**
```sql
-- 1. 创建users表
CREATE TABLE users (...);

-- 2. 插入初始用户
INSERT INTO users ...;

-- 3. 创建账本相关表
CREATE TABLE cost_records (...);
CREATE TABLE customers (...);
CREATE TABLE transactions (...);

-- 4. 创建索引
CREATE INDEX idx_created_at ON cost_records(created_at DESC);
...
```

## 七、开发计划

### 阶段1：后端开发（2天）

**Day 1**
- 创建数据库表和索引
- 实现users表和认证系统重构
- 实现成本管理API
- 实现购买人管理API

**Day 2**
- 实现交易记录API
- 实现统计分析API
- 单元测试和接口测试

### 阶段2：前端开发（2-3天）

**Day 1**
- 重构登录页和认证逻辑
- 实现记账首页布局
- 实现快速记账表单

**Day 2**
- 实现交易记录列表
- 实现统计卡片
- 优化秘钥管理页

**Day 3**
- 实现更多功能页
- 实现成本记录页和购买人管理页
- UI样式优化和动画效果

### 阶段3：联调测试（1天）

- 前后端联调
- 功能测试
- UI/UX优化
- 性能优化

## 八、测试用例

### 8.1 功能测试

**记账流程**
1. 新客记账：输入姓名、金额、选择渠道 → 提交成功 → 列表显示
2. 续费记账：搜索购买人、输入金额、选择渠道 → 提交成功 → 列表显示
3. 手续费计算：闲鱼渠道自动计算7‰手续费，微信渠道手续费为0
4. 统计数据：提交交易后统计数据实时更新

**成本管理**
1. 添加成本：输入金额和备注 → 提交成功 → 列表显示
2. 删除成本：删除成本记录 → 统计数据更新

**购买人管理**
1. 搜索购买人：输入关键词 → 实时过滤列表
2. 修改购买人：修改名称 → 保存成功 → 交易记录同步更新

### 8.2 性能测试

- 交易记录列表加载时间 < 500ms
- 统计数据计算时间 < 200ms
- 页面切换动画流畅（60fps）

### 8.3 兼容性测试

- iOS Safari
- Android Chrome
- 微信内置浏览器
- 桌面端Chrome/Edge

## 九、后续扩展

### 9.1 功能扩展

- 数据导出（Excel）
- 图表统计（收入趋势、渠道占比）
- 购买人详情页（交易历史、消费统计）
- 批量导入交易记录

### 9.2 性能优化

- 交易记录列表虚拟滚动
- 统计数据缓存
- 图片懒加载

---

**文档状态：** 已确认
**下一步：** 准备实施计划
