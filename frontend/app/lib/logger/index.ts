const logger = import.meta.env.SSR
  ? (await import("./node")).default
  : (await import("./browser")).default;

export default logger;
