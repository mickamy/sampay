import jsQR from "jsqr";

type DecoderFunction = (
  data: Uint8Array,
  mimeType: string,
) => Promise<{ pixels: Uint8ClampedArray; width: number; height: number }>;

const canvasDecoders: Record<string, DecoderFunction> = {
  "image/jpeg": (data, mimeType) => decodeWithCanvas(data, mimeType),
  "image/jpg": (data, mimeType) => decodeWithCanvas(data, mimeType),
  "image/png": (data, mimeType) => decodeWithCanvas(data, mimeType),
  "image/gif": (data, mimeType) => decodeWithCanvas(data, mimeType),
  "image/bmp": (data, mimeType) => decodeWithCanvas(data, mimeType),
};

async function decodeWithCanvas(data: Uint8Array, mimeType: string) {
  const blob = new Blob([data], { type: mimeType });
  const img = new Image();
  const objectURL = URL.createObjectURL(blob);

  try {
    await new Promise((resolve, reject) => {
      img.onload = resolve;
      img.onerror = reject;
      img.src = objectURL;
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
  } finally {
    URL.revokeObjectURL(objectURL);
  }
}

export async function parseQRCode(file: File): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();

    reader.onerror = () => {
      return reject(new Error("failed to read file"));
    };

    reader.onload = async () => {
      try {
        const data = new Uint8Array(reader.result as ArrayBuffer);
        const decoder = canvasDecoders[file.type];
        if (!decoder) {
          return reject(new Error("unsupported file type"));
        }

        const { pixels, width, height } = await decoder(data, file.type);
        const decoded = jsQR(pixels, width, height);

        if (decoded) {
          return resolve(decoded.data);
        }

        return reject(Error("failed to parse QR code"));
      } catch (error) {
        return reject(error);
      }
    };

    reader.readAsArrayBuffer(file);
  });
}
