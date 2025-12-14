#!/usr/bin/env python3

# Problem 5 analysis
targets = [184, 136, 19, 48, 143, 50, 68, 3, 53, 165]
buttons = [
    [0,1,4,9],           # Button 0
    [0,2,3,4],           # Button 1
    [0,6,9],             # Button 2
    [0,2,4,6,8,9],       # Button 3
    [0,1,2,3,4,7,8],     # Button 4
    [0,1,2,4,5,6,8,9],   # Button 5
    [0,5,6,9],           # Button 6
    [0,3,8,9],           # Button 7
    [3,5,6,8],           # Button 8
    [1,2],               # Button 9
    [0,1,4,5,6],         # Button 10
    [0,3,7],             # Button 11
    [0,1,4,5,6,9]        # Button 12
]

print("=== PROBLEM 5 ANALYSIS ===")
print(f"Targets: {targets}")
print(f"Total work needed: {sum(targets)}")
print()

# Analyze button coverage
print("Button analysis:")
for i, button in enumerate(buttons):
    total_help = sum(targets[j] for j in button)
    efficiency = total_help / len(button)
    print(f"Button {i:2d}: affects {button} -> total help: {total_help:3d}, efficiency: {efficiency:.1f}")

print()

# Analyze counter coverage
print("Counter analysis:")
for i in range(len(targets)):
    affecting_buttons = [j for j, button in enumerate(buttons) if i in button]
    print(f"Counter {i}: target={targets[i]:3d}, affected by buttons {affecting_buttons}")

print()

# Look for potential bottlenecks
print("=== BOTTLENECK ANALYSIS ===")

# Check if any counter has very limited button options
for i in range(len(targets)):
    affecting_buttons = [j for j, button in enumerate(buttons) if i in button]
    if len(affecting_buttons) <= 2:
        print(f"BOTTLENECK: Counter {i} (target={targets[i]}) only affected by {len(affecting_buttons)} buttons: {affecting_buttons}")

# Check for counters with very high targets but limited options
high_target_counters = [(i, targets[i]) for i in range(len(targets)) if targets[i] > 100]
for i, target in high_target_counters:
    affecting_buttons = [j for j, button in enumerate(buttons) if i in button]
    print(f"High target: Counter {i} (target={target}) affected by {len(affecting_buttons)} buttons: {affecting_buttons}")

print()

# Look for mathematical impossibilities or constraints
print("=== CONSTRAINT ANALYSIS ===")

# Create coefficient matrix (button effects on each counter)
A = []
for i in range(len(targets)):
    row = []
    for j in range(len(buttons)):
        if i in buttons[j]:
            row.append(1)
        else:
            row.append(0)
    A.append(row)

print("Button-Counter Matrix (1 = button affects counter):")
print("     ", " ".join(f"B{i:2d}" for i in range(len(buttons))))
for i in range(len(targets)):
    row = " ".join(f"{A[i][j]:3d}" for j in range(len(buttons)))
    print(f"C{i:2d}: {row} | target={targets[i]}")

print()

# Check for potential linear dependencies or constraints
print("=== MATHEMATICAL CONSTRAINTS ===")

# Look for counters that might force certain button press counts
for i in range(len(targets)):
    affecting_buttons = [j for j, button in enumerate(buttons) if i in button]
    if len(affecting_buttons) == 1:
        print(f"FORCED: Counter {i} requires exactly {targets[i]} presses of button {affecting_buttons[0]}")
    elif len(affecting_buttons) == 2:
        print(f"CONSTRAINED: Counter {i} (target={targets[i]}) can only use buttons {affecting_buttons}")

print()

# Detailed analysis of the critical constraint (Counter 7)
print("=== COUNTER 7 DETAILED ANALYSIS ===")
print("Counter 7 (target=3) can only be incremented by buttons 4 and 11")
print("Button 4 affects:", buttons[4])
print("Button 11 affects:", buttons[11])
print()
print("This creates a constraint: button4_presses + button11_presses = 3")
print("Since other counters need button 4 heavily, this severely constrains the solution space")

# Check for very high efficiency buttons that might be overused
print()
print("Most efficient buttons (sorted by efficiency):")
button_efficiencies = []
for i, button in enumerate(buttons):
    total_help = sum(targets[j] for j in button)
    efficiency = total_help / len(button)
    button_efficiencies.append((i, efficiency, total_help))

button_efficiencies.sort(key=lambda x: x[1], reverse=True)
for i, efficiency, total_help in button_efficiencies:
    print(f"Button {i:2d}: efficiency {efficiency:5.1f} (total help: {total_help:3d}) affects {buttons[i]}")

print()
print("=== WHY THIS IS HARD ===")
print("1. Counter 7 creates a tight constraint (only 2 buttons, target=3)")
print("2. Button 4 is highly efficient but constrained by Counter 7")
print("3. Many high-value targets (184, 136, 143, 165) create large search space")
print("4. Greedy algorithms get trapped because optimal solution requires")
print("   careful balancing of button 4 usage between Counter 7 constraint")
print("   and other high-value counters")