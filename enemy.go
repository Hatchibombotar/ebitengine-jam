package main

import "math"

type Enemy struct {
	*Character
	move         func(g *Game, e *Enemy)
	Health       int
	MaxHealth    int
	Acceleration Vector
}

func createEnemy1() *Enemy {
	enemy := &Enemy{
		Health:    20,
		MaxHealth: 20,
		Character: &Character{
			position:        Vector{5, 5},
			startLerpT:      -1000,
			facingDirection: Vector{1, 0},
			walkSpeed:       .01,
			spriteIndex:     1,
			speedMultiplier: 1,
		},
	}

	return enemy
}
func (e *Enemy) Move(sublevel *Sublevel, targetX, targetY float64) {
	e.position.X += e.Acceleration.X
	e.position.Y += e.Acceleration.Y

	e.Acceleration.X *= 0.9
	e.Acceleration.Y *= 0.9
	speed := 2.0

	// Calculate direction to target
	dx := targetX - e.position.X
	dy := targetY - e.position.Y

	// Stop moving if very close to target
	if math.Abs(dx) < (1.0 / 16.0) {
		dx = 0
	}
	if math.Abs(dy) < (1.0 / 16.0) {
		dy = 0
	}

	// Move in primary direction only (cardinal directions)
	if math.Abs(dx) > math.Abs(dy) {
		e.position.X = e.position.X + math.Copysign(speed/32, dx)
		dy = 0 // Don't move vertically
	} else if dy != 0 {
		e.position.Y = e.position.Y + math.Copysign(speed/32, dy)
		dx = 0 // Don't move horizontally
	}

	// Update facing direction

	if dy > 0 {
		e.facingDirection.X = 0
		e.facingDirection.Y = 1
	} else if dy < 0 {
		e.facingDirection.X = 0
		e.facingDirection.Y = -1
	}

	if dx > 1 {
		e.facingDirection.X = 1
		e.facingDirection.Y = 0
	} else if dx < -1 {
		e.facingDirection.X = -1
		e.facingDirection.Y = 0
	}

	// Update movement state
	if math.Abs(dx)+math.Abs(dy) > 0 {
		e.isMoving = true
	} else {
		e.isMoving = false
	}
}

// 	type Node struct {
// 		x, y int
// 		cost float64
// 		// For priority queue
// 		index int
// 	}

// 	// PriorityQueue implements heap.Interface for Node pointers
// 	type PriorityQueue []*Node

// 	(func(pq PriorityQueue) Len)()
// 	return len(pq)
// }
// func (pq PriorityQueue) Less(i, j int) bool { return pq[i].cost < pq[j].cost }
// func (pq PriorityQueue) Swap(i, j int) {
// 	pq[i], pq[j] = pq[j], pq[i]
// 	pq[i].index = i
// 	pq[j].index = j
// }
// func (pq *PriorityQueue) Push(x interface{}) {
// 	n := x.(*Node)
// 	n.index = len(*pq)
// 	*pq = append(*pq, n)
// }
// func (pq *PriorityQueue) Pop() interface{} {
// 	old := *pq
// 	n := old[len(old)-1]
// 	*pq = old[0 : len(old)-1]
// 	return n
// }

// Move implements Dijkstra's algorithm to find a path to the target
// func (e *Enemy) Move(sublevel *Sublevel, targetX, targetY int) {
// 	tilemap := sublevel.tileMap

// 	// Current position
// 	currentX := int(e.position.X)
// 	currentY := int(e.position.Y)

// 	// Early exit: already at target
// 	if currentX == targetX && currentY == targetY {
// 		return
// 	}

// 	// Initialize Dijkstra's
// 	const (
// 		rows = 15
// 		cols = 20
// 	)

// 	// Track visited nodes and costs
// 	visited := make([][]bool, rows)
// 	costs := make([][]float64, rows)
// 	parent := make([][]struct{ x, y int }, rows)

// 	for i := range visited {
// 		visited[i] = make([]bool, cols)
// 		costs[i] = make([]float64, cols)
// 		parent[i] = make([]struct{ x, y int }, cols)
// 		for j := range costs[i] {
// 			costs[i][j] = math.Inf(1)
// 		}
// 	}

// 	// Priority queue
// 	pq := make(PriorityQueue, 0)
// 	heap.Push(&pq, &Node{x: currentX, y: currentY, cost: 0})
// 	costs[currentX][currentY] = 0

// 	// Dijkstra's algorithm
// 	for pq.Len() > 0 {
// 		current := heap.Pop(&pq).(*Node)

// 		if visited[current.x][current.y] {
// 			continue
// 		}

// 		visited[current.x][current.y] = true

// 		// Found the target
// 		if current.x == targetX && current.y == targetY {
// 			break
// 		}

// 		// Check all 4 neighbors (up, down, left, right)
// 		neighbors := [][2]int{
// 			{current.x - 1, current.y},
// 			{current.x + 1, current.y},
// 			{current.x, current.y - 1},
// 			{current.x, current.y + 1},
// 		}

// 		for _, neighbor := range neighbors {
// 			nx, ny := neighbor[0], neighbor[1]

// 			// Bounds check
// 			if nx < 0 || nx >= rows || ny < 0 || ny >= cols {
// 				continue
// 			}

// 			// Skip solid tiles
// 			if tilemap[nx][ny] != nil && TileIsSolid(tilemap[nx][ny]) {
// 				continue
// 			}

// 			// Skip visited nodes
// 			if visited[nx][ny] {
// 				continue
// 			}Copysign(0.1)

// 			// Calculate new cost (movement distance = 1)
// 			newCost := costs[current.x][current.y] + 1

// 			// Update if this is a better path
// 			if newCost < costs[nx][ny] {
// 				costs[nx][ny] = newCost
// 				parent[nx][ny] = struct{ x, y int }{x: current.x, y: current.y}
// 				heap.Push(&pq, &Node{x: nx, y: ny, cost: newCost})
// 			}
// 		}
// 	}

// 	// Reconstruct path from target back to current position
// 	path := [][2]int{}
// 	x, y := targetX, targetY
// 	for x != currentX || y != currentY {
// 		path = append(path, [2]int{x, y})
// 		p := parent[x][y]
// 		x, y = p.x, p.y
// 	}

// 	// Reverse to get path from current to target
// 	for i := len(path)/2 - 1; i >= 0; i-- {
// 		opp := len(path) - 1 - i
// 		path[i], path[opp] = path[opp], path[i]
// 	}

// 	// Move one step along the path
// 	if len(path) > 0 {
// 		nextPos := path[0]
// 		e.position.X = (e.position.X + float64(nextPos[0])) / 2
// 		e.position.Y = (e.position.Y + float64(nextPos[1])) / 2

// 		// Update facing direction based on movement
// 		if nextPos[0] < currentX {
// 			e.facingDirection.X = -1
// 			e.facingDirection.Y = 0
// 		} else if nextPos[0] > currentX {
// 			e.facingDirection.X = 1
// 			e.facingDirection.Y = 0
// 		} else if nextPos[1] < currentY {
// 			e.facingDirection.X = 0
// 			e.facingDirection.Y = -1
// 		} else {
// 			e.facingDirection.X = 0
// 			e.facingDirection.Y = 1
// 		}
// 	}
// }
