// Copyright (c) 2016-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package store

import (
	"time"

	"github.com/mattermost/mattermost-server/model"
)

type StoreResult struct {
	Data interface{}
	Err  *model.AppError
}

type StoreChannel chan StoreResult

func Do(f func(result *StoreResult)) StoreChannel {
	storeChannel := make(StoreChannel, 1)
	go func() {
		result := StoreResult{}
		f(&result)
		storeChannel <- result
		close(storeChannel)
	}()
	return storeChannel
}

func Must(sc StoreChannel) interface{} {
	r := <-sc
	if r.Err != nil {

		time.Sleep(time.Second)
		panic(r.Err)
	}

	return r.Data
}

type Store interface {
	Team() TeamStore
	Channel() ChannelStore
	Post() PostStore
	User() UserStore
	Bot() BotStore
	Audit() AuditStore
	ClusterDiscovery() ClusterDiscoveryStore
	Compliance() ComplianceStore
	Session() SessionStore
	OAuth() OAuthStore
	System() SystemStore
	Webhook() WebhookStore
	Command() CommandStore
	CommandWebhook() CommandWebhookStore
	Preference() PreferenceStore
	License() LicenseStore
	Token() TokenStore
	Emoji() EmojiStore
	Status() StatusStore
	FileInfo() FileInfoStore
	Reaction() ReactionStore
	Role() RoleStore
	Scheme() SchemeStore
	Job() JobStore
	UserAccessToken() UserAccessTokenStore
	ChannelMemberHistory() ChannelMemberHistoryStore
	Plugin() PluginStore
	TermsOfService() TermsOfServiceStore
	Group() GroupStore
	UserTermsOfService() UserTermsOfServiceStore
	LinkMetadata() LinkMetadataStore
	MarkSystemRanUnitTests()
	Close()
	LockToMaster()
	UnlockFromMaster()
	DropAllTables()
	TotalMasterDbConnections() int
	TotalReadDbConnections() int
	TotalSearchDbConnections() int
}

type TeamStore interface {
	Save(team *model.Team) (*model.Team, *model.AppError)
	Update(team *model.Team) (*model.Team, *model.AppError)
	UpdateDisplayName(name string, teamId string) *model.AppError
	Get(id string) (*model.Team, *model.AppError)
	GetByName(name string) (*model.Team, *model.AppError)
	SearchByName(name string) ([]*model.Team, *model.AppError)
	SearchAll(term string) ([]*model.Team, *model.AppError)
	SearchOpen(term string) StoreChannel
	SearchPrivate(term string) ([]*model.Team, *model.AppError)
	GetAll() ([]*model.Team, *model.AppError)
	GetAllPage(offset int, limit int) ([]*model.Team, *model.AppError)
	GetAllPrivateTeamListing() StoreChannel
	GetAllPrivateTeamPageListing(offset int, limit int) ([]*model.Team, *model.AppError)
	GetAllTeamListing() StoreChannel
	GetAllTeamPageListing(offset int, limit int) ([]*model.Team, *model.AppError)
	GetTeamsByUserId(userId string) StoreChannel
	GetByInviteId(inviteId string) (*model.Team, *model.AppError)
	PermanentDelete(teamId string) *model.AppError
	AnalyticsTeamCount() (int64, *model.AppError)
	SaveMember(member *model.TeamMember, maxUsersPerTeam int) StoreChannel
	UpdateMember(member *model.TeamMember) (*model.TeamMember, *model.AppError)
	GetMember(teamId string, userId string) (*model.TeamMember, *model.AppError)
	GetMembers(teamId string, offset int, limit int, restrictions *model.ViewUsersRestrictions) ([]*model.TeamMember, *model.AppError)
	GetMembersByIds(teamId string, userIds []string, restrictions *model.ViewUsersRestrictions) ([]*model.TeamMember, *model.AppError)
	GetTotalMemberCount(teamId string) (int64, *model.AppError)
	GetActiveMemberCount(teamId string) (int64, *model.AppError)
	GetTeamsForUser(userId string) ([]*model.TeamMember, *model.AppError)
	GetTeamsForUserWithPagination(userId string, page, perPage int) ([]*model.TeamMember, *model.AppError)
	GetChannelUnreadsForAllTeams(excludeTeamId, userId string) StoreChannel
	GetChannelUnreadsForTeam(teamId, userId string) ([]*model.ChannelUnread, *model.AppError)
	RemoveMember(teamId string, userId string) StoreChannel
	RemoveAllMembersByTeam(teamId string) StoreChannel
	RemoveAllMembersByUser(userId string) StoreChannel
	UpdateLastTeamIconUpdate(teamId string, curTime int64) StoreChannel
	GetTeamsByScheme(schemeId string, offset int, limit int) StoreChannel
	MigrateTeamMembers(fromTeamId string, fromUserId string) StoreChannel
	ResetAllTeamSchemes() StoreChannel
	ClearAllCustomRoleAssignments() StoreChannel
	AnalyticsGetTeamCountForScheme(schemeId string) StoreChannel
	GetAllForExportAfter(limit int, afterId string) StoreChannel
	GetTeamMembersForExport(userId string) StoreChannel
	UserBelongsToTeams(userId string, teamIds []string) StoreChannel
	GetUserTeamIds(userId string, allowFromCache bool) StoreChannel
	InvalidateAllTeamIdsForUser(userId string)
	ClearCaches()
}

