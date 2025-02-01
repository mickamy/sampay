import { GetParameterCommand, SSMClient } from "@aws-sdk/client-ssm";

type ParameterName =
  | "SESSION_SECRET"
  | "S3_PUBLIC_BUCKET_NAME"
  | "S3_PRIVATE_BUCKET_NAME";

export async function getParameter({
  name,
}: {
  name: ParameterName;
}): Promise<string> {
  if (process.env.NODE_ENV === "development") {
    return getFromEnvironment(name);
  }
  const client = new SSMClient({ region: "ap-northeast-1" });
  const command = new GetParameterCommand({
    Name: `/sampay/app/${getShortEnvironment()}/${name}`,
    WithDecryption: true,
  });
  return client.send(command).then((response) => {
    if (!response.Parameter?.Value) {
      throw new Error(`missing SSM parameter: ${name}`);
    }
    return response.Parameter.Value;
  });
}

function getFromEnvironment(name: string): string {
  const value = process.env[name];
  if (value == null) {
    throw new Error(`missing env: ${name}`);
  }
  return value;
}

function getShortEnvironment(): string {
  return process.env.ENVIRONMENT === "production"
    ? "prod"
    : process.env.ENVIRONMENT === "staging"
      ? "stg"
      : "dev";
}
