# Navidrome 开发指南

## 开发环境设置

### 系统要求

**必需软件**:
- **Go** 1.24.1+
- **Node.js** (版本参考 `.nvmrc`)
- **Git**
- **FFmpeg** (用于音频转码)

**推荐工具**:
- **Docker** (用于容器化部署)
- **Make** (用于构建命令)
- **IDE** (支持 Go 和 JavaScript/TypeScript)

### 快速开始

```bash
# 克隆仓库
git clone https://github.com/navidrome/navidrome.git
cd navidrome

# 自动设置开发环境
make setup

# 启动开发服务器
make dev
```

## 开发工作流

### 日常开发

**启动开发环境**:
```bash
make dev          # 启动前端+后端开发服务器
make server       # 仅启动后端服务器
make watch        # 启动测试监视模式
```

**代码质量检查**:
```bash
make lint         # Go 代码 lint
make lintall      # Go + JS 代码 lint
make format       # 格式化代码
```

**测试**:
```bash
make test         # 运行 Go 测试
make testrace     # 运行带竞态检测的测试
make testall      # 运行所有测试
```

### 项目结构导航

**核心开发目录**:
- `server/` - 后端 API 和服务器逻辑
- `model/` - 数据模型和数据库层
- `ui/src/` - 前端 React 应用
- `conf/` - 配置管理
- `core/` - 核心业务逻辑
- `utils/` - 工具函数

**配置和构建**:
- `Makefile` - 构建和开发命令
- `go.mod` - Go 依赖管理
- `ui/package.json` - 前端依赖管理
- `db/migrations/` - 数据库迁移

## 后端开发

### Go 开发环境

**依赖管理**:
```bash
go mod tidy                    # 清理依赖
go mod download               # 下载依赖
make download-deps            # 下载所有依赖
```

**构建和运行**:
```bash
make build                     # 构建应用
make debug-build              # 调试版本构建
go run .                      # 直接运行
```

**数据库管理**:
```bash
make migration-sql name=migration_name    # 创建 SQL 迁移
make migration-go name=migration_name     # 创建 Go 迁移
```

### API 开发

**添加新的 Subsonic API 端点**:

1. **在 `server/subsonic/` 中添加处理函数**:
```go
// server/subsonic/example.go
func (api *Router) getExample(c *rest.Context) {
    params := c.Params()

    // 业务逻辑
    result := &ExampleResponse{
        ID:   params["id"],
        Data: "example data",
    }

    c.Respond(result)
}
```

2. **在 `server/subsonic/api.go` 中注册路由**:
```go
// 在 setupRoutes 函数中添加
r.Get("/getExample", api.getExample)
```

3. **在 `server/subsonic/responses/` 中定义响应格式**:
```go
// server/subsonic/responses/example.go
type ExampleResponse struct {
    ID   string `json:"id"`
    Data string `json:"data"`
}

func (r *Router) Example(req *rest.Request, example *ExampleResponse) map[string]interface{} {
    return r.NewSubsonicResponse(example)
}
```

4. **添加测试**:
```go
// server/subsonic/example_test.go
func TestGetExample(t *testing.T) {
    // 测试逻辑
}
```

**添加新的原生 API 端点**:

1. **在 `server/nativeapi/` 中创建处理器**:
```go
// server/nativeapi/example.go
func (api *Router) handleExample(w http.ResponseWriter, r *http.Request) {
    // 处理逻辑
    json.NewEncoder(w).Encode(response)
}
```

2. **在 `server/nativeapi/native_api.go` 中注册路由**:
```go
r.HandleFunc("/example", api.handleExample).Methods("GET")
```

### 数据库操作

**使用 Model 层**:
```go
import "github.com/navidrome/navidrome/model"

// 获取数据
albums, err := ds.Album(ctx).GetAll(model.QueryOptions{
    Max:    10,
    Sort:   "name",
    Order:  "ASC",
})

// 保存数据
err := ds.Album(ctx).Put(album)
```

**事务处理**:
```go
err := ds.WithTx(func(tx model.DataStore) error {
    // 多个操作
    if err := tx.Album(ctx).Put(album1); err != nil {
        return err
    }
    if err := tx.Album(ctx).Put(album2); err != nil {
        return err
    }
    return nil
})
```

### 配置管理

**访问配置**:
```go
import "github.com/navidrome/navidrome/conf"

// 使用配置选项
musicFolder := conf.Server.MusicFolder
port := conf.Server.Port
logLevel := conf.Server.LogLevel
```

**添加新配置**:
1. 在 `conf/configuration.go` 的 `configOptions` 结构体中添加字段
2. 更新配置文件解析逻辑
3. 添加默认值和验证

