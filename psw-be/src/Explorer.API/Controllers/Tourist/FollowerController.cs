using Explorer.Stakeholders.API.Dtos;
using Explorer.Stakeholders.API.Public.Identity;
using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;

namespace Explorer.API.Controllers.Tourist.Identity
{
    [Authorize(Policy = "touristPolicy")]
    [Route("api/tourist/follower")]
    public class FollowerController : BaseApiController
    {
        private readonly IFollowerService _followerService;

        public FollowerController(IFollowerService followerService)
        {
            _followerService = followerService;
        }

        [HttpGet("{id:int}")]
        public ActionResult<List<FollowerDto>> GetFollowersNotifications(int id)
        {
            var result = _followerService.GetFollowersNotifications(id);
            return CreateResponse(result);
        }

        /*[HttpPost]
        public ActionResult<FollowerDto> Create([FromBody] FollowerDto follower)
        {
            var result = _followerService.Create(follower);
            return CreateResponse(result);
        }*/

        [HttpPost]
        public async Task<IActionResult> CreatAync([FromBody] FollowerDto follower)
        {
            try
            {
                await _followerService.CreateAsync(follower);
                return Ok();
            }
            catch(Exception e)
            {
                return StatusCode(500, e.Message);
            }
        }

        /*[HttpDelete("{followerId:int}/{followedId:int}")]
        public ActionResult Delete(int followerId, int followedId)
        {
            var result = _followerService.Delete(followerId, followedId);
            return CreateResponse(result);
        }*/

        [HttpDelete("{followerId:int}/{followedId:int}")]
        public async Task<IActionResult> DeleteAsync(int followerId, int followedId)
        {
            try
            {
                await _followerService.DeleteAsync(followerId, followedId);
                return Ok();
            }
            catch (Exception e)
            {
                return StatusCode(500, e.Message);
            }
        }

        [HttpGet("followings/{id:int}")]
        public ActionResult<List<FollowerDto>> GetFollowings(int id)
        {
            var result = _followerService.GetFollowings(id);
            return CreateResponse(result);
        }
    }
}
