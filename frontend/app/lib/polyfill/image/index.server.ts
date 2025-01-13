import bmp from "bmp-js";
import jpeg from "jpeg-js";
import jsQR from "jsqr";
import { GifReader } from "omggif";
import { PNG } from "pngjs";

type DecoderFunction = (data: Uint8Array) => {
  pixels: Uint8ClampedArray;
  width: number;
  height: number;
};

function decodeJpeg(data: Uint8Array) {
  const jpegData = jpeg.decode(data);
  return {
    pixels: new Uint8ClampedArray(jpegData.data),
    width: jpegData.width,
    height: jpegData.height,
  };
}

const decoders: Record<string, DecoderFunction> = {
  "image/jpeg": (data) => {
    return decodeJpeg(data);
  },
  "image/jpg": (data) => {
    return decodeJpeg(data);
  },
  "image/png": (data) => {
    const png = PNG.sync.read(Buffer.from(data));
    return {
      pixels: new Uint8ClampedArray(png.data),
      width: png.width,
      height: png.height,
    };
  },
  "image/gif": (data) => {
    const reader = new GifReader(Buffer.from(data));
    const pixels = new Uint8ClampedArray(reader.width * reader.height * 4);
    reader.decodeAndBlitFrameRGBA(0, pixels);
    return {
      pixels,
      width: reader.width,
      height: reader.height,
    };
  },
  "image/bmp": (data) => {
    const bmpData = bmp.decode(Buffer.from(data));
    return {
      pixels: new Uint8ClampedArray(bmpData.data),
      width: bmpData.width,
      height: bmpData.height,
    };
  },
};

export async function parseQRCode(file: File): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();

    reader.onerror = () => reject(new Error("failed to read file"));

    reader.onload = () => {
      const data = new Uint8Array(reader.result as ArrayBuffer);

      const decoder = decoders[file.type];
      if (!decoder) {
        return reject(new Error("unsupported file type"));
      }

      try {
        const { pixels, width, height } = decoder(data);
        const decoded = jsQR(pixels, width, height);

        if (decoded) {
          resolve(decoded.data);
        } else {
          reject(new Error("failed to parse QR code"));
        }
      } catch (error) {
        reject(new Error(`failed to decode ${file.type}: ${error}`));
      }
    };

    reader.readAsArrayBuffer(file);
  });
}
