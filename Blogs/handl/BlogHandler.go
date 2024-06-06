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

type BlogHandler struct {
	blogs.UnimplementedBlogServiceServer
	DatabaseConnection *gorm.DB
}

func timeToTimestampp(t time.Time) *timestamppb.Timestamp {
	return timestamppb.New(t)
}

func timestampToTime(t *timestamppb.Timestamp) time.Time {
	return t.AsTime()
}

func (h BlogHandler) GetBlog(ctx context.Context, request *blogs.GetBlogRequest) (*blogs.Blog, error) {
	var blog model.BlogPage
	if err := h.DatabaseConnection.Table(`blog."Blogs"`).Where(`"Blogs"."Id" = ?`, request.Id).First(&blog).Error; err != nil {
		return nil, err
	}

	ratingList := make([]*blogs.Rating, len(blog.Ratings))
	if len(blog.Ratings) > 0 {
		for i, rating := range blog.Ratings {
			ratingList[i] = &blogs.Rating{
				UserId:       int32(rating.UserId),
				CreationDate: timeToTimestampp(rating.CreationDate),
				RatingValue:  int32(rating.RatingValue),
			}
		}
	}

	blogic := &blogs.Blog{
		Id:           int32(blog.Id),
		Title:        blog.Title,
		Description:  blog.Description,
		CreationDate: timeToTimestampp(blog.CreationDate),
		Status:       int32(blog.Status),
		UserId:       int32(blog.UserId),
		RatingSum:    int32(blog.RatingSum),
		Ratings:      ratingList,
	}

	return blogic, nil
}

func (h BlogHandler) CreateBlog(ctx context.Context, request *blogs.Blog) (*blogs.Blog, error) {
	fmt.Printf("Received Blog: %+v\n", request)
	fmt.Printf("Received Blog ID: %d, Title: %s, Description: %s, UserId: %d, RatingSum: %d, Ratings: %+v\n",
		request.Id, request.Title, request.Description, request.UserId, request.RatingSum, request.Ratings)

	var ratingsList model.BlogRatings
	for _, rating := range request.Ratings {
		newRating := model.Rating{
			UserId:       int(rating.UserId),
			CreationDate: timestampToTime(rating.CreationDate),
			RatingValue:  int(rating.RatingValue),
		}
		ratingsList = append(ratingsList, newRating)
	}

	var blog = model.BlogPage{
		Title:        request.Title,
		Description:  request.Description,
		CreationDate: time.Now(),
		Status:       uint(request.Status),
		UserId:       int(request.UserId),
		RatingSum:    int(request.RatingSum),
		Ratings:      ratingsList,
	}

	if len(blog.Ratings) == 0 {
		blog.Ratings = make([]model.Rating, 0)
	}

	if err := h.DatabaseConnection.Table(`blog."Blogs"`).Create(&blog).Error; err != nil {
		return nil, err
	}

	responseRatings := make([]*blogs.Rating, len(blog.Ratings))
	for i, rating := range blog.Ratings {
		responseRatings[i] = &blogs.Rating{
			UserId:       int32(rating.UserId),
			CreationDate: timeToTimestampp(rating.CreationDate),
			RatingValue:  int32(rating.RatingValue),
		}
	}

	response := &blogs.Blog{
		Id:           int32(blog.Id),
		Title:        blog.Title,
		Description:  blog.Description,
		CreationDate: timeToTimestampp(blog.CreationDate),
		Status:       int32(blog.Status),
		UserId:       int32(blog.UserId),
		RatingSum:    int32(blog.RatingSum),
		Ratings:      responseRatings,
	}

	return response, nil
}

