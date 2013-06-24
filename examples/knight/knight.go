//Copyright (c) 2012 Tim Shannon
//
//Permission is hereby granted, free of charge, to any person obtaining a copy
//of this software and associated documentation files (the "Software"), to deal
//in the Software without restriction, including without limitation the rights
//to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//copies of the Software, and to permit persons to whom the Software is
//furnished to do so, subject to the following conditions:
//
//The above copyright notice and this permission notice shall be included in
//all copies or substantial portions of the Software.
//
//THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
//AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
//LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
//THE SOFTWARE.

package main

import (
	"bitbucket.org/tshannon/gohorde/horde3d"
	"fmt"
	"github.com/jteeuwen/glfw"
	"math"
)

const (
	caption   = "Knight - Horde3D Sample (Go Implementation)"
	appWidth  = 1024
	appHeight = 576
)

// Configuration
var (
	fullScreen bool = false
	t0         float64
	mx0, my0   int
	running    bool

	app *Application
)

func windowCloseListener() int {
	running = false
	return 0
}

func keyPressListener(key int, action int) {
	if !running {
		return
	}

	if action == glfw.KeyPress {
		width := appWidth
		height := appHeight

		switch key {
		case glfw.KeyEsc:
			running = false
		case glfw.KeyF1:
			app.release()
			glfw.CloseWindow()

			// Toggle fullscreen mode
			fullScreen = !fullScreen

			if fullScreen {
				mode := glfw.DesktopMode()

				aspect := float32(mode.W) / float32(mode.H)
				if int(aspect*100) == 133 || int(aspect*100) == 125 { // Standard
					width = 1280
					height = 1024
				} else if int(aspect*100) == 177 { // Widescreen 16:9
					width = 1280
					height = 720
				} else if int(aspect*100) == 160 { // Widescreen 16:10
					width = 1280
					height = 800
				} else { // Unknown
					// Use desktop resolution
					width = mode.W
					height = mode.H
				}
			}

			if !setupWindow(width, height, fullScreen) {
				glfw.Terminate()
				//exit( -1 );
			}

			app.init()
			app.resize(width, height)
			t0 = glfw.Time()
		}
	}
}

func mouseMoveListener(x int, y int) {
	if !running {
		mx0 = x
		my0 = y
		return
	}

	app.mouseMoveEvent(float32(x-mx0), float32(my0-y))
	mx0 = x
	my0 = y
}

func setupWindow(width int, height int, fullscreen bool) bool {
	// Create OpenGL window
	var windowType int
	if fullScreen {
		windowType = glfw.Fullscreen
	} else {
		windowType = glfw.Windowed
	}

	if err := glfw.OpenWindow(width, height, 8, 8, 8, 8, 24, 8, windowType); err != nil {
		glfw.Terminate()
		return false
	}

	// Disable vertical synchronization
	glfw.SetSwapInterval(0)

	// Set listeners
	glfw.SetWindowCloseCallback(windowCloseListener)
	glfw.SetKeyCallback(keyPressListener)
	glfw.SetMousePosCallback(mouseMoveListener)

	return true
}

func main() {
	// Initialize GLFW
	glfw.Init()
	glfw.Enable(glfw.StickyKeys)
	if !setupWindow(appWidth, appHeight, fullScreen) {
		return
	}

	// Initialize application and engine
	app = new(Application)
	if !app.init() {
		fmt.Println("Error starting Horde3d")
		horde3d.DumpMessages()
	}
	if !fullScreen {
		glfw.SetWindowTitle(app.title)
	}

	app.resize(appWidth, appHeight)

	glfw.Disable(glfw.MouseCursor)

	var frames float32 = 0
	var fps float32 = 30.0
	t0 = glfw.Time()
	running = true

	// Game loop
	for running {
		// Calc FPS
		frames++
		if frames >= 3 {
			t := glfw.Time()
			fps = frames / float32(t-t0)
			if fps < 5 {
				fps = 30 // Handle breakpoints
			}
			frames = 0
			t0 = t
		}

		// Update key states
		for i := 0; i < 320; i++ {
			app.setKeyState(i, glfw.Key(i) == glfw.KeyPress)
		}
		app.keyStateHandler()

		// Render
		app.mainLoop(fps)
		glfw.SwapBuffers()

	}

	glfw.Enable(glfw.MouseCursor)

	// Quit
	app.release()
	glfw.Terminate()

	return
}

// Convert from degrees to radians
func degToRad(f float32) float64 {
	return float64(f) * (3.1415926 / 180.0)
}

type Application struct {
	keys                          []bool
	prevKeys                      []bool
	x, y, z, rx, ry, rz, velocity float32
	contentDir                    string
	hdrPipeRes, forwardPipeRes, fontMatRes,
	panelMatRes, logoMatRes horde3d.H3DRes
	cam, knight, particleSys horde3d.H3DNode
	animTime, weight, curFps float32
	title                    string
}

