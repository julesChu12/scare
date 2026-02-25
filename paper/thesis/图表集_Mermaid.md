# 论文图表集（Mermaid 语法）

> 本文件包含论文所需的所有图表，使用 Mermaid 语法编写。
> 
> **使用方法**：
> 1. 直接复制到支持 Mermaid 的 Markdown 编辑器（如 Typora、VS Code + Mermaid 插件）
> 2. 或访问 https://mermaid.live 在线编辑器导出为 PNG/SVG 图片
> 3. 或使用 Draw.io 的 Mermaid 导入功能

---

## 图3-1 现有业务流程图

```mermaid
flowchart TD
    A[老年人/家属] -->|电话/微信/线下| B[社区工作人员]
    B --> C[人工整理需求]
    C --> D[人工判断分配]
    D --> E[服务站点]
    E --> F[工作人员上门]
    F --> G[服务完成]
    G -->|电话/微信| B
    B -->|口头反馈| A
    
    style A fill:#e1f5fe
    style B fill:#fff3e0
    style E fill:#e8f5e9
```

---

## 图3-2 目标业务流程图

```mermaid
flowchart TD
    A[老年人/家属] -->|扫码进入| B[C端系统]
    B -->|填写需求+定位| C[需求提交]
    C -->|地理坐标| D[地理围栏匹配]
    D -->|自动分配| E[服务站点任务池]
    E -->|认领| F[工作人员]
    F -->|上门服务| G[完成任务]
    G -->|上传照片| H[B端系统]
    H -->|推送通知| B
    B -->|查看进度/评价| I[服务评价]
    
    style A fill:#e1f5fe
    style B fill:#e1f5fe
    style D fill:#fff9c4
    style E fill:#e8f5e9
    style F fill:#e8f5e9
    style H fill:#e8f5e9
```

---

## 图3-3 系统功能模块图

```mermaid
mindmap
  root((社区养老信息分发平台))
    C端用户端
      快速开通
      需求提交
      需求管理
      服务评价
      个人中心
      消息通知
    B端管理门户
      登录认证
      任务管理
        任务池查看
        任务认领
        任务完成
        任务转派
      需求管理
      站点管理
      围栏管理
      用户管理
      权限管理
      统计报表
      通知管理
```

---

## 图3-4 地理围栏匹配流程图

```mermaid
flowchart TD
    A[用户提交坐标] --> B[获取所有活跃围栏]
    B --> C[按优先级降序排列]
    C --> D{遍历围栏}
    D --> E{BoundingBox检查}
    E -->|不在矩形内| D
    E -->|在矩形内| F{射线法精确匹配}
    F -->|不匹配| D
    F -->|匹配成功| G[返回站点ID]
    D -->|遍历结束无匹配| H{有兜底站点?}
    H -->|是| I[返回兜底站点]
    H -->|否| J[返回无匹配]
    
    style A fill:#e1f5fe
    style G fill:#c8e6c9
    style I fill:#c8e6c9
    style J fill:#ffcdd2
```

---

## 图3-5 系统用例图

```mermaid
flowchart LR
    subgraph C端用户
        UC1[快速开通]
        UC2[提交需求]
        UC3[查看需求]
        UC4[取消需求]
        UC5[评价服务]
        UC6[查看通知]
    end
    
    subgraph B端用户
        UC10[登录系统]
        UC11[查看任务池]
        UC12[认领任务]
        UC13[完成任务]
        UC14[转派任务]
        UC15[管理围栏]
        UC16[管理用户]
        UC17[查看统计]
    end
    
    subgraph 系统
        SYS[地理围栏匹配]
        NOTIFY[通知推送]
    end
    
    UC2 --> SYS
    UC13 --> NOTIFY
```

---

## 图3-6 系统E-R图

