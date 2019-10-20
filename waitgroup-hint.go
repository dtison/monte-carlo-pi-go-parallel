// https://downey.io/notes/dev/openmp-parallel-for-in-golang/

var wg sync.WaitGroup
wg.Add(n)
for i := 0; i < n; i++ {
	go func(i int) {
		defer wg.Done()
		(*list)[i].prefixData += sublists[(*list)[i].sublistHead].prefixData
	}(i)
}
wg.Wait()