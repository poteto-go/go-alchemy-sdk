// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

// Minimal ERC1155 fixture used by the e2e suite. It implements just enough of
// the standard to exercise the SDK's Erc1155 read methods (balanceOf,
// balanceOfBatch, uri) plus a mint helper for test setup.
contract ERC1155 {
    // Shared metadata URI. ERC1155 uses a single URI template with an `{id}`
    // placeholder rather than a per-token URI.
    string private _uri = "https://example.com/erc1155/{id}.json";

    mapping(uint256 => mapping(address => uint256)) private _balances;
    mapping(address => mapping(address => bool)) private _operatorApprovals;

    event TransferSingle(
        address indexed operator, address indexed from, address indexed to, uint256 id, uint256 value
    );
    event ApprovalForAll(address indexed account, address indexed operator, bool approved);

    function uri(uint256) public view returns (string memory) {
        return _uri;
    }

    function balanceOf(address account, uint256 id) public view returns (uint256) {
        require(account != address(0), "ERC1155: address zero is not a valid owner");
        return _balances[id][account];
    }

    function balanceOfBatch(address[] memory accounts, uint256[] memory ids)
        public
        view
        returns (uint256[] memory)
    {
        require(accounts.length == ids.length, "ERC1155: accounts and ids length mismatch");
        uint256[] memory batchBalances = new uint256[](accounts.length);
        for (uint256 i = 0; i < accounts.length; ++i) {
            batchBalances[i] = balanceOf(accounts[i], ids[i]);
        }
        return batchBalances;
    }

    function isApprovedForAll(address account, address operator) public view returns (bool) {
        return _operatorApprovals[account][operator];
    }

    function setApprovalForAll(address operator, bool approved) public {
        require(operator != msg.sender, "ERC1155: setting approval status for self");
        _operatorApprovals[msg.sender][operator] = approved;
        emit ApprovalForAll(msg.sender, operator, approved);
    }

    function mint(address to, uint256 id, uint256 amount) public {
        require(to != address(0), "ERC1155: mint to the zero address");
        _balances[id][to] += amount;
        emit TransferSingle(msg.sender, address(0), to, id, amount);
    }
}
