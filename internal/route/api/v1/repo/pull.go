// Copyright 2016 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package repo

import (
	api "github.com/gogs/go-gogs-client"
	"gogs.io/gogs/internal/conf"
	"gogs.io/gogs/internal/context"
	"gogs.io/gogs/internal/db"
)

func GetPull(c *context.APIContext) {
	issue, err := db.GetIssueByIndex(c.Repo.Repository.ID, c.ParamsInt64(":index"))
	if err != nil {
		c.NotFoundOrError(err, "get pull by index")
		return
	}

	c.JSONSuccess(issue.PullRequest.APIFormat())
}

func listPulls(c *context.APIContext, opts *db.PullsOptions) {
	pulls, err := db.Pulls(opts)
	if err != nil {
		c.Error(err, "list pulls")
		return
	}

	count, err := db.PullsCount(opts)
	if err != nil {
		c.Error(err, "count pulls")
		return
	}

	// FIXME: use IssueList to improve performance.
	apiPulls := make([]*api.PullRequest, len(pulls))
	for i := range pulls {
		if err = pulls[i].LoadAttributes(); err != nil {
			c.Error(err, "load attributes")
			return
		}
		apiPulls[i] = pulls[i].APIFormat()
	}

	c.SetLinkHeader(int(count), conf.UI.IssuePagingNum)
	c.JSONSuccess(&apiPulls)
}

func ListPulls(c *context.APIContext) {
	opts := db.PullsOptions{
		BaseRepoID:   c.Repo.Repository.ID,
		Page:     c.QueryInt("page"),
		HasMerged: false,
		//IsClosed: api.StateType(c.Query("state")) == api.STATE_CLOSED,
	}

	listPulls(c, &opts)
}

func MergePull(c *context.APIContext) {
	issue, err := db.GetIssueByIndex(c.Repo.Repository.ID, c.ParamsInt64(":index"))
	if err != nil {
		c.NotFoundOrError(err, "get issue by index")
		return
	}
	err = db.MergePullRequestAction(c.User, issue.Repo, issue)

	if err != nil {
		c.NotFoundOrError(err, "merging pull request")
		return
	}
	issue, _ = db.GetIssueByIndex(c.Repo.Repository.ID, c.ParamsInt64(":index"))

	c.JSONSuccess(issue.APIFormat())
}
