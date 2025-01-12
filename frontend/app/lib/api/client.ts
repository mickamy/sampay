import type { DescService } from "@bufbuild/protobuf";
import { type Client, createClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";
import {
  createAuthenticateInterceptor,
  createI18NInterceptor,
  loggingInterceptor,
} from "~/lib/api/interceptors";
import { isBrowser } from "~/lib/utils";

export const API_BASE_URL: string = isBrowser()
  ? window.ENV.PUBLIC_API_BASE_URL
  : process.env.PUBLIC_API_BASE_URL;

const transport = createConnectTransport({
  baseUrl: API_BASE_URL,
  interceptors: [loggingInterceptor],
});

export function getClient<T extends DescService>({
  service,
  request,
}: { service: T; request: Request }): Client<T> {
  const transport = createConnectTransport({
    baseUrl: API_BASE_URL,
    interceptors: [loggingInterceptor, createI18NInterceptor(request)],
  });
  return createClient(service, transport);
}
