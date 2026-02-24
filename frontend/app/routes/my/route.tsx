import { Outlet } from "react-router";
import { Header } from "~/components/header";

export default function MyLayout() {
  return (
    <div className="min-h-screen bg-background flex flex-col">
      <Header isLoggedIn={true} />
      <main className="flex-1">
        <div className="container mx-auto px-4 py-8 max-w-2xl">
          <Outlet />
        </div>
      </main>
    </div>
  );
}
