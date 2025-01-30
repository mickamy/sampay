import { join } from "pathe";
import pino from "pino";

import { isBrowser } from "~/lib/utils";

async function init() {
  let level = process.env.NODE_ENV === "development" ? "debug" : "info";
  if (
    typeof import.meta.env !== "undefined" &&
    import.meta.env.VITE_LOG_LEVEL
  ) {
    level = import.meta.env.VITE_LOG_LEVEL;
  } else if (process.env.VITE_LOG_LEVEL) {
    level = process.env.VITE_LOG_LEVEL;
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

  try {
    if (!existsSync(dir)) {
      mkdirSync(dir, { recursive: true });
    }
  } catch (error) {
    console.error(`failed to create log directory: ${dir}`);
    throw error;
  }

  const path = join(dir, "frontend.log");
  const fileStream = createWriteStream(path, {
    flags: "a",
  });

  return pino(
    { level },
    pino.multistream([
      { stream: process.stdout, level },
      { stream: fileStream, level },
    ]),
  );
}

const logger = await init();

export default logger;