type ChannelStore interface {
	Save(channel *model.Channel, maxChannelsPerTeam int64) (*model.Channel, *model.AppError)
	CreateDirectChannel(userId string, otherUserId string) (*model.Channel, *model.AppError)
	SaveDirectChannel(channel *model.Channel, member1 *model.ChannelMember, member2 *model.ChannelMember) (*model.Channel, *model.AppError)
	Update(channel *model.Channel) (*model.Channel, *model.AppError)
	Get(id string, allowFromCache bool) (*model.Channel, *model.AppError)
	InvalidateChannel(id string)
	InvalidateChannelByName(teamId, name string)
	GetFromMaster(id string) (*model.Channel, *model.AppError)
	Delete(channelId string, time int64) *model.AppError
	Restore(channelId string, time int64) *model.AppError
	SetDeleteAt(channelId string, deleteAt int64, updateAt int64) *model.AppError
	PermanentDeleteByTeam(teamId string) StoreChannel
	PermanentDelete(channelId string) StoreChannel
	GetByName(team_id string, name string, allowFromCache bool) (*model.Channel, *model.AppError)
	GetByNames(team_id string, names []string, allowFromCache bool) ([]*model.Channel, *model.AppError)
	GetByNameIncludeDeleted(team_id string, name string, allowFromCache bool) (*model.Channel, *model.AppError)
	GetDeletedByName(team_id string, name string) (*model.Channel, *model.AppError)
	GetDeleted(team_id string, offset int, limit int) (*model.ChannelList, *model.AppError)
	GetChannels(teamId string, userId string, includeDeleted bool) (*model.ChannelList, *model.AppError)
	GetAllChannels(page, perPage int, opts ChannelSearchOpts) (*model.ChannelListWithTeamData, *model.AppError)
	GetMoreChannels(teamId string, userId string, offset int, limit int) (*model.ChannelList, *model.AppError)
	GetPublicChannelsForTeam(teamId string, offset int, limit int) (*model.ChannelList, *model.AppError)
	GetPublicChannelsByIdsForTeam(teamId string, channelIds []string) (*model.ChannelList, *model.AppError)
	GetChannelCounts(teamId string, userId string) (*model.ChannelCounts, *model.AppError)
	GetTeamChannels(teamId string) (*model.ChannelList, *model.AppError)
	GetAll(teamId string) ([]*model.Channel, *model.AppError)
	GetChannelsByIds(channelIds []string) ([]*model.Channel, *model.AppError)
	GetForPost(postId string) (*model.Channel, *model.AppError)
	SaveMember(member *model.ChannelMember) StoreChannel
	UpdateMember(member *model.ChannelMember) StoreChannel
	GetMembers(channelId string, offset, limit int) (*model.ChannelMembers, *model.AppError)
	GetMember(channelId string, userId string) (*model.ChannelMember, *model.AppError)
	GetChannelMembersTimezones(channelId string) ([]model.StringMap, *model.AppError)
	GetAllChannelMembersForUser(userId string, allowFromCache bool, includeDeleted bool) StoreChannel
	InvalidateAllChannelMembersForUser(userId string)
	IsUserInChannelUseCache(userId string, channelId string) bool
	GetAllChannelMembersNotifyPropsForChannel(channelId string, allowFromCache bool) (map[string]model.StringMap, *model.AppError)
	InvalidateCacheForChannelMembersNotifyProps(channelId string)
	GetMemberForPost(postId string, userId string) (*model.ChannelMember, *model.AppError)
	InvalidateMemberCount(channelId string)
	GetMemberCountFromCache(channelId string) int64
	GetMemberCount(channelId string, allowFromCache bool) (int64, *model.AppError)
	GetPinnedPosts(channelId string) StoreChannel
	RemoveMember(channelId string, userId string) *model.AppError
	PermanentDeleteMembersByUser(userId string) StoreChannel
	PermanentDeleteMembersByChannel(channelId string) *model.AppError
	UpdateLastViewedAt(channelIds []string, userId string) StoreChannel
	IncrementMentionCount(channelId string, userId string) StoreChannel
	AnalyticsTypeCount(teamId string, channelType string) (int64, *model.AppError)
	GetMembersForUser(teamId string, userId string) StoreChannel
	GetMembersForUserWithPagination(teamId, userId string, page, perPage int) StoreChannel
	AutocompleteInTeam(teamId string, term string, includeDeleted bool) (*model.ChannelList, *model.AppError)
	AutocompleteInTeamForSearch(teamId string, userId string, term string, includeDeleted bool) StoreChannel
	SearchAllChannels(term string, opts ChannelSearchOpts) StoreChannel
	SearchInTeam(teamId string, term string, includeDeleted bool) (*model.ChannelList, *model.AppError)
	SearchMore(userId string, teamId string, term string) (*model.ChannelList, *model.AppError)
	GetMembersByIds(channelId string, userIds []string) (*model.ChannelMembers, *model.AppError)
	AnalyticsDeletedTypeCount(teamId string, channelType string) (int64, *model.AppError)
	GetChannelUnread(channelId, userId string) (*model.ChannelUnread, *model.AppError)
	ClearCaches()
	GetChannelsByScheme(schemeId string, offset int, limit int) StoreChannel
	MigrateChannelMembers(fromChannelId string, fromUserId string) (map[string]string, *model.AppError)
	ResetAllChannelSchemes() *model.AppError
	ClearAllCustomRoleAssignments() *model.AppError
	MigratePublicChannels() error
	GetAllChannelsForExportAfter(limit int, afterId string) ([]*model.ChannelForExport, *model.AppError)
	GetAllDirectChannelsForExportAfter(limit int, afterId string) ([]*model.DirectChannelForExport, *model.AppError)
	GetChannelMembersForExport(userId string, teamId string) ([]*model.ChannelMemberForExport, *model.AppError)
	RemoveAllDeactivatedMembers(channelId string) *model.AppError
	GetChannelsBatchForIndexing(startTime, endTime int64, limit int) ([]*model.Channel, *model.AppError)
	UserBelongsToChannels(userId string, channelIds []string) (bool, *model.AppError)
}

