package ui2d

import (
	"bufio"
	"fmt"
	"github.com/jackmott/rpg/game"
	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	"image/png"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

const ItemSizeRatio = .033

type mouseState struct {
	leftButton  bool
	rightButton bool
	pos         game.Pos
}

func getMouseState() *mouseState {
	mouseX, mouseY, mouseButtonState := sdl.GetMouseState()
	leftButton := mouseButtonState & sdl.ButtonLMask()
	rightButton := mouseButtonState & sdl.ButtonRMask()
	var result mouseState
	result.pos = game.Pos{int(mouseX), int(mouseY)}
	result.leftButton = !(leftButton == 0)
	result.rightButton = !(rightButton == 0)

	return &result
}

type sounds struct {
	openingDoors []*mix.Chunk
	footsteps    []*mix.Chunk
}

func playRandomSound(chunks []*mix.Chunk, volume int) {
	chunkIndex := rand.Intn(len(chunks))
	chunks[chunkIndex].Volume(volume)
	chunks[chunkIndex].Play(-1, 0)
}

type uiState int

const (
	UIMain uiState = iota
	UIInventory
)

type ui struct {
	state             uiState
	draggedItem       *game.Item
	sounds            sounds
	winWidth          int
	winHeight         int
	renderer          *sdl.Renderer
	window            *sdl.Window
	textureAtlas      *sdl.Texture
	textureIndex      map[rune][]sdl.Rect
	prevKeyboardState []uint8
	keyboardState     []uint8
	centerX           int
	centerY           int
	r                 *rand.Rand
	levelChan         chan *game.Level
	inputChan         chan *game.Input
	fontSmall         *ttf.Font
	fontMedium        *ttf.Font
	fontLarge         *ttf.Font

	eventBackground           *sdl.Texture
	groundInventoryBackground *sdl.Texture
	slotBackground            *sdl.Texture

	str2TexSmall  map[string]*sdl.Texture
	str2TexMedium map[string]*sdl.Texture
	str2TexLarge  map[string]*sdl.Texture

	currentMouseState *mouseState
	prevMouseState    *mouseState
}

func NewUI(inputChan chan *game.Input, levelChan chan *game.Level) *ui {

	ui := &ui{}
	ui.state = UIMain
	ui.str2TexSmall = make(map[string]*sdl.Texture)
	ui.str2TexMedium = make(map[string]*sdl.Texture)
	ui.str2TexLarge = make(map[string]*sdl.Texture)
	ui.inputChan = inputChan
	ui.levelChan = levelChan
	ui.r = rand.New(rand.NewSource(1))
	ui.winHeight = 720
	ui.winWidth = 1280
	window, err := sdl.CreateWindow("RPG!!", 200, 200,
		int32(ui.winWidth), int32(ui.winHeight), sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	ui.window = window

	ui.renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}

	//sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1")

	ui.textureAtlas = ui.imgFileToTexture("ui2d/assets/tiles.png")
	ui.loadTextureIndex()

	ui.keyboardState = sdl.GetKeyboardState()
	ui.prevKeyboardState = make([]uint8, len(ui.keyboardState))
	for i, v := range ui.keyboardState {
		ui.prevKeyboardState[i] = v
	}

	ui.centerX = -1
	ui.centerY = -1

	ui.fontSmall, err = ttf.OpenFont("ui2d/assets/gothic.ttf", int(float64(ui.winWidth)*.015))
	if err != nil {
		panic(err)
	}

	ui.fontMedium, err = ttf.OpenFont("ui2d/assets/gothic.ttf", 32)
	if err != nil {
		panic(err)
	}

	ui.fontLarge, err = ttf.OpenFont("ui2d/assets/gothic.ttf", 64)
	if err != nil {
		panic(err)
	}

	ui.eventBackground = ui.GetSinglePixelTex(sdl.Color{0, 0, 0, 128})
	ui.eventBackground.SetBlendMode(sdl.BLENDMODE_BLEND)

	ui.groundInventoryBackground = ui.GetSinglePixelTex(sdl.Color{149, 84, 19, 200})
	ui.groundInventoryBackground.SetBlendMode(sdl.BLENDMODE_BLEND)

	ui.slotBackground = ui.GetSinglePixelTex(sdl.Color{0, 0, 0, 0})

	//if( Mix_OpenAudio( 22050, MIX_DEFAULT_FORMAT, 2, 4096 ) == -1 )
	err = mix.OpenAudio(22050, mix.DEFAULT_FORMAT, 2, 4096)
	if err != nil {
		panic(err)
	}
	mus, err := mix.LoadMUS("ui2d/assets/ambient.ogg")
	if err != nil {
		panic(err)
	}
	mus.Play(-1)

	footstepBase := "ui2d/assets/footstep0"
	for i := 0; i < 10; i++ {
		footstepFile := footstepBase + strconv.Itoa(i) + ".ogg"
		footstepSound, err := mix.LoadWAV(footstepFile)
		if err != nil {
			panic(err)
		}
		ui.sounds.footsteps = append(ui.sounds.footsteps, footstepSound)
	}

	doorOpen1, err := mix.LoadWAV("ui2d/assets/doorOpen_1.ogg")
	if err != nil {
		panic(err)
	}
	ui.sounds.openingDoors = append(ui.sounds.openingDoors, doorOpen1)
	doorOpen2, err := mix.LoadWAV("ui2d/assets/doorOpen_2.ogg")
	if err != nil {
		panic(err)
	}
	ui.sounds.openingDoors = append(ui.sounds.openingDoors, doorOpen2)

	return ui
}

type FontSize int

const (
	FontSmall FontSize = iota
	FontMedium
	FontLarge
)

func (ui *ui) stringToTexture(s string, color sdl.Color, size FontSize) *sdl.Texture {

	var font *ttf.Font
	switch size {
	case FontSmall:
		font = ui.fontSmall
		tex, exists := ui.str2TexSmall[s]
		if exists {
			return tex
		}
	case FontMedium:
		font = ui.fontMedium
		tex, exists := ui.str2TexMedium[s]
		if exists {
			return tex
		}
	case FontLarge:
		font = ui.fontLarge
		tex, exists := ui.str2TexLarge[s]
		if exists {
			return tex
		}
	}

	fontSurface, err := font.RenderUTF8Blended(s, color)
	if err != nil {
		panic(err)
	}

	tex, err := ui.renderer.CreateTextureFromSurface(fontSurface)
	if err != nil {
		panic(err)
	}

	switch size {
	case FontSmall:
		ui.str2TexSmall[s] = tex
	case FontMedium:
		ui.str2TexMedium[s] = tex
	case FontLarge:
		ui.str2TexLarge[s] = tex
	}

	return tex
}

func (ui *ui) loadTextureIndex() {
	ui.textureIndex = make(map[rune][]sdl.Rect)
	infile, err := os.Open("ui2d/assets/atlas-index.txt")
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(infile)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		tileRune := rune(line[0])
		xy := line[1:]
		splitXYC := strings.Split(xy, ",")
		x, err := strconv.ParseInt(strings.TrimSpace(splitXYC[0]), 10, 64)
		if err != nil {
			panic(err)
		}
		y, err := strconv.ParseInt(strings.TrimSpace(splitXYC[1]), 10, 64)
		if err != nil {
			panic(err)
		}

		variationCount, err := strconv.ParseInt(strings.TrimSpace(splitXYC[2]), 10, 64)
		if err != nil {
			panic(err)
		}

		var rects []sdl.Rect
		for i := int64(0); i < variationCount; i++ {
			rects = append(rects, sdl.Rect{int32(x * 32), int32(y * 32), 32, 32})
			x++
			if x > 62 {
				x = 0
				y++
			}
		}
		fmt.Println("rectLen:", len(rects))
		ui.textureIndex[tileRune] = rects
	}

}
func (ui *ui) imgFileToTexture(filename string) *sdl.Texture {
	infile, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer infile.Close()

	img, err := png.Decode(infile)
	if err != nil {
		panic(err)
	}

	w := img.Bounds().Max.X
	h := img.Bounds().Max.Y

	pixels := make([]byte, w*h*4)
	bIndex := 0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			pixels[bIndex] = byte(r / 256)
			bIndex++
			pixels[bIndex] = byte(g / 256)
			bIndex++
			pixels[bIndex] = byte(b / 256)
			bIndex++
			pixels[bIndex] = byte(a / 256)
			bIndex++
		}
	}

	tex, err := ui.renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STATIC, int32(w), int32(h))
	if err != nil {
		panic(err)
	}
	tex.Update(nil, pixels, w*4)

	err = tex.SetBlendMode(sdl.BLENDMODE_BLEND)
	if err != nil {
		panic(err)
	}
	return tex
}

