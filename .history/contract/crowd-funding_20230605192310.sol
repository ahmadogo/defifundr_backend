// SPDX-License-Identifier: SEE LICENSE IN LICENSE

pragma solidity ^0.8.0;

contract CrowdFunding {
    struct Campaign {
        address owner;
        string campaignType;
        string title;
        string description;
        uint256 goal;
        uint256 deadline;
        uint256 totalFunds;
        uint256 totalContributors;
        string image;
        address[] donators;
        uint256[] donations;
    }

    mapping(uint256 => Campaign) public campaigns;

    uint public campaignCount = 0;

    function createCampaign() public{}

    function getCampaign(uint256 _campaignId) public view returns (Campaign memory) {}

    function donate(uint256 _campaignId) public payable {}

    function getCampaignDonators(uint256 _campaignId) public view returns (address[] memory) {}

    function getCampaignDonations(uint256 _campaignId) public view returns (uint256[] memory) {}

    function getCampaignDonationsSum(uint256 _campaignId) public view returns (uint256) {}

    
}
