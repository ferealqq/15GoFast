# Viikko 6 

## Suorituskyvyn analysointi 

Mennellää viikolla käytin suurimman osan ajasta algoritmin prosessointi kyvyn parantamiseen.

Minulla oli selvästi mielessä mitä ohjelmassa oli optimoitavaa jotta se toimisi vielä tehokkaammin.

Kaksi optimisoinnin näkökulmaa jota lähdin tutkimaan:
- Helpompien taulujen ratkaisu tehokkuuden optimointi.
- Kompleksimpien taulujen ratkaisuiden tehokkuudn optimointi, jotta algoritmi pystyisi ratkaisemaan myös erittäin moni mutkaiset taulut.

Suoritus kyky testaukseen löysin sopivaksi työkaluksi [pprof](https://pkg.go.dev/runtime/pprof#pkg-overview). Työlalun avulla kirjoitin kaksi testi skenaariota yllä mainittujen näkökulmien tutkimiseen.

Monimutkaisten taulujen ratkaisu kyvyn optimointiin:
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
`board` muuttujan ratkaisemiseen tarvitaan 48 siirtoa. Ratkaisun löytämiseen huonosti optimoidulla IDASearch algoritmillä menee noin 10 sekunttia. 


Helpompien taulujen ratkaisu kyvyn optimointiin testi:
```golang
func TestPerformanceAverage(t *testing.T) {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
	}
	maxRuntimeMS := time.Duration(10600)
	boards := [][16]t_cell{
		{3, 4, 6, 5, 1, 7, 2, 14, 13, 15, 11, 8, 10, 0, 9, 12},
		{3, 6, 8, 7, 5, 0, 9, 2, 1, 4, 14, 15, 13, 10, 12, 11},
	}
	perfList := make(map[int][]time.Duration)
	for id, board := range boards {
		for i := 0; i < 60; i++ {
			n := time.Now()
			srh := NewSearch(&State{
				size:       BOARD_ROW_SIZE,
				board:      board,
				complexity: 0,
			})
			node, _ := srh.IDAStar(maxRuntimeMS)
			dur := time.Since(n)
			perfList[id] = append(perfList[id], dur)
			assert.True(t, node.state.debugIsSuccess())
		}
	}
	if *cpuprofile != "" {
		pprof.StopCPUProfile()
	}
	for i, v := range perfList {
		sum := time.Duration(0)
		for _, j := range v {
			sum += j
		}
		avg := sum / time.Duration(len(v))
		fmt.Printf("board index %d ", i)
		fmt.Printf("average %s \n", avg)
	}
}
```

Testaus suoritettiin siten, että testien ja `golang` kielen käyttämät `cachet` olivat puhtaita. Tämä saavutettiin komennolla `go clean --cache ; go clean -testcache`. 

Loppu tuloksia analysoin `pprof` työkalun luomalla `svg` diagrammilla. Testien diagrammit ovat tallennettuina repository:n [kansioon documentation/graphs](https://github.com/ferealqq/15GoFast/blob/main/documentation/graphs). 

### `hash` funktion deprikointi

[Commit](https://github.com/ferealqq/15GoFast/commit/acc9e682af0079d3de72921a6bda3f409d119b31)

Käytetty komento:
```terminal
go clean --cache ; go clean -testcache ; go test -run "TestPerformanceOne" -cpuprofile once.prof ; go tool pprof -web once.prof
```
Käytetty testi: `TestPerformanceOnce`

Kuva diagrammista ennen `hash`
![image](https://user-images.githubusercontent.com/22598325/206879795-c09a8c38-6488-4942-874a-444b90d48ac4.png)

Diagrammista näkee, että oikean ratkaisun löytämiseen menee noin 8 sekunttia, josta 1.21 sekunttia käytetään `hash` funktiossa. Joka ei tee ohjelman kannalta mitään erikoista. Funktiota käytettiin `board` lista muuttujien vertailuun koska jokaisen `board` muuttujassa olevan alkion erillinen vertailu olisi turhan raskasta. Tämän takia päätin tehdä funktion joka käyttää base64 enkoodausta esittämään `board` muuttujasta uniikkia avainta jonka avulla voin helposti vertailla `board` muuttujia. base64 enkoodaus on kuitenkin hyvin hidas tähän käyttö tarkoitukseen. Siitä syystä, että sitä kutsuttiin joka kerta kun algoritmi tarkasteli uutta polkua.

Keksin, että tehokkaampi tapa esittää `board` muuttuja avaimena olisi käyttää jo olemassa olevaa `code` funktiota joka luo `int` pohjaisen avaimen. 

Kuva diagrammista jossa käytetään `code` funktiota `hash` funktion sijaan.
![image](https://user-images.githubusercontent.com/22598325/206880069-c1b50169-3ea2-4ebf-8992-745813907151.png)

Diagrammista havaitaan, että `code` funktion suoritukseen menee vain `0.26` sekunttia kun taas `hash` funktion suorittamiseen meni `1.21` sekunttia. 

### `GetValidStates` tehokkuuden parannus

[Commit](https://github.com/ferealqq/15GoFast/commit/901882fa6c9c25075b50451932021cc17931e5c3)

Käytetty komento:
```terminal
go clean --cache ; go clean -testcache ; go test -run "TestPerformanceOne" -cpuprofile once.prof ; go tool pprof -web once.prof
```
Käytetty testi: `TestPerformanceOnce`

Yllä olevista diagrammeista (hash funktio deprikaatioon liittyvät diagrammit) havaitaan, että iso osa suoritusajasta menee `GetValidStates` funktion suorittamisesta; noin `1.5` sekunttia. 

Tiedostin, että ohjelmaan ei näillä näkymin tule muiden kuin 4x4 lautojen ratkaisevia implementaatiota. Joten suorituskyvyn parantamiseksi muutin kaikki `board` muuttujaa käsittelevät funktiota ottamaan vain 16:sta olion pituisia `board` muuttujia. Tämän päivityksen ansiosta suorituskyky parani noin viisitoista prosenttia. `board` muuttujan staattisella pituudella oli myös vaikutus muiden funktioiden suorituskykyyn.

Diagrammi päivityksen jälkeen:
![image](https://user-images.githubusercontent.com/22598325/206880447-3f2d4840-0b85-4cff-a54f-676e5e921c4c.png)



With hash commit   					=> acc9e682af0079d3de72921a6bda3f409d119b31\
without hash commit 				=> 1772cfb422683a1f2263822c9ba8c7d29cf43996\
GetValidStates improvements => 901882fa6c9c25075b50451932021cc17931e5c3\
code improvements 					=> b2819343a6871c6be710120c3798572fffae6a8c\
memoization 								=> 75e6634be9a880dcae586de56fa01da906736814\
without SearchState.staes   => 6cb4194badddf5f93fc281904c9b8b48cdad55b4\
without Move								=> 888be03c76dd06c3f02efac6aa9decdd505f6038

Vertaile without hash commitin graaffia GetValidStates improvementtiin

memoa => without SearchState.states

Static GetValidStates return ja getElementIndex ei käytä generikkejä. noin (-1s) parannus.


Memoization teki algoritmista hitaamman helpommpilla algoritmeillä mutta nopeamman vaikeamilla algoritmeilla.



Vertaa miten, viikko 5 suoriutui tietyistä benchmarkeista ja miten viikko 6 lopussa ohjelma suoriutui benchmarkeista. 



