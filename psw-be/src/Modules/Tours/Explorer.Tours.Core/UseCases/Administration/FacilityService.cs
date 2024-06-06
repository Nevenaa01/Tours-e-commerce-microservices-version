using AutoMapper;
using Explorer.BuildingBlocks.Core.UseCases;
using Explorer.Tours.API.Dtos;
using Explorer.Tours.API.Public.Administration;
using Explorer.Tours.Core.Domain;
using System.Text.Json;
using System.Text;

namespace Explorer.Tours.Core.UseCases.Administration
{
    public class FacilityService : CrudService<FacilityDto, Facility>, IFacilityService
    {
        public FacilityService(ICrudRepository<Facility> crudRepository, IMapper mapper) : base(crudRepository, mapper) { }

        public async Task<string> CreateAsync(FacilityDto facilityDto, HttpClient _httpClient)
        {
            using StringContent jsonContent = new(JsonSerializer.Serialize(facilityDto), Encoding.UTF8, "application/json");
            using HttpResponseMessage response = await _httpClient.PostAsync("http://tours_service:8080/facilities", jsonContent);
            response.EnsureSuccessStatusCode();
            var jsonResponse = await response.Content.ReadAsStringAsync();
            return jsonResponse;
        }

        public async Task<string> DeleteAsync(int id, HttpClient _httpClient)
        {
            using HttpResponseMessage response = await _httpClient.DeleteAsync("http://tours_service:8080/facilities/" + id);
            //response.EnsureSuccessStatusCode();
            return "works";
        }
    }
}
