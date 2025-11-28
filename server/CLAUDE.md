# Server 模块文档

## 概述

Server 模块是 Navidrome 的核心 HTTP 服务器层，负责处理所有传入的请求，包括 Subsonic API、原生 API、静态资源服务等。

## 目录结构

```
server/
├── server.go              # 主服务器实现
├── server_suite_test.go   # 测试套件
├── serve_index.go         # 索引页面服务
├── auth.go               # 认证相关
├── middlewares.go        # 中间件
├── initial_setup.go      # 初始化设置
├── subsonic/             # Subsonic API 实现
│   ├── api.go            # 主 API 路由
│   ├── browsing.go       # 浏览相关端点
│   ├── media_retrieval.go # 媒体检索
│   ├── stream.go         # 流媒体
│   ├── playlists.go      # 播放列表
│   ├── searching.go      # 搜索
│   ├── users.go          # 用户管理
│   ├── media_annotation.go # 媒体标注
│   ├── system.go         # 系统信息
│   ├── jukebox.go        # 点唱机功能
│   ├── radio.go          # 电台功能
│   ├── sharing.go        # 分享功能
│   ├── opensubsonic.go   # OpenSubsonic 扩展
│   ├── bookmarks.go      # 书签
│   ├── helpers.go        # 辅助函数
│   ├── middlewares.go    # Subsonic 特定中间件
│   ├── album_lists.go    # 专辑列表
│   ├── library_scanning.go # 库扫描
│   └── responses/       # API 响应格式
│       ├── responses.go  # 响应构建器
│       └── errors.go    # 错误处理
├── nativeapi/            # 原生 API
│   ├── native_api.go     # 主路由
│   ├── translations.go   # 翻译
│   ├── playlists.go      # 播放列表管理
│   ├── inspect.go        # 检查端点
│   └── missing.go        # 缺失处理
├── public/               # 公共端点
│   ├── public.go         # 路由定义
│   ├── handle_streams.go # 流媒体处理
│   ├── handle_downloads.go # 下载处理
│   ├── handle_images.go  # 图像处理
│   ├── handle_shares.go  # 分享处理
│   └── encode_id.go      # ID 编码
├── events/               # 事件处理
│   ├── events.go         # 事件管理
│   ├── sse.go           # Server-Sent Events
│   └── events_suite_test.go
└── backgrounds/          # 背景图像
    └── handler.go       # 背景处理
```

## 核心组件

### Server (server/server.go)

主要的 HTTP 服务器结构体，负责：
- 初始化路由
- 挂载各种 API 端点
- 配置中间件
- 管理认证

```go
type Server struct {
    router   chi.Router
    ds       model.DataStore
    appRoot  string
    broker   events.Broker
    insights metrics.Insights
}
```

### 路由架构

服务器使用 Chi v5 路由器，分层路由结构：

1. **根级别路由** (`/`):
   - UI 索引页
   - 健康检查
   - 静态资源

2. **Subsonic API** (`/rest/*`):
   - 完整的 Subsonic API 兼容性
   - 支持 OpenSubsonic 扩展

3. **原生 API** (`/api/*`):
   - Navidrome 特定功能
   - 不限于 Subsonic 标准的扩展

4. **公共资源** (`/*`):
   - 流媒体 (`/stream`, `/download`)
   - 图像服务 (`/img`, `/cover`)
   - 分享链接 (`/share`)

5. **SSE 事件** (`/events`):
   - 实时事件推送
   - 进度通知
   - 库扫描状态

### 认证与授权

**认证流程** (auth.go, initial_setup.go):
- JWT Token 认证
- 用户管理
- 初始管理员设置
- 会话管理

**中间件** (middlewares.go):
- 请求日志
- CORS 处理
- 速率限制
- 用户认证检查
- 错误处理

## Subsonic API 实现

### API 结构

`subsonic/api.go` 定义了主要的 Subsonic API 路由结构：

```go
func (api *Router) setupRoutes(r chi.Router) {
    // 子路由分组
    r.Route("/", func(r chi.Router) {
        r.Get("/", api.ping)
        // 各种 Subsonic 端点...
    })
}
```

### 主要端点类别

**浏览和检索** (browsing.go):
- `getMusicFolders` - 音乐文件夹
- `getArtist` / `getArtists` - 艺术家信息
- `getAlbum` / `getAlbumList` - 专辑信息
- `getSong` / `getSongs` - 歌曲信息