func (app *Application) init() bool {
	app.title = "Horde3D Knight Sample - Go Implementation"
	app.contentDir = "../content"
	app.keys = make([]bool, 320)
	app.prevKeys = make([]bool, 320)

	app.x = 5
	app.y = 3
	app.z = 19
	app.rx = 7
	app.ry = 15
	app.velocity = 10.0
	app.curFps = 30

	app.animTime = 0
	app.weight = 1.0
	app.cam = 0

	// Initialize engine
	if !horde3d.Init() {
		horde3d.DumpMessages()
		return false
	}

	// Set options
	horde3d.SetOption(horde3d.Options_LoadTextures, 1)
	horde3d.SetOption(horde3d.Options_TexCompression, 0)
	horde3d.SetOption(horde3d.Options_FastAnimation, 0)
	horde3d.SetOption(horde3d.Options_MaxAnisotropy, 4)
	horde3d.SetOption(horde3d.Options_ShadowMapSize, 2048)
	//horde3d.SetOption(horde3d.Options_DebugViewMode, 1)

	// Add resources
	// Pipelines
	app.hdrPipeRes = horde3d.AddResource(horde3d.ResTypes_Pipeline, "pipelines/hdr.pipeline.xml", 0)
	app.forwardPipeRes = horde3d.AddResource(horde3d.ResTypes_Pipeline, "pipelines/forward.pipeline.xml", 0)
	// Overlays
	app.fontMatRes = horde3d.AddResource(horde3d.ResTypes_Material, "overlays/font.material.xml", 0)
	app.panelMatRes = horde3d.AddResource(horde3d.ResTypes_Material, "overlays/panel.material.xml", 0)
	app.logoMatRes = horde3d.AddResource(horde3d.ResTypes_Material, "overlays/logo.material.xml", 0)
	// Environment
	envRes := horde3d.AddResource(horde3d.ResTypes_SceneGraph, "models/sphere/sphere.scene.xml", 0)
	// Knight
	knightRes := horde3d.AddResource(horde3d.ResTypes_SceneGraph, "models/knight/knight.scene.xml", 0)
	knightAnim1Res := horde3d.AddResource(horde3d.ResTypes_Animation, "animations/knight_order.anim", 0)
	knightAnim2Res := horde3d.AddResource(horde3d.ResTypes_Animation, "animations/knight_attack.anim", 0)
	// Particle system
	particleSysRes := horde3d.AddResource(horde3d.ResTypes_SceneGraph,
		"particles/particleSys1/particleSys1.scene.xml", 0)
	// Load resources
	horde3d.LoadResourcesFromDisk(app.contentDir)

	// Add scene nodes
	// Add camera
	app.cam = horde3d.RootNode.AddCameraNode("Camera", app.hdrPipeRes)
	app.cam.SetNodeParamI(horde3d.Camera_OccCullingI, 0)
	// Add environment
	env := horde3d.RootNode.AddNodes(envRes)
	env.SetTransform(0, -20, 0, 0, 0, 0, 20, 20, 20)
	// Add knight
	app.knight = horde3d.RootNode.AddNodes(knightRes)
	app.knight.SetTransform(0, 0, 0, 0, 180, 0, 0.1, 0.1, 0.1)
	horde3d.SetupModelAnimStage(app.knight, 0, knightAnim1Res, 0, "", false)
	horde3d.SetupModelAnimStage(app.knight, 1, knightAnim2Res, 0, "", false)
	// Attach particle system to hand joint
	horde3d.FindNodes(app.knight, "Bip01_R_Hand", horde3d.NodeTypes_Joint)
	hand := horde3d.GetNodeFindResult(0)
	app.particleSys = hand.AddNodes(particleSysRes)
	app.particleSys.SetTransform(0, 40, 0, 90, 0, 0, 1, 1, 1)

	// Add light source
	light := horde3d.RootNode.AddLightNode("Light1", 0, "LIGHTING", "SHADOWMAP")
	light.SetTransform(0, 15, 10, -60, 0, 0, 1, 1, 1)
	light.SetNodeParamF(horde3d.Light_RadiusF, 0, 30)
	light.SetNodeParamF(horde3d.Light_FovF, 0, 90)
	light.SetNodeParamI(horde3d.Light_ShadowMapCountI, 1)
	light.SetNodeParamF(horde3d.Light_ShadowMapBiasF, 0, 0.01)
	light.SetNodeParamF(horde3d.Light_ColorF3, 0, 1.0)
	light.SetNodeParamF(horde3d.Light_ColorF3, 1, 0.8)
	light.SetNodeParamF(horde3d.Light_ColorF3, 2, 0.7)
	light.SetNodeParamF(horde3d.Light_ColorMultiplierF, 0, 1.0)

	// Customize post processing effects
	matRes := horde3d.FindResource(horde3d.ResTypes_Material, "pipelines/postHDR.material.xml")
	horde3d.SetMaterialUniform(matRes, "hdrExposure", 2.5, 0, 0, 0)
	horde3d.SetMaterialUniform(matRes, "hdrBrightThres", 0.5, 0, 0, 0)
	horde3d.SetMaterialUniform(matRes, "hdrBrightOffset", 0.08, 0, 0, 0)
	return true
}

