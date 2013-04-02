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

	sdl.WM_SetCaption("Horde3d Go SDL Example", "")

	//set sdl video mode
	if sdl.SetVideoMode(width, height, 32, sdl.OPENGL) == nil {
		panic(sdl.GetError())
	}

	horde3d.Init()
	//horde3d.SetOption(horde3d.Options_DebugViewMode, 1)
	//horde3d.SetOption(horde3d.Options_WireframeMode, 1)
	fmt.Println("Version: ", horde3d.VersionString())

	//pipeline
	pipeRes := horde3d.AddResource(horde3d.ResTypes_Pipeline, "forward.pipeline.xml", 0)
	modelRes := horde3d.AddResource(horde3d.ResTypes_SceneGraph, "platform.scene.xml", 0)

	horde3d.LoadResourcesFromDisk("../content|" +
		"../content/pipelines|" +
		"../content/models|" +
		"../content/materials|" +
		"../content/shaders|" +
		"../content/textures|" +
		"../content/animations|" +
		"../content/particles|" +
		"../content/models/platform|" +
		"../content/effects")

	//add camera
	cam := horde3d.RootNode.AddCameraNode("Camera", pipeRes)
	//Setup Camera Viewport
	cam.SetNodeParamI(horde3d.Camera_ViewportXI, 0)
	cam.SetNodeParamI(horde3d.Camera_ViewportYI, 0)
	cam.SetNodeParamI(horde3d.Camera_ViewportWidthI, width)
	cam.SetNodeParamI(horde3d.Camera_ViewportHeightI, height)

	//add model
	model := horde3d.RootNode.AddNodes(modelRes)
	model.SetTransform(0, -30, -150, 0, 0, 0, 1, 1, 1)
	//add light
	light := horde3d.RootNode.AddLightNode("Light1", 0, "LIGHTING", "SHADOWMAP")
	light.SetTransform(0, 20, 0, 0, 0, 0, 1, 1, 1)
	light.SetNodeParamF(horde3d.Light_RadiusF, 0, 150)
	light.SetNodeParamF(horde3d.Light_FovF, 0, 90)
	//horde3d.SetNodeParamI(light, horde3d.Light_ShadowMapCountI, 3)
	light.SetNodeParamF(horde3d.Light_ShadowSplitLambdaF, 0, 0.9)
	//horde3d.SetNodeParamF(light, horde3d.Light_ShadowMapBiasF, 0, 0.001)
	light.SetNodeParamF(horde3d.Light_ColorF3, 0, 1.9)
	light.SetNodeParamF(horde3d.Light_ColorF3, 1, 1.7)
	light.SetNodeParamF(horde3d.Light_ColorF3, 2, 1.75)

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
		//horde3d.SetNodeTransform(model,
		//t*10, 0, 0,
		//0, 0, 0,
		//1, 1, 1)
		horde3d.Render(cam)
		horde3d.FinalizeFrame()
		horde3d.DumpMessages()
		sdl.GL_SwapBuffers()
	}
	horde3d.Release()
	sdl.Quit()

}
