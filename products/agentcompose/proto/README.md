# Agent Compose protobuf snapshot

These files are generated snapshots of the Agent Compose V2 API. The source of
truth is the `agent-compose` repository at commit
`99bd5586d03d4e5c95b5143f779196d9f0dd3fca` (feat/cli-logs-event-id-filter; pending v2609.3.1+).

Run `task agent-compose:proto:sync` from the repository root to refresh and
verify the snapshot, then run `task agent-compose:proto:generate` to regenerate
the Go messages and Connect clients. Do not edit snapshot or generated files by
hand.
