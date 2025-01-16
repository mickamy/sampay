package i18n
// Code generated by generate.go; DO NOT EDIT.

type MessageID string

func (id MessageID) String() string { return string(id) }

const (
	RegistrationUsecaseCommonErrorEmail_already_exists MessageID = "registration.usecase.common.error.email_already_exists"
	RegistrationEmailRequest_email_verificationHeader MessageID = "registration.email.request_email_verification.header"
	RegistrationEmailRequest_email_verificationBody MessageID = "registration.email.request_email_verification.body"
	UserModelUser_attributeErrorDuplicated MessageID = "user.model.user_attribute.error.duplicated"
	UserHandlerUser_linkErrorInvalid_provider_type MessageID = "user.handler.user_link.error.invalid_provider_type"
	UserUsecaseGet_userErrorNot_found MessageID = "user.usecase.get_user.error.not_found"
	AuthUsecaseErrorInvalid_email_password MessageID = "auth.usecase.error.invalid_email_password"
	AuthUsecaseErrorInvalid_refresh_token MessageID = "auth.usecase.error.invalid_refresh_token"
	AuthUsecaseErrorInvalid_access_refresh_token MessageID = "auth.usecase.error.invalid_access_refresh_token"
	CommonHandlerDirect_upload_urlErrorInvalid_s3_object MessageID = "common.handler.direct_upload_url.error.invalid_s3_object"
	CommonHandlerErrorInternal MessageID = "common.handler.error.internal"
)
