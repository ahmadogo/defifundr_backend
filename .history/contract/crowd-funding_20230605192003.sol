// SPDX-License-Identifier: SEE LICENSE IN LICENSE

pragma solidity ^0.8.0;

contract CrowdFunding {
    struct Campaign {
       address  owner;
    string public campaignType;
    string title;
    string description;
    uint256 goal;
    uint256 deadline;
    uint256 totalFunds;
    uint256 totalContributors;
    string image;
    
    }



}