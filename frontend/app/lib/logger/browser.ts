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
  let level: Level = import.meta.env.DEV ? "debug" : "info";
  const envLevel = import.meta.env.VITE_LOG_LEVEL;
  if (envLevel && isPinoLevel(envLevel)) {
    level = envLevel;
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