func (h BlogHandler) GetAllBlog(ctx context.Context, request *blogs.Emptyyy) (*blogs.ListBlog, error) {
	fmt.Println("Usao u getAllBlog")
	var blogsFromDB []model.BlogPage
	if err := h.DatabaseConnection.Table(`blog."Blogs"`).Find(&blogsFromDB).Error; err != nil {
		return nil, err
	}

	var protoBlogs []*blogs.Blog
	for _, blog := range blogsFromDB {
		ratingList := make([]*blogs.Rating, len(blog.Ratings))
		if len(blog.Ratings) > 0 {
			for i, rating := range blog.Ratings {
				ratingList[i] = &blogs.Rating{
					UserId:       int32(rating.UserId),
					CreationDate: timeToTimestampp(rating.CreationDate),
					RatingValue:  int32(rating.RatingValue),
				}
			}
		}

		protoBlog := &blogs.Blog{
			Id:           int32(blog.Id),
			Title:        blog.Title,
			Description:  blog.Description,
			CreationDate: timeToTimestampp(blog.CreationDate),
			Status:       int32(blog.Status),
			UserId:       int32(blog.UserId),
			RatingSum:    int32(blog.RatingSum),
			Ratings:      ratingList,
		}

		protoBlogs = append(protoBlogs, protoBlog)
	}

	response := &blogs.ListBlog{
		Blogs: protoBlogs,
	}

	return response, nil
}

func (h BlogHandler) UpdateOneBlog(ctx context.Context, request *blogs.Blog) (*blogs.Blog, error) {
	fmt.Println("Usao u updateOneBlog")
	var existingBlog model.BlogPage
	if err := h.DatabaseConnection.Table(`blog."Blogs"`).Where(`"Blogs"."Id" = ?`, request.Id).First(&existingBlog).Error; err != nil {
		return nil, err
	}

	existingBlog.Title = request.Title
	existingBlog.Description = request.Description
	existingBlog.Status = uint(request.Status)
	existingBlog.UserId = int(request.UserId)
	existingBlog.RatingSum = int(request.RatingSum)

	existingBlog.Ratings = make([]model.Rating, len(request.Ratings))
	for i, rating := range request.Ratings {

		newRating := model.Rating{
			UserId:       int(rating.UserId),
			CreationDate: timestampToTime(rating.CreationDate),
			RatingValue:  int(rating.RatingValue),
		}
		existingBlog.Ratings[i] = newRating
	}

	if err := h.DatabaseConnection.Table(`blog."Blogs"`).Save(&existingBlog).Error; err != nil {
		return nil, err
	}

	response := &blogs.Blog{
		Id:           int32(existingBlog.Id),
		Title:        existingBlog.Title,
		Description:  existingBlog.Description,
		CreationDate: timeToTimestampp(existingBlog.CreationDate),
		Status:       int32(existingBlog.Status),
		UserId:       int32(existingBlog.UserId),
		RatingSum:    int32(existingBlog.RatingSum),
		Ratings:      request.Ratings,
	}

	return response, nil
}

func (h BlogHandler) GetAllBlogsByStatus(ctx context.Context, request *blogs.GetBlogStatus) (*blogs.ListBlog, error) {
	fmt.Println("Usao u GetAllBlogsByStatus")
	state := request.State

	var blogsFromDB []model.BlogPage

	if err := h.DatabaseConnection.Table(`blog."Blogs"`).Where(`"Status" = ?`, state).Find(&blogsFromDB).Error; err != nil {
		return nil, err
	}

	protoBlogs := make([]*blogs.Blog, len(blogsFromDB))
	for i, blog := range blogsFromDB {
		ratingList := make([]*blogs.Rating, len(blog.Ratings))
		if len(blog.Ratings) > 0 {
			for j, rating := range blog.Ratings {
				ratingList[j] = &blogs.Rating{
					UserId:       int32(rating.UserId),
					CreationDate: timeToTimestampp(rating.CreationDate),
					RatingValue:  int32(rating.RatingValue),
				}
			}
		}

		protoBlogs[i] = &blogs.Blog{
			Id:           int32(blog.Id),
			Title:        blog.Title,
			Description:  blog.Description,
			CreationDate: timeToTimestampp(blog.CreationDate),
			Status:       int32(blog.Status),
			UserId:       int32(blog.UserId),
			RatingSum:    int32(blog.RatingSum),
			Ratings:      ratingList,
		}
	}

	response := &blogs.ListBlog{
		Blogs: protoBlogs,
	}

	return response, nil
}

