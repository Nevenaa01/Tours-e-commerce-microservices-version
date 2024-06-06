package repository

import (
	"blogs_service/model"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type BlogRepository struct {
	DatabaseConnection *gorm.DB
}

func (repo *BlogRepository) CreateBlog(blog *model.BlogPage) error {

	dbResult := repo.DatabaseConnection.Table(`blog."Blogs"`).Create(blog)

	if dbResult.Error != nil {
		return dbResult.Error
	}
	println("Rows affected: ", dbResult.RowsAffected)
	return nil
}

func (repo *BlogRepository) GetAll() ([]model.BlogPage, error) {
	var blogs []model.BlogPage
	if err := repo.DatabaseConnection.Table(`blog."Blogs"`).Find(&blogs).Error; err != nil {
		return nil, err
	}
	return blogs, nil
}

func (repo *BlogRepository) FindByID(id int) (*model.BlogPage, error) {
	var blog model.BlogPage
	if err := repo.DatabaseConnection.Table(`blog."Blogs"`).First(&blog, id).Error; err != nil {
		return nil, err
	}
	fmt.Println(blog)
	return &blog, nil
}

func (repo *BlogRepository) UpdateOneBlog(blog *model.BlogPage) error {
	dbResult := repo.DatabaseConnection.Table(`blog."Blogs"`).Save(blog)

	if dbResult.Error != nil {
		return dbResult.Error
	}

	println("Rows affected: ", dbResult.RowsAffected)
	return nil
}

func (repo *BlogRepository) GetAllByStatus(status int) (*[]model.BlogPage, error) {
	var blogs []model.BlogPage
	if err := repo.DatabaseConnection.Table(`blog."Blogs"`).Where(`"Status" = ?`, status).Find(&blogs).Error; err != nil {
		return nil, err
	}
	return &blogs, nil
}

func (repo *BlogRepository) UpdateRating(blogId int, userId int, value int) (*model.BlogPage, error) {
	blog, err := repo.FindByID(blogId)

	if err != nil {
		return nil, err
	}

	ratings := blog.Ratings

	n := 0
	var newRatings []model.Rating
	for _, r := range ratings {
		if r.UserId != userId {
			newRatings = append(newRatings, r)
			n += r.RatingValue
		}
	}
	newRate := model.Rating{
		UserId:       userId,
		CreationDate: time.Now(),
		RatingValue:  value,
	}
	newRatings = append(newRatings, newRate)
	blog.Ratings = newRatings

	n += newRate.RatingValue
	blog.RatingSum = n

	err2 := repo.UpdateOneBlog(blog)
	if err2 != nil {
		return nil, err2
	}
	return blog, nil
}

func (repo *BlogRepository) DeleteRating(userId int, blogId int) error {
	blog, err := repo.FindByID(blogId)
	if err != nil {
		return err
	}
	ratings := blog.Ratings

	n := 0
	var newRatings []model.Rating
	for _, r := range ratings {
		if r.UserId != userId {
			newRatings = append(newRatings, r)
		}
		if r.UserId == userId {
			n = r.RatingValue
		}
	}
	blog.Ratings = newRatings
	if blog.Ratings == nil {
		blog.Ratings = make(model.BlogRatings, 0)
	}

	blog.RatingSum -= n

	err2 := repo.UpdateOneBlog(blog)
	if err2 != nil {
		return err2
	}
	return nil
}
