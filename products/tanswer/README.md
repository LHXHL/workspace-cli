# 全悉 AI 可读 CLI

`chaitin-cli tanswer` 是全悉面向人类操作者和 AI Agent 的命令行入口。它将高频安全运营动作封装为语义命令，并保留 `api <METHOD> <PATH>` 作为 Open API fallback。

## 快速开始

```bash
chaitin-cli tanswer --url https://quanxi.example.com --api-key "$TANSWER_API_KEY" auth check
chaitin-cli tanswer manifest
chaitin-cli tanswer system +status
chaitin-cli tanswer alarm +overview --time today
```

## 安全边界

- 请通过环境变量、配置文件或命令行参数传入访问凭证，不要在脚本、日志或共享文档中暴露 API key。
- 查询命令直接执行；会修改产品配置或处置对象的命令必须先预览变更，再使用确认令牌执行。
- 使用 `chaitin-cli tanswer manifest` 查看当前版本支持的命令、参数、输出字段和确认要求。