## 前端开发

### React 开发环境

**前端开发设置**:
```bash
cd ui
npm ci              # 安装依赖
npm run dev        # 开发模式
npm run build      # 构建生产版本
```

**组件开发**:
```bash
cd ui
npm run test       # 运行测试
npm run lint       # 代码检查
npm run prettier   # 格式化代码
```

### 添加新组件

**创建资源组件**:

1. **创建组件文件**:
```jsx
// ui/src/example/ExampleList.jsx
import React from 'react'
import { List, Datagrid, TextField } from 'react-admin'

export const ExampleList = props => (
  <List {...props}>
    <Datagrid>
      <TextField source="name" />
      <TextField source="description" />
    </Datagrid>
  </List>
)
```

2. **创建资源定义**:
```jsx
// ui/src/example/index.jsx
import ExampleList from './ExampleList'

export default {
  name: 'example',
  list: ExampleList,
  icon: ExampleIcon,
}
```

3. **在 App.jsx 中注册资源**:
```jsx
import example from './example'

// 在 <Admin> 组件中添加
<Resource name="example" {...example} />
```

**状态管理**:

1. **创建 Action**:
```javascript
// ui/src/actions/example.js
export const setExampleData = data => ({
  type: 'SET_EXAMPLE_DATA',
  payload: data
})
```

2. **创建 Reducer**:
```javascript
// ui/src/reducers/example.js
const exampleReducer = (state = {}, action) => {
  switch (action.type) {
    case 'SET_EXAMPLE_DATA':
      return { ...state, ...action.payload }
    default:
      return state
  }
}

export default exampleReducer
```

3. **在 Store 中注册**:
```javascript
// ui/src/store/createAdminStore.js
import exampleReducer from '../reducers/example'

const rootReducer = combineReducers({
  // 其他 reducers...
  example: exampleReducer
})
```

### API 集成

**使用 DataProvider**:
```javascript
// ui/src/example/useExampleData.js
import { useQuery } from 'react-query'
import dataProvider from '../dataProvider'

export const useExampleData = id => {
  return useQuery(
    ['example', id],
    () => dataProvider('GET_ONE', 'example', { id })
  )
}
```

**自定义 API 调用**:
```javascript
// ui/src/example/api.js
import { fetchUtils } from 'react-admin'
import { httpClient } from '../dataProvider'

export const getCustomData = async params => {
  const { json } = await httpClient('/api/custom-endpoint', {
    method: 'GET',
  })
  return json
}
```

## 测试

### 后端测试

**单元测试**:
```go
// 使用 Ginkgo/Gomega
import (
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
)

var _ = Describe("AlbumRepository", func() {
    It("should create album", func() {
        album := &model.Album{
            Name: "Test Album",
        }
        err := ds.Album(ctx).Put(album)
        Expect(err).To(BeNil())
    })
})
```

**集成测试**:
```go
// 测试 API 端点
func TestGetAlbum(t *testing.T) {
    // 设置测试数据库
    ds := setupTestDB(t)
    defer cleanupTestDB(ds)

    // 创建测试数据
    album := createTestAlbum(ds)

    // 测试 API
    req := httptest.NewRequest("GET", "/rest/getAlbum?id="+album.ID, nil)
    w := httptest.NewRecorder()

    // 执行请求
    router.ServeHTTP(w, req)

    // 验证响应
    assert.Equal(t, 200, w.Code)
}
```

### 前端测试

**组件测试**:
```jsx
// ui/src/example/ExampleList.test.jsx
import { render, screen } from '@testing-library/react'
import { ExampleList } from './ExampleList'

describe('ExampleList', () => {
  it('renders example list', () => {
    render(<ExampleList />)
    expect(screen.getByText('Example')).toBeInTheDocument()
  })
})
```

**Hook 测试**:
```jsx
// ui/src/example/useExampleData.test.js
import { renderHook } from '@testing-library/react-hooks'
import { useExampleData } from './useExampleData'

describe('useExampleData', () => {
  it('should fetch data', async () => {
    const { result, waitFor } = renderHook(() => useExampleData('test-id'))

    await waitFor(() => {
      expect(result.current.isSuccess).toBe(true)
    })
  })
})
```

### 运行测试

**Go 测试**:
```bash
make test              # 运行所有测试
make testrace          # 运行带竞态检测的测试
make watch            # 测试监视模式
```

**JavaScript 测试**:
```bash
cd ui
npm test              # 运行测试
npm run test:ci       # CI 模式测试
npm run test:coverage # 覆盖率测试
```

