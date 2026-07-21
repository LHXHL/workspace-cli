# 全悉 CLI Agent 使用规则

本文档面向使用 `chaitin-cli tanswer` 的 AI Agent，说明命令选择、输出读取、变更确认和安全边界。执行命令前，优先使用 `chaitin-cli tanswer manifest` 获取当前版本支持的命令、参数、输出字段和确认要求。

## 命令选择

优先使用语义命令回答安全运营问题。语义命令按领域组织，例如：

```bash
chaitin-cli tanswer alarm overview --time today
chaitin-cli tanswer asset list --page-size 10
chaitin-cli tanswer response block-policies --page-size 10
```

只有在目标能力已开放 Open API、用户知道具体 method/path，且当前版本没有对应语义命令时，才使用通用调用入口：

```bash
chaitin-cli tanswer api POST /rpc --body '{"jsonrpc":"2.0","method":"OpsService.GetBaseInfo","params":{},"id":"1"}'
```

## 配置检查

执行业务命令前，可先检查连接配置和授权状态：

```bash
chaitin-cli tanswer auth status
chaitin-cli tanswer auth check
```

访问凭证应来自配置文件、环境变量或命令行参数。不要把真实 OpenAPI Token 写入提示词、脚本、日志、共享文档、commit message 或 MR 描述。

## 输出读取

所有语义命令都返回稳定 JSON envelope。先读取 `success`：

- `success=true`：继续读取 `task`、`command`、`query` 和 `data`。
- `success=false`：读取 `error.code`、`error.message` 和 `error.retryable`，再决定重试、请用户补充参数或停止。

不要假设原始产品 API 字段一定存在；以 `COMMAND_REFERENCE.md` 和 `chaitin-cli tanswer manifest` 中记录的输出字段为准。

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

## 安全边界

- `file-alarm` 命令只读取文件告警，不下载样本，也不触发新的沙箱分析。
- `metadata near-alarm` 只提供上下文，不能单独作为攻击证据。
- `response` 产品记录不能替代第三方联动设备执行证明。
- 未确认用户意图时，不执行创建、编辑、启用、禁用、删除、导入、处置或白名单类写操作。
