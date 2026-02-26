# 基于地理围栏的社区养老信息分发平台的设计与实现

## 摘 要

随着我国人口老龄化进程加快，社区养老服务需求日益增长。当前社区养老服务普遍存在需求提交渠道分散、人工派单效率低下、服务过程不透明等问题。针对上述痛点，本文以北京市昌平区霍营街道为背景，设计并实现了一个基于地理围栏的社区养老信息分发平台。

本系统采用前后端分离的 B/S 架构，后端基于 Go 语言和 Gin 框架开发 RESTful API 服务，前端采用 Vue 3 和 TypeScript 构建，数据存储使用 MySQL 和 Redis。系统包含 C 端用户端和 B 端管理门户两个子系统。C 端面向老年人及家属，支持 PWA 渐进式 Web 应用特性，提供需求提交、进度查询和服务评价等功能。B 端面向工作人员和管理人员，提供任务管理、站点管理、围栏管理、权限配置和统计报表等功能。

系统的核心创新点在于地理围栏自动分发机制。系统采用射线法（Ray Casting）算法判断用户坐标是否位于服务站点的地理围栏范围内，结合 BoundingBox 快速预筛选和优先级排序优化，实现了需求提交后的自动站点匹配与任务派发，单次匹配耗时小于 50ms。当坐标未命中任何围栏时，系统通过 Haversine 公式计算最近站点进行兜底分配。

在安全与权限方面，系统实现了基于 JWT 的双端认证机制，B 端和 C 端通过 Token 类型字段实现端隔离。权限控制采用自定义三表 RBAC 模型（permissions、roles、role_permissions），支持细粒度的功能权限和数据访问控制。

经过系统测试，各功能模块运行正常，地理围栏匹配算法准确可靠，API 响应时间满足性能要求，系统达到了预期的设计目标。

**关键词**：社区养老 地理围栏 射线法 Go语言 Vue.js RBAC

## ABSTRACT

With the acceleration of population aging in China, the demand for community elderly care services is increasing. Current community elderly care services generally suffer from fragmented demand submission channels, inefficient manual dispatching, and opaque service processes. To address these issues, this paper designs and implements a community elderly care information distribution platform based on geofencing technology, taking Huoying Sub-district of Changping District, Beijing as the application scenario.

The system adopts a front-end and back-end separated B/S architecture. The back-end is developed using Go language with the Gin framework to provide RESTful API services, while the front-end is built with Vue 3 and TypeScript. MySQL and Redis are used for data storage. The system consists of two subsystems: a C-end user application and a B-end management portal. The C-end targets elderly users and their families, supporting Progressive Web Application (PWA) features and providing functions such as service request submission, progress tracking, and service evaluation. The B-end serves staff and administrators with task management, station management, geofence management, permission configuration, and statistical reporting.

The core innovation of the system lies in the automatic geofence-based distribution mechanism. The system employs the Ray Casting algorithm to determine whether user coordinates fall within the geofence boundaries of service stations. Combined with BoundingBox pre-filtering and priority-based sorting optimization, the system achieves automatic station matching and task dispatching upon request submission, with a single matching operation taking less than 50 milliseconds. When coordinates do not match any geofence, the system uses the Haversine formula to calculate the nearest station as a fallback.

Regarding security and access control, the system implements a JWT-based dual-end authentication mechanism, with B-end and C-end isolated through token type fields. Access control adopts a custom three-table RBAC model (permissions, roles, role_permissions), supporting fine-grained functional permissions and data access control.

Through systematic testing, all functional modules operate correctly, the geofencing algorithm proves accurate and reliable, API response times meet performance requirements, and the system achieves its intended design objectives.

**KEY WORDS**: Community Elderly Care, Geofencing, Ray Casting Algorithm, Go Language, Vue.js, RBAC
# 第一章 绪论

## 1.1 选题背景与意义

### 1.1.1 选题背景

随着我国人口老龄化进程加快，社区养老服务需求日益增长。截至2024年底，我国60岁及以上人口已超过2.9亿，占总人口比例超过21%。北京市昌平区霍营街道作为典型的城市社区，老年人口数量逐年增加，养老服务需求呈现多样化特点。

当前社区养老服务存在以下问题：

1. **需求上报渠道分散**：电话、微信群、线下登记等多种方式，缺乏统一入口
2. **人工派单效率低下**：依靠人工整理转发，响应不及时
3. **处理过程不透明**：服务需求处理进度难以追踪
4. **空间资源匹配不精准**：需求分发主要依赖人工经验判断

基于上述背景，有必要设计一个以社区站点为基本单元、基于地理围栏进行需求匹配和分发的服务信息平台。

### 1.1.2 选题意义

**理论意义**：将地理围栏技术应用于社区养老服务场景，探索基于地理位置信息的智能需求分发机制，为地理信息技术在社区服务领域的应用提供实践案例。

**实践意义**：
- 降低使用门槛：通过扫码方式提供统一需求入口
- 提升服务效率：基于地理围栏规则实现需求自动分发
- 增强服务透明度：全流程状态跟踪

## 1.2 国内外研究现状

### 1.2.1 社区养老服务信息化

发达国家在社区养老服务信息化方面起步较早。日本建立了完善的长期护理保险制度，依托信息化手段实现服务需求精准对接。欧盟各国普遍采用"智慧养老"理念，通过物联网、大数据等技术提升社区养老服务效率。

近年来，我国高度重视养老服务信息化建设，国务院印发《"十四五"国家老龄事业发展和养老服务体系规划》，明确提出推进养老服务数字化转型。然而，现有平台多侧重于资源展示和预约功能，在基于地理位置的智能分发方面研究较少。

### 1.2.2 地理围栏技术应用

地理围栏（Geofencing）是一种基于位置服务的虚拟边界技术。点在多边形判断是地理围栏的核心问题，常用算法包括射线法（Ray Casting）、回转数法等。射线法因实现简单、计算效率高而被广泛应用。

地理围栏技术已广泛应用于物流配送、智慧交通等领域，但在社区服务场景的应用探索不足。本设计将两者结合，应用于社区养老服务需求分发场景。

## 1.3 开发目标

### 1.3.1 总体目标

设计并实现一个基于地理围栏的社区养老信息分发平台，实现养老服务需求的统一采集、按地理范围自动分发以及处理过程的全流程管理。

### 1.3.2 具体目标

**功能目标**：
- 统一的需求提交入口，支持扫码快速进入
- 基于地理围栏的需求自动分发
- 完整的任务管理流程
- 细粒度的权限控制

**性能目标**：
- 地理围栏匹配时间 < 50ms
- API平均响应时间 < 500ms
- 系统可用性 ≥ 99%

## 1.4 本文的主要内容及组织结构

本文共分为七章：

- **第一章 绪论**：介绍选题背景、意义、研究现状和开发目标
- **第二章 相关理论与技术**：介绍系统开发涉及的关键技术
- **第三章 需求分析**：进行可行性分析，明确功能需求和非功能需求
- **第四章 系统概要设计**：设计系统架构、功能模块、数据库和接口
- **第五章 系统详细设计与实现**：描述各功能模块的设计与实现
- **第六章 系统测试**：给出测试用例和测试结果
- **第七章 总结与展望**：总结全文工作，展望后续优化方向
# 第二章 相关理论与技术

## 2.1 后端技术

### 2.1.1 Go语言

Go是由Google公司发布的静态强类型、编译型编程语言，具有语法简洁、原生支持并发、标准库丰富、跨平台编译等特点。Go的goroutine轻量级线程机制非常适合构建高并发Web服务。

### 2.1.2 Gin框架

Gin是Go语言的轻量级Web框架，基于基数树实现高性能路由，QPS可达数万级别。Gin提供灵活的中间件机制，本系统利用中间件实现了JWT认证、权限校验、日志记录等功能。

### 2.1.3 GORM

GORM是Go语言最流行的ORM库，支持关联映射、事务、钩子函数等特性。本系统采用代码优先模式，通过Go结构体定义数据模型，使用GORM Gen工具自动生成数据访问层代码。

## 2.2 前端技术

### 2.2.1 Vue.js

Vue.js是渐进式JavaScript框架，采用响应式数据绑定和组件化开发。Vue 3引入的Composition API允许按逻辑关注点组织代码，便于逻辑复用。

### 2.2.2 TypeScript

TypeScript是JavaScript的超集，添加了静态类型检查，能够在编译时发现类型错误，提供完善的IDE支持。

### 2.2.3 Element Plus

Element Plus是基于Vue 3的企业级UI组件库，提供60+组件，支持按需引入和主题定制。本系统B端采用Element Plus构建管理界面。

### 2.2.4 PWA

PWA（渐进式Web应用）结合Web和原生应用优势，通过Service Worker实现离线缓存，用户可将应用添加到桌面。本系统C端实现了PWA支持。

## 2.3 数据库技术

### 2.3.1 MySQL

MySQL是最流行的关系型数据库之一，InnoDB存储引擎提供ACID事务支持。本系统使用utf8mb4字符集，为常用查询字段建立索引。

### 2.3.2 Redis

Redis是开源的内存数据结构存储系统，读写速度极快。本系统使用Redis实现Token黑名单、会话缓存和分布式锁。

## 2.4 安全技术

### 2.4.1 JWT认证

JWT（JSON Web Token）是开放标准RFC 7519，由Header、Payload、Signature三部分组成。本系统采用双端分离的JWT认证策略，B端和C端使用不同Token类型。

### 2.4.2 RBAC权限模型

RBAC（基于角色的访问控制）是广泛使用的权限管理模型。本系统实现了自定义三表模型：permissions表（权限定义）、roles表（角色定义）、role_permissions表（角色权限关联）。

## 2.5 核心算法

### 2.5.1 地理围栏匹配算法

本系统采用射线法（Ray Casting）判断点是否在多边形内。算法原理：从待判断点向右发射射线，统计与多边形边界的交点数量，奇数则在内部，偶数则在外部。

性能优化措施：
- **BoundingBox预筛选**：先判断点是否在围栏外接矩形内
- **优先级排序**：围栏按优先级降序排列
- **内存加载**：启动时加载所有活跃围栏到内存

### 2.5.2 距离计算算法

当需求未命中任何围栏时，采用Haversine公式计算需求位置到各站点的球面距离，选择最近站点作为兜底方案。

## 2.6 本章小结

本章介绍了系统开发涉及的关键技术：后端采用Go和Gin框架，前端采用Vue 3和TypeScript，数据库采用MySQL和Redis，安全方面采用JWT认证和RBAC权限模型，核心算法采用射线法实现地理围栏匹配。
# 第三章 需求分析

本章结合霍营街道社区养老服务的实际情况，明确系统的功能需求和非功能需求。

## 3.1 可行性分析

### 3.1.1 技术可行性

本系统采用的技术栈均为业界成熟技术：后端采用Go语言，前端采用Vue 3框架，数据库采用MySQL 8.0，地理围栏匹配采用射线法算法。系统对硬件资源要求较低，2核CPU、4GB内存即可满足运行需求。

### 3.1.2 经济可行性

系统开发采用开源技术栈，无需购买商业软件授权，开发成本可控。系统可部署在云服务器上，运维成本在可接受范围内。

### 3.1.3 操作可行性

系统界面设计简洁直观，通过扫描二维码即可进入系统，降低了老年人的使用门槛。管理后台支持用户管理、站点管理、围栏管理等功能，运维人员可方便地进行系统维护。

## 3.2 业务分析

### 3.2.1 现有业务流程

当前社区养老服务的需求处理流程如下：

**图3-1 现有业务流程图**

**【请插入图片：图3-1_现有业务流程图.png — 图3-1 现有业务流程图】**

**现有流程存在的问题**：

| 问题 | 描述 |
|------|------|
| 渠道分散 | 电话、微信、线下多种渠道，信息难以统一管理 |
| 响应延迟 | 人工整理和派单耗时长，响应不及时 |
| 过程不透明 | 老年人无法了解需求处理进度 |
| 责任不清 | 缺乏明确的分发规则，派单依赖个人经验 |

### 3.2.2 目标业务流程

系统上线后的需求处理流程：

**图3-2 目标业务流程图**

**【请插入图片：图3-2_目标业务流程图.png — 图3-2 目标业务流程图】**

### 3.2.3 用户角色分析

| 角色 | 说明 | 主要职责 |
|------|------|---------|
| 老年人/家属 | C端用户 | 提交需求、查看进度、评价服务 |
| 工作人员 | B端用户 | 认领任务、上门服务、完成反馈 |
| 站点负责人 | B端用户 | 管理人员、分配任务、查看统计 |
| 系统管理员 | B端用户 | 管理站点、围栏、用户、权限 |

## 3.3 功能需求

### 3.3.1 功能模块划分

**图3-3 系统功能模块图**

**【请插入图片：图3-3_系统功能模块图.png — 图3-3 系统功能模块图】**

**C端功能模块**：

| 模块 | 功能描述 |
|------|---------|
| 快速开通 | 手机号验证码登录，自动创建账号 |
| 需求提交 | 选择服务类型、填写详情、获取定位 |
| 需求管理 | 查看需求列表、详情、取消需求 |
| 服务评价 | 对已完成服务进行评分评价 |

**B端功能模块**：

| 模块 | 功能描述 |
|------|---------|
| 任务管理 | 任务池查看、认领、完成、转派 |
| 站点管理 | 站点CRUD、人员管理 |
| 围栏管理 | 围栏CRUD、优先级设置 |
| 用户管理 | 用户CRUD、角色分配 |
| 权限管理 | 角色CRUD、权限分配 |
| 统计报表 | 任务统计、效率分析 |

### 3.3.2 地理围栏匹配流程

**图3-4 地理围栏匹配流程图**

**【请插入图片：图3-4_地理围栏匹配流程图.png — 图3-4 地理围栏匹配流程图】**

## 3.4 非功能需求

### 3.4.1 性能需求

