package boundary

// Ref: https://developer.hashicorp.com/boundary/docs/api-clients/go-sdk

import (
	"context"
	"fmt"
	"github.com/hashicorp/boundary/api"
	"github.com/hashicorp/boundary/api/accounts"
	"github.com/hashicorp/boundary/api/authmethods"
	"github.com/hashicorp/boundary/api/groups"
	"github.com/hashicorp/boundary/api/users"
)

const ManagedGroupDescription = "This group is being managed by automation following authoritative approach from Gsuite." +
	"Any changes will be overwritten periodically."

type Boundary struct {
	Ctx context.Context

	//
	Client *api.Client

	//
	groupsClient   *groups.Client
	usersClient    *users.Client
	accountsClient *accounts.Client

	// General
	Address string

	// ScopeId represent the ID of the Boundary scope containing 'groups' and 'users' to keep in sync
	ScopeId string

	// AuthMethodOidcId represent the ID of the OIDC auth method used to retrieve users accounts
	// to be compared against G. Workspace
	AuthMethodOidcId string

	// Login related
	AuthMethodPasswordId string
	Username             string
	Password             string
}

// InitBoundary TODO
// Ref: https://github.com/hashicorp/boundary/blob/main/api/
func (a *Boundary) InitBoundary() (err error) {

	// Create a client from the Boundary API and set the address to reach Boundary
	config := api.Config{Addr: a.Address}
	a.Client, err = api.NewClient(&config)
	if err != nil {
		return err
	}

	// Create an auth method client
	amClient := authmethods.NewClient(a.Client)

	authenticationResult, err := amClient.Authenticate(a.Ctx,
		a.AuthMethodPasswordId,
		"login",
		map[string]interface{}{
			"login_name": a.Username,
			"password":   a.Password,
		})

	if err != nil {
		return err
	}

	// Update the original client with the token we got from the Authenticate() call
	a.Client.SetToken(fmt.Sprint(authenticationResult.Attributes["token"]))

	//
	a.groupsClient = groups.NewClient(a.Client)
	a.usersClient = users.NewClient(a.Client)
	a.accountsClient = accounts.NewClient(a.Client)

	return err
}

// GetGroups TODO
func (a *Boundary) GetGroups() (groupMap map[string]*Group, err error) {

	groupMap = make(map[string]*Group)

	groupsResult, err := a.groupsClient.List(a.Ctx, "global")
	if err != nil {
		return groupMap, err
	}

	for _, group := range groupsResult.Items {

		// Read group members
		groupContent, err := a.groupsClient.Read(a.Ctx, group.Id)
		if err != nil {
			return groupMap, err
		}

		groupMap[group.Name] = &Group{
			Id:      group.Id,
			Name:    group.Name,
			Version: group.Version,
			Members: groupContent.Item.MemberIds,
		}
	}

	return groupMap, err
}

// CreateGroup TODO
func (a *Boundary) CreateGroup(name string) (groupCreationResult *groups.GroupCreateResult, err error) {

	var options []groups.Option

	options = append(options, groups.WithName(name))
	options = append(options, groups.WithDescription(ManagedGroupDescription))

	groupCreationResult, err = a.groupsClient.Create(a.Ctx, "global", options...)

	return groupCreationResult, err
}

// SetGroupMembers TODO
func (a *Boundary) SetGroupMembers(groupId string, groupVersion uint32, userIds []string) (err error) {

	_, err = a.groupsClient.SetMembers(a.Ctx, groupId, groupVersion, userIds)

	return err
}

// GetUsers TODO
func (a *Boundary) GetUsers(personMap *map[string]*Person) (err error) {

	usersResult, err := a.usersClient.List(a.Ctx, a.ScopeId)
	if err != nil {
		return err
	}

	for _, user := range usersResult.Items {
		if user.Email == "" {
			continue
		}

		if (*personMap)[user.Email] == nil {
			(*personMap)[user.Email] = &Person{}
		}

		(*personMap)[user.Email].UserId = user.Id
		(*personMap)[user.Email].Email = fmt.Sprintf("%v", user.Email)

	}

	return err
}

// GetAccounts TODO
func (a *Boundary) GetAccounts(personMap *map[string]*Person) (err error) {

	accountsResult, err := a.accountsClient.List(a.Ctx, a.AuthMethodOidcId)
	if err != nil {
		return err
	}

	for _, account := range accountsResult.Items {
		emailString := account.Attributes["email"].(string)

		if emailString == "" {
			continue
		}

		if (*personMap)[emailString] == nil {
			(*personMap)[emailString] = &Person{}
		}

		(*personMap)[emailString].Email = account.Attributes["email"].(string)
		(*personMap)[emailString].AccountId = account.Id
		(*personMap)[emailString].Subject = account.Attributes["subject"].(string)
	}

	return err
}

// GetPersonMap TODO
func (a *Boundary) GetPersonMap() (personMap map[string]*Person, err error) {

	personMap = make(map[string]*Person)

	err = a.GetAccounts(&personMap)
	if err != nil {
		return personMap, err
	}
	err = a.GetUsers(&personMap)

	return personMap, err
}
