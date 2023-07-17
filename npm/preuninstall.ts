#!/usr/bin/env node
"use strict";
import * as fs from "node:fs/promises";

async function uninstall(): Promise<void> {
  const binaryName = process.platform === "win32" ? "altv-cli.exe" : "altv-cli";

  await fs.unlink(binaryName);
}

uninstall();
