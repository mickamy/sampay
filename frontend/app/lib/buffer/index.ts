function bufferToString(buf: ArrayBuffer): string {
  const decoder = new TextDecoder("utf-16");
  return decoder.decode(buf);
}

export function arrayBufferToString(buf: ArrayBuffer): string {
  const tmp: string[] = [];
  const len = 1024;
  for (let offset = 0; offset < buf.byteLength; offset += len) {
    tmp.push(bufferToString(buf.slice(offset, offset + len)));
  }
  return tmp.join("");
}
