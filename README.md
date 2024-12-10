# SpacedAce
An E-learning platform utilizing flashcards, AI, and spaced repetition.

## Platform Access
You can access and try out the platform here: [SpacedAce](https://spacedace.hu/)

## Example contexts for quiz generation

### English

#### Solar system, short

```
The solar system is a fascinating collection of celestial bodies bound together by the Sun’s gravitational pull. At its center is the Sun, a massive star that provides light and heat to the system. Surrounding the Sun are eight planets, each following its own orbit. The inner planets, Mercury, Venus, Earth, and Mars, are rocky and smaller, while the outer planets, Jupiter, Saturn, Uranus, and Neptune, are much larger and primarily composed of gases. In addition to planets, the solar system contains dwarf planets, such as Pluto, countless moons, and small objects like asteroids and comets. These objects move through the vast expanse of space, which is also filled with interplanetary dust and charged particles from the solar wind. The outer boundary of the solar system is defined by the heliosphere, a bubble-like region influenced by the Sun’s magnetic field, stretching well beyond the orbit of Neptune.
```

#### Middle ages, long

```
The Middle Ages, or the medieval period, spanned roughly from the 5th to the late 15th century, following the fall of the Western Roman Empire. This era is often divided into three periods: the Early Middle Ages, the High Middle Ages, and the Late Middle Ages. It was a time of significant social, political, and cultural transformation in Europe. Feudalism became the dominant political and economic system, with lords, vassals, and serfs forming a hierarchical society. The Catholic Church held immense power and influence, shaping education, art, and governance across the continent. 

Castles and fortified cities were built to defend against invasions, while knights followed the code of chivalry in warfare and conduct. This period also saw the rise of powerful monarchies and the establishment of universities, which became centers of learning. Major events such as the Crusades and the Black Death had profound impacts on medieval society. The Crusades were religious wars aimed at reclaiming the Holy Land, while the Black Death devastated Europe’s population in the 14th century. Art and literature flourished, with Gothic architecture and epic tales like *The Divine Comedy* emerging during this time. 

The Middle Ages gradually gave way to the Renaissance, marking the beginning of a new era of exploration, science, and cultural revival. Despite being often characterized as a "dark age," the medieval period laid important foundations for the modern world.
```

### Hungarian

#### Naprendszer, short

```
A Naprendszer egy lenyűgöző égi testek gyűjteménye, amelyeket a Nap gravitációs vonzása tart össze. Középpontjában a Nap áll, egy hatalmas csillag, amely fényt és hőt biztosít a rendszer számára. A Nap körül nyolc bolygó kering, mindegyik saját pályáján haladva. A belső bolygók, mint a Merkúr, a Vénusz, a Föld és a Mars, kicsik és sziklásak, míg a külső bolygók, a Jupiter, a Szaturnusz, az Uránusz és a Neptunusz, sokkal nagyobbak és főként gázokból állnak. A bolygókon kívül a Naprendszerben törpebolygók, például a Plútó, számtalan hold, valamint kisebb objektumok, mint aszteroidák és üstökösök találhatók. Ezek az objektumok a tér hatalmas kiterjedésében mozognak, amelyet bolygóközi por és a napszél tölt meg. A Naprendszer külső határát a helioszféra jelöli ki, amely egy buborékszerű régió, amelyet a Nap mágneses tere befolyásol, és messze túlér a Neptunusz pályáján.
```

#### Középkor, long

```
A középkor, más néven a középkori időszak, nagyjából a 5. századtól a 15. század végéig tartott, a Nyugat-Római Birodalom bukása után. Ezt a korszakot gyakran három részre osztják: a korai középkorra, a virágzó középkorra és a késő középkorra. Ez az időszak jelentős társadalmi, politikai és kulturális átalakulásokkal járt Európában. A feudális rendszer vált uralkodó politikai és gazdasági rendszerré, amelyben urak, vazallusok és jobbágyok alkották a hierarchikus társadalmat. A katolikus egyház óriási hatalommal és befolyással bírt, amely az oktatást, a művészetet és a kormányzást is formálta.

Várakat és megerősített városokat építettek a támadások elleni védelem érdekében, miközben a lovagok a lovagiasság kódexét követték a harcban és viselkedésükben. Ebben az időszakban erős királyságok emelkedtek fel, és egyetemek jöttek létre, amelyek a tanulás központjaivá váltak. A középkor jelentős eseményei közé tartoznak a keresztes háborúk és a fekete halál, amelyek mély hatást gyakoroltak a társadalomra. A keresztes háborúk vallási háborúk voltak a Szentföld visszahódításáért, míg a fekete halál a 14. században Európa lakosságának nagy részét pusztította el. A művészet és az irodalom virágzott, ekkor jelentek meg például a gótikus építészet és olyan epikus művek, mint Az Isteni színjáték.

A középkor fokozatosan átadta helyét a reneszánsznak, amely az újkori felfedezések, tudományos fejlődés és kulturális megújulás korszakát hozta el. Bár a középkort gyakran "sötét korszaknak" nevezik, fontos alapokat fektetett le a modern világ számára.
```

## Repository Structure

This repository contains the source code for the components of the platform.

1. backend - Backend services's code
2. frontend - Frontend service's code
3. llm-api - The LLM integration's code
4. llm - Modelfile and initialization scripts
5. postgres - Docerfiles for the database