## 调试

### 后端调试

**日志配置**:
```bash
# 设置日志级别
export ND_LOGLEVEL=debug

# 启用详细日志
export ND_ENABLEINSIGHTSCOLLECTOR=false
```

**调试构建**:
```bash
make debug-build      # 构建调试版本
```

**性能分析**:
```bash
# 启用 pprof
go tool pprof http://localhost:4533/debug/pprof/profile
```

### 前端调试

**React DevTools**:
- 安装浏览器扩展
- 检查组件状态
- 性能分析

**Redux DevTools**:
- 状态变化追踪
- 时间旅行调试
- Action 日志

**浏览器开发者工具**:
- 网络请求分析
- 控制台错误检查
- 性能监控

## 部署

### 本地部署

**构建生产版本**:
```bash
make build          # 构建应用
./navidrome         # 运行
```

**Docker 部署**:
```bash
make docker-image   # 构建 Docker 镜像
docker run deluan/navidrome:develop
```

### 环境配置

**配置文件** (navidrome.toml):
```toml
MusicFolder = "/path/to/music"
DataFolder = "/path/to/data"
LogLevel = "info"
Address = "0.0.0.0"
Port = 4533
```

**环境变量**:
```bash
export ND_MUSICFOLDER="/path/to/music"
export ND_DATAFOLDER="/path/to/data"
export ND_LOGLEVEL="debug"
```

## 代码规范

### Go 代码规范

**格式化**:
```bash
go fmt ./...          # 格式化 Go 代码
goimports -w ./...    # 更新导入
```

**Lint 检查**:
```bash
make lint            # 运行 golangci-lint
```

**命名规范**:
- 包名：小写，简短
- 变量名：驼峰式
- 常量：大写或驼峰式
- 接口名：以 -er 结尾

### JavaScript 代码规范

**格式化**:
```bash
cd ui
npm run prettier     # 格式化代码
npm run lint         # ESLint 检查
```

**代码风格**:
- 使用 Prettier 格式化
- 遵循 ESLint 规则
- 组件使用 PascalCase
- 变量使用 camelCase

## 性能优化

### 后端优化

**数据库优化**:
- 合理使用索引
- 批量操作
- 连接池配置
- 查询优化

**内存管理**:
- 避免内存泄漏
- 合理使用缓存
- 垃圾回收优化

### 前端优化

**代码分割**:
```jsx
// 懒加载组件
const LazyComponent = React.lazy(() => import('./LazyComponent'))
```

**状态优化**:
- 避免不必要的重渲染
- 使用 useMemo 和 useCallback
- 合理的状态结构

## 贡献指南

### 提交代码

**Git 工作流**:
```bash
git checkout -b feature/new-feature
# 开发...
make testall make lintall  # 运行测试和检查
git commit -m "feat: add new feature"
git push origin feature/new-feature
# 创建 Pull Request
```

**提交信息规范**:
- `feat:` 新功能
- `fix:` 错误修复
- `docs:` 文档更新
- `style:` 代码格式
- `refactor:` 重构
- `test:` 测试
- `chore:` 构建/工具

### 代码审查

**检查清单**:
- [ ] 测试通过
- [ ] 代码格式化
- [ ] 文档更新
- [ ] 性能考虑
- [ ] 安全检查

## 故障排除

### 常见问题

**编译错误**:
```bash
# 清理模块缓存
go clean -modcache
go mod download

# 重新构建
make clean
make build
```

**前端构建错误**:
```bash
cd ui
rm -rf node_modules package-lock.json
npm ci
npm run build
```

**数据库问题**:
```bash
# 重新初始化数据库
rm -f data/navidrome.db
./navidrome
```

### 调试技巧

**后端调试**:
- 使用 `log.Debugf()` 添加调试日志
- 检查配置文件
- 验证数据库连接

**前端调试**:
- 使用 `console.log()` 调试
- 检查网络请求
- 验证状态变化

## 资源和文档

**官方文档**:
- [Navidrome 官网](https://www.navidrome.org/)
- [开发文档](https://www.navidrome.org/docs/developers/)
- [API 文档](https://www.navidrome.org/docs/developers/subsonic-api/)

**相关资源**:
- [React Admin 文档](https://marmelab.com/react-admin/)
- [Material-UI 文档](https://material-ui.com/)
- [Go 语言官方文档](https://golang.org/doc/)

**社区支持**:
- [GitHub Issues](https://github.com/navidrome/navidrome/issues)
- [Discord 聊天](https://discord.gg/xh7j7yF)
- [Subreddit](https://www.reddit.com/r/navidrome/)