#!/usr/bin/env node
"use strict";
import * as fs from "node:fs/promises";
import { createWriteStream } from "node:fs";
import * as path from "node:path";
import * as http from "node:http";
import * as https from "node:https";

interface PackageJSON {
  version: string;
}

const binaryURL = (version: string, bin: string): string =>
  `https://github.com/Timo972/altv-cli/releases/download/${version}/${bin}`;

async function get(
  url: string,
  options: http.RequestOptions
): Promise<http.IncomingMessage> {
  const resp = await new Promise<http.IncomingMessage>((resolve, reject) => {
    const req = https.get(url, options, resolve);
    req.on("error", reject);
  });

  if (!resp.statusCode) return resp;
  if (!resp.headers.location) return resp;

  return get(resp.headers.location, options);
}

async function install(): Promise<void> {
  const pkgJSON = await fs.readFile(path.join(process.cwd(), "package.json"), {
    encoding: "utf-8",
  });
  const pkg: PackageJSON = JSON.parse(pkgJSON);

  const version = `v${pkg.version}`;
  const binaryName = process.platform === "win32" ? "altv.exe" : "altv";

  const url = binaryURL(version, binaryName);
  const resp = await get(url, {});

  await new Promise((resolve, reject) =>
    resp
      .pipe(createWriteStream(binaryName))
      .on("finish", resolve)
      .on("error", reject)
      .on("close", resolve)
  );
}

install();
