package recommendation

const maxConsecutiveSameValue = 2

func applyDiversity(scored []scoredCandidate, limit int) []scoredCandidate {
	selected := make([]scoredCandidate, 0, limit)
	deferred := make([]scoredCandidate, 0)

	for _, item := range scored {
		if len(selected) >= limit {
			break
		}

		if violatesDiversity(selected, item) {
			deferred = append(deferred, item)
			continue
		}

		selected = append(selected, item)
	}

	for len(selected) < limit && len(deferred) > 0 {
		index := firstDiverseCandidateIndex(selected, deferred)
		if index == -1 {
			index = 0
		}

		selected = append(selected, deferred[index])
		deferred = append(deferred[:index], deferred[index+1:]...)
	}

	return selected
}

func firstDiverseCandidateIndex(selected []scoredCandidate, candidates []scoredCandidate) int {
	for index, candidate := range candidates {
		if !violatesDiversity(selected, candidate) {
			return index
		}
	}

	return -1
}

func violatesDiversity(selected []scoredCandidate, candidate scoredCandidate) bool {
	return violatesAuthorDiversity(selected, candidate) || violatesCategoryDiversity(selected, candidate)
}

func violatesAuthorDiversity(selected []scoredCandidate, candidate scoredCandidate) bool {
	if len(selected) < maxConsecutiveSameValue {
		return false
	}

	for i := len(selected) - maxConsecutiveSameValue; i < len(selected); i++ {
		if selected[i].Candidate.Video.AuthorID != candidate.Candidate.Video.AuthorID {
			return false
		}
	}

	return true
}

func violatesCategoryDiversity(selected []scoredCandidate, candidate scoredCandidate) bool {
	if candidate.Candidate.Video.CategoryID == nil || len(selected) < maxConsecutiveSameValue {
		return false
	}

	for i := len(selected) - maxConsecutiveSameValue; i < len(selected); i++ {
		if selected[i].Candidate.Video.CategoryID == nil ||
			*selected[i].Candidate.Video.CategoryID != *candidate.Candidate.Video.CategoryID {
			return false
		}
	}

	return true
}
