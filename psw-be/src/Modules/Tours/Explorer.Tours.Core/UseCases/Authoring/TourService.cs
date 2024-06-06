using AutoMapper;
using Explorer.BuildingBlocks.Core.UseCases;
using Explorer.Payments.API.Dtos;
using Explorer.Payments.API.Internal;
using Explorer.Tours.API.Dtos;
using Explorer.Tours.API.Internal;
using Explorer.Tours.API.Public;
using Explorer.Tours.API.Public.Authoring;
using Explorer.Tours.Core.Domain;
using Explorer.Tours.Core.Domain.RepositoryInterfaces;
using Explorer.Tours.Core.Domain.Tours;
using FluentResults;
using Microsoft.AspNetCore.Mvc.RazorPages;
using Microsoft.AspNetCore.Components.Web;
using Microsoft.AspNetCore.SignalR;
using System.Dynamic;
using System.Linq;
using System.Xml.Linq;
using System.Text.Json;
using System.Text;
using Newtonsoft.Json;
using Explorer.Stakeholders.Core.Domain;
using Newtonsoft.Json.Converters;
using System.Globalization;
using System.Reflection.Metadata;
using Newtonsoft.Json.Linq;

namespace Explorer.Tours.Core.UseCases.Authoring
{
    public class TourService : CrudService<TourDto, Tour>, ITourService, IInternalTourService
    {
        private readonly ITourRepository _tourRepository;
        private readonly ITourKeyPointService _keyPointService;
        private readonly IInternalBoughtItemService _internalBoughtItemService;
        private readonly IInternalPersonService _internalPersonService;
        private readonly HttpClient _httpClient;

        public TourService(ITourRepository repository, IMapper mapper, ITourKeyPointService tourKeyPointService,
            IInternalBoughtItemService internalBoughtItemService,
            IInternalPersonService internalPersonService, HttpClient httpClient) : base(repository, mapper)
        {
            _tourRepository = repository;
            _keyPointService = tourKeyPointService;
            _internalBoughtItemService = internalBoughtItemService;
            _internalPersonService = internalPersonService;
            _httpClient = httpClient;
        }

        /*private static HttpClient sharedClient = new()
        {
            BaseAddress = new Uri("http://tours_service:8080")
        };*/
        private static HttpClient sharedClient = new()
        {
            BaseAddress = new Uri("http://api_gateway:8000")
        };

        public async Task<Result<TourDto>> CreateAsync(TourDto tour)
        {
            using StringContent jsonContent = new(System.Text.Json.JsonSerializer.Serialize(tour), Encoding.UTF8, "application/json");

            using HttpResponseMessage response = await sharedClient.PostAsync("/createTour", jsonContent);

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

            TourDto tourDto = new TourDto
            {
                Id = jsonObject.GetProperty("Id").GetInt32(),
                Name = jsonObject.GetProperty("Name").GetString(),
                Description = jsonObject.GetProperty("Description").GetString(),
                Difficulty = jsonObject.GetProperty("Difficulty").GetInt32(),
                Tags = ParseTags(jsonObject.GetProperty("Tags")),
                Status = jsonObject.GetProperty("Status").GetInt32(),
                Price = jsonObject.GetProperty("Price").GetDouble(),
                AuthorId = jsonObject.GetProperty("AuthorId").GetInt32(),
                Equipment = ParseEquipment(jsonObject.GetProperty("Equipment")),
                DistanceInKm = jsonObject.GetProperty("DistanceInKm").GetDouble(),
                ArchivedDate = jsonObject.TryGetProperty("ArchivedDate", out var archivedDate) && archivedDate.ValueKind != JsonValueKind.Null ?
        (DateTime?)DateTime.Parse(archivedDate.GetString()) : null,
                PublishedDate = jsonObject.TryGetProperty("PublishedDate", out var publishedDate) && publishedDate.ValueKind != JsonValueKind.Null ?
        (DateTime?)DateTime.Parse(publishedDate.GetString()) : null,
                Durations = { },
                KeyPoints = { },
                Image = imageUri, 
            };

            return tourDto;
        }

        private List<string> ParseTags(JsonElement tagsElement)
        {
            List<string> tags = new List<string>();
            if (tagsElement.ValueKind == JsonValueKind.Array)
            {
                foreach (var tag in tagsElement.EnumerateArray())
                {
                    tags.Add(tag.GetString());
                }
            }
            return tags;
        }

