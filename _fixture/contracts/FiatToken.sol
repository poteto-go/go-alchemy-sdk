// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract FiatToken {
    string public name;
    string public symbol;
    string public currency;
    uint8 public decimals;
    uint256 public totalSupply;

    mapping(address => uint256) public balanceOf;
    mapping(address => mapping(address => uint256)) public allowance;
    mapping(address => bool) public isMinter;
    mapping(address => uint256) public minterAllowance;
    mapping(address => bool) public isBlacklisted;

    address public owner;
    address public masterMinter;
    address public pauser;
    address public blacklister;

    bool public paused;

    event Transfer(address indexed from, address indexed to, uint256 value);
    event Approval(address indexed owner, address indexed spender, uint256 value);
    event Mint(address indexed minter, address indexed to, uint256 amount);
    event Burn(address indexed burner, uint256 amount);
    event MinterConfigured(address indexed minter, uint256 minterAllowedAmount);
    event MinterRemoved(address indexed oldMinter);
    event BlackListed(address indexed account);
    event UnBlacklisted(address indexed account);
    event Pause();
    event Unpause();
    event OwnershipTransferred(address indexed previousOwner, address indexed newOwner);

    modifier onlyOwner() {
        require(msg.sender == owner, "FiatToken: caller is not the owner");
        _;
    }

    modifier onlyMasterMinter() {
        require(msg.sender == masterMinter, "FiatToken: caller is not the masterMinter");
        _;
    }

    modifier onlyMinter() {
        require(isMinter[msg.sender], "FiatToken: caller is not a minter");
        _;
    }

    modifier onlyPauser() {
        require(msg.sender == pauser, "FiatToken: caller is not the pauser");
        _;
    }

    modifier onlyBlacklister() {
        require(msg.sender == blacklister, "FiatToken: caller is not the blacklister");
        _;
    }

    modifier notBlacklisted(address account) {
        require(!isBlacklisted[account], "FiatToken: account is blacklisted");
        _;
    }

    modifier whenNotPaused() {
        require(!paused, "FiatToken: paused");
        _;
    }

    constructor(
        string memory _name,
        string memory _symbol,
        string memory _currency,
        uint8 _decimals,
        address _masterMinter,
        address _pauser,
        address _blacklister,
        address _owner
    ) {
        name = _name;
        symbol = _symbol;
        currency = _currency;
        decimals = _decimals;
        masterMinter = _masterMinter;
        pauser = _pauser;
        blacklister = _blacklister;
        owner = _owner;
    }

    function configureMinter(address minter, uint256 minterAllowedAmount)
        external
        onlyMasterMinter
        whenNotPaused
        returns (bool)
    {
        isMinter[minter] = true;
        minterAllowance[minter] = minterAllowedAmount;
        emit MinterConfigured(minter, minterAllowedAmount);
        return true;
    }

    function removeMinter(address minter) external onlyMasterMinter returns (bool) {
        isMinter[minter] = false;
        minterAllowance[minter] = 0;
        emit MinterRemoved(minter);
        return true;
    }

    function mint(address to, uint256 amount)
        external
        onlyMinter
        whenNotPaused
        notBlacklisted(msg.sender)
        notBlacklisted(to)
        returns (bool)
    {
        require(amount > 0, "FiatToken: mint amount must be greater than 0");
        require(minterAllowance[msg.sender] >= amount, "FiatToken: mint amount exceeds minterAllowance");

        minterAllowance[msg.sender] -= amount;
        totalSupply += amount;
        balanceOf[to] += amount;
        emit Mint(msg.sender, to, amount);
        emit Transfer(address(0), to, amount);
        return true;
    }

    function burn(uint256 amount)
        external
        onlyMinter
        whenNotPaused
        notBlacklisted(msg.sender)
    {
        require(amount > 0, "FiatToken: burn amount must be greater than 0");
        require(balanceOf[msg.sender] >= amount, "FiatToken: burn amount exceeds balance");

        totalSupply -= amount;
        balanceOf[msg.sender] -= amount;
        emit Burn(msg.sender, amount);
        emit Transfer(msg.sender, address(0), amount);
    }

    function transfer(address to, uint256 amount)
        external
        whenNotPaused
        notBlacklisted(msg.sender)
        notBlacklisted(to)
        returns (bool)
    {
        require(to != address(0), "FiatToken: transfer to the zero address");
        require(balanceOf[msg.sender] >= amount, "FiatToken: transfer amount exceeds balance");

        balanceOf[msg.sender] -= amount;
        balanceOf[to] += amount;
        emit Transfer(msg.sender, to, amount);
        return true;
    }

    function approve(address spender, uint256 amount)
        external
        whenNotPaused
        notBlacklisted(msg.sender)
        notBlacklisted(spender)
        returns (bool)
    {
        allowance[msg.sender][spender] = amount;
        emit Approval(msg.sender, spender, amount);
        return true;
    }

    function transferFrom(address from, address to, uint256 amount)
        external
        whenNotPaused
        notBlacklisted(msg.sender)
        notBlacklisted(from)
        notBlacklisted(to)
        returns (bool)
    {
        require(to != address(0), "FiatToken: transfer to the zero address");
        require(balanceOf[from] >= amount, "FiatToken: transfer amount exceeds balance");
        require(allowance[from][msg.sender] >= amount, "FiatToken: transfer amount exceeds allowance");

        allowance[from][msg.sender] -= amount;
        balanceOf[from] -= amount;
        balanceOf[to] += amount;
        emit Transfer(from, to, amount);
        return true;
    }

    function blacklist(address account) external onlyBlacklister {
        isBlacklisted[account] = true;
        emit BlackListed(account);
    }

    function unBlacklist(address account) external onlyBlacklister {
        isBlacklisted[account] = false;
        emit UnBlacklisted(account);
    }

    function pause() external onlyPauser {
        paused = true;
        emit Pause();
    }

    function unpause() external onlyPauser {
        paused = false;
        emit Unpause();
    }

    function transferOwnership(address newOwner) external onlyOwner {
        require(newOwner != address(0), "FiatToken: new owner is the zero address");
        emit OwnershipTransferred(owner, newOwner);
        owner = newOwner;
    }

    function updateMasterMinter(address newMasterMinter) external onlyOwner {
        require(newMasterMinter != address(0), "FiatToken: new masterMinter is the zero address");
        masterMinter = newMasterMinter;
    }

    function updatePauser(address newPauser) external onlyOwner {
        require(newPauser != address(0), "FiatToken: new pauser is the zero address");
        pauser = newPauser;
    }

    function updateBlacklister(address newBlacklister) external onlyOwner {
        require(newBlacklister != address(0), "FiatToken: new blacklister is the zero address");
        blacklister = newBlacklister;
    }

    function version() public pure returns (string memory) {
        return "1";
    }
}
