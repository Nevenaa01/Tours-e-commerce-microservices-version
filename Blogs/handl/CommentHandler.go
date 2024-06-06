package handl

import (
	"blogs_service/model"
	"blogs_service/proto/blogs"
	"context"
	"fmt"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type CommentHandler struct {
	blogs.UnimplementedBlogServiceServer
	DatabaseConnection *gorm.DB
}

func timeToTimestamp(t time.Time) *timestamppb.Timestamp {
	return timestamppb.New(t)
}

func timestampToTimee(t *timestamppb.Timestamp) time.Time {
	return t.AsTime()
}

func (h BlogHandler) CreateComment(ctx context.Context, request *blogs.Comment) (*blogs.Comment, error) {
	fmt.Println("usao u create commentara")

	var comment = model.Comment{
		Id:           int(request.Id),
		UserId:       int(request.UserId),
		CreationDate: request.CreationDate.AsTime(),
		Description:  request.Description,
		LastEditDate: request.LastEditDate.AsTime(),
		BlogId:       int(request.BlogId),
	}
	if err := h.DatabaseConnection.Table(`blog."Comments"`).Create(&comment).Error; err != nil {
		return nil, err
	}

	response := &blogs.Comment{
		Id:           int32(comment.Id),
		UserId:       int32(comment.UserId),
		CreationDate: request.CreationDate,
		Description:  comment.Description,
		LastEditDate: request.LastEditDate,
		BlogId:       int32(comment.BlogId),
	}

	return response, nil

}

func (h BlogHandler) GetCommentsByBlogIdAsync(ctx context.Context, request *blogs.GetCommentRequest) (*blogs.ListComment, error) {
	fmt.Println("usao u GetCommentsByBlogIdAsync")
	blogId := request.BlogId

	var comments []model.Comment

	if err := h.DatabaseConnection.Table(`blog."Comments"`).Where(`"BlogId" = ?`, blogId).Find(&comments).Error; err != nil {
		return nil, err
	}

	protoComments := make([]*blogs.Comment, len(comments))
	for i, comment := range comments {

		protoComments[i] = &blogs.Comment{
			Id:           int32(comment.Id),
			UserId:       int32(comment.UserId),
			CreationDate: timeToTimestamp(comment.CreationDate),
			Description:  comment.Description,
			LastEditDate: timeToTimestamp(comment.LastEditDate),
			BlogId:       int32(comment.BlogId),
		}
	}

	response := &blogs.ListComment{
		Comments: protoComments,
	}

	return response, nil
}

func (h BlogHandler) UpdateComment(ctx context.Context, request *blogs.Comment) (*blogs.Comment, error) {
	fmt.Println("usao u update komentara")

	var existingComment model.Comment
	if err := h.DatabaseConnection.Table(`blog."Comments"`).Where(`"Id" = ?`, request.Id).First(&existingComment).Error; err != nil {
		return nil, err
	}

	existingComment.UserId = int(request.UserId)
	existingComment.CreationDate = timestampToTimee(request.CreationDate)
	existingComment.Description = request.Description
	existingComment.LastEditDate = timestampToTimee(request.LastEditDate)
	existingComment.BlogId = int(request.BlogId)

	if err := h.DatabaseConnection.Table(`blog."Comments"`).Save(&existingComment).Error; err != nil {
		return nil, err
	}

	response := &blogs.Comment{
		Id:           int32(existingComment.Id),
		UserId:       int32(existingComment.UserId),
		CreationDate: request.CreationDate,
		Description:  existingComment.Description,
		LastEditDate: request.LastEditDate,
		BlogId:       int32(existingComment.BlogId),
	}

	return response, nil

}

func (h BlogHandler) DeleteComment(ctx context.Context, request *blogs.DeleteCommentRequest) (*blogs.Comment, error) {
	fmt.Println("usao u delete komentara")

	var existingComment model.Comment
	if err := h.DatabaseConnection.Table(`blog."Comments"`).Where(`"Id" = ?`, request.Id).First(&existingComment).Error; err != nil {
		return nil, err
	}

	if err := h.DatabaseConnection.Table(`blog."Comments"`).Delete(&existingComment).Error; err != nil {
		return nil, err
	}

	response := &blogs.Comment{
		Id:           int32(existingComment.Id),
		UserId:       int32(existingComment.UserId),
		CreationDate: timeToTimestamp(existingComment.CreationDate),
		Description:  existingComment.Description,
		LastEditDate: timeToTimestamp(existingComment.LastEditDate),
		BlogId:       int32(existingComment.BlogId),
	}

	return response, nil

}
