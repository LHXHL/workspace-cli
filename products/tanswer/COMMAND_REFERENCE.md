# 全悉 AI 可读 CLI 命令参考

本文档记录 `chaitin-cli tanswer` 的完整命令表、使用边界和长示例。命令名、flag、RPC method、JSON field、enum 和 path 保持英文。

## 全局配置

```bash
export TANSWER_URL='https://<全悉 Web 端 IP>'
export TANSWER_API_KEY='<全悉 OpenAPI Token>'
```

```bash
chaitin-cli tanswer --url 'https://<全悉 Web 端 IP>' --api-key "$TANSWER_API_KEY" auth check
```

## Foundation 和 Open API 通用调用

| 命令 | 用途 |
| --- | --- |
| `chaitin-cli tanswer auth status` | 查看当前连接配置，不校验 Token 权限。 |
| `chaitin-cli tanswer auth check` | 校验 OpenAPI Token 和基础 RPC 访问能力。 |
| `chaitin-cli tanswer manifest` | 输出 AI 可读命令清单。 |
| `chaitin-cli tanswer api <METHOD> <PATH>` | 访问用户已知且已授权的全悉 Open API。 |

## System

| 命令 | 用途 |
| --- | --- |
| `chaitin-cli tanswer system status` | 查询版本、License、节点状态和自检摘要。 |

## Threat alarm

| 命令 | 用途 |
| --- | --- |
| `chaitin-cli tanswer alarm overview` | 查询威胁告警概览。 |
| `chaitin-cli tanswer alarm timeline` | 查询告警趋势和严重性分桶。 |
| `chaitin-cli tanswer alarm list` | 查询威胁告警列表和分页信息。 |
| `chaitin-cli tanswer alarm high-priority` | 查询 critical/high 且成功或失陷的优先告警。 |
| `chaitin-cli tanswer alarm detail --id <doc_id>` | 查询单条威胁告警详情。 |
| `chaitin-cli tanswer alarm by-attacker --attacker <ip>` | 按攻击者汇总告警。 |
| `chaitin-cli tanswer alarm by-victim --victim <ip>` | 按受害对象汇总告警。 |
| `chaitin-cli tanswer alarm by-threat` | 按威胁名称、类型或阶段汇总告警。 |
| `chaitin-cli tanswer alarm important-assets` | 查询重要资产相关告警。 |
| `chaitin-cli tanswer alarm attacker-rank` | 查询攻击者排行。 |
| `chaitin-cli tanswer alarm victim-rank` | 查询受害对象排行。 |
| `chaitin-cli tanswer alarm phase-distribution` | 查询攻击阶段分布。 |
| `chaitin-cli tanswer alarm related --id <doc_id>` | 基于源告警查找相近相关告警。 |

## File alarm

| 命令 | 用途 |
| --- | --- |
| `chaitin-cli tanswer file-alarm overview` | 查询文件告警概览。 |
| `chaitin-cli tanswer file-alarm malicious` | 查询恶意文件告警。 |
| `chaitin-cli tanswer file-alarm webshell` | 查询 Webshell 文件告警。 |
| `chaitin-cli tanswer file-alarm sandbox` | 查询已有沙箱检测结果告警。 |
| `chaitin-cli tanswer file-alarm detail --id <doc_id>` | 查询单条文件告警详情。 |

## Asset

| 命令 | 用途 |
| --- | --- |
| `chaitin-cli tanswer asset list` | 查询资产配置列表。 |
| `chaitin-cli tanswer asset detail --id <asset_id>` | 查询单个资产详情。 |
| `chaitin-cli tanswer asset group-tree` | 查询资产组树和分组数量。 |
| `chaitin-cli tanswer asset download-template` | 下载资产导入模板。 |
| `chaitin-cli tanswer asset export` | 导出资产配置。 |
| `chaitin-cli tanswer asset create` | preview/confirm 创建资产。 |
| `chaitin-cli tanswer asset update --id <asset_id>` | preview/confirm 更新资产。 |
| `chaitin-cli tanswer asset delete --id-list <asset_ids>` | preview/confirm 删除资产。 |
| `chaitin-cli tanswer asset batch-maintain --id-list <asset_ids>` | preview/confirm 批量维护资产字段。 |
| `chaitin-cli tanswer asset batch-tag --id-list <asset_ids> --tag-id <tag_ids>` | preview/confirm 批量维护资产标签。 |
| `chaitin-cli tanswer asset group-create --name <name>` | preview/confirm 创建资产组。 |
| `chaitin-cli tanswer asset group-rename --id <group_id> --name <name>` | preview/confirm 重命名资产组。 |
| `chaitin-cli tanswer asset group-delete --id-list <group_ids>` | preview/confirm 删除资产组。 |
| `chaitin-cli tanswer asset tree-move --id <id> --type <group|asset> --prev-id <id> --prev-type <group|asset>` | preview/confirm 调整资产树层级。 |
| `chaitin-cli tanswer asset import --file <xlsx>` | preview/confirm 上传资产导入文件。 |

