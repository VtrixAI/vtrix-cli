#!/usr/bin/env bash
set -euo pipefail

GO="${GO:-$(which go 2>/dev/null || echo /opt/homebrew/bin/go)}"

APP="vtrix"
VERSION="${VERSION:-$(git describe --tags --always --dirty 2>/dev/null || echo "unknown")}"
DIST="dist"

# Production URLs — use online defaults, allow override via env when needed.
VTRIX_BASE_URL="${VTRIX_BASE_URL:-https://vtrix.ai}"
VTRIX_MODELS_URL="${VTRIX_MODELS_URL:-https://seacloud-cloud-model-spec.api.seaart.ai}"
VTRIX_GENERATION_URL="${VTRIX_GENERATION_URL:-$VTRIX_BASE_URL}"
VTRIX_SKILLHUB_URL="${VTRIX_SKILLHUB_URL:-https://seacloud-skill-hub.api.seaart.ai/api/v1}"

LDFLAGS="-s -w \
  -X github.com/VtrixAI/vtrix-cli/internal/buildinfo.Version=${VERSION} \
  -X github.com/VtrixAI/vtrix-cli/internal/auth.BaseURL=${VTRIX_BASE_URL} \
  -X github.com/VtrixAI/vtrix-cli/internal/models.BaseURL=${VTRIX_MODELS_URL} \
  -X github.com/VtrixAI/vtrix-cli/internal/generation.BaseURL=${VTRIX_GENERATION_URL} \
  -X github.com/VtrixAI/vtrix-cli/internal/skillhub.BaseURL=${VTRIX_SKILLHUB_URL}"

TARGETS=(
  "darwin/amd64"
  "darwin/arm64"
  "linux/amd64"
  "linux/arm64"
  "windows/amd64"
)

rm -rf "$DIST"
mkdir -p "$DIST"

echo "Building $APP $VERSION (prod)"
echo "  BaseURL:          $VTRIX_BASE_URL"
echo "  ModelsBaseURL:    $VTRIX_MODELS_URL"
echo "  GenerationBaseURL: $VTRIX_GENERATION_URL"
echo "  SkillhubBaseURL:  $VTRIX_SKILLHUB_URL"
echo ""

for target in "${TARGETS[@]}"; do
  OS="${target%/*}"
  ARCH="${target#*/}"

  BIN="$APP"
  [[ "$OS" == "windows" ]] && BIN="${APP}.exe"

  OUT_DIR="$DIST/${APP}_${OS}_${ARCH}"
  mkdir -p "$OUT_DIR"

  echo "  -> $OS/$ARCH"
  GOOS="$OS" GOARCH="$ARCH" CGO_ENABLED=0 "$GO" build \
    -ldflags="${LDFLAGS}" \
    -o "$OUT_DIR/$BIN" .

  if [[ "$OS" == "windows" ]]; then
    (cd "$DIST" && zip -q "${APP}_${OS}_${ARCH}.zip" "${APP}_${OS}_${ARCH}/${BIN}")
  else
    tar -czf "$DIST/${APP}_${OS}_${ARCH}.tar.gz" -C "$DIST" "${APP}_${OS}_${ARCH}"
  fi

  rm -rf "$OUT_DIR"
done

echo ""
echo "Artifacts in ./$DIST/:"
ls -lh "$DIST/"