        private int[] ParseEquipment(JsonElement equipmentElement)
        {
            List<int> equipmentList = new List<int>();
            if (equipmentElement.ValueKind == JsonValueKind.Array)
            {
                foreach (var item in equipmentElement.EnumerateArray())
                {
                    if (item.ValueKind == JsonValueKind.Number)
                    {
                        equipmentList.Add(item.GetInt32());
                    }
                }
            }
            return equipmentList.ToArray();
        }

        public async Task<Result<TourDto>> UpdateAsync(TourDto tour)
        {
            using StringContent jsonContent = new(System.Text.Json.JsonSerializer.Serialize(tour), Encoding.UTF8, "application/json");
            using HttpResponseMessage response = await sharedClient.PutAsync("/updateTour", jsonContent);
            response.EnsureSuccessStatusCode();

            var jsonResponse = await response.Content.ReadAsStringAsync();
            var tourDto = JsonConvert.DeserializeObject<TourDto>(jsonResponse);

            return tourDto;
        }


        public async Task<string> ArchiveAsync(int id, int userId, HttpClient _httpClient)
        {
            using HttpResponseMessage responseGet = await sharedClient.GetAsync("/tours/" + id.ToString());
            responseGet.EnsureSuccessStatusCode();
            var jsonResponseGet = await responseGet.Content.ReadAsStringAsync();
            var tour = Newtonsoft.Json.JsonConvert.DeserializeObject<Tour>(jsonResponseGet);

            tour.Id = id;
            tour.Archive(userId);
            using StringContent jsonContent = new(System.Text.Json.JsonSerializer.Serialize(tour), Encoding.UTF8, "application/json");
            using HttpResponseMessage response = await sharedClient.PutAsync("/updateTour", jsonContent);
            response.EnsureSuccessStatusCode();
            var jsonResponse = await response.Content.ReadAsStringAsync();
            return jsonResponse;
        }
        public Result<TourDto> Archive(int id, int userId)
        {
            try
            {
                var tour = _tourRepository.Get(id);
                tour.Archive(userId);
                _tourRepository.SaveChanges();
                return MapToDto(tour);
            }
            catch (ArgumentException e)
            {
                return Result.Fail(FailureCode.InvalidArgument).WithError(e.Message);
            }
            catch (UnauthorizedAccessException e)
            {
                return Result.Fail(FailureCode.Forbidden).WithError(e.Message);
            }
        }
        
        public async Task<string> PublishAsync(int id,int userId, HttpClient _httpClient)
        {
            using HttpResponseMessage responseGet = await sharedClient.GetAsync("/tours/" + id.ToString());
            responseGet.EnsureSuccessStatusCode();

            var jsonResponseGet = await responseGet.Content.ReadAsStringAsync();
            var tour= Newtonsoft.Json.JsonConvert.DeserializeObject<Tour>(jsonResponseGet);

            tour.Id = id;
            tour.Publish(userId);
            using StringContent jsonContent = new(System.Text.Json.JsonSerializer.Serialize(tour), Encoding.UTF8, "application/json");
            using HttpResponseMessage response = await sharedClient.PutAsync("/updateTour", jsonContent);
            response.EnsureSuccessStatusCode();
            var jsonResponse = await response.Content.ReadAsStringAsync();
            return jsonResponse;
        }

        public Result<TourDto> Publish(int id, int userId)
        {
            try
            {
                var tour = _tourRepository.Get(id);
                tour.Publish(userId);
                _tourRepository.SaveChanges();
                return MapToDto(tour);
            }
            catch (ArgumentException e)
            {
                return Result.Fail(FailureCode.InvalidArgument).WithError(e.Message);
            }
            catch (UnauthorizedAccessException e)
            {
                return Result.Fail(FailureCode.Forbidden).WithError(e.Message);
            }
        }

        public Result<PagedResult<TourDto>> GetPagedByAuthorId(int authorId, int page, int pageSize)
        {
            var result = _tourRepository.GetPagedByAuthorId(authorId, page, pageSize);
            return MapToDto(result);
        }
        public async Task<Result<List<TourDto>>> GetPagedByAuthorIdAsync(int authorId)
        {


            using HttpResponseMessage response = await sharedClient.GetAsync("/tours/author/" + authorId.ToString());
            response.EnsureSuccessStatusCode();
            var jsonResponseGet = await response.Content.ReadAsStringAsync();

            var trimmedJsonResponse = jsonResponseGet.Replace("{\"tour\":", "").TrimEnd('}');

            var settings = new JsonSerializerSettings
            {
                Converters = new List<JsonConverter> { new CustomDateTimeConverter() }
            };

            var tourDtos = JsonConvert.DeserializeObject<List<TourDto>>(trimmedJsonResponse, settings);

            var listResult = new List<TourDto>(tourDtos);


            return tourDtos;
        }

