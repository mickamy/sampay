import type { UserLink as UserLinkPB } from "@buf/mickamy_sampay.bufbuild_es/user/v1/user_link_pb";
import {
  type UserLinkDisplayAttribute,
  convertToUserLinkDisplayAttributes,
} from "~/models/user/user-link-display-attribute-model";
import type { UserLinkProviderType } from "~/models/user/user-link-provider-type-model";

export interface UserLink {
  id: string;
  providerType: UserLinkProviderType;
  uri: string;
  displayAttribute: UserLinkDisplayAttribute;
}

export function convertToUserLink(pb: UserLinkPB): UserLink {
  if (!pb.displayAttribute) {
    throw new Error("displayAttribute is required");
  }
  return {
    id: pb.id,
    providerType: pb.providerType,
    uri: pb.uri,
    displayAttribute: convertToUserLinkDisplayAttributes(pb.displayAttribute),
  };
}

export function convertToUserLinks(pbs: UserLinkPB[]): UserLink[] {
  return pbs.map(convertToUserLink);
}
