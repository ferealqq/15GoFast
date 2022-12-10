# Viikko 6 

## Suorituskyvyn analysointi 

`hash` funktio vei 0-30ms suorittaa ja 30ms pelkästään lookup käyttöön on erittäin paljon. Tämän korjaaminen on tärkeää.

Tähän kuva molemmista, ja näytä siitä kuvasta, että tällä testikoodilla saatiin noin 1.5s säästö. data setti on kyllä vaan 1 mut silti.

```golang
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func TestPerformanceOne(t *testing.T) {
	flag.Parse()
	maxRuntimeMS := time.Duration(3000)
	board := [16]t_cell{2, 5, 8, 9, 4, 1, 14, 7, 6, 10, 3, 15, 13, 0, 11, 12}
	n := time.Now()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
	}
	srh := NewSearch(&State{
		size:       BOARD_ROW_SIZE,
		board:      board,
		complexity: 0,
	})
	node, _ := srh.IDAStar(maxRuntimeMS)
	if *cpuprofile != "" {
		pprof.StopCPUProfile()
	}
	dur := time.Since(n)
	fmt.Printf("since n %s \n", dur)
	fmt.Printf("complexity %d \n", node.state.complexity)
	assert.True(t, node.state.isSuccess())
}
```

With hash commit   					=> acc9e682af0079d3de72921a6bda3f409d119b31
without hash commit 				=> 1772cfb422683a1f2263822c9ba8c7d29cf43996
GetValidStates improvements => 901882fa6c9c25075b50451932021cc17931e5c3
code improvements 					=> b2819343a6871c6be710120c3798572fffae6a8c
memoization 								=> 75e6634be9a880dcae586de56fa01da906736814

Vertaile without hash commitin graaffia GetValidStates improvementtiin

memoa => without SearchState.states

Static GetValidStates return ja getElementIndex ei käytä generikkejä. noin (-1s) parannus.


Memoization teki algoritmista hitaamman helpommpilla algoritmeillä mutta nopeamman vaikeamilla algoritmeilla.



Vertaa miten, viikko 5 suoriutui tietyistä benchmarkeista ja miten viikko 6 lopussa ohjelma suoriutui benchmarkeista. 