## Metadata

| 命令 | 用途 |
| --- | --- |
| `chaitin-cli tanswer metadata protocol --protocol <protocol>` | 按协议查询流量元数据。 |
| `chaitin-cli tanswer metadata search --protocol <protocol> --advanced-query <query>` | 使用高级条件查询流量元数据。 |
| `chaitin-cli tanswer metadata detail --id <metadata_id> --timestamp <ms> --protocol <protocol>` | 查询单条元数据详情。 |
| `chaitin-cli tanswer metadata near-alarm --id <alarm_doc_id>` | 查询告警附近的流量上下文。 |
| `chaitin-cli tanswer metadata config` | 查询协议采集配置。 |
| `chaitin-cli tanswer metadata config-update` | preview/confirm 更新协议采集配置。 |

## Policy

| 命令 | 用途 |
| --- | --- |
| `chaitin-cli tanswer policy detection-whitelist` | 查询检测白名单。 |
| `chaitin-cli tanswer policy detection-whitelist-create` | preview/confirm 创建检测白名单。 |
| `chaitin-cli tanswer policy detection-whitelist-update` | preview/confirm 更新检测白名单。 |
| `chaitin-cli tanswer policy detection-whitelist-enable` | preview/confirm 启用检测白名单。 |
| `chaitin-cli tanswer policy detection-whitelist-disable` | preview/confirm 禁用检测白名单。 |
| `chaitin-cli tanswer policy detection-whitelist-delete` | preview/confirm 删除检测白名单。 |
| `chaitin-cli tanswer policy detection-whitelist-from-alarm` | preview/confirm 从误报告警生成检测白名单。 |
| `chaitin-cli tanswer policy detection-whitelist-export` | 导出检测白名单。 |
| `chaitin-cli tanswer policy detection-whitelist-import` | preview/confirm 上传检测白名单导入文件。 |
| `chaitin-cli tanswer policy custom-intelligence` | 查询自定义 IOC 情报。 |
| `chaitin-cli tanswer policy custom-intelligence-create` | preview/confirm 创建自定义 IOC 情报。 |
| `chaitin-cli tanswer policy custom-intelligence-update` | preview/confirm 更新自定义 IOC 情报。 |
| `chaitin-cli tanswer policy custom-intelligence-enable` | preview/confirm 启用自定义 IOC 情报。 |
| `chaitin-cli tanswer policy custom-intelligence-disable` | preview/confirm 禁用自定义 IOC 情报。 |
| `chaitin-cli tanswer policy custom-intelligence-delete` | preview/confirm 删除自定义 IOC 情报。 |
| `chaitin-cli tanswer policy custom-intelligence-export` | 导出自定义 IOC 情报。 |
| `chaitin-cli tanswer policy custom-intelligence-import` | preview/confirm 上传自定义 IOC 情报导入文件。 |

## Response

