using Explorer.BuildingBlocks.Core.UseCases;
using Explorer.Tours.API.Dtos;
using Explorer.Tours.API.Public.Administration;
using Explorer.Tours.Core.Domain;
using Explorer.Tours.Core.UseCases.Administration;
using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;
using System.Text;
using System.Text.Json;

namespace Explorer.API.Controllers.Author
{
    
    [Authorize(Policy = "authorPolicy")]
    [Route("api/author/facilities")]
    public class FacilityController : BaseApiController
    {
        private readonly IFacilityService _facilityService;
        private readonly IPublicFacilityService _publicFacilityService;
        private readonly HttpClient _httpClient;
        public FacilityController(IFacilityService facilityService, IPublicFacilityService publicFacilityService)
        {
            _facilityService = facilityService;
            _publicFacilityService = publicFacilityService;
            _httpClient = new HttpClient();
        }

        [HttpGet]
        public ActionResult<PagedResult<FacilityDto>> GetAll([FromQuery] int page, [FromQuery] int pageSize)
        {
            var result = _facilityService.GetPaged(page, pageSize);
            return CreateResponse(result);
        }

        [HttpGet("{id:int}")]
        public ActionResult<FacilityDto> Get(int id)
        {
            var result = _facilityService.Get(id);
            return CreateResponse(result);
        }

        [HttpPost]
        public Task<string> Create([FromBody] FacilityDto facilityDto)
        {
            var result = _facilityService.CreateAsync(facilityDto, _httpClient);
            return result;
        }

        [HttpPut("{id:int}")]
        public ActionResult<FacilityDto> Update([FromBody] FacilityDto facility)
        {
            var result = _facilityService.Update(facility);
            return CreateResponse(result);
        }
        //public async Task<string> CreateAsync(FacilityDto facilityDto, HttpClient _httpClient)
        //{
        //    using StringContent jsonContent = new(JsonSerializer.Serialize(facilityDto), Encoding.UTF8, "application/json");
        //    using HttpResponseMessage response = await _httpClient.PostAsync("http://localhost:8080/facilities", jsonContent);
        //    response.EnsureSuccessStatusCode();
        //    var jsonResponse = await response.Content.ReadAsStringAsync();
        //    return jsonResponse;
        //}
        [HttpDelete("{id:int}")]
        public async Task<string> DeleteAsync(int id, HttpClient _httpClient)
        {
            var result = _facilityService.DeleteAsync(id, _httpClient);
            return "works";
        }
        
        //public ActionResult Delete(int id)
        //{
        //    var result = _facilityService.Delete(id);
        //    return CreateResponse(result);
        //}



        [HttpGet("public")]
        public ActionResult<PagedResult<PublicFacilityDto>> GetAllPublic([FromQuery] int page, [FromQuery] int pageSize)
        {
            var result = _publicFacilityService.GetPaged(page, pageSize);
            return CreateResponse(result);
        }

        [HttpPost("public")]
        public ActionResult<PublicFacilityDto> Create([FromBody] PublicFacilityDto facility)
        {
            var result = _publicFacilityService.Create(facility);
            return CreateResponse(result);
        }
        [HttpPut("public/{facilityId}/{status}")]
        public ActionResult<PublicFacilityDto> ChangeStatus(int facilityId, String status)
        {
            var result = _publicFacilityService.ChangeStatus(facilityId, status);
            return CreateResponse(result);
        }

        [HttpGet("public/{status}")]
        public ActionResult<PagedResult<PublicFacilityDto>> GetByStatus(String status)
        {
            var result = _publicFacilityService.GetByStatus(status);
            return CreateResponse(result);
        }
    }
}

