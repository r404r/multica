package main

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/multica-ai/multica/server/internal/integrations/email"
	"github.com/multica-ai/multica/server/internal/projectaccess"
	"github.com/multica-ai/multica/server/internal/util"
	db "github.com/multica-ai/multica/server/pkg/db/generated"
)

// emailNotifierQueries bridges email.NotifierQueries onto *db.Queries by
// projecting onto the single fields the notifier needs. Keeps the email
// package independent of the db package's row types.
type emailNotifierQueries struct {
	q *db.Queries
}

func (a emailNotifierQueries) GetUserEmail(ctx context.Context, userID pgtype.UUID) (string, error) {
	u, err := a.q.GetUser(ctx, userID)
	if err != nil {
		return "", err
	}
	return u.Email, nil
}

func (a emailNotifierQueries) GetNotificationPreference(ctx context.Context,
	arg email.GetNotificationPreferenceParams) ([]byte, error) {
	row, err := a.q.GetNotificationPreference(ctx, db.GetNotificationPreferenceParams{
		WorkspaceID: arg.WorkspaceID,
		UserID:      arg.UserID,
	})
	if err != nil {
		return nil, err
	}
	return row.Preferences, nil
}

func (a emailNotifierQueries) GetWorkspaceSlug(ctx context.Context, wsID pgtype.UUID) (string, error) {
	w, err := a.q.GetWorkspace(ctx, wsID)
	if err != nil {
		return "", err
	}
	return w.Slug, nil
}

// CanSeeIssue applies the same project visibility the inbox applies when it
// is read (projectaccess.ForMember): a guest sees only its projects' issues
// and nothing outside a project, anyone else every issue outside the private
// projects they are not in.
func (a emailNotifierQueries) CanSeeIssue(ctx context.Context, wsID, userID pgtype.UUID, issueID string) (bool, error) {
	member, err := a.q.GetMemberByUserAndWorkspace(ctx, db.GetMemberByUserAndWorkspaceParams{
		UserID:      userID,
		WorkspaceID: wsID,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	hidden, err := projectaccess.ForMember(ctx, a.q, member)
	if err != nil {
		return false, err
	}
	if hidden.Empty() {
		return true, nil
	}
	var projectID pgtype.UUID // invalid: outside any project
	if issueID != "" {
		id, err := util.ParseUUID(issueID)
		if err != nil {
			return false, err
		}
		issue, err := a.q.GetIssue(ctx, id)
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		if err != nil {
			return false, err
		}
		if issue.WorkspaceID != wsID {
			return false, nil
		}
		projectID = issue.ProjectID
	}
	return hidden.CanSeeProject(projectID), nil
}
