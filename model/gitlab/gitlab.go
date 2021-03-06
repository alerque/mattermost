// Copyright (c) 2015 Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package oauthgitlab

import (
	"encoding/json"
	"github.com/mattermost/platform/einterfaces"
	"github.com/mattermost/platform/model"
	"io"
	"strconv"
	"strings"
)

const (
	USER_AUTH_SERVICE_GITLAB = "gitlab"
)

type GitLabProvider struct {
}

type GitLabUser struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Login    string `json:"login"`
	Email    string `json:"email"`
	Name     string `json:"name"`
}

func init() {
	provider := &GitLabProvider{}
	einterfaces.RegisterOauthProvider(USER_AUTH_SERVICE_GITLAB, provider)
}

func userFromGitLabUser(glu *GitLabUser) *model.User {
	user := &model.User{}
	username := glu.Username
	if username == "" {
		username = glu.Login
	}
	user.Username = model.CleanUsername(username)
	splitName := strings.Split(glu.Name, " ")
	if len(splitName) == 2 {
		user.FirstName = splitName[0]
		user.LastName = splitName[1]
	} else if len(splitName) >= 2 {
		user.FirstName = splitName[0]
		user.LastName = strings.Join(splitName[1:], " ")
	} else {
		user.FirstName = glu.Name
	}
	user.Email = glu.Email
	user.AuthData = strconv.FormatInt(glu.Id, 10)
	user.AuthService = USER_AUTH_SERVICE_GITLAB

	return user
}

func gitLabUserFromJson(data io.Reader) *GitLabUser {
	decoder := json.NewDecoder(data)
	var glu GitLabUser
	err := decoder.Decode(&glu)
	if err == nil {
		return &glu
	} else {
		return nil
	}
}

func (glu *GitLabUser) IsValid() bool {
	if glu.Id == 0 {
		return false
	}

	if len(glu.Email) == 0 {
		return false
	}

	return true
}

func (glu *GitLabUser) getAuthData() string {
	return strconv.FormatInt(glu.Id, 10)
}

func (m *GitLabProvider) GetIdentifier() string {
	return USER_AUTH_SERVICE_GITLAB
}

func (m *GitLabProvider) GetUserFromJson(data io.Reader) *model.User {
	glu := gitLabUserFromJson(data)
	if glu.IsValid() {
		return userFromGitLabUser(glu)
	}

	return &model.User{}
}

func (m *GitLabProvider) GetAuthDataFromJson(data io.Reader) string {
	glu := gitLabUserFromJson(data)

	if glu.IsValid() {
		return glu.getAuthData()
	}

	return ""
}
