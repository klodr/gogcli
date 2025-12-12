import { mkdirSync } from "node:fs";
import { join } from "node:path";
import { spawnSync } from "node:child_process";

function run(cmd, args) {
  const res = spawnSync(cmd, args, { stdio: "inherit" });
  if (res.error) throw res.error;
  if (res.status !== 0) {
    process.exit(typeof res.status === "number" ? res.status : 1);
  }
  return res.status;
}

const repoRoot = process.cwd();
const binDir = join(repoRoot, "bin");
mkdirSync(binDir, { recursive: true });

const exe = process.platform === "win32" ? "gog.exe" : "gog";
const binPath = join(binDir, exe);

run("go", ["build", "-o", binPath, "./cmd/gog"]);

const final = spawnSync(binPath, process.argv.slice(2), { stdio: "inherit" });
if (final.error) throw final.error;
process.exit(typeof final.status === "number" ? final.status : 1);
