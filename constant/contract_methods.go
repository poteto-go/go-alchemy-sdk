package constant

var (
	// ERC-20 function signatures
	TransferFnSignature     = []byte("transfer(address,uint256)")
	TransferFromFnSignature = []byte("transferFrom(address,address,uint256)")
	ApproveFnSignature      = []byte("approve(address,uint256)")
	BalanceOfFnSignature    = []byte("balanceOf(address)")
	TotalSupplyFnSignature  = []byte("totalSupply()")
	AllowanceFnSignature    = []byte("allowance(address,address)")
	NameFnSignature         = []byte("name()")
	SymbolFnSignature       = []byte("symbol()")
	DecimalsFnSignature     = []byte("decimals()")

	// Stablecoin function signatures
	MintFnSignature               = []byte("mint(address,uint256)")
	BurnFnSignature               = []byte("burn(uint256)")
	BlacklistFnSignature          = []byte("blacklist(address)")
	UnBlacklistFnSignature        = []byte("unBlacklist(address)")
	IsBlacklistedFnSignature      = []byte("isBlacklisted(address)")
	PauseFnSignature              = []byte("pause()")
	UnpauseFnSignature            = []byte("unpause()")
	PausedFnSignature             = []byte("paused()")
	TransferOwnershipFnSignature  = []byte("transferOwnership(address)")
	OwnerFnSignature              = []byte("owner()")
	CurrencyFnSignature           = []byte("currency()")
	VersionFnSignature            = []byte("version()")
	MasterMinterFnSignature       = []byte("masterMinter()")
	PauserFnSignature             = []byte("pauser()")
	BlacklisterFnSignature        = []byte("blacklister()")
	ConfigureMinterFnSignature    = []byte("configureMinter(address,uint256)")
	RemoveMinterFnSignature       = []byte("removeMinter(address)")
	MinterAllowanceFnSignature    = []byte("minterAllowance(address)")
	IsMinterFnSignature           = []byte("isMinter(address)")
	UpdateMasterMinterFnSignature = []byte("updateMasterMinter(address)")
	UpdateBlacklisterFnSignature  = []byte("updateBlacklister(address)")
	UpdatePauserFnSignature       = []byte("updatePauser(address)")

	// ERC-721 function signatures
	OwnerOfFnSignature                  = []byte("ownerOf(uint256)")
	TokenURIFnSignature                 = []byte("tokenURI(uint256)")
	GetApprovedFnSignature              = []byte("getApproved(uint256)")
	IsApprovedForAllFnSignature         = []byte("isApprovedForAll(address,address)")
	SafeTransferFromFnSignature         = []byte("safeTransferFrom(address,address,uint256)")
	SafeTransferFromWithDataFnSignature = []byte("safeTransferFrom(address,address,uint256,bytes)")
	SetApprovalForAllFnSignature        = []byte("setApprovalForAll(address,bool)")
	// NOTE: ERC-721 approve(to, tokenId) has the same selector as ERC-20's
	// approve(address,uint256), so it reuses ApproveFnSignature above.

	// ERC-1155 function signatures
	// NOTE: balanceOfToken(address,uint256) differs from ERC-20/721's
	// balanceOf(address), so it has its own selector.
	BalanceOfTokenFnSignature = []byte("balanceOf(address,uint256)")
	BalanceOfBatchFnSignature = []byte("balanceOfBatch(address[],uint256[])")
	UriFnSignature            = []byte("uri(uint256)")
	// ERC-1155 safeTransferFrom takes (from,to,id,amount,data) — a distinct selector
	// from ERC-721's safeTransferFrom(from,to,id) and safeTransferFrom(from,to,id,bytes).
	Erc1155SafeTransferFromFnSignature = []byte("safeTransferFrom(address,address,uint256,uint256,bytes)")
	SafeBatchTransferFromFnSignature   = []byte("safeBatchTransferFrom(address,address,uint256[],uint256[],bytes)")
	// NOTE: ERC-1155 setApprovalForAll(operator,bool) / isApprovedForAll(account,
	// operator) share the ERC-721 selectors, so they reuse the signatures above.

	// EIP-2612
	PermitFnSignature          = []byte("permit(address,address,uint256,uint256,uint8,bytes32,bytes32)")
	NoncesFnSignature          = []byte("nonces(address)")
	DomainSeparatorFnSignature = []byte("DOMAIN_SEPARATOR()")

	// EIP-3009
	AuthorizationStateFnSignature        = []byte("authorizationState(address,bytes32)")
	TransferWithAuthorizationFnSignature = []byte("transferWithAuthorization(address,address,uint256,uint256,uint256,bytes32,uint8,bytes32,bytes32)")
	ReceiveWithAuthorizationFnSignature  = []byte("receiveWithAuthorization(address,address,uint256,uint256,uint256,bytes32,uint8,bytes32,bytes32)")
	CancelAuthorizationFnSignature       = []byte("cancelAuthorization(address,bytes32,uint8,bytes32,bytes32)")
)
