#!/usr/bin/env bash
set -euo pipefail

echo "Checking docs/examples build..."
go build ./docs/examples/...

echo "Checking docs/01-getting-started.md complete snippet..."
tmpdir=$(mktemp -d)
trap 'rm -rf "$tmpdir"' EXIT

python3 - <<'PY' "$tmpdir"
import pathlib, re, sys
out_dir = pathlib.Path(sys.argv[1])
text = pathlib.Path("docs/01-getting-started.md").read_text(encoding="utf-8")
blocks = re.findall(r"```go\n(.*?)```", text, re.DOTALL)
compiled = 0
for i, b in enumerate(blocks):
    if "package main" in b:
        p = out_dir / f"doc01_{i}.go"
        p.write_text(b, encoding="utf-8")
        compiled += 1
if compiled == 0:
    raise SystemExit("No complete package main Go block found in docs/01-getting-started.md")
print(f"Extracted {compiled} complete Go snippet(s) from docs/01-getting-started.md")
PY

for f in "$tmpdir"/*.go; do
  go build -o /dev/null "$f"
done

echo "Docs snippets OK."
