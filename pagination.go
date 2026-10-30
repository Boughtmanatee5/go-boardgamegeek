package boardgamegeek

import (
	"context"
	"errors"
	"iter"
)

// fetchPage requests one page and returns its items and the total number of items across all pages
type fetchPage[T any] func(ctx context.Context, page int) (items []T, total int, err error)

// paginate yields every item from startPage (1 if 0) onwards, requesting one page at a time.
// It stops after an empty or short page, once the total BGG reports has been reached, or at the
// first error, which is yielded with a nil item. Stopping the range loop early stops further requests.
func paginate[T any](ctx context.Context, startPage int, fetch fetchPage[T]) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		firstPage := max(startPage, 1)
		page := firstPage
		// pageSize is taken from the first page; items on pages before firstPage count towards total
		var pageSize, yielded int

		for {
			items, total, err := fetch(ctx, page)
			if err != nil {
				var zero T
				yield(zero, err)
				return
			}
			if len(items) == 0 {
				return
			}
			if pageSize == 0 {
				pageSize = len(items)
			}

			for _, item := range items {
				if !yield(item, nil) {
					return
				}
			}

			yielded += len(items)
			if len(items) < pageSize || (total > 0 && (firstPage-1)*pageSize+yielded >= total) {
				return
			}
			page++
		}
	}
}

// AllForumThreads yields every thread in a forum, starting at opts.Page
func (c *Client) AllForumThreads(ctx context.Context, opts *ForumOptions) iter.Seq2[*ThreadSummary, error] {
	return paginate(ctx, opts.Page, func(ctx context.Context, page int) ([]*ThreadSummary, int, error) {
		pageOpts := *opts
		pageOpts.Page = page
		resp, err := c.GetForum(ctx, &pageOpts)
		if err != nil {
			return nil, 0, err
		}
		return resp.Threads, resp.NumThreads, nil
	})
}

// AllGuildMembers yields every member of a guild, starting at opts.Page. opts.Members is set for you.
func (c *Client) AllGuildMembers(ctx context.Context, opts *GuildOptions) iter.Seq2[*GuildMember, error] {
	return paginate(ctx, opts.Page, func(ctx context.Context, page int) ([]*GuildMember, int, error) {
		pageOpts := *opts
		pageOpts.Members = true
		pageOpts.Page = page
		resp, err := c.GetGuild(ctx, &pageOpts)
		if err != nil || resp.Members == nil {
			return nil, 0, err
		}
		return resp.Members.Members, resp.Members.Count, nil
	})
}

// AllUserBuddies yields every buddy of a user, starting at opts.Page. opts.Buddies is set for you.
func (c *Client) AllUserBuddies(ctx context.Context, opts *UserOptions) iter.Seq2[*UserBuddy, error] {
	return paginate(ctx, opts.Page, func(ctx context.Context, page int) ([]*UserBuddy, int, error) {
		pageOpts := *opts
		pageOpts.Buddies = true
		pageOpts.Page = page
		resp, err := c.GetUser(ctx, &pageOpts)
		if err != nil || resp.Buddies == nil {
			return nil, 0, err
		}
		return resp.Buddies.Buddies, resp.Buddies.Total, nil
	})
}

// AllUserGuilds yields every guild a user belongs to, starting at opts.Page. opts.Guilds is set for you.
func (c *Client) AllUserGuilds(ctx context.Context, opts *UserOptions) iter.Seq2[*UserGuild, error] {
	return paginate(ctx, opts.Page, func(ctx context.Context, page int) ([]*UserGuild, int, error) {
		pageOpts := *opts
		pageOpts.Guilds = true
		pageOpts.Page = page
		resp, err := c.GetUser(ctx, &pageOpts)
		if err != nil || resp.Guilds == nil {
			return nil, 0, err
		}
		return resp.Guilds.Guilds, resp.Guilds.Total, nil
	})
}

// errThingCommentsID is yielded by AllThingComments when opts does not name exactly one thing
var errThingCommentsID = errors.New("AllThingComments needs exactly one id in ThingsOptions.IDs")

// AllThingComments yields every comment on a single thing, starting at opts.Page and requesting
// opts.PageSize comments per page. opts.Comments is set for you unless opts.RatingComments is,
// which yields every rating instead.
func (c *Client) AllThingComments(ctx context.Context, opts *ThingsOptions) iter.Seq2[*Comment, error] {
	return paginate(ctx, opts.Page, func(ctx context.Context, page int) ([]*Comment, int, error) {
		if len(opts.IDs) != 1 {
			return nil, 0, errThingCommentsID
		}
		pageOpts := *opts
		pageOpts.Comments = !opts.RatingComments
		pageOpts.Page = page
		resp, err := c.GetThings(ctx, &pageOpts)
		if err != nil || len(resp.Items) == 0 || resp.Items[0].Comments == nil {
			return nil, 0, err
		}
		comments := resp.Items[0].Comments
		return comments.Comments, int(comments.TotalItems), nil
	})
}

// AllPlays yields every play matching opts, starting at opts.Page
func (c *Client) AllPlays(ctx context.Context, opts *PlaysOptions) iter.Seq2[*Play, error] {
	return paginate(ctx, opts.Page, func(ctx context.Context, page int) ([]*Play, int, error) {
		pageOpts := *opts
		pageOpts.Page = page
		resp, err := c.GetPlays(ctx, &pageOpts)
		if err != nil {
			return nil, 0, err
		}
		return resp.Plays, resp.Total, nil
	})
}
