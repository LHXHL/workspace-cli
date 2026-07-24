# 全悉 CLI 命令索引

这是供仓库维护者快速浏览的静态索引，不是安装包交付依赖，也不是命令事实来源。

在实际操作或让 AI Agent 操作前，必须使用当前安装版本的运行时契约：

```bash
chaitin-cli tanswer --help
chaitin-cli tanswer <领域或命令> --help
chaitin-cli tanswer manifest
```

尤其是参数、枚举、输出字段、风险等级和 `--confirm` token，只能以运行时 help 与 manifest 为准。

## Foundation

```text
tanswer auth status
tanswer auth check
tanswer system status
tanswer manifest
tanswer api <METHOD> <PATH>
```

`api` 仅用于已知且已授权、又没有语义命令覆盖的 Open API；不要据此猜测 RPC 方法或路径。

## Threat alarm

```text
tanswer alarm overview
tanswer alarm list
tanswer alarm timeline
tanswer alarm high-priority
tanswer alarm detail
tanswer alarm by-attacker
tanswer alarm by-victim
tanswer alarm by-threat
tanswer alarm important-assets
tanswer alarm attacker-rank
tanswer alarm victim-rank
tanswer alarm phase-distribution
tanswer alarm related
```

## File alarm

```text
tanswer file-alarm overview
tanswer file-alarm malicious
tanswer file-alarm webshell
tanswer file-alarm sandbox
tanswer file-alarm detail
```

## Asset

```text
tanswer asset list
tanswer asset detail
tanswer asset group-tree
tanswer asset download-template
tanswer asset export
tanswer asset create
tanswer asset update
tanswer asset delete
tanswer asset batch-maintain
tanswer asset batch-tag
tanswer asset group-create
tanswer asset group-rename
tanswer asset group-delete
tanswer asset tree-move
tanswer asset import
```

`create`、`update`、`delete`、批量操作、资产组操作和导入属于受保护写操作。

## Metadata

```text
tanswer metadata protocol
tanswer metadata search
tanswer metadata detail
tanswer metadata near-alarm
tanswer metadata config
tanswer metadata config-update
```

`config-update` 属于受保护写操作。

## Security policy

```text
tanswer policy detection-whitelist
tanswer policy detection-whitelist-create
tanswer policy detection-whitelist-update
tanswer policy detection-whitelist-enable
tanswer policy detection-whitelist-disable
tanswer policy detection-whitelist-delete
tanswer policy detection-whitelist-from-alarm
tanswer policy detection-whitelist-export
tanswer policy detection-whitelist-import

tanswer policy custom-intelligence
tanswer policy custom-intelligence-create
tanswer policy custom-intelligence-update
tanswer policy custom-intelligence-enable
tanswer policy custom-intelligence-disable
tanswer policy custom-intelligence-delete
tanswer policy custom-intelligence-export
tanswer policy custom-intelligence-import
```

除查询和导出外，策略命令通常会改变产品状态，必须先查看对应命令的 help 或 manifest。

## Response

```text
tanswer response block-policies
tanswer response block-policy-create
tanswer response block-policy-update
tanswer response block-policy-enable
tanswer response block-policy-disable
tanswer response block-policy-delete
tanswer response block-records
tanswer response whitelist
tanswer response whitelist-create
tanswer response whitelist-update
tanswer response whitelist-enable
tanswer response whitelist-disable
tanswer response whitelist-delete
tanswer response block-policy-from-alarm
tanswer response whitelist-from-alarm
tanswer response devices
tanswer response device-records
tanswer response auto-policies
tanswer response auto-list
```

创建、更新、启停、删除和从告警创建策略/白名单属于受保护写操作。

## 受保护写操作

统一流程如下：

```bash
# 1. 用实际命令及参数查看影响；该步骤不执行产品变更。
chaitin-cli tanswer <write-command> ... --preview

# 2. 读取 preview 后，使用该命令 help 或 manifest 给出的精确 token 执行。
chaitin-cli tanswer <write-command> ... --confirm CONFIRM_<OPERATION>
```

不要复用或猜测 `CONFIRM_<OPERATION>`；每个操作需要的精确 token 以运行时返回内容为准。