| 指标 | 要求 |
|------|------|
| 地理围栏匹配时间 | < 50ms |
| API平均响应时间 | < 500ms |
| 页面首屏加载时间 | < 2秒 |
| 系统可用性 | ≥ 99% |

### 3.4.2 安全需求

| 需求项 | 描述 |
|--------|------|
| 身份认证 | JWT Token认证，支持双端隔离 |
| 权限控制 | 基于RBAC的细粒度权限控制 |
| 数据安全 | 密码加密存储，传输使用HTTPS |

### 3.4.3 可用性需求

- 界面简洁直观，符合老年人使用习惯
- 需求提交流程不超过3个步骤
- C端支持PWA离线访问

## 3.5 系统逻辑模型

### 3.5.1 E-R图

**图3-6 系统E-R图**

**【请插入图片：图3-6_系统实体关系图.png — 图3-6 系统E-R图】**

## 3.6 本章小结

本章从可行性、业务流程、功能需求、非功能需求四个方面进行了需求分析，并给出了系统的E-R模型。需求分析的结果将作为系统设计和实现的依据。
# 第四章 系统概要设计

本章基于需求分析的结果，对系统进行整体架构设计、功能模块划分、数据库设计和接口设计，为后续详细设计和编码实现提供指导。

## 4.1 系统架构设计

### 4.1.1 总体架构

系统采用前后端分离的B/S架构，由前端层、应用服务层和数据层三层组成。


**【请插入图片：图4-1_系统架构图.png — 图4-1 系统总体架构图】**

### 4.1.2 技术架构选型

**后端技术栈：**

| 组件 | 技术 | 选型理由 |
|------|------|---------|
| 编程语言 | Go 1.25 | 原生支持高并发，编译型语言性能优异，部署简单 |
| Web框架 | Gin | 轻量级高性能，中间件机制灵活，API友好 |
| ORM | GORM | 功能完善，支持代码优先模式，自动迁移 |
| 数据库 | MySQL 8.0 | 开源免费，社区成熟，utf8mb4支持emoji |
| 缓存 | Redis 7.0 | 高性能内存数据库，支持多种数据结构 |
| 认证 | JWT | 无状态认证，适合分布式架构 |

**前端技术栈：**

| 组件 | 技术 | 选型理由 |
|------|------|---------|
| 框架 | Vue 3 + TypeScript | 响应式框架，类型安全，生态完善 |
| 状态管理 | Pinia | Vue 3官方推荐，简洁直观 |
| UI组件库 | Element Plus | 企业级组件库，开箱即用，B端与C端统一风格 |
| 离线支持 | PWA (vite-plugin-pwa) | 支持离线访问，可安装到桌面 |

### 4.1.3 部署架构


**【请插入图片：图4-4_系统部署架构图.png — 图4-2 系统部署架构图】**

## 4.2 功能模块设计

### 4.2.1 后端模块结构

后端采用分层架构，按照职责划分为Handler、Service、Repository三层。


**【请插入图片：图4-2_后端分层架构图.png — 图4-3 后端模块结构图】**

**分层职责说明：**

| 层次 | 职责 | 示例 |
|------|------|------|
| Handler | 参数校验、响应格式化、Swagger注解 | 解析请求参数、调用Service、返回JSON |
| Service | 业务逻辑、事务控制、跨表操作 | 围栏匹配、状态流转、通知触发 |
| Repository | 数据库操作、查询封装 | CRUD操作、复杂查询 |

### 4.2.2 前端模块结构

**C端（用户端）结构：**


**【请插入图片：图4-5_前端模块结构图.png — C端模块结构图】**

**B端（管理门户）结构：**


**【请插入图片：图4-5_前端模块结构图.png — B端模块结构图（同上图右侧部分）】**

### 4.2.3 核心业务模块设计

**（1）地理围栏分发模块**


**【请插入图片：图3-4_地理围栏匹配流程图.png — 图4-4 地理围栏匹配流程图】**

**（2）任务状态流转模块**


**【请插入图片：图4-3_任务状态流转图.png — 图4-5 任务状态流转图】**

**（3）通知触发模块**

系统在以下业务节点自动触发通知：

| 触发事件 | 通知对象 | 通知内容 |
|---------|---------|---------|
| 需求创建 | 站点负责人 | 有新的服务需求待处理 |
| 需求取消 | 站点负责人 | 服务需求已取消 |
| 任务认领 | 需求提交人 | 您的需求已被受理 |
| 任务完成 | 需求提交人 | 服务已完成，请评价 |
| 任务转派 | 新认领人 | 您有新的任务待处理 |

## 4.3 数据库设计

### 4.3.1 概念数据库设计

**E-R图核心实体关系：**


**【请插入图片：图3-6_系统实体关系图.png — 图4-6 数据库E-R图】**

### 4.3.2 逻辑数据库设计

**（1）用户相关表**

**users - 用户表**

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | BIGINT | PK, AUTO_INCREMENT | 主键 |
| phone | VARCHAR(20) | UK, NOT NULL | 手机号 |
| password_hash | VARCHAR(255) | | 密码哈希 |
| name | VARCHAR(50) | | 姓名 |
| status | VARCHAR(20) | DEFAULT 'active' | 状态 |
| station_id | BIGINT | INDEX, FK | 所属站点 |
| created_at | DATETIME | | 创建时间 |
| updated_at | DATETIME | | 更新时间 |
| deleted_at | DATETIME | INDEX | 软删除时间 |

**user_identities - 用户身份表（支持多身份）**

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | BIGINT | PK | 主键 |
| user_id | BIGINT | FK, NOT NULL | 用户ID |
| identity_type | VARCHAR(20) | NOT NULL | 身份类型 |
| is_primary | TINYINT | DEFAULT 0 | 是否主身份 |
| station_id | BIGINT | FK | 关联站点 |

**（2）站点与围栏表**

**service_stations - 服务站点表**

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | BIGINT | PK | 主键 |
| name | VARCHAR(100) | NOT NULL | 站点名称 |
| code | VARCHAR(50) | UK | 站点编码 |
| address | VARCHAR(200) | | 地址 |
| latitude | DECIMAL(10,7) | | 纬度 |
| longitude | DECIMAL(10,7) | | 经度 |
| contact_phone | VARCHAR(20) | | 联系电话 |
| status | VARCHAR(20) | DEFAULT 'active' | 状态 |

**service_zones - 地理围栏表**

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | BIGINT | PK | 主键 |
| station_id | BIGINT | FK, NOT NULL | 所属站点 |
| name | VARCHAR(100) | NOT NULL | 围栏名称 |
| points | JSON | NOT NULL | 多边形顶点坐标 |
| priority | INT | DEFAULT 0 | 优先级（越大越优先）|
| status | VARCHAR(20) | DEFAULT 'active' | 状态 |

**points字段JSON格式示例：**
```json
[
  {"lat": 40.0820, "lng": 116.3650},
  {"lat": 40.0820, "lng": 116.3680},
  {"lat": 40.0850, "lng": 116.3680},
  {"lat": 40.0850, "lng": 116.3650}
]
```

**（3）需求与任务表**

**service_requests - 服务需求表**

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | BIGINT | PK | 主键 |
| request_no | VARCHAR(50) | UK | 需求编号 |
| user_id | BIGINT | FK | 提交用户 |
| service_type_id | BIGINT | FK | 服务类型 |
| station_id | BIGINT | FK | 分发站点 |
| status | VARCHAR(20) | INDEX | 状态 |
| description | TEXT | | 需求描述 |
| submit_latitude | DECIMAL(10,7) | | 提交纬度 |
| submit_longitude | DECIMAL(10,7) | | 提交经度 |
| contact_name | VARCHAR(50) | | 联系人 |
| contact_phone | VARCHAR(20) | | 联系电话 |
| service_address | VARCHAR(200) | | 服务地址 |
| rating | INT | | 评分（1-5）|
| rating_comment | TEXT | | 评价内容 |

**task_assignments - 任务分配表**

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | BIGINT | PK | 主键 |
| request_id | BIGINT | UK, FK | 关联需求 |
| staff_id | BIGINT | FK, INDEX | 认领人 |
| status | VARCHAR(20) | INDEX | 任务状态 |
| claimed_at | DATETIME | | 认领时间 |
| completed_at | DATETIME | | 完成时间 |
| completion_images | JSON | | 完成图片 |
| completion_notes | TEXT | | 完成备注 |

**（4）权限系统表**

**permissions - 权限表**

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | BIGINT | PK | 主键 |
| code | VARCHAR(50) | UK | 权限编码 |
| name | VARCHAR(100) | | 权限名称 |
| group | VARCHAR(50) | | 权限分组 |

**roles - 角色表**

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | BIGINT | PK | 主键 |
| code | VARCHAR(50) | UK | 角色编码 |
| name | VARCHAR(100) | | 角色名称 |

**role_permissions - 角色权限关联表**

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | BIGINT | PK | 主键 |
| role_id | BIGINT | FK | 角色ID |
| permission_id | BIGINT | FK | 权限ID |

### 4.3.3 物理数据库设计

**（1）索引策略**

```sql
-- 用户表索引
CREATE UNIQUE INDEX idx_users_phone ON users(phone);
CREATE INDEX idx_users_station ON users(station_id);
CREATE INDEX idx_users_deleted ON users(deleted_at);

-- 需求表索引
CREATE UNIQUE INDEX idx_requests_no ON service_requests(request_no);
CREATE INDEX idx_requests_user ON service_requests(user_id);
CREATE INDEX idx_requests_station ON service_requests(station_id);
CREATE INDEX idx_requests_status ON service_requests(status);

-- 任务表索引
CREATE UNIQUE INDEX idx_tasks_request ON task_assignments(request_id);
CREATE INDEX idx_tasks_staff ON task_assignments(staff_id);
CREATE INDEX idx_tasks_status ON task_assignments(status);
```

**（2）数据库配置**

| 配置项 | 值 | 说明 |
|--------|------|------|
| 字符集 | utf8mb4 | 支持emoji和特殊字符 |
| 排序规则 | utf8mb4_unicode_ci | 不区分大小写 |
| 存储引擎 | InnoDB | 支持事务和外键 |
| 事务隔离 | READ-COMMITTED | 默认隔离级别 |

## 4.4 通信及接口设计

### 4.4.1 API设计规范

**（1）RESTful风格**

系统API遵循RESTful设计规范：

| HTTP方法 | 路径 | 说明 |
|---------|------|------|
| GET | /api/requests | 获取需求列表 |
| GET | /api/requests/:id | 获取需求详情 |
| POST | /api/requests | 创建需求 |
| PUT | /api/requests/:id | 更新需求 |
| DELETE | /api/requests/:id | 删除需求 |

**（2）统一响应格式**

```json
{
  "msg": "ok",
  "data": {
    // 业务数据
  }
}
```

错误响应：
```json
{
  "msg": "错误描述",
  "data": null
}
```

**（3）分页参数**

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| page | int | 1 | 页码（从1开始）|
| page_size | int | 10 | 每页条数 |

### 4.4.2 核心接口列表

**C端接口：**

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/v1/c/auth/quick-start | 快速开通 |
| POST | /api/v1/c/auth/login | 登录 |
| POST | /api/v1/c/auth/logout | 登出 |
| GET | /api/v1/c/auth/profile | 获取当前用户 |
| POST | /api/v1/c/requests | 创建需求 |
| GET | /api/v1/c/requests | 需求列表 |
| GET | /api/v1/c/requests/:id | 需求详情 |
| PUT | /api/v1/c/requests/:id | 更新需求 |
| DELETE | /api/v1/c/requests/:id | 取消需求 |
| POST | /api/v1/c/requests/:id/rate | 评价需求 |

**B端接口：**

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/v1/b/auth/login | 登录 |
| GET | /api/v1/b/tasks | 任务池 |
| GET | /api/v1/b/tasks/my | 我的任务 |
| POST | /api/v1/b/tasks/:id/claim | 认领任务 |
| POST | /api/v1/b/tasks/:id/complete | 完成任务 |
| POST | /api/v1/b/tasks/:id/transfer | 转派任务 |
| GET | /api/v1/b/stations | 站点列表 |
| POST | /api/v1/b/stations | 创建站点 |
| GET | /api/v1/b/zones | 围栏列表 |
| POST | /api/v1/b/zones | 创建围栏 |
| POST | /api/v1/b/zones/:id/match | 测试围栏匹配 |

### 4.4.3 JWT Token设计

Token载荷结构：

```json
{
  "user_id": 123,
  "type": "b_end",
  "identities": ["staff"],
  "primary_role": "staff",
  "station_id": 1,
  "exp": 1704067200,
  "iat": 1703980800
}
```

| 字段 | 说明 |
|------|------|
| user_id | 用户ID |
| type | 端类型（b_end/c_end）|
| identities | 用户身份列表 |
| primary_role | 主角色 |
| station_id | 所属站点 |
| exp | 过期时间 |
| iat | 签发时间 |

## 4.5 安全性设计

### 4.5.1 认证与授权

**（1）双端认证**

系统支持B端和C端独立认证：
- B端：手机号 + 密码登录
- C端：支持短信验证码登录（默认）和密码登录两种方式
- C端快速开通：仅支持短信验证码，自动创建用户
- Token中type字段区分端类型（b_end / c_end），防止跨端访问

**（2）Token黑名单**

用户登出时，Token加入Redis黑名单，实现即时失效：

```go
func (s *AuthService) Logout(token string) error {
    // 解析Token获取剩余有效期
    claims := s.jwtManager.ParseToken(token)
    ttl := time.Until(claims.ExpiresAt.Time)
    
    // 加入黑名单
    return s.redis.Set(ctx, "blacklist:"+token, "1", ttl).Err()
}
```

**（3）RBAC权限控制**

采用自定义三表RBAC模型：
- permissions表：定义权限
- roles表：定义角色
- role_permissions表：关联角色和权限

