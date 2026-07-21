# 全悉 AI 可读 CLI

`chaitin-cli tanswer` 是全悉面向人类操作者和 AI Agent 的命令行入口。它把高频安全运营动作封装为语义命令，并提供 `api <METHOD> <PATH>` 通用调用入口，用于访问用户已知且已授权的全悉 Open API。

## 为什么使用全悉 CLI？

- Agent 原生：命令按安全运营任务命名，输出稳定 JSON。
- 读写边界清晰：查询命令直接执行，写操作必须 preview 和 confirm。
- 覆盖核心安全运营场景：系统状态、威胁告警、文件告警、流量元数据、资产配置、安全策略和响应处置。
- 行为可预期：命令、参数、输出字段和变更确认要求均可通过 manifest 查询。

## 快速开始

```bash
export TANSWER_URL='https://<全悉 Web 端 IP>'
export TANSWER_API_KEY='<全悉 OpenAPI Token>'

chaitin-cli tanswer auth status
chaitin-cli tanswer auth check
chaitin-cli tanswer system status
chaitin-cli tanswer alarm overview --time today
```

命令默认输出 JSON。若目标环境使用自签或非标准证书，可设置：

```bash
chaitin-cli tanswer --insecure auth check
```

## 快速开始：AI Agent

Agent 应先读取命令清单，再按用户意图选择语义命令：

```bash
chaitin-cli tanswer manifest
chaitin-cli tanswer alarm high-priority --time 24h
chaitin-cli tanswer asset list --page-size 10
chaitin-cli tanswer response block-policies --page-size 10
```

当语义命令已覆盖目标任务时，不要优先使用 `tanswer api`。只有目标能力已开放 Open API、且当前版本没有对应语义命令时，才使用通用调用入口：

```bash
chaitin-cli tanswer api POST /rpc --body '{"jsonrpc":"2.0","method":"OpsService.GetBaseInfo","params":{},"id":"1"}'
chaitin-cli tanswer api GET /api/example --query '{"count":10,"offset":0}'
chaitin-cli tanswer api POST /rpc --body @./request.json
```

`api <METHOD> <PATH>` 输出 `status_code` 和 `raw` 原始响应；它不提供专属参数解释或业务摘要。

## 配置和认证

`chaitin-cli tanswer` 支持通过配置文件、环境变量和全局 flag 配置连接信息。

| 配置 | 说明 |
| --- | --- |
| `tanswer.url` / `TANSWER_URL` | 全悉控制台地址，例如 `https://<全悉 Web 端 IP>`。 |
| `tanswer.api_key` / `TANSWER_API_KEY` | OpenAPI Token。不要在脚本、日志或共享文档中暴露真实值。 |
| `--timeout` | 请求超时时间，默认 `30s`。 |
| `--insecure` | 是否跳过 TLS 证书校验，默认 `false`。 |

常用 flag：

```bash
chaitin-cli tanswer --url 'https://<全悉 Web 端 IP>' --api-key "$TANSWER_API_KEY" auth check
chaitin-cli tanswer --insecure auth check
```

## 命令层级

| 层级 | 命令形态 | 用途 |
| --- | --- | --- |
| Foundation | `chaitin-cli tanswer auth status` | 查看当前目标环境、Token 是否已配置和输出设置。 |
| Foundation | `chaitin-cli tanswer auth check` | 校验 OpenAPI Token 能否访问当前全悉环境。 |
| Manifest | `chaitin-cli tanswer manifest` | 输出 AI 可读命令清单、风险等级、字段和确认要求。 |
| Semantic shortcut | `chaitin-cli tanswer alarm overview` | 用稳定业务字段回答安全运营问题。 |
| Protected write | `chaitin-cli tanswer asset create --preview` | 先返回写操作预览，不直接修改产品状态。 |
| Open API 通用调用 | `chaitin-cli tanswer api <METHOD> <PATH>` | 访问用户已知且已授权的全悉 Open API。 |

语义命令按领域组织；Agent 可通过 `chaitin-cli tanswer manifest` 区分语义命令、风险等级和 Open API 通用调用入口。

## 常用场景

```bash
chaitin-cli tanswer system status
chaitin-cli tanswer alarm overview --time 24h --severity critical,high
chaitin-cli tanswer alarm detail --id '<doc_id>'
chaitin-cli tanswer file-alarm malicious --time 7d --page-size 10
```

```bash
chaitin-cli tanswer asset list --page-size 10
chaitin-cli tanswer asset group-tree --depth 2
chaitin-cli tanswer metadata protocol --protocol http --time today
chaitin-cli tanswer policy detection-whitelist --page-size 10
```

```bash
chaitin-cli tanswer response block-policies --page-size 10
chaitin-cli tanswer response block-records --time 24h --page-size 10
chaitin-cli tanswer asset create --name core-db --ip 192.0.2.10 --preview
chaitin-cli tanswer asset create --name core-db --ip 192.0.2.10 --confirm CONFIRM_ASSET_CREATE
```

完整命令表、参数和长示例见 [COMMAND_REFERENCE.md](./COMMAND_REFERENCE.md)。

## JSON 输出契约

成功输出：

```json
{
  "success": true,
  "task": "查看威胁告警概览",
  "command": "chaitin-cli tanswer alarm overview",
  "query": {},
  "data": {},
  "warnings": []
}
```

失败输出：

```json
{
  "success": false,
  "task": "查看威胁告警概览",
  "command": "chaitin-cli tanswer alarm overview",
  "error": {
    "code": "ALARM_OVERVIEW_FAILED",
    "message": "error detail",
    "retryable": true
  }
}
```

文件类命令返回文件路径、大小、格式和导入/导出摘要，不在日志中输出 Token。

## 受保护写操作

高影响写命令不能静默执行。首次调用应使用 `--preview`，输出必须包含 `requires_confirmation=true`、`confirmed=false`、`operation_type`、`target`、`change_summary`、`impact`、`risk_warnings` 和命令专属 `confirmation_token`。

确认执行时必须重新运行命令并传入完全一致的 `--confirm` 值，例如：

```bash
chaitin-cli tanswer response block-policy-create --name block-bad-ip --object 192.0.2.30 --duration 3600 --preview
chaitin-cli tanswer response block-policy-create --name block-bad-ip --object 192.0.2.30 --duration 3600 --confirm CONFIRM_RESPONSE_BLOCK_POLICY_CREATE
```

确认大小写敏感。执行结果必须包含 `confirmed=true`、`result`、`object` 和 `audit` 字段。

## 安全边界

- 不在脚本、日志、共享文档、commit message 或 MR 描述中暴露真实 Token、真实环境地址或真实内网地址。
- 查询命令可以直接执行；资产、安全策略、元数据配置和响应处置写命令必须先预览变更，再使用确认令牌执行，并返回审计信息。
- 使用 `chaitin-cli tanswer manifest` 查看当前版本支持的命令、参数、输出字段和确认要求。
- `file-alarm` 命令只读取文件告警，不下载样本，也不触发新的沙箱分析。
- `metadata near-alarm` 只提供上下文，不能单独作为攻击证据。
- `response` 产品记录不能替代第三方联动设备执行证明。

## 文档索引

- [COMMAND_REFERENCE.md](./COMMAND_REFERENCE.md)：完整命令参考和长示例。
- [agent-skill.md](./agent-skill.md)：AI Agent 使用规则。
- [manifest.go](./manifest.go)：`chaitin-cli tanswer manifest` 使用的 AI 可读命令清单。
