#!/usr/bin/env bash
# Jenkins host orchestration only: every toolchain runs in a disposable,
# labelled container. Source is copied into a per-build Docker volume (no
# workspace bind mount, no Git metadata, no credentials).
#   CI_RUN_ID=ci-<unique> bash ci/run.sh <command>
set -Eeuo pipefail
cd "$(dirname "$0")/.."
: "${CI_RUN_ID:?Set a unique CI_RUN_ID (ci- followed by lowercase letters, digits or hyphens)}"
[[ "$CI_RUN_ID" =~ ^ci-[a-z0-9][a-z0-9-]{0,58}$ ]] || { echo 'Invalid CI_RUN_ID' >&2; exit 2; }
label_key=com.redghoul.ci.run
label="$label_key=$CI_RUN_ID"
volume="${CI_RUN_ID}-work"
network="${CI_RUN_ID}-net"

hardening=(--cap-drop ALL --security-opt no-new-privileges)

# task <name> <image> <cpus> <memory> <script> [extra docker run args...]
# Steps that only need the already-populated volume pass --network none.
task() {
  local name=$1 image=$2 cpus=$3 memory=$4 script=$5; shift 5
  docker run --rm --name "${CI_RUN_ID}-${name}" --label "$label" --init "${hardening[@]}" \
    --cpus "$cpus" --memory "$memory" --memory-swap "$memory" --pids-limit 2048 \
    --mount "type=volume,src=$volume,dst=/work" -w /work -e CI=true "$@" \
    "$image" bash -Eeuo pipefail -c "$script"
}

init() {
  docker version --format 'docker client {{.Client.Version}} / server {{.Server.Version}}'
  # Fail on accidental ID reuse; never clean up a possibly active build first.
  if docker volume inspect "$volume" >/dev/null 2>&1; then
    echo "Build volume already exists: $volume" >&2; exit 1
  fi
  for image in "$@"; do docker pull -q "$image"; done
  docker volume create --label "$label" "$volume" >/dev/null
  git archive --format=tar HEAD | docker run -i --rm --label "$label" --network none "${hardening[@]}" \
    --mount "type=volume,src=$volume,dst=/work" "$1" tar -xf - -C /work
  git rev-parse HEAD
}

cleanup() {
  # Exact build label only; never a broad prune or another build's resources.
  local ids failed=0
  ids=$(docker ps -aq --filter "label=$label") || return 1
  if [[ -n "$ids" ]]; then
    # Docker emits whitespace-delimited IDs, not user-supplied arguments.
    docker rm -f $ids >/dev/null || failed=1
  fi
  if docker network inspect "$network" >/dev/null 2>&1; then
    [[ $(docker network inspect -f "{{index .Labels \"$label_key\"}}" "$network") == "$CI_RUN_ID" ]] || return 1
    docker network rm "$network" >/dev/null || failed=1
  fi
  if docker volume inspect "$volume" >/dev/null 2>&1; then
    [[ $(docker volume inspect -f "{{index .Labels \"$label_key\"}}" "$volume") == "$CI_RUN_ID" ]] || return 1
    docker volume rm "$volume" >/dev/null || failed=1
  fi
  return "$failed"
}
trap 'trap - INT TERM; cleanup || true; exit 143' INT TERM

GO_IMAGE=golang:1.24-bookworm
go_env=(-e GOTOOLCHAIN=local -e GOFLAGS=-mod=readonly -e GOMODCACHE=/work/.cache/gomod -e GOCACHE=/work/.cache/gobuild)

case "${1:-}" in
  init) init "$GO_IMAGE" ;;
  deps) task deps "$GO_IMAGE" 2 2g 'go version && go mod download' "${go_env[@]}" ;;
  vet) task vet "$GO_IMAGE" 2 3g 'go vet ./...' "${go_env[@]}" --network none ;;
  build) task build "$GO_IMAGE" 2 3g 'go build ./...' "${go_env[@]}" --network none ;;
  test) task test "$GO_IMAGE" 2 3g 'go test ./... -v' "${go_env[@]}" --network none ;;
  cleanup) cleanup ;;
  *) echo 'Usage: ci/run.sh init|deps|vet|build|test|cleanup' >&2; exit 2 ;;
esac
