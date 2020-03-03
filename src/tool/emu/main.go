package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	_ "image/png"
	"io"
	"os"
	"pkg/shaders"
	s "pkg/shared"
	scrn "pkg/screen"
	vm "pkg/vm"
	"time"

	"github.com/faiface/glhf"
	"github.com/faiface/mainthread"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

// WIDTH of LCD screen
const WIDTH = 64

// HEIGHT of LCD screen
const HEIGHT = 32

// SCALE of LCD screen
const SCALE = 6

var _VM = new(vm.VirtualMachine)

func main() {
	var opcodes []uint16

	romFile := flag.String("rom", "", "the name of the rom to load and emulate")
	flag.Parse()

	fmt.Println(romFile)
	if *romFile == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	fp, err := os.Open(*romFile)

	if err != nil {
		fmt.Printf("Could not find the '%s' rom file\n", *romFile)
		os.Exit(1)
	}

	defer fp.Close()

	fi, err := fp.Stat()
	if err != nil {
		fmt.Printf("Problem stating rom file.\n")
		os.Exit(1)
	}

	romLength := fi.Size()
	romBuffer := make([]byte, romLength)

	// read the ROM contents
	_, err = fp.Read(romBuffer)

	if err == io.EOF {
		fmt.Printf("Problem readiing rom file.\n")
		os.Exit(1)
	}

	for i := 0; i < len(romBuffer)-1; i += 2 {
		wVal := uint16(binary.BigEndian.Uint16(romBuffer[i:]))
		opcodes = append(opcodes, wVal)
	}

	_VM.Init(opcodes)

	// move this "run" function into another file in the namespace
	mainthread.Run(run)
}

func glfwKeyToKeyCode(key glfw.Key) uint8 {
	gotKey := s.KEYCODE_UNKNOWN
	switch key {
	case glfw.KeyX:
		gotKey = s.KEYCODE_0
	case glfw.Key1:
		gotKey = s.KEYCODE_1
	case glfw.KeyA:
		gotKey = s.KEYCODE_2
	case glfw.Key3:
		gotKey = s.KEYCODE_3
	case glfw.Key2:
		gotKey = s.KEYCODE_4
	case glfw.KeyQ:
		gotKey = s.KEYCODE_5
	case glfw.KeyE:
		gotKey = s.KEYCODE_6
	case glfw.KeyW:
		gotKey = s.KEYCODE_7
	case glfw.KeyS:
		gotKey = s.KEYCODE_8
	case glfw.KeyD:
		gotKey = s.KEYCODE_9
	case glfw.KeyZ:
		gotKey = s.KEYCODE_A
	case glfw.KeyC: //
		gotKey = s.KEYCODE_B
	case glfw.Key4:
		gotKey = s.KEYCODE_C
	case glfw.KeyR:
		gotKey = s.KEYCODE_D
	case glfw.KeyF:
		gotKey = s.KEYCODE_E
	case glfw.KeyV:
		gotKey = s.KEYCODE_F
	}

	return gotKey
}

