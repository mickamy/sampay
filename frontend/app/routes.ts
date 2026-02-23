import { index, type RouteConfig, route } from "@react-router/dev/routes";

export default [
  index("routes/home.tsx"),

  route("oauth/:provider", "routes/oauth/provider/route.tsx"),
  route("oauth/callback", "routes/oauth/callback/route.tsx"),
  route("auth/logout", "routes/auth/logout/route.tsx"),
] satisfies RouteConfig;
