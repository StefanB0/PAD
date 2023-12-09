package service

type LoadBalancer struct {
	items []string
	shift int
}

func NewLoadBalancerService() *LoadBalancer {
	return &LoadBalancer{
		items: make([]string, 0),
		shift: 0,
	}
}

func (s *LoadBalancer) SetItems(items []string) {
	s.items = items
}

func (s *LoadBalancer) UpdateItems(items []string) {
	for _, item := range items {
		s.AddItem(item)
	}

	uniqueMap := make(map[string]bool)
	newArray := make([]string, 0)

	for _, item := range s.items {
		if _, value := uniqueMap[item]; !value {
			uniqueMap[item] = true
			newArray = append(newArray, item)
		}
	}

	s.items = newArray
}

func (s *LoadBalancer) AddItem(item string) {
	s.items = append(s.items, item)
}

func (s *LoadBalancer) RemoveItem(item string) {
	for i, v := range s.items {
		if v == item {
			s.items = append(s.items[:i], s.items[i+1:]...)
			return
		}
	}
}

func (s *LoadBalancer) GetItem() string {
	if len(s.items) == 0 {
		return ""
	}

	// first := s.items[0]
	// for i := 0; i < len(s.items)-1; i++ {
	// 	s.items[i] = s.items[i+1]
	// }
	// s.items[len(s.items)-1] = first

	first := s.items[s.shift]
	s.shift = (s.shift + 1) % len(s.items)

	return first
}
