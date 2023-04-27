// Code generated by cmd/lexgen (see Makefile's lexgen); DO NOT EDIT.

package bsky

// schema: app.bsky.actor.defs

import (
	comatprototypes "github.com/bluesky-social/indigo/api/atproto"
)

// ActorDefs_ProfileView is a "profileView" in the app.bsky.actor.defs schema.
type ActorDefs_ProfileView struct {
	Avatar      *string                            `json:"avatar,omitempty" cborgen:"avatar,omitempty"`
	Description *string                            `json:"description,omitempty" cborgen:"description,omitempty"`
	Did         string                             `json:"did" cborgen:"did"`
	DisplayName *string                            `json:"displayName,omitempty" cborgen:"displayName,omitempty"`
	Handle      string                             `json:"handle" cborgen:"handle"`
	IndexedAt   *string                            `json:"indexedAt,omitempty" cborgen:"indexedAt,omitempty"`
	Labels      []*comatprototypes.LabelDefs_Label `json:"labels,omitempty" cborgen:"labels,omitempty"`
	Viewer      *ActorDefs_ViewerState             `json:"viewer,omitempty" cborgen:"viewer,omitempty"`
}

// ActorDefs_ProfileViewBasic is a "profileViewBasic" in the app.bsky.actor.defs schema.
type ActorDefs_ProfileViewBasic struct {
	Avatar      *string                            `json:"avatar,omitempty" cborgen:"avatar,omitempty"`
	Did         string                             `json:"did" cborgen:"did"`
	DisplayName *string                            `json:"displayName,omitempty" cborgen:"displayName,omitempty"`
	Handle      string                             `json:"handle" cborgen:"handle"`
	Labels      []*comatprototypes.LabelDefs_Label `json:"labels,omitempty" cborgen:"labels,omitempty"`
	Viewer      *ActorDefs_ViewerState             `json:"viewer,omitempty" cborgen:"viewer,omitempty"`
}

// ActorDefs_ProfileViewDetailed is a "profileViewDetailed" in the app.bsky.actor.defs schema.
type ActorDefs_ProfileViewDetailed struct {
	Avatar         *string                            `json:"avatar,omitempty" cborgen:"avatar,omitempty"`
	Banner         *string                            `json:"banner,omitempty" cborgen:"banner,omitempty"`
	Description    *string                            `json:"description,omitempty" cborgen:"description,omitempty"`
	Did            string                             `json:"did" cborgen:"did"`
	DisplayName    *string                            `json:"displayName,omitempty" cborgen:"displayName,omitempty"`
	FollowersCount *int64                             `json:"followersCount,omitempty" cborgen:"followersCount,omitempty"`
	FollowsCount   *int64                             `json:"followsCount,omitempty" cborgen:"followsCount,omitempty"`
	Handle         string                             `json:"handle" cborgen:"handle"`
	IndexedAt      *string                            `json:"indexedAt,omitempty" cborgen:"indexedAt,omitempty"`
	Labels         []*comatprototypes.LabelDefs_Label `json:"labels,omitempty" cborgen:"labels,omitempty"`
	PostsCount     *int64                             `json:"postsCount,omitempty" cborgen:"postsCount,omitempty"`
	Viewer         *ActorDefs_ViewerState             `json:"viewer,omitempty" cborgen:"viewer,omitempty"`
}

// ActorDefs_ViewerState is a "viewerState" in the app.bsky.actor.defs schema.
type ActorDefs_ViewerState struct {
	FollowedBy *string `json:"followedBy,omitempty" cborgen:"followedBy,omitempty"`
	Following  *string `json:"following,omitempty" cborgen:"following,omitempty"`
	Muted      *bool   `json:"muted,omitempty" cborgen:"muted,omitempty"`
}
