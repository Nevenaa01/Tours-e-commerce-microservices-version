using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using AutoMapper;
using Explorer.BuildingBlocks.Core.UseCases;
using Explorer.Tours.API.Dtos;
using Explorer.Tours.API.Public;
using Explorer.Tours.Core.Domain;
using Explorer.Tours.Core.Domain.RepositoryInterfaces;
using FluentResults;
using Explorer.Stakeholders.API.Internal;
using Explorer.Tours.API.Public.Authoring;
using System.Text.Json;
using Explorer.Blog.Core.Domain;

namespace Explorer.Tours.Core.UseCases
{
    public class TourProblemService : CrudService<TourProblemDto, TourProblem>, ITourProblemService
    {
        private readonly ITourProblemRepository _tourProblemRepository;
        private readonly IUserNames _userNamesService;
        private readonly ITourService _tourService;
        private readonly HttpClient _httpClient;

        public TourProblemService(ICrudRepository<TourProblem> repository, 
            IMapper mapper, ITourProblemRepository tourProblemRepository, 
            IUserNames userNamesService, ITourService tourService,
            HttpClient httpClient) : base(repository, mapper)
        {
            _tourProblemRepository = tourProblemRepository;
            _userNamesService = userNamesService;
            _tourService = tourService;
            _httpClient = httpClient;
        }

        private static HttpClient sharedClient = new()
        {
            BaseAddress = new Uri("http://tours_service:8080")
        };
        public Result<List<TourProblemDto>> GetByTouristId(long touristId)
        {
            List<TourProblemDto> result = new List<TourProblemDto>();
            List<TourProblem> tourProblems = _tourProblemRepository.GetByTouristId(touristId);

            tourProblems.ForEach(t => result.Add(MapToDto(t)));

            return result;
        }
        public Result<List<TourProblemDto>> GetByAuthorId(long authorId)
        {
            var tours = _tourService.GetPaged(0, 0).Value.Results;
            List<TourProblemDto> result = new List<TourProblemDto>();
            List<TourProblem> tourProblems = new List<TourProblem>();
            foreach (var t in tours)
            {
                if (authorId == t.AuthorId)
                {
                    tourProblems.AddRange(_tourProblemRepository.GetByTourId(t.Id));
                }
            }
            tourProblems.ForEach(t => result.Add(MapToDto(t)));
            return result;
        }

        public async Task<List<TourProblemDto>> GetByAuthorIdAsync(long authorId)
        {
            using HttpResponseMessage response = await sharedClient.GetAsync("/getByAuthorId/" +  authorId);
            response.EnsureSuccessStatusCode();

            string responseData = await response.Content.ReadAsStringAsync();
            var jsonObject = JsonDocument.Parse(responseData).RootElement;

            List<TourProblemDto> tourProblems = new List<TourProblemDto>();

            foreach (var problemJson in jsonObject.EnumerateArray())
            {
                TourProblemDto tourProblem = new TourProblemDto
                {
                    TouristId = problemJson.GetProperty("touristId").GetInt32(),
                    TourId = problemJson.GetProperty("tourId").GetInt32(),
                    Category = (API.Dtos.TourProblemCategory)problemJson.GetProperty("category").GetInt32(),
                    Priority = (API.Dtos.TourProblemPriority)problemJson.GetProperty("priority").GetInt32(),
                    Description = problemJson.GetProperty("description").GetString(),
                    Time = DateTime.Parse(problemJson.GetProperty("time").GetString()),
                    IsSolved = problemJson.GetProperty("isSolved").GetBoolean(),
                    Messages = ParseMessages(problemJson.GetProperty("messages")),
                    Deadline = (DateTime)(problemJson.TryGetProperty("deadline", out var deadline) ?
                        DateTime.Parse(deadline.GetString()) : (DateTime?)null) ,
                    TouristUsername = problemJson.GetProperty("touristUsername").GetString(),
                    AuthorUsername = problemJson.GetProperty("authorUsername").GetString()
                };

                tourProblems.Add(tourProblem);
            }

            return tourProblems;
        }

        private List<TourProblemMessageDto> ParseMessages(JsonElement messagesJson)
        {
            List<TourProblemMessageDto> messages = new List<TourProblemMessageDto>();

            foreach (var messageJson in messagesJson.EnumerateArray())
            {
                TourProblemMessageDto message = new TourProblemMessageDto
                {
                    SenderId = messageJson.GetProperty("SenderId").GetInt64(),
                    RecipientId = messageJson.GetProperty("RecipientId").GetInt64(),
                    CreationTime = DateTime.Parse(messageJson.GetProperty("CreationTime").GetString()),
                    Description = messageJson.GetProperty("Description").GetString(),
                    SenderName = messageJson.GetProperty("SenderName").GetString(),
                    IsRead = messageJson.GetProperty("IsRead").GetBoolean()
                };

                messages.Add(message);
            }

            return messages;
        }

        public void FindNames(List<TourProblemDto> result)
        {
            var tours = _tourService.GetPaged(0, 0).Value.Results;

            foreach (var r in result)
            {
                long authorId = tours.Find(t => t.Id == r.TourId).AuthorId;
                r.AuthorUsername = _userNamesService.GetName(authorId).Username;
                r.TouristUsername = _userNamesService.GetName(r.TouristId).Username;
                foreach (var m in r.Messages)
                {
                    if(m.SenderId != 0)
                        m.SenderName = _userNamesService.GetName(m.SenderId).Username;
                }
            }
        }

        public Result<List<TourProblemMessageDto>> GetUnreadMessages(long id)
        {
            var tourProblems = GetPaged(0, 0).Value.Results;
            FindNames(tourProblems);
            List<TourProblemMessageDto> unreadMessages = new List<TourProblemMessageDto>();

            foreach (var t in tourProblems)
            {
                foreach (var m in t.Messages)
                {
                    if (m.RecipientId == id && !m.IsRead)
                    {
                        unreadMessages.Add(m);
                    }
                }
            }
             return unreadMessages;
           
        }
        public Result<TourProblemDto> GiveDeadline(DateTime deadline, long tourProblemId)
        {
            var tourProblem = _tourProblemRepository.GiveDeadline(deadline, tourProblemId);
            TourProblemDto dto = new TourProblemDto();
            dto.Deadline=deadline;
            return dto;
        }

        public Result<TourProblemDto> PunishAuthor(string authorUsername, long tourId, long tourProblemId)
        {
            var tourProblem=_tourProblemRepository.PunishAuthor(authorUsername, tourId, tourProblemId);
            TourProblemDto tpdto= new TourProblemDto();
            tpdto.IsSolved = true;
            return tpdto;
        }
    }
}