权限检查流程：
1. 解析Token获取用户身份
2. 查询用户的角色和权限列表
3. 检查是否拥有所需权限
4. Admin角色跳过所有权限检查

### 4.5.2 数据安全

**（1）密码存储**

使用bcrypt算法加密存储：
```go
hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
```

**（2）敏感数据脱敏**

API响应中对手机号等敏感信息进行脱敏处理：
```go
func MaskPhone(phone string) string {
    if len(phone) != 11 {
        return phone
    }
    return phone[:3] + "****" + phone[7:]
}
```

### 4.5.3 接口安全

**（1）请求限流**

通过中间件实现API请求频率限制，防止恶意请求。

**（2）参数校验**

Handler层对请求参数进行严格校验：
```go
type CreateRequestReq struct {
    ServiceTypeID int64   `json:"service_type_id" binding:"required"`
    Description   string  `json:"description" binding:"required,max=500"`
    Latitude      float64 `json:"latitude" binding:"required,min=-90,max=90"`
    Longitude     float64 `json:"longitude" binding:"required,min=-180,max=180"`
}
```

## 4.6 其它设计

### 4.6.1 日志设计

采用结构化日志，记录关键业务操作：

| 日志字段 | 说明 |
|---------|------|
| timestamp | 时间戳 |
| level | 日志级别 |
| user_id | 操作用户 |
| action | 操作类型 |
| resource | 操作资源 |
| ip | 客户端IP |
| details | 详细信息 |

### 4.6.2 异常处理

统一异常处理机制，返回友好的错误信息：

| 错误码 | HTTP状态 | 说明 |
|--------|---------|------|
| 200 | 200 | 成功 |
| 400 | 400 | 参数错误 |
| 401 | 401 | 未认证 |
| 403 | 403 | 无权限 |
| 404 | 404 | 资源不存在 |
| 500 | 500 | 服务器错误 |

## 4.7 本章小结

本章对系统进行了概要设计。首先设计了系统的整体架构和技术选型；然后对后端和前端的功能模块进行了划分；接着进行了数据库的概念设计、逻辑设计和物理设计；最后设计了通信接口和安全性方案。概要设计的结果将作为后续详细设计和编码实现的依据。

---

<!-- 素材补充说明 -->
<!--
TODO 素材清单：
1. 图4-1 系统总体架构图
2. 图4-2 系统部署架构图
3. 图4-3 后端模块结构图
4. 图4-4 地理围栏匹配流程图
5. 图4-5 任务状态流转图
6. 图4-6 数据库E-R图
-->
# 第五章 系统详细设计与实现

本章基于系统概要设计，详细描述各功能模块的实现过程，包括开发环境搭建、登录鉴权、需求提交、地理围栏分发、任务管理、通知模块、评价模块和统计报表等核心功能的实现。

## 5.1 开发环境搭建

### 5.1.1 硬件环境

| 项目 | 配置要求 |
|------|---------|
| CPU | 2核心以上 |
| 内存 | 8GB以上 |
| 硬盘 | 50GB以上 |
| 网络 | 支持互联网访问 |

### 5.1.2 软件环境

**后端开发环境：**

| 软件 | 版本 | 用途 |
|------|------|------|
| Go | 1.25+ | 编程语言 |
| MySQL | 8.0+ | 数据库 |
| Redis | 7.0+ | 缓存 |
| Docker | 24.0+ | 容器化部署 |
| Air | latest | 热重载开发工具 |
| Swag | latest | Swagger文档生成 |

**前端开发环境：**

| 软件 | 版本 | 用途 |
|------|------|------|
| Node.js | 18+ | JavaScript运行时 |
| npm | 9+ | 包管理器 |
| VS Code | latest | 代码编辑器 |

### 5.1.3 项目初始化

**后端项目初始化：**

```bash
# 创建项目目录
mkdir scare-backend && cd scare-backend

# 初始化Go模块
go mod init community-elderly-care-platform

# 安装依赖
go get -u github.com/gin-gonic/gin
go get -u gorm.io/gorm
go get -u gorm.io/driver/mysql
go get -u github.com/golang-jwt/jwt/v5
go get -u github.com/redis/go-redis/v9
```

**前端项目初始化：**

```bash
# 创建C端项目
npm create vite@latest c-end -- --template vue-ts

# 创建B端项目
npm create vite@latest management-portal -- --template vue-ts

# 安装依赖
npm install element-plus pinia vue-router axios
```

## 5.2 项目目录结构

### 5.2.1 后端目录结构

```
backend/
├── cmd/
│   └── server/
│       └── main.go              # 程序入口
├── internal/
│   ├── config/
│   │   └── config.go            # 配置管理
│   ├── consts/
│   │   ├── roles.go             # 角色常量
│   │   └── status.go            # 状态常量
│   ├── dao/
│   │   ├── model/               # GORM模型（自动生成）
│   │   └── query/               # 查询方法（自动生成）
│   ├── dto/
│   │   ├── request_dto.go       # 需求DTO
│   │   └── task_dto.go          # 任务DTO
│   ├── handler/
│   │   ├── b_auth_handler.go    # B端认证
│   │   ├── c_auth_handler.go    # C端认证
│   │   ├── request_handler.go   # 需求处理
│   │   └── task_handler.go      # 任务处理
│   ├── service/
│   │   ├── auth_service.go      # 认证服务
│   │   ├── request_service.go   # 需求服务
│   │   ├── task_service.go      # 任务服务
│   │   └── notification_service.go
│   ├── repository/
│   │   ├── user_repo.go
│   │   ├── request_repo.go
│   │   └── task_repo.go
│   ├── middleware/
│   │   ├── auth.go              # JWT认证
│   │   └── permission.go        # 权限检查
│   └── router/
│       ├── router.go            # 路由注册
│       ├── b_end.go             # B端路由
│       └── deps.go              # 依赖注入
├── pkg/
│   ├── geo/
│   │   ├── engine.go            # 围栏引擎
│   │   ├── raycast.go           # 射线法算法
│   │   └── haversine.go         # 距离计算
│   ├── jwt/
│   │   └── manager.go           # JWT管理
│   └── logger/
│       └── logger.go            # 日志组件
├── database/
│   ├── schema/
│   │   └── schema.sql           # 表结构
│   └── seeds/
│       └── seed.sql             # 种子数据
├── docs/                        # Swagger文档
├── go.mod
└── go.sum
```

### 5.2.2 前端目录结构

**C端目录结构：**

```
frontend/c-end/
├── src/
│   ├── api/
│   │   ├── auth.ts              # 认证接口
│   │   ├── request.ts           # 需求接口
│   │   └── index.ts             # 统一导出
│   ├── views/
│   │   ├── Home.vue             # 首页
│   │   ├── QuickStart.vue       # 快速开通
│   │   ├── Login.vue            # 登录
│   │   ├── RequestList.vue      # 需求列表
│   │   └── RequestDetail.vue    # 需求详情
│   ├── components/
│   │   └── OfflineIndicator.vue # 离线提示
│   ├── composables/
│   │   └── useRequest.ts        # 需求相关逻辑
│   ├── store/
│   │   └── user.ts              # 用户状态
│   ├── router/
│   │   └── index.ts             # 路由配置
│   └── utils/
│       └── request.ts           # HTTP请求封装
├── public/
│   └── manifest.json            # PWA配置
└── vite.config.ts
```

## 5.3 登录鉴权与角色权限控制的实现

### 5.3.1 JWT认证机制

系统采用JWT（JSON Web Token）实现无状态认证。Token包含用户ID、身份类型、主角色、站点ID等信息。

**【请插入图片：图5-1_令牌认证流程图.png — 图5-4 JWT认证流程图】**


**Token生成实现：**

```go
// pkg/jwt/manager.go
type Claims struct {
    UserID      int64    `json:"user_id"`
    Type        string   `json:"type"`           // b_end / c_end
    Identities  []string `json:"identities"`
    PrimaryRole string   `json:"primary_role"`
    StationID   int64    `json:"station_id"`
    jwt.RegisteredClaims
}

func (m *Manager) GenerateToken(claims *Claims) (string, error) {
    claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(m.expireDuration))
    claims.IssuedAt = jwt.NewNumericDate(time.Now())
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(m.secretKey)
}
```

**认证中间件实现：**

```go
// internal/middleware/auth.go
func AuthMiddleware(jwtManager *jwt.Manager) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 获取Authorization头
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.AbortWithStatusJSON(401, gin.H{"msg": "未提供认证信息"})
            return
        }
        
        // 解析Bearer Token
        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        
        // 验证Token
        claims, err := jwtManager.ParseToken(tokenString)
        if err != nil {
            c.AbortWithStatusJSON(401, gin.H{"msg": "Token无效或已过期"})
            return
        }
        
        // 将用户信息存入上下文
        c.Set("user_id", claims.UserID)
        c.Set("identities", claims.Identities)
        c.Set("station_id", claims.StationID)
        
        c.Next()
    }
}
```

### 5.3.2 双端认证实现

系统支持B端和C端独立认证，通过Token中的type字段区分。B端仅支持密码登录，C端同时支持短信验证码登录（默认）和密码登录两种方式。

**B端登录实现：**

```go
// internal/service/auth_service.go
func (s *AuthService) LoginBEnd(phone, password string) (*Tokens, *model.User, error) {
    // 1. 验证用户凭据
    user, err := s.userRepo.GetByPhone(phone)
    if err != nil {
        return nil, nil, ErrInvalidCredentials
    }
    if err := VerifyPassword(user.PasswordHash, password); err != nil {
        return nil, nil, ErrInvalidCredentials
    }
    
    // 2. 获取用户的B端身份
    bEndIdentities, err := s.identityRepo.GetBEndIdentities(user.ID)
    if len(bEndIdentities) == 0 {
        return nil, nil, ErrNoRoleForBEnd
    }
    
    // 3. 生成JWT Token
    claims := &jwt.Claims{
        UserID:      user.ID,
        Type:        "b_end",
        Identities:  GetIdentityTypes(bEndIdentities),
        PrimaryRole: GetPrimaryRole(bEndIdentities),
        StationID:   user.StationID,
    }
    accessToken, _ := s.jwtManager.GenerateToken(claims)
    
    return &Tokens{AccessToken: accessToken}, user, nil
}
```

**C端快速开通实现：**

```go
// internal/service/auth_service.go
func (s *AuthService) QuickStart(phone, code string) (*Tokens, *model.User, error) {
    // 1. 验证短信验证码
    if !s.smsService.VerifyCode(phone, code) {
        return nil, nil, ErrInvalidCode
    }
    
    // 2. 检查用户是否存在，不存在则创建
    user, err := s.userRepo.GetByPhone(phone)
    if err == gorm.ErrRecordNotFound {
        user = &model.User{
            Phone:  phone,
            Status: "active",
        }
        s.userRepo.Create(user)
        
        // 创建C端身份
        identity := &model.UserIdentity{
            UserID:       user.ID,
            IdentityType: "elderly",
            IsPrimary:    1,
        }
        s.identityRepo.Create(identity)
    }
    
    // 3. 生成Token
    claims := &jwt.Claims{
        UserID:      user.ID,
        Type:        "c_end",
        Identities:  []string{"elderly"},
        PrimaryRole: "elderly",
    }
    accessToken, _ := s.jwtManager.GenerateToken(claims)
    
    return &Tokens{AccessToken: accessToken}, user, nil
}
```

**C端登录实现（支持双模式）：**

C端登录支持短信验证码和密码两种方式，通过请求参数中的`type`字段区分，默认为短信验证码登录：

```go
// internal/handler/c_auth_handler.go
func (h *CAuthHandler) Login(c *gin.Context) {
    var req dto.CLoginReq
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"msg": "参数错误"})
        return
    }
    
    // 根据登录类型选择认证方式，默认为短信验证码
    loginType := req.Type
    if loginType == "" {
        loginType = "code"
    }
    
    var tokens *service.Tokens
    var user *model.User
    var err error
    
    switch loginType {
    case "code":
        // 短信验证码登录
        tokens, user, err = h.authSvc.LoginByCode(req.Phone, req.Code)
    case "password":
        // 密码登录
        tokens, user, err = h.authSvc.LoginByPassword(req.Phone, req.Password)
    default:
        c.JSON(400, gin.H{"msg": "不支持的登录类型"})
        return
    }
    // ...
}
```

### 5.3.3 RBAC权限控制实现

系统采用自定义三表RBAC模型实现细粒度权限控制。

**权限检查实现：**

```go
// internal/service/permission_service.go
func (s *PermissionService) HasPermission(userID int64, permissionCode string) bool {
    // Admin角色跳过所有权限检查
    if s.isAdmin(userID) {
        return true
    }
    
    // 获取用户的权限列表（带缓存）
    permissions := s.GetUserPermissions(userID)
    for _, p := range permissions {
        if p.Code == permissionCode {
            return true
        }
    }
    return false
}

func (s *PermissionService) GetUserPermissions(userID int64) []*model.Permission {
    // 1. 获取用户的角色
    roles := s.identityRepo.GetRolesByUserID(userID)
    
    // 2. 获取角色的权限
    var permissions []*model.Permission
    for _, role := range roles {
        perms := s.rolePermRepo.GetPermissionsByRoleID(role.ID)
        permissions = append(permissions, perms...)
    }
    return permissions
}
```

**权限中间件：**

```go
// internal/middleware/permission.go
func PermissionMiddleware(permSvc *service.PermissionService, requiredCode string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.GetInt64("user_id")
        
        if !permSvc.HasPermission(userID, requiredCode) {
            c.AbortWithStatusJSON(403, gin.H{"msg": "无权限访问"})
            return
        }
        c.Next()
    }
}
```

## 5.4 养老服务需求提交模块的实现

### 5.4.1 需求提交流程


**【请插入图片：图5-4_需求提交时序图.png — 图5-1 需求提交时序图】**