type ChannelMemberHistoryStore interface {
	LogJoinEvent(userId string, channelId string, joinTime int64) StoreChannel
	LogLeaveEvent(userId string, channelId string, leaveTime int64) StoreChannel
	GetUsersInChannelDuring(startTime int64, endTime int64, channelId string) StoreChannel
	PermanentDeleteBatch(endTime int64, limit int64) StoreChannel
}

type PostStore interface {
	Save(post *model.Post) (*model.Post, *model.AppError)
	Update(newPost *model.Post, oldPost *model.Post) (*model.Post, *model.AppError)
	Get(id string) (*model.PostList, *model.AppError)
	GetSingle(id string) (*model.Post, *model.AppError)
	Delete(postId string, time int64, deleteByID string) *model.AppError
	PermanentDeleteByUser(userId string) *model.AppError
	PermanentDeleteByChannel(channelId string) *model.AppError
	GetPosts(channelId string, offset int, limit int, allowFromCache bool) (*model.PostList, *model.AppError)
	GetFlaggedPosts(userId string, offset int, limit int) (*model.PostList, *model.AppError)
	GetFlaggedPostsForTeam(userId, teamId string, offset int, limit int) (*model.PostList, *model.AppError)
	GetFlaggedPostsForChannel(userId, channelId string, offset int, limit int) (*model.PostList, *model.AppError)
	GetPostsBefore(channelId string, postId string, numPosts int, offset int) (*model.PostList, *model.AppError)
	GetPostsAfter(channelId string, postId string, numPosts int, offset int) (*model.PostList, *model.AppError)
	GetPostsSince(channelId string, time int64, allowFromCache bool) (*model.PostList, *model.AppError)
	GetEtag(channelId string, allowFromCache bool) string
	Search(teamId string, userId string, params *model.SearchParams) StoreChannel
	AnalyticsUserCountsWithPostsByDay(teamId string) (model.AnalyticsRows, *model.AppError)
	AnalyticsPostCountsByDay(teamId string) (model.AnalyticsRows, *model.AppError)
	AnalyticsPostCount(teamId string, mustHaveFile bool, mustHaveHashtag bool) (int64, *model.AppError)
	ClearCaches()
	InvalidateLastPostTimeCache(channelId string)
	GetPostsCreatedAt(channelId string, time int64) ([]*model.Post, *model.AppError)
	Overwrite(post *model.Post) (*model.Post, *model.AppError)
	GetPostsByIds(postIds []string) ([]*model.Post, *model.AppError)
	GetPostsBatchForIndexing(startTime int64, endTime int64, limit int) ([]*model.PostForIndexing, *model.AppError)
	PermanentDeleteBatch(endTime int64, limit int64) (int64, *model.AppError)
	GetOldest() (*model.Post, *model.AppError)
	GetMaxPostSize() int
	GetParentsForExportAfter(limit int, afterId string) ([]*model.PostForExport, *model.AppError)
	GetRepliesForExport(parentId string) ([]*model.ReplyForExport, *model.AppError)
	GetDirectPostParentsForExportAfter(limit int, afterId string) ([]*model.DirectPostForExport, *model.AppError)
}

