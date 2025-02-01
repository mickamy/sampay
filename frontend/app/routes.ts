import { type RouteConfig, index, route } from "@react-router/dev/routes";

export default [
  index("routes/static/landing-page-route.tsx"),
  route("privacy", "routes/static/privacy-route.tsx"),
  route("terms", "routes/static/terms-route.tsx"),

  route("admin", "routes/admin/admin-route.tsx"),
  route("onboarding", "routes/onboarding/onboarding-route.tsx"),
  route("oauth/callback/google", "routes/oauth/oauth-callback-route.tsx"),
  route("oauth/google", "routes/oauth/google/oauth-google-route.tsx"),
  route("reset-password", "routes/reset-password/reset-password-route.tsx"),
  route("sign-in", "routes/sign-in/sign-in-route.tsx"),
  route("sign-out", "routes/sign-out/sign-out-route.tsx"),
  route("sign-up", "routes/sign-up/sign-up-route.tsx"),
  route("u/:slug", "routes/user/user-route.tsx"),
] satisfies RouteConfig;
