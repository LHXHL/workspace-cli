# Agent Compose protobuf snapshot

These files are generated snapshots of the Agent Compose V2 API. The source of
truth is the `agent-compose` repository at commit
`85ba5a4185e907c1306beb588d3800e570a610b0` (feat/cli-logs-event-id-filter; pending v2609.3.1+).

Run `task agent-compose:proto:sync` from the repository root to refresh and
verify the snapshot, then run `task agent-compose:proto:generate` to regenerate
the Go messages and Connect clients. Do not edit snapshot or generated files by
hand.
