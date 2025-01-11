import {
  type RouteConfig,
  index,
  prefix,
  route,
} from "@react-router/dev/routes";

export default [
  index("routes/index.tsx"),

  route("admin", "routes/admin/admin-route.tsx"),
  route(
    "registration/sign-up",
    "routes/registration/sign-up/sign-up-route.tsx",
  ),
  ...prefix("onboarding", [index("routes/onboarding/onboarding-index.tsx")]),
] satisfies RouteConfig;
