import type { DescService } from "@bufbuild/protobuf";
import {
  type Client,
  createClient,
  type Interceptor,
} from "@connectrpc/connect";
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
  interceptors = [],
}: {
  service: T;
  request: Request;
  interceptors?: Interceptor[];
}): Client<T> {
  const transport = createConnectTransport({
    baseUrl: API_BASE_URL,
    interceptors: [
      loggingInterceptor,
      createI18NInterceptor(request),
      ...interceptors,
    ],
  });
  return createClient(service, transport);
}