type UserStore interface {
	Save(user *model.User) StoreChannel
	Update(user *model.User, allowRoleUpdate bool) (*model.UserUpdate, *model.AppError)
	UpdateLastPictureUpdate(userId string) StoreChannel
	ResetLastPictureUpdate(userId string) StoreChannel
	UpdateUpdateAt(userId string) StoreChannel
	UpdatePassword(userId, newPassword string) StoreChannel
	UpdateAuthData(userId string, service string, authData *string, email string, resetMfa bool) StoreChannel
	UpdateMfaSecret(userId, secret string) StoreChannel
	UpdateMfaActive(userId string, active bool) StoreChannel
	Get(id string) (*model.User, *model.AppError)
	GetAll() StoreChannel
	ClearCaches()
	InvalidateProfilesInChannelCacheByUser(userId string)
	InvalidateProfilesInChannelCache(channelId string)
	GetProfilesInChannel(channelId string, offset int, limit int) StoreChannel
	GetProfilesInChannelByStatus(channelId string, offset int, limit int) StoreChannel
	GetAllProfilesInChannel(channelId string, allowFromCache bool) StoreChannel
	GetProfilesNotInChannel(teamId string, channelId string, groupConstrained bool, offset int, limit int, viewRestrictions *model.ViewUsersRestrictions) StoreChannel
	GetProfilesWithoutTeam(offset int, limit int, viewRestrictions *model.ViewUsersRestrictions) StoreChannel
	GetProfilesByUsernames(usernames []string, viewRestrictions *model.ViewUsersRestrictions) StoreChannel
	GetAllProfiles(options *model.UserGetOptions) StoreChannel
	GetProfiles(options *model.UserGetOptions) StoreChannel
	GetProfileByIds(userId []string, allowFromCache bool, viewRestrictions *model.ViewUsersRestrictions) StoreChannel
	InvalidatProfileCacheForUser(userId string)
	GetByEmail(email string) (*model.User, *model.AppError)
	GetByAuth(authData *string, authService string) (*model.User, *model.AppError)
	GetAllUsingAuthService(authService string) StoreChannel
	GetByUsername(username string) StoreChannel
	GetForLogin(loginId string, allowSignInWithUsername, allowSignInWithEmail bool) StoreChannel
	VerifyEmail(userId, email string) (string, *model.AppError)
	GetEtagForAllProfiles() StoreChannel
	GetEtagForProfiles(teamId string) StoreChannel
	UpdateFailedPasswordAttempts(userId string, attempts int) StoreChannel
	GetSystemAdminProfiles() StoreChannel
	PermanentDelete(userId string) *model.AppError
	AnalyticsActiveCount(time int64) StoreChannel
	GetUnreadCount(userId string) StoreChannel
	GetUnreadCountForChannel(userId string, channelId string) StoreChannel
	GetAnyUnreadPostCountForChannel(userId string, channelId string) StoreChannel
	GetRecentlyActiveUsersForTeam(teamId string, offset, limit int, viewRestrictions *model.ViewUsersRestrictions) StoreChannel
	GetNewUsersForTeam(teamId string, offset, limit int, viewRestrictions *model.ViewUsersRestrictions) StoreChannel
	Search(teamId string, term string, options *model.UserSearchOptions) StoreChannel
	SearchNotInTeam(notInTeamId string, term string, options *model.UserSearchOptions) StoreChannel
	SearchInChannel(channelId string, term string, options *model.UserSearchOptions) StoreChannel
	SearchNotInChannel(teamId string, channelId string, term string, options *model.UserSearchOptions) StoreChannel
	SearchWithoutTeam(term string, options *model.UserSearchOptions) StoreChannel
	AnalyticsGetInactiveUsersCount() StoreChannel
	AnalyticsGetSystemAdminCount() StoreChannel
	GetProfilesNotInTeam(teamId string, groupConstrained bool, offset int, limit int, viewRestrictions *model.ViewUsersRestrictions) StoreChannel
	GetEtagForProfilesNotInTeam(teamId string) StoreChannel
	ClearAllCustomRoleAssignments() StoreChannel
	InferSystemInstallDate() StoreChannel
	GetAllAfter(limit int, afterId string) StoreChannel
	GetUsersBatchForIndexing(startTime, endTime int64, limit int) StoreChannel
	Count(options model.UserCountOptions) StoreChannel
	GetTeamGroupUsers(teamID string) StoreChannel
	GetChannelGroupUsers(channelID string) StoreChannel
}

