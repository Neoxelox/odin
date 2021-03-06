package repository

import (
	"context"
	"fmt"

	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/class"
	"github.com/neoxelox/odin/internal/core"
	"github.com/neoxelox/odin/internal/database"
	"github.com/neoxelox/odin/pkg/model"
)

const (
	POST_TABLE         = "post"
	POST_HISTORY_TABLE = "post_history"
)

var ErrPostGeneric = internal.NewError("Post query failed")

type PostRepository struct {
	class.Repository
}

func NewPostRepository(configuration internal.Configuration, logger core.Logger, database database.Database) *PostRepository {
	return &PostRepository{
		Repository: *class.NewRepository(configuration, logger, database),
	}
}

func (self *PostRepository) CreatePost(ctx context.Context, post model.Post) (*model.Post, error) {
	var p model.Post

	query := fmt.Sprintf(`INSERT INTO "%s"
						  ("id", "thread_id", "creator_id", "last_history_id", "type", "priority", "recipient_ids", "voter_ids", "created_at")
						  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
						  RETURNING *;`, POST_TABLE)

	err := self.Database.Query(
		ctx, query, post.ID, post.ThreadID, post.CreatorID, post.LastHistoryID, post.Type, post.Priority, post.RecipientIDs, post.VoterIDs, post.CreatedAt).Scan(&p)
	if err != nil {
		return nil, ErrPostGeneric().Wrap(err)
	}

	return &p, nil
}

func (self *PostRepository) CreateHistory(ctx context.Context, history model.PostHistory) (*model.PostHistory, error) {
	var h model.PostHistory

	query := fmt.Sprintf(`INSERT INTO "%s"
						  ("id", "post_id", "updator_id", "message", "categories", "state", "media", "widgets", "created_at")
						  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
						  RETURNING *;`, POST_HISTORY_TABLE)

	err := self.Database.Query(
		ctx, query, history.ID, history.PostID, history.UpdatorID, history.Message, history.Categories, history.State, history.Media, history.Widgets, history.CreatedAt).Scan(&h)
	if err != nil {
		return nil, ErrPostGeneric().Wrap(err)
	}

	return &h, nil
}

func (self *PostRepository) GetByID(ctx context.Context, id string) (*model.Post, *model.PostHistory, error) {
	var p model.Post
	var h model.PostHistory

	query := fmt.Sprintf(`SELECT * FROM "%s"
						  WHERE "id" = $1;`, POST_TABLE)

	err := self.Database.Query(ctx, query, id).Scan(&p)
	switch {
	case database.ErrNoRows().Is(err):
		return nil, nil, nil
	case err != nil:
		return nil, nil, ErrPostGeneric().Wrap(err)
	}

	query = fmt.Sprintf(`SELECT * FROM "%s"
						 WHERE "id" = $1;`, POST_HISTORY_TABLE)

	err = self.Database.Query(ctx, query, p.LastHistoryID).Scan(&h)
	switch {
	case err == nil:
		return &p, &h, nil
	case database.ErrNoRows().Is(err):
		return &p, nil, nil
	default:
		return nil, nil, ErrPostGeneric().Wrap(err)
	}
}

func (self *PostRepository) GetHistory(ctx context.Context, id string) ([]model.PostHistory, error) {
	var hs []model.PostHistory

	query := fmt.Sprintf(`SELECT * FROM "%s"
						  WHERE "post_id" = $1;`, POST_HISTORY_TABLE)

	err := self.Database.Query(ctx, query, id).Scan(&hs)
	switch {
	case err == nil:
		return hs, nil
	case database.ErrNoRows().Is(err):
		return []model.PostHistory{}, nil
	default:
		return nil, ErrPostGeneric().Wrap(err)
	}
}

func (self *PostRepository) GetSubposts(ctx context.Context, id string) (int, error) {
	var s int

	query := fmt.Sprintf(`SELECT COUNT(*) FROM "%s"
						  WHERE "thread_id" = $1;`, POST_TABLE)

	err := self.Database.Query(ctx, query, id).Scan(&s)
	switch {
	case err == nil:
		return s, nil
	case database.ErrNoRows().Is(err):
		return 0, nil
	default:
		return 0, ErrPostGeneric().Wrap(err)
	}
}

