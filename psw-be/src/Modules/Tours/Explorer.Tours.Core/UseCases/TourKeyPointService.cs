using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using AutoMapper;
using Explorer.BuildingBlocks.Core.UseCases;
using Explorer.Tours.API.Dtos;
using Explorer.Tours.API.Public;
using Explorer.Tours.Core.Domain.RepositoryInterfaces;
using Explorer.Tours.Core.Domain.Tours;
using FluentResults;
using System.Text.Json;
using Microsoft.AspNetCore.DataProtection;
using Newtonsoft.Json;
using System.Net.Http;
using Newtonsoft.Json.Linq;

namespace Explorer.Tours.Core.UseCases
{
    public class TourKeyPointService : CrudService<TourKeyPointDto, TourKeyPoint>, ITourKeyPointService
    {
        private readonly ITourKeyPointsRepository _tourKeyPointsRepository;

        public TourKeyPointService(ICrudRepository<TourKeyPoint> repository, IMapper mapper, ITourKeyPointsRepository tourKeyPointsRepository) : base(repository, mapper)
        {
            _tourKeyPointsRepository = tourKeyPointsRepository;
        }

        private static HttpClient _httpClient = new()
        {
            BaseAddress = new Uri("http://api_gateway:8000")
        };
        public async Task<Result<TourKeyPointDto>> CreateAsync(TourKeyPointDto tourKeypointDto)
        {
            using StringContent jsonContent = new(System.Text.Json.JsonSerializer.Serialize(tourKeypointDto), Encoding.UTF8, "application/json");
            using HttpResponseMessage response = await _httpClient.PostAsync("/createTourKeypoint", jsonContent);
            response.EnsureSuccessStatusCode();
            var jsonString = await response.Content.ReadAsStringAsync();
            var jsonObject = JsonDocument.Parse(jsonString).RootElement;

            var imageString = jsonObject.GetProperty("Image").GetString();
            Uri imageUri = null;

            if (Uri.TryCreate(imageString, UriKind.Absolute, out var uriResult) &&
                (uriResult.Scheme == Uri.UriSchemeHttp || uriResult.Scheme == Uri.UriSchemeHttps))
            {
                imageUri = uriResult;
            }

            TourKeyPointDto tourKeyPointDto = new TourKeyPointDto
            {
                Id = jsonObject.GetProperty("Id").GetInt32(),
                Name = jsonObject.GetProperty("Name").GetString(),
                Description = jsonObject.GetProperty("Description").GetString(),
                Latitude = jsonObject.GetProperty("Latitude").GetDouble(),
                Longitude = jsonObject.GetProperty("Longitude").GetDouble(),
                TourId = jsonObject.GetProperty("TourId").GetInt32(),
                PositionInTour = jsonObject.GetProperty("PositionInTour").GetInt32(),
                Secret = jsonObject.GetProperty("Secret").GetString(),
                Image = imageUri,
            };

            return tourKeyPointDto;
        }
        public Result<List<TourKeyPointDto>> GetAllByPublicKeypointId(long publicId)
        {
            List<TourKeyPointDto> tourKeyPointDtos = new List<TourKeyPointDto>();
            var tourKeyPoints = _tourKeyPointsRepository.GetAllByPublicId(publicId);
            foreach (var tourKeyPoint in tourKeyPoints)
            {
                TourKeyPointDto tourKeyPointDto = new TourKeyPointDto
                {
                    Id = (int)tourKeyPoint.Id,
                    Name = tourKeyPoint.Name,
                    Description = tourKeyPoint.Description,
                    Image = tourKeyPoint.Image,
                    Latitude = tourKeyPoint.Latitude,
                    Longitude = tourKeyPoint.Longitude,
                    TourId = tourKeyPoint.TourId,
                    PositionInTour = tourKeyPoint.PositionInTour,
                    PublicPointId = tourKeyPoint.PublicPointId
                };
                tourKeyPointDtos.Add(tourKeyPointDto);
            }

            return tourKeyPointDtos;
        }

        public Result<List<TourKeyPointDto>> GetByTourId(long tourId)
        {
            List<TourKeyPointDto> tourKeyPointDtos = new List<TourKeyPointDto>();
           var tourKeyPoints = _tourKeyPointsRepository.GetByTourId(tourId);
           foreach (var tourKeyPoint in tourKeyPoints)
           {
               TourKeyPointDto tourKeyPointDto = new TourKeyPointDto
               {
                   Id = (int)tourKeyPoint.Id,
                   Name = tourKeyPoint.Name,
                   Description = tourKeyPoint.Description,
                   Image = tourKeyPoint.Image,
                   Latitude = tourKeyPoint.Latitude,
                   Longitude = tourKeyPoint.Longitude,
                   TourId = tourKeyPoint.TourId
               };
                tourKeyPointDtos.Add(tourKeyPointDto);
           }

           return tourKeyPointDtos;
        }

        public async Task<Result<List<TourKeyPointDto>>> GetByTourIdAsync(int tourId)
        {
            using HttpResponseMessage response = await _httpClient.GetAsync("/getByTourId/" + tourId.ToString());
            response.EnsureSuccessStatusCode();

            var jsonResponse = await response.Content.ReadAsStringAsync();

            var trimmedJsonResponse = jsonResponse.Replace("{\"results\":", "").TrimEnd('}');


            var tourKeypointDtos = JsonConvert.DeserializeObject<List<TourKeyPointDto>>(trimmedJsonResponse);

            var listResult = new List<TourKeyPointDto>(tourKeypointDtos);


            return listResult;

        }

    }
}
