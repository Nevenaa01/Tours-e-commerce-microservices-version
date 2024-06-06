using AutoMapper;
using Explorer.Blog.API.Dtos;
using Explorer.Blog.API.Public;
using Explorer.Blog.Core.Domain;
using Explorer.Blog.Core.Domain.RepositoryInterfaces;
using Explorer.BuildingBlocks.Core.UseCases;
using FluentResults;
using Newtonsoft.Json;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Net.Http;
using System.Text;
using System.Threading.Tasks;

namespace Explorer.Blog.Core.UseCases
{
    public class CommentService : CrudService<CommentDto, Comment>, ICommentService
    {
        public ICommentRepository _commentRepository;
        private static HttpClient _httpClient;

        public CommentService(ICrudRepository<Comment> repository, IMapper mapper, ICommentRepository commentRepository) : base(repository, mapper)
        {
            _commentRepository = commentRepository;
            _httpClient = new HttpClient()
            {
                BaseAddress = new Uri("http://api_gateway:8000")
            };
        }

        public Result<List<CommentDto>> GetCommentsByBlogId(int blogId)
        {
            return MapToDto(_commentRepository.GetCommentsByBlogId(blogId));
        }

        public async Task<Result<CommentDto>> CreateCommentAsync(CommentDto comment)
        {
            using StringContent jsonContent = new(System.Text.Json.JsonSerializer.Serialize(comment), Encoding.UTF8, "application/json");

            using HttpResponseMessage response = await _httpClient.PostAsync("/comment", jsonContent);

            response.EnsureSuccessStatusCode();
            var jsonResponse = await response.Content.ReadAsStringAsync();
            var commentDto = JsonConvert.DeserializeObject<CommentDto>(jsonResponse);

            return commentDto;
        }

        public async Task<Result> DeleteCommentAsync(int commentId)
        {
            using HttpResponseMessage response = await _httpClient.DeleteAsync("/comment/" + commentId.ToString());
            response.EnsureSuccessStatusCode();

            var jsonResponse = await response.Content.ReadAsStringAsync();
            var commentDto = JsonConvert.DeserializeObject<CommentDto>(jsonResponse);

            return Result.Ok();
        }

        public async Task<Result<CommentDto>> UpdateCommentAsync(CommentDto commentDto)
        {
            using StringContent jsonContent = new(System.Text.Json.JsonSerializer.Serialize(commentDto), Encoding.UTF8, "application/json");
            using HttpResponseMessage response = await _httpClient.PutAsync("/comment", jsonContent);
            response.EnsureSuccessStatusCode();

            var jsonResponse = await response.Content.ReadAsStringAsync();
            var comment = JsonConvert.DeserializeObject<CommentDto>(jsonResponse);

            return comment;
        }
    }
}
