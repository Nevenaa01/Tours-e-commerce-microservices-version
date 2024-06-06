using AutoMapper;
using Explorer.Blog.API.Dtos;
using Explorer.Blog.API.Public;
using Explorer.Blog.Core.Domain;
using Explorer.Blog.Core.Domain.RepositoryInterfaces;
using Explorer.BuildingBlocks.Core.Domain;
using Explorer.BuildingBlocks.Core.UseCases;
using Explorer.Stakeholders.API.Internal;
using Explorer.Stakeholders.Core.Domain;
using Explorer.Stakeholders.Core.Domain.RepositoryInterfaces;
using Explorer.Tours.API.Dtos;
using FluentResults;
using Microsoft.AspNetCore.Http.HttpResults;
using Newtonsoft.Json;
using Newtonsoft.Json.Converters;
using System;
using System.Collections.Generic;
using System.Globalization;
using System.Linq;
using System.Reflection.Metadata;
using System.Runtime.InteropServices;
using System.Runtime.Serialization;
using System.Text;
using System.Text.Json;
using System.Threading.Tasks;
using static System.Reflection.Metadata.BlobBuilder;

namespace Explorer.Blog.Core.UseCases
{
    public class BlogService:CrudService<BlogDto, BlogPage>, IBlogService
    {
        public IInternalBlogService _internalBlogService;
        public IInternalCommentService _internalCommentService;
        public IBlogRepository _blogRepository;
        public ICommentService _commentService;
        private static HttpClient _httpClient;
        public BlogService(ICrudRepository<BlogPage> repository , IMapper mapper, IInternalBlogService internalBlogService, IBlogRepository blogRepository, IInternalCommentService internalCommentService, ICommentService commentService):base(repository,mapper) {
            _internalBlogService = internalBlogService;
            _internalCommentService = internalCommentService;
            _blogRepository = blogRepository;
            _commentService = commentService;
            /*_httpClient = new HttpClient()
            {
                BaseAddress = new Uri("http://blogs_service:8081")
            };*/
            _httpClient = new HttpClient()
            {
                BaseAddress = new Uri("http://api_gateway:8000")
            };
        }

        public Result<BlogDto> Get(int id)
        {
            try
            {
                var blog = _blogRepository.Get(id);
                if (blog == null)return Result.Fail(FailureCode.NotFound).WithError("Blog not found.");
                
                var user = _internalBlogService.GetByUserId(blog.UserId);
                if (user == null)return Result.Fail(FailureCode.NotFound).WithError("User not found.");
                
                var result= MapToDto(blog);
                result.Username = user.Username;
                return result;
            }
            catch (KeyNotFoundException e)
            {
                return Result.Fail(FailureCode.NotFound).WithError(e.Message);
            }
        }

        public Result<List<BlogDto>> GetAll()
        {
            try
            {
                var blogs = MapToDto(_blogRepository.GetAll());

                foreach (var blog in blogs.Value)
                {
                    var user = _internalBlogService.GetByUserId(blog.UserId);
                    blog.Username = user.Username;
                }

                return blogs;
            }
            catch (KeyNotFoundException e)
            {
                return Result.Fail(FailureCode.NotFound).WithError(e.Message);
            }
        }

        public Result<CommentDto> CreateComment(CommentDto comment)
        {
            var result = _commentService.Create(comment);
            return result;
        }

        public Result<CommentDto> UpdateComment(CommentDto comment)
        {
            var result = _commentService.Update(comment);
            return result;
        }

        public Result DeleteComment(int id)
        {
            var result = _commentService.Delete(id);
            return result;
        }

        public Result<CommentDto> GetComment(int id)
        {
            var result = _commentService.Get(id);
            if (result.IsSuccess == true)
            {
                var user = _internalBlogService.GetByUserId(result.Value.UserId);
                result.Value.Username = user.Username;
                var person = _internalCommentService.GetByUserId(user.Id);
                result.Value.ProfilePic = person.ProfilePic;
            }
            return result;
        }

        public Result<PagedResult<CommentDto>> GetPagedComments(int page, int pageSize)
        {
            var result = _commentService.GetPaged(page, pageSize);
            return result;
        }

        public Result<List<CommentDto>> GetCommentsByBlogId(int blogId)
        {
            var result = _commentService.GetCommentsByBlogId(blogId);

            foreach (var comment in result.Value)
            {
                var user = _internalBlogService.GetByUserId(comment.UserId);
                comment.Username = user.Username;
                var person = _internalCommentService.GetByUserId(user.Id);
                comment.ProfilePic = person.ProfilePic;
            }

            return result;
        }

