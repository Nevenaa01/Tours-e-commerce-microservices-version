package service

import (
	"blogs_service/model"
	"blogs_service/repository"
	"time"
)

type CommentService struct {
	CommentRepository *repository.CommentRepository
}

func (service *CommentService) CreateComment(comment *model.Comment) error {
	comment.CreationDate = time.Now()
	err := service.CommentRepository.CreateComment(comment)
	if err != nil {
		return err
	}
	return nil
}

func (service *CommentService) UpdateComment(comment *model.Comment) error {
	err := service.CommentRepository.UpdateComment(comment)
	if err != nil {
		return err
	}
	return nil
}

func (service *CommentService) Delete(commentId int) error {
	err := service.CommentRepository.DeleteComment(commentId)

	if err != nil {
		return err
	}
	return nil
}

func (service *CommentService) GetByBlogId(blogID int) (*[]model.Comment, error) {
	return service.CommentRepository.GetCommentsByBlogId(blogID)
}
