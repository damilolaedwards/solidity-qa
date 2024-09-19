"""
This script runs Slither on a directory and writes the output to a specified file.
"""

import json
import os
import argparse
from typing import List

from slither import Slither
from slither.core.declarations import Function, Contract


class SlitherHelper:
    @staticmethod
    def _parse_array(arrString: str):
        return arrString.split(',')

    @staticmethod
    def _parse_dict(dictString: str):
        return json.loads(dictString)

    @staticmethod
    def _filter_contracts(contracts: List[Contract], target: str, exclude_paths: List[str]):
        # Normalize target_dir path
        target = os.path.abspath(target)

        # Normalize and filter out non-existing exclude paths
        exclude_paths = [os.path.abspath(
            path) for path in exclude_paths if os.path.exists(path)]

        # Helper function to determine if a contract's path falls under target_dir
        def is_under_target_dir(contract_path):
            contract_path = os.path.abspath(contract_path)
            return contract_path.startswith(target)

        # Helper function to determine if a contract's path falls under any exclude path
        def is_under_exclude_paths(contract_path):
            contract_path = os.path.abspath(contract_path)
            return any(contract_path.startswith(exclude_path) for exclude_path in exclude_paths)

        # Filter documents based on the target_dir and exclude_paths
        filtered_documents = [
            c for c in contracts
            if is_under_target_dir(c.source_mapping.filename.relative) and not is_under_exclude_paths(c.source_mapping.filename.relative)
        ]

        return filtered_documents

    @staticmethod
    def _get_inheritance_tree(contract: Contract):
        inheritance_tree = {
            "id": contract.id,
            "name": contract.name,
            "code": SlitherHelper._get_contract_code(contract),
            "is_abstract": contract.is_abstract,
            "is_interface": contract.is_interface,
            "is_library": contract.is_library,
            "functions": SlitherHelper._get_functions_data(contract.functions),
            "inherited_contracts": []
        }
        for inherited_contract in contract.inheritance:
            inheritance_tree["inherited_contracts"].append(
                SlitherHelper._get_inheritance_tree(inherited_contract))

        return inheritance_tree

    @staticmethod
    def _get_functions_data(functions: List[Function]):
        functions_data = []

        for function in functions:
            functions_data.append({
                "id": function.id,
                "name": function.name,
                "visibility": function.visibility,
                "view": function.view,
                "pure": function.pure,
                "returns": [str(r.type) for r in function.returns] if function.returns != None else [],
                "parameters": [
                    {
                        "name": p.name,
                        "is_constant": p.is_constant,
                        "is_storage": p.is_storage,
                        "type": str(p.type),
                    } for p in function.parameters
                ] if function.parameters is not None else [],
                "modifiers": [m.name for m in function.modifiers],
            })

        return functions_data

    @staticmethod
    def _get_contract_code(c: Contract) -> str:
        """Extract the source code of a smart contract"""
        src_mapping = c.source_mapping
        content: str = c.compilation_unit.core.source_code[src_mapping.filename.absolute]
        start = src_mapping.start
        end = src_mapping.start + src_mapping.length
        return content[start:end]

    @staticmethod
    def is_address(address: str) -> bool:
        if not isinstance(address, str):
            return False
        return address.startswith('0x') and len(address) == 42 and all(c in '0123456789abcdefABCDEF' for c in address[2:])

    @staticmethod
    def is_supported_network_prefix(network_prefix: str) -> bool:
        return network_prefix in ["mainet", "arbi", "poly", "mumbai", "avax", "ftm", "bsc", "optim"]

    @staticmethod
    def generate_api_key_dict(network_prefix: str, api_key: str) -> dict:
        out_dict = {}
        match network_prefix:
            case "mainet":
                out_dict["etherscan_api_key"] = api_key
            case "arbi":
                out_dict["arbiscan_api_key"] = api_key
            case "poly":
                out_dict["polygonscan_api_key"] = api_key
            case "mumbai":
                out_dict["test_polygonscan_api_key"] = api_key
            case "avax":
                out_dict["avax_api_key"] = api_key
            case "ftm":
                out_dict["ftmscan_api_key"] = api_key
            case "bsc":
                out_dict["bscan_api_key"] = api_key
            case "optim":
                out_dict["optim_api_key"] = api_key
            case _:
                out_dict["etherscan_api_key"] = api_key

        return out_dict

    @staticmethod
    def get_slither_from_address(address: str, network_prefix: str, api_key: str, args: List[str]) -> Slither:
        # Validate network prefix
        if not SlitherHelper.is_supported_network_prefix(network_prefix):
            raise ValueError(f"Unsupported network prefix: {network_prefix}")

        # Validate address
        if not SlitherHelper.is_address(address):
            raise ValueError(f"Invalid address: {address}")

        args = SlitherHelper.generate_api_key_dict(
            network_prefix, api_key)

        # Add extra slither args
        if args is not None:
            args.extend(args)

        s = Slither(
            f"{network_prefix}:{address}", **args)
        return s

    @staticmethod
    def get_contracts(slither: Slither):
        contracts_data = []
        contracts = slither.contracts

        if os.path.isdir(args.target):
            target_contracts = SlitherHelper._filter_contracts(
                contracts, args.contracts_dir, args.exclude_contract_paths)
            test_contracts = []

            if args.tests_dir is not None:
                test_contracts = SlitherHelper._filter_contracts(
                    contracts, args.tests_dir, args.exclude_test_paths)

            contracts = target_contracts + test_contracts

        for c in contracts:
            contracts_data.append(SlitherHelper._get_inheritance_tree(c))

        return contracts_data


if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        description='This program runs slither on a directory and writes the output to a specified file')

    parser.add_argument('--target', type=str, required=True,
                        help='The target directory')
    parser.add_argument('--out', type=str, required=True,
                        help='The file the slither output will be written to')
    parser.add_argument(
        '--slither-args', type=SlitherHelper._parse_dict, required=False, default={}, help='Extra arguments to be passed to slither')
    parser.add_argument('--onchain', action='store_true',
                        help='Whether the target is an onchain contract')
    parser.add_argument('--network-prefix', type=str, required=False,
                        help='The network prefix of the onchain contract')
    parser.add_argument('--api-key', type=str, required=False,
                        help='The API key to use for the onchain contract')
    parser.add_argument('--contracts-dir', type=str, required=False,
                        help='The directory containing your target contracts')
    parser.add_argument('--exclude-contract-paths', type=SlitherHelper._parse_array, required=False, default=[
    ], help='Paths to be excluded from the target contracts')
    parser.add_argument('--tests-dir', type=str, required=False,
                        help='The directory containing your target contract tests')
    parser.add_argument('--exclude-test-paths', type=SlitherHelper._parse_array,
                        required=False, default=[], help='Paths to be excluded from the tests')

    args = parser.parse_args()

    target = args.target
    output_file = args.out
    kwargs = args.slither_args

    # Validate args
    if not args.onchain and not os.path.exists(target):
        raise ValueError(f"Target directory '{target}' does not exist.")
    if args.onchain and not args.api_key:
        raise ValueError(
            f"API key must be specified if target is an onchain contract.")

    # Get slither instance
    slither = None
    if args.onchain:
        slither = SlitherHelper.get_slither_from_address(
            address=target, network_prefix=args.network_prefix if args.network_prefix != "" else "mainet", api_key=args.api_key, args=kwargs)
    else:
        slither = Slither(
            target=target, kwargs=kwargs)

    contracts_data = SlitherHelper.get_contracts(slither=slither)

    # Write contracts data to output file
    with open(output_file, "w") as file:
        json.dump({"contracts": contracts_data}, file, indent=4)