        public Result DeleteRating(int blogId,int userId)
        {
            var result=_blogRepository.DeleteRating(userId, blogId);

            return result;
        }

        public Result<BlogDto> UpdateRating(int blogId,int userId,int value)
        {
            var result = _blogRepository.UpdateRating(blogId,userId,value);

            return MapToDto(result);
        }

        public Result<List<BlogDto>> GetBlogsByStatus(int state)
        {
            var blogs = _blogRepository.GetBlogsByStatus(state);
            if(blogs==null)return Result.Fail(FailureCode.NotFound).WithError("Blog not found.");  

            var result= MapToDto(blogs);

            foreach (var blog in result.Value)
            {
                var user = _internalBlogService.GetByUserId(blog.UserId);
                if (user == null) return Result.Fail(FailureCode.NotFound).WithError("User not found.");
                blog.Username = user.Username;
            }

            return result;
        }

        public Result<List<BlogDto>> GetBlogsByAuthor(int authorId)
        {
            var blogs = _blogRepository.GetBlogsByAuthor(authorId);
            if (blogs == null) return Result.Fail(FailureCode.NotFound).WithError("Blog not found.");

            var result = MapToDto(blogs);

            foreach (var blog in result.Value)
            {
                var user = _internalBlogService.GetByUserId(blog.UserId);
                if (user == null) return Result.Fail(FailureCode.NotFound).WithError("User not found.");
                blog.Username = user.Username;
            }

            return result;
        }


        public async Task<Result<BlogDto>> CreateBlogAsync(BlogDto blog)
        {
            var jsonString = System.Text.Json.JsonSerializer.Serialize(blog);
            Console.WriteLine("Serialized JSON: " + jsonString);

            using StringContent jsonContent = new(
                jsonString,
                Encoding.UTF8,
                "application/json"
            );

            using HttpResponseMessage response = await _httpClient.PostAsync("/blogs", jsonContent);

            response.EnsureSuccessStatusCode();
            var jsonResponse = await response.Content.ReadAsStringAsync();
            var blogDto = JsonConvert.DeserializeObject<BlogDto>(jsonResponse);

            return blogDto;

        }

        public async Task<Result<List<BlogDto>>> GetAllBlogsAsync()
        {

            using HttpResponseMessage response = await _httpClient.GetAsync("/blogs");
            response.EnsureSuccessStatusCode();
            
            var jsonResponse = await response.Content.ReadAsStringAsync();
 
            var trimmedJsonResponse = jsonResponse.Replace("{\"Blogs\":", "").TrimEnd('}');

            var settings = new JsonSerializerSettings
            {
                Converters = new List<JsonConverter> { new CustomDateTimeConverter() }
            };

            var blogDtos = JsonConvert.DeserializeObject<List<BlogDto>>(trimmedJsonResponse,settings);

            var listResult = new List<BlogDto>(blogDtos);

            foreach (var blog in listResult)
            {
                var user = _internalBlogService.GetByUserId(blog.UserId);
                blog.Username = user.Username;
            }


            
            return listResult;

        }

        public async Task<Result<BlogDto>> GetBlogByIdAsync(int id)
        {
            
            using HttpResponseMessage response = await _httpClient.GetAsync("/blogs/" + id.ToString());
            response.EnsureSuccessStatusCode();

            var jsonResponse = await response.Content.ReadAsStringAsync();
            var blogDto = JsonConvert.DeserializeObject<BlogDto>(jsonResponse);
            
            if(blogDto != null)
            {
                var user = _internalBlogService.GetByUserId(blogDto.UserId);
                blogDto.Username = user.Username;
                return blogDto;
            }
            return blogDto;    
            
        }

        public async Task<Result<BlogDto>> UpdateBlogAsync(BlogDto blog)
        {
            using StringContent jsonContent = new(System.Text.Json.JsonSerializer.Serialize(blog), Encoding.UTF8, "application/json");
            using HttpResponseMessage response = await _httpClient.PutAsync("/blogs/updateOneBlog", jsonContent);        
            response.EnsureSuccessStatusCode();

            var jsonResponse = await response.Content.ReadAsStringAsync();
            var blogDto = JsonConvert.DeserializeObject<BlogDto>(jsonResponse);

            return blogDto;
        }

