package main

import (
	"fmt"
	"io"
	"log"
	"math/rand/v2"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type bits int

func (b bits) String() string {
	return fmt.Sprintf("%03x", int(b))
}

func (b bits) asBits(n int) string {
	var result strings.Builder
	for i := 0; i < n; i++ {
		if b&(1<<i) != 0 {
			result.WriteByte('#')
		} else {
			result.WriteByte('.')
		}
	}
	return result.String()
}

// we call them switches in part 1 and buttons in part 2
type machine struct {
	nbits    int
	lamps    bits
	buttons  [][]int // each button is a list of indices
	switches []bits  // the bit patterns for each switch
	joltages []int
	state    bits
	children []*machine
	jstate   []int
	presses  []int
}

func (m machine) String() string {
	return fmt.Sprintf("lamps: %s, switches: %v, joltages: %v", m.lamps.asBits(m.nbits), m.switches, m.joltages)
}

func (m machine) JString() string {
	return fmt.Sprintf("jstate: %v, presses: %v", m.jstate, m.presses)
}

func (m machine) print() {
	fmt.Printf("Lamps - %s\n", m.lamps.asBits(m.nbits))
	for i, s := range m.switches {
		fmt.Printf(" Sw%02d   %s\n", i, s.asBits(m.nbits))
	}
	fmt.Printf(" J      {%v}\n", m.joltages)
}

func (m machine) pressSwitch(i int) (*machine, bool) {
	newState := m.state ^ m.switches[i]
	if newState == m.lamps {
		return nil, true
	}
	child := &machine{
		nbits:    m.nbits,
		lamps:    m.lamps,
		joltages: m.joltages,
		state:    newState,
	}
	// remove the switch we just used from the list
	for j := range m.switches {
		if j == i {
			continue
		}
		child.switches = append(child.switches, m.switches[j])
	}
	// sort the switches so we can hash them
	sort.Slice(child.switches, func(i, j int) bool {
		return child.switches[i] < child.switches[j]
	})
	return child, false
}

// This is a breadth-first search of the solution space.
// we push each of the switches and clone children, but
// stop if pushing any of them solves the machine.
// If we've tried all and we're not done, we go deeper by
// visiting each child.
func (m machine) search(depth int) int {
	// fmt.Printf("%*s at depth %d: %s\n", depth*2, "", depth, m)
	for i := range m.switches {
		n, done := m.pressSwitch(i)
		if done {
			return depth
		}
		m.children = append(m.children, n)
	}
	dmin := 1000
	for _, child := range m.children {
		if d := child.search(depth + 1); d != -1 {
			if d < dmin {
				dmin = d
			}
		}
	}
	if dmin != 1000 {
		return dmin
	}
	return -1
}

func (m *machine) pressButton() bool {
	deltas := make([]int, len(m.joltages))
	for i, st := range m.jstate {
		deltas[i] = m.joltages[i] - st
	}

	best := -1
	bestScore := -1000000 // Allow negative scores

	for i, button := range m.buttons {
		score := 0
		for _, ix := range button {
			delta := deltas[ix]
			if delta > 0 {
				// Weight smaller targets higher to avoid getting stuck
				weight := float64(delta) / float64(m.joltages[ix])
				score += int(10 * weight)
			} else if delta == 0 {
				score -= 5 // Moderate penalty
			} else {
				score -= 20 // Heavy penalty for overshooting
			}
		}

		if score > bestScore {
			bestScore = score
			best = i
		}
	}

	if best == -1 {
		return false // No useful buttons available
	}

	// Press the best button
	m.presses = append(m.presses, best)
	for _, i := range m.buttons[best] {
		m.jstate[i]++
	}
	return true
}

func (m *machine) hybridOptimalSolve() int {
	// Strategy 1: Try deterministic greedy variants
	best := m.multiGreedySolve()

	// Strategy 2: Try randomized greedy variants (fast and often effective)
	if candidate := m.randomizedGreedySolve(best); candidate < best {
		best = candidate
	}

	// Strategy 3: For hard cases, give randomized greedy more attempts
	tries := 0
	for best == 1000000 { // Only if all strategies completely failed
		best = m.intensiveRandomizedGreedy()
		tries++
		if tries > 10 {
			break
		}
	}

	if best == 1000000 {
		// Could not find a solution
		return -1
	}
	return best
}

func (m *machine) multiGreedySolve() int {
	strategies := []func() int{
		m.greedyByEfficiency,
		m.greedyByBottleneck,
		m.greedyByRatio,
		m.greedyByLargestFirst,
	}
	strategyNames := []string{
		"efficiency", "bottleneck", "ratio", "largest-first",
	}

	best := 1000000
	for i, strategy := range strategies {
		result := strategy()
		fmt.Printf("  %s strategy: %d presses\n", strategyNames[i], result)
		if result < best {
			best = result
		}
	}
	return best
}

func (m *machine) randomizedGreedySolve(currentBest int) int {
	fmt.Printf("  Trying randomized greedy strategies...\n")
	best := currentBest

	// Try multiple randomized runs of each strategy
	strategies := []func(int) int{
		m.randomizedGreedyByEfficiency,
		m.randomizedGreedyByWeightedChoice,
		m.randomizedGreedyWithNoise,
		m.randomizedGreedyWithExploration,
	}

	strategyNames := []string{
		"rand-efficiency", "weighted-choice", "noisy-greedy", "exploration",
	}

	for i, strategy := range strategies {
		attempts := 10
		if currentBest == 1000000 { // If deterministic greedy completely failed
			attempts = 25 // Give more attempts for hard cases
		}

		for attempt := 0; attempt < attempts; attempt++ {
			result := strategy(best)
			if result < best {
				best = result
				fmt.Printf("  %s found better solution: %d\n", strategyNames[i], best)
				break // Move to next strategy once we find improvement
			}
		}
	}

	if best < currentBest {
		fmt.Printf("  Randomized greedy improved to %d\n", best)
	} else {
		fmt.Printf("  Randomized greedy found no improvement\n")
	}

	return best
}

func (m *machine) randomizedGreedyByEfficiency(targetBest int) int {
	m.reset()
	maxIterations := 5000
	iterations := 0

	for m.check() != Equal && iterations < maxIterations {
		iterations++
		if m.check() == TooHigh {
			m.unpressButton()
			continue
		}

		// Build candidate list with scores
		type candidate struct {
			button int
			score  int
		}
		var candidates []candidate

		deltas := make([]int, len(m.joltages))
		for i, st := range m.jstate {
			deltas[i] = m.joltages[i] - st
		}

		for i, button := range m.buttons {
			score := 0
			useful := false

			for _, ix := range button {
				delta := deltas[ix]
				if delta > 0 {
					score += delta
					useful = true
				} else if delta == 0 {
					score -= 5
				} else {
					score -= 20
				}
			}

			if useful && score > 0 {
				candidates = append(candidates, candidate{i, score})
			}
		}

		if len(candidates) == 0 {
			return 1000000
		}

		// Choose randomly from top 30% of candidates
		sort.Slice(candidates, func(i, j int) bool {
			return candidates[i].score > candidates[j].score
		})

		topCount := max(1, len(candidates)*3/10) // Top 30%
		chosen := candidates[rand.IntN(topCount)]

		m.presses = append(m.presses, chosen.button)
		for _, counter := range m.buttons[chosen.button] {
			m.jstate[counter]++
		}

		// Early termination if we're not improving
		if len(m.presses) > targetBest {
			return 1000000
		}
	}

	if iterations >= maxIterations {
		return 1000000
	}
	return len(m.presses)
}

func (m *machine) randomizedGreedyByWeightedChoice(targetBest int) int {
	m.reset()
	maxIterations := 5000
	iterations := 0

	for m.check() != Equal && iterations < maxIterations {
		iterations++
		if m.check() == TooHigh {
			m.unpressButton()
			continue
		}

		// Build weighted candidate list
		var candidates []int
		var weights []int

		for i, button := range m.buttons {
			weight := 0
			useful := false

			for _, counter := range button {
				need := m.joltages[counter] - m.jstate[counter]
				if need > 0 {
					weight += need * need // Square to emphasize urgent needs
					useful = true
				} else if need < 0 {
					weight -= 10
				}
			}

			if useful && weight > 0 {
				candidates = append(candidates, i)
				weights = append(weights, weight)
			}
		}

		if len(candidates) == 0 {
			return 1000000
		}

		// Weighted random selection
		totalWeight := 0
		for _, w := range weights {
			totalWeight += w
		}

		r := rand.IntN(totalWeight)
		chosen := 0
		for i, w := range weights {
			r -= w
			if r <= 0 {
				chosen = candidates[i]
				break
			}
		}

		m.presses = append(m.presses, chosen)
		for _, counter := range m.buttons[chosen] {
			m.jstate[counter]++
		}

		if len(m.presses) > targetBest {
			return 1000000
		}
	}

	if iterations >= maxIterations {
		return 1000000
	}
	return len(m.presses)
}

func (m *machine) randomizedGreedyWithNoise(targetBest int) int {
	m.reset()
	maxIterations := 5000
	iterations := 0

	for m.check() != Equal && iterations < maxIterations {
		iterations++
		if m.check() == TooHigh {
			m.unpressButton()
			continue
		}

		bestButton := -1
		bestScore := -1000000

		for i, button := range m.buttons {
			score := 0
			useful := false

			for _, counter := range button {
				need := m.joltages[counter] - m.jstate[counter]
				if need > 0 {
					score += need
					useful = true
				} else if need == 0 {
					score -= 5
				} else {
					score -= 20
				}
			}

			if useful {
				// Add random noise (-5 to +5)
				noise := rand.IntN(11) - 5
				score += noise

				if score > bestScore {
					bestScore = score
					bestButton = i
				}
			}
		}

		if bestButton == -1 {
			return 1000000
		}

		m.presses = append(m.presses, bestButton)
		for _, counter := range m.buttons[bestButton] {
			m.jstate[counter]++
		}

		if len(m.presses) > targetBest {
			return 1000000
		}
	}

	if iterations >= maxIterations {
		return 1000000
	}
	return len(m.presses)
}

func (m *machine) randomizedGreedyWithExploration(targetBest int) int {
	m.reset()
	maxIterations := 5000
	iterations := 0

	for m.check() != Equal && iterations < maxIterations {
		iterations++
		if m.check() == TooHigh {
			m.unpressButton()
			continue
		}

		// 80% exploitation, 20% exploration
		if rand.Float64() < 0.8 {
			// Exploitation: choose best button
			bestButton := -1
			bestScore := -1000000

			for i, button := range m.buttons {
				score := 0
				useful := false

				for _, counter := range button {
					need := m.joltages[counter] - m.jstate[counter]
					if need > 0 {
						score += need
						useful = true
					} else {
						score -= 10
					}
				}

				if useful && score > bestScore {
					bestScore = score
					bestButton = i
				}
			}

			if bestButton == -1 {
				return 1000000
			}

			m.presses = append(m.presses, bestButton)
			for _, counter := range m.buttons[bestButton] {
				m.jstate[counter]++
			}
		} else {
			// Exploration: choose randomly from useful buttons
			var useful []int

			for i, button := range m.buttons {
				hasUse := false
				for _, counter := range button {
					if m.joltages[counter] > m.jstate[counter] {
						hasUse = true
						break
					}
				}
				if hasUse {
					useful = append(useful, i)
				}
			}

			if len(useful) == 0 {
				return 1000000
			}

			chosen := useful[rand.IntN(len(useful))]
			m.presses = append(m.presses, chosen)
			for _, counter := range m.buttons[chosen] {
				m.jstate[counter]++
			}
		}

		if len(m.presses) > targetBest {
			return 1000000
		}
	}

	if iterations >= maxIterations {
		return 1000000
	}
	return len(m.presses)
}

func (m *machine) greedyByEfficiency() int {
	m.reset()
	maxIterations := 10000
	iterations := 0
	for m.check() != Equal && iterations < maxIterations {
		iterations++
		if m.check() == TooHigh {
			m.unpressButton()
			continue
		}

		bestButton := -1
		bestScore := -1000000

		deltas := make([]int, len(m.joltages))
		for i, st := range m.jstate {
			deltas[i] = m.joltages[i] - st
		}

		for i, button := range m.buttons {
			score := 0
			useful := false

			for _, ix := range button {
				delta := deltas[ix]
				if delta > 0 {
					score += delta
					useful = true
				} else if delta == 0 {
					score -= 5
				} else {
					score -= 20
				}
			}

			if useful && score > bestScore {
				bestScore = score
				bestButton = i
			}
		}

		if bestButton == -1 {
			return 1000000
		}

		m.presses = append(m.presses, bestButton)
		for _, counter := range m.buttons[bestButton] {
			m.jstate[counter]++
		}
	}
	if iterations >= maxIterations {
		return 1000000 // Failed due to timeout
	}
	return len(m.presses)
}

func (m *machine) greedyByBottleneck() int {
	m.reset()
	maxIterations := 10000
	iterations := 0
	for m.check() != Equal && iterations < maxIterations {
		iterations++
		if m.check() == TooHigh {
			m.unpressButton()
			continue
		}

		bestButton := -1
		bestScore := -1.0

		for i, button := range m.buttons {
			score := 0.0
			useful := false

			for _, counter := range button {
				remaining := float64(m.joltages[counter] - m.jstate[counter])
				if remaining > 0 {
					ratio := remaining / float64(m.joltages[counter])
					score += ratio * ratio
					useful = true
				} else if remaining < 0 {
					score -= 10.0
				}
			}

			if useful && score > bestScore {
				bestScore = score
				bestButton = i
			}
		}

		if bestButton == -1 {
			return 1000000
		}

		m.presses = append(m.presses, bestButton)
		for _, counter := range m.buttons[bestButton] {
			m.jstate[counter]++
		}
	}
	if iterations >= maxIterations {
		return 1000000
	}
	return len(m.presses)
}

func (m *machine) greedyByRatio() int {
	m.reset()
	maxIterations := 10000
	iterations := 0
	for m.check() != Equal && iterations < maxIterations {
		iterations++
		if m.check() == TooHigh {
			m.unpressButton()
			continue
		}

		bestButton := -1
		bestRatio := -1.0

		for i, button := range m.buttons {
			totalHelp := 0
			totalWaste := 0

			for _, counter := range button {
				remaining := m.joltages[counter] - m.jstate[counter]
				if remaining > 0 {
					totalHelp += remaining
				} else {
					totalWaste++
				}
			}

			if totalHelp > 0 {
				ratio := float64(totalHelp) / float64(len(button)+totalWaste)
				if ratio > bestRatio {
					bestRatio = ratio
					bestButton = i
				}
			}
		}

		if bestButton == -1 {
			return 1000000
		}

		m.presses = append(m.presses, bestButton)
		for _, counter := range m.buttons[bestButton] {
			m.jstate[counter]++
		}
	}
	if iterations >= maxIterations {
		return 1000000
	}
	return len(m.presses)
}

func (m *machine) greedyByLargestFirst() int {
	m.reset()
	maxIterations := 10000
	iterations := 0
	for m.check() != Equal && iterations < maxIterations {
		iterations++
		if m.check() == TooHigh {
			m.unpressButton()
			continue
		}

		// Find counter with largest remaining need
		maxNeed := 0
		for i, target := range m.joltages {
			need := target - m.jstate[i]
			if need > maxNeed {
				maxNeed = need
			}
		}

		bestButton := -1
		bestScore := -1

		for i, button := range m.buttons {
			score := 0
			affectsLargest := false

			for _, counter := range button {
				need := m.joltages[counter] - m.jstate[counter]
				if need == maxNeed {
					affectsLargest = true
					score += 100
				} else if need > 0 {
					score += need
				} else {
					score -= 10
				}
			}

			if affectsLargest && score > bestScore {
				bestScore = score
				bestButton = i
			}
		}

		if bestButton == -1 {
			return 1000000
		}

		m.presses = append(m.presses, bestButton)
		for _, counter := range m.buttons[bestButton] {
			m.jstate[counter]++
		}
	}
	if iterations >= maxIterations {
		return 1000000
	}
	return len(m.presses)
}

func (m *machine) intensiveRandomizedGreedy() int {
	fmt.Printf("  Trying intensive randomized greedy for hard case...\n")
	best := 1000000

	// Try all strategies with many more attempts and longer iterations
	strategies := []func(int) int{
		m.randomizedGreedyByEfficiency,
		m.randomizedGreedyByWeightedChoice,
		m.randomizedGreedyWithNoise,
		m.randomizedGreedyWithExploration,
		m.randomizedGreedyWithBacktrack,
		m.randomizedGreedyWithSimulatedAnnealing,
	}

	strategyNames := []string{
		"intensive-efficiency", "intensive-weighted", "intensive-noise",
		"intensive-exploration", "backtrack", "simulated-annealing",
	}

	for i, strategy := range strategies {
		for attempt := 0; attempt < 50; attempt++ { // Many more attempts
			result := strategy(1000000) // No early termination limit
			if result < best {
				best = result
				fmt.Printf("  %s found solution: %d\n", strategyNames[i], best)
			}
		}
	}

	if best < 1000000 {
		fmt.Printf("  Intensive search found solution: %d\n", best)
	} else {
		fmt.Printf("  Intensive search failed to find solution\n")
	}

	return best
}

func (m *machine) randomizedGreedyWithBacktrack(targetBest int) int {
	m.reset()
	maxIterations := 10000 // Longer for hard cases
	iterations := 0

	for m.check() != Equal && iterations < maxIterations {
		iterations++
		if m.check() == TooHigh {
			m.unpressButton()
			continue
		}

		// Every 100 steps, try backtracking if we're stuck
		if iterations%100 == 0 && len(m.presses) > 10 {
			// Remove last 3-5 button presses and try different path
			backtrack := min(5, len(m.presses)/2)
			for i := 0; i < backtrack; i++ {
				if len(m.presses) > 0 {
					m.unpressButton()
				}
			}
		}

		// Use noisy greedy selection
		bestButton := -1
		bestScore := -1000000

		for i, button := range m.buttons {
			score := 0
			useful := false

			for _, counter := range button {
				need := m.joltages[counter] - m.jstate[counter]
				if need > 0 {
					score += need
					useful = true
				} else if need == 0 {
					score -= 3
				} else {
					score -= 15
				}
			}

			if useful {
				// Add larger random noise for more exploration
				noise := rand.IntN(21) - 10 // -10 to +10
				score += noise

				if score > bestScore {
					bestScore = score
					bestButton = i
				}
			}
		}

		if bestButton == -1 {
			return 1000000
		}

		m.presses = append(m.presses, bestButton)
		for _, counter := range m.buttons[bestButton] {
			m.jstate[counter]++
		}
	}

	if iterations >= maxIterations {
		return 1000000
	}
	return len(m.presses)
}

func (m *machine) randomizedGreedyWithSimulatedAnnealing(targetBest int) int {
	m.reset()
	maxIterations := 8000
	iterations := 0

	for m.check() != Equal && iterations < maxIterations {
		iterations++
		if m.check() == TooHigh {
			m.unpressButton()
			continue
		}

		// Temperature decreases over time (starts hot, cools down)
		temperature := float64(maxIterations-iterations) / float64(maxIterations)

		bestButton := -1
		bestScore := -1000000
		var candidates []struct {
			button int
			score  int
		}

		for i, button := range m.buttons {
			score := 0
			useful := false

			for _, counter := range button {
				need := m.joltages[counter] - m.jstate[counter]
				if need > 0 {
					score += need
					useful = true
				} else if need == 0 {
					score -= 2
				} else {
					score -= 8
				}
			}

			if useful {
				candidates = append(candidates, struct {
					button int
					score  int
				}{i, score})

				if score > bestScore {
					bestScore = score
					bestButton = i
				}
			}
		}

		if len(candidates) == 0 {
			return 1000000
		}

		// With high temperature, accept worse choices; with low temperature, be greedy
		if temperature > 0.3 && len(candidates) > 1 {
			// Accept a random candidate with probability based on temperature
			if rand.Float64() < temperature {
				chosen := candidates[rand.IntN(len(candidates))]
				bestButton = chosen.button
			}
		}

		m.presses = append(m.presses, bestButton)
		for _, counter := range m.buttons[bestButton] {
			m.jstate[counter]++
		}
	}

	if iterations >= maxIterations {
		return 1000000
	}
	return len(m.presses)
}

func (m *machine) boundedIterativeDeepening(maxBound int) int {
	maxDepth := maxBound * 4 / 5
	if maxDepth > 30 { // Reduced from 50
		maxDepth = 30
	}

	fmt.Printf("  Trying iterative deepening up to depth %d...\n", maxDepth)
	for depth := 1; depth <= maxDepth; depth++ {
		deadline := time.Now().Add(2 * time.Second) // Reduced from 5
		if m.dfsWithTimeLimit(make([]int, len(m.joltages)), 0, depth, deadline) {
			fmt.Printf("  Found solution at depth %d\n", depth)
			return depth
		}
	}
	fmt.Printf("  Iterative deepening failed\n")
	return 1000000
}

func (m *machine) dfsWithTimeLimit(current []int, pressesUsed int, maxPresses int, deadline time.Time) bool {
	if time.Now().After(deadline) {
		return false
	}

	if m.isAtTarget(current) {
		return true
	}

	if pressesUsed >= maxPresses {
		return false
	}

	for _, button := range m.buttons {
		newState := make([]int, len(current))
		copy(newState, current)
		valid := true

		for _, counter := range button {
			newState[counter]++
			if newState[counter] > m.joltages[counter] {
				valid = false
				break
			}
		}

		if valid && m.dfsWithTimeLimit(newState, pressesUsed+1, maxPresses, deadline) {
			return true
		}
	}

	return false
}

func (m *machine) isAtTarget(current []int) bool {
	for i, target := range m.joltages {
		if current[i] != target {
			return false
		}
	}
	return true
}

func (m *machine) randomizedSearchWithRestarts(currentBest int) int {
	fmt.Printf("  Trying randomized search with target < %d...\n", currentBest)
	best := currentBest

	for restart := 0; restart < 3; restart++ { // Further reduced
		candidate := m.randomizedSearch(best, 200) // Further reduced
		if candidate < best {
			best = candidate
			fmt.Printf("  Randomized search found solution: %d\n", best)
		}
	}

	if best < currentBest {
		fmt.Printf("  Randomized search improved to %d\n", best)
	} else {
		fmt.Printf("  Randomized search found no improvement\n")
	}
	return best
}

func (m *machine) randomizedSearch(maxPresses int, maxIterations int) int {
	for iter := 0; iter < maxIterations; iter++ {
		m.reset()

		failed := false
		for len(m.presses) < maxPresses && !failed {
			switch m.check() {
			case Equal:
				return len(m.presses)
			case TooHigh:
				if len(m.presses) > 0 {
					m.unpressButton()
				} else {
					failed = true // Can't make progress, restart
				}
			case TooLow:
				candidates := m.getGoodCandidates()
				if len(candidates) == 0 {
					failed = true // No valid moves, restart
				} else {
					button := candidates[rand.IntN(len(candidates))]
					m.presses = append(m.presses, button)
					for _, counter := range m.buttons[button] {
						m.jstate[counter]++
					}
				}
			}
		}
	}

	return 1000000
}

func (m *machine) getGoodCandidates() []int {
	var candidates []int

	for i, button := range m.buttons {
		score := 0
		hasPositive := false

		for _, counter := range button {
			diff := m.joltages[counter] - m.jstate[counter]
			if diff > 0 {
				score += diff
				hasPositive = true
			} else if diff < 0 {
				score -= 5
			}
		}

		if hasPositive && score > 0 {
			weight := score
			if weight > 10 {
				weight = 10
			}
			for w := 0; w < weight; w++ {
				candidates = append(candidates, i)
			}
		}
	}

	return candidates
}

func (m *machine) unpressButton() {
	ix := rand.IntN(len(m.presses))
	button := m.presses[ix]
	m.presses = append(m.presses[:ix], m.presses[ix+1:]...)
	for _, i := range m.buttons[button] {
		m.jstate[i]--
	}
}

type joltstate int

const (
	TooHigh joltstate = 1
	TooLow  joltstate = -1
	Equal   joltstate = 0
)

func (m *machine) check() joltstate {
	isEqual := true
	for i := range m.joltages {
		if m.jstate[i] > m.joltages[i] {
			return TooHigh
		}
		if m.jstate[i] < m.joltages[i] {
			isEqual = false
		}
	}
	if isEqual {
		return Equal
	}
	return TooLow
}

func (m *machine) reset() {
	m.jstate = make([]int, len(m.joltages))
	m.presses = nil
}

func parseLamps(s string) (int, bits) {
	// string is . (0) and # (1), and low order bit is first
	// so ".#.#" should be A and "#..###" should be 56
	var result int
	for i, c := range s {
		b := 0
		if c == '#' {
			b = 1
		}
		result = result | (b << i)
	}
	return len(s), bits(result)
}

func parseNumberList(s string) []int {
	// takes a comma-separated list of numbers and returns a slice of ints
	parts := strings.Split(s, ",")
	var result []int
	for _, p := range parts {
		num, err := strconv.Atoi(p)
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, num)
	}
	return result
}

