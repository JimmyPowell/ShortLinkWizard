# 短链接生成器

基于Gin和GORM实现的短链接生成服务，采用MVC架构设计。

## 项目结构

```
.
├── config/         # 配置文件
├── controller/     # 控制器层
├── model/          # 数据模型层
├── service/        # 业务逻辑层
├── router/         # 路由配置
├── main.go         # 主入口文件
└── README.md       # 项目文档
```

## API 接口

### 1. 创建短链接

**请求**
- 方法: POST
- URL: /api/shorten
- 请求体:
```json
{
  "original_url": "长链接URL"
}
```

**响应**
```json
{
  "short_url": "短链接代码"
}
```

### 2. 访问短链接

**请求**
- 方法: GET
- URL: /api/{code}

**响应**
- 状态码: 301
- 重定向到原始URL

## 快速开始

1. 安装依赖
```bash
go mod tidy
```

2. 启动服务
```bash
go run main.go
```

3. 测试API
```bash
# 创建短链接
curl -X POST http://localhost:8080/api/shorten \
  -H "Content-Type: application/json" \
  -d '{"original_url":"https://example.com"}'

# 访问短链接
curl -v http://localhost:8080/api/{short_code}
```

## 依赖

- Gin: Web框架
- GORM: ORM库
- MySQL: 数据库
