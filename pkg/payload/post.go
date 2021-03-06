package payload

import (
	"time"

	"github.com/neoxelox/odin/internal/class"
)

type Post struct {
	ID           string    `json:"id"`
	ThreadID     *string   `json:"thread_id"`
	CreatorID    string    `json:"creator_id"`
	Type         string    `json:"type"`
	Priority     *int      `json:"priority"`
	RecipientIDs *[]string `json:"recipient_ids"`
	VoterIDs     []string  `json:"voter_ids"`
	Subposts     int       `json:"subposts"`
	CreatedAt    time.Time `json:"created_at"`
	PostHistory
}

type PostWidgets struct {
	Poll *map[string][]string `json:"poll,omitempty"`
}

type PostHistory struct {
	UpdatorID  string      `json:"updator_id"`
	Message    string      `json:"message"`
	Categories []string    `json:"categories"`
	State      *string     `json:"state"`
	Media      []string    `json:"media"`
	Widgets    PostWidgets `json:"widgets"`
	CreatedAt  time.Time   `json:"created_at,omitempty"`
}

type GetPostRequest struct {
	class.Payload
	CommunityID string `param:"community_id" validate:"required"`
	PostID      string `param:"post_id" validate:"required"`
}

type GetPostResponse struct {
	class.Payload
	Post
}

type GetPostHistoryRequest struct {
	class.Payload
	CommunityID string `param:"community_id" validate:"required"`
	PostID      string `param:"post_id" validate:"required"`
}

type GetPostHistoryResponse struct {
	class.Payload
	History []PostHistory `json:"history"`
}

type GetPostThreadRequest struct {
	class.Payload
	CommunityID string `param:"community_id" validate:"required"`
	PostID      string `param:"post_id" validate:"required"`
}

type GetPostThreadResponse struct {
	class.Payload
	Thread []Post `json:"thread"`
}

type GetPostListRequest struct {
	class.Payload
	CommunityID string  `param:"id" validate:"required"`
	Type        *string `query:"type" validate:"omitempty,required"`
}

type GetPostListResponse struct {
	class.Payload
	Posts []Post `json:"posts"`
}

type PostPostRequest struct {
	class.Payload
	CommunityID  string    `param:"id" validate:"required"`
	Type         string    `json:"type" validate:"required"`
	ThreadID     *string   `json:"thread_id" validate:"omitempty,required"`
	Priority     *int      `json:"priority" validate:"omitempty,required"`
	RecipientIDs *[]string `json:"recipient_ids" validate:"omitempty,required"`
	Message      string    `json:"message" validate:"required"`
	Categories   *[]string `json:"categories" validate:"omitempty,required"`
	State        *string   `json:"state" validate:"omitempty,required"`
	Media        *[]string `json:"media" validate:"omitempty,required"`
	Widgets      *struct {
		PollOptions *[]string `json:"poll_options" validate:"omitempty,required"`
	} `json:"widgets" validate:"omitempty,required"`
}

type PostPostResponse struct {
	class.Payload
	Post
}

type PutPostRequest struct {
	class.Payload
	CommunityID string    `param:"community_id" validate:"required"`
	PostID      string    `param:"post_id" validate:"required"`
	Message     *string   `json:"message" validate:"omitempty,required"`
	Categories  *[]string `json:"categories" validate:"omitempty,required"`
	State       *string   `json:"state" validate:"omitempty,required"`
	Media       *[]string `json:"media" validate:"omitempty,required"`
	Widgets     *struct {
		PollOptions *[]string `json:"poll_options" validate:"omitempty,required"`
	} `json:"widgets" validate:"omitempty,required"`
}

type PutPostResponse struct {
	class.Payload
	Post
}

type PostVotePostRequest struct {
	class.Payload
	CommunityID string `param:"community_id" validate:"required"`
	PostID      string `param:"post_id" validate:"required"`
}

type PostVotePostResponse struct {
	class.Payload
	Post
}

type PostUnvotePostRequest struct {
	class.Payload
	CommunityID string `param:"community_id" validate:"required"`
	PostID      string `param:"post_id" validate:"required"`
}

type PostUnvotePostResponse struct {
	class.Payload
	Post
}

type PostVotePostPollRequest struct {
	class.Payload
	CommunityID string `param:"community_id" validate:"required"`
	PostID      string `param:"post_id" validate:"required"`
	Option      string `json:"option" validate:"required"`
}

type PostVotePostPollResponse struct {
	class.Payload
	Post
}

type PostPinPostRequest struct {
	class.Payload
	CommunityID string `param:"community_id" validate:"required"`
	PostID      string `param:"post_id" validate:"required"`
}

type PostPinPostResponse struct {
	class.Payload
	Community
}

type PostUnpinPostRequest struct {
	class.Payload
	CommunityID string `param:"community_id" validate:"required"`
	PostID      string `param:"post_id" validate:"required"`
}

type PostUnpinPostResponse struct {
	class.Payload
	Community
}
