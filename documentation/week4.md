# Viikko 3 & 4 

## Tutkimustiedon kerääminen

Aloitin viikon 3 sillä, että keräsin tutkimus materiaalin pelin algoritmiselle ratkaisulle. Valitettavasti en löytänyt pelistä paljoa hyvää tutkimus materiaalia joka olisi auttanut minua implementoimaan ratkaisu algoritmin. Joten jouduin soveltamaan eri paikoista löytämääni tietoa jotta pystyin aloittamaan algoritmin toteuttamisen. 

Löytämiäni lähteitä: 
https://michael.kim/blog/puzzle
https://web.archive.org/web/20141224035932/http://juropollo.xe0.ru:80/stp_wd_translation_en.htm
http://kevingong.com/Math/SixteenPuzzle.html
https://mediatum.ub.tum.de/doc/1283911/1283911.pdf



## Ratkaisu algoritmin toteuttaminen

Ratkaisu algoritminä käytän sovellettua [IDA*:ta](https://en.wikipedia.org/wiki/Iterative_deepening_A*). Heurestiikkaan käytän ["invert distance"](https://michael.kim/blog/puzzle) laskentaa. Heuristiikka ja algoritmi ei ole optimaalinen tapa löytää ratkaisu 15Puzzle peliin mutta ottaen huomioon projektin aikataulun on valitsemani heuristiikka ja algoritmi parhaat mahdolliset. Tavoitteena on, että ratkaisu implementaatio pystyy ratkaisemaan noin 80% kaikista mahdollisisat 15puzzle peleistä alle viiden sekunnin. Mikäli haluaisin, että algoritmi pystyisi parempaan tulokseen tulisi minun implementoida "pattern database" ja "walking distance". 


Projektissa on tällä hetkellä toimiva "invert distance" heuristiikka mutta IDAStar algoritmissa on vielä tehtävää. 

## Haasteet

Haasteita implementaatiossa alkoi näkyä kun olin ehtinyt tutkia pelin ratkaisu mahdollisuuksia. Hyvin nopeasti kävi selväksi, että kaikki 15Puzzle pelit eivät olet ratkaistavissa, tämä taas tuotti ylimääräisen päänsäryn kun minun piti 

## Dokumentaatio

Projektissa on automaattinen dokumentaatio. Koodin dokumentaatiota voi käydä katselemassa [README:ssa](https://github.com/ferealqq/15GoFast/blob/main/README.md) olevien ohjeiden mukaan.


## Seuraava viikko

Tavoitteena on, että seuraavan viikon aikana IDAStar hakemis algoritmi on toimiva.