# Viikko 5 

## Haasteet

Viime viikolla algoritmi jäi tilanteeseen jossa luulin, että "invert distance" heuristiikka toimisi. Tällä viikolla iso osa ajasta meni invert distance heuristiikan kanssa taisteluun. Ongelma oli moni syinen 
mutta pääsäntöisesti invert distance heuristiikka ei toiminut huonon implementaation takia. Turhauduin jatkuvaan debuggaamiseen ja päädyin etsimään vaihtoehtoisia heuristiikka menetelmiä.

## Heuristiikan vaihto 

Huomattuain, että invert distance heuristiikka ei ole paras mahdollinen ja sen toteuttaminen oli vaikeampaa kun ennustin. Niin päätin vaihtaa heursitiikkaa "walking distance". "Walking distance" heursitiikka on yhdistelmä "invert distance" ja "manhattan distance" heuristiikkaa jonka isoin etu on se, että sillä saa aina paremman heuristiikan kun kummallakaan toisella vaihtoehdolla. Mutta "walking distance" heuristiikka vaatii paljon enemmän  prosessointia kun vertaisensa. Tähän löytyy kuitenkin ratkaisu. Implementoimalla yksin kertaisen "tietokannan" jonne laskennan osa-alueita (niin sanottuja "patterneja") tallennetaan voimme laskea "walking distance" heursitiikan laskenta-aikaa huomattavasti. Implementaationi perustuu tähän [artikkeliin](https://web.archive.org/web/20141224035932/http://juropollo.xe0.ru:80/stp_wd_translation_en.htm).

## IDAStar toteutus 

Ymmärtämättäni olin viimeviikolla toteuttanut toimivan IDAStar algoritmin mutta ongelmana oli silloin käyttämäni heuristiikka "invert distance" tarkemmin sanottuna minun implementaatio heuristiikasta.

## Seuraava viikko

Vaikkakin algoritmi on valmis, siinä on paljon optimoitavaa vielä. Joten seuraavalla viikolla keskityn algoritmin optimointiin ja algoritmin liittämistä osaksi käyttöliittymää.