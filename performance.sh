go clean -testcache ; go test -run "TestPerformance" -cpuprofile cpu.prof ; go tool pprof -web cpu.prof