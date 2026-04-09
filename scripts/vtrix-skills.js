#!/usr/bin/env node

const fs = require("fs");
const path = require("path");
const { spawn } = require("child_process");

const rootDir = path.resolve(__dirname, "..");
const exeName = process.platform === "win32" ? "vtrix.exe" : "vtrix";
const exePath = path.join(rootDir, "vendor", exeName);

if (!fs.existsSync(exePath)) {
  console.error("vtrix binary is not installed.");
  console.error("Reinstall the package: npm install -g @vtrixai/vtrix-cli");
  process.exit(1);
}

const child = spawn(exePath, ["skills", ...process.argv.slice(2)], {
  stdio: "inherit",
});

child.on("error", (err) => {
  console.error(`failed to start vtrix skills: ${err.message}`);
  process.exit(1);
});

child.on("exit", (code, signal) => {
  if (signal) {
    process.kill(process.pid, signal);
    return;
  }
  process.exit(code ?? 0);
});
