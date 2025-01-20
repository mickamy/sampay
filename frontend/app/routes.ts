import {
  type RouteConfig,
  index,
  prefix,
  route,
} from "@react-router/dev/routes";

export default [
  index("routes/static/landing-page-route.tsx"),

  route("account/sign-up", "routes/account/sign-up/sign-up-route.tsx"),
  ...prefix("auth", [
    route(
      "/reset-password",
      "routes/auth/reset-password/reset-password-route.tsx",
    ),
    route("/sign-in", "routes/auth/sign-in/sign-in-route.tsx"),
    route("/sign-out", "routes/auth/sign-out/sign-out-route.tsx"),
  ]),
  route("admin", "routes/admin/admin-route.tsx"),
  route("onboarding", "routes/onboarding/onboarding-route.tsx"),
  route("privacy", "routes/static/privacy-route.tsx"),
  route("terms", "routes/static/terms-route.tsx"),
  route("u/:slug", "routes/user/user-route.tsx"),
] satisfies RouteConfig;
