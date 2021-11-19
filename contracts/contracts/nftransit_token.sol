// SPDX-License-Identifier: MIT

pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC1155/ERC1155.sol";
import "@openzeppelin/contracts/utils/cryptography/draft-EIP712.sol";
import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";

contract NFTransitToken is ERC1155, EIP712 {
    bytes32 private constant mintRequestTypeHash =
        keccak256(
            "MintRequestType(bytes32[] _uris,address[] _tos,uint256 nonce)"
        );

    bytes32 private constant burnRequestTypeHash =
        keccak256("BurnRequestType(uint256 nonce)");

    address public admin;
    uint256 public tokenIds;
    uint256 public mintNonce;
    uint256 public burnNonce;

    mapping(uint256 => string) private uris;

    event tokenBurn(uint256[] ids, uint256 burnNonce, address user);

    modifier onlyAdmin(
        string[] memory _uris,
        address[] memory _tos,
        bytes memory signature
    ) {
        if (msg.sender != admin) {
            require(
                admin ==
                    ECDSA.recover(
                        hashMintRequest(_uris, _tos, mintNonce),
                        signature
                    ),
                "tx signature mismatch"
            );
        }
        _;
    }

    constructor(
        address _admin,
        string memory _name,
        string memory _version
    ) ERC1155("") EIP712(_name, _version) {
        admin = _admin;
        tokenIds = 0;
        mintNonce = 0;
        burnNonce = 0;
    }

    function mint(
        string[] memory _uris,
        address[] memory _tos,
        bytes memory signature
    ) external onlyAdmin(_uris, _tos, signature) {
        require(
            _uris.length == _tos.length,
            "uri and recipient addr length mistmatch"
        );
        for (uint256 i = 0; i < _uris.length; i++) {
            _mint(_tos[i], tokenIds, 1, "");
            uris[tokenIds] = _uris[i];
            tokenIds++;
        }
        mintNonce++;
    }

    function burn(uint256[] memory ids) public {
        for (uint256 i = 0; i < ids.length; i++) {
            _burn(msg.sender, ids[i], 1);
        }
        emit tokenBurn(ids, burnNonce, msg.sender);
        burnNonce++;
    }

    function user_tokens(address account)
        public
        view
        returns (string[] memory)
    {
        string[] memory meta_uris = new string[](tokenIds);
        for (uint256 i = 0; i < tokenIds; i++) {
            if (balanceOf(account, i) != 0) {
                meta_uris[i] = uris[i];
            } else {
                meta_uris[i] = "";
            }
        }
        return meta_uris;
    }

    function uri(uint256 id) public view override returns (string memory) {
        return uris[id];
    }

    function hashMintRequest(
        string[] memory _uris,
        address[] memory _tos,
        uint256 nonce
    ) public view returns (bytes32) {
        return
            _hashTypedDataV4(
                keccak256(
                    abi.encode(
                        mintRequestTypeHash,
                        hashStrArray(_uris),
                        keccak256(abi.encodePacked(_tos)),
                        nonce
                    )
                )
            );
    }

    function hashStrArray(string[] memory array)
        internal
        pure
        returns (bytes32)
    {
        bytes32[] memory _array = new bytes32[](array.length);
        for (uint256 i = 0; i < array.length; ++i) {
            _array[i] = keccak256(abi.encodePacked(array[i]));
        }
        return keccak256(abi.encodePacked(_array));
    }

    function hashBurnRequest(uint256 nonce) public view returns (bytes32) {
        return
            _hashTypedDataV4(keccak256(abi.encode(burnRequestTypeHash, nonce)));
    }
}