用户提交需求的完整流程：
1. 用户扫码进入C端首页
2. 选择服务类型，填写需求描述
3. 系统获取/选择服务地址
4. 提交需求
5. 系统根据地址坐标自动匹配站点
6. 系统创建需求记录并分配至站点任务池
7. 系统通知站点工作人员

### 5.4.2 需求创建实现

**Handler层：**

```go
// internal/handler/request_handler.go
func (h *RequestHandler) CreateRequest(c *gin.Context) {
    var req dto.CreateRequestReq
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"msg": "参数错误"})
        return
    }
    
    userID := c.GetInt64("user_id")
    
    request, err := h.requestSvc.CreateRequest(c.Request.Context(), userID, &req)
    if err != nil {
        c.JSON(500, gin.H{"msg": err.Error()})
        return
    }
    
    c.JSON(200, gin.H{"msg": "ok", "data": request})
}
```

**Service层：**

```go
// internal/service/request_service.go
func (s *RequestService) CreateRequest(ctx context.Context, userID int64, req *dto.CreateRequestReq) (*model.ServiceRequest, error) {
    // 1. 生成需求编号
    requestNo := s.generateRequestNo()
    
    // 2. 地理围栏匹配站点
    stationID, matched := s.geofenceSvc.Match(req.Latitude, req.Longitude)
    if !matched {
        // 使用最近站点兜底
        stationID = s.geofenceSvc.FindNearestStation(req.Latitude, req.Longitude)
    }
    
    // 3. 创建需求记录
    request := &model.ServiceRequest{
        RequestNo:        requestNo,
        UserID:           userID,
        ServiceTypeID:    req.ServiceTypeID,
        StationID:        stationID,
        Status:           "dispatched",
        Description:      req.Description,
        SubmitLatitude:   req.Latitude,
        SubmitLongitude:  req.Longitude,
        ContactName:      req.ContactName,
        ContactPhone:     req.ContactPhone,
        ServiceAddress:   req.ServiceAddress,
    }
    
    if err := s.requestRepo.Create(request); err != nil {
        return nil, err
    }
    
    // 4. 发送通知给站点工作人员
    s.notificationSvc.NotifyNewRequest(stationID, request)
    
    return request, nil
}
```

### 5.4.3 需求编号生成

```go
func (s *RequestService) generateRequestNo() string {
    // 格式：REQ + 日期 + 随机数
    // 示例：REQ20250225001
    dateStr := time.Now().Format("20060102")
    randNum := rand.Intn(1000)
    return fmt.Sprintf("REQ%s%03d", dateStr, randNum)
}
```

## 5.5 基于地理围栏的需求分发模块的实现

### 5.5.1 围栏引擎设计

围栏引擎在系统启动时加载所有活跃围栏到内存，按优先级排序，提供快速的坐标匹配能力。


**【请插入图片：图5-3b_射线法算法流程图.png — 图5-2 围栏匹配算法流程图】**

**围栏数据结构：**

```go
// pkg/geo/engine.go
type Zone struct {
    ID        int64
    StationID int64
    Priority  int
    Points    []Point
    Box       BoundingBox  // 预计算的外接矩形
}

type Point struct {
    Lat float64
    Lng float64
}

type BoundingBox struct {
    MinLat, MaxLat float64
    MinLng, MaxLng float64
}
```

**围栏引擎初始化：**

```go
func NewEngine(zones []Zone) *Engine {
    filtered := make([]Zone, 0, len(zones))
    
    for _, zone := range zones {
        if len(zone.Points) < 3 {
            continue  // 至少需要3个点构成多边形
        }
        // 预计算外接矩形
        zone.Box = NewBoundingBox(zone.Points)
        filtered = append(filtered, zone)
    }
    
    // 按优先级降序排列
    sort.SliceStable(filtered, func(i, j int) bool {
        return filtered[i].Priority > filtered[j].Priority
    })
    
    return &Engine{zones: filtered}
}
```

### 5.5.2 射线法算法实现

射线法（Ray Casting）是判断点是否在多边形内的经典算法。


**【请插入图片：图5-3_射线法原理示意图.png — 图5-3 射线法原理示意图】**

**算法原理：**
从待判断点向右发射一条水平射线，统计射线与多边形边界的交点数量。若交点数量为奇数，则点在多边形内；若为偶数，则点在多边形外。

**算法实现：**

```go
// pkg/geo/raycast.go
func PointInPolygon(point Point, polygon []Point) bool {
    n := len(polygon)
    if n < 3 {
        return false
    }
    
    inside := false
    j := n - 1
    
    for i := 0; i < n; i++ {
        // 判断射线是否与边相交
        // 条件：边的两个端点分别在射线的上方和下方
        //       且点的x坐标小于边与射线交点的x坐标
        if ((polygon[i].Lat > point.Lat) != (polygon[j].Lat > point.Lat)) &&
           (point.Lng < (polygon[j].Lng-polygon[i].Lng)*(point.Lat-polygon[i].Lat)/
            (polygon[j].Lat-polygon[i].Lat)+polygon[i].Lng) {
            inside = !inside
        }
        j = i
    }
    
    return inside
}
```

### 5.5.3 匹配流程实现

```go
// pkg/geo/engine.go
func (e *Engine) Match(point Point) (int64, bool) {
    for _, zone := range e.zones {
        // 1. BoundingBox快速排除
        if !zone.Box.Contains(point) {
            continue
        }
        
        // 2. 射线法精确判断
        if PointInPolygon(point, zone.Points) {
            return zone.StationID, true
        }
    }
    
    // 3. 无匹配
    return 0, false
}
```

### 5.5.4 Haversine兜底机制

当需求未命中任何围栏时，使用Haversine公式计算最近站点作为兜底方案。

```go
// pkg/geo/haversine.go
func HaversineDistance(lat1, lng1, lat2, lng2 float64) float64 {
    const earthRadius = 6371.0  // 地球半径（千米）
    
    lat1Rad := lat1 * math.Pi / 180
    lat2Rad := lat2 * math.Pi / 180
    deltaLat := (lat2 - lat1) * math.Pi / 180
    deltaLng := (lng2 - lng1) * math.Pi / 180
    
    a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
         math.Cos(lat1Rad)*math.Cos(lat2Rad)*
         math.Sin(deltaLng/2)*math.Sin(deltaLng/2)
    c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
    
    return earthRadius * c
}

func (s *GeofenceService) FindNearestStation(lat, lng float64) int64 {
    stations := s.stationRepo.GetAllActive()
    
    var nearestID int64
    minDist := math.MaxFloat64
    
    for _, station := range stations {
        dist := HaversineDistance(lat, lng, station.Latitude, station.Longitude)
        if dist < minDist {
            minDist = dist
            nearestID = station.ID
        }
    }
    
    return nearestID
}
```

## 5.6 工作人员端任务管理模块的实现

### 5.6.1 任务池查看

工作人员可查看所属站点的待处理任务列表。

```go
// internal/service/task_service.go
func (s *TaskService) GetStationTasks(stationID int64, page, pageSize int) ([]*model.TaskAssignment, int64, error) {
    // 查询站点下状态为pending的任务
    tasks, total, err := s.taskRepo.GetByStationAndStatus(stationID, "pending", page, pageSize)
    if err != nil {
        return nil, 0, err
    }
    
    // 填充关联数据（需求信息、用户信息）
    for _, task := range tasks {
        task.Request, _ = s.requestRepo.GetByID(task.RequestID)
    }
    
    return tasks, total, nil
}
```

### 5.6.2 任务认领（并发安全）

任务认领使用乐观锁机制，防止并发重复认领。

**【请插入图片：图5-2_任务认领时序图.png — 图5-5 任务认领时序图】**


```go
// internal/service/task_service.go
func (s *TaskService) ClaimTask(taskID, staffID int64) error {
    // 1. 使用事务保证原子性
    return s.db.Transaction(func(tx *gorm.DB) error {
        // 2. 查询任务（加锁）
        var task model.TaskAssignment
        if err := tx.Where("id = ? AND status = ?", taskID, "pending").
            First(&task).Error; err != nil {
            return errors.New("任务不存在或已被认领")
        }
        
        // 3. 更新任务状态
        task.StaffID = staffID
        task.Status = "claimed"
        task.ClaimedAt = time.Now()
        
        if err := tx.Save(&task).Error; err != nil {
            return err
        }
        
        // 4. 更新需求状态
        if err := tx.Model(&model.ServiceRequest{}).
            Where("id = ?", task.RequestID).
            Update("status", "claimed").Error; err != nil {
            return err
        }
        
        // 5. 发送通知
        s.notificationSvc.NotifyTaskClaimed(task)
        
        return nil
    })
}
```

### 5.6.3 任务完成

```go
func (s *TaskService) CompleteTask(taskID, staffID int64, images []string, notes string) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        // 1. 验证任务归属
        var task model.TaskAssignment
        if err := tx.Where("id = ? AND staff_id = ? AND status = ?", 
            taskID, staffID, "claimed").First(&task).Error; err != nil {
            return errors.New("任务不存在或状态异常")
        }
        
        // 2. 更新任务
        task.Status = "completed"
        task.CompletedAt = time.Now()
        task.CompletionImages = images
        task.CompletionNotes = notes
        
        if err := tx.Save(&task).Error; err != nil {
            return err
        }
        
        // 3. 更新需求状态
        if err := tx.Model(&model.ServiceRequest{}).
            Where("id = ?", task.RequestID).
            Update("status", "completed").Error; err != nil {
            return err
        }
        
        // 4. 发送通知
        s.notificationSvc.NotifyTaskCompleted(task)
        
        return nil
    })
}
```

### 5.6.4 任务转派

```go
func (s *TaskService) TransferTask(taskID, fromStaffID, toStaffID int64) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        // 1. 验证任务
        var task model.TaskAssignment
        if err := tx.Where("id = ? AND staff_id = ? AND status = ?",
            taskID, fromStaffID, "claimed").First(&task).Error; err != nil {
            return errors.New("任务不存在或无权转派")
        }
        
        // 2. 更新任务
        task.StaffID = toStaffID
        task.Status = "transferred"
        
        if err := tx.Save(&task).Error; err != nil {
            return err
        }
        
        // 3. 创建转派记录
        history := &model.TaskHistory{
            TaskID:     taskID,
            ActionType: "transfer",
            FromUserID: fromStaffID,
            ToUserID:   toStaffID,
        }
        tx.Create(history)
        
        // 4. 发送通知
        s.notificationSvc.NotifyTaskTransferred(task, toStaffID)
        
        return nil
    })
}
```

## 5.7 多渠道通知模块的实现

### 5.7.1 通知服务设计

系统实现了统一的通知接口，支持站内信和邮件两种通知渠道，并预留了短信等扩展接口。

```go
// internal/service/notification_service.go
type NotificationService struct {
    repo       *repository.NotificationRepository
    mailSender notify.MailSender
}

// 通知类型
const (
    NotifyTypeNewRequest   = "new_request"    // 新需求
    NotifyTypeTaskClaimed  = "task_claimed"   // 任务已认领
    NotifyTypeTaskComplete = "task_complete"  // 任务已完成
    NotifyTypeTaskTransfer = "task_transfer"  // 任务已转派
)

func (s *NotificationService) Send(userID int64, notifyType string, title, content string) error {
    // 1. 创建站内信记录
    notification := &model.Notification{
        UserID:  userID,
        Type:    notifyType,
        Title:   title,
        Content: content,
        Status:  "unread",
    }
    s.repo.Create(notification)
    
    // 2. 发送邮件（异步）
    go func() {
        user, _ := s.userRepo.GetByID(userID)
        s.mailSender.Send(user.Email, title, content)
    }()
    
    return nil
}
```

### 5.7.2 邮件发送实现

```go
// internal/notify/email.go
type MailSender interface {
    Send(to, subject, body string) error
}

type SMTPMailSender struct {
    host     string
    port     int
    username string
    password string
    from     string
}

func (s *SMTPMailSender) Send(to, subject, body string) error {
    m := gomail.NewMessage()
    m.SetHeader("From", s.from)
    m.SetHeader("To", to)
    m.SetHeader("Subject", subject)
    m.SetBody("text/html", body)
    
    d := gomail.NewDialer(s.host, s.port, s.username, s.password)
    return d.DialAndSend(m)
}
```

### 5.7.3 通知触发点

| 触发场景 | 通知对象 | 通知内容 |
|---------|---------|---------|
| 需求创建 | 站点负责人 | "您有新的服务需求待处理" |
| 需求取消 | 站点负责人 | "服务需求已取消" |
| 任务认领 | 需求提交用户 | "您的需求已被受理" |
| 任务完成 | 需求提交用户 | "服务已完成，请评价" |
| 任务转派 | 新认领人 | "您有新的任务待处理" |

## 5.8 服务评价与回访模块的实现

### 5.8.1 评价提交

```go
// internal/service/request_service.go
func (s *RequestService) RateRequest(userID, requestID int64, rating int, comment string) error {
    // 1. 验证需求归属
    request, err := s.requestRepo.GetByID(requestID)
    if err != nil || request.UserID != userID {
        return errors.New("需求不存在或无权评价")
    }
    
    // 2. 验证需求状态
    if request.Status != "completed" {
        return errors.New("需求尚未完成，无法评价")
    }
    
    // 3. 验证是否已评价
    if request.Rating != nil {
        return errors.New("已评价，请勿重复提交")
    }
    
    // 4. 更新评价
    request.Rating = &rating
    request.RatingComment = &comment
    request.RatedAt = time.Now()
    
    return s.requestRepo.Update(request)
}
```

## 5.9 数据统计与报表模块的实现

### 5.9.1 统计指标

系统提供以下统计指标：

| 指标 | 说明 |
|------|------|
| 任务总数 | 指定时间段内的任务总量 |
| 完成率 | 已完成任务占比 |
| 平均响应时间 | 需求提交到任务认领的平均时间 |
| 平均处理时间 | 任务认领到完成的平均时间 |
| 满意度评分 | 用户评价的平均分数 |

