package service

import "errors"

var (
	DuplicateUserErr           = errors.New("Duplicate user")
	UserNotFoundErr            = errors.New("User no found")
	PasswordNotMatch           = errors.New("Password doesn't match")
	DuplicateWorkspaceErr      = errors.New("Duplicate workspace")
	DuplicateWorkspaceEmailErr = errors.New("Duplicate user in workspace")
	WorkspaceNotExistsErr      = errors.New("Workspace didn`t exist")
	WorkspaceAccessDeniedErr   = errors.New("Workspace access denied")
	EmptyUrlParamsErr          = errors.New("Empty URL parameters in http address")
)
