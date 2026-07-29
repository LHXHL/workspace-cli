# T-Answer Final Acceptance Matrix

This is the frozen acceptance baseline for the T-Answer migration. It defines the release gate for the current scope; an evaluator must not report the canonical product as passed unless every requirement below has fresh evidence.

The authoritative runtime interfaces are the installed binary's `--help` and `chaitin-cli tanswer manifest`. This matrix is for maintainers and independent evaluators; it is not an additional runtime dependency for an installed user.

## Candidate and evidence rules

1. Create a candidate commit before final acceptance. The independent evaluator must record the repository, branch, commit SHA, and a clean worktree.
2. Run all commands without a real endpoint, credential, or customer data. Use `127.0.0.1:9` and `placeholder` only where an attempted connection is required as evidence.
3. An evaluator must report one result and one concrete command/test evidence for every requirement ID. An untested requirement is a failure, not a pass by inference.
4. Jupiter's retired prototype is outside this product gate. Its deletion is separately gated on publishing a canonical CLI release.

## Requirement matrix

| ID | Requirement | Automated evidence | Independent evidence |
| --- | --- | --- | --- |
| R-01 Configuration contract | The URL, API Key, timeout, and TLS interfaces are discoverable and consistent: `--url`/`TANSWER_URL`/`tanswer.url`; `--api-key`/`TANSWER_API_KEY`/`tanswer.api_key`; `--timeout`/`TANSWER_TIMEOUT`/`tanswer.timeout`; `--insecure`/`TANSWER_INSECURE`/`tanswer.insecure`. Configuration lookup is `flags > environment/.env > recognized ./config.yaml > ~/.chaitin-cli/config.yaml`. | `TestTAnswerRuntimeAndRepositoryGuidanceListAllConfigurationSources`; `TestCompanySkillMatchesRuntimeConfigDiscovery`; `TestTAnswerGlobalConfigFallbackAtCLIEntry`; configuration unit tests. | Run `chaitin-cli tanswer --help` and `TANSWER_URL=http://127.0.0.1:9 TANSWER_API_KEY=placeholder TANSWER_TIMEOUT=7s TANSWER_INSECURE=true chaitin-cli --config /dev/null tanswer auth status`. |
| R-02 Runtime discovery | A person or AI can find T-Answer from root help, then find authentication, semantic domains, raw API guidance, manifest, and the fact that root `--dry-run` does not apply. | `TestTAnswerHelpGuidesUsersWithoutProductDocuments`; `TestTAnswerRootHelpDirectsUsersToProductHelp`; `TestManifestHelpExplainsPurpose`. | Run root, product, domain, leaf, and raw API help without reading product documents. |
| R-03 Command inventory | Runtime manifest commands are registered in Cobra help and the human command reference covers every manifest command using command-boundary matching. | `TestManifestCommandsAreDiscoverableFromHelp`; `TestHumanCommandReferenceMatchesRuntimeManifest`. | Count manifest commands and report any command-reference misses. |
| R-04 Semantic protected writes | Every semantic command with `requires_confirmation: true` has a complete preview contract: `requires_confirmation`, `confirmed`, `operation_type`, `risk_level`, `target`, `change_summary`, `impact`, `risk_warnings`, `confirmation_token`, and `confirmation_note`. It must state a confirmation condition and mention preview/confirm in help. | `TestSemanticProtectedWritesDeclareCompletePreviewContract`; protected-write unit tests. | Enumerate every protected semantic manifest command, verify its fields and help, then run an `asset create` preview and wrong-confirm case against the safe local fixture. |
| R-05 Raw API safety | GET/HEAD may request directly. Every other HTTP method defaults to preview, rejects wrong confirmation structurally without sending a request, and sends only after `CONFIRM_TANSWER_RAW_API_WRITE` following explicit user confirmation. The manifest must declare `--preview`, `--confirm`, its exact token, condition, and preview fields. | `TestRawAPIPotentiallyMutatingRequestReturnsPreviewWithoutSending`; `TestRawAPIPotentiallyMutatingRequestRejectsWrongConfirmStructurally`; `TestRawAPIPotentiallyMutatingRequestRequiresExactConfirm`; `TestManifestCommandOutputsAIReadableCommandMetadata`. | Run GET, HEAD, POST without confirm, POST with a wrong confirm, and POST with the exact token against `127.0.0.1:9`; record whether a connection is attempted. |
| R-06 Human and AI task journey | With only the root README and installed binary, a person and an AI can discover T-Answer, configure it, run auth status/check, select a read command, and explain the protected-write and raw API workflow without guessing. | Help and manifest discovery tests listed in R-01 to R-05. | Perform the journey without reading `products/tanswer/README.md`, `COMMAND_REFERENCE.md`, company skill, or Jupiter documents. Record the help/manifest path used for each step. |
| R-07 Documentation and secret hygiene | Root README (Chinese and English), company skill, product README, command reference, help, and manifest describe the same configuration and safety model. They contain no real credentials, customer addresses, or customer data. | `TestTAnswerRuntimeAndRepositoryGuidanceListAllConfigurationSources`; `TestEnglishRepositoryGuideIncludesTAnswerOnboarding`. | Run a scoped secret scan over the listed guidance files and manually compare each configuration tuple with R-01. |
| R-08 Negative controls | The most important acceptance tests must prove that they fail when their protected contract regresses. | Execute the disposable-copy procedure below before release sign-off. | Report each intentional mutation, its selected test, and the observed non-zero exit status. |

## Required negative controls

Run these only in a disposable copy of the candidate commit. They must never modify the candidate worktree.

```sh
candidate_dir="$(mktemp -d)"
git clone --no-local . "$candidate_dir/chaitin-cli"
cd "$candidate_dir/chaitin-cli"
```

For each control, restore the disposable copy from the candidate commit before applying the next mutation.

| Mutation in the disposable copy | Required failing test |
| --- | --- |
| Replace every `TANSWER_TIMEOUT` with `TANSWER_TIMEOUT_REGRESSED` in `products/tanswer/root.go`. | `go test ./products/tanswer -run TestTAnswerRuntimeAndRepositoryGuidanceListAllConfigurationSources` |
| Replace `RAW_API_CONFIRMATION_REQUIRED` with `RAW_API_CONFIRMATION_REGRESSED` in `products/tanswer/raw_api.go`. | `go test ./products/tanswer -run TestRawAPIPotentiallyMutatingRequestRejectsWrongConfirmStructurally` |
| Rename raw API manifest flag `--confirm` to `--confirm-regressed` in `products/tanswer/manifest.go`. | `go test ./products/tanswer -run TestManifestCommandOutputsAIReadableCommandMetadata` |
| Replace the semantic preview field `risk_warnings` in `semanticProtectedWritePreviewFields` with `risk_warnings_regressed`. | `go test ./products/tanswer -run TestSemanticProtectedWritesDeclareCompletePreviewContract` |

Every command in the right column must fail because of the intended mutation. A test that remains green is a release blocker.

## Final independent-AI protocol

1. Confirm the candidate commit and clean worktree.
2. Run `GOCACHE=/private/tmp/chaitin-cli-final-gocache go test -count=1 ./...`, build the binary, and run `git diff --check`.
3. Evaluate R-01 through R-08 in order. The final report must include all eight IDs, commands, exit codes, and observed evidence.
4. Mark canonical acceptance as passed only if all eight IDs pass. Do not combine a Jupiter retirement observation with the canonical result.