func init() {

	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		panic(err)
	}

	err = ttf.Init()
	if err != nil {
		panic(err)
	}

	err = mix.Init(mix.INIT_OGG)
	//SDL Bug here, ignoring error
	/*if err != nil {
		panic(err)
	}*/
}

func (ui *ui) Draw(level *game.Level) {

	if ui.centerX == -1 && ui.centerY == -1 {
		ui.centerX = level.Player.X
		ui.centerY = level.Player.Y
	}

	limit := 5
	if level.Player.X > ui.centerX+limit {
		diff := level.Player.X - (ui.centerX + limit)
		ui.centerX += diff
	} else if level.Player.X < ui.centerX-limit {
		diff := (ui.centerX - limit) - level.Player.X
		ui.centerX -= diff
	} else if level.Player.Y > ui.centerY+limit {
		diff := level.Player.Y - (ui.centerY + limit)
		ui.centerY += diff
	} else if level.Player.Y < ui.centerY-limit {
		diff := (ui.centerY - limit) - level.Player.Y
		ui.centerY -= diff
	}

	offsetX := int32((ui.winWidth / 2) - ui.centerX*32)
	offsetY := int32((ui.winHeight / 2) - ui.centerY*32)

	ui.renderer.Clear()
	ui.r.Seed(1)
	// Render Map Tiles
	for y, row := range level.Map {
		for x, tile := range row {
			if tile.Rune != game.Blank {
				srcRects := ui.textureIndex[tile.Rune]
				srcRect := srcRects[ui.r.Intn(len(srcRects))]
				if tile.Visible || tile.Seen {
					dstRect := sdl.Rect{int32(x*32) + offsetX, int32(y*32) + offsetY, 32, 32}
					pos := game.Pos{x, y}
					if level.Debug[pos] {
						ui.textureAtlas.SetColorMod(128, 0, 0)
					} else if tile.Seen && !tile.Visible {
						ui.textureAtlas.SetColorMod(128, 128, 128)
					} else {
						ui.textureAtlas.SetColorMod(255, 255, 255)
					}
					ui.renderer.Copy(ui.textureAtlas, &srcRect, &dstRect)

					if tile.OverlayRune != game.Blank {
						// Todo what if there are multiple variants for overlay images?
						srcRect := ui.textureIndex[tile.OverlayRune][0]
						ui.renderer.Copy(ui.textureAtlas, &srcRect, &dstRect)
					}

				}
			}

		}
	}
	ui.textureAtlas.SetColorMod(255, 255, 255)

	// Render Monsters
	for pos, monster := range level.Monsters {
		if level.Map[pos.Y][pos.X].Visible {
			monsterSrcRect := ui.textureIndex[monster.Rune][0]
			ui.renderer.Copy(ui.textureAtlas, &monsterSrcRect, &sdl.Rect{int32(pos.X)*32 + offsetX, int32(pos.Y)*32 + offsetY, 32, 32})
		}
	}

	// Render Items
	for pos, items := range level.Items {
		if level.Map[pos.Y][pos.X].Visible {
			for _, item := range items {
				itemSrcRect := ui.textureIndex[item.Rune][0]
				ui.renderer.Copy(ui.textureAtlas, &itemSrcRect, &sdl.Rect{int32(pos.X)*32 + offsetX, int32(pos.Y)*32 + offsetY, 32, 32})
			}
		}
	}

	// Render Player
	playerSrcRect := ui.textureIndex[level.Player.Rune][0]
	ui.renderer.Copy(ui.textureAtlas, &playerSrcRect, &sdl.Rect{int32(level.Player.X)*32 + offsetX, int32(level.Player.Y)*32 + offsetY, 32, 32})

	// Event UI Begin
	textStart := int32(float64(ui.winHeight) * .68)
	textWidth := int32(float64(ui.winWidth) * .25)

	ui.renderer.Copy(ui.eventBackground, nil, &sdl.Rect{0, textStart, textWidth, int32(ui.winHeight) - textStart})
	i := level.EventPos
	count := 0
	_, fontSizeY, _ := ui.fontSmall.SizeUTF8("A")
	for {
		event := level.Events[i]
		if event != "" {
			tex := ui.stringToTexture(event, sdl.Color{255, 0, 0, 0}, FontSmall)
			_, _, w, h, _ := tex.Query()
			ui.renderer.Copy(tex, nil, &sdl.Rect{5, int32(count*fontSizeY) + textStart, w, h})
		}
		i = (i + 1) % (len(level.Events))
		count++
		if i == level.EventPos {
			break
		}
	}
	// Event UI End

	// Inventory UI
	groundInvStart := int32(float64(ui.winWidth) * .9)
	groundInvWidth := int32(ui.winWidth) - groundInvStart
	itemSize := int32(ItemSizeRatio * float32(ui.winWidth))
	ui.renderer.Copy(ui.groundInventoryBackground, nil, &sdl.Rect{groundInvStart, int32(ui.winHeight) - itemSize, groundInvWidth, itemSize})
	items := level.Items[level.Player.Pos]

	for i, item := range items {
		itemSrcRect := ui.textureIndex[item.Rune][0]
		ui.renderer.Copy(ui.textureAtlas, &itemSrcRect, ui.getGroundItemRect(i))
	}
	// Inventory UI END

}

