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
	"code.google.com/p/gohorde/horde3d"
	"fmt"
	"github.com/jteeuwen/glfw"
	"os"
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
	benchmarkLength int = 600;
	t0 float64
	mx0, my0 int
	running bool

	app Application
)



func windowCloseListener() {
	running = false;
	return 0;
}


func keyPressListener(key int, action  int) {
	if !running {
		return
	}

	if action == glfw.GLFW_PRESS {
		width := appWidth
		height := appHeight;
		
		switch key {
		case GLFW_KEY_ESC:
			running = false;
		case GLFW_KEY_F1:
			app.release();
			glfw.CloseWindow()

			// Toggle fullscreen mode
			fullScreen = !fullScreen;

			if fullScreen  {
				mode := glfw.VidMode
				glfw.GetDesktopMode( &mode )
				
				aspect := float32(mode.Width) / (float32)mode.Height;
				if int(aspect * 100) == 133 || int(aspect * 100) == 125 {   // Standard
					width = 1280
					height = 1024;
				} else if int(aspect * 100) == 177 {                           // Widescreen 16:9
					width = 1280
					height = 720;
				} else if int(aspect * 100) == 160 {                           // Widescreen 16:10
					width = 1280
					height = 800
				} else {                                                            // Unknown
					// Use desktop resolution
					width = mode.Width; 
					height = mode.Height;
				}
			}
			
			if !setupWindow( width, height, fullScreen ) {
				glfw.Terminate();
				//exit( -1 );
			}
			
			app.init();
			app.resize( width, height );
			t0 = glfw.GetTime();
		}
	}
}


func mouseMoveListener( x int, y int ) {
	if !running { 
		mx0 = x; 
		my0 = y;
		return;
	}

	app.mouseMoveEvent( float32(x - mx0), float32(my0 - y) )
	mx0 = x; my0 = y;
}


