# 全悉 AI 可读 CLI

`chaitin-cli tanswer` 是全悉面向人类操作者和 AI Agent 的命令行入口。它以语义命令覆盖一期安全运营场景，并保留 `api <METHOD> <PATH>` 作为 Open API fallback。

## 快速开始

```bash
chaitin-cli tanswer --url https://quanxi.example.com --api-key "$TANSWER_API_KEY" auth check
chaitin-cli tanswer manifest
chaitin-cli tanswer system +status
chaitin-cli tanswer alarm +overview --time today
```

## 安全边界

- 只提交脱敏示例，不提交真实 token、真实测试地址或真实内网地址。
- 查询命令直接执行；写命令必须 preview/confirm/audit。
- `risk`、`risk-host`、`asset-risk`、`vulnerability-risk` 不在一期范围。
