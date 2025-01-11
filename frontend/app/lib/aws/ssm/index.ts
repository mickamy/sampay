import { GetParameterCommand, SSMClient } from "@aws-sdk/client-ssm";

type ParameterName = "API_BASE_URL" | "SESSION_SECRET";

export async function getParameter({
  name,
}: {
  name: ParameterName;
}): Promise<string> {
  if (process.env.NODE_ENV === "development") {
    return getEnv(name);
  }
  const client = new SSMClient({ region: "ap-northeast-1" });
  const command = new GetParameterCommand({
    Name: `/sampay/config/${name}`,
    WithDecryption: true,
  });
  return client.send(command).then((response) => {
    if (!response.Parameter?.Value) {
      throw new Error(`missing SSM parameter: ${name}`);
    }
    return response.Parameter.Value;
  });
}

function getEnv(name: string): string {
  const value = process.env[name];
  if (value == null) {
    throw new Error(`missing env: ${name}`);
  }
  return value;
}
