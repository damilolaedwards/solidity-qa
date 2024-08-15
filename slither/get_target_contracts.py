"""Reader for smart contracts in Solidity projects through Slither"""
import sys
import time

from typing import Any, Iterator, List

from llama_index.core.readers.base import Document
from slither import Slither
from slither.core.declarations import Contract


def _get_contract_code(c: Contract) -> str:
    """Extract the source code of a smart contract"""
    src_mapping = c.source_mapping
    content: str = c.compilation_unit.core.source_code[src_mapping.filename.absolute]
    start = src_mapping.start
    end = src_mapping.start + src_mapping.length
    return content[start:end]


def _get_abstract_code(c: Contract) -> str:
    """Extract the "abstract" code of a smart contract (replaces function bodies with ...)"""
    # TODO: implement me
    return ""


def _get_contracts() -> Iterator[Document]:
    """Get an iterator of llama_index Documents based on contracts' codes"""
    for c in slither.contracts:
        # Get an interface-like code for the whole contract
        abstract_code = _get_abstract_code(c)
        yield Document(
            text=abstract_code,
            metadata={
                "contract_name": c.name,
                "filepath": c.file_scope.filename.relative,
            },
        )

        code = _get_contract_code(c)
        print(code)
        yield Document(
            text=code,
            metadata={
                "contract_name": c.name,
                "filepath": c.file_scope.filename.relative,
            },
        )


def load_data() -> List[Document]:
    """Load contracts from the Solidity projects"""
    return list(_get_contracts())


if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: python solidity-qa.py <target>")
        sys.exit(1)

    # Obtain target from CLI arguments
    target = sys.argv[1]
    slither = Slither(target)

    load_data()
