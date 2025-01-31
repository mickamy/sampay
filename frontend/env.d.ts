declare namespace NodeJS {
  interface ProcessEnv {
    ENVIRONMENT: "development" | "staging" | "production";
    PUBLIC_API_BASE_URL: string;
    PUBLIC_BUCKET_NAME: string;
  }
}
