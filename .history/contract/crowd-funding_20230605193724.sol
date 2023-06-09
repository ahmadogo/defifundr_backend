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

    modifier isTimePassed(uint256 _campaignId) {
        require(
            block.timestamp > campaigns[_campaignId].deadline,
            "Deadline is not passed"
        );
        _;
    }

    modifier isCampaignOwner(uint256 _campaignId) {
        require(
            msg.sender == campaigns[_campaignId].owner,
            "You are not the owner of this campaign"
        );
        _;
    }

    function createCampaign(
        string memory _campaignType,
        string memory _title,
        string memory _description,
        uint256 _goal,
        uint256 _deadline,
        string memory _image
    ) public 

    function getCampaign(
        uint256 _campaignId
    ) public view returns (Campaign memory) {}

    function donate(uint256 _campaignId) public payable {}

    function getCampaignDonators(
        uint256 _campaignId
    ) public view returns (address[] memory) {}

    function getCampaignDonations(
        uint256 _campaignId
    ) public view returns (uint256[] memory) {}

    function getCampaignDonationsSum(
        uint256 _campaignId
    ) public view returns (uint256) {}

    function getCampaignDonationsCount(
        uint256 _campaignId
    ) public view returns (uint256) {}

    function getCampaignDonatorsCount(
        uint256 _campaignId
    ) public view returns (uint256) {}

    function getCampaignsCount() public view returns (uint256) {}

    function getCampaigns() public view returns (Campaign[] memory) {}

    function getCampaignsByOwner(
        address _owner
    ) public view returns (Campaign[] memory) {}

    function getCampaignsByType(
        string memory _campaignType
    ) public view returns (Campaign[] memory) {}
}