        public Result<TourDto> CreateCampaign(List<TourDto> tours, string name, string description, int touristId)
        {
            if (tours.Count < 2)
            {
                throw new ArgumentException("In order to create campaign, atleast 2 Tours have to be picked.");
            }

            int difficulty = 0;
            List<string> tags = new List<string>();
            TourStatus status = TourStatus.Draft;
            double price = 0;
            int[] equipment = { };
            double distanceInKm = 0;
            DateTime? archivedDate = null;
            DateTime? publishedDate = null;
            List<TourDuration> durations = new List<TourDuration>();

            int counter = 0;

            foreach (var tour in tours)
            {
                distanceInKm += tour.DistanceInKm;
                difficulty += tour.Difficulty;
                tags = CombineCampaignTags(tags, tour.Tags);
                equipment = CombineCampaignEquipment(equipment, tour.Equipment);
                counter++;
            }

            TourDifficulty difficultyTemp = GetCampaignDifficulty(difficulty, counter);

            Tour campaign = new Tour(name, description, difficultyTemp, tags, status, price, touristId, equipment,
                distanceInKm, archivedDate, publishedDate, durations);

            var createdCampaign = Create(MapToDto(campaign));

            CreateDuplicateKeypoints(tours, createdCampaign.Value.Id);
            CreateBoughtItemForCampaign(touristId, createdCampaign.Value.Id);

            return MapToDto(campaign);
        }

        public void CreateBoughtItemForCampaign(int touristId, int campaignId)
        {
            BoughtItemDto boughtItemDto = new BoughtItemDto();
            boughtItemDto.UserId = touristId;
            boughtItemDto.TourId = campaignId;
            boughtItemDto.IsUsed = false;

            _internalBoughtItemService.CreateBoughtItem(boughtItemDto);
        }

        public void CreateDuplicateKeypoints(List<TourDto> tours, int campaignId)
        {
            List<TourKeyPointDto> keypoints = new List<TourKeyPointDto>();
            int positionInCampaign = 1;
            foreach (var tour in tours)
            {
                keypoints = _keyPointService.GetByTourId(tour.Id).Value;
                keypoints = keypoints.OrderBy(kp => kp.PositionInTour).ToList();
                foreach (var keypoint in keypoints)
                {
                    keypoint.PositionInTour = positionInCampaign;
                    keypoint.TourId = campaignId;
                    keypoint.Id = 0;
                    _keyPointService.Create(keypoint);
                    positionInCampaign++;
                }
            }
        }

        public int CalculateCampaignDifficulty(int difficultySum, int counter)
        {
            return difficultySum / counter;
        }

        public List<string> CombineCampaignTags(List<string> campaignTags, List<string> tourTags)
        {
            foreach (var tag in tourTags)
            {
                campaignTags.Add(tag);
            }

            return campaignTags;
        }

        public TourDifficulty GetCampaignDifficulty(int difficulty, int counter)
        {
            difficulty = CalculateCampaignDifficulty(difficulty, counter);

            if (difficulty == 0)
            {
                return TourDifficulty.Beginner;
            }
            else if (difficulty == 1)
            {
                return TourDifficulty.Intermediate;
            }
            else if (difficulty == 2)
            {
                return TourDifficulty.Advanced;
            }
            else
            {
                return TourDifficulty.Pro;
            }
        }

        public int[] CombineCampaignEquipment(int[] campaignEquipment, int[] tourEquipment)
        {
            List<int> campaignEquipmentList = new List<int>(campaignEquipment);

            foreach (var equipmentId in tourEquipment)
            {
                if (!campaignEquipmentList.Contains(equipmentId))
                {
                    campaignEquipmentList.Add(equipmentId);
                }
            }

            return campaignEquipmentList.ToArray();
        }