type BotStore interface {
	Get(userId string, includeDeleted bool) (*model.Bot, *model.AppError)
	GetAll(options *model.BotGetOptions) ([]*model.Bot, *model.AppError)
	Save(bot *model.Bot) (*model.Bot, *model.AppError)
	Update(bot *model.Bot) (*model.Bot, *model.AppError)
	PermanentDelete(userId string) *model.AppError
}

type SessionStore interface {
	Get(sessionIdOrToken string) (*model.Session, *model.AppError)
	Save(session *model.Session) (*model.Session, *model.AppError)
	GetSessions(userId string) ([]*model.Session, *model.AppError)
	GetSessionsWithActiveDeviceIds(userId string) ([]*model.Session, *model.AppError)
	Remove(sessionIdOrToken string) *model.AppError
	RemoveAllSessions() *model.AppError
	PermanentDeleteSessionsByUser(teamId string) *model.AppError
	UpdateLastActivityAt(sessionId string, time int64) *model.AppError
	UpdateRoles(userId string, roles string) (string, *model.AppError)
	UpdateDeviceId(id string, deviceId string, expiresAt int64) (string, *model.AppError)
	AnalyticsSessionCount() (int64, *model.AppError)
	Cleanup(expiryTime int64, batchSize int64)
}

type AuditStore interface {
	Save(audit *model.Audit) *model.AppError
	Get(user_id string, offset int, limit int) (model.Audits, *model.AppError)
	PermanentDeleteByUser(userId string) *model.AppError
	PermanentDeleteBatch(endTime int64, limit int64) (int64, *model.AppError)
}

type ClusterDiscoveryStore interface {
	Save(discovery *model.ClusterDiscovery) *model.AppError
	Delete(discovery *model.ClusterDiscovery) (bool, *model.AppError)
	Exists(discovery *model.ClusterDiscovery) (bool, *model.AppError)
	GetAll(discoveryType, clusterName string) ([]*model.ClusterDiscovery, *model.AppError)
	SetLastPingAt(discovery *model.ClusterDiscovery) *model.AppError
	Cleanup() *model.AppError
}

type ComplianceStore interface {
	Save(compliance *model.Compliance) (*model.Compliance, *model.AppError)
	Update(compliance *model.Compliance) (*model.Compliance, *model.AppError)
	Get(id string) (*model.Compliance, *model.AppError)
	GetAll(offset, limit int) (model.Compliances, *model.AppError)
	ComplianceExport(compliance *model.Compliance) ([]*model.CompliancePost, *model.AppError)
	MessageExport(after int64, limit int) ([]*model.MessageExport, *model.AppError)
}

type OAuthStore interface {
	SaveApp(app *model.OAuthApp) (*model.OAuthApp, *model.AppError)
	UpdateApp(app *model.OAuthApp) (*model.OAuthApp, *model.AppError)
	GetApp(id string) (*model.OAuthApp, *model.AppError)
	GetAppByUser(userId string, offset, limit int) ([]*model.OAuthApp, *model.AppError)
	GetApps(offset, limit int) ([]*model.OAuthApp, *model.AppError)
	GetAuthorizedApps(userId string, offset, limit int) ([]*model.OAuthApp, *model.AppError)
	DeleteApp(id string) *model.AppError
	SaveAuthData(authData *model.AuthData) (*model.AuthData, *model.AppError)
	GetAuthData(code string) (*model.AuthData, *model.AppError)
	RemoveAuthData(code string) *model.AppError
	PermanentDeleteAuthDataByUser(userId string) *model.AppError
	SaveAccessData(accessData *model.AccessData) (*model.AccessData, *model.AppError)
	UpdateAccessData(accessData *model.AccessData) (*model.AccessData, *model.AppError)
	GetAccessData(token string) (*model.AccessData, *model.AppError)
	GetAccessDataByUserForApp(userId, clientId string) ([]*model.AccessData, *model.AppError)
	GetAccessDataByRefreshToken(token string) (*model.AccessData, *model.AppError)
	GetPreviousAccessData(userId, clientId string) (*model.AccessData, *model.AppError)
	RemoveAccessData(token string) *model.AppError
}

