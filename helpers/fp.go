package helpers

func Filter[T any](data []T, callback func(T) bool) []T {
	filtered := make([]T, 0, len(data))

	for _, element := range data {
		if callback(element) {
			filtered = append(filtered, element)
		}
	}

	return filtered
}

func Find[T any](data []T, callback func(T) bool) *T {
	for _, element := range data {
		if callback(element) {
			return &element
		}
	}

	return nil
}

func Map[T, U any](data []T, callback func(T) U) []U {
	mapped := make([]U, 0, len(data))

	for _, element := range data {
		mapped = append(mapped, callback(element))
	}

	return mapped
}
