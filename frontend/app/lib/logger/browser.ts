import pino, { type Level } from "pino";

async function init() {
  let level: Level = process.env.NODE_ENV === "development" ? "debug" : "info";
  if (import.meta.env.VITE_LOG_LEVEL) {
    level = import.meta.env.VITE_LOG_LEVEL;
  }

  return pino({
    level,
    browser: {
      asObject: true,
    },
  });
}

const logger = await init();

export default logger;
