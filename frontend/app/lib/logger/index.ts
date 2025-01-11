import { join } from "pathe";
import pino from "pino";

import { isBrowser } from "~/lib/utils";

async function init() {
  let level = process.env.NODE_ENV === "development" ? "debug" : "info";
  if (import.meta.env.VITE_LOG_LEVEL) {
    level = import.meta.env.VITE_LOG_LEVEL;
  }

  if (isBrowser()) {
    return pino({
      level,
      browser: {
        asObject: true,
      },
    });
  }
  const dir = "/var/log/sampay";
  const { createWriteStream, existsSync, mkdirSync } = await import("node:fs");
  if (!existsSync(dir)) {
    mkdirSync(dir, { recursive: true });
  }

  const path = join(dir, "web.log");
  const fileStream = createWriteStream(path, {
    flags: "a",
  });

  return pino(
    { level },
    pino.multistream([{ stream: process.stdout }, { stream: fileStream }]),
  );
}

const logger = await init();

export default logger;
