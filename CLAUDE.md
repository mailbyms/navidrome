# CLAUDE.md

这个文件为 Claude Code (claude.ai/code) 在此代码库中工作提供指导。

## 项目概述

Navidrome 是一个开源的基于网络的音乐收藏服务器和流媒体服务器。它允许用户从任何浏览器或移动设备收听自己的音乐收藏，就像个人版的 Spotify。

### 技术栈

**后端 (Go):**
- Go 1.24.1+
- Chi v5 路由框架
- SQLite 数据库 (通过 go-sqlite3)
- Wire 依赖注入
- Ginkgo/Gomega 测试框架

**前端 (JavaScript/React):**
- React 17 + Material-UI v4
- React Admin 框架
- Vite 构建工具
- Redux 状态管理
- Vitest 测试框架

**其他工具:**
- FFmpeg 用于音频转码
- Docker 容器化部署
- Prometheus 监控指标

## 项目架构

### 核心目录结构

```
navidrome/
├── server/          # 后端服务器代码
│   ├── subsonic/    # Subsonic API 兼容层
│   ├── nativeapi/   # 原生 API
│   ├── public/      # 公共端点 (下载、流媒体等)
│   └── events/      # SSE 事件处理
├── model/           # 数据模型和仓储层
├── conf/            # 配置管理
├── core/            # 核心业务逻辑
├── resources/       # 资源文件 (logo等)
├── ui/              # 前端 React 应用
│   └── src/         # 前端源码
├── db/              # 数据库迁移
└── utils/           # 工具函数
```

### 后端架构

**服务器层 (server/):**
- `server.go`: 主 HTTP 服务器
- `subsonic/`: Subsonic API 兼容实现
- `nativeapi/`: Navidrome 特定 API
- `public/`: 公共资源访问 (流媒体、图像、下载)
- `events/`: Server-Sent Events 和实时事件

**数据层 (model/):**
- `datastore.go`: 数据存储接口定义
- 实体模型: `album.go`, `artist.go`, `mediafile.go`, `playlist.go` 等
- 仓储模式: 每个实体类型都有对应的 Repository 接口
- `criteria/`: 查询条件和过滤

**核心业务层 (core/):**
- 认证和授权
- 音频转码
- 元数据处理
- 扫描和导入

### 前端架构

基于 React Admin 框架，主要模块：
- `album/`: 专辑相关组件
- `artist/`: 艺术家相关组件
- `audioplayer/`: 音频播放器
- `common/`: 通用组件
- `dataProvider/`: API 数据提供者
- `layout/`: 布局和导航

## 常用开发命令

### 环境准备
```bash
make setup          # 安装所有依赖 (Go + Node.js)
make setup-dev      # 同上，别名
```

### 开发模式
```bash
make dev            # 启动完整开发环境 (前端+后端热重载)
make server         # 仅启动后端开发服务器
make watch          # 启动测试监视模式
```

### 测试
```bash
make test           # 运行 Go 测试
make testrace       # 运行带竞态检测的测试
make testall        # 运行 Go + JS 测试
make watch          # 测试监视模式
```

### 代码质量
```bash
make lint           # Go 代码 lint
make lintall        # Go + JS 代码 lint
make format         # 格式化代码 (Go + JS)
```

### 构建
```bash
make build          # 构建应用 (包含前端构建)
make buildjs        # 仅构建前端
make debug-build    # 调试版本构建
```

### 数据库迁移
```bash
make migration-sql name=migration_name    # 创建 SQL 迁移文件
make migration-go name=migration_name     # 创建 Go 迁移文件
```

### Docker 和部署
```bash
make docker-build   # 跨平台 Docker 构建
make docker-image   # 构建 Docker 镜像
make package        # 创建所有平台的二进制包
```

## 开发工作流

### 后端开发
1. 使用 `make server` 启动后端热重载
2. 测试使用 `make test` 或 `make testrace`
3. 数据库变更需要创建迁移文件
4. API 端点主要在 `server/subsonic/` (Subsonic 兼容) 和 `server/nativeapi/`

### 前端开发
1. 使用 `make dev` 启动完整开发环境
2. 前端代码位于 `ui/src/`
3. 组件测试使用 Vitest
4. UI 基于 Material-UI 和 React Admin

### 添加新功能
1. 后端: 在相应的 service 层添加业务逻辑，在 server 层暴露 API
2. 前端: 在对应的资源模块下添加组件，更新 dataProvider
3. 数据库变更: 创建迁移文件，更新 model 层

### 调试
- 后端调试: 使用 `make debug-build` 生成带调试信息的二进制
- 日志级别通过配置文件或环境变量控制
- 前端调试使用浏览器开发者工具

## 配置

配置文件格式支持 TOML/JSON/YAML，主要配置项：
- `MusicFolder`: 音乐文件夹路径
- `DataFolder`: 数据文件夹
- `LogLevel`: 日志级别
- `Address`/`Port`: 服务器监听地址
- `EnableTranscodingConfig`: 转码配置

完整的配置选项参考 `conf/configuration.go` 中的 `configOptions` 结构体。

## API 接口

Navidrome 提供:
1. **Subsonic API** (`/rest/*`) - 兼容现有音乐播放器客户端
2. **Native API** (`/api/*`) - Navidrome 特定功能
3. **公共端点** - 流媒体、下载、图像访问

API 文档: https://www.navidrome.org/docs/developers/subsonic-api/

## 测试策略

- **单元测试**: 使用 Ginkgo/Gomega (Go) 和 Vitest (JS)
- **集成测试**: API 端点测试
- **快照测试**: Subsonic API 响应验证 (`make snapshots`)
- **数据库测试**: 使用内存 SQLite 进行测试

## 常见任务

### 添加新的 Subsonic API 端点
1. 在 `server/subsonic/` 中添加处理函数
2. 在 `server/subsonic/responses/` 中定义响应格式
3. 添加对应的测试
4. 更新 API 文档

### 添加新的前端页面
1. 在 `ui/src/` 对应模块下创建组件
2. 在 `App.jsx` 中添加路由
3. 更新 dataProvider 以支持新的数据访问
4. 添加相应的测试

### 数据库结构变更
1. 使用 `make migration-sql` 或 `make migration-go` 创建迁移
2. 更新 `model/` 中对应的实体定义
3. 调整相关的 Repository 实现
4. 测试迁移的向前和向后兼容性