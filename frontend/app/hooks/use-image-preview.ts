import { useCallback, useState } from "react";

export default function useImagePreview(initialImageURL?: string) {
  const [imageURL, setImageURL] = useState<string | undefined>(initialImageURL);

  const onImageChange = useCallback((file: File | null) => {
    if (file) {
      const reader = new FileReader();
      reader.onloadend = () => {
        if (reader.result) {
          setImageURL(reader.result as string);
        }
      };
      reader.readAsDataURL(file);
    }
  }, []);

  return { imageURL, onImageChange };
}
