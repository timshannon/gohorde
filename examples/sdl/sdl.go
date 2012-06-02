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

/* 
SDL Example from http://horde3d.org/wiki/index.php5?title=Tutorial_-_Setup_Horde_with_SDL
implemented in Go
*/

import (
	"code.google.com/p/gohorde/horde3d"
	"fmt"
	"github.com/banthar/Go-SDL/sdl"
)

func main() {
	var running bool = true
	var width int = 800
	var height int = 600
	var t float32

	if sdl.Init(sdl.INIT_VIDEO) != 0 {
		panic(sdl.GetError())
	}

	sdl.WM_SetCaption("Simple SDL Frame", "")

	//set sdl video mode
	if sdl.SetVideoMode(width, height, 32, sdl.OPENGL) == nil {
		panic(sdl.GetError())
	}

	horde3d.Init()
	fmt.Println("Version: ", horde3d.GetVersionString())

	//pipeline
	pipeRes := horde3d.AddResource(horde3d.H3DResTypes_Pipeline, "standard.pipeline.xml", 0)
	modelRes := horde3d.AddResource(horde3d.H3DResTypes_SceneGraph, "character.scene.xml", 0)
	animRes := horde3d.AddResource(horde3d.H3DResTypes_Animation, "walk.anim.xml", 0)

	horde3d.LoadResourcesFromDisk("")

	//add camera
	cam := horde3d.AddCameraNode(horde3d.RootNode, "Camera", pipeRes)
	//Setup Camera Viewport
	horde3d.SetNodeParamI(cam, horde3d.H3DCamera_ViewportXI, 0)
	horde3d.SetNodeParamI(cam, horde3d.H3DCamera_ViewportYI, 0)
	horde3d.SetNodeParamI(cam, horde3d.H3DCamera_ViewportWidthI, width)
	horde3d.SetNodeParamI(cam, horde3d.H3DCamera_ViewportHeightI, height)

	//add model
	model := horde3d.AddNodes(horde3d.RootNode, modelRes)
	horde3d.SetupModelAnimStage(model, 0, animRes, 0, "", false)
	//add light
	light := horde3d.AddLightNode(horde3d.RootNode, "Light1", 0, "LIGHTING", "SHADOWMAP")
	horde3d.SetNodeTransform(light, 0, 20, 0, 0, 0, 0, 1, 1, 1)
	horde3d.SetNodeParamF(light, horde3d.H3DLight_RadiusF, 0, 50)
	horde3d.SetNodeParamF(light, horde3d.H3DLight_FovF, 0, 90)
	horde3d.SetNodeParamI(light, horde3d.H3DLight_ShadowMapCountI, 3)
	horde3d.SetNodeParamF(light, horde3d.H3DLight_ShadowSplitLambdaF, 0, 0.9)
	horde3d.SetNodeParamF(light, horde3d.H3DLight_ShadowMapBiasF, 0, 0.001)
	horde3d.SetNodeParamF(light, horde3d.H3DLight_ColorF3, 0, 0.9)
	horde3d.SetNodeParamF(light, horde3d.H3DLight_ColorF3, 1, 0.7)
	horde3d.SetNodeParamF(light, horde3d.H3DLight_ColorF3, 2, 0.75)

	for running {
		t = 0
		//increase anim time
		t = t + 10.0*(1/60)
		//process SDL events / input
		switch event := sdl.PollEvent(); event.(type) {
		case *sdl.QuitEvent:
			running = false
			break
		}
		horde3d.SetModelAnimParams(model, 0, t, 1.0)
		horde3d.SetNodeTransform(model,
			t*10, 0, 0,
			0, 0, 0,
			1, 1, 1)
		horde3d.Render(cam)
		sdl.GL_SwapBuffers()
	}
	horde3d.Release()
	sdl.Quit()

}