| 命令 | 用途 |
| --- | --- |
| `chaitin-cli tanswer response block-policies` | 查询阻断策略。 |
| `chaitin-cli tanswer response block-policy-create` | preview/confirm 创建阻断策略。 |
| `chaitin-cli tanswer response block-policy-update` | preview/confirm 更新阻断策略。 |
| `chaitin-cli tanswer response block-policy-enable` | preview/confirm 启用阻断策略。 |
| `chaitin-cli tanswer response block-policy-disable` | preview/confirm 禁用阻断策略。 |
| `chaitin-cli tanswer response block-policy-delete` | preview/confirm 删除阻断策略。 |
| `chaitin-cli tanswer response block-records` | 查询阻断命中记录。 |
| `chaitin-cli tanswer response whitelist` | 查询响应白名单。 |
| `chaitin-cli tanswer response whitelist-create` | preview/confirm 创建响应白名单。 |
| `chaitin-cli tanswer response whitelist-update` | preview/confirm 更新响应白名单。 |
| `chaitin-cli tanswer response whitelist-enable` | preview/confirm 启用响应白名单。 |
| `chaitin-cli tanswer response whitelist-disable` | preview/confirm 禁用响应白名单。 |
| `chaitin-cli tanswer response whitelist-delete` | preview/confirm 删除响应白名单。 |
| `chaitin-cli tanswer response block-policy-from-alarm` | preview/confirm 从告警生成阻断策略。 |
| `chaitin-cli tanswer response whitelist-from-alarm` | preview/confirm 从告警生成响应白名单。 |
| `chaitin-cli tanswer response devices` | 查询联动设备配置。 |
| `chaitin-cli tanswer response device-records --device-id <device_id>` | 查询联动设备发送记录。 |
| `chaitin-cli tanswer response auto-policies` | 查询自动响应策略。 |
| `chaitin-cli tanswer response auto-list` | 查询自动响应生成列表。 |

## 长示例

```bash
chaitin-cli tanswer auth status
chaitin-cli tanswer auth check
chaitin-cli tanswer manifest
chaitin-cli tanswer system status
chaitin-cli tanswer alarm overview --time today
chaitin-cli tanswer alarm overview --time 24h --severity critical,high --result success,control
chaitin-cli tanswer alarm timeline --time 24h --interval 1h --severity critical,high
chaitin-cli tanswer alarm list --time today --page-size 10
chaitin-cli tanswer alarm high-priority --time today
chaitin-cli tanswer alarm detail --id <doc_id>
chaitin-cli tanswer alarm by-attacker --attacker 192.0.2.10 --time today
chaitin-cli tanswer alarm by-victim --victim 192.0.2.20 --time today
chaitin-cli tanswer alarm by-threat --name SQL注入 --time today
chaitin-cli tanswer alarm important-assets --time today
chaitin-cli tanswer alarm attacker-rank --time today --top 10
chaitin-cli tanswer alarm victim-rank --time today --top 10
chaitin-cli tanswer alarm phase-distribution --time today
chaitin-cli tanswer alarm related --id <doc_id>
```

```bash
chaitin-cli tanswer file-alarm overview --time today
chaitin-cli tanswer file-alarm malicious --time today --page-size 10
chaitin-cli tanswer file-alarm webshell --time today --page-size 10
chaitin-cli tanswer file-alarm sandbox --time today --page-size 10
chaitin-cli tanswer file-alarm detail --id <doc_id>
chaitin-cli tanswer metadata protocol --protocol http --time today --page-size 10
chaitin-cli tanswer metadata search --protocol dns --advanced-query "dns_rrname = 'example.com'"
chaitin-cli tanswer metadata detail --id <metadata_id> --timestamp 1784282400000 --protocol http
chaitin-cli tanswer metadata near-alarm --id <doc_id> --window 30m --page-size 10
chaitin-cli tanswer metadata config
chaitin-cli tanswer metadata config-update --node-id <node_id> --enable http,dns --preview
chaitin-cli tanswer metadata config-update --node-id <node_id> --disable tcp,udp --confirm CONFIRM_METADATA_CONFIG_UPDATE
```

```bash
chaitin-cli tanswer asset list --page-size 10
chaitin-cli tanswer asset detail --id <asset_id>
chaitin-cli tanswer asset group-tree --depth 2
chaitin-cli tanswer asset download-template --output ./asset-template.xlsx
chaitin-cli tanswer asset import --file ./assets.xlsx --preview
chaitin-cli tanswer asset import --file ./assets.xlsx --confirm CONFIRM_ASSET_IMPORT
chaitin-cli tanswer asset export --id-list 3,7 --output ./selected-assets.xlsx
chaitin-cli tanswer asset create --name core-db --ip 192.0.2.10 --preview
chaitin-cli tanswer asset create --name core-db --ip 192.0.2.10 --confirm CONFIRM_ASSET_CREATE
chaitin-cli tanswer asset update --id 9 --name core-db-new --ip 192.0.2.11 --preview
chaitin-cli tanswer asset update --id 9 --name core-db-new --ip 192.0.2.11 --confirm CONFIRM_ASSET_UPDATE
chaitin-cli tanswer asset delete --id-list 9 --preview
chaitin-cli tanswer asset delete --id-list 9,10 --confirm CONFIRM_ASSET_DELETE
chaitin-cli tanswer asset batch-maintain --id-list 9,10 --contact secops --preview
chaitin-cli tanswer asset batch-tag --id-list 9,10 --tag-id 3,7 --confirm CONFIRM_ASSET_BATCH_TAG
chaitin-cli tanswer asset group-create --name 核心区 --parent-id 2 --preview
chaitin-cli tanswer asset group-rename --id 3 --name 核心区 --confirm CONFIRM_ASSET_GROUP_RENAME
chaitin-cli tanswer asset group-delete --id-list 3 --confirm CONFIRM_ASSET_GROUP_DELETE
chaitin-cli tanswer asset tree-move --id 9 --type asset --prev-id 3 --prev-type group --top-layer --preview
```

