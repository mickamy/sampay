import type { UserLinkDisplayAttribute as UserLinkDisplayAttributePB } from "@buf/mickamy_sampay.bufbuild_es/user/v1/user_link_pb";

export interface UserLinkDisplayAttribute {
  name: string;
  displayOrder: number;
}

export function convertToUserLinkDisplayAttributes(
  pb: UserLinkDisplayAttributePB,
): UserLinkDisplayAttribute {
  return {
    name: pb.name,
    displayOrder: pb.displayOrder,
  };
}
