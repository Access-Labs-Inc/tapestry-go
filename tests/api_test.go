package tests

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/tapestry"
)

var (
	client      tapestry.TapestryClient
	testProfile *tapestry.ProfileResponse
)

func TestMain(m *testing.M) {
	apiKey := os.Getenv("TAPESTRY_API_KEY")
	baseURL := os.Getenv("TAPESTRY_API_BASE_URL")
	if apiKey == "" || baseURL == "" {
		panic("TAPESTRY_API_KEY and TAPESTRY_API_BASE_URL must be set")
	}

	client = tapestry.NewTapestryClient(apiKey, baseURL, tapestry.ExecutionConfirmedParsed, "SOLANA")

	var err error
	testProfile, err = client.FindOrCreateProfile(tapestry.FindOrCreateProfileParameters{
		WalletAddress: "97QsK6DFcUZFz8tkRTcYypysyWsrGuC5CcHJuZMWAQhH",
		Username:      "test_user_20241108143421",
		Bio:           "Test bio",
		Image:         "https://example.com/image.jpg",
	})
	if err != nil {
		panic("Failed to get test profile: " + err.Error())
	}

	os.Exit(m.Run())
}

func TestProfileOperations(t *testing.T) {
	// Test GetProfileByID
	profile, err := client.GetProfileByID(testProfile.Profile.ID)

	if err != nil {
		t.Fatalf("GetProfileByID failed: %v", err)
	}
	if profile.Profile.Username != testProfile.Profile.Username {
		t.Errorf("Expected username %s, got %s", testProfile.Profile.Username, profile.Profile.Username)
	}

	// Test UpdateProfile
	newUsername := "updated_user_" + time.Now().Format("20060102150405")
	err = client.UpdateProfile(testProfile.Profile.ID, tapestry.UpdateProfileParameters{
		Username: newUsername,
		Bio:      "Updated bio",
	})
	if err != nil {
		t.Fatalf("UpdateProfile failed: %v", err)
	}

	// Verify update
	// updatedProfile, err := client.GetProfileByID(testProfile.Profile.ID)
	// if err != nil {
	// 	t.Fatalf("GetProfileByID after update failed: %v", err)
	// }
	// if updatedProfile.Profile.Username != newUsername {
	// 	t.Errorf("Expected updated username %s, got %s", newUsername, updatedProfile.Profile.Username)
	// }
}

func TestContentOperations(t *testing.T) {
	// Test FindOrCreateContent
	contentProps := []tapestry.ContentProperty{
		{Key: "title", Value: "Test Content"},
		{Key: "description", Value: "Test Description"},
	}
	randomContentId := "test_content_" + time.Now().Format("20060102150405")
	content, err := client.FindOrCreateContent(testProfile.Profile.ID, randomContentId, contentProps)
	if err != nil {
		t.Fatalf("FindOrCreateContent failed: %v", err)
	}

	// Test GetContentByID
	retrievedContent, err := client.GetContentByID(randomContentId)
	if err != nil {
		t.Fatalf("GetContentByID failed: %v", err)
	}
	if retrievedContent.Content.ID != content.Content.ID {
		t.Errorf("Expected content ID %s, got %s", content.Content.ID, retrievedContent.Content.ID)
	}

	// Test UpdateContent
	updatedProps := []tapestry.ContentProperty{
		{Key: "title", Value: "Updated Title"},
		{Key: "description", Value: "Updated Description"},
	}
	_, err = client.UpdateContent(randomContentId, updatedProps)
	if err != nil {
		t.Fatalf("UpdateContent failed: %v", err)
	}

	// Test GetContents
	contents, err := client.GetContents(
		tapestry.WithProfileID(testProfile.Profile.ID),
		tapestry.WithPagination("1", "10"),
		tapestry.WithOrderBy("created_at", tapestry.GetContentsSortDirectionDesc),
	)
	if err != nil {
		t.Fatalf("GetContents failed: %v", err)
	}
	if len(contents.Contents) == 0 {
		t.Error("Expected at least one content item")
	}

	// Test DeleteContent
	err = client.DeleteContent(randomContentId)
	if err != nil {
		t.Fatalf("DeleteContent failed: %v", err)
	}
}

