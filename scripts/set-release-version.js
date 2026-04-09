#!/usr/bin/env node

const fs = require("fs");
const path = require("path");

const rootDir = path.resolve(__dirname, "..");
const npmPackagePath = path.join(rootDir, "package.json");

const input = process.argv[2];
if (!input) {
  console.error("usage: node scripts/set-release-version.js <version>");
  console.error("example: node scripts/set-release-version.js v0.1.1");
  process.exit(1);
}

const normalizedVersion = normalizeVersion(input);
const npmVersion = normalizedVersion.replace(/^v/, "");

updateJsonVersion(npmPackagePath, npmVersion);

console.log(`updated npm package version to ${npmVersion}`);

function normalizeVersion(version) {
  const trimmed = version.trim();
  const withPrefix = trimmed.startsWith("v") ? trimmed : `v${trimmed}`;
  if (!/^v\d+\.\d+\.\d+(?:[-+][0-9A-Za-z.-]+)?$/.test(withPrefix)) {
    throw new Error(`invalid version: ${version}`);
  }
  return withPrefix;
}

function updateJsonVersion(filePath, version) {
  const data = JSON.parse(fs.readFileSync(filePath, "utf8"));
  data.version = version;
  fs.writeFileSync(filePath, `${JSON.stringify(data, null, 2)}\n`);
}