```bash
chaitin-cli tanswer policy detection-whitelist --page-size 10
chaitin-cli tanswer policy detection-whitelist-create --name 登录误报 --src-ip 192.0.2.10 --preview
chaitin-cli tanswer policy detection-whitelist-create --name 登录误报 --src-ip 192.0.2.10 --confirm CONFIRM_POLICY_DETECTION_WHITELIST_CREATE
chaitin-cli tanswer policy detection-whitelist-from-alarm --id <doc_id> --remark 已确认误报 --confirm CONFIRM_POLICY_DETECTION_WHITELIST_FROM_ALARM
chaitin-cli tanswer policy detection-whitelist-export --id-list 21,22 --output ./detection-whitelist.xlsx
chaitin-cli tanswer policy detection-whitelist-import --file ./detection-whitelist.xlsx --preview
chaitin-cli tanswer policy custom-intelligence --ioc evil.example --type domain --status enabled
chaitin-cli tanswer policy custom-intelligence-create --name 恶意域名 --ioc evil.example --type domain --confirm CONFIRM_POLICY_CUSTOM_INTELLIGENCE_CREATE
chaitin-cli tanswer policy custom-intelligence-export --id-list 12,13 --output ./custom-intelligence.xlsx
```

```bash
chaitin-cli tanswer response block-policies --page-size 10
chaitin-cli tanswer response block-policy-create --name block-bad-ip --object 192.0.2.30 --preview
chaitin-cli tanswer response block-policy-create --name block-bad-ip --object 192.0.2.30 --duration 3600 --confirm CONFIRM_RESPONSE_BLOCK_POLICY_CREATE
chaitin-cli tanswer response block-records --time 24h --page-size 10
chaitin-cli tanswer response whitelist --object 192.0.2.40 --type ip
chaitin-cli tanswer response whitelist-create --type ip --object 192.0.2.40 --expire 1784277612410 --preview
chaitin-cli tanswer response whitelist-create --type ip --object 192.0.2.40 --expire 1784277612410 --confirm CONFIRM_RESPONSE_WHITELIST_CREATE
chaitin-cli tanswer response block-policy-from-alarm --id <doc_id> --target attacker --preview
chaitin-cli tanswer response whitelist-from-alarm --id <doc_id> --target victim --expire 1784277612410 --confirm CONFIRM_RESPONSE_WHITELIST_FROM_ALARM
chaitin-cli tanswer response devices --page-size 10
chaitin-cli tanswer response device-records --device-id <device_id> --page-size 10
chaitin-cli tanswer response auto-policies --page-size 10
chaitin-cli tanswer response auto-list --time 7d --page-size 10
```

## Agent 使用规则

- 优先使用语义命令；只有没有语义覆盖时才使用 `tanswer api`。
- 写操作必须先给用户展示 preview，再用精确 `confirmation_token` 执行。
- `alarm overview` 在聚合 API 不可用时可 fallback 到列表 API，并通过 `summary.source=list_fallback` 标记。
- `file-alarm` 命令只读，不下载样本，不触发新沙箱分析。
- `asset tree-move` 只有在节点 ID、节点类型和产品 `top_layer` 语义明确时才使用。
- `metadata near-alarm` 只用于上下文调查，不单独构成攻击结论。
- `response block-policy-from-alarm` 和 `response whitelist-from-alarm` 只在已有 `doc_id` 且用户明确选择处置动作后使用。
- 使用 `chaitin-cli tanswer manifest` 查看当前版本支持的命令、参数、输出字段和确认要求。
