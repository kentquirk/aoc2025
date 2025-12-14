#!/usr/bin/env python3

import re
import sys

try:
    import pulp
    USE_PULP = True
except ImportError:
    print("PuLP not installed. Installing...")
    import subprocess
    subprocess.check_call([sys.executable, '-m', 'pip', 'install', 'pulp'])
    import pulp
    USE_PULP = True

def parse_machine(line):
    """Parse a machine line from the input file"""
    # Extract buttons (parentheses) and targets (braces)
    button_pattern = r'\(([^)]+)\)'
    target_pattern = r'\{([^}]+)\}'

    button_matches = re.findall(button_pattern, line)
    target_match = re.search(target_pattern, line)

    buttons = []
    for match in button_matches:
        button = [int(x.strip()) for x in match.split(',')]
        buttons.append(button)

    targets = []
    if target_match:
        targets = [int(x.strip()) for x in target_match.group(1).split(',')]

    return buttons, targets

def solve_machine_ilp(buttons, targets, machine_id=None):
    """Solve a single machine using PuLP ILP solver"""
    # Create the problem
    prob = pulp.LpProblem(f"Machine_{machine_id}", pulp.LpMinimize)

    # Decision variables (number of times each button is pressed)
    x = [pulp.LpVariable(f"button_{i}", lowBound=0, cat='Integer')
         for i in range(len(buttons))]

    # Objective: minimize total button presses
    prob += pulp.lpSum(x)

    # Constraints: each counter must reach exactly its target
    for counter in range(len(targets)):
        affecting_buttons = [i for i, button in enumerate(buttons)
                           if counter in button]

        if not affecting_buttons:
            print(f"Error: Counter {counter} cannot be affected by any button!")
            return None, None

        # Sum of button presses for this counter equals target
        prob += pulp.lpSum([x[i] for i in affecting_buttons]) == targets[counter]

    # Solve
    if machine_id:
        print(f"Solving machine {machine_id}... ", end="", flush=True)

    # Use CBC solver with time limit
    prob.solve(pulp.PULP_CBC_CMD(msg=0, timeLimit=30))

    if prob.status == pulp.LpStatusOptimal:
        solution = [int(x[i].varValue) for i in range(len(buttons))]
        total_presses = sum(solution)

        if machine_id:
            print(f"{total_presses} presses")

        # Verify solution
        verify_solution(buttons, targets, solution, machine_id)

        return total_presses, solution

    elif prob.status == pulp.LpStatusNotSolved:
        if machine_id:
            print(f"Time limit reached")
        return None, None

    else:
        if machine_id:
            print(f"No solution found (status: {pulp.LpStatus[prob.status]})")
        return None, None

def verify_solution(buttons, targets, solution, machine_id=None):
    """Verify that a solution is correct"""
    achieved = [0] * len(targets)

    for button_idx, presses in enumerate(solution):
        for counter in buttons[button_idx]:
            achieved[counter] += presses

    for i, (target, actual) in enumerate(zip(targets, achieved)):
        if target != actual:
            print(f"ERROR in machine {machine_id}: Counter {i} target={target} but achieved={actual}")
            return False

    return True

def solve_all_machines(filename):
    """Solve all machines in the input file"""
    total_presses = 0
    solved_count = 0
    failed_machines = []

    try:
        with open(filename, 'r') as f:
            lines = [line.strip() for line in f if line.strip()]
    except FileNotFoundError:
        print(f"Error: Could not find file {filename}")
        return None

    print(f"Found {len(lines)} machines to solve\n")

    for i, line in enumerate(lines, 1):
        try:
            buttons, targets = parse_machine(line)

            if not buttons or not targets:
                print(f"Machine {i}: Failed to parse input")
                failed_machines.append(i)
                continue

            presses, solution = solve_machine_ilp(buttons, targets, i)

            if presses is not None:
                total_presses += presses
                solved_count += 1
            else:
                failed_machines.append(i)

        except Exception as e:
            print(f"Machine {i}: Error - {e}")
            failed_machines.append(i)

    print(f"\n=== RESULTS ===")
    print(f"Solved: {solved_count}/{len(lines)} machines")
    print(f"Total button presses: {total_presses}")

    if failed_machines:
        print(f"Failed machines: {failed_machines}")
    else:
        print("All machines solved successfully!")

    return total_presses if not failed_machines else None

def solve_sample():
    """Test with the sample data from sample.txt"""
    print("=== TESTING WITH SAMPLE DATA ===")

    try:
        with open("data/sample.txt", 'r') as f:
            sample_lines = [line.strip() for line in f if line.strip()]
    except FileNotFoundError:
        print("Warning: data/sample.txt not found, skipping sample test")
        return True  # Don't fail if sample file doesn't exist

    total = 0
    for i, line in enumerate(sample_lines, 1):
        buttons, targets = parse_machine(line)
        presses, solution = solve_machine_ilp(buttons, targets, i)
        if presses:
            total += presses

    print(f"Sample total: {total} (expected: 33)")
    return total == 33

def main():
    if len(sys.argv) > 1:
        filename = sys.argv[1]
        print(f"Solving {filename}...")
        result = solve_all_machines(filename)
    else:
        # Test with sample first if it exists
        print("Testing with sample data...")
        if solve_sample():
            print("✅ Sample test passed!\n")
        else:
            print("❌ Sample test failed!\n")

        # Then solve the main input file
        filename = "data/input.txt"
        print(f"Solving {filename}...")
        result = solve_all_machines(filename)

    if result is not None:
        print(f"\n🎉 Final Answer: {result}")
    else:
        print(f"\n❌ Some machines could not be solved")

if __name__ == "__main__":
    main()