// The handler for keystroke handling
// Input is done with a hex keyboard that has 16 keys which range from 0 to F.
// The '8', '4', '6', and '2' keys are typically used for directional input.
func keyHandlerCallback(_ *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {

	keyCode := glfwKeyToKeyCode(key)
	if keyCode == s.KEYCODE_UNKNOWN {
		return
	}

	if action == s.GLFW_PRESS {
		_VM.KeyPress(keyCode)
	}

	if action == s.GLFW_RELEASE {
		_VM.KeyRelease(keyCode)
	}
}

func run() {

	// pointer to the GL window struct
	var win *glfw.Window

	// initialize component modules
	screen := new(scrn.Screen)
	screen.Init()
	_VM.AttachScreen(screen)

	// cpu etc

	// window close
	defer func() {
		mainthread.Call(func() {
			glfw.Terminate()
		})
	}()

	// Initialize View
	mainthread.Call(func() {
		glfw.Init()

		glfw.WindowHint(glfw.ContextVersionMajor, 3)
		glfw.WindowHint(glfw.ContextVersionMinor, 3)
		glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
		glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
		glfw.WindowHint(glfw.Resizable, glfw.False)
		glfw.WindowHint(glfw.Samples, 1)

		var err error
		win, err = glfw.CreateWindow(WIDTH*SCALE, HEIGHT*SCALE, "CHIP-8", nil, nil)
		if err != nil {
			panic(err)
		}

		win.MakeContextCurrent()

		glhf.Init()
	})

	// Register a keyboard handler
	win.SetKeyCallback(keyHandlerCallback)

	// set up the shader info
	var (
		vertexFormat = glhf.AttrFormat{
			{Name: "position", Type: glhf.Vec2},
			{Name: "texture", Type: glhf.Vec2},
		}

		fragmentFormat = glhf.AttrFormat{
			{Name: "resolution", Type: glhf.Vec2},
			{Name: "time", Type: glhf.Float},
		}

		shader  *glhf.Shader
		texture *glhf.Texture
		slice   *glhf.VertexSlice
	)

	// Every OpenGL call needs to be done inside the main thread.
	mainthread.Call(func() {
		var err error

		// Create the shader
		shader, err = glhf.NewShader(vertexFormat, fragmentFormat, shaders.VertexShader, shaders.FragmentShader)

		if err != nil {
			panic(err)
		}

		// Now create a texture from these pixels
		// TODO: Move this texture into the Screen
		// https://github.com/faiface/glhf/blob/master/texture.go#L20
		smooth := false 
		texture = glhf.NewTexture(
			screen.Width,
			screen.Height,
			smooth,             // do not smooth
			screen.GetPixels(), //[]uint8
		)

		// And finally, we make a vertex slice, which is basically a dynamically sized
		// vertex array. The length of the slice is 6 and the capacity is the same.
		// The slice inherits the vertex format of the supplied shader. Also, it should
		// only be used with that shader.
		slice = glhf.MakeVertexSlice(shader, 6, 6)

		// Before we use a slice, we need to Begin it. The same holds for all objects in
		// GLHF.
		slice.Begin()

		// We assign data to the vertex slice. The values are in the order as in the vertex
		// format of the slice (shader). Each two floats correspond to an attribute of type
		// glhf.Vec2.
		slice.SetVertexData([]float32{
			-1, -1, 0, 1,
			+1, -1, 1, 1,
			+1, +1, 1, 0,

			-1, -1, 0, 1,
			+1, +1, 1, 0,
			-1, +1, 0, 0,
		})

		// When we're done with the slice, we End it.
		slice.End()
	})

	/********************************************************
	 * VM Loop
	 ********************************************************/

	shouldQuit := false
	thisTime := time.Now()
	lastTime := thisTime
	timeSinceRefresh := float32(0.0)
	
	// view render loop
	for !shouldQuit {

		lastPixelBuffer := screen.GetPixels()

		// do CPU emulation here
		_VM.FetchDecodeExecute()

		// update delatime
		thisTime := time.Now()
		deltaTime := float32(thisTime.Sub(lastTime)) * 0.0000000001
		lastTime = thisTime

		// Refresh screen interval
		timeSinceRefresh += deltaTime

		if timeSinceRefresh >= 0.05 {
			timeSinceRefresh = 0
			screen.RefreshPixelBytes();
		}

		// should sound be playing?
		// cpu.SoundTimer > 0
		// speaker.playTone();

		// Render
		// Texture: https://github.com/faiface/glhf/blob/master/texture.go#L83
		mainthread.CallNonBlock(func() {
			if win.ShouldClose() {
				shouldQuit = true
			}
			glhf.Clear(1, 1, 1, 1)
			shader.Begin()
			res := mgl32.Vec2{64 * (SCALE*2), 32 * (SCALE * 2)} 
			shader.SetUniformAttr(0, res) 			
			shader.SetUniformAttr(1, float32(deltaTime)) // {Name: "time", Type: glhf.Float},
			texture.Begin()
			texture.SetPixels(0, 0, texture.Width(), texture.Height(), lastPixelBuffer)
			slice.Begin()
			slice.Draw()
			slice.End()
			texture.End()
			shader.End()
			win.SwapBuffers()
			glfw.PollEvents()
		})
	}
}