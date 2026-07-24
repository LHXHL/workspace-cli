# 全悉 CLI Agent 使用规则

本文档面向使用 `chaitin-cli tanswer` 的 AI Agent，说明标准流程、命令选择、输出读取、变更确认和安全边界。执行命令前，优先使用 `chaitin-cli tanswer manifest` 获取当前版本支持的命令、参数、输出字段和确认要求。

## 标准执行流程

首次接入一个全悉环境，或用户没有明确说明当前 CLI 已验证时，按下面顺序检查：

```bash
chaitin-cli tanswer auth status
chaitin-cli tanswer auth check
chaitin-cli tanswer manifest
```

`auth status` 只检查本地是否配置了目标地址和 Token，不校验权限。`auth check` 用只读接口验证 TokenAuth 链路。`manifest` 是命令层级、参数、枚举、风险等级、输出字段和确认要求的当前版本事实来源。

访问凭证应来自配置文件、环境变量或命令行参数。不要把真实 OpenAPI Token 写入提示词、脚本、日志、共享文档、commit message 或 MR 描述。

## 命令选择优先级

优先使用语义命令回答安全运营问题。语义命令按领域组织，例如：

```bash
chaitin-cli tanswer alarm overview --time today
chaitin-cli tanswer asset list --page-size 10
chaitin-cli tanswer response block-policies --page-size 10
```

只有在目标能力已开放 Open API、用户知道具体 method/path，且当前版本没有对应语义命令时，才使用通用调用入口：

```bash
chaitin-cli tanswer api POST /rpc --body '{"jsonrpc":"2.0","method":"OpsService.GetBaseInfo","params":{},"id":"1"}'
chaitin-cli tanswer api GET /api/example --query '{"count":10,"offset":0}'
chaitin-cli tanswer api POST /rpc --body @./request.json
```

`api` 输出 `status_code` 和 `raw`。不要假设 CLI 内置完整 Open API 文档；调用前必须由用户输入或已授权文档确认 method、path、query 和 body。不要根据源码片段、猜测或相似产品经验拼 RPC method。

## 输出读取

所有语义命令都返回稳定 JSON envelope。先读取 `success`：

- `success=true`：继续读取 `task`、`command`、`query` 和 `data`。
- `success=false`：读取 `error.code`、`error.message` 和 `error.retryable`，再决定重试、请用户补充参数或停止。

不要假设原始产品 API 字段一定存在；以 `COMMAND_REFERENCE.md` 和 `chaitin-cli tanswer manifest` 中记录的输出字段为准。

## 错误处理

遇到失败时先判断错误属于哪一类：

| 错误类型 | 处理方式 |
| --- | --- |
| 缺少地址或 Token | 提醒用户配置 `TANSWER_URL`、`TANSWER_API_KEY`，或在命令中传入 `--url`、`--api-key`。 |
| `TOKEN_CHECK_FAILED` | 检查地址、Token、证书、有效期、IP 访问策略、频率限制和角色权限。 |
| 参数错误 | 根据 `error.message` 和命令 `--help` 补齐必填参数或修正枚举值。 |
| 权限不足 | 停止执行写操作，请用户确认 OpenAPI Token 绑定角色是否具备对应权限。 |
| 可重试错误 | 只有 `error.retryable=true` 且操作是只读命令时，才考虑重试。 |

## 变更确认

资产、安全策略、元数据配置和响应处置等会修改产品状态的命令必须先预览变更，再使用确认令牌执行。

预览输出应包含：

- `requires_confirmation=true`
- `confirmed=false`
- `operation_type`
- `target`
- `change_summary`
- `impact`
- `risk_warnings`
- `confirmation_token`

不要自动推断或替用户填写确认令牌。只有用户或上游系统明确提供与预览一致的 `confirmation_token` 时，才能使用 `--confirm` 执行。确认令牌大小写敏感。

确认前需要向用户复述目标对象、影响范围、风险提示和确认令牌来源。用户只表达“可以”“执行吧”但没有提供准确确认令牌时，不执行写操作。

执行输出应包含：

- `confirmed=true`
- `result`
- `object`
- `audit`

## 领域路由

| 用户意图 | 推荐命令 |
| --- | --- |
| 检查 CLI 连接和授权 | `auth status`, `auth check` |
| 查看系统版本、License、节点和自检状态 | `system status` |
| 查看威胁告警概览、趋势、列表、详情和关联调查 | `alarm ...` |
| 查看恶意文件、Webshell、沙箱检测等文件告警 | `file-alarm ...` |
| 查询、导入、导出或维护资产配置 | `asset ...` |
| 查询流量元数据、告警附近上下文和采集配置 | `metadata ...` |
| 查询或维护检测白名单和自定义 IOC 情报 | `policy ...` |
| 查询或维护阻断策略、响应白名单、联动设备和自动响应记录 | `response ...` |

