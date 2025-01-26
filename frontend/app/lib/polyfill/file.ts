export interface FileLike {
  name: string;
  type: string;

  arrayBuffer(): Promise<ArrayBuffer>;
}

export function isFileLike(value: unknown) {
  return (
    value instanceof File ||
    (value &&
      typeof value === "object" &&
      "arrayBuffer" in value &&
      "name" in value &&
      "type" in value)
  );
}
