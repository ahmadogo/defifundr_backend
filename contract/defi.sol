// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity ^0.8.0;

import "./IERC20.sol";

contract CrowdFunding {
    struct Campaign {
        address owner;
        string campaignType;
        string title;
        string description;
        uint256 goal;
        uint256 deadline;
        string image;
        uint256 id;
        mapping(address => uint256) totalFundsPerToken; 
        address[] supportedTokens;
        bool isDeleted;
    }

    mapping(uint256 => Campaign) private campaigns;
    uint256 public campaignCount;

    modifier onlyOwner(uint256 _campaignId) {
        require(
            msg.sender == campaigns[_campaignId].owner,
            "Not campaign owner"
        );
        _;
    }

    modifier campaignActive(uint256 _campaignId) {
        require(!campaigns[_campaignId].isDeleted, "Campaign deleted");
        require(
            block.timestamp <= campaigns[_campaignId].deadline,
            "Campaign expired"
        );
        _;
    }

    function createCampaign(
        string calldata _campaignType,
        string calldata _title,
        string calldata _description,
        uint256 _goal,
        uint256 _deadline,
        string calldata _image,
        address[] calldata _supportedTokens
    ) external returns (uint256) {
        require(_goal > 0, "Goal must be greater than 0");
        require(_deadline > block.timestamp, "Invalid deadline");
        require(
            _supportedTokens.length > 0,
            "At least one token must be supported"
        );

        Campaign storage campaign = campaigns[campaignCount];
        campaign.owner = msg.sender;
        campaign.campaignType = _campaignType;
        campaign.title = _title;
        campaign.description = _description;
        campaign.goal = _goal;
        campaign.deadline = _deadline;
        campaign.image = _image;
        campaign.id = campaignCount;
        campaign.supportedTokens = _supportedTokens;

        return campaignCount++;
    }

    function donate(
        uint256 _campaignId,
        address _token,
        uint256 _amount
    ) external campaignActive(_campaignId) {
        Campaign storage campaign = campaigns[_campaignId];

        require(isTokenSupported(_campaignId, _token), "Token not supported");
        require(_amount > 0, "Amount must be greater than 0");

        IERC20 token = IERC20(_token);
        require(
            token.transferFrom(msg.sender, address(this), _amount),
            "Token transfer failed"
        );

        campaign.totalFundsPerToken[_token] += _amount;
    }

    function isTokenSupported(
        uint256 _campaignId,
        address _token
    ) public view returns (bool) {
        Campaign storage campaign = campaigns[_campaignId];
        for (uint256 i = 0; i < campaign.supportedTokens.length; i++) {
            if (campaign.supportedTokens[i] == _token) {
                return true;
            }
        }
        return false;
    }

    function getFundsPerToken(
        uint256 _campaignId,
        address _token
    ) public view returns (uint256) {
        return campaigns[_campaignId].totalFundsPerToken[_token];
    }

    function withdrawFunds(
        uint256 _campaignId,
        address _token
    ) external onlyOwner(_campaignId) {
        Campaign storage campaign = campaigns[_campaignId];
        uint256 amount = campaign.totalFundsPerToken[_token];

        require(amount > 0, "No funds to withdraw");

        IERC20 token = IERC20(_token);
        require(token.transfer(campaign.owner, amount), "Withdraw failed");

        campaign.totalFundsPerToken[_token] = 0;
    }
}