func TestCommentOperations(t *testing.T) {
	// Create test content first
	contentProps := []tapestry.ContentProperty{
		{Key: "title", Value: "Test Content for Comments"},
	}
	randomContentId := "test_content_" + time.Now().Format("20060102150405")
	fmt.Println("profile id", testProfile.Profile.ID)
	content, err := client.FindOrCreateContent(testProfile.Profile.ID, randomContentId, contentProps)
	if err != nil {
		t.Fatalf("Failed to create test content: %v", err)
	}

	// Verify initial comment count is 0
	initialContent, err := client.GetContentByID(content.Content.ID)
	if err != nil {
		t.Fatalf("GetContentByID failed: %v", err)
	}
	if initialContent.SocialCounts.CommentCount != 0 {
		t.Errorf("Expected initial comment count 0, got %d", initialContent.SocialCounts.CommentCount)
	}

	// Test CreateComment
	comment, err := client.CreateComment(tapestry.CreateCommentOptions{
		ContentID: content.Content.ID,
		ProfileID: testProfile.Profile.ID,
		Text:      "Test comment",
		Properties: []tapestry.CommentProperty{
			{Key: "test", Value: "property"},
		},
	})
	if err != nil {
		t.Fatalf("CreateComment failed: %v", err)
	}

	// Test UpdateComment
	newProperty := "new property"
	_, err = client.UpdateComment(comment.Comment.ID, []tapestry.CommentProperty{
		{Key: "test", Value: newProperty},
	})
	if err != nil {
		t.Fatalf("UpdateComment failed: %v", err)
	}
	// Verify comment count increased to 1
	contentAfterComment, err := client.GetContentByID(content.Content.ID)
	if err != nil {
		t.Fatalf("GetContentByID failed: %v", err)
	}
	if contentAfterComment.SocialCounts.CommentCount != 1 {
		t.Errorf("Expected comment count 1, got %d", contentAfterComment.SocialCounts.CommentCount)
	}

	// Test GetCommentByID - verify initial like count
	commentDetail, err := client.GetCommentByID(comment.Comment.ID, testProfile.Profile.ID)
	if err != nil {
		t.Fatalf("GetCommentByID failed: %v", err)
	}
	if commentDetail.SocialCounts.LikeCount != 0 {
		t.Errorf("Expected initial comment like count 0, got %d", commentDetail.SocialCounts.LikeCount)
	}

	// Test liking the comment
	err = client.CreateLike(comment.Comment.ID, testProfile.Profile)
	if err != nil {
		t.Fatalf("CreateLike on comment failed: %v", err)
	}

	// Verify like count increased to 1
	commentAfterLike, err := client.GetCommentByID(comment.Comment.ID, testProfile.Profile.ID)
	if err != nil {
		t.Fatalf("GetCommentByID after like failed: %v", err)
	}
	if commentAfterLike.SocialCounts.LikeCount != 1 {
		t.Errorf("Expected comment like count 1, got %d", commentAfterLike.SocialCounts.LikeCount)
	}
	// if !commentAfterLike.RequestingProfileSocialInfo["hasLiked"].(bool) {
	// 	t.Error("Expected hasLiked to be true")
	// }

	// Test unliking the comment
	err = client.DeleteLike(comment.Comment.ID, testProfile.Profile)
	if err != nil {
		t.Fatalf("DeleteLike on comment failed: %v", err)
	}

	// Verify like count back to 0
	commentAfterUnlike, err := client.GetCommentByID(comment.Comment.ID, testProfile.Profile.ID)
	if err != nil {
		t.Fatalf("GetCommentByID after unlike failed: %v", err)
	}
	if commentAfterUnlike.SocialCounts.LikeCount != 0 {
		t.Errorf("Expected comment like count 0, got %d", commentAfterUnlike.SocialCounts.LikeCount)
	}
	// if commentAfterUnlike.RequestingProfileSocialInfo["hasLiked"].(bool) {
	// 	t.Error("Expected hasLiked to be false")
	// }

	// Test GetComments
	comments, err := client.GetComments(tapestry.GetCommentsOptions{
		ContentID:           content.Content.ID,
		RequestingProfileID: testProfile.Profile.ID,
		Page:                1,
		PageSize:            10,
	})
	if err != nil {
		t.Fatalf("GetComments failed: %v", err)
	}
	if len(comments.Comments) == 0 {
		t.Error("Expected at least one comment")
	}

	// Test DeleteComment
	err = client.DeleteComment(comment.Comment.ID)
	if err != nil {
		t.Fatalf("DeleteComment failed: %v", err)
	}

	// Verify comment count back to 0
	contentAfterDelete, err := client.GetContentByID(content.Content.ID)
	if err != nil {
		t.Fatalf("GetContentByID failed: %v", err)
	}
	if contentAfterDelete.SocialCounts.CommentCount != 0 {
		t.Errorf("Expected comment count 0 after delete, got %d", contentAfterDelete.SocialCounts.CommentCount)
	}
}

