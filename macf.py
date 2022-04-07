#!/usr/bin/env python3
# coding: utf-8

#
# macf: A dumb little MAC address formatting utility tool by Jake.
#
# Usage: macf <mac address> or 'cat foo.txt | macf'
#

import os
import re
import sys


# Assumes valid input
def clean(unclean):
    return re.sub(r'\W+', '', unclean)

# aa:bb:cc:dd:11:22


def get_colon_style_lower(input):
    return f"{input[:2]}:{input[2:4]}:{input[4:6]}:{input[6:8]}:{input[8:10]}:{input[10:12]}".lower()

# AA:BB:CC:DD:11:22


def get_colon_style_upper(input):
    return f"{input[:2]}:{input[2:4]}:{input[4:6]}:{input[6:8]}:{input[8:10]}:{input[10:12]}".upper()

# aa-bb-cc-dd-11-22


def get_hyphen_style_lower(input):
    return f"{input[:2]}-{input[2:4]}-{input[4:6]}-{input[6:8]}-{input[8:10]}-{input[10:12]}".lower()

# AA-BB-CC-DD-11-22


def get_hyphen_style_upper(input):
    return f"{input[:2]}-{input[2:4]}-{input[4:6]}-{input[6:8]}-{input[8:10]}-{input[10:12]}".upper()

# aabb.ccdd.1122


def get_cisco_style(input):
    return f"{input[:4]}.{input[4:8]}.{input[8:12]}".lower()


def main():
    if len(sys.argv) == 2:
        clean_addr = clean(sys.argv[1])
        print("")
        print(get_colon_style_lower(clean_addr), end="\n\n")
        print(clean_addr, end="\n\n")
        print(clean_addr.upper(), end="\n\n")
        print(get_colon_style_upper(clean_addr), end="\n\n")
        print(get_cisco_style(clean_addr), end="\n\n")
        print(get_hyphen_style_lower(clean_addr), end="\n\n")
        print(get_hyphen_style_upper(clean_addr), end="\n\n")
    else:
        for line in sys.stdin:
            clean_addr = clean(line)
            print(get_colon_style_lower(clean_addr), end=' ')
            print(clean_addr, end=' ')
            print(get_colon_style_upper(clean_addr), end=' ')
            print(get_hyphen_style_lower(clean_addr), end=' ')
            print(get_hyphen_style_upper(clean_addr), end=' ')
            print(get_cisco_style(clean_addr))


if __name__ == "__main__":
    main()
