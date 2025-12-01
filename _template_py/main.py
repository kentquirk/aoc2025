#! /usr/bin/env python3

import sys
import os
import re
import itertools

def part1(data):
    return 0

def part2(data):
    return 0

if __name__ == "__main__":
    file = "sample"
    if len(sys.argv) > 1:
        filename = sys.argv[1]

    filename = f"data/{filename}.txt"

    f = open(sys.argv[1])
    lines = f.readlines()
    data = [int(l) for l in lines]
    print(part1(data))