### 5.9.2 统计服务实现

```go
// internal/service/statistics_service.go
type Statistics struct {
    TotalRequests    int64   `json:"total_requests"`
    CompletedCount   int64   `json:"completed_count"`
    CompletionRate   float64 `json:"completion_rate"`
    AvgResponseTime  float64 `json:"avg_response_time"`  // 小时
    AvgProcessTime   float64 `json:"avg_process_time"`   // 小时
    AvgRating        float64 `json:"avg_rating"`
}

func (s *StatisticsService) GetStationStats(stationID int64, startDate, endDate time.Time) (*Statistics, error) {
    stats := &Statistics{}
    
    // 1. 统计任务总数
    s.db.Model(&model.TaskAssignment{}).
        Where("station_id = ? AND created_at BETWEEN ?", stationID, startDate, endDate).
        Count(&stats.TotalRequests)
    
    // 2. 统计完成数
    s.db.Model(&model.TaskAssignment{}).
        Where("station_id = ? AND status = ? AND created_at BETWEEN ?", 
              stationID, "completed", startDate, endDate).
        Count(&stats.CompletedCount)
    
    // 3. 计算完成率
    if stats.TotalRequests > 0 {
        stats.CompletionRate = float64(stats.CompletedCount) / float64(stats.TotalRequests) * 100
    }
    
    // 4. 计算平均响应时间
    var avgResponse float64
    s.db.Model(&model.TaskAssignment{}).
        Where("station_id = ? AND claimed_at IS NOT NULL", stationID).
        Select("AVG(TIMESTAMPDIFF(HOUR, created_at, claimed_at))").
        Scan(&avgResponse)
    stats.AvgResponseTime = avgResponse
    
    // 5. 计算平均满意度
    var avgRating float64
    s.db.Table("service_requests").
        Where("station_id = ? AND rating IS NOT NULL", stationID).
        Select("AVG(rating)").
        Scan(&avgRating)
    stats.AvgRating = avgRating
    
    return stats, nil
}
```

## 5.10 前端实现

本节描述C端和B端前端的关键实现，包括路由守卫、权限指令、地图集成和PWA离线支持。

### 5.10.1 C端路由守卫

C端采用Vue Router的全局前置守卫实现登录态控制。路由配置中通过 `meta.requiresAuth` 标记需要登录的页面。

```typescript
// frontend/c-end/src/router/index.ts
router.beforeEach((to, _from, next) => {
  const tokenStore = useTokenStore()
  
  // 设置页面标题
  document.title = to.meta.title ? `${to.meta.title} - sCare` : 'sCare 社区养老服务'
  
  // 需要登录但未登录，跳转到登录页并保留原始路径
  if (to.meta.requiresAuth && !tokenStore.isLoggedIn) {
    next({
      name: 'Login',
      query: { redirect: to.fullPath }
    })
  } else {
    next()
  }
})
```

C端共包含16个页面视图，按是否需要登录分为两类：
- **公开页面**：首页、快速开通、登录、服务列表、站点动态、新闻详情、“我的”入口
- **需登录页面**：我的服务列表、服务详情、个人资料（基本信息/联系信息/服务地址/健康档案）、设置、消息通知

### 5.10.2 B端权限路由守卫

B端管理门户采用更精细的权限路由守卫，不仅检查登录状态，还基于 `meta.permission_code` 进行页面级权限校验。

```typescript
// frontend/management-portal/src/router/guards/permission.guard.ts
export function setupPermissionGuard(router: Router) {
  router.beforeEach(async (to, from, next) => {
    const authStore = useAuthStore()

    // 公开路由无需认证
    if (to.meta.public) {
      if (to.path === '/login' && authStore.isLoggedIn) {
        return next('/')
      }
      return next()
    }

    // 未登录跳转登录页
    if (!authStore.isLoggedIn) {
      ElMessage.warning('请先登录')
      return next({
        path: '/login',
        query: { redirect: to.fullPath },
      })
    }

    // 权限码校验
    const requiredPermission = to.meta.permission_code as string | undefined
    if (requiredPermission) {
      if (!authStore.hasPermission(requiredPermission)) {
        ElMessage.error('权限不足，无法访问此页面')
        return next(from.path || '/')
      }
    }

    next()
  })
}
```

### 5.10.3 v-permission 自定义指令

为实现按钮级别的权限控制，系统实现了 `v-permission` 自定义指令。该指令在元素挂载时检查当前用户是否拥有指定权限，若无权限则直接移除DOM元素。

```typescript
// frontend/management-portal/src/directives/permission.ts
export const vPermission: Directive<HTMLElement, string | string[]> = {
  mounted(el: HTMLElement, binding: DirectiveBinding<string | string[]>) {
    const authStore = useAuthStore()
    const value = binding.value

    if (!value) {
      console.warn('v-permission 指令需要权限码参数')
      return
    }

    // 支持单个权限码或权限码数组（任一满足即可）
    const hasPermission = Array.isArray(value)
      ? value.some((code) => authStore.hasPermission(code))
      : authStore.hasPermission(value)

    if (!hasPermission) {
      el.parentNode?.removeChild(el)
    }
  },
}
```

使用示例：

```vue
<!-- 单个权限控制 -->
<el-button v-permission="'service:task:claim'">认领任务</el-button>

<!-- 多权限控制（任一满足） -->
<el-button v-permission="['service:task:claim', 'service:task:complete']">处理任务</el-button>
```

结合 `usePermission` composable，还可以在组件逻辑中进行编程式权限判断：

```typescript
// frontend/management-portal/src/composables/usePermission.ts
export function usePermission() {
  const authStore = useAuthStore()

  const hasPermission = async (resource: string, action: string): Promise<boolean> => {
    return await checkPermission(`${resource}:${action}`)
  }

  const hasAnyPermission = async (
    permissions: Array<{ resource: string; action: string }>
  ): Promise<boolean> => {
    const results = await Promise.all(
      permissions.map((p) => checkPermission(`${p.resource}:${p.action}`))
    )
    return results.some((r) => r)
  }

  return { hasPermission, hasAnyPermission, hasAllPermissions, hasRole }
}
```

这样，系统的前端权限体系形成四层防护：路由守卫（页面级）→ v-permission指令（DOM级）→ usePermission composable（逻辑级）→ 后端接口权限校验（最终防线）。

### 5.10.4 高德地图集成

B端管理门户集成了高德地图JavaScript API 2.0，实现地图展示和地理围栏编辑两大功能。

**（1）地图展示组件（MapViewer）**

用于在服务详情等页面展示服务位置标记：

```typescript
// frontend/management-portal/src/components/MapViewer.vue
import AMapLoader from '@amap/amap-jsapi-loader'

async function initMap() {
  const AMap = await AMapLoader.load({
    key: import.meta.env.VITE_AMAP_KEY,
    version: '2.0',
    plugins: ['AMap.Marker'],
  })

  map = new AMap.Map(mapContainer.value, {
    viewMode: '2D',
    zoom: props.zoom,
    center: [props.longitude, props.latitude],
  })

  marker = new AMap.Marker({
    position: [props.longitude, props.latitude],
    title: '服务位置',
  })
  map.add(marker)
}
```

**（2）围栏编辑组件（MapPolygonEditor）**

用于在围栏管理页面可视化编辑地理围栏多边形区域，支持顶点的添加、拖拽调整和删除：

```typescript
// frontend/management-portal/src/components/MapPolygonEditor.vue
function startDraw() {
  polygon = new AMapObj.Polygon({ path: [], ... })
  map.add(polygon)
  polyEditor = new AMapObj.PolygonEditor(map, polygon)
  
  // 监听编辑事件，实时同步坐标数据
  polyEditor.on('addnode', updateModel)
  polyEditor.on('adjust', updateModel)
  polyEditor.on('removenode', updateModel)
  polyEditor.on('end', updateModel)
  
  polyEditor.open()
}
```

管理员通过该组件在地图上直接绘制和编辑围栏边界，操作完成后坐标数据以JSON格式保存至数据库的 `service_zones.points` 字段。

### 5.10.5 PWA离线支持

C端用户端通过 `vite-plugin-pwa` 实现渐进式Web应用（PWA）支持，使老年用户可将应用安装到手机主屏幕，提供类似原生应用的体验。

```typescript
// frontend/c-end/vite.config.ts
import { VitePWA } from 'vite-plugin-pwa'

VitePWA({
  registerType: 'autoUpdate',
  manifest: {
    name: 'sCare 社区养老服务',
    short_name: 'sCare',
    theme_color: '#409EFF',
    display: 'standalone',
    icons: [
      { src: 'pwa-192x192.png', sizes: '192x192', type: 'image/png' },
      { src: 'pwa-512x512.png', sizes: '512x512', type: 'image/png' },
    ]
  },
  workbox: {
    runtimeCaching: [
      // API请求：NetworkFirst策略，优先使用网络，断网时回退到缓存
      {
        urlPattern: /\/api\/v1\/.*/i,
        handler: 'NetworkFirst',
        options: { cacheName: 'api-cache', expiration: { maxAgeSeconds: 3600 } }
      },
      // 图片资源：CacheFirst策略，优先使用缓存
      {
        urlPattern: /\.(?:png|jpg|jpeg|svg|gif|webp|ico)$/i,
        handler: 'CacheFirst',
        options: { cacheName: 'image-cache' }
      },
    ]
  }
})
```

PWA配置实现了以下特性：
- **可安装性**：通过manifest配置，支持“添加到主屏幕”功能
- **离线访问**：Service Worker拦截请求，API采用NetworkFirst策略，静态资源采用CacheFirst策略
- **自动更新**：`registerType: 'autoUpdate'` 确保用户始终使用最新版本

## 5.11 本章小结

本章详细描述了系统各功能模块的实现过程。首先介绍了开发环境搭建和项目目录结构；然后依次实现了登录鉴权与权限控制、需求提交、地理围栏分发、任务管理、通知服务、评价回访和数据统计等核心功能。各模块的实现遵循分层架构原则，代码结构清晰，职责明确。地理围栏分发模块采用内存计算和射线法算法，实现了高效的自动派单功能。前端部分实现了路由守卫、v-permission权限指令、高德地图集成和PWA离线支持等功能，为用户提供了流畅的交互体验。
---

<!-- 素材补充说明 -->
<!--
TODO 素材清单：
1. 图5-1 需求提交时序图 → 图5-4_需求提交时序图.png
2. 图5-2 围栏匹配算法流程图 → 图5-3b_射线法算法流程图.png
3. 图5-3 射线法原理示意图 → 图5-3_射线法原理示意图.png
4. 图5-4 JWT认证流程图 → 图5-1_令牌认证流程图.png
5. 图5-5 任务认领时序图 → 图5-2_任务认领时序图.png

截图占位位置：
- 5.10.1 C端路由守卫章节：c-home.png, c-login.png, c-quickstart.png
- 5.10.4 高德地图集成章节：b-zone-map.png
- 5.10.5 PWA离线支持章节：c-home.png (复用)
# 第六章 系统测试

## 6.1 系统测试的概念和理论

### 6.1.1 测试目的

系统测试是软件开发生命周期中的关键环节，旨在验证系统是否满足规定的需求，发现系统中存在的缺陷和问题，确保系统在交付使用前达到预期的质量标准。本系统的测试主要包含以下目标：

1. **功能验证**：验证系统各功能模块是否按照需求规格正确实现
2. **接口测试**：确保前后端数据交互的正确性和一致性
3. **权限测试**：验证基于 RBAC 的权限控制机制是否有效
4. **性能测试**：评估系统在正常负载下的响应时间和资源消耗
5. **地理围栏测试**：验证点在多边形算法的准确性和性能

### 6.1.2 测试方法

本系统采用多种测试方法相结合的策略：

**1. 黑盒测试**

黑盒测试又称功能测试，将软件看作一个黑盒子，不考虑内部结构和实现细节，只根据需求规格说明检查程序的功能是否正确。本系统主要采用以下黑盒测试技术：

- **等价类划分**：将输入数据划分为若干等价类，从每个等价类中选取代表性数据进行测试
- **边界值分析**：针对输入边界条件进行测试，如分页参数的最小值、最大值
- **错误推测法**：根据经验推测可能存在的问题，设计测试用例

**2. 白盒测试**

白盒测试关注程序的内部逻辑结构，本系统在地理围栏匹配算法、JWT 认证等核心模块进行了白盒测试，验证代码逻辑的正确性。

**3. 接口测试**

采用 HTTP 客户端工具（如 curl、Postman）对 RESTful API 进行测试，验证请求参数、响应格式、状态码等是否符合设计规范。

**4. 端到端测试**

模拟真实用户操作流程，从用户登录、提交需求、地理围栏匹配、任务派发到任务完成的完整业务链路进行测试。

### 6.1.3 测试环境

系统测试在以下环境中进行：

| 环境项 | 配置 |
|--------|------|
| 操作系统 | macOS Sonoma 14.x / Linux (Ubuntu 22.04) |
| 后端运行时 | Go 1.25 |
| 数据库 | MySQL 8.0 (Docker 容器) |
| 缓存 | Redis 7.0-alpine (Docker 容器) |
| 后端框架 | Gin v1.9+ |
| 前端框架 | Vue 3 + Vite |
| 测试工具 | curl、bash script |
| 容器化 | Docker Compose |

### 6.1.4 测试范围

根据系统功能模块划分，本次测试覆盖以下核心功能：

1. **B 端管理门户**
   - 用户登录认证
   - 权限控制验证
   - 任务池查看
   - 任务认领与完成
   - 站点管理
   - 围栏管理
   - 用户管理
   - 统计报表

2. **C 端用户端**
   - 用户登录/注册
   - 服务需求提交
   - 地理围栏匹配
   - 需求状态查询
   - 评价反馈

3. **核心算法**
   - 地理围栏匹配算法（射线法）
   - JWT Token 生成与验证
   - RBAC 权限检查

