interface MetaEntry {
  title?: string;
  name?: string;
  property?: string;
  content?: string;
}

const COMMON_META: MetaEntry[] = [
  { property: "og:type", content: "website" },
  { property: "og:site_name", content: "Sampay" },
  { property: "og:locale", content: "ja_JP" },
  { name: "twitter:card", content: "summary_large_image" },
];

export function buildMeta({
  title,
  description,
  image,
}: {
  title: string;
  description: string;
  image?: string;
}): MetaEntry[] {
  const meta: MetaEntry[] = [
    { title },
    { name: "description", content: description },
    { property: "og:title", content: title },
    { property: "og:description", content: description },
    ...COMMON_META,
  ];
  if (image) {
    meta.push(
      { property: "og:image", content: image },
      { name: "twitter:image", content: image },
    );
  }
  return meta;
}
