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

Tiedostin, että ohjelmaan ei näillä näkymin tule muiden kuin 4x4 lautojen ratkaisevia implementaatiota. Joten suorituskyvyn parantamiseksi muutin kaikki `board` muuttujaa käsittelevät funktiota ottamaan vain 16:sta olion pituisia `board` muuttujia. Tämän päivityksen ansiosta suorituskyky parani noin viisitoista prosenttia. `board` muuttujan pituuden muuttaminen staattiseksi vaikutti positiivisesti myös muiden funktioiden suoritusaikaan.

Diagrammi päivityksen jälkeen:
![image](https://user-images.githubusercontent.com/22598325/206880447-3f2d4840-0b85-4cff-a54f-676e5e921c4c.png)

### Lisää staattisia muuttujia

[Commit](https://github.com/ferealqq/15GoFast/commit/b2819343a6871c6be710120c3798572fffae6a8c)

Käytetty komento:
```terminal
go clean --cache ; go clean -testcache ; go test -run "TestPerformanceOne" -cpuprofile once.prof ; go tool pprof -web once.prof
```
Käytetty testi: `TestPerformanceOnce`

`code` funktiossa laskettiin joka suorituksella `BOARD_ROW_SIZE` muuttujan bittien määrä sekä `board` muuttujan pituus. Muutin äsken mainitut dynaamiset arvot staattisiksi arvoiksi siten, että ne lasketaan vain kerran kun ohjelma ensimmäistä kertaa konstruktoidaan. Päivitys ei vaikuttanut merkittävästi sovelluksen suoritusaikaan näin pienellä suoritus ajalla. Päivitys paransi suoritusaikaa n 3 prosenttia.


![image](https://user-images.githubusercontent.com/22598325/206880952-8284fc61-6610-49b5-9b2b-f2325243f640.png)


### `WalkingDistance.Calculate` funktion memoizointi

[Commit](https://github.com/ferealqq/15GoFast/commit/75e6634be9a880dcae586de56fa01da906736814)

Käytetty komento:
```terminal
go clean --cache ; go clean -testcache ; go test -run "TestPerformanceOne" -cpuprofile once.prof ; go tool pprof -web once.prof
```
Käytetty testi: `TestPerformanceOnce`

On ilmiselvää, että IDA* algoritmi kutsuu heuristiikka funktiota useamman kerran yhden. Diagrammeista nähdään, että `Calculate` funktion suorittaminen on noin 30 prosenttia koko IDA* algoritmin suoritusajasta. [Memoization](https://en.wikipedia.org/wiki/Memoization) tekniikka on juurikin tälläisiin tarkoituksiin erittäin hyödyllinen. Käytännössä `SearchState` struktuurin funktio `Heuristic` tarkistaa onko se laskenut vielä tietylle `board` muuttujalle heuristiikkaa jos se on laskenut palauttaa funktio arvon suoraan muistista mikäli arvoa ei ole vielä laskettu. 

Päivitys paransi algoritmin suoritusaikaa lähes 50 prosenttia. (Vertaa aikaisempaan diagrammiin, 6s/3s) 

![image](https://user-images.githubusercontent.com/22598325/206881215-0d948f7b-28de-4d0e-b511-b63dd72c6bfb.png)


### `SearchState.states` muuttujan deprikointi 

[Commit](https://github.com/ferealqq/15GoFast/commit/6cb4194badddf5f93fc281904c9b8b48cdad55b4)

Käytetty komento:
```terminal
go clean --cache ; go clean -testcache ; go test -run "TestPerformanceAverage" -cpuprofile cpu.prof ; go tool pprof -web cpu.prof
```
Käytetty testi: `TestPerformanceAverage`


En ollut vieläkään tyytyväinen `IDASearch` funktion suoritusaikaan. Diagrammeja tarkastellesse minulle heräsi ajatus, että vaikka `SearchState.states` on oleellinen muuttuja ohjleman toiminnan kannalta se ei pidä sisällään mitään uniikkia tietoa koska kaikki `State` pointterit ovat jo `SearchState.hasSeen` muuttujan sisällä vaikkakin ei oikeassa järjestykssä. `hasSeen` arvot voidaan kuitenkin myöhemmin järjestää oikeaan järjestykseen `State.complexity` olion avulla. 

Kuvassa nähdään kuinka kauan keskiarvoltaan (N = 60) kahden eri `board` muuttujan ratkaisemiseen meni kun `SearchState.states` oli käytössä.

![image](https://user-images.githubusercontent.com/22598325/206881420-588ee18d-1f60-4c83-be0d-c108f2befb9b.png)


Ilman `SearchState.states` muuttujaa.

![image](https://user-images.githubusercontent.com/22598325/206881500-89a3a5bf-1520-4105-ae48-acb74c08169b.png)

Kuvankaappauksista näkyy, helpommissa ratkaisuissa suorituskyky parani n 50 prosenttia.

### `State.move` muuttujan deprikointi

[Commit](https://github.com/ferealqq/15GoFast/commit/888be03c76dd06c3f02efac6aa9decdd505f6038)

Käytetty komento:
```terminal
go clean --cache ; go clean -testcache ; go test -run "TestPerformanceOne" -cpuprofile once.prof ; go tool pprof -web once.prof
```
Käytetty testi: `TestPerformanceOnce`

Ennen.
![image](https://user-images.githubusercontent.com/22598325/206881577-9a584e8c-0bad-402e-b72f-971778ff10d2.png)

Diagram! Minun lempilapsi ärsyttää jälleen. `GetValidStates` suoritukseen menee noin 13 prosenttia IDASearch funktion suoritusajasta. Ratkaisun löyty jälleen turhan muuttujan heivaamisesta ojaan. Tällä kertaa listalta löytyi `State` struktuurin muuttujassa `move` elänyt struktuuri `Move`. `Move` struktuuri piti sisällään tiedon mihin suuntaan `boardia` oltiin siirtämässä. Tämä tieto on kuitenkin irrelevantti koska `walking distance` heuristiikka pitää kuitenkin laskea koko boardista. `Move` muuttuja oli jäänne niiltä ajoilta kun IDAStar käytti `invert distance` heuristiikkaa.

Simsala bim! `Move` on poistettu ja suorituskyky parani noin 14 prosenttia. (2.6s/3s)

![image](https://user-images.githubusercontent.com/22598325/206881749-e88503bb-4d68-4cb4-be6e-9030a03d29e9.png)

Diagrammista nähdään, että `GetValidState` on enään 8.30 prosenttia `IDASearch` funktion suoritusajasta (ennen 13.6 prosenttia)


### Recap

Viiko aloitetiin siitä, että testillä `TestPerformanceOnce` `IDASearch` funktion suoritusaika oli 8 sekunttia ja lopetettiin siihen, että suoritusaika on 3 sekunttia. Suoritusaika parani noin 270 prosenttia.


## Mitä viikolla opin?

Viikolla tuli aika paljon opittua `golang` kielen sisäisitä toiminnallisuuksista ja suoritusajan optimoinnista.


## Seuraavaks

- Siivotaan koodi
- Lisätään dokumentointia
- Hiotaan käyttöliittymästä komia
- Mennään demoon esittämään siisti softa.


Kiitos ja anteeks.
  


