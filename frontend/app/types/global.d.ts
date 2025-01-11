export {};

declare global {
  // eslint-disable-next-line no-var
  var environment: Environment | undefined;

  interface Environment {
    API_BASE_URL?: string;
    SESSION_SECRET?: string;
  }
}
