#!/bin/sh

# make
(command -v make 1>/dev/null && printf '[✓] make\n') ||
  (printf '⚠ error: "make" required but not found - please install; exiting' && exit 1)

# pre-commit hook
(
  (test -e .git/hooks/pre-commit || test -L .git/hooks/pre-commit) &&
    printf '[✓] pre-commit hook\n'
) ||
  (
    printf '\n==> installing pre-commit hook...\n' &&
      cp .dev/pre-commit.sh .git/hooks/pre-commit &&
      printf '[✓] pre-commit hook\n'
  )

# pre-push hook
(
  (test -e .git/hooks/pre-push || test -L .git/hooks/pre-push) &&
    printf '[✓] pre-push hook\n'
) ||
  (
    printf '\n==> installing pre-push hook...\n' &&
      cp .dev/pre-push.sh .git/hooks/pre-push &&
      printf '[✓] pre-push hook\n'
  )
