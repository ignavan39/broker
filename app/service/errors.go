package service

import "errors"

var (
	DuplicateUserErr            = errors.New("Duplicate user")
	UserNotFoundErr             = errors.New("User no found")
	PasswordNotMatch            = errors.New("Password doesn't match")
	DuplicateWorkspaceErr       = errors.New("Duplicate workspace")
	DuplicateWorkspaceEmailErr  = errors.New("Duplicate user in workspace")
	WorkspaceAccessDeniedErr    = errors.New("Workspace access denied")
	EmptyUrlParamsErr           = errors.New("Empty URL parameters in http address")
	DuplicateWorkspaceAccessErr = errors.New("Duplicate workspace access")
	DuplicateInvitationErr      = errors.New("Duplicate invitation to this workspace for user")
)
