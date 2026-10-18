#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
VERSION="${1:-}"
OUTPUT_DIR="${2:-${ROOT_DIR}/dist}"
PROVIDER_ADDRESS="registry.terraform.io/repinoid/nubes"
MANIFEST_SOURCE="${ROOT_DIR}/terraform-registry-manifest.json"

if [[ ! "$VERSION" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  echo "Usage: bash scripts/build-release.sh <X.Y.Z> [output-dir]" >&2
  exit 2
fi
if [[ ! -f "${ROOT_DIR}/go.mod" || ! -f "${ROOT_DIR}/main.go" ]]; then
  echo "Error: run this script from an exported provider repository" >&2
  exit 2
fi
if [[ -z "${GPG_KEY_FILE:-}" || ! -f "${GPG_KEY_FILE}" ]]; then
  echo "Error: GPG_KEY_FILE must point to the dedicated release signing key" >&2
  exit 2
fi
if [[ -z "${GPG_KEY_ID:-}" ]]; then
  echo "Error: GPG_KEY_ID must identify the public key registered in Terraform Registry" >&2
  exit 2
fi
if [[ -n "${GPG_PASSPHRASE_FILE:-}" && ! -f "${GPG_PASSPHRASE_FILE}" ]]; then
  echo "Error: GPG_PASSPHRASE_FILE does not exist" >&2
  exit 2
fi

mkdir -p "$(dirname "$OUTPUT_DIR")"
OUTPUT_DIR="$(cd "$(dirname "$OUTPUT_DIR")" && pwd)/$(basename "$OUTPUT_DIR")"
case "$OUTPUT_DIR" in
  "${ROOT_DIR}"/*) ;;
  *)
    echo "Error: output directory must stay inside the provider repository" >&2
    exit 2
    ;;
esac
if [[ -e "$OUTPUT_DIR" && ( ! -d "$OUTPUT_DIR" || -n "$(find "$OUTPUT_DIR" -mindepth 1 -print -quit)" ) ]]; then
  echo "Error: output directory must be new or empty: $OUTPUT_DIR" >&2
  exit 2
fi
if [[ ! -f "$MANIFEST_SOURCE" ]]; then
  echo "Error: Terraform Registry manifest missing: $MANIFEST_SOURCE" >&2
  exit 2
fi

BUILD_DIR="$(mktemp -d "${OUTPUT_DIR}.tmp.XXXXXX")"
GNUPGHOME="$(mktemp -d)"
chmod 700 "$GNUPGHOME"
cleanup() {
  rm -rf "$BUILD_DIR" "$GNUPGHOME"
}
trap cleanup EXIT

export GNUPGHOME
gpg --batch --import "$GPG_KEY_FILE" >/dev/null
if ! gpg --batch --list-secret-keys "$GPG_KEY_ID" >/dev/null 2>&1; then
  echo "Error: GPG_KEY_ID is not present in the imported key" >&2
  exit 2
fi

build_archive() {
  local os="$1"
  local arch="$2"
  local binary="terraform-provider-nubes_v${VERSION}"
  if [[ "$os" == "windows" ]]; then
    binary+=".exe"
  fi

  echo "Building ${os}/${arch}"
  (
    cd "$ROOT_DIR"
    CGO_ENABLED=0 GOOS="$os" GOARCH="$arch" go build \
      -trimpath \
      -ldflags="-s -w -X main.version=${VERSION} -X main.address=${PROVIDER_ADDRESS}" \
      -o "${BUILD_DIR}/${binary}" .
  )
  (
    cd "$BUILD_DIR"
    python3 -m zipfile -c \
      "terraform-provider-nubes_${VERSION}_${os}_${arch}.zip" \
      "$binary"
  )
  rm -f "${BUILD_DIR:?}/${binary}"
}

build_archive linux amd64
build_archive windows amd64
build_archive darwin amd64

cp "$MANIFEST_SOURCE" "$BUILD_DIR/terraform-registry-manifest.json"
(cd "$BUILD_DIR" && sha256sum ./*.zip | sed 's#  \./#  #' > "terraform-provider-nubes_${VERSION}_SHA256SUMS")

SIGN_ARGS=(--batch --yes --local-user "$GPG_KEY_ID")
if [[ -n "${GPG_PASSPHRASE_FILE:-}" ]]; then
  SIGN_ARGS+=(--pinentry-mode loopback --passphrase-file "$GPG_PASSPHRASE_FILE")
fi
gpg "${SIGN_ARGS[@]}" --detach-sign \
  --output "${BUILD_DIR}/terraform-provider-nubes_${VERSION}_SHA256SUMS.sig" \
  "${BUILD_DIR}/terraform-provider-nubes_${VERSION}_SHA256SUMS"
gpg "${SIGN_ARGS[@]}" --verify \
  "${BUILD_DIR}/terraform-provider-nubes_${VERSION}_SHA256SUMS.sig" \
  "${BUILD_DIR}/terraform-provider-nubes_${VERSION}_SHA256SUMS"

if [[ -d "$OUTPUT_DIR" ]]; then
  rmdir "$OUTPUT_DIR"
fi
mv "$BUILD_DIR" "$OUTPUT_DIR"
echo "Release assets prepared in $OUTPUT_DIR"
