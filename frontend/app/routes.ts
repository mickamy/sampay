import {
  type RouteConfig,
  index,
  prefix,
  route,
} from "@react-router/dev/routes";

export default [
  index("routes/index.tsx"),

  route("admin", "routes/admin/admin-route.tsx"),
  route("account/sign-up", "routes/account/sign-up/sign-up-route.tsx"),
  ...prefix("auth", [
    route("/sign-in", "routes/auth/sign-in/sign-in-route.tsx"),
    route("/sign-out", "routes/auth/sign-out/sign-out-route.tsx"),
  ]),
  route("onboarding", "routes/onboarding/onboarding-route.tsx"),
] satisfies RouteConfig;
