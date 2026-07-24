# 全悉 CLI（`tanswer`）

`chaitin-cli tanswer` 是全悉面向安全运营人员和 AI Agent 的命令行入口，覆盖系统确认、威胁与文件告警、资产、流量元数据、安全策略和响应处置。

## 文档定位

本文档是全悉 CLI 的人工可读入门，供产品用户、开发者和仓库维护者快速了解产品边界、常见操作与本地验证方式。

它不替代运行时命令契约。安装后的二进制才是当前版本的唯一事实来源：

```bash
chaitin-cli tanswer --help
chaitin-cli tanswer manifest
```

AI Agent 必须先通过 `--help` 发现命令；需要完整的机器可读参数、输出、风险和确认契约时读取 `manifest`，不能以本文档猜测命令或 API。

## 快速使用

连接信息可通过 flag、环境变量或 `config.yaml` 的 `tanswer` 节配置；优先级为 flag、环境变量、配置文件、默认值。

```bash
export TANSWER_URL='https://<全悉控制台地址>'
export TANSWER_API_KEY='<OpenAPI Token>'

chaitin-cli tanswer auth check
chaitin-cli tanswer system status
chaitin-cli tanswer alarm overview --time 24h
```

在证书校验必须绕过的受控测试环境，可额外设置 `TANSWER_INSECURE=true` 或传入 `--insecure`。不要把真实 Token 写入仓库、脚本、日志或共享文档。

## 使用规则

1. 优先使用按业务领域组织的语义命令。
2. 先通过根命令或领域命令的 `--help` 了解可用操作和必填参数。
3. 对受保护写操作，先执行 `--preview`，确认影响后再使用该命令 help 或 manifest 指定的 `--confirm` token。
4. 仅在目标能力没有语义命令覆盖，且调用者已知并获授权访问对应端点时，才使用 `tanswer api <METHOD> <PATH>`。

常见领域入口：

```bash
chaitin-cli tanswer auth --help
chaitin-cli tanswer system --help
chaitin-cli tanswer alarm --help
chaitin-cli tanswer file-alarm --help
chaitin-cli tanswer asset --help
chaitin-cli tanswer metadata --help
chaitin-cli tanswer policy --help
chaitin-cli tanswer response --help
```

可供人工快速浏览的领域命令索引见 [COMMAND_REFERENCE.md](./COMMAND_REFERENCE.md)。具体 flag、枚举、输出字段及写操作确认要求仍以运行时 help 和 manifest 为准。

## 开发验证

在仓库根目录执行：

```bash
go test ./products/tanswer
go build -buildvcs=false -o /tmp/chaitin-cli .
/tmp/chaitin-cli tanswer manifest
```

产品测试覆盖运行时引导、manifest 契约、help 可发现性、配置、结构化输出和受保护写操作。