func (h BlogHandler) UpdateRating(ctx context.Context, request *blogs.UpdateRatingRequest) (*blogs.Blog, error) {
	fmt.Println("Usao u UpdateRating")
	var existingBlog model.BlogPage
	if err := h.DatabaseConnection.Table(`blog."Blogs"`).Where(`"Id" = ?`, request.BlogId).First(&existingBlog).Error; err != nil {
		return nil, err
	}

	newRating := model.Rating{
		UserId:       int(request.UserId),
		CreationDate: time.Now(),
		RatingValue:  int(request.Value),
	}
	existingBlog.Ratings = append(existingBlog.Ratings, newRating)

	existingBlog.RatingSum += int(request.Value)

	if err := h.DatabaseConnection.Table(`blog."Blogs"`).Save(&existingBlog).Error; err != nil {
		return nil, err
	}

	responseRatings := make([]*blogs.Rating, len(existingBlog.Ratings))
	for i, rating := range existingBlog.Ratings {
		responseRatings[i] = &blogs.Rating{
			UserId:       int32(rating.UserId),
			CreationDate: timeToTimestampp(rating.CreationDate),
			RatingValue:  int32(rating.RatingValue),
		}
	}

	response := &blogs.Blog{
		Id:           int32(existingBlog.Id),
		Title:        existingBlog.Title,
		Description:  existingBlog.Description,
		CreationDate: timeToTimestampp(existingBlog.CreationDate),
		Status:       int32(existingBlog.Status),
		UserId:       int32(existingBlog.UserId),
		RatingSum:    int32(existingBlog.RatingSum),
		Ratings:      responseRatings,
	}

	return response, nil

}

func (h BlogHandler) DeleteRating(ctx context.Context, request *blogs.DeleteRatingRequest) (*blogs.Blog, error) {
	fmt.Println("Usao u DeleteRating")
	var existingBlog model.BlogPage
	if err := h.DatabaseConnection.Table(`blog."Blogs"`).Where(`"Id" = ?`, request.BlogId).First(&existingBlog).Error; err != nil {
		return nil, err
	}

	var ratingToDelete *model.Rating
	var updatedRatings []model.Rating
	for _, rating := range existingBlog.Ratings {
		if rating.UserId == int(request.UserId) {
			ratingToDelete = &rating
			existingBlog.RatingSum -= rating.RatingValue
		} else {
			updatedRatings = append(updatedRatings, rating)
		}
	}

	if ratingToDelete == nil {
		return nil, fmt.Errorf("rating not found for user %d in blog %d", request.UserId, request.BlogId)
	}

	if len(updatedRatings) == 0 {
		existingBlog.Ratings = make([]model.Rating, 0)
	} else {
		existingBlog.Ratings = updatedRatings
	}

	if err := h.DatabaseConnection.Table(`blog."Blogs"`).Save(&existingBlog).Error; err != nil {
		return nil, err
	}

	responseRatings := make([]*blogs.Rating, len(existingBlog.Ratings))
	for i, rating := range existingBlog.Ratings {
		responseRatings[i] = &blogs.Rating{
			UserId:       int32(rating.UserId),
			CreationDate: timeToTimestampp(rating.CreationDate),
			RatingValue:  int32(rating.RatingValue),
		}
	}

	response := &blogs.Blog{
		Id:           int32(existingBlog.Id),
		Title:        existingBlog.Title,
		Description:  existingBlog.Description,
		CreationDate: timeToTimestampp(existingBlog.CreationDate),
		Status:       int32(existingBlog.Status),
		UserId:       int32(existingBlog.UserId),
		RatingSum:    int32(existingBlog.RatingSum),
		Ratings:      responseRatings,
	}

	return response, nil

}
