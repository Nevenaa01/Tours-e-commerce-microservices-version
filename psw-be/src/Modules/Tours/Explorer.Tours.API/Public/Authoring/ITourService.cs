using Explorer.BuildingBlocks.Core.UseCases;
using Explorer.Tours.API.Dtos;
using FluentResults;
using Microsoft.AspNetCore.Mvc.RazorPages;
using System.Xml.Linq;

namespace Explorer.Tours.API.Public.Authoring
{
    public interface ITourService
    {
        Result<PagedResult<TourDto>> GetPaged(int page, int pageSize);
        Result<TourDto> Create(TourDto tour);
        Task<Result<TourDto>> CreateAsync(TourDto tour);
        Result<TourDto> Update(TourDto tour);
        Task<Result<TourDto>> UpdateAsync(TourDto tour);
        Result Delete(int id);
        Result<TourDto> Get(int id);
        Result<TourDto> Publish(int id, int userId);
        Result<TourDto> Archive(int id, int userId);
        Result<PagedResult<TourDto>> GetPagedByAuthorId(int authorId, int page, int pageSize);
        Result<PagedResult<TourDto>> GetPagedForSearch(string name, string[] tags, int page, int pageSize);
        Result<TourDto> CreateCampaign(List<TourDto> tours, string name, string description, int touristId);
        Result<PagedResult<TourDto>> GetPagedForSearchByLocation(int page, int pageSize, int touristId);
        List<TourDto> GetAllByAuthorId(int authorId);
        Task<string> UpdateAsync(TourDto tour, HttpClient _httpClient);
        Task<string> PublishAsync(int id, int userId, HttpClient _httpClient);
        Task<string> ArchiveAsync(int id, int userId, HttpClient _httpClient);
        Task<string> CreateAsync(TourDto tour, HttpClient _httpClient);

        Task<Result<List<TourDto>>> GetPagedByAuthorIdAsync(int authorId);
        Task<Result<TourDto>> GetAsync(int id);

        Task<Result<PagedResult<TourDto>>> GetAllAsync();
    }
}
