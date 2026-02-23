import { mkdirSync } from "node:fs";
import { dirname } from "node:path";
import pino, { type Level } from "pino";

const validLevels: Level[] = [
  "fatal",
  "error",
  "warn",
  "info",
  "debug",
  "trace",
];

function isPinoLevel(value: string): value is Level {
  return validLevels.includes(value as Level);
}

async function init() {
  let level: Level = process.env.NODE_ENV === "development" ? "debug" : "info";
  const envLevel = import.meta.env.VITE_LOG_LEVEL;
  if (envLevel && isPinoLevel(envLevel)) {
    level = envLevel;
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
