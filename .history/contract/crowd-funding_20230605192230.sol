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

    function createCampaign(
        string memory _campaignType,
        string memory _title,
        string memory _description,
        uint256 _goal,
        uint256 _deadline,
        string memory _image
    ) public {
        campaignCount++;
        campaigns[campaignCount] = Campaign(
            msg.sender,
            _campaignType,
            _title,
            _description,
            _goal,
            _deadline,
            0,
            0,
            _image,
            new address[](0),
            new uint256[](0)
        );
    }
}
