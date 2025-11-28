# Model 模块文档

## 概述

Model 模块是 Navidrome 的数据访问层，负责定义数据模型、数据库操作、查询构建和业务实体管理。采用仓储模式 (Repository Pattern) 提供统一的数据访问接口。

## 目录结构

```
model/
├── datastore.go              # 数据存储主接口
├── errors.go                 # 自定义错误类型
├── searchable.go             # 可搜索实体接口
├── request/                  # 请求数据传输对象
│   └── request.go
├── metadata/                 # 元数据处理
│   ├── metadata.go           # 元数据管理器
│   ├── map_mediafile.go       # 音频文件元数据映射
│   ├── map_participants.go    # 参与者信息映射
│   ├── legacy_ids.go         # 传统 ID 处理
│   └── persistent_ids.go     # 持久化 ID
├── criteria/                 # 查询条件
│   ├── criteria.go           # 条件构建器
│   ├── fields.go             # 字段定义
│   ├── operators.go          # 操作符定义
│   └── json.go              # JSON 序列化
├── core_entities.go          # 核心实体定义
├── library.go               # 音乐库
├── folder.go                # 文件夹
├── artist.go                # 艺术家
├── artist_info.go           # 艺术家详细信息
├── album.go                 # 专辑
├── mediafile.go             # 媒体文件
├── genre.go                 # 音乐类型
├── tag.go                   # 标签
├── playlist.go              # 播放列表
├── playqueue.go             # 播放队列
├── transcoding.go           # 转码配置
├── player.go                # 播放器
├── radio.go                 # 电台
├── share.go                 # 分享
├── user.go                  # 用户
├── user_props.go            # 用户属性
├── scrobble_buffer.go       # Scrobble 缓冲
├── annotation.go            # 标注（收藏、评分等）
├── properties.go            # 系统属性
├── lyrics.go                # 歌词
├── bookmark.go              # 书签
├── tag_mappings.go          # 标签映射
├── file_types.go            # 文件类型定义
├── artwork_id.go            # 封面 ID 生成
└── participants.go          # 参与者信息
```

## 核心设计原则

### 仓储模式

所有数据访问通过 Repository 接口进行，提供统一的 CRUD 操作：

```go
type Repository interface {
    Count(options ...QueryOptions) (int64, error)
    Exists(id string) (bool, error)
    Put(entity interface{}) error
    Get(id string, options ...QueryOptions) (interface{}, error)
    GetAll(options ...QueryOptions) ([]interface{}, error)
    Delete(id string) error
    // 特定于实体的方法...
}
```

### 查询选项

统一的查询接口，支持排序、分页、过滤：

```go
type QueryOptions struct {
    Sort    string
    Order   string
    Max     int
    Offset  int
    Filters squirrel.Sqlizer
    Seed    string // 用于随机排序
}
```

### 数据库抽象

通过 `DataStore` 接口提供对所有仓储的访问：

```go
type DataStore interface {
    // 各实体仓储访问器
    Album(ctx context.Context) AlbumRepository
    Artist(ctx context.Context) ArtistRepository
    MediaFile(ctx context.Context) MediaFileRepository
    // ... 其他实体

    // 事务支持
    WithTx(block func(tx DataStore) error, scope ...string) error
    WithTxImmediate(block func(tx DataStore) error, scope ...string) error

    // 资源仓储（通用）
    Resource(ctx context.Context, model interface{}) ResourceRepository

    // 垃圾回收
    GC(ctx context.Context) error
}
```

## 核心实体

### 音乐实体层次

**Library (音乐库)**
```go
type Library struct {
    ID    string     `json:"id"`
    Name  string     `json:"name"`
    Path  string     `json:"path"`
    // ...
}
```

**Folder (文件夹)**
```go
type Folder struct {
    ID       string `json:"id"`
    Path     string `json:"path"`
    Name     string `json:"name"`
    LibraryID string `json:"libraryId"`
    // ...
}
```

**Artist (艺术家)**
```go
type Artist struct {
    ID              string    `json:"id"`
    Name            string    `json:"name"`
    AlbumCount      int       `json:"albumCount"`
    SongCount       int       `json:"songCount"`
    Size            int64     `json:"size"`
    Duration        float64   `json:"duration"`
    PlayCount       int64     `json:"playCount"`
    PlayDuration    float64   `json:"playDuration"`
    FullText        string    `json:"fullText"`
    OrderArtistName string    `json:"orderArtistName"`
    ArtistArt       string    `json:"artistArt"`
    ExternalInfoUrl string    `json:"externalInfoUrl"`
    ExternalInfoUpdatedAt time.Time `json:"-"`
    // ...
}
```

