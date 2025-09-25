#!/bin/sh

# go
(command -v go 1>/dev/null && printf '[✓] go\n') ||
  (printf 'error: "go" required but not found - please install; exiting' && exit 1)

# curl
(command -v curl 1>/dev/null && printf '[✓] curl\n') ||
  (printf '⚠ error: "curl" required but not found - please install; exiting' && exit 1)

# .tools directory
(command -v mkdir 1>/dev/null && mkdir -p .tools && printf '[✓] ".tools" directory\n') ||
  (printf '⚠ error: "mkdir" required but not found' && exit 1)

# golangci-lint
(
  printf '\n==> installing golangci-lint...\n' &&
    (curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b .tools latest) &&
    printf '[✓] golangci-lint\n'
)

# make
(command -v make 1>/dev/null && printf '[✓] make\n') ||
  (printf '⚠ error: "make" required but not found - please install; exiting' && exit 1)
