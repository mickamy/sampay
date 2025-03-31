package seed

import (
	"context"
	"fmt"
	"os"
	"path"

	"mickamy.com/sampay/config"
	"mickamy.com/sampay/internal/di"
	authFixture "mickamy.com/sampay/internal/domain/auth/fixture"
	authModel "mickamy.com/sampay/internal/domain/auth/model"
	commonModel "mickamy.com/sampay/internal/domain/common/model"
	commonRepository "mickamy.com/sampay/internal/domain/common/repository"
	userModel "mickamy.com/sampay/internal/domain/user/model"
	"mickamy.com/sampay/internal/infra/storage/database"
	"mickamy.com/sampay/internal/lib/ptr"
)

func seedUser(ctx context.Context, writer *database.Writer, env config.Env) error {
	if env != config.Development {
		fmt.Println("do not seed user because env is not development")
		return nil
	}

	imgKey := "mickamy.jpeg"
	img := path.Join(config.Common().PackageRoot, "internal", "cli", "db", "seed", "testdata", imgKey)
	file, err := os.Open(img)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	if err := di.InitLibs().S3.PutObject(ctx, config.AWS().S3PublicBucket, imgKey, file); err != nil {
		return fmt.Errorf("failed to put object: %w", err)
	}

	imgObj := commonModel.S3Object{
		Bucket: config.AWS().S3PublicBucket,
		Key:    imgKey,
	}

	if err := writer.WriterTransaction(ctx, func(tx database.WriterTransactional) error {
		if err := commonRepository.NewS3Object(tx.WriterDB()).Upsert(ctx, &imgObj); err != nil {
			return fmt.Errorf("failed to upsert s3 object: %w", err)
		}
		email := "test@sampay.link"
		user := userModel.User{
			Slug:  "mickamy",
			Email: email,
			Attribute: userModel.UserAttribute{
				UsageCategoryType: "personal",
			},
			Profile: userModel.UserProfile{
				Name: "Tetsuro Mikami",
				Bio: ptr.Of(`Hello, I'm Tetsuro Mikami.
I'm a software engineer with over 10 years of experience in developing scalable web applications.
I specialize in Go, SQL, and cloud infrastructure.
In my free time, I enjoy contributing to open-source projects and exploring new technologies.`),
				ImageID: &imgObj.ID,
			},
		}

		if err := tx.WriterDB().WithContext(ctx).FirstOrCreate(&user, "slug = ?", user.Slug).Error; err != nil {
			return fmt.Errorf("failed to upsert user: %w", err)
		}

		auth := authFixture.AuthenticationEmailPassword(func(m *authModel.Authentication) {
			m.UserID = user.ID
			m.Identifier = email
		})

		if err := tx.WriterDB().WithContext(ctx).FirstOrCreate(&auth, "user_id = ? AND type = ? AND identifier = ?", auth.UserID, auth.Type, auth.Identifier).Error; err != nil {
			return fmt.Errorf("failed to create authentication: %w", err)
		}

		links := []userModel.UserLink{
			{
				UserID:       user.ID,
				ProviderType: userModel.UserLinkProviderTypeKyash,
				URI:          "kyash://qr/u/8617830531998519755",
				DisplayAttribute: userModel.UserLinkDisplayAttribute{
					Name:         "Kyash",
					DisplayOrder: 1,
				},
			},
			{
				UserID:       user.ID,
				ProviderType: userModel.UserLinkProviderTypePayPay,
				URI:          "https://qr.paypay.ne.jp/p2p01_saRwMniDTpCabvPb",
				DisplayAttribute: userModel.UserLinkDisplayAttribute{
					Name:         "PayPay",
					DisplayOrder: 2,
				},
			},
			{
				UserID:       user.ID,
				ProviderType: userModel.UserLinkProviderTypeAmazon,
				URI:          "https://www.amazon.jp/hz/wishlist/ls/3QLIZ5HDJSCCC?ref_=wl_share",
				DisplayAttribute: userModel.UserLinkDisplayAttribute{
					Name:         "Amazon",
					DisplayOrder: 3,
				},
			},
			{
				UserID:       user.ID,
				ProviderType: userModel.UserLinkProviderTypeOther,
				URI:          "https://mickamy.com/",
				DisplayAttribute: userModel.UserLinkDisplayAttribute{
					Name:         "mickamy.com",
					DisplayOrder: 4,
				},
			},
		}

		for _, link := range links {
			if err := tx.WriterDB().WithContext(ctx).FirstOrCreate(&link, "uri = ?", link.URI).Error; err != nil {
				return fmt.Errorf("failed to create user link: %w", err)
			}
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
