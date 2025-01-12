import { type RouteConfig, index, route } from "@react-router/dev/routes";

export default [
  index("routes/index.tsx"),

  route("admin", "routes/admin/admin-route.tsx"),
  route(
    "registration/sign-up",
    "routes/registration/sign-up/sign-up-route.tsx",
  ),
  route("onboarding", "routes/onboarding/onboarding-route.tsx"),
] satisfies RouteConfig;
