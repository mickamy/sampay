import { mkdirSync } from "node:fs";
import { dirname } from "node:path";
import pino, { type Level } from "pino";

async function init() {
  let level: Level = process.env.NODE_ENV === "development" ? "debug" : "info";
  if (import.meta.env.VITE_LOG_LEVEL) {
    level = import.meta.env.VITE_LOG_LEVEL;
  }

  if (process.env.NODE_ENV === "development") {
    return pino({
      level,
      transport: {
        target: "pino-pretty",
      },
    });
  }

  const streams: pino.StreamEntry[] = [{ stream: process.stdout, level }];

  if (process.env.LOG_FILE) {
    mkdirSync(dirname(process.env.LOG_FILE), { recursive: true });
    streams.push({ stream: pino.destination(process.env.LOG_FILE), level });
  }

  return pino({ level }, pino.multistream(streams));
}

const logger = await init();

export default logger;