func (ui *ui) getGroundItemRect(i int) *sdl.Rect {
	itemSize := int32(ItemSizeRatio * float32(ui.winWidth))
	return &sdl.Rect{int32(ui.winWidth) - itemSize - int32(i)*itemSize, int32(ui.winHeight) - itemSize, itemSize, itemSize}
}

func (ui *ui) keyDownOnce(key uint8) bool {
	return ui.keyboardState[key] == 1 && ui.prevKeyboardState[key] == 0
}

// Check for key pressed then released
func (ui *ui) keyPressed(key uint8) bool {
	return ui.keyboardState[key] == 0 && ui.prevKeyboardState[key] == 1
}

func (ui *ui) GetSinglePixelTex(color sdl.Color) *sdl.Texture {
	tex, err := ui.renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STATIC, 1, 1)
	if err != nil {
		panic(err)
	}
	pixels := make([]byte, 4)
	pixels[0] = color.R
	pixels[1] = color.G
	pixels[2] = color.B
	pixels[3] = color.A
	tex.Update(nil, pixels, 4)
	return tex
}

func (ui *ui) Run() {
	var newLevel *game.Level
	ui.prevMouseState = getMouseState()

	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				ui.inputChan <- &game.Input{Typ: game.QuitGame}
			case *sdl.WindowEvent:
				if e.Event == sdl.WINDOWEVENT_CLOSE {
					ui.inputChan <- &game.Input{Typ: game.CloseWindow, LevelChannel: ui.levelChan}
				}
			}
		}
		ui.currentMouseState = getMouseState()

		// Suspect quick keypresses sometimes cause channel gridlock
		var ok bool
		select {
		case newLevel, ok = <-ui.levelChan:
			if ok {
				switch newLevel.LastEvent {
				case game.Move:
					playRandomSound(ui.sounds.footsteps, 10)
				case game.DoorOpen:
					playRandomSound(ui.sounds.openingDoors, 32)
				default:
					//add more sounds
				}

			}
		default:
		}

		ui.Draw(newLevel)
		var input game.Input
		if ui.state == UIInventory {

			//have we stopped dragging?
			if ui.draggedItem != nil && !ui.currentMouseState.leftButton && ui.prevMouseState.leftButton {

				item := ui.CheckEquippedItem()
				if item != nil {
					input.Typ = game.EquipItem
					input.Item = item
					ui.draggedItem = nil
				}
				if ui.draggedItem != nil {
					item := ui.CheckDroppedItem()
					if item != nil {
						input.Typ = game.DropItem
						input.Item = item
						ui.draggedItem = nil
					}
				}

			}
			if !ui.currentMouseState.leftButton || ui.draggedItem == nil {
				ui.draggedItem = ui.CheckInventoryItems(newLevel)
			}
			ui.DrawInventory(newLevel)
		}
		ui.renderer.Present()

		item := ui.CheckGroundItems(newLevel)
		if item != nil {
			input.Typ = game.TakeItem
			input.Item = item
		}
		if sdl.GetKeyboardFocus() == ui.window || sdl.GetMouseFocus() == ui.window {

			if ui.keyDownOnce(sdl.SCANCODE_UP) {
				input.Typ = game.Up
			} else if ui.keyDownOnce(sdl.SCANCODE_DOWN) {
				input.Typ = game.Down
			} else if ui.keyDownOnce(sdl.SCANCODE_LEFT) {
				input.Typ = game.Left
			} else if ui.keyDownOnce(sdl.SCANCODE_RIGHT) {
				input.Typ = game.Right
			} else if ui.keyDownOnce(sdl.SCANCODE_T) {
				input.Typ = game.TakeAll
			} else if ui.keyDownOnce(sdl.SCANCODE_I) {
				fmt.Println("I")
				if ui.state == UIMain {
					ui.state = UIInventory
				} else {
					ui.state = UIMain
				}
			}

			for i, v := range ui.keyboardState {
				ui.prevKeyboardState[i] = v
			}

			if input.Typ != game.None {
				ui.inputChan <- &input
			}
		}
		ui.prevMouseState = ui.currentMouseState
		sdl.Delay(10)

	}

}
