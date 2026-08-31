# Product Goals (No-PRD)

**Project**: gzh-cli-core
**Doc Type**: Goals + Constraints + Quality Gates
**Status**: Active
**Last Updated**: 2026-08-31

______________________________________________________________________

## Product Intent

gzh-cli-core is the **dependency-free foundation library** that the whole gzh-cli
family builds on. It:

- provides shared utilities (logger, testutil, errors, config, cli, version),
- depends on no other `gzh-cli-*` repo — it is the root of the dependency graph,
- and stays import-only: it ships no CLI binary.

This is a library project — a single PRODUCT.md is sufficient. It replaces a PRD.

| 제공하는 것 (Is)                                        | 되지 않을 것 (Is Not)                     |
| ------------------------------------------------------- | ----------------------------------------- |
| 패밀리 공용 유틸 (logger·errors·config·cli·testutil·version) | 특정 도메인 기능 (git/quality/shell 등)   |
| 무의존 기반 — 다른 `gzh-cli-*` 의존 0건                 | feature 라이브러리를 역참조               |
| 최소 서드파티 (stdlib + yaml.v3 + cobra/pflag)          | 무거운 프레임워크·DI 컨테이너·CGO 의존    |
| import로 소비하는 라이브러리                            | 독립 실행 CLI 바이너리                    |

______________________________________________________________________

## Goals (Measurable Targets)

G1. **Zero family dependency**

- Target: `go.mod`의 require에 `gzh-cli-*` 0건, CGO 사용 0건

G2. **Thin dependency surface**

- Target: 직접 서드파티는 stdlib 외 yaml.v3 + cobra/pflag만; 신규 추가는 본 문서에 기록

G3. **Test reliability**

- Target: 전 패키지 테스트 통과, 커버리지 >= 85% (기반 라이브러리이므로 상향)

G4. **Cross-compile friendliness**

- Target: `CGO_ENABLED=0`으로 linux/darwin/windows × amd64/arm64 빌드 성공

G5. **API stability**

- Target: 공개 API의 파괴적 변경은 semver major 범프 + 마이그레이션 노트 동반

______________________________________________________________________

## Non-Goals (Explicitly Out of Scope)

- No 도메인 기능 (git/quality/pm/env 등은 feature 리포의 몫)
- No CLI 바이너리 (core는 라이브러리 전용)
- No 다른 `gzh-cli-*` 의존 — 역방향/상호 의존 전면 금지
- No 무거운 프레임워크나 DI 컨테이너 도입
- No CGO 의존 (크로스 컴파일성 훼손)

______________________________________________________________________

## Guardrails and Technical Constraints

**Architecture**

- 패키지당 단일 관심사 (one concern per package); 공개 함수는 테스트·GoDoc 필수

**Dependency Boundaries**

- **core는 다른 `gzh-cli-*`를 절대 의존하지 않는다 (무의존).** 이 규칙이 패밀리
  의존 그래프의 뿌리이며, 위반 시 릴리스 순서(core → feature → gzh-cli)가 깨진다
  (GUIDELINES §2)

**Compatibility**

- 소비자 최소 Go 1.26 (`go.mod`의 `go` directive); CGO 미사용. `toolchain go1.26.7`은
  이 모듈을 main module로 개발할 때의 권장 toolchain이며 CI도 1.26.7로 고정한다.
  소비자 하한 변경은 compatibility review 대상이고, toolchain 변경은 소비자 계약이 아니라
  개발·CI 재현성 관점에서 검토한다.

**Safety**

- 공개 API의 파괴적 변경은 semver major에서만; 하위 소비자 신뢰를 최우선한다

**Baseline (진행 중)**

- `Makefile`·`.golangci.yml`을 보유한다. `make check`은 fmt·lint·race test를 실행하며,
  build·vet·format-only 검증은 각각 `make build`·`make vet`·`make fmt-check`으로 수행한다.

______________________________________________________________________

## Quality Gates (Release Readiness)

**Build and Lint**

- `make check`(fmt·lint·race test), `make build`, `make vet`, `make fmt-check` 통과

**Testing**

- `make check`의 race test pass; 전체 커버리지 >= 85%

**Compatibility**

- `CGO_ENABLED=0` 크로스 빌드 (linux/darwin/windows × amd64/arm64) 성공

**Docs**

- 모든 공개 타입/함수 GoDoc 문서화

______________________________________________________________________

## Decision Rules

- 새 유틸은 **여러** `gzh-cli-*` 리포에서 공용으로 필요할 때만 추가한다
  (단일 리포 전용이면 그 리포에 둔다)
- `gzh-cli-*` 의존을 요구하는 코드는 core에 들어올 수 없다 (무의존 게이트)
- 서드파티 추가는 Guardrails에 정당화를 기록해야 한다
- Quality Gates 미충족 시 릴리스는 차단된다

______________________________________________________________________

**End of Document**
