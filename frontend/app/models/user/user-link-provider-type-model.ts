export const UserLinkProviderTypes = ["kyash", "paypay", "amazon", "other"];

export type UserLinkProviderType = (typeof UserLinkProviderTypes)[number];

export function getUserLinkProviderTypeImage(
  provider: UserLinkProviderType,
): string {
  switch (provider) {
    case "kyash":
      return "/provider/kyash.png";
    case "paypay":
      return "/provider/paypay.jpg";
    case "amazon":
      return "/provider/amazon.png";
    default:
      return "/provider/other.svg";
  }
}
