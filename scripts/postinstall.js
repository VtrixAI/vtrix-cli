#!/usr/bin/env node

const crypto = require("crypto");
const fs = require("fs");
const os = require("os");
const path = require("path");
const AdmZip = require("adm-zip");
const tar = require("tar");

const rootDir = path.resolve(__dirname, "..");
const pkg = require(path.join(rootDir, "package.json"));

const bundledAssetsDir = path.join(
  rootDir,
  pkg.vtrix?.bundledAssetsDir || "npm-bundles"
);
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
    const releaseSources = resolveReleaseSources();
    const projectName = pkg.vtrix?.projectName || "vtrix";
    const version = pkg.version;

    const assetName = `${projectName}_${version}_${target.os}_${target.arch}.${target.ext}`;
    const archivePath = path.join(tmpDir, assetName);
    const extractDir = path.join(tmpDir, "extract");

    fs.mkdirSync(vendorDir, { recursive: true });
    fs.rmSync(vendorDir, { recursive: true, force: true });
    fs.mkdirSync(vendorDir, { recursive: true });
    fs.mkdirSync(extractDir, { recursive: true });

    log(`downloading ${assetName}`);
    const expectedSha = await resolveExpectedSha(releaseSources, assetName);
    await materializeArchive(releaseSources, assetName, archivePath, expectedSha);

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
  return resolveBaseUrl(
    process.env.VTRIX_RELEASE_BASE_URL,
    pkg.vtrix?.releaseBaseUrlTemplate,
    "missing vtrix.releaseBaseUrlTemplate in package.json"
  );
}

function resolveReleaseMirrorBaseUrl() {
  return resolveBaseUrl(
    process.env.VTRIX_RELEASE_MIRROR_BASE_URL,
    pkg.vtrix?.releaseMirrorBaseUrlTemplate || null,
    null
  );
}

function resolveBaseUrl(fromEnv, template, missingMessage) {
  if (fromEnv) {
    return stripTrailingSlash(fromEnv);
  }
  if (!template) {
    if (missingMessage) {
      throw new Error(missingMessage);
    }
    return null;
  }
  return stripTrailingSlash(
    template
      .replaceAll("{version}", pkg.version)
      .replaceAll("{tag}", `v${pkg.version}`)
  );
}

function resolveReleaseSources() {
  const sources = [
    { name: "GitHub Release", type: "remote", baseUrl: resolveReleaseBaseUrl() }
  ];
  const mirrorBaseUrl = resolveReleaseMirrorBaseUrl();

  if (mirrorBaseUrl && mirrorBaseUrl !== sources[0].baseUrl) {
    sources.push({ name: "Release mirror", type: "remote", baseUrl: mirrorBaseUrl });
  }

  if (fs.existsSync(bundledAssetsDir)) {
    sources.push({ name: "bundled npm package", type: "local", basePath: bundledAssetsDir });
  }

  return sources;
}

function getChecksumLocation(source) {
  return source.type === "local"
    ? path.join(source.basePath, "SHA256SUMS")
    : `${source.baseUrl}/SHA256SUMS`;
}

function getAssetLocation(source, assetName) {
  return source.type === "local"
    ? path.join(source.basePath, assetName)
    : `${source.baseUrl}/${assetName}`;
}

async function resolveExpectedSha(sources, assetName) {
  const errors = [];

  for (const source of sources) {
    const checksumLocation = getChecksumLocation(source);
    try {
      const checksumText = source.type === "local"
        ? fs.readFileSync(checksumLocation, "utf8")
        : await fetchText(checksumLocation);
      const expectedSha = parseChecksum(checksumText, assetName);
      if (!expectedSha) {
        throw new Error(`checksum for ${assetName} not found in SHA256SUMS`);
      }
      log(`using checksums from ${source.name}`);
      return expectedSha;
    } catch (err) {
      errors.push(`${source.name}: ${err.message}`);
    }
  }

  throw new Error(`failed to resolve checksums for ${assetName}\n${errors.join("\n")}`);
}

async function materializeArchive(sources, assetName, archivePath, expectedSha) {
  const errors = [];

  for (const source of sources) {
    const assetLocation = getAssetLocation(source, assetName);
    try {
      log(`trying ${source.name}: ${assetLocation}`);
      if (source.type === "local") {
        fs.copyFileSync(assetLocation, archivePath);
      } else {
        await downloadFile(assetLocation, archivePath);
      }
      const actualSha = sha256File(archivePath);
      if (actualSha !== expectedSha) {
        throw new Error(`checksum mismatch for ${assetName}`);
      }
      log(`downloaded from ${source.name}`);
      return;
    } catch (err) {
      fs.rmSync(archivePath, { force: true });
      errors.push(`${source.name}: ${err.message}`);
    }
  }

  throw new Error(`failed to obtain ${assetName}\n${errors.join("\n")}`);
}

function stripTrailingSlash(value) {
  return value.replace(/\/+$/, "");
}

async function fetchText(url) {
  try {
    const response = await fetch(url);
    if (!response.ok) {
      throw new Error(`HTTP ${response.status}`);
    }
    return response.text();
  } catch (err) {
    throw new Error(`failed to download ${url}: ${err.message}`);
  }
}

async function downloadFile(url, destination) {
  try {
    const response = await fetch(url);
    if (!response.ok) {
      throw new Error(`HTTP ${response.status}`);
    }
    const buffer = Buffer.from(await response.arrayBuffer());
    fs.writeFileSync(destination, buffer);
  } catch (err) {
    throw new Error(`failed to download ${url}: ${err.message}`);
  }
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
  console.error(
    "[vtrix installer] hint: GitHub download failed; package will fall back to bundled archives when available."
  );
  process.exit(1);
});
