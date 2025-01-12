import { DirectUploadURLService } from "@buf/mickamy_sampay.bufbuild_es/common/v1/direct_upload_url_pb";
import type { getClientType } from "~/lib/api/request";
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
    const obj = {
      bucket: bucket(type),
      key: await key(type),
      contentType: file.type,
    };
    const { url } = await getClient(DirectUploadURLService).request({
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
      logger.error("failed to upload file", res);
      return Promise.reject(new Error("failed to upload file"));
    }
    return obj;
  } catch (e) {
    console.error("failed to upload file", e);
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

function publicBucketName(): string {
  if (process.env.PUBLIC_BUCKET_NAME) {
    return process.env.PUBLIC_BUCKET_NAME;
  }
  throw new Error("missing PUBLIC_BUCKET_NAME");
}