**媒体检索** (media_retrieval.go):
- `stream` - 音频流
- `download` - 下载
- `getCoverArt` - 封面图像

**播放功能** (stream.go, jukebox.go):
- 流媒体传输
- 转码支持
- 点唱机控制

**用户管理** (users.go):
- 用户认证
- 权限管理
- 设置同步

**搜索** (searching.go):
- `search` / `search2` / `search3`
- 高级过滤
- 分页支持

**标注功能** (media_annotation.go):
- 星级评分
- 收藏
- 播放计数
- 新增：歌曲评论 (`getSongComments`)

### 响应格式

**统一响应结构** (responses/responses.go):
```go
type SubsonicResponse struct {
    Status  string      `json:"status"`
    Version string      `json:"version"`
    Type    string      `json:"type,omitempty"`
    Error   *Error     `json:"error,omitempty"`
    // 动态数据字段...
}
```

**错误处理** (responses/errors.go):
- 标准 Subsonic 错误代码
- 统一错误格式
- 详细错误信息

## 原生 API

**特定扩展** (nativeapi/):
- 不受 Subsonic 标准限制
- 现代 RESTful 设计
- 前端特定功能

**主要端点**:
- `/api/translations` - 多语言支持
- `/api/playlists` - 播放列表 CRUD
- `/api/inspect` - 系统检查
- `/api/missing` - 缺失文件分析

## 公共资源服务

### 流媒体处理

**音频流** (handle_streams.go):
- 支持多种音频格式
- 实时转码
- 范围请求 (支持断点续传)
- 缓存控制

**转码配置**:
- FFmpeg 集成
- 多种输出格式
- 质量设置
- 客户端特定配置

### 图像服务

**封面和图像** (handle_images.go):
- 自动封面提取
- 多尺寸缩放
- 格式转换 (JPEG/PNG)
- 缓存优化
- 专辑封面优先级配置

### 下载服务

**批量下载** (handle_downloads.go):
- ZIP 打包
- 转码下载
- 权限检查
- 进度跟踪

### 分享功能

**公开分享** (handle_shares.go):
- 临时链接
- 访问控制
- 下载权限
- 过期管理

## 事件系统

### Server-Sent Events

**实时通信** (events/sse.go):
- 进度推送
- 库扫描状态
- 播放同步
- 系统通知

**事件类型** (events/events.go):
- 扫描开始/结束
- 文件添加/删除
- 播放状态变化
- 错误通知

## 开发指南

### 添加新 API 端点

**Subsonic API**:
1. 在相应的处理文件中添加函数
2. 在 `api.go` 中注册路由
3. 在 `responses/` 中定义响应格式
4. 添加测试用例

**原生 API**:
1. 在 `nativeapi/` 中添加处理器
2. 在 `native_api.go` 中注册路由
3. 使用现代 RESTful 设计
4. 添加相应的认证检查

### 错误处理

统一使用 `responses` 包的错误处理机制：
```go
return resp.Error(req, code, message, details)
```

### 中间件开发

在 `middlewares.go` 中添加新的中间件函数，并在 `New()` 中注册。

### 测试

每个功能都应该有对应的测试：
- 单元测试测试逻辑
- 集成测试测试端点
- 快照测试测试响应格式

## 性能优化

### 缓存策略
- 图像缓存 (`ImageCacheSize`)
- 转码缓存 (`TranscodingCacheSize`)
- 响应缓存
- 静态资源缓存

### 并发处理
- Goroutine 池管理
- 数据库连接池
- 流媒体并发限制
- 速率限制

## 配置选项

服务器相关的配置项（参考 `conf/configuration.go`）：
- `Address`, `Port` - 服务器地址
- `TLSCert`, `TLSKey` - TLS 配置
- `SessionTimeout` - 会话超时
- `LogLeve`l - 日志级别
- `EnableTranscodingConfig` - 转码启用
- `FFmpegPath` - FFmpeg 路径
- 各种缓存大小设置

## 调试和监控

### 日志
- 结构化日志
- 可配置级别
- 请求追踪
- 性能指标

### 监控
- Prometheus 指标
- 健康检查端点
- 性能分析
- 错误统计