```mermaid
erDiagram
    USER ||--o{ SERVICE_REQUEST : "提交"
    USER ||--o{ NOTIFICATION : "接收"
    USER ||--o{ USER_IDENTITY : "拥有"
    USER ||--o| CUSTOMER_PROFILE : "关联"
    
    SERVICE_STATION ||--o{ SERVICE_ZONE : "定义"
    SERVICE_STATION ||--o{ TASK_ASSIGNMENT : "分配"
    SERVICE_STATION ||--o{ USER : "所属"
    
    SERVICE_ZONE }o--|| SERVICE_STATION : "属于"
    
    SERVICE_REQUEST ||--|| TASK_ASSIGNMENT : "对应"
    SERVICE_REQUEST }o--|| USER : "属于"
    SERVICE_REQUEST }o--o| SERVICE_STATION : "分配至"
    
    TASK_ASSIGNMENT }o--|| SERVICE_REQUEST : "关联"
    TASK_ASSIGNMENT }o--|| SERVICE_STATION : "所属"
    TASK_ASSIGNMENT }o--o| USER : "认领"
    
    ROLE ||--o{ ROLE_PERMISSION : "拥有"
    PERMISSION ||--o{ ROLE_PERMISSION : "被分配"
    
    USER {
        bigint id PK
        varchar phone UK
        varchar password_hash
        varchar name
        bigint station_id FK
        varchar status
    }
    
    SERVICE_STATION {
        bigint id PK
        varchar name
        varchar code UK
        varchar address
        decimal latitude
        decimal longitude
        varchar status
    }
    
    SERVICE_ZONE {
        bigint id PK
        bigint station_id FK
        varchar name
        json points
        int priority
        varchar status
    }
    
    SERVICE_REQUEST {
        bigint id PK
        varchar request_no UK
        bigint user_id FK
        varchar service_type
        varchar status
        decimal submit_location_lat
        decimal submit_location_lng
        bigint station_id FK
    }
    
    TASK_ASSIGNMENT {
        bigint id PK
        bigint request_id FK
        bigint station_id FK
        bigint staff_id FK
        varchar status
        datetime claimed_at
        datetime completed_at
    }
    
    NOTIFICATION {
        bigint id PK
        bigint user_id FK
        varchar title
        text body
        varchar type
        varchar channel
        varchar send_status
    }
    
    ROLE {
        bigint id PK
        varchar code UK
        varchar name
        varchar status
    }
    
    PERMISSION {
        bigint id PK
        varchar code UK
        varchar name
        varchar module
        varchar type
    }
    
    ROLE_PERMISSION {
        bigint id PK
        bigint role_id FK
        bigint permission_id FK
    }
```

---

## 图4-1 系统架构图

```mermaid
flowchart TB
    subgraph 客户端
        C1[C端用户<br/>Vue3 + PWA]
        C2[B端管理门户<br/>Vue3 + Element Plus]
    end
    
    subgraph 网关层
        NG[Nginx<br/>反向代理]
    end
    
    subgraph 应用层
        API[Go + Gin<br/>RESTful API]
        AUTH[JWT认证]
        RBAC[RBAC权限]
        GEO[地理围栏引擎]
    end
    
    subgraph 数据层
        MYSQL[(MySQL 8.0<br/>业务数据)]
        REDIS[(Redis 7.0<br/>缓存/Token黑名单)]
    end
    
    subgraph 外部服务
        AMAP[高德地图API]
        SMTP[邮件服务]
    end
    
    C1 --> NG
    C2 --> NG
    NG --> API
    API --> AUTH
    API --> RBAC
    API --> GEO
    API --> MYSQL
    API --> REDIS
    API --> AMAP
    API --> SMTP
    
    style C1 fill:#e1f5fe
    style C2 fill:#e8f5e9
    style GEO fill:#fff9c4
```

---

## 图4-2 后端分层架构图

```mermaid
flowchart LR
    subgraph 表示层
        H[Handler<br/>参数校验/响应格式化]
    end
    
    subgraph 业务层
        S[Service<br/>业务逻辑/事务处理]
    end
    
    subgraph 数据层
        R[Repository<br/>GORM数据库操作]
        M[Model<br/>GORM Gen生成]
    end
    
    subgraph 数据库
        DB[(MySQL)]
    end
    
    H --> S --> R --> M --> DB
    
    style H fill:#e3f2fd
    style S fill:#e8f5e9
    style R fill:#fff3e0
    style M fill:#fce4ec
```

---

## 图4-3 任务状态流转图

```mermaid
stateDiagram-v2
    [*] --> dispatched: 系统派单
    dispatched --> claimed: 工作人员认领
    dispatched --> cancelled: 取消需求
    claimed --> completed: 服务完成
    claimed --> transferred: 转派他人
    transferred --> claimed: 重新认领
    completed --> [*]
    cancelled --> [*]
```

---

## 图4-4 需求状态流转图

