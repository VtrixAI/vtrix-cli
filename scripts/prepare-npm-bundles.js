#!/usr/bin/env node

const crypto = require("crypto");
const fs = require("fs");
const path = require("path");

const rootDir = path.resolve(__dirname, "..");
const pkg = require(path.join(rootDir, "package.json"));

const version = pkg.version;
const projectName = pkg.vtrix?.projectName || "vtrix";
const distDir = path.join(rootDir, "dist");
const outputDir = path.join(
  rootDir,
  pkg.vtrix?.bundledAssetsDir || "npm-bundles"
);

const artifacts = [
  {
    source: `${projectName}_darwin_amd64.tar.gz`,
    target: `${projectName}_${version}_darwin_amd64.tar.gz`,
  },
  {
    source: `${projectName}_darwin_arm64.tar.gz`,
    target: `${projectName}_${version}_darwin_arm64.tar.gz`,
  },
  {
    source: `${projectName}_linux_amd64.tar.gz`,
    target: `${projectName}_${version}_linux_amd64.tar.gz`,
  },
  {
    source: `${projectName}_linux_arm64.tar.gz`,
    target: `${projectName}_${version}_linux_arm64.tar.gz`,
  },
  {
    source: `${projectName}_windows_amd64.zip`,
    target: `${projectName}_${version}_windows_amd64.zip`,
  },
];

if (!fs.existsSync(distDir)) {
  throw new Error(`dist directory not found: ${distDir}`);
}

fs.rmSync(outputDir, { recursive: true, force: true });
fs.mkdirSync(outputDir, { recursive: true });

const checksums = [];

for (const artifact of artifacts) {
  const sourcePath = path.join(distDir, artifact.source);
  const targetPath = path.join(outputDir, artifact.target);

  if (!fs.existsSync(sourcePath)) {
    throw new Error(`missing build artifact: ${artifact.source}`);
  }

  fs.copyFileSync(sourcePath, targetPath);
  checksums.push(`${sha256File(targetPath)}  ${artifact.target}`);
}

fs.writeFileSync(path.join(outputDir, "SHA256SUMS"), `${checksums.join("\n")}\n`);

console.log(`prepared npm bundles in ${path.relative(rootDir, outputDir)}`);

function sha256File(filePath) {
  const hash = crypto.createHash("sha256");
  hash.update(fs.readFileSync(filePath));
  return hash.digest("hex");
}
