// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

interface IERC1155Receiver {
    function onERC1155Received(address operator, address from, uint256 id, uint256 value, bytes calldata data)
        external
        returns (bytes4);

    function onERC1155BatchReceived(
        address operator,
        address from,
        uint256[] calldata ids,
        uint256[] calldata values,
        bytes calldata data
    ) external returns (bytes4);
}

// Minimal ERC1155 fixture used by the e2e suite. It implements enough of
// the standard to exercise the SDK's ERC1155 read and write methods.
contract ERC1155 {
    string private _uri = "https://example.com/erc1155/{id}.json";

    mapping(uint256 => mapping(address => uint256)) private _balances;
    mapping(address => mapping(address => bool)) private _operatorApprovals;

    event TransferSingle(
        address indexed operator, address indexed from, address indexed to, uint256 id, uint256 value
    );
    event TransferBatch(
        address indexed operator,
        address indexed from,
        address indexed to,
        uint256[] ids,
        uint256[] values
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

    function safeTransferFrom(address from, address to, uint256 id, uint256 amount, bytes memory data) public {
        require(
            from == msg.sender || isApprovedForAll(from, msg.sender),
            "ERC1155: caller is not token owner or approved"
        );
        require(to != address(0), "ERC1155: transfer to the zero address");
        require(_balances[id][from] >= amount, "ERC1155: insufficient balance for transfer");

        _balances[id][from] -= amount;
        _balances[id][to] += amount;
        emit TransferSingle(msg.sender, from, to, id, amount);

        if (to.code.length > 0) {
            try IERC1155Receiver(to).onERC1155Received(msg.sender, from, id, amount, data) returns (bytes4 retval) {
                require(
                    retval == IERC1155Receiver.onERC1155Received.selector,
                    "ERC1155: transfer to non-ERC1155Receiver"
                );
            } catch {
                revert("ERC1155: transfer to non-ERC1155Receiver");
            }
        }
    }

    function safeBatchTransferFrom(
        address from,
        address to,
        uint256[] memory ids,
        uint256[] memory amounts,
        bytes memory data
    ) public {
        require(
            from == msg.sender || isApprovedForAll(from, msg.sender),
            "ERC1155: caller is not token owner or approved"
        );
        require(to != address(0), "ERC1155: transfer to the zero address");
        require(ids.length == amounts.length, "ERC1155: ids and amounts length mismatch");

        for (uint256 i = 0; i < ids.length; ++i) {
            require(_balances[ids[i]][from] >= amounts[i], "ERC1155: insufficient balance for transfer");
            _balances[ids[i]][from] -= amounts[i];
            _balances[ids[i]][to] += amounts[i];
        }
        emit TransferBatch(msg.sender, from, to, ids, amounts);

        if (to.code.length > 0) {
            try IERC1155Receiver(to).onERC1155BatchReceived(msg.sender, from, ids, amounts, data) returns (
                bytes4 retval
            ) {
                require(
                    retval == IERC1155Receiver.onERC1155BatchReceived.selector,
                    "ERC1155: transfer to non-ERC1155Receiver"
                );
            } catch {
                revert("ERC1155: transfer to non-ERC1155Receiver");
            }
        }
    }

    function mint(address to, uint256 id, uint256 amount) public {
        require(to != address(0), "ERC1155: mint to the zero address");
        _balances[id][to] += amount;
        emit TransferSingle(msg.sender, address(0), to, id, amount);
    }
}
