Fire Emblem team randomizer, creating a team for you and assigning classes to the units.

Currently supports the following editions:
- FE11 (Shadow Dragon)
- FE12 (New Mystery of the Emblem)
- FE16 (Three Houses)

SETTINGS FILE:

You MUST add a file named "settings.txt" to the directory.

Each option must be listed in one line of a settings file.

The following settings are needed: (irrelevant settings can be left out)

"game: <FE11|FE12|FE16>"
Choose the game

"route: <route-abbreviation>"
FE16 only
Choose which route you will play on (in FE16, characters that aren't recruitable in that route are filtered out).
Accepts the following abbreviations:
- AM for Azure Moon
- CF for Crimson Flower
- VW for Verdant Wind
- SS for Silver Snow

"force_dancer: <yes|no>"
FE12 and FE16 only
Choose if the randomizer always gives you a dancer (random team member in FE16, Feena in FE12).
Picking yes is highly recommended.

"male_crossover: <yes|no>"
FE12 only
Choose if male units in FE12 can receive classes from the other class set (between A- and B-set classes).

Example:
Barst starts out with a B-set, so without crossover classes, he can only become a Warrior, General, Berserker, Horseman, Hero or Sorcerer.
With crossover classes, he can be randomized into the A-set classes aswell, such as Swordmaster, Paladin and Dracoknight.

"units: <number>"
total number of team members (excluding automatically free units like Nagi/Gotoh in FE11 and endgame bishops in FE12)

"same_class_limit: <number>"
maximum number of units that can receive the same class.

Example:
You choose "2" for this option and the randomizer already gave Dracoknight to 2 characters.
It then notices that the limit has been reached and makes sure all following units get a different class.

Exception:
If the limit has been reached for all eligible classes, a random class will be picked from the list regardless of the limit.
This can occur if you choose a high number of team members, paired with a low same-class limit.

"gaidens: <yes|no>"
FE11 only
Decides if gaiden characters (Athena, Horace, Etzel, Ymir) can be assigned to the team.
Only recommended if you use a patch that removes the gaiden requirements, or you'll have to sacrifice everyone else to get all team members.


EXAMPLE CONTENT FOR A SETTINGS FILE:

------
filename: settings.txt
------
game: FE16
route: AM
force_dancer: yes
male_crossover: yes
units: 15
same_class_limit: 2
gaidens: yes
------

When reading this file, the randomizer notices that the game is FE16, so the irrelevant settings "male_crossover" and "gaidens" are ignored.
It will throw an error if the 
They don't have to be removed from the file.


CLASS NAMES & IDEAS

The randomizer always gives puts out the names of the highest-tier class.

In FE16, you can choose your own rules.
I like the following ruleset: allow all classes with matching or lower weapon requirements. 

E. g. Mortal Savant has A Swords and B+ Reason
-> if a unit receives Mortal Savant, the unit can use all classes which...
- only require sword and/or reason rank
- require A or less in sword rank and 

In FE11 and FE12, the unpromoted classes correspond pretty exactly, aside from male dracoknight, female Paladin (FE11) and female General (FE12),
so the lower-tier classes are pretty obvious.

In FE12, I like to give female Generals both Archer and Cav as options before promotions, while male Dracoknights can choose between Fighter or Cav.
In FE11, Cavalier->M!Dracoknight and Pegasus Knight->F!Paladin are obvious.