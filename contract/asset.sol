pragma solidity >=0.4.22 <0.6;

contract Asset {
	address public owner;
	address public pendingAuction;

	constructor() public {
		owner = tx.origin;
	}	

	function startAuction(address auction) public  {
		pendingAuction = auction;
	}

	function endAuction(address winer) public {
		// add signature verification of some kind (e.g, threshold sig)
		delete pendingAuction;
		owner = winer;
	}
}