type SystemStore interface {
	Save(system *model.System) *model.AppError
	SaveOrUpdate(system *model.System) *model.AppError
	Update(system *model.System) *model.AppError
	Get() (model.StringMap, *model.AppError)
	GetByName(name string) (*model.System, *model.AppError)
	PermanentDeleteByName(name string) (*model.System, *model.AppError)
}

type WebhookStore interface {
	SaveIncoming(webhook *model.IncomingWebhook) (*model.IncomingWebhook, *model.AppError)
	GetIncoming(id string, allowFromCache bool) (*model.IncomingWebhook, *model.AppError)
	GetIncomingList(offset, limit int) ([]*model.IncomingWebhook, *model.AppError)
	GetIncomingByTeam(teamId string, offset, limit int) ([]*model.IncomingWebhook, *model.AppError)
	UpdateIncoming(webhook *model.IncomingWebhook) (*model.IncomingWebhook, *model.AppError)
	GetIncomingByChannel(channelId string) ([]*model.IncomingWebhook, *model.AppError)
	DeleteIncoming(webhookId string, time int64) *model.AppError
	PermanentDeleteIncomingByChannel(channelId string) *model.AppError
	PermanentDeleteIncomingByUser(userId string) *model.AppError

	SaveOutgoing(webhook *model.OutgoingWebhook) (*model.OutgoingWebhook, *model.AppError)
	GetOutgoing(id string) (*model.OutgoingWebhook, *model.AppError)
	GetOutgoingByChannel(channelId string, offset, limit int) ([]*model.OutgoingWebhook, *model.AppError)
	GetOutgoingList(offset, limit int) ([]*model.OutgoingWebhook, *model.AppError)
	GetOutgoingByTeam(teamId string, offset, limit int) ([]*model.OutgoingWebhook, *model.AppError)
	DeleteOutgoing(webhookId string, time int64) *model.AppError
	PermanentDeleteOutgoingByChannel(channelId string) *model.AppError
	PermanentDeleteOutgoingByUser(userId string) *model.AppError
	UpdateOutgoing(hook *model.OutgoingWebhook) (*model.OutgoingWebhook, *model.AppError)

	AnalyticsIncomingCount(teamId string) (int64, *model.AppError)
	AnalyticsOutgoingCount(teamId string) (int64, *model.AppError)
	InvalidateWebhookCache(webhook string)
	ClearCaches()
}

type CommandStore interface {
	Save(webhook *model.Command) (*model.Command, *model.AppError)
	GetByTrigger(teamId string, trigger string) (*model.Command, *model.AppError)
	Get(id string) (*model.Command, *model.AppError)
	GetByTeam(teamId string) ([]*model.Command, *model.AppError)
	Delete(commandId string, time int64) *model.AppError
	PermanentDeleteByTeam(teamId string) *model.AppError
	PermanentDeleteByUser(userId string) *model.AppError
	Update(hook *model.Command) (*model.Command, *model.AppError)
	AnalyticsCommandCount(teamId string) (int64, *model.AppError)
}

type CommandWebhookStore interface {
	Save(webhook *model.CommandWebhook) StoreChannel
	Get(id string) StoreChannel
	TryUse(id string, limit int) StoreChannel
	Cleanup()
}

type PreferenceStore interface {
	Save(preferences *model.Preferences) *model.AppError
	GetCategory(userId string, category string) (model.Preferences, *model.AppError)
	Get(userId string, category string, name string) (*model.Preference, *model.AppError)
	GetAll(userId string) (model.Preferences, *model.AppError)
	Delete(userId, category, name string) *model.AppError
	DeleteCategory(userId string, category string) *model.AppError
	DeleteCategoryAndName(category string, name string) *model.AppError
	PermanentDeleteByUser(userId string) *model.AppError
	IsFeatureEnabled(feature, userId string) (bool, *model.AppError)
	CleanupFlagsBatch(limit int64) (int64, *model.AppError)
}

type LicenseStore interface {
	Save(license *model.LicenseRecord) (*model.LicenseRecord, *model.AppError)
	Get(id string) (*model.LicenseRecord, *model.AppError)
}

type TokenStore interface {
	Save(recovery *model.Token) StoreChannel
	Delete(token string) StoreChannel
	GetByToken(token string) (*model.Token, *model.AppError)
	Cleanup()
	RemoveAllTokensByType(tokenType string) StoreChannel
}