func parseSwitch(s string) bits {
	// takes a comma-separated list of positions and sets them to true in the result
	var result int
	values := parseNumberList(s)
	for _, pos := range values {
		result |= 1 << pos
	}
	return bits(result)
}

func parseButtons(s string) []int {
	// takes a comma-separated list of positions and sets them to true in the result
	values := parseNumberList(s)
	return values
}

func parseJolts(s string) []int {
	return parseNumberList(s)
}

var lampPat = regexp.MustCompile(`\[([^\]]+)\]`)
var switchPat = regexp.MustCompile(`\(([^\)]+)\)`)
var joltsPat = regexp.MustCompile(`\{([^\}]+)\}`)

func parseMachine(line string) machine {
	machine := machine{}
	ml := lampPat.FindStringSubmatch(line)
	if ml == nil {
		log.Fatalf("invalid format %s", line)
	}
	machine.nbits, machine.lamps = parseLamps(ml[1])

	ms := switchPat.FindAllStringSubmatch(line, -1)
	for _, sub := range ms {
		machine.switches = append(machine.switches, parseSwitch(sub[1]))
		machine.buttons = append(machine.buttons, parseButtons(sub[1]))
	}

	mj := joltsPat.FindStringSubmatch(line)
	if mj != nil {
		machine.joltages = parseJolts(mj[1])
	}
	machine.jstate = make([]int, len(machine.joltages))
	return machine
}

