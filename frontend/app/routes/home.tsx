import { m } from "~/paraglide/messages";

export function meta() {
  return [
    { title: "New React Router App" },
    { name: "description", content: "Welcome to React Router!" },
  ];
}

export default function Home() {
  return <div>{m.example_message({ username: "mickamy" })}</div>;
}
