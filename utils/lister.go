package utils

type ListenerType string
type SubscribeFunc[T any] func(ListenerType, func(*T))
