import { index, type RouteConfig, route } from "@react-router/dev/routes";

export default [
  index("routes/home.tsx"),

  route("my", "routes/my/route.tsx", [
    index("routes/my/index/route.tsx"),
    route("edit", "routes/my/edit/route.tsx"),
    route("events/new", "routes/my/events/new/route.tsx"),
    route("events/:id", "routes/my/events/$id/route.tsx"),
    route("events/:id/edit", "routes/my/events/$id_.edit/route.tsx"),
  ]),

  route("e/:id", "routes/e/$id/route.tsx"),
  route("u/:slug", "routes/u/slug/route.tsx"),
  route("og/u/:slug.png", "routes/og/u/slug.tsx"),
  route("og/e/:id.png", "routes/og/e/$id.tsx"),

  route("enter", "routes/enter/route.tsx"),
  route("api/check-slug", "routes/api/check-slug.ts"),

  route("oauth/:provider", "routes/oauth/provider/route.tsx"),
  route("oauth/callback", "routes/oauth/callback/route.tsx"),
  route("auth/logout", "routes/auth/logout/route.tsx"),
] satisfies RouteConfig;