// TODO: ADD LOTS OF FILTER AND ORDER PARAMS
// TODO: DISCUSS WHETHER HAVING COMMUNITY_ID DIRECTLY ON THE POST IN ORDER TO AVOID JOINING TABLES...
// TODO: USE A QUERY BUILDER, MYGOD...
func (self *PostRepository) ListByCommunityID(ctx context.Context, communityID string, typee *string) ([]model.Post, []model.PostHistory, error) {
	var err error
	var ps []model.Post
	var hs []model.PostHistory

	query := fmt.Sprintf(`SELECT "post".* FROM "%s" as "post"
						  INNER JOIN "%s" as "membership" ON "post"."creator_id" = "membership"."id"
						  WHERE "membership"."community_id" = $1 AND "post"."thread_id" is NULL`, POST_TABLE, MEMBERSHIP_TABLE)

	if typee != nil {
		query = fmt.Sprintf(`%s AND "post"."type" = $2;`, query)
		err = self.Database.Query(ctx, query, communityID, typee).Scan(&ps)
	} else {
		query = query + ";"
		err = self.Database.Query(ctx, query, communityID).Scan(&ps)
	}

	switch {
	case database.ErrNoRows().Is(err):
		return []model.Post{}, []model.PostHistory{}, nil
	case err != nil:
		return nil, nil, ErrPostGeneric().Wrap(err)
	}

	historyIDs := []string{}
	for _, post := range ps {
		historyIDs = append(historyIDs, *post.LastHistoryID)
	}

	query = fmt.Sprintf(`SELECT * FROM "%s"
						 WHERE "id" = ANY ($1);`, POST_HISTORY_TABLE)

	err = self.Database.Query(ctx, query, historyIDs).Scan(&hs)
	switch {
	case err == nil:
		return ps, hs, nil
	case database.ErrNoRows().Is(err):
		return ps, []model.PostHistory{}, nil
	default:
		return nil, nil, ErrPostGeneric().Wrap(err)
	}
}

func (self *PostRepository) ListByThreadID(ctx context.Context, threadID string) ([]model.Post, []model.PostHistory, error) {
	var ps []model.Post
	var hs []model.PostHistory

	query := fmt.Sprintf(`SELECT * FROM "%s"
						  WHERE "thread_id" = $1;`, POST_TABLE)

	err := self.Database.Query(ctx, query, threadID).Scan(&ps)
	switch {
	case database.ErrNoRows().Is(err):
		return []model.Post{}, []model.PostHistory{}, nil
	case err != nil:
		return nil, nil, ErrPostGeneric().Wrap(err)
	}

	historyIDs := []string{}
	for _, post := range ps {
		historyIDs = append(historyIDs, *post.LastHistoryID)
	}

	query = fmt.Sprintf(`SELECT * FROM "%s"
						 WHERE "id" = ANY ($1);`, POST_HISTORY_TABLE)

	err = self.Database.Query(ctx, query, historyIDs).Scan(&hs)
	switch {
	case err == nil:
		return ps, hs, nil
	case database.ErrNoRows().Is(err):
		return ps, []model.PostHistory{}, nil
	default:
		return nil, nil, ErrPostGeneric().Wrap(err)
	}
}

func (self *PostRepository) UpdateHistory(ctx context.Context, id string, historyID string) error {
	query := fmt.Sprintf(`UPDATE "%s"
						  SET "last_history_id" = $1
						  WHERE "id" = $2;`, POST_TABLE)

	affected, err := self.Database.Exec(ctx, query, historyID, id)
	if err != nil {
		return ErrPostGeneric().Wrap(err)
	}

	if affected != 1 {
		return ErrPostGeneric()
	}

	return nil
}

func (self *PostRepository) UpdateVoters(ctx context.Context, id string, voterIDs []string) error {
	query := fmt.Sprintf(`UPDATE "%s"
						  SET "voter_ids" = $1
						  WHERE "id" = $2;`, POST_TABLE)

	affected, err := self.Database.Exec(ctx, query, voterIDs, id)
	if err != nil {
		return ErrPostGeneric().Wrap(err)
	}

	if affected != 1 {
		return ErrPostGeneric()
	}

	return nil
}

func (self *PostRepository) UpdateVotersAndPriority(ctx context.Context, id string, voterIDs []string, priority int) error {
	query := fmt.Sprintf(`UPDATE "%s"
						  SET "voter_ids" = $1, "priority" = $2
						  WHERE "id" = $3;`, POST_TABLE)

	affected, err := self.Database.Exec(ctx, query, voterIDs, priority, id)
	if err != nil {
		return ErrPostGeneric().Wrap(err)
	}

	if affected != 1 {
		return ErrPostGeneric()
	}

	return nil
}

func (self *PostRepository) UpdateWidgets(ctx context.Context, historyID string, widgets model.PostWidgets) error {
	query := fmt.Sprintf(`UPDATE "%s"
						  SET "widgets" = $1
						  WHERE "id" = $2;`, POST_HISTORY_TABLE)

	affected, err := self.Database.Exec(ctx, query, widgets, historyID)
	if err != nil {
		return ErrPostGeneric().Wrap(err)
	}

	if affected != 1 {
		return ErrPostGeneric()
	}

	return nil
}
