package system

import (
	"fmt"
	"sort"

	"github.com/SimonBaeumer/goss/util"
	"github.com/opencontainers/runc/libcontainer/user"
)

type User interface {
	Username() string
	Exists() (bool, error)
	UID() (int, error)
	GID() (int, error)
	Groups() ([]string, error)
	Home() (string, error)
	Shell() (string, error)
}

// DefUser represents a user in the system package
type DefUser struct {
	username string
}

// NewDefUser is the constructor for creating a new user
func NewDefUser(username string, system *System, config util.Config) User {
	return &DefUser{username: username}
}

//Username returns the username
func (u *DefUser) Username() string {
	return u.username
}

//Exists checks if a user exists on the system
func (u *DefUser) Exists() (bool, error) {
	_, err := user.LookupUser(u.username)
	if err != nil {
		return false, nil
	}
	return true, nil
}

//UID returns the UID by username
func (u *DefUser) UID() (int, error) {
	user, err := user.LookupUser(u.username)
	if err != nil {
		return 0, err
	}

	return user.Uid, nil
}

// GID returns the group id for the user
func (u *DefUser) GID() (int, error) {
	user, err := user.LookupUser(u.username)
	if err != nil {
		return 0, err
	}

	return user.Gid, nil
}

func (u *DefUser) Home() (string, error) {
	user, err := user.LookupUser(u.username)
	if err != nil {
		return "", err
	}

	return user.Home, nil
}

func (u *DefUser) Shell() (string, error) {
	user, err := user.LookupUser(u.username)
	if err != nil {
		return "", err
	}

	return user.Shell, nil
}

func (u *DefUser) Groups() ([]string, error) {
	user, err := user.LookupUser(u.username)
	if err != nil {
		return nil, err
	}

	var groupList []string
	groups, err := lookupUserGroups(user)
	if err != nil {
		return nil, err
	}

	for _, g := range groups {
		groupList = append(groupList, g.Name)
	}

	sort.Strings(groupList)
	return groupList, nil
}

func lookupUserGroups(userS user.User) ([]user.Group, error) {
	// Get operating system-specific group reader-closer.
	group, err := user.GetGroup()
	if err != nil {
		return []user.Group{user.Group{}}, err
	}
	defer group.Close()

	groups, err := user.ParseGroupFilter(group, func(g user.Group) bool {
		// Primary group
		if g.Gid == userS.Gid {
			return true
		}

		// Check if user is a member.
		for _, u := range g.List {
			if u == userS.Name {
				return true
			}
		}

		return false
	})

	if err != nil {
		return []user.Group{user.Group{}}, fmt.Errorf("Unable to find groups for user %v: %v", userS, err)
	}

	return groups, nil
}
