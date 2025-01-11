import { type RouteConfig, index, route } from "@react-router/dev/routes";

export default [
  index("routes/index.tsx"),

  route("admin", "routes/admin/admin-route.tsx"),
  route("auth/sign-up", "routes/auth/sign-up/sign-up-route.tsx"),
] satisfies RouteConfig;
