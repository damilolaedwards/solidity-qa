import sys
import json
from typing import List
from slither.slither import Slither
from slither.core.declarations import Contract, Function
import argparse

parser = argparse.ArgumentParser(
    description='This program runs slither on a directory and writes the output to a specified file')

parser.add_argument('--target', type=str, required=True,
                    help='The target directory')
parser.add_argument('--out', type=str, required=True,
                    help='The file the slither output will be written to')
parser.add_argument('--include-interfaces', type=bool, required=False, default=False,
                    help='Whether interfaces should be included in the output')
parser.add_argument('--include-libraries', type=bool, required=False, default=False,
                    help='Whether libraries should be included in the output')
parser.add_argument('--include-abstract', type=bool, required=False, default=False,
                    help='Whether abstract contracts should be included in the output')

args = parser.parse_args()

target = args.target
output_file = args.out
include_interfaces = args.include_interfaces
include_libraries = args.include_libraries
include_abstract = args.include_abstract


def get_functions_data(functions: List[Function]):
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


def get_inheritance_tree(contract: Contract):
    inheritance_tree = {
        "id": contract.id,
        "name": contract.name,
        "path": contract.file_scope.filename.relative,
        "is_abstract": contract.is_abstract,
        "is_interface": contract.is_interface,
        "is_library": contract.is_library,
        "functions": get_functions_data(contract.functions),
        "inherited_contracts": []
    }
    for inherited_contract in contract.inheritance:
        if (not inherited_contract.is_interface or include_interfaces) and (not inherited_contract.is_library or include_libraries) and (not inherited_contract.is_abstract or include_abstract):
            inheritance_tree["inherited_contracts"].append(
                get_inheritance_tree(inherited_contract))

    return inheritance_tree


slither = None

try:
    slither = Slither(target)

    contracts_data = []
    for contract in slither.contracts:
        contract.is_abstract
        if (not contract.is_interface or include_interfaces) and (not contract.is_library or include_libraries) and (not contract.is_abstract or include_abstract):
            contracts_data.append(get_inheritance_tree(contract))

    with open(output_file, "w") as file:
        json.dump(contracts_data, file, indent=4)


except Exception as e:
    print(e.__str__())
    sys.exit(-1)
