GravSim
============

GravSim is an extremely basic game or tech demo that simulates N-body gravitational force.
It uses Ebiten, a 2D game engine for Go.


![screenshot](https://user-images.githubusercontent.com/1165651/151712567-ab8ac289-517d-4193-8415-e0e5ba5dbe56.png)

GravSim isn't a great game to play; that wasn't really my goal. It is very simplistic and boring.
I only wrote it to refresh my linear algebra and graphics coding skills. The code is messy, poorly organized,
and there's a lot of half implemented ideas and stubs here and there, as the goals of this project kept shifting.

But it is still a game.  It has some standard game features for this sort of retro shmup:

* Toroidal play field geometry: Ship, bullets, planets, and gravity wrap around screen edges (like Asteroids)
* Player ship with 3 lives and basic controls: thrust, rotation, emergency stop, and shoot
* Shoot planets, 3 per level, to increase your score
  * Bullets are short range only
  * Only 3 bullets may in play at once
  * Earn extra lives after reaching certain score thresholds
  * Achieve the high score (max of all playthroughs since game start)
  * Once all planets are destroyed, player reaches a new level with 3 new planets
* Collision detection:
  * If the ship collides with a planet, the ship is destroyed, costing one life
    * If lives remain after ship is destroyed, ship position is reset
    * If no lives remain, game is over and everything resets except high score
  * If two planets collide, they bounce off each other (elastic collisions)
  * If a ship's bullet collides with a planet, the planet is destroyed and the player earns score.

Now that I have written this I would be much more comfortable writing a real game...