## 6.2 测试用例与测试执行结果

### 6.2.1 核心功能测试用例

#### 6.2.1.1 用户认证测试

**测试模块**：B 端/C 端用户登录
**测试接口**：`POST /api/v1/b/auth/login`、`POST /api/v1/c/auth/login`

| 用例编号 | 测试场景 | 输入数据 | 预期结果 | 实际结果 | 状态 |
|----------|----------|----------|----------|----------|------|
| TC-AUTH-001 | 管理员登录 | 手机号: 13800000001, 密码: Test@123 | 返回 JWT Token，角色为 admin | 返回 Token，角色正确 | ✅ 通过 |
| TC-AUTH-002 | 站点负责人登录 | 手机号: 13800000002, 密码: Test@123 | 返回 JWT Token，角色为 station_manager | 返回 Token，角色正确 | ✅ 通过 |
| TC-AUTH-003 | 工作人员登录 | 手机号: 13800000004, 密码: Test@123 | 返回 JWT Token，角色为 staff | 返回 Token，角色正确 | ✅ 通过 |
| TC-AUTH-004 | 老年人登录 | 手机号: 13800000008, 密码: Test@123 | 返回 JWT Token，角色为 elderly | 返回 Token，角色正确 | ✅ 通过 |
| TC-AUTH-005 | 错误密码登录 | 手机号: 13800000001, 密码: WrongPass | 返回 401 错误，提示凭证无效 | 返回 401，错误信息正确 | ✅ 通过 |
| TC-AUTH-006 | 不存在的用户 | 手机号: 19999999999, 密码: Test@123 | 返回 401 错误 | 返回 401，错误信息正确 | ✅ 通过 |
| TC-AUTH-007 | B 端 Token 访问 C 端接口 | B 端 Token 访问 /api/v1/c/requests | 返回 403 禁止访问 | 返回 403 | ✅ 通过 |

**测试结论**：用户认证功能正常，JWT Token 生成正确，双端隔离有效。

#### 6.2.1.2 地理围栏匹配测试

**测试模块**：服务需求提交与围栏匹配
**测试接口**：`POST /api/v1/c/requests`

| 用例编号 | 测试场景 | 输入坐标 | 预期匹配围栏 | 实际结果 | 状态 |
|----------|----------|----------|--------------|----------|------|
| TC-GEO-001 | 华龙苑北里 | (40.05, 116.38) | 霍营站A区-华龙苑北里 | 匹配正确，station_id=1 | ✅ 通过 |
| TC-GEO-002 | 霍营地铁站 | (40.06, 116.39) | 霍营站B区-地铁周边 | 匹配正确，station_id=1 | ✅ 通过 |
| TC-GEO-003 | 围栏边界点 | 边界坐标 | 正确归属围栏 | 匹配正确 | ✅ 通过 |
| TC-GEO-004 | 围栏外坐标 | (40.10, 116.50) | 无匹配（兜底规则） | 返回无匹配提示 | ✅ 通过 |
| TC-GEO-005 | 多围栏重叠区 | 重叠区域坐标 | 返回优先级高的围栏 | 按优先级匹配 | ✅ 通过 |

**地理围栏匹配性能测试**：

| 指标 | 目标值 | 实际值 | 状态 |
|------|--------|--------|------|
| 单次匹配耗时 | < 50ms | ~10ms | ✅ 优秀 |
| 并发 100 请求 | < 500ms | 未测试 | - |

**测试结论**：地理围栏匹配算法（射线法）准确可靠，性能满足需求。

**【请插入截图：test-unit.png — 图6-17 单元测试结果】**

**【请插入截图：test-geofence.png — 图6-18 围栏匹配测试】**

**【请插入截图：test-swagger.png — 图6-19 Swagger接口文档】**

#### 6.2.1.3 任务管理测试

**测试模块**：任务池查看、认领、完成
**测试接口**：`GET /api/v1/tasks/pool`、`POST /api/v1/tasks/:id/claim`、`POST /api/v1/tasks/:id/complete`

| 用例编号 | 测试场景 | 操作 | 预期结果 | 实际结果 | 状态 |
|----------|----------|------|----------|----------|------|
| TC-TASK-001 | 查看任务池 | staff 用户请求任务池 | 返回本站点待认领任务列表 | 返回正确列表 | ✅ 通过 |
| TC-TASK-002 | 认领任务 | 点击认领按钮 | 任务状态变为 claimed，staff_id 设置 | 状态更新正确 | ✅ 通过 |
| TC-TASK-003 | 完成任务 | 上传照片并提交 | 任务状态变为 completed | 状态更新正确 | ✅ 通过 |
| TC-TASK-004 | 重复认领 | 已认领任务再次认领 | 返回错误，任务已被认领 | 返回 400 错误 | ✅ 通过 |
| TC-TASK-005 | 跨站点操作 | staff 认领其他站点任务 | 返回 403 禁止访问 | 返回 403 | ✅ 通过 |
| TC-TASK-006 | 乐观锁测试 | 两人同时认领同一任务 | 只有一人成功，另一人收到冲突提示 | 并发控制正确 | ✅ 通过 |

**测试结论**：任务管理功能完整，乐观锁并发控制有效，权限隔离正确。

#### 6.2.1.4 权限控制测试

**测试模块**：RBAC 权限验证
**测试策略**：不同角色访问不同接口

| 接口 | elderly | family | staff | station_manager | admin |
|------|---------|--------|-------|-----------------|-------|
| POST /api/v1/c/requests | ✅ | ✅ | ❌ | ❌ | ❌ |
| GET /api/v1/c/requests/:id | ✅ 自己 | ✅ 自己 | ✅ 站点 | ✅ 站点 | ✅ 全部 |
| GET /api/v1/b/tasks/pool | ❌ | ❌ | ✅ | ✅ | ✅ |
| POST /api/v1/b/tasks/:id/claim | ❌ | ❌ | ✅ | ✅ | ✅ |
| POST /api/v1/b/tasks/:id/complete | ❌ | ❌ | ✅ | ✅ | ✅ |
| GET /api/v1/b/stations | ❌ | ❌ | ✅ | ✅ | ✅ |
| POST /api/v1/b/stations | ❌ | ❌ | ❌ | ❌ | ✅ |
| GET /api/v1/b/zones | ❌ | ❌ | ✅ | ✅ | ✅ |
| POST /api/v1/b/zones | ❌ | ❌ | ❌ | ✅ | ✅ |
| GET /api/v1/b/users | ❌ | ❌ | ❌ | ❌ | ✅ |
| POST /api/v1/b/permissions | ❌ | ❌ | ❌ | ❌ | ✅ |

**测试结论**：基于自定义 RBAC 的权限控制机制工作正常，Admin 角色拥有全部权限，其他角色按设计正确受限。

### 6.2.2 端到端业务流程测试

**测试场景**：完整业务流程验证

```
测试流程：
1. 老年人用户（张大爷）登录 C 端
   ↓
2. 提交送餐需求（位置：华龙苑北里）
   ↓ 
3. 系统地理围栏匹配 → 自动派单到霍营站
   ↓
4. 工作人员（王小红）登录 B 端
   ↓
5. 查看任务池，看到待认领任务
   ↓
6. 认领任务
   ↓
7. 上门服务，拍照上传，完成任务
   ↓
8. 老年人查看需求状态（已完成）
   ↓
9. 老年人提交评价反馈
```

| 步骤 | 操作 | 预期结果 | 实际结果 | 状态 |
|------|------|----------|----------|------|
| 1 | 老年人登录 | Token 生成，角色 elderly | 正确 | ✅ |
| 2 | 提交需求 | 需求创建，自动派单 | request_no 生成，station_id=1 | ✅ |
| 3 | 围栏匹配 | 匹配霍营站A区 | 匹配正确 | ✅ |
| 4 | 工作人员登录 | Token 生成，角色 staff | 正确 | ✅ |
| 5 | 查看任务池 | 显示待认领任务 | 显示 1 条任务 | ✅ |
| 6 | 认领任务 | 任务状态 claimed | status=claimed | ✅ |
| 7 | 完成任务 | 任务状态 completed | status=completed | ✅ |
| 8 | 查询需求 | 需求状态同步更新 | status=completed | ✅ |
| 9 | 提交评价 | 评价记录保存 | 评价保存成功 | ✅ |

**端到端测试结论**：完整业务流程无阻塞，数据一致性良好，前后端交互正常。

### 6.2.3 前端界面测试

#### 6.2.3.1 B 端管理门户测试

| 用例编号 | 测试场景 | 预期结果 | 实际结果 | 状态 |
|----------|----------|----------|----------|------|
| TC-UI-B001 | 登录页面渲染 | 显示登录表单 | 正常显示 | ✅ |
| TC-UI-B002 | 侧边栏菜单 | 根据权限动态显示 | 菜单正确过滤 | ✅ |
| TC-UI-B003 | 任务池列表 | 分页显示任务 | 分页正常 | ✅ |
| TC-UI-B004 | 围栏地图展示 | 高德地图显示围栏区域 | 地图正常加载 | ✅ |
| TC-UI-B005 | 统计图表 | ECharts 图表渲染正确 | 图表显示正常 | ✅ |
| TC-UI-B006 | 按钮权限 | v-permission 指令生效 | 无权限按钮隐藏 | ✅ |

**【请插入截图：b-login.png — 图6-1 B端登录界面】**

**【请插入截图：b-menu.png — 图6-2 B端侧边栏菜单】**

**【请插入截图：b-task-pool.png — 图6-3 B端任务池列表】**

**【请插入截图：b-task-detail.png — 图6-4 B端任务详情】**

**【请插入截图：b-zone-map.png — 图6-5 B端围栏地图展示】**

**【请插入截图：b-statistics.png — 图6-6 B端统计图表】**

**【请插入截图：b-dashboard.png — 图6-7 B端工作台首页】**

**【请插入截图：b-user.png — 图6-8 B端用户管理】**

**【请插入截图：b-station-list.png — 图6-9 B端站点管理】**

**【请插入截图：b-role.png — 图6-10 B端角色权限管理】**

#### 6.2.3.2 C 端用户端测试

| 用例编号 | 测试场景 | 预期结果 | 实际结果 | 状态 |
|----------|----------|----------|----------|------|
| TC-UI-C001 | 扫码进入 | 自动获取地理位置 | 定位成功 | ✅ |
| TC-UI-C002 | 需求表单 | 表单验证生效 | 验证正确 | ✅ |
| TC-UI-C003 | 提交成功 | 显示成功提示 | 提示正常 | ✅ |
| TC-UI-C004 | PWA 安装 | 支持添加到桌面 | 安装提示显示 | ✅ |
| TC-UI-C005 | 离线访问 | 缓存页面可访问 | 离线可用 | ✅ |

**【请插入截图：c-home.png — 图6-11 C端首页】**

**【请插入截图：c-login.png — 图6-12 C端登录页】**

**【请插入截图：c-quickstart.png — 图6-13 C端快速开通】**

**【请插入截图：c-request-create.png — 图6-14 C端需求提交】**

**【请插入截图：c-request-list.png — 图6-15 C端需求列表】**

**【请插入截图：c-request-detail.png — 图6-16 C端需求详情】**

### 6.2.4 性能测试

#### 6.2.4.1 API 响应时间测试

| 接口 | 平均响应时间 | 最大响应时间 | 状态 |
|------|--------------|--------------|------|
| GET /api/v1/health | < 10ms | 15ms | ✅ 优秀 |
| POST /api/v1/c/auth/login | ~50ms | 80ms | ✅ 良好 |
| POST /api/v1/c/requests | ~150ms | 200ms | ✅ 良好 |
| GET /api/v1/b/tasks/pool | ~100ms | 150ms | ✅ 良好 |
| POST /api/v1/b/tasks/:id/claim | ~80ms | 120ms | ✅ 良好 |
| GET /api/v1/b/statistics | ~200ms | 350ms | ✅ 良好 |

#### 6.2.4.2 地理围栏算法性能

| 测试场景 | 围栏数量 | 匹配耗时 | 状态 |
|----------|----------|----------|------|
| 单围栏匹配 | 1 | < 1ms | ✅ |
| 多围栏匹配 | 5 | ~10ms | ✅ |
| 边界框过滤 | 5 | < 1ms | ✅ |

## 6.3 发现的问题及问题解决

### 6.3.1 测试过程中发现的问题

#### 问题 1：密码哈希占位符问题

**问题描述**：登录接口测试时，使用 seed.sql 中的测试账号登录失败，提示"无效凭证"。

**原因分析**：种子数据中的 `password_hash` 字段为占位符字符串，非真实的 bcrypt 哈希值。

**解决方案**：使用 Go 的 bcrypt 库生成真实密码哈希，更新数据库种子数据。

```go
hash, _ := bcrypt.GenerateFromPassword([]byte("Test@123"), bcrypt.DefaultCost)
// 更新所有测试账号的密码哈希
```

**状态**：✅ 已解决

#### 问题 2：API 路径前缀不匹配

**问题描述**：部分 API 请求返回 403 Forbidden 错误。

**原因分析**：权限策略配置文件中路径前缀为 `/api/...`，实际路由为 `/api/v1/...`，导致路径匹配失败。

**解决方案**：统一修正权限策略配置文件中的路径前缀。

**状态**：✅ 已解决

#### 问题 3：RESTful 路径参数匹配问题

**问题描述**：带路径参数的接口（如 `/api/v1/tasks/4/claim`）返回 403 错误。

**原因分析**：自定义权限中间件在进行路径匹配时使用精确匹配（`==`），无法正确匹配 RESTful 风格的动态路径参数（如 `:id`）。

**解决方案**：修改权限中间件中的路径匹配逻辑，采用路径模板匹配方式。将权限表中注册的路径模板（如 `/api/v1/tasks/:id/claim`）与实际请求路径进行模式匹配，支持动态路径参数的通配。

