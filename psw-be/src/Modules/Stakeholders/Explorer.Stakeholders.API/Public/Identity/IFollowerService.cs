using Explorer.Stakeholders.API.Dtos;
using FluentResults;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace Explorer.Stakeholders.API.Public.Identity
{
    public interface IFollowerService
    {
        public Result<List<SavedNotificationDto>> GetFollowersNotifications(int id);
        public Result<FollowerDto> Create(FollowerDto follower);
        public Task CreateAsync(FollowerDto follower);
        public Result Delete(int followerId, int followedId);
        public Result<List<FollowerDto>> GetFollowings(int id);

        public Task DeleteAsync(int followerId, int followedId);
    }
}
