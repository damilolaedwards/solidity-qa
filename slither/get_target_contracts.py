"""Reader for smart contracts in Solidity projects through Slither"""
import json
import sys
import os
import argparse
from typing import List
from io import StringIO

from slither import Slither
from slither.core.declarations import Function, Contract


def _parse_exclude_paths(exclude_paths):
    return exclude_paths.split(',')


def _filter_contracts(contracts: List[Contract], target_dir: str, exclude_paths: List[str]):
    # Check if target_dir exists and is a directory
    if not os.path.exists(target_dir):
        raise ValueError(f"Target directory '{target_dir}' does not exist.")
    if not os.path.isdir(target_dir):
        raise ValueError(f"Target path '{target_dir}' is not a directory.")

    # Normalize target_dir path
    target_dir = os.path.abspath(target_dir)

    # Normalize and filter out non-existing exclude paths
    exclude_paths = [os.path.abspath(
        path) for path in exclude_paths if os.path.exists(path)]

    # Helper function to determine if a contract's path falls under target_dir
    def is_under_target_dir(contract_path):
        contract_path = os.path.abspath(contract_path)
        return contract_path.startswith(target_dir)

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


def _get_inheritance_tree(contract: Contract):
    inheritance_tree = {
        "id": contract.id,
        "name": contract.name,
        "code": _get_contract_code(contract),
        "is_abstract": contract.is_abstract,
        "is_interface": contract.is_interface,
        "is_library": contract.is_library,
        "functions": _get_functions_data(contract.functions),
        "inherited_contracts": []
    }
    for inherited_contract in contract.inheritance:
        inheritance_tree["inherited_contracts"].append(
            _get_inheritance_tree(inherited_contract))

    return inheritance_tree


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


def _get_contract_code(c: Contract) -> str:
    """Extract the source code of a smart contract"""
    src_mapping = c.source_mapping
    content: str = c.compilation_unit.core.source_code[src_mapping.filename.absolute]
    start = src_mapping.start
    end = src_mapping.start + src_mapping.length
    return content[start:end]


def _get_contracts():
    contracts_data = []
    contracts = slither.contracts

    target_contracts = _filter_contracts(
        contracts, args.contracts_dir, args.exclude_contract_paths)
    test_contracts = []

    if args.tests_dir is not None:
        test_contracts = _filter_contracts(
            contracts, args.tests_dir, args.exclude_test_paths)

    contracts = target_contracts + test_contracts

    for c in contracts:
        contracts_data.append(_get_inheritance_tree(c))

    return {"contracts": contracts_data}


if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        description='This program runs slither on a directory and writes the output to a specified file')

    parser.add_argument('--target', type=str, required=True,
                        help='The target directory')
    parser.add_argument('--out', type=str, required=True,
                        help='The file the slither output will be written to')
    parser.add_argument("--contracts-dir", type=str, required=True,
                        help='The directory containing your target contracts')
    parser.add_argument("--exclude-contract-paths", type=_parse_exclude_paths, required=False, default=[
    ], help='Paths to be excluded from the target contracts')
    parser.add_argument("--tests-dir", type=str, required=False,
                        help='The directory containing your target contract tests')
    parser.add_argument("--exclude-test-paths", type=_parse_exclude_paths,
                        required=False, default=[], help='Paths to be excluded from the tests')

    args = parser.parse_args()

    target = args.target
    output_file = args.out

    slither = Slither(target)
    with open(output_file, "w") as file:
        json.dump(_get_contracts(), file, indent=4)
