export {};

declare global {
  // eslint-disable-next-line no-var
  var environment: Environment | undefined;

  interface Environment {
    API_BASE_URL?: string;
    S3_PUBLIC_BUCKET_NAME?: string;
    SESSION_SECRET?: string;
  }

  interface Window {
    ENV: {
      PUBLIC_API_BASE_URL: string;
    };
  }
}