type EmojiStore interface {
	Save(emoji *model.Emoji) (*model.Emoji, *model.AppError)
	Get(id string, allowFromCache bool) (*model.Emoji, *model.AppError)
	GetByName(name string) (*model.Emoji, *model.AppError)
	GetMultipleByName(names []string) StoreChannel
	GetList(offset, limit int, sort string) StoreChannel
	Delete(id string, time int64) *model.AppError
	Search(name string, prefixOnly bool, limit int) StoreChannel
}

type StatusStore interface {
	SaveOrUpdate(status *model.Status) StoreChannel
	Get(userId string) StoreChannel
	GetByIds(userIds []string) StoreChannel
	GetOnlineAway() StoreChannel
	GetOnline() StoreChannel
	GetAllFromTeam(teamId string) StoreChannel
	ResetAll() StoreChannel
	GetTotalActiveUsersCount() StoreChannel
	UpdateLastActivityAt(userId string, lastActivityAt int64) StoreChannel
}

type FileInfoStore interface {
	Save(info *model.FileInfo) (*model.FileInfo, *model.AppError)
	Get(id string) (*model.FileInfo, *model.AppError)
	GetByPath(path string) (*model.FileInfo, *model.AppError)
	GetForPost(postId string, readFromMaster bool, allowFromCache bool) ([]*model.FileInfo, *model.AppError)
	GetForUser(userId string) ([]*model.FileInfo, *model.AppError)
	InvalidateFileInfosForPostCache(postId string)
	AttachToPost(fileId string, postId string, creatorId string) *model.AppError
	DeleteForPost(postId string) (string, *model.AppError)
	PermanentDelete(fileId string) *model.AppError
	PermanentDeleteBatch(endTime int64, limit int64) (int64, *model.AppError)
	PermanentDeleteByUser(userId string) (int64, *model.AppError)
	ClearCaches()
}

type ReactionStore interface {
	Save(reaction *model.Reaction) (*model.Reaction, *model.AppError)
	Delete(reaction *model.Reaction) (*model.Reaction, *model.AppError)
	GetForPost(postId string, allowFromCache bool) ([]*model.Reaction, *model.AppError)
	DeleteAllWithEmojiName(emojiName string) *model.AppError
	PermanentDeleteBatch(endTime int64, limit int64) (int64, *model.AppError)
	BulkGetForPosts(postIds []string) ([]*model.Reaction, *model.AppError)
}

type JobStore interface {
	Save(job *model.Job) (*model.Job, *model.AppError)
	UpdateOptimistically(job *model.Job, currentStatus string) (bool, *model.AppError)
	UpdateStatus(id string, status string) (*model.Job, *model.AppError)
	UpdateStatusOptimistically(id string, currentStatus string, newStatus string) (bool, *model.AppError)
	Get(id string) (*model.Job, *model.AppError)
	GetAllPage(offset int, limit int) ([]*model.Job, *model.AppError)
	GetAllByType(jobType string) ([]*model.Job, *model.AppError)
	GetAllByTypePage(jobType string, offset int, limit int) ([]*model.Job, *model.AppError)
	GetAllByStatus(status string) ([]*model.Job, *model.AppError)
	GetNewestJobByStatusAndType(status string, jobType string) (*model.Job, *model.AppError)
	GetCountByStatusAndType(status string, jobType string) (int64, *model.AppError)
	Delete(id string) (string, *model.AppError)
}

type UserAccessTokenStore interface {
	Save(token *model.UserAccessToken) (*model.UserAccessToken, *model.AppError)
	Delete(tokenId string) StoreChannel
	DeleteAllForUser(userId string) StoreChannel
	Get(tokenId string) (*model.UserAccessToken, *model.AppError)
	GetAll(offset int, limit int) ([]*model.UserAccessToken, *model.AppError)
	GetByToken(tokenString string) (*model.UserAccessToken, *model.AppError)
	GetByUser(userId string, page, perPage int) ([]*model.UserAccessToken, *model.AppError)
	Search(term string) StoreChannel
	UpdateTokenEnable(tokenId string) StoreChannel
	UpdateTokenDisable(tokenId string) *model.AppError
}

type PluginStore interface {
	SaveOrUpdate(keyVal *model.PluginKeyValue) StoreChannel
	CompareAndSet(keyVal *model.PluginKeyValue, oldValue []byte) (bool, *model.AppError)
	Get(pluginId, key string) StoreChannel
	Delete(pluginId, key string) StoreChannel
	DeleteAllForPlugin(PluginId string) StoreChannel
	DeleteAllExpired() StoreChannel
	List(pluginId string, page, perPage int) StoreChannel
}