func TestLikeOperations(t *testing.T) {
	// Create test content first
	contentProps := []tapestry.ContentProperty{
		{Key: "title", Value: "Test Content for Likes"},
	}
	randomContentId := "test_content_" + time.Now().Format("20060102150405")
	content, err := client.FindOrCreateContent(testProfile.Profile.ID, randomContentId, contentProps)
	if err != nil {
		t.Fatalf("Failed to create test content: %v", err)
	}

	// Verify initial like count is 0
	initialContent, err := client.GetContentByID(content.Content.ID)
	if err != nil {
		t.Fatalf("GetContentByID failed: %v", err)
	}
	if initialContent.SocialCounts.LikeCount != 0 {
		t.Errorf("Expected initial like count 0, got %d", initialContent.SocialCounts.LikeCount)
	}

	// Test CreateLike
	err = client.CreateLike(content.Content.ID, testProfile.Profile)
	if err != nil {
		t.Fatalf("CreateLike failed: %v", err)
	}

	// Verify like count increased to 1
	contentAfterLike, err := client.GetContentByID(content.Content.ID)
	if err != nil {
		t.Fatalf("GetContentByID failed: %v", err)
	}
	if contentAfterLike.SocialCounts.LikeCount != 1 {
		t.Errorf("Expected like count 1, got %d", contentAfterLike.SocialCounts.LikeCount)
	}

	// Test DeleteLike
	err = client.DeleteLike(content.Content.ID, testProfile.Profile)
	if err != nil {
		t.Fatalf("DeleteLike failed: %v", err)
	}

	// Verify like count back to 0
	contentAfterDelete, err := client.GetContentByID(content.Content.ID)
	if err != nil {
		t.Fatalf("GetContentByID failed: %v", err)
	}
	if contentAfterDelete.SocialCounts.LikeCount != 0 {
		t.Errorf("Expected like count 0 after delete, got %d", contentAfterDelete.SocialCounts.LikeCount)
	}
}

func TestFollowerOperations(t *testing.T) {
	// Create followee profile
	followee, err := client.FindOrCreateProfile(tapestry.FindOrCreateProfileParameters{
		WalletAddress: solana.NewWallet().PublicKey().String(),
		Username:      "followee_" + time.Now().Format("20060102150405"),
	})
	if err != nil {
		t.Fatalf("Failed to create followee: %v", err)
	}

	// Create two additional test profiles with random Solana addresses
	follower1, err := client.FindOrCreateProfile(tapestry.FindOrCreateProfileParameters{
		WalletAddress: solana.NewWallet().PublicKey().String(),
		Username:      "follower1_" + time.Now().Format("20060102150405"),
	})
	if err != nil {
		t.Fatalf("Failed to create follower1: %v", err)
	}

	follower2, err := client.FindOrCreateProfile(tapestry.FindOrCreateProfileParameters{
		WalletAddress: solana.NewWallet().PublicKey().String(),
		Username:      "follower2_" + time.Now().Format("20060102150405"),
	})
	if err != nil {
		t.Fatalf("Failed to create follower2: %v", err)
	}

	// Add followers to followee profile
	err = client.AddFollower(follower1.Profile.ID, followee.Profile.ID)
	if err != nil {
		t.Fatalf("Failed to add follower1: %v", err)
	}

	err = client.AddFollower(follower2.Profile.ID, followee.Profile.ID)
	if err != nil {
		t.Fatalf("Failed to add follower2: %v", err)
	}

	// Also make follower1 follow follower2
	err = client.AddFollower(follower1.Profile.ID, follower2.Profile.ID)
	if err != nil {
		t.Fatalf("Failed to make follower1 follow follower2: %v", err)
	}

	// Verify followers of followee profile
	followers, err := client.GetFollowers(followee.Profile.ID)
	if err != nil {
		t.Fatalf("GetFollowers failed: %v", err)
	}
	if len(followers.Profiles) != 2 {
		t.Errorf("Expected 2 followers, got %d", len(followers.Profiles))
	}

	// Verify following for follower1
	following, err := client.GetFollowing(follower1.Profile.ID)
	if err != nil {
		t.Fatalf("GetFollowing failed: %v", err)
	}
	if len(following.Profiles) != 2 {
		t.Errorf("Expected follower1 to be following 2 profiles, got %d", len(following.Profiles))
	}

	// Remove one follower and verify
	err = client.RemoveFollower(follower1.Profile.ID, followee.Profile.ID)
	if err != nil {
		t.Fatalf("RemoveFollower failed: %v", err)
	}

	// Verify updated follower count
	updatedFollowers, err := client.GetFollowers(followee.Profile.ID)
	if err != nil {
		t.Fatalf("GetFollowers after removal failed: %v", err)
	}
	if len(updatedFollowers.Profiles) != 1 {
		t.Errorf("Expected 1 follower after removal, got %d", len(updatedFollowers.Profiles))
	}

	// Verify updated following count for follower1
	updatedFollowing, err := client.GetFollowing(follower1.Profile.ID)
	if err != nil {
		t.Fatalf("GetFollowing after removal failed: %v", err)
	}
	if len(updatedFollowing.Profiles) != 1 {
		t.Errorf("Expected follower1 to be following 1 profile after removal, got %d", len(updatedFollowing.Profiles))
	}

	// Test GetFollowingWhoFollow
	_, err = client.GetFollowingWhoFollow(follower2.Profile.ID, follower1.Profile.ID)
	if err != nil {
		t.Fatalf("GetFollowingWhoFollow failed: %v", err)
	}

	// TODO: assert output

	// Test GetSuggestedProfiles
	_, err = client.GetSuggestedProfiles(follower1.WalletAddress, true)
	if err != nil {
		t.Fatalf("GetSuggestedProfiles failed: %v", err)
	}

	// TODO: assert output
}
