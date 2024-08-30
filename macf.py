#!/usr/bin/env python3
# coding: utf-8

"""
macf: A simple MAC address formatting utility tool by Jake Shaw
      (https://github.com/codejake).

Usage: macf <mac address> or 'cat foo.txt | macf'
"""

import re
import sys
from typing import Callable, List


def clean(unclean: str) -> str:
    """Remove any non-hexadecimal characters from the input string."""
    return re.sub(r'[^0-9a-fA-F]', '', unclean)


def format_mac(mac: str, separator: str, case: str = 'lower') -> str:
    """Format a MAC address with the given separator and case."""
    formatted = f"{mac[:2]}{separator}{mac[2:4]}{separator}{mac[4:6]}{separator}" \
                f"{mac[6:8]}{separator}{mac[8:10]}{separator}{mac[10:12]}"
    return formatted.lower() if case == 'lower' else formatted.upper()


def get_cisco_style(mac: str) -> str:
    """Format a MAC address in Cisco style."""
    return f"{mac[:4]}.{mac[4:8]}.{mac[8:12]}".lower()


def get_formatted_macs(mac: str) -> List[str]:
    """Generate all formatted versions of a MAC address."""
    clean_mac = clean(mac)
    return [
        format_mac(clean_mac, ':'),
        clean_mac,
        clean_mac.upper(),
        format_mac(clean_mac, ':', 'upper'),
        get_cisco_style(clean_mac),
        format_mac(clean_mac, '-'),
        format_mac(clean_mac, '-', 'upper')
    ]


def process_mac(mac: str, output_func: Callable[[List[str]], None]) -> None:
    """Process a single MAC address and output the results."""
    formatted_macs = get_formatted_macs(mac)
    output_func(formatted_macs)


def print_formatted(formatted_macs: List[str]) -> None:
    """Print formatted MAC addresses with newlines between them."""
    print('\n\n'.join(formatted_macs))


def print_formatted_inline(formatted_macs: List[str]) -> None:
    """Print formatted MAC addresses on a single line."""
    print(' '.join(formatted_macs))


def main() -> None:
    if len(sys.argv) == 2:
        process_mac(sys.argv[1], print_formatted)
    else:
        for line in sys.stdin:
            process_mac(line.strip(), print_formatted_inline)


if __name__ == "__main__":
    main()
