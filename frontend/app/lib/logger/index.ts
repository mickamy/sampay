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

  return pino({ level }, pino.multistream([{ stream: process.stdout, level }]));
}

const logger = await init();

export default logger;
