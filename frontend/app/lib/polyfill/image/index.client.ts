import jsQR from "jsqr";

type DecoderFunction = (
  data: Uint8Array,
  mimeType: string,
) => Promise<{ pixels: Uint8ClampedArray; width: number; height: number }>;

async function decodeWithCanvas(data: Uint8Array, mimeType: string) {
  const blob = new Blob([data], { type: mimeType });
  const img = new Image();
  img.src = URL.createObjectURL(blob);

  await new Promise((resolve, reject) => {
    img.onload = resolve;
    img.onerror = reject;
  });

  const canvas = document.createElement("canvas");
  const context = canvas.getContext("2d");
  if (!context) {
    throw new Error("Failed to get canvas context");
  }

  canvas.width = img.width;
  canvas.height = img.height;
  context.drawImage(img, 0, 0);

  return {
    pixels: context.getImageData(0, 0, img.width, img.height).data,
    width: img.width,
    height: img.height,
  };
}

const canvasDecoders: Record<string, DecoderFunction> = {
  "image/jpeg": (data, mimeType) => decodeWithCanvas(data, mimeType),
  "image/jpg": (data, mimeType) => decodeWithCanvas(data, mimeType),
  "image/png": (data, mimeType) => decodeWithCanvas(data, mimeType),
  "image/gif": (data, mimeType) => decodeWithCanvas(data, mimeType),
  "image/bmp": (data, mimeType) => decodeWithCanvas(data, mimeType),
};

export async function parseQRCode(file: File): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();

    reader.onerror = () => reject(new Error("failed to read file"));

    reader.onload = async () => {
      const data = new Uint8Array(reader.result as ArrayBuffer);
      const decoder = canvasDecoders[file.type];
      if (!decoder) {
        return reject(new Error("unsupported file type"));
      }

      try {
        const { pixels, width, height } = await decoder(data, file.type);
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
