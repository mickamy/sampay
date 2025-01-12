export interface FileLike {
  arrayBuffer(): Promise<ArrayBuffer>;

  name: string;
  type: string;
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