**Album (专辑)**
```go
type Album struct {
    ID                string     `json:"id"`
    Name              string     `json:"name"`
    Artist            string     `json:"artist"`
    ArtistID          string     `json:"artistId"`
    AlbumArtist       string     `json:"albumArtist"`
    AlbumArtistID     string     `json:"albumArtistId"`
    MaxYear           int        `json:"maxYear"`
    MinYear           int        `json:"minYear"`
    SongCount         int        `json:"songCount"`
    Duration          float64    `json:"duration"`
    PlayCount         int64      `json:"playCount"`
    PlayDuration      float64    `json:"playDuration"`
    Size              int64      `json:"size"`
    CoverArtId        string     `json:"coverArtId"`
    CoverArtPath      string     `json:"-"`
    CoverArtPriority  string     `json:"coverArtPriority"`
    ImagePaths        string     `json:"imagePaths"`
    Genre             string     `json:"genre"`
    FullText          string     `json:"fullText"`
    OrderAlbumName    string     `json:"orderAlbumName"`
    OrderArtistName   string     `json:"orderArtistName"`
    Compilation       bool       `json:"compilation"`
    Comment           string     `json:"comment"`
    ReleaseDate       string     `json:"releaseDate"`
    Discs             int        `json:"discs"`
    // ...
}
```

**MediaFile (媒体文件)**
```go
type MediaFile struct {
    ID               string     `json:"id"`
    Path             string     `json:"path"`
    Title            string     `json:"title"`
    Album            string     `json:"album"`
    Artist           string     `json:"artist"`
    AlbumArtist      string     `json:"albumArtist"`
    AlbumID          string     `json:"albumId"`
    ArtistID         string     `json:"artistId"`
    AlbumArtistID    string     `json:"albumArtistId"`
    Genre            string     `json:"genre"`
    Year             int        `json:"year"`
    TrackNumber      int        `json:"trackNumber"`
    DiscNumber       int        `json:"discNumber"`
    DiscSubtitle     string     `json:"discSubtitle"`
    Duration         float64    `json:"duration"`
    BitRate          int        `json:"bitRate"`
    Channels         int        `json:"channels"`
    Size             int64      `json:"size"`
    Suffix           string     `json:"suffix"`
    HasCoverArt      bool       `json:"hasCoverArt"`
    CoverArtId       string     `json:"coverArtId"`
    CoverArtPath     string     `json:"-"`
    PathHash         string     `json:"-"`
    FullText         string     `json:"fullText"`
    SortTitle        string     `json:"sortTitle"`
    SortAlbumName    string     `json:"sortAlbumName"`
    SortArtistName   string     `json:"sortArtistName"`
    OrderAlbumName   string     `json:"orderAlbumName"`
    OrderArtistName  string     `json:"orderArtistName"`
    Compilation      bool       `json:"compilation"`
    Comment          string     `json:"comment"`
    Lyrics           string     `json:"lyrics"`
    Bpm              int        `json:"bpm"`
    ReleaseDate      string     `json:"releaseDate"`
    SampleRate       int        `json:"sampleRate"`
    // ReplayGain 信息
    RgAlbumGain      float64    `json:"rgAlbumGain"`
    RgAlbumPeak      float64    `json:"rgAlbumPeak"`
    RgTrackGain      float64    `json:"rgTrackGain"`
    RgTrackPeak      float64    `json:"rgTrackPeak"`
    // 外部 ID
    MbzReleaseTrackId string     `json:"mbzReleaseTrackId"`
    MbzReleaseId     string     `json:"mbzReleaseId"`
    MbzArtistId      string     `json:"mbzArtistId"`
    MbzAlbumId       string     `json:"mbzAlbumId"`
    MbzRecordingId   string     `json:"mbzRecordingId"`
    // ...
}
```

### 用户和权限

**User (用户)**
```go
type User struct {
    ID              string   `json:"id"`
    UserName        string   `json:"userName"`
    Name            string   `json:"name"`
    Email           string   `json:"email"`
    Password        string   `json:"-"`
    IsAdmin         bool     `json:"isAdmin"`
    LastLoginAt     *time.Time `json:"lastLoginAt"`
    LastAccessAt    *time.Time `json:"lastAccessAt"`
    CreatedAt       time.Time `json:"createdAt"`
    UpdatedAt       time.Time `json:"updatedAt"`
    // ...
}
```

### 播放和交互

**Playlist (播放列表)**
```go
type Playlist struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    Comment     string    `json:"comment"`
    Duration    float64   `json:"duration"`
    Size        int64     `json:"size"`
    SongCount   int       `json:"songCount"`
    OwnerID     string    `json:"ownerId"`
    Public      bool      `json:"public"`
    Synced      bool      `json:"synced"`
    Path        string    `json:"path"`
    CreatedAt   time.Time `json:"createdAt"`
    UpdatedAt   time.Time `json:"updatedAt"`
    // ...
}
```

**Annotation (标注)**
```go
type Annotation struct {
    ItemID      string     `json:"itemId"`
    ItemType    string     `json:"itemType"`
    PlayCount   int64      `json:"playCount"`
    PlayDuration float64   `json:"playDuration"`
    Starred     bool       `json:"starred"`
    StarredAt   *time.Time `json:"starredAt"`
    Rating      int         `json:"rating"`
    Favorite    bool       `json:"favorite"`
    // ...
}
```

## 查询系统

