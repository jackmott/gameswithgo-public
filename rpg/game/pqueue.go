package game

type priorityPos struct {
	Pos
	priority int
}

type pqueue []priorityPos

func (pq pqueue) push(pos Pos, priority int) pqueue {
	newNode := priorityPos{pos, priority}
	pq = append(pq, newNode)
	newNodeIndex := len(pq) - 1
	parentIndex, parent := pq.parent(newNodeIndex)

	for newNode.priority < parent.priority && newNodeIndex != 0 {
		pq.swap(newNodeIndex, parentIndex)
		newNodeIndex = parentIndex
		parentIndex, parent = pq.parent(newNodeIndex)
	}
	return pq
}

func (pq pqueue) pop() (pqueue, Pos) {
	result := pq[0].Pos
	pq[0] = pq[len(pq)-1]
	pq = pq[:len(pq)-1]

	if len(pq) == 0 {
		return pq, result
	}

	index := 0
	node := pq[index]

	leftExists, leftIndex, left := pq.left(index)
	rightExists, rightIndex, right := pq.right(index)

	for (leftExists && node.priority > left.priority) ||
		(rightExists && node.priority > right.priority) {

		if !rightExists || left.priority <= right.priority {
			pq.swap(index, leftIndex)
			index = leftIndex
		} else {
			pq.swap(index, rightIndex)
			index = rightIndex
		}

		leftExists, leftIndex, left = pq.left(index)
		rightExists, rightIndex, right = pq.right(index)
	}

	return pq, result

}

func (pq pqueue) swap(i, j int) {
	tmp := pq[i]
	pq[i] = pq[j]
	pq[j] = tmp
}

func (pq pqueue) parent(i int) (int, priorityPos) {
	index := (i - 1) / 2
	return index, pq[index]
}

func (pq pqueue) left(i int) (bool, int, priorityPos) {
	index := i*2 + 1
	if index < len(pq) {
		return true, index, pq[index]
	}
	return false, 0, priorityPos{}
}

func (pq pqueue) right(i int) (bool, int, priorityPos) {
	index := i*2 + 2
	if index < len(pq) {
		return true, index, pq[index]
	}
	return false, 0, priorityPos{}
}