        public async Task<Result<List<BlogDto>>> GetBlogsByStatusAsync(int state)
        {
            using HttpResponseMessage response = await _httpClient.GetAsync("/blogs/getByStatus/" + state.ToString());
            response.EnsureSuccessStatusCode();

            var jsonResponse = await response.Content.ReadAsStringAsync();

            var trimmedJsonResponse = jsonResponse.Replace("{\"Blogs\":", "").TrimEnd('}');

            var settings = new JsonSerializerSettings
            {
                Converters = new List<JsonConverter> { new CustomDateTimeConverter() }
            };

            var blogDtos = JsonConvert.DeserializeObject<List<BlogDto>>(trimmedJsonResponse, settings);

            var listResult = new List<BlogDto>(blogDtos);

            foreach (var blog in listResult)
            {
                var user = _internalBlogService.GetByUserId(blog.UserId);
                blog.Username = user.Username;
            }
            return listResult;
        }

        public async Task<Result<BlogDto>> UpdateRatingAsync(int blogId, int userId, int value)
        {
            
            using HttpResponseMessage response = await _httpClient.PutAsync("/blogs/rating/" + userId.ToString() + "/" + blogId.ToString() + "/" +  value.ToString(),null);
            response.EnsureSuccessStatusCode();

            var jsonResponse = await response.Content.ReadAsStringAsync();
            var blogDto = JsonConvert.DeserializeObject<BlogDto>(jsonResponse);

            return blogDto;
        }

        public async Task<Result> DeleteRatingAsync(int userId, int blogId)
        {

            using HttpResponseMessage response = await _httpClient.DeleteAsync("/blogs/rating/" + userId.ToString() + "/" + blogId.ToString());
            response.EnsureSuccessStatusCode();

            var jsonResponse = await response.Content.ReadAsStringAsync();
            var blogDto = JsonConvert.DeserializeObject<BlogDto>(jsonResponse);

            return Result.Ok();
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

        public async Task<Result<List<CommentDto>>> GetCommentsByBlogIdAsync(int blogId)
        {
            using HttpResponseMessage response = await _httpClient.GetAsync("/comment/" + blogId.ToString());
            response.EnsureSuccessStatusCode();

            var jsonResponse = await response.Content.ReadAsStringAsync();

            var trimmedJsonResponse = jsonResponse.Replace("{\"Comments\":", "").TrimEnd('}');

            var settings = new JsonSerializerSettings
            {
                Converters = new List<JsonConverter> { new CustomDateTimeConverter() }
            };

            var commentDtos = JsonConvert.DeserializeObject<List<CommentDto>>(trimmedJsonResponse, settings);


            var listResult = new List<CommentDto>(commentDtos);

            foreach (var comment in listResult)
            {
                var user = _internalBlogService.GetByUserId(comment.UserId);
                comment.Username = user.Username;
                var person = _internalCommentService.GetByUserId(user.Id);
                comment.ProfilePic = person.ProfilePic;
            }
            return listResult;
        }
    }
}

public class CustomDateTimeConverter : DateTimeConverterBase
{
    public override object? ReadJson(JsonReader reader, Type objectType, object? existingValue, Newtonsoft.Json.JsonSerializer serializer)
    {
        if (reader.Value == null)
            return DateTime.MinValue;

        // Attempt to parse the value as a DateTimeOffset first
        if (DateTimeOffset.TryParse(reader.Value.ToString(), out DateTimeOffset dateTimeOffset))
            return dateTimeOffset.DateTime;

        // If parsing as DateTimeOffset fails, try using specific formats
        string[] formats =
        {
            "yyyy-MM-dd HH:mm:ss.fff zzz",
            "yyyy-MM-dd HH:mm:ss.fff z",
            "yyyy-MM-dd HH:mm:ss.fff zz",
            "yyyy-MM-dd HH:mm:ss.fff zzzz",
            "M/d/yyyy h:mm:ss tt",
            "yyyy-MM-ddTHH:mm:ss.fffffffZ",
            "yyyy-MM-ddTHH:mm:ss.fffZ",
            "yyyy-MM-ddTHH:mm:ss.fffffffK",
            "yyyy-MM-ddTHH:mm:ss.fffK"
        };

        if (DateTime.TryParseExact(reader.Value.ToString(), formats, CultureInfo.InvariantCulture, DateTimeStyles.AdjustToUniversal, out DateTime result))
            return result;

        throw new Newtonsoft.Json.JsonException($"Could not convert string to DateTime: {reader.Value}");
    }

    public override void WriteJson(JsonWriter writer, object? value, Newtonsoft.Json.JsonSerializer serializer)
    {
        writer.WriteValue(((DateTime)value).ToString("yyyy-MM-dd HH:mm:ss.fffK"));
    }
}