func parse(filename string) []machine {
	f, err := os.Open(fmt.Sprintf("./data/%s.txt", filename))
	if err != nil {
		log.Fatal(err)
	}
	b, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(b), "\n")

	var machines []machine
	for _, line := range lines {
		if line == "" {
			continue
		}
		machine := parseMachine(line)
		// machine.print()
		machines = append(machines, machine)
	}

	return machines
}

func part1(data []machine) int {
	sum := 0
	for _, m := range data {
		if d := m.search(1); d != -1 {
			sum += d
			fmt.Printf("solved in %d steps at %s: %s\n", d, time.Now().Format(time.RFC3339), m)
		} else {
			fmt.Printf("failed to solve machine at %s: %s\n", time.Now().Format(time.RFC3339), m)
		}
	}
	return sum
}

func part2(data []machine) int {
	total := 0
	failCount := 0
	for i, m := range data {
		fmt.Printf("Solving machine %d...\n", i+1)
		best := m.hybridOptimalSolve()
		if best == -1 {
			failCount++
		} else {
			fmt.Printf("Machine %d solved with %d presses\n", i+1, best)
			total += best
		}
	}
	if failCount > 0 {
		fmt.Printf("%d machines failed to solve optimally.\n", failCount)
	}
	return total
}

func main() {
	args := os.Args[1:]
	filename := "sample"
	if len(args) > 0 {
		switch args[0] {
		case "-s":
			filename = "sample"
		case "-i":
			filename = "input"
		default:
			filename = args[0]
		}
	}
	data := parse(filename)
	// fmt.Println(part1(data))
	fmt.Println(part2(data))
}
