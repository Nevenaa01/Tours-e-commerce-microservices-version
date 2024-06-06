using Explorer.Blog.API.Dtos;
using Explorer.BuildingBlocks.Core.Domain;
using Explorer.BuildingBlocks.Core.UseCases;
using FluentResults;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace Explorer.Blog.API.Public
{
    public interface IBlogService
    {
        Result<PagedResult<BlogDto>> GetPaged(int page, int pageSize);
        Result<BlogDto> Create(BlogDto blog);
        Result<BlogDto> Update(BlogDto blog);
        Result Delete(int id);
        Result<BlogDto> Get(int id);
        Result<CommentDto> CreateComment(CommentDto comment);
        Result<CommentDto> UpdateComment(CommentDto comment);
        Result DeleteComment(int id);
        Result<CommentDto> GetComment(int id);
        Result<PagedResult<CommentDto>> GetPagedComments(int page, int pageSize);
        Result<List<CommentDto>> GetCommentsByBlogId(int blogId);
        Result<List<BlogDto>> GetAll();
        Result DeleteRating(int blogId, int userId);
        Result<BlogDto> UpdateRating(int blogId, int userId,int value);
        Result<List<BlogDto>> GetBlogsByStatus(int state);
        Result<List<BlogDto>> GetBlogsByAuthor(int authorId);


        Task<Result<BlogDto>> CreateBlogAsync(BlogDto blog);
        Task<Result<List<BlogDto>>> GetAllBlogsAsync();
        Task<Result<BlogDto>> GetBlogByIdAsync(int id);
        Task<Result<BlogDto>> UpdateBlogAsync(BlogDto blog);
        Task<Result<List<BlogDto>>> GetBlogsByStatusAsync(int state);
        Task<Result<BlogDto>> UpdateRatingAsync(int blogId, int userId, int value);
        Task<Result> DeleteRatingAsync(int userId, int blogId);

        Task<Result<CommentDto>> CreateCommentAsync(CommentDto commentDto);
        Task<Result<CommentDto>> UpdateCommentAsync(CommentDto commentDto);
        Task<Result> DeleteCommentAsync(int commentId);
    //    Task<Result<CommentDto>> GetCommentAsync(int id);
        Task<Result<List<CommentDto>>> GetCommentsByBlogIdAsync(int blogId);
    }
}
