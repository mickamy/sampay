import type { DescService } from "@bufbuild/protobuf";
import { type Client, createClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";
import {
  createI18NInterceptor,
  loggingInterceptor,
} from "~/lib/api/interceptors.server";

export function getClient<T extends DescService>({
  service,
  request,
}: { service: T; request: Request }): Client<T> {
  const baseURL = process.env.API_BASE_URL;
  if (!baseURL) {
    throw new Error("API_BASE_URL is not set");
  }
  const transport = createConnectTransport({
    baseUrl: baseURL,
    interceptors: [createI18NInterceptor(request), loggingInterceptor],
  });
  return createClient(service, transport);
}
