import type { DescService } from "@bufbuild/protobuf";
import { type Client, createClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";
import {
  createI18NInterceptor,
  loggingInterceptor,
} from "~/lib/api/interceptors.server";

const baseURL = process.env.API_BASE_URL;
if (!baseURL) {
  throw new Error("API_BASE_URL is not set");
}
export const API_BASE_URL = baseURL;

export function getClient<T extends DescService>({
  service,
  request,
}: { service: T; request: Request }): Client<T> {
  const transport = createConnectTransport({
    baseUrl: API_BASE_URL,
    interceptors: [createI18NInterceptor(request), loggingInterceptor],
  });
  return createClient(service, transport);
}
