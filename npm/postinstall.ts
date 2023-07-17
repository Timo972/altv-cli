#!/usr/bin/env node
"use strict";
import * as fs from "node:fs/promises";
import { createWriteStream } from "node:fs";
import * as path from "node:path";
import * as http from "node:http";

interface PackageJSON {
  version: string;
}

const binaryURL = (version: string, bin: string): string =>
  `https://github.com/Timo972/altv-cli/releases/download/${version}/${bin}`;

async function install(): Promise<void> {
  const pkgJSON = await fs.readFile(path.join(process.cwd(), "package.json"), {
    encoding: "utf-8",
  });
  const pkg: PackageJSON = JSON.parse(pkgJSON);

  const version = `v${pkg.version}`;
  const binaryName = process.platform === "win32" ? "altv-cli.exe" : "altv-cli";

  const resp = await new Promise<http.IncomingMessage>((resolve) =>
    http.get(binaryURL(version, binaryName), resolve)
  );

  await new Promise((resolve, reject) =>
    resp
      .pipe(createWriteStream(binaryName))
      .on("finish", resolve)
      .on("error", reject)
      .on("close", resolve)
  );
}

install();