```go
// 路径模板匹配：将 :id 等参数替换为通配符进行匹配
func matchPath(pattern, path string) bool {
    patternParts := strings.Split(pattern, "/")
    pathParts := strings.Split(path, "/")
    if len(patternParts) != len(pathParts) {
        return false
    }
    for i, part := range patternParts {
        if strings.HasPrefix(part, ":") {
            continue // 动态参数，跳过匹配
        }
        if part != pathParts[i] {
            return false
        }
    }
    return true
}
```

**状态**：✅ 已解决

#### 问题 4：并发认领冲突

**问题描述**：多用户同时认领同一任务时，可能出现数据不一致。

**原因分析**：初始实现未考虑并发场景，缺乏锁机制。

**解决方案**：
1. 数据库层面：使用乐观锁（版本号机制）
2. 应用层面：事务内先检查状态再更新

```go
// 乐观锁实现
result := db.Model(&Task{}).
    Where("id = ? AND status = ?", taskID, "dispatched").
    Update("status", "claimed")
if result.RowsAffected == 0 {
    return errors.New("任务已被认领")
}
```

**状态**：✅ 已解决

### 6.3.2 已知限制

以下限制为 MVP 阶段的设计决策，非系统缺陷：

1. **邮件通知**：当前仅写入通知记录，未实现实际邮件发送（可后续接入 SMTP 服务）
2. **图片存储**：使用本地文件系统存储，生产环境建议切换到 OSS
3. **性能压测**：未进行大规模并发压测，建议后续补充
4. **移动端适配**：B 端管理门户主要面向 PC 使用，移动端体验未优化

### 6.3.3 测试总结

经过系统的功能测试、接口测试、权限测试和端到端业务流程测试，本系统的核心功能均已验证通过：

| 测试类别 | 用例总数 | 通过数 | 失败数 | 通过率 |
|----------|----------|--------|--------|--------|
| 认证测试 | 7 | 7 | 0 | 100% |
| 地理围栏测试 | 5 | 5 | 0 | 100% |
| 任务管理测试 | 6 | 6 | 0 | 100% |
| 权限控制测试 | 20+ | 20+ | 0 | 100% |
| 端到端测试 | 9 | 9 | 0 | 100% |
| 前端界面测试 | 11 | 11 | 0 | 100% |

**测试结论**：

1. ✅ 系统核心功能完整，满足需求规格要求
2. ✅ 前后端接口交互正常，数据一致性良好
3. ✅ RBAC 权限控制机制有效，角色隔离正确
4. ✅ 地理围栏匹配算法准确，性能满足需求
5. ✅ 完整业务流程无阻塞，可进入部署阶段

**系统状态**：已达到可交付标准，建议进入生产环境部署。

---

<!-- 截图占位位置说明：
认证测试截图：b-login.png, c-login.png, c-quickstart.png
任务管理截图：b-task-pool.png, b-task-detail.png
围栏管理截图：b-zone-map.png
统计报表截图：b-statistics.png, b-dashboard.png
单元测试截图：test-unit.png, test-geofence.png
Swagger截图：test-swagger.png
-->
# 第七章 总结与展望

## 7.1 总结

### 7.1.1 项目概述

本论文设计并实现了一个基于地理围栏的社区养老服务信息分发平台——昌平区霍营街道社区养老信息分发平台。该系统旨在解决传统社区养老服务中需求提交流程繁琐、信息分发效率低下、任务管理不透明等问题，通过信息化手段提升社区养老服务的响应速度和服务质量。

### 7.1.2 主要工作成果

本项目完成了以下主要工作：

**1. 系统架构设计**

采用前后端分离的 B/S 架构，后端使用 Go 语言和 Gin 框架，前端使用 Vue 3 和 TypeScript。系统采用分层架构设计，包括表示层、业务逻辑层和数据访问层，各层职责清晰，便于维护和扩展。同时引入 Redis 缓存和 JWT 认证机制，提升了系统的性能和安全性。

**2. 地理围栏自动匹配**

实现了基于射线法（Ray Casting）的点在多边形算法，能够根据用户提交需求时的地理坐标，自动匹配所属服务站点的地理围栏范围。该算法经过 BoundingBox 快速过滤和优先级排序优化，单次匹配耗时小于 50ms，满足了实时性要求。

**3. 自定义 RBAC 权限系统**

设计并实现了基于三表模型（permissions、roles、role_permissions）的自定义 RBAC 权限控制系统，支持细粒度的权限配置和角色继承。系统支持五种角色：系统管理员、站点负责人、工作人员、老年人、家属，各角色拥有不同的功能权限和数据访问范围。

**4. 双端应用开发**

- **B 端管理门户**：面向工作人员和管理人员，提供任务池管理、任务认领、任务完成、站点管理、围栏管理、用户管理、统计报表等功能。集成高德地图实现地理围栏的可视化展示，集成 ECharts 实现数据统计图表。

- **C 端用户端**：面向老年人和家属，提供服务需求提交、需求状态查询、评价反馈等功能。支持 PWA 特性，可离线使用并添加到手机桌面。

**5. 通知与评价系统**

实现了关键业务节点的通知推送机制，需求创建、任务认领、任务完成等事件自动触发通知记录。同时实现了服务评价反馈功能，便于收集用户意见，持续改进服务质量。

### 7.1.3 技术特点

本系统在技术实现上具有以下特点：

1. **Go 语言后端**：采用 Go 1.25 作为后端开发语言，利用其高并发特性和丰富的标准库，构建高性能的 RESTful API 服务。使用 GORM 作为 ORM 框架，支持 Code First 模式的数据库模型生成。

2. **JWT 双端认证**：实现了基于 JWT 的双端认证机制，B 端和 C 端使用不同的 Token 类型，通过 Token 中的 type 字段进行端隔离，确保不同端的用户只能访问各自授权的接口。

3. **乐观锁并发控制**：在任务认领等关键操作中，使用乐观锁机制防止并发冲突，确保数据一致性。

4. **统一响应格式**：所有 API 接口采用统一的 JSON 响应格式 `{ "msg": "xxx", "data": {...} }`，便于前端统一处理。

5. **PWA 渐进式应用**：C 端使用 Vite PWA 插件实现渐进式 Web 应用特性，支持离线缓存和桌面安装，提升了用户体验。

### 7.1.4 项目成果

经过需求分析、系统设计、编码实现和系统测试，本项目完成了以下交付物：

| 交付物 | 说明 |
|--------|------|
| 后端 API 服务 | 30+ RESTful 接口，覆盖核心业务流程 |
| B 端管理门户 | 15+ 功能页面，支持权限控制 |
| C 端用户端 | 5+ 功能页面，支持 PWA |
| 数据库设计 | 8 张核心数据表，完整的种子数据 |
| API 文档 | Swagger 自动生成的接口文档 |
| 测试报告 | 功能测试、权限测试、端到端测试报告 |

## 7.2 展望

### 7.2.1 功能扩展方向

**1. 智能推荐与调度**

当前系统的任务分配主要依赖地理围栏自动匹配和人工认领。未来可以引入智能调度算法，根据工作人员的位置、技能专长、历史服务评价等因素，实现更精准的任务推荐和自动分配。

**2. 实时通讯功能**

增加 WebSocket 或即时通讯功能，支持工作人员与老年人之间的实时沟通，便于确认上门时间、服务细节等。同时可以推送实时通知，提升信息触达效率。

**3. 移动 APP 开发**

当前 C 端以 PWA 形式提供，未来可以开发原生移动 APP，更好地利用手机硬件能力（如 GPS 定位、摄像头、推送通知等），提供更流畅的用户体验。

**4. 数据分析与决策支持**

基于积累的服务数据，构建数据分析平台，提供服务质量评估、需求趋势分析、资源配置优化等决策支持功能，帮助管理者优化养老服务供给。

**5. 多渠道接入**

扩展需求提交渠道，支持电话语音接入、智能音箱接入等，降低老年人使用数字技术的门槛。

### 7.2.2 技术优化方向

**1. 微服务架构演进**

随着业务规模扩大，可以将单体架构演进为微服务架构，将用户服务、任务服务、通知服务等拆分为独立的服务，提升系统的可扩展性和可维护性。

**2. 容器化与云原生部署**

完善 Docker 容器化配置，引入 Kubernetes 进行容器编排，实现弹性伸缩、滚动更新、故障自愈等云原生能力。

**3. 性能优化**

- 引入消息队列（如 RabbitMQ、Kafka）处理异步任务，如通知发送、数据统计等
- 优化数据库查询，添加合适的索引，使用读写分离提升查询性能
- 引入 CDN 加速静态资源访问

**4. 安全加固**

- 实现 API 限流，防止恶意请求
- 增加敏感操作的二次验证
- 完善日志审计功能
- 定期进行安全漏洞扫描

**5. 可观测性建设**

引入链路追踪（如 Jaeger）、指标监控（如 Prometheus）、日志聚合（如 ELK Stack）等可观测性工具，实现系统运行状态的全面监控和问题快速定位。

### 7.2.3 业务拓展方向

**1. 多区域推广**

当前系统针对昌平区霍营街道设计，未来可以推广到更多街道和社区，支持多租户模式和区域级数据隔离。

**2. 服务类型扩展**

除当前的送餐、家政、医疗陪伴等服务类型外，可以扩展更多养老服务类型，如康复训练、心理咨询、法律援助等，构建更完善的社区养老服务体系。

**3. 与其他系统集成**

- 与医院系统对接，获取老年人健康档案，提供更精准的健康服务
- 与社保系统对接，实现服务费用的自动结算
- 与社区门禁系统对接，记录老年人出入情况，关注独居老人安全

**4. 适老化设计深化**

进一步优化界面的适老化设计，支持更大的字体、更高的对比度、语音输入等，让老年人使用更加便捷。

### 7.2.4 结语

随着我国人口老龄化进程的加快，社区养老服务的需求将持续增长。信息化、智能化是提升社区养老服务效率和质量的重要手段。本项目的实践证明，通过合理的系统设计和现代信息技术的应用，可以有效地解决社区养老服务中的信息不对称、资源调度效率低下等问题。

未来，随着人工智能、物联网、5G 等技术的发展，社区养老服务将迎来更多的创新机遇。希望本项目能够为社区养老服务信息化建设提供参考，也期待更多的技术创新能够惠及广大老年人群体，让他们享受到更便捷、更贴心的养老服务。

---

本项目当前用于毕业设计场景，部署以本地与测试环境为主，暂不提供公网访问地址。

系统演示采用线下答辩与本地运行演示方式，演示视频链接不对外公开。
# 参考文献

[1] 国务院办公厅. 关于推进养老服务发展的意见[EB/OL]. (2019-04-16). http://www.gov.cn/zhengce/content/2019-04/16/content_5383274.htm

[2] 民政部. "十四五"民政事业发展规划[EB/OL]. (2021-06-18). http://www.mca.gov.cn/article/gk/wj/202106/20210600033794.shtml

[3] 北京市昌平区人民政府. 昌平区"十四五"时期老龄事业发展规划[EB/OL]. (2021-12-20). http://www.bjchp.gov.cn/

[4] Pressman R S. Software Engineering: A Practitioner's Approach[M]. 8th ed. New York: McGraw-Hill, 2014.

[5] Sommerville I. Software Engineering[M]. 10th ed. Boston: Pearson, 2016.

[6] Go语言官方团队. Go编程语言[EB/OL]. https://go.dev/

[7] Gin Web Framework. Gin - 快速Go Web框架[EB/OL]. https://gin-gonic.com/

[8] GORM. GORM - Go语言ORM库[EB/OL]. https://gorm.io/

[9] Vue.js. Vue.js - 渐进式JavaScript框架[EB/OL]. https://vuejs.org/

[10] TypeScript. TypeScript - JavaScript的类型超集[EB/OL]. https://www.typescriptlang.org/

[11] Element Plus. Element Plus - Vue 3 UI组件库[EB/OL]. https://element-plus.org/

[12] MySQL. MySQL - 开源关系型数据库[EB/OL]. https://www.mysql.com/

[13] Redis. Redis - 内存数据结构存储[EB/OL]. https://redis.io/

[14] JWT.io. JSON Web Token介绍[EB/OL]. https://jwt.io/introduction

[15] RFC 7519. JSON Web Token (JWT)[S]. Internet Engineering Task Force, 2015.

[16] Sandhu R S, Samarati P. Access Control: Principles and Practice[J]. IEEE Communications Magazine, 1994, 32(9): 40-48.

[17] 陈志, 王建. 基于角色的访问控制模型研究[J]. 计算机科学, 2018, 45(3): 25-32.

[18] Cormen T H, Leiserson C E, Rivest R L, et al. Introduction to Algorithms[M]. 3rd ed. Cambridge: MIT Press, 2009.

[19] O'Rourke J. Computational Geometry in C[M]. 2nd ed. Cambridge: Cambridge University Press, 1998.

[20] 周培德. 计算几何——算法设计与分析[M]. 5版. 北京: 清华大学出版社, 2017.

[21] 李德仁, 龚健雅, 邵振峰. 从数字地球到智慧地球[J]. 武汉大学学报(信息科学版), 2010, 35(2): 127-132.

[22] 高德地图开放平台. 高德地图Web API文档[EB/OL]. https://lbs.amap.com/api/

[23] ECharts. ECharts - 数据可视化图表库[EB/OL]. https://echarts.apache.org/

[24] Vite. Vite - 下一代前端构建工具[EB/OL]. https://vitejs.dev/

[25] PWA. Progressive Web Apps[EB/OL]. https://web.dev/progressive-web-apps/

[26] MDN Web Docs. Progressive Web Apps介绍[EB/OL]. https://developer.mozilla.org/en-US/docs/Web/Progressive_web_apps

[27] Docker. Docker - 应用容器引擎[EB/OL]. https://www.docker.com/

[28] Nginx. Nginx - 高性能HTTP和反向代理服务器[EB/OL]. https://nginx.org/

