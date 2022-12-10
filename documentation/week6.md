# Viikko 6 

## Suorituskyvyn analysointi 

`hash` funktio vei 0-30ms suorittaa ja 30ms pelkästään lookup käyttöön on erittäin paljon. Tämän korjaaminen on tärkeää.

Tähän kuva molemmista, ja näytä siitä kuvasta, että tällä testikoodilla saatiin noin 1.5s säästö. data setti on kyllä vaan 1 mut silti.

```golang
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
			f, err := os.Create(*cpuprofile)
			if err != nil {
					log.Fatal(err)
			}
			pprof.StartCPUProfile(f)
			defer pprof.StopCPUProfile()
	}
	maxRuntimeMS := time.Duration(3500)
	state := &State{
		size:       BOARD_ROW_SIZE,
		board:			[16]t_cell{3, 4, 6, 5, 1, 7, 2, 14, 13, 15, 11, 8, 10, 0, 9, 12},
		complexity: 0,
	}
	srh := NewSearch(state)
	node, _ := srh.IDAStar(maxRuntimeMS)
	if !node.state.isSuccess() {
		panic("algo broke");
	}

	maxRuntimeMS = time.Duration(7600)
	state = &State{
		size:       BOARD_ROW_SIZE,
		board:      [16]t_cell{2, 5, 8, 9, 4, 1, 14, 7, 6, 10, 3, 15, 13, 0, 11, 12},
		complexity: 0,
	}

	srh = NewSearch(state)
	node, _ = srh.IDAStar(maxRuntimeMS)
	if !node.state.isSuccess() {
		panic("algo broke");
	}
}
```

Static GetValidStates return ja getElementIndex ei käytä generikkejä. noin (-1s) parannus.


Vertaa miten, viikko 5 suoriutui tietyistä benchmarkeista ja miten viikko 6 lopussa ohjelma suoriutui benchmarkeista.

