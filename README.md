# pixelext - Game Engine Extension for pixel

**Important: this is under heavy development and anything might change at any time**

This library is an extension to the pixel 2D game engine (github.com/faiface/pixel) and adds the following features:

* Nodes
  * BaseNode
  * Canvas
  * Sprite
  * Sub Scene: render node tree in "sub window"
* Scene management
* GUI elements incl. styling
* Resource management
  * Pictures
  * Spritesheets
  * Fonts (TTF support is planned)
  * Sound samples (WAV) and music (MP3) (uses github.com/gpayer/go-audio-service/snd)
* Tile Map Editor support (https://www.mapeditor.org/)
  * basic tile maps
  * *Planned: offsets and other advanced tile map features*
  * *Planned: general object support*
* *Planned: integration of collision and physics framework*
* *Planned: event system for e.g. collisions*

## Installation

This is based on `faiface/pixel`, so all its requirements must be fullfilled first, see https://github.com/faiface/pixel#requirements for details.

> go get -u github.com/gpayer/pixelext

## Example

Go into the example directory and run
> go run pixelext.go