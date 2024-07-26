import sys
import json
import slither.slither
from slither.slither import Slither

if len(sys.argv) != 3:
    print("Usage: python3 slither.py target output_file")
    sys.exit(-1)

target = sys.argv[1]
output_file = sys.argv[2]


def get_functions_data(functions):
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


def get_inheritance_tree(contract):
    inheritance_tree = {
        "id": contract.id,
        "name": contract.name,
        "functions": get_functions_data(contract.functions),
        "inherited_contracts": []
    }
    for inherited_contract in contract.inheritance:
        inheritance_tree["inherited_contracts"].append(
            get_inheritance_tree(inherited_contract))

    return inheritance_tree


slither = None

try:
    slither = Slither(target)

    contracts_data = []
    for contract in slither.contracts:
        if not contract.is_interface:
            contracts_data.append(get_inheritance_tree(contract))

    with open(output_file, "w") as file:
        json.dump(contracts_data, file, indent=4)


except Exception as e:
    print(e.__str__())
    sys.exit(-1)
