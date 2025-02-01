import { DirectUploadURLService } from "@buf/mickamy_sampay.bufbuild_es/common/v1/direct_upload_url_pb";
import type { getClientType } from "~/lib/api/request.server";
import { getParameter } from "~/lib/aws/ssm";
import logger from "~/lib/logger";
import { randomUUID } from "~/lib/polyfill/crypto";
import type { S3Object } from "~/models/common/s3-object-model";

type DirectUploadFileType = "profile_image" | "qr_code";

export async function directUpload({
  type,
  file,
  getClient,
}: {
  type: DirectUploadFileType;
  file: File;
  getClient: getClientType;
}): Promise<S3Object> {
  try {
    const obj: S3Object = {
      bucket: await bucket(type),
      key: await key(type),
      contentType: file.type,
    };
    const { url } = await getClient(
      DirectUploadURLService,
    ).createDirectUploadURL({
      s3Object: obj,
    });
    const res = await fetch(url, {
      method: "PUT",
      body: file,
      headers: {
        "Content-Type": file.type,
      },
    });
    if (!res.ok) {
      logger.error({ response: res }, "failed to upload file");
      return Promise.reject(new Error("failed to upload file"));
    }
    return obj;
  } catch (e) {
    logger.error({ error: e }, "failed to upload file");
    throw e;
  }
}

function bucket(type: DirectUploadFileType) {
  switch (type) {
    case "profile_image":
      return publicBucketName();
    case "qr_code":
      return publicBucketName();
    default:
      throw new Error(`unknown type: ${type}`);
  }
}

async function key(type: DirectUploadFileType) {
  switch (type) {
    case "profile_image":
      return `profile_images/${Date.now()}_${await randomUUID()}`;
    case "qr_code":
      return `qr_codes/${Date.now()}_${await randomUUID()}`;
    default:
      throw new Error(`unknown type: ${type}`);
  }
}

async function publicBucketName(): Promise<string> {
  if (!global.environment?.S3_PUBLIC_BUCKET_NAME) {
    try {
      const sessionSecret = await getParameter({ name: "SESSION_SECRET" });
      global.environment = {
        ...global.environment,
        SESSION_SECRET: sessionSecret,
      };
    } catch (e) {
      logger.error({ error: e }, "failed to retrieve SSM parameters");
      throw e;
    }
  }

  if (!global.environment.S3_PUBLIC_BUCKET_NAME) {
    throw new Error("S3_PUBLIC_BUCKET_NAME is not set");
  }
  return global.environment.S3_PUBLIC_BUCKET_NAME;
}
