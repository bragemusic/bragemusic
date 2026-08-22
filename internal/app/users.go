package app

import (
	"os"

	"github.com/bragemusic/bragemusic/pkg/types"
	"github.com/gofrs/uuid/v5"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) CreateUser(email, username, password string, roles []string) {
	proles := []types.UserRole{}
	for _, r := range roles {
		proles = append(proles, types.UserRole(r))
	}

	err := a.client.CreateUser(a.ctx, email, username, password, proles)
	if err != nil {
		a.handleError(err)
		return
	}
}

func (a *App) DeleteUser(id string) {
	uid, err := uuid.FromString(id)
	if err != nil {
		a.handleError(err)
		return
	}

	err = a.client.DeleteUser(a.ctx, uid)
	if err != nil {
		a.handleError(err)
		return
	}
}

func (a *App) ListUsers(includeMachineUsers bool) []types.UserDetails {
	users, err := a.client.ListUsers(a.ctx, includeMachineUsers)
	if err != nil {
		a.handleError(err)
		return []types.UserDetails{}
	}

	return users
}

func (a *App) ListUserRoles() []types.UserRole {
	users, err := a.client.ListUserRoles(a.ctx)
	if err != nil {
		a.handleError(err)
		return []types.UserRole{}
	}

	return users
}

func (a *App) UpdateUser(id, email, username string, password *string, roles []string) {
	proles := []types.UserRole{}
	for _, r := range roles {
		proles = append(proles, types.UserRole(r))
	}

	uid, err := uuid.FromString(id)
	if err != nil {
		a.handleError(err)
		return
	}

	err = a.client.UpdateUser(a.ctx, uid, email, username, password, proles)
	if err != nil {
		a.handleError(err)
		return
	}
}

func (a *App) UpdateProfile(email, username, password, newPassword, NewPasswordConfirm *string) {
	data := types.UpdateProfileReq{
		Email:              email,
		Username:           username,
		Password:           password,
		NewPassword:        newPassword,
		NewPasswordConfirm: NewPasswordConfirm,
	}

	err := a.client.UpdateProfile(a.ctx, data)
	if err != nil {
		a.handleError(err)
		return
	}
}

func (a *App) UploadUserImage() {
	filename, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select profile image",
		Filters: []runtime.FileFilter{{
			DisplayName: "Image",
			Pattern:     "*.jpg",
		}},
	})
	if err != nil {
		a.handleError(err)
		return
	}

	if filename == "" {
		return
	}

	r, err := os.Open(filename)
	if err != nil {
		a.handleError(err)
		return
	}
	defer r.Close()

	err = a.client.UploadUserImage(a.ctx, r, filename)
	if err != nil {
		a.handleError(err)
		return
	}
}
