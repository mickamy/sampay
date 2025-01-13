import type { UserProfile as UserProfilePB } from "@buf/mickamy_sampay.bufbuild_es/user/v1/user_profile_pb";

export interface UserProfile {
  name: string;
  bio?: string;
  imageURL?: string;
}

export function convertToUserProfile(pb: UserProfilePB): UserProfile {
  return {
    name: pb.name,
    bio: pb.bio,
    imageURL: pb.imageUrl,
  };
}