func setupWindow(width int,height int , fullscreen bool ) bool {
{
	// Create OpenGL window
	var windowType int
	if fullScreen {
		windowType = glfw.GLFW_FULLSCREEN
	} else {
		windowType = glfw.GLFW_WINDOW
	}

	if !glfw.OpenWindow( width, height, 8, 8, 8, 8, 24, 8, windowType )  {
		glfw.Terminate();
		return false;
	}
	
	// Disable vertical synchronization
	glfw.SwapInterval( 0 );

	// Set listeners
	glfw.SetWindowCloseCallback( windowCloseListener );
	glfw.SetKeyCallback( keyPressListener );
	glfw.SetMousePosCallback( mouseMoveListener );
	
	return true;
}


func main()
{
	// Initialize GLFW
	glfw.Init();
	glfw.Enable( glfw.GLFW_STICKY_KEYS );
	if !setupWindow( appWidth, appHeight, fullScreen ) {
		return
	}

	// Initialize application and engine
	app = new Application();
	if !fullScreen {
		glfw.SetWindowTitle( app.getTitle() );
	}
	
	app.resize( appWidth, appHeight );

	glfw.Disable( GLFW_MOUSE_CURSOR );

	frames := 0;
	fps := 30.0f;
	t0 = glfw.GetTime();
	running = true;

	// Game loop
	for running {	
		// Calc FPS
		frames++;
		if !benchmark && frames >= 3 {
			double t = glfw.GetTime();
			fps = frames / float32(t - t0);
			if fps < 5 {
				fps = 30;  // Handle breakpoints
			}
			frames = 0;
			t0 = t;
		}

		// Update key states
		for  i := 0; i < 320; ++i {
			app.setKeyState( i, glfw.GetKey( i ) == glfw.GLFW_PRESS );
		}
		app.keyStateHandler();

		// Render
		app.mainLoop( fps );
		glfw.SwapBuffers();

	}

	glfw.Enable( glfw.GLFW_MOUSE_CURSOR );

		// Quit
	app.release();
	glfw.Terminate();

	return 
}



// Convert from degrees to radians
func degToRad( f float32 ) float32 {
	return f * (3.1415926 / 180.0);
}

type Application struct {
	keys map [int]bool
	prevKeys map [int]bool
	x,y,z, rx, ry, rz, velocity float32
	fps int  = 30
	animTime = 0
	weight = 1.0
	cam = 0
	contentDir "../content"
	hdrPipeRes, forwardPipeRes, fontMatRes, 
		panelMatRes, logoMatRes horde3d.H3DRes
	knight, particleSys horde3d.H3DNode
	animTime, weight, curFps float32


}

func (app *Application) init() bool {
	// Initialize engine
	if !horde3d.Init() {
		horde3d.DumpMessages();
		return false;
	}

	// Set options
	horde3d.SetOption( horde3d.Options_LoadTextures, 1 );
	horde3d.SetOption( horde3d.Options_TexCompression, 0 );
	horde3d.SetOption( horde3d.Options_FastAnimation, 0 );
	horde3d.SetOption( horde3d.Options_MaxAnisotropy, 4 );
	horde3d.SetOption( horde3d.Options_ShadowMapSize, 2048 );

	// Add resources
	// Pipelines
	app.hdrPipeRes = horde3d.AddResource( horde3d.ResTypes_Pipeline, "pipelines/hdr.pipeline.xml", 0 );
	app.forwardPipeRes = horde3d.AddResource( horde3d.ResTypes_Pipeline, "pipelines/forward.pipeline.xml", 0 );
	// Overlays
	app.fontMatRes = horde3d.AddResource( horde3d.ResTypes_Material, "overlays/font.material.xml", 0 );
	app.panelMatRes = horde3d.AddResource( horde3d.ResTypes_Material, "overlays/panel.material.xml", 0 );
	app.logoMatRes = horde3d.AddResource( horde3d.ResTypes_Material, "overlays/logo.material.xml", 0 );
	// Environment
	envRes := horde3d.AddResource( horde3d.ResTypes_SceneGraph, "models/sphere/sphere.scene.xml", 0 );
	// Knight
	knightRes := horde3d.AddResource( horde3d.ResTypes_SceneGraph, "models/knight/knight.scene.xml", 0 );
	knightAnim1Res := horde3d.AddResource( horde3d.ResTypes_Animation, "animations/knight_order.anim", 0 );
	knightAnim2Res := horde3d.AddResource( horde3d.ResTypes_Animation, "animations/knight_attack.anim", 0 );
	// Particle system
	particleSysRes := horde3d.AddResource( horde3d.ResTypes_SceneGraph, "particles/particleSys1/particleSys1.scene.xml", 0 );
	// Load resources
	horde3d.LoadResourcesFromDisk( app.contentDir );

	// Add scene nodes
	// Add camera
	app.cam = horde3d.AddCameraNode( horde3d.RootNode, "Camera", app.hdrPipeRes );
	//horde3d.SetNodeParamI( _cam, horde3d.Camera::OccCullingI, 1 );
	// Add environment
	env := horde3d.AddNodes( horde3d.RootNode, envRes );
	horde3d.SetNodeTransform( env, 0, -20, 0, 0, 0, 0, 20, 20, 20 );
	// Add knight
	app.knight = horde3d.AddNodes( horde3d.RootNode, knightRes );
	horde3d.SetNodeTransform( app.knight, 0, 0, 0, 0, 180, 0, 0.1, 0.1, 0.1 );
	horde3d.SetupModelAnimStage( app.knight, 0, knightAnim1Res, 0, "", false );
	horde3d.SetupModelAnimStage( app.knight, 1, knightAnim2Res, 0, "", false );
	// Attach particle system to hand joint
	horde3d.FindNodes( _knight, "Bip01_R_Hand", horde3d.NodeTypes_Joint );
	hand := horde3d.GetNodeFindResult( 0 );
	app.particleSys = horde3d.AddNodes( hand, particleSysRes );
	horde3d.SetNodeTransform( app.particleSys, 0, 40, 0, 90, 0, 0, 1, 1, 1 );

	// Add light source
	light := horde3d.AddLightNode( horde3d.RootNode, "Light1", 0, "LIGHTING", "SHADOWMAP" );
	horde3d.SetNodeTransform( light, 0, 15, 10, -60, 0, 0, 1, 1, 1 );
	horde3d.SetNodeParamF( light, horde3d.Light_RadiusF, 0, 30 );
	horde3d.SetNodeParamF( light, horde3d.Light_FovF, 0, 90 );
	horde3d.SetNodeParamI( light, horde3d.Light_ShadowMapCountI, 1 );
	horde3d.SetNodeParamF( light, horde3d.Light_ShadowMapBiasF, 0, 0.01 );
	horde3d.SetNodeParamF( light, horde3d.Light_ColorF3, 0, 1.0 );
	horde3d.SetNodeParamF( light, horde3d.Light_ColorF3, 1, 0.8 );
	horde3d.SetNodeParamF( light, horde3d.Light_ColorF3, 2, 0.7 );
	horde3d.SetNodeParamF( light, horde3d.Light_ColorMultiplierF, 0, 1.0 );

	// Customize post processing effects
	matRes := horde3d.FindResource( horde3d.ResTypes_Material, "pipelines/postHDR.material.xml" );
	horde3d.SetMaterialUniform( matRes, "hdrExposure", 2.5, 0, 0, 0 );
	horde3d.SetMaterialUniform( matRes, "hdrBrightThres", 0.5, 0, 0, 0 );
	horde3d.SetMaterialUniform( matRes, "hdrBrightOffset", 0.08, 0, 0, 0 );
	return true
}


func (app *Application) mainLoop( fps float32 ) {
	app.curFPS = fps;

		app.animTime += 1.0f / _curFPS;

		// Do animation blending
		horde3d.SetModelAnimParams( app.knight, 0, app.animTime * 24.0, app.weight );
		horde3d.SetModelAnimParams( app.knight, 1, app.animTime * 24.0, 1.0 - app.weight );

		// Animate particle systems (several emitters in a group node)
		cnt := horde3d.FindNodes( app.particleSys, "", horde3d.NodeTypes_Emitter );
		for i := 0; i < cnt; ++i {
			horde3d.AdvanceEmitterTime( horde3d.GetNodeFindResult( i ), 1.0 / app.curFPS );
		}
	
	// Set camera parameters
	horde3d.SetNodeTransform( app.cam, app.x, app.y, app.z, app.rx ,app.ry, 0, 1, 1, 1 );
	
	// Show stats
		// Show logo
		ww := float32(horde3d.GetNodeParamI( app.cam, horde3d.Camera_ViewportWidthI ) /
	                 float32(horde3d.GetNodeParamI( app.cam, horde3d.Camera_ViewportHeightI );
	var ovLogo []float32 = { ww-0.4, 0.8, 0, 1,  ww-0.4, 1, 0, 0,  ww, 1, 1, 0,  ww, 0.8, 1, 1 };
	horde3d.ShowOverlays( ovLogo, 4, 1.0, 1.0, 1.0, 1.0, app.logoMatRes, 0 );
	
	// Render scene
	horde3d.Render( app.cam );

	// Finish rendering of frame
	horde3d.FinalizeFrame();

	// Remove all overlays
	horde3d.ClearOverlays();

	// Write all messages to log file
	horde3d.DumpMessages();
}


func (app *Application) release() {
	// Release engine
	horde3d.Release();
}


func (app *Application) resize(  width int, height int ) {
	// Resize viewport
	horde3d.SetNodeParamI( app.cam, horde3d.Camera_ViewportXI, 0 );
	horde3d.SetNodeParamI( app.cam, horde3d.Camera_ViewportYI, 0 );
	horde3d.SetNodeParamI( app.cam, horde3d.Camera_ViewportWidthI, width );
	horde3d.SetNodeParamI( app.cam, horde3d.Camera_ViewportHeightI, height );
	
	// Set virtual camera parameters
	horde3d.SetupCameraView( app.cam, 45.0, float32(width) / height, 0.1, 1000.0 );
	horde3d.ResizePipelineBuffers( app.hdrPipeRes, width, height );
	horde3d.ResizePipelineBuffers( app.forwardPipeRes, width, height );
}


func (app *Application) keyStateHandler() {
	// ----------------
	// Key-press events
	// ----------------
	

	if( app.keys[260] && !app.prevKeys[260] )  // F3
	{
		if horde3d.GetNodeParamI( app.cam, horde3d.Camera_PipeResI ) == app.hdrPipeRes  {
			horde3d.SetNodeParamI( app.cam, horde3d.Camera_PipeResI, app.forwardPipeRes );
		} else {
			horde3d.SetNodeParamI( app.cam, hord3d.Camera_PipeResI, app.hdrPipeRes );
		}
	}
	
	// --------------
	// Key-down state
	// --------------
	
	curVel := app.velocity / app.curFPS;
	
	if app.keys[287] { curVel *= 5;	// LShift}
	
	if app.keys['W']  {
		// Move forward
		app.x -= math.Sin(degToRad( app.ry ) ) * math.Cos( -degToRad( app.rx ) ) * curVel;
		app.y -= math.Sin( -degToRad( app.rx ) ) * curVel;
		app.z -= math.Cosf( degToRad( app.ry ) ) * math.Cosf( -degToRad( app.rx ) ) * curVel;
	}
	if app.keys['S'] {
		// Move backward
		app.x += math.Sin( degToRad( app.ry ) ) * math.Cos( -degToRad( app.rx ) ) * curVel;
		app.y += math.Sin( -degToRad( app.rx ) ) * curVel;
		app.z += math.Cos( degToRad( app.ry ) ) * math.Cos( -degToRad( app.rx ) ) * curVel;
	}
	if app.keys['A'] {
		// Strafe left
		app.x += math.Sin( degToRad( app.ry - 90) ) * curVel;
		app.z += math.Cos( degToRad( app.ry - 90 ) ) * curVel;
	}
	if app.keys['D'] {
		// Strafe right
		app.x += math.Sin( degToRad( app.ry + 90 ) ) * curVel;
		app.z += math.Cos( degToRad( app.ry + 90 ) ) * curVel;
	}
	if app.keys['1'] {
		// Change blend weight
		app.weight += 2 / app.curFPS;
		if app.weight > 1 { app.weight = 1}
	}
	if app.keys['2'] {
		// Change blend weight
		app.weight -= 2 / app.curFPS;
		if app.weight < 0 { app.weight = 0 }
	}
}


func (app *Application) mouseMoveEvent(  dX float32, dY float32 ) {
	
	// Look left/right
	app.ry -= dX / 100 * 30;
	
	// Loop up/down but only in a limited range
	app.rx += dY / 100 * 30;
	if( app.rx > 90 ) app.rx = 90; 
	if( app.rx < -90 ) app.rx = -90;
}
