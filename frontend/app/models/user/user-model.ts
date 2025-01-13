import type { User as UserPB } from "@buf/mickamy_sampay.bufbuild_es/user/v1/user_pb";
import {
  type UserProfile,
  convertToUserProfile,
} from "~/models/user/user-profile-model";

export interface User {
  id: string;
  slug: string;
  profile: UserProfile;
}

export function convertToUser(pb: UserPB): User {
  if (!pb.profile) {
    throw new Error("profile is required");
  }
  return {
    id: pb.id,
    slug: pb.slug,
    profile: convertToUserProfile(pb.profile),
  };
}
