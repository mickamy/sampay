import type { DescService } from "@bufbuild/protobuf";
import { type Client, createClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";
import {
  createI18NInterceptor,
  loggingInterceptor,
} from "~/lib/api/interceptors.server";

export const API_BASE_URL: string = process.env.PUBLIC_API_BASE_URL;

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
