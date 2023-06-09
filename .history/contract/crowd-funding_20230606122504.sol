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

    //? Create a new campaign
    function createCampaign(
        string memory _campaignType,
        string memory _title,
        string memory _description,
        uint256 _goal,
        uint256 _deadline,
        string memory _image
    ) public returns (uint256) {
        Campaign storage campaign = campaigns[campaignCount];

        require(campaign.deadline < block.timestamp, "Deadline is not passed");

        campaign.owner = msg.sender;
        campaign.campaignType = _campaignType;
        campaign.title = _title;
        campaign.description = _description;
        campaign.goal = _goal;
        campaign.deadline = _deadline;
        campaign.image = _image;

        campaignCount++;

        return campaignCount - 1;
    }

    //? Get campaign by id
    function getCampaign(
        uint256 _campaignId
    ) public view returns (Campaign memory) {
        return campaigns[_campaignId];
    }

    //? Donate to a campaign
    function donate(uint256 _campaignId) public payable {
        uint256 amount = msg.value;
        Campaign storage campaign = campaigns[_campaignId];

        require(campaign.deadline > block.timestamp, "Deadline is passed");
        require(msg.value > 0, "Donation amount must be greater than 0");

        campaign.totalFunds += msg.value;
        campaign.totalContributors += 1;
        campaign.donators.push(msg.sender);
        campaign.donations.push(msg.value);

        (bool sent, ) = payable(campaign.owner).call{value: amount}("");

        if (sent) {
            campaign.totalFunds += amount;
        } else {
            revert("Failed to send Ether");
        }
    }

    //? Get campaign donations
    function getCampaignDonators(
        uint256 _campaignId
    ) public view returns (address[] memory) {
        return campaigns[_campaignId].donators;
    }

    //? Get campaign donations
    function getCampaignDonations(
        uint256 _campaignId
    ) public view returns (uint256[] memory) {
        return campaigns[_campaignId].donations;
    }

    //? Get campaign donations sum
    function getCampaignDonationsSum(
        uint256 _campaignId
    ) public view returns (uint256) {
        uint256 sum = 0;

        for (uint256 i = 0; i < campaigns[_campaignId].donations.length; i++) {
            sum += campaigns[_campaignId].donations[i];
        }

        return sum;
    }

    //? Get campaign donations count
    function getCampaignDonationsCount(
        uint256 _campaignId
    ) public view returns (uint256) {
        return campaigns[_campaignId].donations.length;
    }

    //? Get campaign donators count
    function getCampaignDonatorsCount(
        uint256 _campaignId
    ) public view returns (uint256) {
        return campaigns[_campaignId].donators.length;
    }

    //? Get campaign donations count
    function getCampaignsCount() public view returns (uint256) {
        return campaignCount;
    }

    //? Get campaign donations count
    function getCampaigns() public view returns (Campaign[] memory) {
        Campaign[] memory _campaigns = new Campaign[](campaignCount);

        for (uint256 i = 0; i < campaignCount; i++) {
            _campaigns[i] = campaigns[i];
        }

        return _campaigns;
    }

    //? Get campaign donations count
    function getCampaignsByOwner(
        address _owner
    ) public view returns (Campaign[] memory) {
        Campaign[] memory _campaigns = new Campaign[](campaignCount);

        uint256 count = 0;

        for (uint256 i = 0; i < campaignCount; i++) {
            if (campaigns[i].owner == _owner) {
                _campaigns[count] = campaigns[i];
                count++;
            }
        }

        return _campaigns;
    }

    //? Get campaign By Type
    function getCampaignsByType(
        string memory _campaignType
    ) public view returns (Campaign[] memory) {
        Campaign[] memory _campaigns = new Campaign[](campaignCount);

        uint256 count = 0;

        for (uint256 i = 0; i < campaignCount; i++) {
            if (
                keccak256(abi.encodePacked(campaigns[i].campaignType)) ==
                keccak256(abi.encodePacked(_campaignType))
            ) {
                _campaigns[count] = campaigns[i];
                count++;
            }
        }

        return _campaigns;
    }

    //? Get Donors Addresses And Amounts
    function getDonorsAddressesAndAmounts(
        uint256 _campaignId
    ) public view returns (address[] memory, uint256[] memory, uint256) {
        return (
            campaigns[_campaignId].donators,
            campaigns[_campaignId].donations,
            campaigns[_campaignId].totalFunds
        );
    }

    //? Pay out to campaign owner if goal is reached
    function payOut(
        uint256 _campaignId
    ) public isCampaignOwner(_campaignId) isTimePassed(_campaignId) {
        Campaign storage campaign = campaigns[_campaignId];

        require(
            campaign.totalFunds >= campaign.goal,
            "Campaign goal is not reached"
        );

        (bool sent, ) = payable(campaign.owner).call{
            value: campaign.totalFunds
        }("");

        if (sent) {
            campaign.totalFunds = 0;
        } else {
            revert("Failed to send Ether");
        }
    }

    //? S
}