### Criteria 系统

强大的查询构建器，支持复杂的过滤条件：

```go
type Criteria struct {
    Expression Expression `json:"expression"`
    Order     string     `json:"order"`
    Limit     int        `json:"limit"`
    Offset    int        `json:"offset"`
}

type Expression interface {
    ToSql() (string, []interface{}, error)
}
```

**支持的运算符** (operators.go):
- `And`, `Or` - 逻辑运算
- `Equals`, `NotEquals` - 相等性
- `Contains`, `StartsWith`, `EndsWith` - 字符串匹配
- `GreaterThan`, `LessThan`, `Between` - 数值比较
- `In`, `NotIn` - 集合操作
- `IsNull`, `IsNotNull` - 空值检查
- `NewMatches` - 全文搜索

**字段定义** (fields.go):
每个实体都有对应的字段定义，支持类型检查和验证。

### 搜索支持

**Searchable 接口**:
```go
type Searchable interface {
    Search(q string, offset, max int) (interface{}, error)
}
```

**全文搜索**:
- 使用 SQLite FTS (Full-Text Search)
- 支持中文分词
- 相关性排序
- 多字段搜索

## 元数据处理

### 音频元数据

**MetadataManager** (metadata/metadata.go):
- 支持多种音频格式 (MP3, FLAC, OGG, M4A, WMA 等)
- 使用 `dhowden/tag` 库读取标签
- 自动提取封面图像
- 轨道信息解析

**元数据映射** (map_mediafile.go):
- 标准化字段映射
- 多语言支持
- 字符编码检测
- 数据清洗

### 参与者信息

**艺术家解析** (map_participants.go):
- 自动识别表演者、作曲者、编曲者等
- 支持多艺术家格式
- 角色分类管理

### ID 管理

**MusicBrainz 集成** (legacy_ids.go):
- 外部 ID 映射
- 重复数据合并
- 关联关系维护

## 数据库操作

### 事务支持

```go
// 简单事务
err := ds.WithTx(func(tx DataStore) error {
    // 执行多个操作
    return nil
})

// 立即事务
err := ds.WithTxImmediate(func(tx DataStore) error {
    // 执行操作并立即提交
    return nil
})
```

### 批量操作

**高效插入**:
```go
// 使用预编译语句批量插入
func (r *mediaFileRepository) insertAll(mfs []MediaFile) error
```

**索引优化**:
- 自动创建复合索引
- 查询性能监控
- 索引使用分析

### 数据完整性

**外键约束**:
- 参照完整性检查
- 级联删除配置
- 数据一致性验证

## 性能优化

### 缓存策略

**应用层缓存**:
- 查询结果缓存
- ID 映射缓存
- 统计数据缓存

**数据库缓存**:
- 连接池管理
- 预编译语句缓存
- 查询计划缓存

### 查询优化

**分页实现**:
```go
// 使用 OFFSET/LIMIT
// 大数据集使用游标分页
// 随机排序支持
```

**索引使用**:
- 复合索引优化
- 覆盖索引设计
- 查询计划分析

## 错误处理

### 自定义错误类型

```go
var (
    ErrNotFound      = errors.New("resource not found")
    ErrInvalidID     = errors.New("invalid id format")
    ErrDuplicate     = errors.New("duplicate resource")
    ErrInvalidData   = errors.New("invalid data")
)
```

### 错误包装

使用 Go 1.13+ 错误包装机制：
```go
return fmt.Errorf("failed to save artist %s: %w", artist.Name, err)
```

## 测试策略

### 测试数据管理

**测试夹具**:
```go
func TestArtistRepository(t *testing.T) {
    ds := &mockDataStore{}
    repo := NewArtistRepository(ds)
    // 测试逻辑...
}
```

**数据隔离**:
- 每个测试使用独立数据库
- 自动清理测试数据
- 并发测试支持

### 集成测试

**数据库集成**:
- 使用内存 SQLite
- 迁移脚本测试
- 性能基准测试

## 扩展指南

### 添加新实体

1. **定义实体结构**:
```go
type NewEntity struct {
    ID        string    `json:"id"`
    Name      string    `json:"name"`
    CreatedAt time.Time `json:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt"`
}
```

2. **实现 Repository 接口**:
```go
type NewEntityRepository interface {
    ResourceRepository
    // 特定方法
    FindByName(name string) (*NewEntity, error)
}
```

3. **添加到 DataStore**:
```go
type DataStore interface {
    // 现有接口...
    NewEntity(ctx context.Context) NewEntityRepository
}
```

4. **创建数据库表**:
- 编写迁移脚本
- 定义索引
- 设置约束

5. **添加测试用例**:
- 单元测试
- 集成测试
- 性能测试

### 查询功能扩展

1. **添加新字段到 Criteria**:
```go
const FieldNewField = "newField"
```

2. **实现字段验证器**:
```go
func validateNewField(value interface{}) error {
    // 验证逻辑
}
```

3. **更新查询构建器**:
- 支持新的过滤条件
- 添加排序选项
- 优化查询计划