### 告警路由

- 总览、等级分布、攻击结果分布、攻击源 Top、受害对象 Top：使用 `alarm overview`。
- 趋势、峰值、按时间观察等级变化：使用 `alarm timeline`。
- 原始告警行、分页、按攻击源/受害对象/资产 IP/威胁名过滤：使用 `alarm list`。
- 值班优先处置、critical/high、成功或失陷告警：优先使用 `alarm high-priority`。
- 单条告警详情：只有拿到 `doc_id` 后使用 `alarm detail`。
- 围绕一个攻击源、受害对象或威胁名做归因摘要：使用 `alarm by-attacker`、`alarm by-victim`、`alarm by-threat`。
- 查同一告警附近是否有相关事件：只有拿到源告警 `doc_id` 后使用 `alarm related`。

### 文件告警路由

- 恶意文件、Webshell、沙箱检测风险概览：使用 `file-alarm overview`。
- 恶意文件列表：使用 `file-alarm malicious`。
- Webshell 结果：使用 `file-alarm webshell`。
- 已有沙箱检测结果、沙箱分数、运行环境：使用 `file-alarm sandbox`。
- 单条文件告警详情：只有拿到 `doc_id` 后使用 `file-alarm detail`。
- 不使用 `file-alarm` 下载样本、提交样本或触发新的沙箱分析。

### 资产路由

- 查配置资产、按名称/IP/MAC/类型/重要性/标签/资产组过滤：使用 `asset list`。
- 单资产详情：只有拿到资产 ID 后使用 `asset detail`。
- 查资产组、组 ID、层级和资产数量：使用 `asset group-tree`。
- 导入前需要模板：使用 `asset download-template`。
- 导出资产配置：使用 `asset export`。
- 创建、更新、删除、导入、批量维护、批量打标签、资产组管理和树层级调整都属于写操作，必须先预览并等待准确确认令牌。
- 不用资产命令回答漏洞风险、资产风险或风险主机问题；这些不在当前全悉 CLI 范围内。

### 元数据路由

- 按 HTTP、DNS、TCP、UDP 等协议查流量元数据：使用 `metadata protocol`。
- 用户给出高级查询表达式或精确条件：使用 `metadata search`。
- 单条元数据详情：只有拿到 `id`、`timestamp` 和 `protocol` 后使用 `metadata detail`。
- 告警附近流量上下文：只有拿到告警 `doc_id` 后使用 `metadata near-alarm`，结果只能作为调查上下文，不能单独当作攻击证据。
- 查询或调整协议采集配置：使用 `metadata config` 和 `metadata config-update`；调整配置属于写操作。

### 策略路由

- 查询误报抑制、检测白名单状态、源/目的过滤、域名、URL、User-Agent、XFF、响应码、威胁类型或检测规则 ID：使用 `policy detection-whitelist`。
- 从已确认误报告警生成检测白名单：只有拿到告警 `doc_id` 且用户明确确认误报后，使用 `policy detection-whitelist-from-alarm --preview`。
- 查询自定义 IOC 情报、IOC 类型、启停状态和备注：使用 `policy custom-intelligence`。
- 创建、更新、启用、禁用、删除、导入检测白名单或自定义 IOC 都是写操作，必须先预览并等待准确确认令牌。
- 不用 policy 命令维护响应白名单、旁路阻断策略或自动响应。

### 响应路由

- 查询旁路阻断策略：使用 `response block-policies`。
- 查询阻断命中记录：使用 `response block-records`。
- 查询响应白名单：使用 `response whitelist`。
- 查询联动设备：使用 `response devices`；不知道设备 ID 时先查设备，再查 `response device-records`。
- 查询自动响应策略和自动响应生成列表：使用 `response auto-policies`、`response auto-list`。
- 从告警生成阻断策略或响应白名单：只有拿到告警 `doc_id` 且用户明确选择处置动作后，使用 `response block-policy-from-alarm --preview` 或 `response whitelist-from-alarm --preview`。
- 不把全悉产品记录当作第三方联动设备已经执行成功的唯一证明。

## 安全边界

- `file-alarm` 命令只读取文件告警，不下载样本，也不触发新的沙箱分析。
- `metadata near-alarm` 只提供上下文，不能单独作为攻击证据。
- `response` 产品记录不能替代第三方联动设备执行证明。
- 未确认用户意图时，不执行创建、编辑、启用、禁用、删除、导入、处置或白名单类写操作。
- 不绕过 `chaitin-cli tanswer` 直接用 `curl` 执行全悉写操作。
- 不根据猜测调用未在 `manifest`、`COMMAND_REFERENCE.md` 或用户授权 Open API 文档中出现的接口。
