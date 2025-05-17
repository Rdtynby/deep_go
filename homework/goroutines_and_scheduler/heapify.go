package main

func heapifyStep(arr []Task, n, i int) {
	for {
		left := 2*i + 1
		right := 2*i + 2
		largest := i

		if left < n && arr[left].Priority > arr[largest].Priority {
			largest = left
		}

		if right < n && arr[right].Priority > arr[largest].Priority {
			largest = right
		}

		if largest == i {
			break
		}

		arr[i], arr[largest] = arr[largest], arr[i]
		i = largest
	}
}

func heapify(arr []Task) {
	n := len(arr)

	for i := n/2 - 1; i >= 0; i-- {
		heapifyStep(arr, n, i)
	}
}
