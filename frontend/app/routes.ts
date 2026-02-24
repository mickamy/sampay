import { index, type RouteConfig, route } from "@react-router/dev/routes";

export default [
  index("routes/home.tsx"),

  route("my", "routes/my/route.tsx", [
    index("routes/my/index/route.tsx"),
    route("edit", "routes/my/edit/route.tsx"),
  ]),

  route("u/:slug", "routes/u/slug/route.tsx"),

  route("enter", "routes/enter/route.tsx"),
  route("api/check-slug", "routes/api/check-slug.ts"),

  route("oauth/:provider", "routes/oauth/provider/route.tsx"),
  route("oauth/callback", "routes/oauth/callback/route.tsx"),
  route("auth/logout", "routes/auth/logout/route.tsx"),
] satisfies RouteConfig;