        public Result<PagedResult<TourDto>> GetPagedForSearch(string name, string[] tags, int page, int pageSize)
        {
            var tours = _tourRepository.GetPaged(page, pageSize);
            PagedResult<TourDto> filteredTours = new PagedResult<TourDto>(new List<TourDto>(), 0);

            if (tags[0].Contains(","))
            {
                tags = tags[0].Split(",");
            }

            filteredTours.Results.AddRange(
                tours.Results
                    .Where(tour => (tour.Name.ToLower().Contains(name.ToLower()) || name.Equals("empty")) &&
                                   tags.All(tag => tour.Tags.Any(tourTag =>
                                       tourTag.ToLower() == tag.ToLower() || tag == "empty")))
                    .Select(MapToDto)
            );

            return filteredTours;
        }
        public async Task<string> CreateAsync(TourDto tour, HttpClient _httpClient)
        {
            using StringContent jsonContent = new(System.Text.Json.JsonSerializer.Serialize(tour), Encoding.UTF8, "application/json");
            using HttpResponseMessage response = await sharedClient.PostAsync("/tours", jsonContent);
            response.EnsureSuccessStatusCode();
            var jsonResponse = await response.Content.ReadAsStringAsync();
            return jsonResponse;
        }
        public Result<PagedResult<TourDto>> GetPagedForSearchByLocation(int page, int pageSize, int touristId)
        {
            double radius = 40000;
            var person = _internalPersonService.Get(touristId);
            var tours = _tourRepository.GetPaged(page, pageSize);
            var publishedTours = tours.Results.Where(tour => tour.Status == Domain.Tours.TourStatus.Published).ToList();
            var resultTours = new PagedResult<TourDto>(new List<TourDto>(), 0);

            if (person.Value.Latitude == null && person.Value.Longitude == null)
            {
                foreach (Tour tour in publishedTours)
                {
                    resultTours.Results.Add(MapToDto(tour));
                }

                return resultTours;
            }

            PagedResult<TourDto> filteredTours = new PagedResult<TourDto>(new List<TourDto>(), 0);
            foreach (var tour in tours.Results)
            {
                if (tour.Status == TourStatus.Published && CheckIfAnyKeyPointInRange(tour.KeyPoints,
                        person.Value.Latitude, person.Value.Longitude, radius))
                {
                    filteredTours.Results.Add(MapToDto(tour));
                }
            }

            return filteredTours;
        }

        public bool CheckIfAnyKeyPointInRange(List<TourKeyPoint> keyPoints, double? lat, double? lon, double radius)
        {
            return keyPoints.Any(keyPoint => IsInRange(keyPoint, lat, lon, radius));
        }

        public bool IsInRange(TourKeyPoint keyPoint, double? lat, double? lon, double radius)
        {
            double distance;
            int earthRadius = 6371000;
            double radiusInDegrees = radius * 360 / (2 * Math.PI * earthRadius);
            distance = Math.Sqrt(Math.Pow((double)(keyPoint.Latitude - lat), 2) +
                                 Math.Pow((double)(keyPoint.Longitude - lon), 2));
            return distance <= radiusInDegrees;
        }

        public Result<PagedResult<TourDto>> GetPagedByIds(List<int> ids, int page, int pageSize)
        {
            var result = _tourRepository.GetPagedByIds(ids, page, pageSize);
            return MapToDto(result);
        }
        public List<TourDto> GetAllByAuthorId(int authorId)
        {
            var result = _tourRepository.GetPagedByAuthorId(authorId, 0, 0);
            return MapToDto(result.Results).Value;
        }
        public Result<PagedResult<TourDto>> GetAllPagedByAuthorId(int authorId)
        {
            return GetPagedByAuthorId(authorId, 0, 0);
        }
        public async Task<string> UpdateAsync(TourDto tour, HttpClient _httpClient)
        {
            using StringContent jsonContent = new(System.Text.Json.JsonSerializer.Serialize(tour), Encoding.UTF8, "application/json");
            using HttpResponseMessage response = await sharedClient.PutAsync("/tours", jsonContent);
            response.EnsureSuccessStatusCode();
            var jsonResponse = await response.Content.ReadAsStringAsync();
            return jsonResponse;
        }
        public async Task<Result<TourDto>> GetAsync(int id)
        {
            using HttpResponseMessage response = await sharedClient.GetAsync("/tours/" + id.ToString());
            response.EnsureSuccessStatusCode();
            var jsonResponseGet = await response.Content.ReadAsStringAsync();

            var tourDto = JsonConvert.DeserializeObject<TourDto>(jsonResponseGet);
            return tourDto;
        }

        public async Task<Result<PagedResult<TourDto>>> GetAllAsync()
        {
            using HttpResponseMessage response = await sharedClient.GetAsync("/getAll");
            response.EnsureSuccessStatusCode();

            var jsonResponse = await response.Content.ReadAsStringAsync();

            var settings = new JsonSerializerSettings
            {
                Converters = new List<JsonConverter> { new CustomDateTimeConverter() }
            };
            // Parse the JSON response
            var parsedResponse = JObject.Parse(jsonResponse);

            // Extract the "results" array
            var results = parsedResponse["results"].ToString();
            var totalCount = parsedResponse["totalCount"].ToString();

            // Deserialize the results array to List<TourDto>
            var tourDtos = JsonConvert.DeserializeObject<List<TourDto>>(results, settings);

            var listResult = new PagedResult<TourDto>(tourDtos,int.Parse(totalCount));

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