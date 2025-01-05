import { createClient } from "@connectrpc/connect";
import { SessionService } from "@buf/mickamy_sampay.connectrpc_es/auth/v1/session_connect";
import { createConnectTransport } from "@connectrpc/connect-web";

const client = createClient(
  SessionService,
  createConnectTransport({
    baseUrl: "http://localhost:8080/api",
  })
);

export default function Index() {
  const signIn = async () => {
    try {
      const response = await client.signIn({
        emailOrSlug: "",
        password: "",
      });
      console.log(response);
    } catch (error) {
      console.error(error);
    }
  };

  return (
    <div>
      <button onClick={signIn}>Sign In</button>
    </div>
  );
}