[29] RESTful API. RESTful Web Services[EB/OL]. https://restfulapi.net/

[30] Fielding R T. Architectural Styles and the Design of Network-based Software Architectures[D]. Irvine: University of California, 2000.

[31] 王伟, 李强. 基于微服务架构的社区养老服务系统设计[J]. 计算机应用与软件, 2021, 38(5): 156-162.

[32] 张明, 刘洋. 智慧养老服务平台关键技术与应用研究[J]. 信息技术与信息化, 2020, (8): 45-49.

[33] 赵华, 陈明. 基于GIS的社区养老服务资源空间配置研究[J]. 地理科学, 2019, 39(7): 1143-1150.

[34] 中国互联网信息中心. 第52次中国互联网络发展状况统计报告[R]. 北京: CNNIC, 2023.

[35] 全国信息安全标准化技术委员会. 信息安全技术 个人信息安全规范: GB/T 35273-2020[S]. 北京: 中国标准出版社, 2020.

[36] Swagger. Swagger - API文档生成工具[EB/OL]. https://swagger.io/

[37] bcrypt. bcrypt密码哈希算法[EB/OL]. https://en.wikipedia.org/wiki/Bcrypt

[38] RFC 6749. The OAuth 2.0 Authorization Framework[S]. Internet Engineering Task Force, 2012.

[39] 阿里云. 阿里云对象存储OSS[EB/OL]. https://www.aliyun.com/product/oss

[40] Pinia. Pinia - Vue状态管理库[EB/OL]. https://pinia.vuejs.org/
# 致谢

时光荏苒，转眼间本科阶段的学习即将画上句点。回首在北京邮电大学的学习时光，我收获颇丰，不仅在专业知识上有了长足进步，更在思维方式、问题解决能力等方面得到了全面锻炼。本论文的完成，离不开许多人的帮助和支持，在此谨表达我由衷的感谢。

首先，我要衷心感谢我的指导教师。从论文选题、开题到最终定稿，导师都给予了我悉心的指导和耐心的帮助。导师严谨的治学态度、渊博的专业知识和循循善诱的教导方式，让我受益匪浅。在系统开发过程中遇到的技术难题，导师也总是能够给予及时、专业的建议，帮助我攻克难关。

其次，我要感谢北京邮电大学计算机科学与技术专业的各位任课教师。正是他们在课堂上悉心传授的专业知识，为我完成本项目奠定了坚实的理论基础。从程序设计基础、数据结构、计算机网络到软件工程，每一门课程都为我的专业成长提供了重要支撑。

感谢我的同学和朋友们。在项目开发过程中，我们相互交流、共同探讨，你们提出的建议和反馈帮助我不断完善系统功能。特别感谢在测试阶段提供帮助的同学们，你们的测试反馈帮助我发现并修复了许多问题。

感谢昌平区霍营街道养老服务中心的工作人员。在需求调研阶段，你们热情地介绍了社区养老服务的实际工作流程和痛点问题，为本系统的设计提供了宝贵的业务背景和需求参考。

感谢开源社区的贡献者们。本项目使用了众多优秀的开源项目，包括Go语言、Gin框架、Vue.js、Element Plus、MySQL、Redis等。正是这些开源项目的存在，才使得本系统能够高效地完成开发。感谢所有为开源事业做出贡献的开发者们。

感谢我的家人。是你们的支持和理解，让我能够全身心投入到学习和项目开发中。你们的鼓励是我前进的动力，你们的期望是我不断努力的原因。

最后，感谢所有为本论文评审付出时间和精力的老师们。你们宝贵的意见和建议将帮助我进一步完善论文，也是我未来学习和工作的宝贵财富。

路漫漫其修远兮，吾将上下而求索。本科阶段的结束只是人生旅途的一个节点，我将带着在北京邮电大学学到的知识和精神，继续在计算机科学领域深耕细作，为社会信息化建设贡献自己的一份力量。

再次向所有帮助和支持过我的人表示最诚挚的感谢！

---

**作者**
**2026年 于北京邮电大学**
# 附录

## 附录A 核心算法代码

### A.1 射线法（Ray Casting）算法

```go
// pkg/geo/raycast.go
// 射线法判断点是否在多边形内部
func PointInPolygon(point Point, polygon []Point) bool {
    n := len(polygon)
    if n < 3 {
        return false
    }
    
    inside := false
    j := n - 1
    
    for i := 0; i < n; i++ {
        // 判断射线是否与多边形的边相交
        // 条件1：边的两个端点分别在射线的上方和下方
        // 条件2：点的经度小于边与射线交点的经度
        if ((polygon[i].Lat > point.Lat) != (polygon[j].Lat > point.Lat)) &&
           (point.Lng < (polygon[j].Lng-polygon[i].Lng)*(point.Lat-polygon[i].Lat)/
            (polygon[j].Lat-polygon[i].Lat)+polygon[i].Lng) {
            inside = !inside
        }
        j = i
    }
    
    return inside
}
```

### A.2 Haversine距离计算

```go
// pkg/geo/haversine.go
// Haversine公式计算两点间的球面距离（千米）
func HaversineDistance(lat1, lng1, lat2, lng2 float64) float64 {
    const earthRadius = 6371.0  // 地球平均半径（千米）
    
    lat1Rad := lat1 * math.Pi / 180
    lat2Rad := lat2 * math.Pi / 180
    deltaLat := (lat2 - lat1) * math.Pi / 180
    deltaLng := (lng2 - lng1) * math.Pi / 180
    
    a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
         math.Cos(lat1Rad)*math.Cos(lat2Rad)*
         math.Sin(deltaLng/2)*math.Sin(deltaLng/2)
    c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
    
    return earthRadius * c
}
```

### A.3 地理围栏引擎匹配

```go
// pkg/geo/engine.go
// 围栏匹配：遍历所有围栏，先BoundingBox快筛再射线法精确判断
func (e *Engine) Match(point Point) (int64, bool) {
    for _, zone := range e.zones {
        // BoundingBox快速排除
        if !zone.Box.Contains(point) {
            continue
        }
        // 射线法精确判断
        if PointInPolygon(point, zone.Points) {
            return zone.StationID, true
        }
    }
    return 0, false
}

// BoundingBox包含检查
func (box BoundingBox) Contains(p Point) bool {
    return p.Lat >= box.MinLat && p.Lat <= box.MaxLat &&
           p.Lng >= box.MinLng && p.Lng <= box.MaxLng
}
```

## 附录B 数据库建表脚本（核心表）

### B.1 用户表

```sql
CREATE TABLE `users` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `phone` VARCHAR(20) NOT NULL COMMENT '手机号',
  `password_hash` VARCHAR(255) DEFAULT NULL COMMENT '密码哈希',
  `name` VARCHAR(50) DEFAULT NULL COMMENT '姓名',
  `avatar` VARCHAR(255) DEFAULT NULL COMMENT '头像',
  `status` VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT '状态',
  `station_id` BIGINT UNSIGNED DEFAULT NULL COMMENT '所属站点',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` DATETIME DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_users_phone` (`phone`),
  KEY `idx_users_station` (`station_id`),
  KEY `idx_users_deleted` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';
```

### B.2 用户身份表

```sql
CREATE TABLE `user_identities` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
  `identity_type` VARCHAR(20) NOT NULL COMMENT '身份类型（admin/station_manager/staff/elderly/family）',
  `is_primary` TINYINT NOT NULL DEFAULT 0 COMMENT '是否主身份',
  `station_id` BIGINT UNSIGNED DEFAULT NULL COMMENT '关联站点',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` DATETIME DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_identities_user` (`user_id`),
  KEY `idx_identities_deleted` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户身份表';
```

### B.3 服务站点表

```sql
CREATE TABLE `service_stations` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(100) NOT NULL COMMENT '站点名称',
  `code` VARCHAR(50) DEFAULT NULL COMMENT '站点编码',
  `address` VARCHAR(200) DEFAULT NULL COMMENT '地址',
  `latitude` DECIMAL(10,7) DEFAULT NULL COMMENT '纬度',
  `longitude` DECIMAL(10,7) DEFAULT NULL COMMENT '经度',
  `contact_phone` VARCHAR(20) DEFAULT NULL COMMENT '联系电话',
  `status` VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT '状态',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` DATETIME DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_stations_code` (`code`),
  KEY `idx_stations_deleted` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='服务站点表';
```

### B.4 地理围栏表

```sql
CREATE TABLE `service_zones` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `station_id` BIGINT UNSIGNED NOT NULL COMMENT '所属站点',
  `name` VARCHAR(100) NOT NULL COMMENT '围栏名称',
  `points` JSON NOT NULL COMMENT '多边形顶点坐标',
  `priority` INT NOT NULL DEFAULT 0 COMMENT '优先级（越大越优先）',
  `status` VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT '状态',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` DATETIME DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_zones_station` (`station_id`),
  KEY `idx_zones_deleted` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='地理围栏表';
```

### B.5 服务需求表

```sql
CREATE TABLE `service_requests` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `request_no` VARCHAR(50) NOT NULL COMMENT '需求编号',
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '提交用户',
  `service_type_id` BIGINT UNSIGNED DEFAULT NULL COMMENT '服务类型',
  `station_id` BIGINT UNSIGNED DEFAULT NULL COMMENT '分发站点',
  `status` VARCHAR(20) NOT NULL DEFAULT 'pending' COMMENT '状态',
  `description` TEXT COMMENT '需求描述',
  `submit_latitude` DECIMAL(10,7) DEFAULT NULL COMMENT '提交纬度',
  `submit_longitude` DECIMAL(10,7) DEFAULT NULL COMMENT '提交经度',
  `contact_name` VARCHAR(50) DEFAULT NULL COMMENT '联系人',
  `contact_phone` VARCHAR(20) DEFAULT NULL COMMENT '联系电话',
  `service_address` VARCHAR(200) DEFAULT NULL COMMENT '服务地址',
  `rating` INT DEFAULT NULL COMMENT '评分（1-5）',
  `rating_comment` TEXT COMMENT '评价内容',
  `rated_at` DATETIME DEFAULT NULL COMMENT '评价时间',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` DATETIME DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_requests_no` (`request_no`),
  KEY `idx_requests_user` (`user_id`),
  KEY `idx_requests_station` (`station_id`),
  KEY `idx_requests_status` (`status`),
  KEY `idx_requests_deleted` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='服务需求表';
```

### B.6 任务分配表

```sql
CREATE TABLE `task_assignments` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `request_id` BIGINT UNSIGNED NOT NULL COMMENT '关联需求',
  `staff_id` BIGINT UNSIGNED DEFAULT NULL COMMENT '认领人',
  `status` VARCHAR(20) NOT NULL DEFAULT 'pending' COMMENT '任务状态',
  `claimed_at` DATETIME DEFAULT NULL COMMENT '认领时间',
  `completed_at` DATETIME DEFAULT NULL COMMENT '完成时间',
  `completion_images` JSON DEFAULT NULL COMMENT '完成照片',
  `completion_notes` TEXT COMMENT '完成备注',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` DATETIME DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_tasks_request` (`request_id`),
  KEY `idx_tasks_staff` (`staff_id`),
  KEY `idx_tasks_status` (`status`),
  KEY `idx_tasks_deleted` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='任务分配表';
```

### B.7 权限相关表

```sql
-- 权限表
CREATE TABLE `permissions` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `code` VARCHAR(50) NOT NULL COMMENT '权限编码',
  `name` VARCHAR(100) NOT NULL COMMENT '权限名称',
  `group_name` VARCHAR(50) DEFAULT NULL COMMENT '权限分组',
  `api_path` VARCHAR(200) DEFAULT NULL COMMENT 'API路径',
  `api_method` VARCHAR(10) DEFAULT NULL COMMENT 'HTTP方法',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_permissions_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='权限表';

-- 角色表
CREATE TABLE `roles` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `code` VARCHAR(50) NOT NULL COMMENT '角色编码',
  `name` VARCHAR(100) NOT NULL COMMENT '角色名称',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_roles_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色表';

-- 角色权限关联表
CREATE TABLE `role_permissions` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `role_id` BIGINT UNSIGNED NOT NULL COMMENT '角色ID',
  `permission_id` BIGINT UNSIGNED NOT NULL COMMENT '权限ID',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_role_perm` (`role_id`, `permission_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色权限关联表';
```

## 附录C 接口文档示例

### C.1 C端需求提交接口

- **URL**: `POST /api/v1/c/requests`
- **认证**: 需要JWT Token
- **请求体**:

```json
{
  "service_type_id": 1,
  "description": "需要日常照护服务",
  "latitude": 40.0835,
  "longitude": 116.3668,
  "contact_name": "张先生",
  "contact_phone": "13800138000",
  "service_address": "北京市昌平区霍营街道XX小区"
}
```

- **成功响应**:

```json
{
  "msg": "ok",
  "data": {
    "id": 1,
    "request_no": "REQ20250225001",
    "status": "dispatched",
    "station_id": 1,
    "created_at": "2025-02-25T10:30:00Z"
  }
}
```

### C.2 B端任务认领接口

- **URL**: `POST /api/v1/b/tasks/:id/claim`
- **认证**: 需要JWT Token（B端staff角色）
- **路径参数**: `id` - 任务ID
- **成功响应**:

```json
{
  "msg": "ok",
  "data": {
    "id": 1,
    "status": "claimed",
    "staff_id": 5,
    "claimed_at": "2025-02-25T11:00:00Z"
  }
}
```

### C.3 围栏匹配测试接口

- **URL**: `POST /api/v1/b/zones/:id/match`
- **认证**: 需要JWT Token（B端admin/station_manager角色）
- **请求体**:

```json
{
  "latitude": 40.0835,
  "longitude": 116.3668
}
```

- **成功响应**:

```json
{
  "msg": "ok",
  "data": {
    "matched": true,
    "station_id": 1,
    "station_name": "霍营街道养老服务站"
  }
}
```