type RoleStore interface {
	Save(role *model.Role) (*model.Role, *model.AppError)
	Get(roleId string) (*model.Role, *model.AppError)
	GetAll() ([]*model.Role, *model.AppError)
	GetByName(name string) (*model.Role, *model.AppError)
	GetByNames(names []string) ([]*model.Role, *model.AppError)
	Delete(roldId string) (*model.Role, *model.AppError)
	PermanentDeleteAll() *model.AppError
}

type SchemeStore interface {
	Save(scheme *model.Scheme) StoreChannel
	Get(schemeId string) StoreChannel
	GetByName(schemeName string) StoreChannel
	GetAllPage(scope string, offset int, limit int) StoreChannel
	Delete(schemeId string) StoreChannel
	PermanentDeleteAll() StoreChannel
}

type TermsOfServiceStore interface {
	Save(termsOfService *model.TermsOfService) StoreChannel
	GetLatest(allowFromCache bool) StoreChannel
	Get(id string, allowFromCache bool) StoreChannel
}

type UserTermsOfServiceStore interface {
	GetByUser(userId string) StoreChannel
	Save(userTermsOfService *model.UserTermsOfService) StoreChannel
	Delete(userId, termsOfServiceId string) StoreChannel
}

type GroupStore interface {
	Create(group *model.Group) StoreChannel
	Get(groupID string) StoreChannel
	GetByIDs(groupIDs []string) ([]*model.Group, *model.AppError)
	GetByRemoteID(remoteID string, groupSource model.GroupSource) StoreChannel
	GetAllBySource(groupSource model.GroupSource) StoreChannel
	Update(group *model.Group) StoreChannel
	Delete(groupID string) StoreChannel

	GetMemberUsers(groupID string) StoreChannel
	GetMemberUsersPage(groupID string, offset int, limit int) StoreChannel
	GetMemberCount(groupID string) StoreChannel
	CreateOrRestoreMember(groupID string, userID string) StoreChannel
	DeleteMember(groupID string, userID string) StoreChannel

	CreateGroupSyncable(groupSyncable *model.GroupSyncable) (*model.GroupSyncable, *model.AppError)
	GetGroupSyncable(groupID string, syncableID string, syncableType model.GroupSyncableType) (*model.GroupSyncable, *model.AppError)
	GetAllGroupSyncablesByGroupId(groupID string, syncableType model.GroupSyncableType) ([]*model.GroupSyncable, *model.AppError)
	UpdateGroupSyncable(groupSyncable *model.GroupSyncable) (*model.GroupSyncable, *model.AppError)
	DeleteGroupSyncable(groupID string, syncableID string, syncableType model.GroupSyncableType) (*model.GroupSyncable, *model.AppError)

	TeamMembersToAdd(since int64) ([]*model.UserTeamIDPair, *model.AppError)
	ChannelMembersToAdd(since int64) ([]*model.UserChannelIDPair, *model.AppError)

	TeamMembersToRemove() ([]*model.TeamMember, *model.AppError)
	ChannelMembersToRemove() ([]*model.ChannelMember, *model.AppError)

	GetGroupsByChannel(channelId string, opts model.GroupSearchOpts) ([]*model.Group, *model.AppError)
	CountGroupsByChannel(channelId string, opts model.GroupSearchOpts) (int64, *model.AppError)

	GetGroupsByTeam(teamId string, opts model.GroupSearchOpts) ([]*model.Group, *model.AppError)
	CountGroupsByTeam(teamId string, opts model.GroupSearchOpts) (int64, *model.AppError)

	GetGroups(page, perPage int, opts model.GroupSearchOpts) ([]*model.Group, *model.AppError)

	TeamMembersMinusGroupMembers(teamID string, groupIDs []string, page, perPage int) ([]*model.UserWithGroups, *model.AppError)
	CountTeamMembersMinusGroupMembers(teamID string, groupIDs []string) (int64, *model.AppError)
	ChannelMembersMinusGroupMembers(channelID string, groupIDs []string, page, perPage int) ([]*model.UserWithGroups, *model.AppError)
	CountChannelMembersMinusGroupMembers(channelID string, groupIDs []string) (int64, *model.AppError)
}

type LinkMetadataStore interface {
	Save(linkMetadata *model.LinkMetadata) StoreChannel
	Get(url string, timestamp int64) StoreChannel
}

// ChannelSearchOpts contains options for searching channels.
//
// NotAssociatedToGroup will exclude channels that have associated, active GroupChannels records.
// IncludeDeleted will include channel records where DeleteAt != 0.
// ExcludeChannelNames will exclude channels from the results by name.
//
type ChannelSearchOpts struct {
	NotAssociatedToGroup string
	IncludeDeleted       bool
	ExcludeChannelNames  []string
}
