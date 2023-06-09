// SPDX-License-Identifier: SEE LICENSE IN LICENSE

pragma solidity ^0.8.0;

contract CrowdFunding {
    mapping(address => uint) public contributors;
    address public admin;
    uint public noOfContributors;
    uint public minimumContribution;
    uint public deadline;
    uint public goal;
    uint public raisedAmount = 0;

    struct Request {
        string description;
        address payable recipient;
        uint value;
        bool completed;
        uint noOfVoters;
        mapping(address => bool) voters;
    }

    Request[] public requests;

    event ContributeEvent(address sender, uint value);
    event CreateRequestEvent(
        string _description,
        address _recipient,
        uint _value
    );
    event MakePaymentEvent(address recipient, uint value);

    constructor(uint _goal, uint _deadline) {
        goal = _goal;
        deadline = block.timestamp + _deadline;
        admin = msg.sender;
        minimumContribution = 100 wei;
    }

    modifier onlyAdmin() {
        require(msg.sender == admin, "Only admin can call this function");
        _;
    }

    function contribute() public payable {
        require(block.timestamp < deadline, "The deadline has passed");
        require(
            msg.value >= minimumContribution,
            "The minimum contribution not met"
        );

        if (contributors[msg.sender] == 0) {
            noOfContributors++;
        }

        contributors[msg.sender] += msg.value;
        raisedAmount += msg.value;

        emit ContributeEvent(msg.sender, msg.value);
    }

    function getBalance() public view returns (uint) {
        return address(this).balance;
    }

    function getRefund() public {
        require(block.timestamp > deadline && raisedAmount < goal);
        require(contributors[msg.sender] > 0);

        address payable recipient = payable(msg.sender);
        uint value = contributors[msg.sender];

        recipient.transfer(value);
        contributors[msg.sender] = 0;
    }

    function createRequest(
        string memory _description,
        address payable _recipient,
        uint _value
    ) public onlyAdmin {
        Request memory newRequest = Request({
            description: _description,
            recipient: _recipient,
            value: _value,
            completed: false,
            noOfVoters: 0
        });

        requests.push(newRequest);

        emit CreateRequestEvent(_description, _recipient, _value);
    }

    function voteRequest(uint index) public {
        Request storage thisRequest = requests[index];

        require(
            contributors[msg.sender] > 0,
            "You must be a contributor to vote"
        );
        require(
            thisRequest.voters[msg.sender] == false,
            "You have already voted"
        );

        thisRequest.voters[msg.sender] = true;
        thisRequest.noOfVoters++;
    }
}