func (app *Application) mainLoop(fps float32) {
	app.curFps = fps
	//fmt.Println(app.curFps)

	app.animTime += 1.0 / app.curFps

	// Do animation blending
	horde3d.SetModelAnimParams(app.knight, 0, app.animTime*24.0, app.weight)
	horde3d.SetModelAnimParams(app.knight, 1, app.animTime*24.0, 1.0-app.weight)

	// Animate particle systems (several emitters in a group node)
	cnt := horde3d.FindNodes(app.particleSys, "", horde3d.NodeTypes_Emitter)
	for i := 0; i < cnt; i++ {
		horde3d.UpdateEmitter(horde3d.GetNodeFindResult(i), 1.0/app.curFps)
	}

	// Set camera parameters
	app.cam.SetTransform(app.x, app.y, app.z, app.rx, app.ry, 0, 1, 1, 1)

	// Show stats
	// Show logo
	ww := float32(app.cam.NodeParamI(horde3d.Camera_ViewportWidthI)) /
		float32(app.cam.NodeParamI(horde3d.Camera_ViewportHeightI))
	ovLogo := []float32{ww - 0.4, 0.8, 0, 1, ww - 0.4, 1, 0, 0, ww, 1, 1, 0, ww, 0.8, 1, 1}
	horde3d.ShowOverlays(ovLogo, 4, 1.0, 1.0, 1.0, 1.0, app.logoMatRes, 0)
	// Render scene
	horde3d.Render(app.cam)

	// Finish rendering of frame
	horde3d.FinalizeFrame()

	// Remove all overlays
	horde3d.ClearOverlays()

	// Write all messages to log file
	//horde3d.DumpMessages()
}

func (app *Application) release() {
	// Release engine
	horde3d.Release()
}

func (app *Application) resize(width int, height int) {
	// Resize viewport
	app.cam.SetNodeParamI(horde3d.Camera_ViewportXI, 0)
	app.cam.SetNodeParamI(horde3d.Camera_ViewportYI, 0)
	app.cam.SetNodeParamI(horde3d.Camera_ViewportWidthI, width)
	app.cam.SetNodeParamI(horde3d.Camera_ViewportHeightI, height)

	// Set virtual camera parameters
	horde3d.SetupCameraView(app.cam, 45.0, float32(width)/float32(height), 0.1, 1000.0)
	horde3d.ResizePipelineBuffers(app.hdrPipeRes, width, height)
	horde3d.ResizePipelineBuffers(app.forwardPipeRes, width, height)
}

func (app *Application) keyStateHandler() {
	// ----------------
	// Key-press events
	// ----------------

	if app.keys[260] && !app.prevKeys[260] { // F3
		if horde3d.H3DRes(app.cam.NodeParamI(horde3d.Camera_PipeResI)) == app.hdrPipeRes {
			app.cam.SetNodeParamI(horde3d.Camera_PipeResI, int(app.forwardPipeRes))
		} else {
			app.cam.SetNodeParamI(horde3d.Camera_PipeResI, int(app.hdrPipeRes))
		}
	}

	// --------------
	// Key-down state
	// --------------

	curVel := float64(app.velocity / app.curFps)

	if app.keys[287] {
		curVel *= 5 // LShift
	}

	if app.keys['W'] {
		// Move forward
		app.x -= float32(math.Sin(degToRad(app.ry)) * math.Cos(-degToRad(app.rx)) * curVel)
		app.y -= float32(math.Sin(-degToRad(app.rx)) * float64(curVel))
		app.z -= float32(math.Cos(degToRad(app.ry)) * math.Cos(-degToRad(app.rx)) * curVel)
	}
	if app.keys['S'] {
		// Move backward
		app.x += float32(math.Sin(degToRad(app.ry)) * math.Cos(-degToRad(app.rx)) * curVel)
		app.y += float32(math.Sin(-degToRad(app.rx)) * float64(curVel))
		app.z += float32(math.Cos(degToRad(app.ry)) * math.Cos(-degToRad(app.rx)) * curVel)
	}
	if app.keys['A'] {
		// Strafe left
		app.x += float32(math.Sin(degToRad(app.ry-90)) * curVel)
		app.z += float32(math.Cos(degToRad(app.ry-90)) * curVel)
	}
	if app.keys['D'] {
		// Strafe right
		app.x += float32(math.Sin(degToRad(app.ry+90)) * curVel)
		app.z += float32(math.Cos(degToRad(app.ry+90)) * curVel)
	}
	if app.keys['1'] {
		// Change blend weight
		app.weight += 2 / app.curFps
		if app.weight > 1 {
			app.weight = 1
		}
	}
	if app.keys['2'] {
		// Change blend weight
		app.weight -= 2 / app.curFps
		if app.weight < 0 {
			app.weight = 0
		}
	}
}

func (app *Application) mouseMoveEvent(dX float32, dY float32) {

	// Look left/right
	app.ry -= dX / 100 * 30

	// Loop up/down but only in a limited range
	app.rx -= dY / 100 * 30
	if app.rx > 90 {
		app.rx = 90
	}
	if app.rx < -90 {
		app.rx = -90
	}
}

func (app *Application) setKeyState(key int, state bool) {
	app.prevKeys[key] = app.keys[key]
	app.keys[key] = state
}
