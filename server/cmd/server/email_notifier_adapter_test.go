package main

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	dbfx "github.com/multica-ai/multica/server/internal/testutil"
	"github.com/multica-ai/multica/server/internal/util"
	db "github.com/multica-ai/multica/server/pkg/db/generated"
)

// The email notifier sends what the recipient's inbox would show once read:
// a guest is emailed only about issues in its own projects, and nobody about
// a private project they are not in.
func TestEmailNotifierCanSeeIssue(t *testing.T) {
	if testPool == nil {
		t.Skip("no database connection")
	}
	f := dbfx.New(testPool, testWorkspaceID, testUserID)
	f.WorkspaceID = f.Workspace(t, "Email visibility", "email-visibility-"+uuid.NewString())
	f.Member(t, f.WorkspaceID, f.UserID, "owner")

	guest := f.User(t, "Guest", "guest-"+uuid.NewString()[:8]+"@example.test")
	f.Member(t, f.WorkspaceID, guest, "guest")
	member := f.User(t, "Member", "member-"+uuid.NewString()[:8]+"@example.test")
	f.Member(t, f.WorkspaceID, member, "member")
	outsider := f.User(t, "Outsider", "outsider-"+uuid.NewString()[:8]+"@example.test")

	guestProject := f.Project(t, "Shared with guest")
	f.InsertNoID(t, "project_member", dbfx.Cols{"workspace_id": f.WorkspaceID, "project_id": guestProject, "user_id": guest},
		"project_id = $1 AND user_id = $2", guestProject, guest)
	otherProject := f.Project(t, "Not shared")
	privateProject := f.Project(t, "Private", dbfx.Cols{"visibility": "private"})

	inGuestProject := f.Issue(t, "In the guest's project", dbfx.Cols{"project_id": guestProject})
	inOtherProject := f.Issue(t, "In another project", dbfx.Cols{"project_id": otherProject})
	inPrivateProject := f.Issue(t, "In a private project", dbfx.Cols{"project_id": privateProject})
	unfiled := f.Issue(t, "Outside any project")

	a := emailNotifierQueries{q: db.New(testPool)}
	ws := mustParseUUID(t, f.WorkspaceID)
	cases := []struct {
		name    string
		user    string
		issueID string
		want    bool
	}{
		{"guest, own project", guest, inGuestProject, true},
		{"guest, other project", guest, inOtherProject, false},
		{"guest, unfiled issue", guest, unfiled, false},
		{"guest, no issue", guest, "", false},
		{"member, other project", member, inOtherProject, true},
		{"member, unfiled issue", member, unfiled, true},
		{"member, no issue", member, "", true},
		{"member, private project", member, inPrivateProject, false},
		{"non-member", outsider, inOtherProject, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := a.CanSeeIssue(context.Background(), ws, mustParseUUID(t, tc.user), tc.issueID)
			if err != nil {
				t.Fatalf("CanSeeIssue: %v", err)
			}
			if got != tc.want {
				t.Fatalf("CanSeeIssue = %v, want %v", got, tc.want)
			}
		})
	}
}

func mustParseUUID(t *testing.T, s string) pgtype.UUID {
	t.Helper()
	id, err := util.ParseUUID(s)
	if err != nil {
		t.Fatalf("parse uuid %q: %v", s, err)
	}
	return id
}