```mermaid
stateDiagram-v2
    [*] --> pending: 用户提交
    pending --> dispatched: 围栏匹配成功
    pending --> cancelled: 用户取消
    dispatched --> claimed: 任务被认领
    claimed --> processing: 开始服务
    processing --> completed: 服务完成
    processing --> cancelled: 服务取消
    completed --> [*]
    cancelled --> [*]
```

---

## 图5-1 JWT认证流程图

```mermaid
sequenceDiagram
    participant U as 用户
    participant C as 客户端
    participant S as 服务端
    participant R as Redis
    participant D as MySQL
    
    U->>C: 输入手机号/密码
    C->>S: POST /auth/login
    S->>D: 查询用户信息
    D-->>S: 返回用户数据
    S->>S: 验证密码(bcrypt)
    S->>S: 生成JWT Token
    S-->>C: 返回Token
    C->>C: 存储Token
    
    Note over C,S: 后续请求
    C->>S: API请求 + Authorization Header
    S->>S: 解析验证Token
    S->>R: 检查Token黑名单
    R-->>S: 未在黑名单
    S->>S: 执行业务逻辑
    S-->>C: 返回响应
```

---

## 图5-2 任务认领时序图

```mermaid
sequenceDiagram
    participant W as 工作人员
    participant B as B端前端
    participant S as 后端服务
    participant D as MySQL
    
    W->>B: 点击认领按钮
    B->>S: POST /tasks/:id/claim
    S->>S: 验证JWT Token
    S->>S: 检查权限(staff角色)
    S->>D: 查询任务状态
    D-->>S: status=dispatched
    S->>D: 更新任务(乐观锁)
    Note over S,D: WHERE id=? AND status='dispatched'
    D-->>S: 更新成功(affected=1)
    S->>D: 创建任务历史
    S->>D: 创建通知记录
    S-->>B: 返回认领成功
    B-->>W: 显示认领成功
```

---

## 使用说明

### 方法1：在线渲染导出

1. 访问 https://mermaid.live
2. 复制上面任意图表代码粘贴
3. 点击 "Download PNG" 或 "Download SVG" 导出图片
4. 插入论文

### 方法2：Typora 渲染

1. 打开 Typora
2. 粘贴图表代码（需要用 ```mermaid 包裹）
3. 右键图表 → 导出为图片

### 方法3：VS Code 预览

1. 安装 "Markdown Preview Mermaid Support" 插件
2. 预览时即可看到渲染效果
3. 右键可导出图片

### 方法4：Draw.io 导入

1. 打开 https://app.diagrams.net
2. 高级 → 插入 → 高级
3. 粘贴 Mermaid 代码
4. 编辑后导出

---

## 图5-3 射线法原理示意图

```mermaid
flowchart TD
    subgraph 说明
        L1["● = 待判断点 P"]
        L2["→ = 从P向右发射的水平射线"]
        L3["× = 射线与边的交点"]
    end

    subgraph 示例1["点在多边形内部（交点数=奇数）"]
        direction LR
        P1["● P"] -->|射线→| X1["× 交点1"]
        X1 --> X2["× 交点2"]
        X2 --> X3["× 交点3"]
        X3 --> INF1["→ ∞"]
    end

    subgraph 示例2["点在多边形外部（交点数=偶数）"]
        direction LR
        P2["● P"] -->|射线→| X4["× 交点1"]
        X4 --> X5["× 交点2"]
        X5 --> INF2["→ ∞"]
    end

    style P1 fill:#ff6b6b,color:#fff
    style P2 fill:#4ecdc4,color:#fff
    style X1 fill:#ffd93d
    style X2 fill:#ffd93d
    style X3 fill:#ffd93d
    style X4 fill:#ffd93d
    style X5 fill:#ffd93d
```

---

## 图5-3b 射线法算法流程图

```mermaid
flowchart TD
    A[输入: 点P 和 多边形顶点数组] --> B{顶点数 >= 3?}
    B -->|否| C[返回 false]
    B -->|是| D[初始化 inside = false<br/>j = n - 1]
    D --> E[i = 0]
    E --> F{i < n?}
    F -->|否| G[返回 inside]
    F -->|是| H{边的两端点分别在<br/>P的上方和下方?}
    H -->|否| I[j = i<br/>i = i + 1]
    H -->|是| J{P在射线与边<br/>交点的左侧?}
    J -->|否| I
    J -->|是| K[inside = !inside]
    K --> I
    I --> F

    style A fill:#e1f5fe
    style C fill:#ffcdd2
    style G fill:#c8e6c9
    style K fill:#fff9c4
```
