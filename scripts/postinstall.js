#!/usr/bin/env node

const crypto = require("crypto");
const fs = require("fs");
const os = require("os");
const path = require("path");
const AdmZip = require("adm-zip");
const tar = require("tar");

const rootDir = path.resolve(__dirname, "..");
const pkg = require(path.join(rootDir, "package.json"));

const vendorDir = path.join(rootDir, "vendor");
const tmpDir = fs.mkdtempSync(path.join(os.tmpdir(), "vtrix-npm-"));

const SUPPORTED_TARGETS = {
  "darwin-arm64": { os: "darwin", arch: "arm64", ext: "tar.gz", bin: "vtrix" },
  "darwin-x64": { os: "darwin", arch: "amd64", ext: "tar.gz", bin: "vtrix" },
  "linux-arm64": { os: "linux", arch: "arm64", ext: "tar.gz", bin: "vtrix" },
  "linux-x64": { os: "linux", arch: "amd64", ext: "tar.gz", bin: "vtrix" },
  "win32-x64": { os: "windows", arch: "amd64", ext: "zip", bin: "vtrix.exe" }
};

async function main() {
  try {
    if (process.env.VTRIX_SKIP_POSTINSTALL === "1") {
      log("skip postinstall because VTRIX_SKIP_POSTINSTALL=1");
      return;
    }

    const target = resolveTarget();
    const releaseBaseUrl = resolveReleaseBaseUrl();
    const projectName = pkg.vtrix?.projectName || "vtrix";
    const version = pkg.version;

    const assetName = `${projectName}_${version}_${target.os}_${target.arch}.${target.ext}`;
    const checksumUrl = `${releaseBaseUrl}/SHA256SUMS`;
    const assetUrl = `${releaseBaseUrl}/${assetName}`;
    const archivePath = path.join(tmpDir, assetName);
    const extractDir = path.join(tmpDir, "extract");

    fs.mkdirSync(vendorDir, { recursive: true });
    fs.rmSync(vendorDir, { recursive: true, force: true });
    fs.mkdirSync(vendorDir, { recursive: true });
    fs.mkdirSync(extractDir, { recursive: true });

    log(`downloading ${assetName}`);
    const checksumText = await fetchText(checksumUrl);
    const expectedSha = parseChecksum(checksumText, assetName);
    if (!expectedSha) {
      throw new Error(`checksum for ${assetName} not found in SHA256SUMS`);
    }

    await downloadFile(assetUrl, archivePath);
    const actualSha = sha256File(archivePath);
    if (actualSha !== expectedSha) {
      throw new Error(`checksum mismatch for ${assetName}`);
    }

    await extractArchive(archivePath, extractDir, target.ext);
    const extractedBinary = findFileRecursive(extractDir, target.bin);
    if (!extractedBinary) {
      throw new Error(`failed to locate ${target.bin} in extracted archive`);
    }

    const finalBinary = path.join(vendorDir, target.bin);
    fs.copyFileSync(extractedBinary, finalBinary);
    if (process.platform !== "win32") {
      fs.chmodSync(finalBinary, 0o755);
    }

    log(`installed ${target.bin} to vendor directory`);
  } finally {
    fs.rmSync(tmpDir, { recursive: true, force: true });
  }
}

function resolveTarget() {
  const key = `${process.platform}-${process.arch}`;
  const target = SUPPORTED_TARGETS[key];
  if (!target) {
    throw new Error(`unsupported platform: ${process.platform}/${process.arch}`);
  }
  return target;
}

function resolveReleaseBaseUrl() {
  const fromEnv = process.env.VTRIX_RELEASE_BASE_URL;
  if (fromEnv) {
    return stripTrailingSlash(fromEnv);
  }

  const template = pkg.vtrix?.releaseBaseUrlTemplate;
  if (!template) {
    throw new Error("missing vtrix.releaseBaseUrlTemplate in package.json");
  }
  return stripTrailingSlash(
    template
      .replaceAll("{version}", pkg.version)
      .replaceAll("{tag}", `v${pkg.version}`)
  );
}

function stripTrailingSlash(value) {
  return value.replace(/\/+$/, "");
}

async function fetchText(url) {
  const response = await fetch(url);
  if (!response.ok) {
    throw new Error(`failed to download ${url}: HTTP ${response.status}`);
  }
  return response.text();
}

async function downloadFile(url, destination) {
  const response = await fetch(url);
  if (!response.ok) {
    throw new Error(`failed to download ${url}: HTTP ${response.status}`);
  }
  const buffer = Buffer.from(await response.arrayBuffer());
  fs.writeFileSync(destination, buffer);
}

function parseChecksum(content, fileName) {
  for (const line of content.split(/\r?\n/)) {
    const trimmed = line.trim();
    if (!trimmed) {
      continue;
    }
    const match = trimmed.match(/^([a-fA-F0-9]{64})\s+\*?(.+)$/);
    if (match && path.basename(match[2]) === fileName) {
      return match[1].toLowerCase();
    }
  }
  return null;
}

function sha256File(filePath) {
  const hash = crypto.createHash("sha256");
  hash.update(fs.readFileSync(filePath));
  return hash.digest("hex");
}

async function extractArchive(archivePath, destination, extension) {
  if (extension === "zip") {
    const zip = new AdmZip(archivePath);
    zip.extractAllTo(destination, true);
    return;
  }
  if (extension === "tar.gz") {
    await tar.x({
      file: archivePath,
      cwd: destination
    });
    return;
  }
  throw new Error(`unsupported archive format: ${extension}`);
}

function findFileRecursive(dir, fileName) {
  const entries = fs.readdirSync(dir, { withFileTypes: true });
  for (const entry of entries) {
    const fullPath = path.join(dir, entry.name);
    if (entry.isFile() && entry.name === fileName) {
      return fullPath;
    }
    if (entry.isDirectory()) {
      const found = findFileRecursive(fullPath, fileName);
      if (found) {
        return found;
      }
    }
  }
  return null;
}

function log(message) {
  console.log(`[vtrix installer] ${message}`);
}

main().catch((err) => {
  console.error(`[vtrix installer] ${err.message}`);
  process.exit(1);
});
