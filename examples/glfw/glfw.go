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
)

const (
	caption   = "GLFW Sample (Go Implementation)"
	appWidth  = 800
	appHeight = 600
)

var pipeRes horde3d.H3DRes

func main() {
	var running bool = true

	if err := glfw.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "[e] %v\n", err)
		return
	}

	//ensure glfw is cleaned up
	defer glfw.Terminate()

	if err := glfw.OpenWindow(appWidth, appHeight, 8, 8, 8, 8,
		24, 8, glfw.Windowed); err != nil {
		fmt.Fprintf(os.Stderr, "[e] %v\n", err)
		return
	}
	defer glfw.CloseWindow()

	glfw.SetSwapInterval(1)
	glfw.SetWindowTitle(caption)
	if !horde3d.Init() {
		fmt.Println("Error starting Horde3D. Check Horde3d_log.html for details.")
		horde3d.DumpMessages()
		return
	}

	horde3d.SetOption(horde3d.Options_DebugViewMode, 1)
	// Add resources
	//pipeline
	pipeRes = horde3d.AddResource(horde3d.ResTypes_Pipeline, "hdr.pipeline.xml", 0)

	knightRes := horde3d.AddResource(horde3d.ResTypes_SceneGraph, "knight/knight.scene.xml", 0)

	//load resources paths separated by |
	horde3d.LoadResourcesFromDisk("../content|" +
		"../content/pipelines|" +
		"../content/models|" +
		"../content/materials|" +
		"../content/shaders|" +
		"../content/textures|" +
		"../content/effects")

	model := horde3d.AddNodes(horde3d.RootNode, knightRes)
	horde3d.SetNodeTransform(model, 0, 0, 300, 0, 180, 0, 0.1, 0.1, 0.1)

	// Add light source
	light := horde3d.AddLightNode(horde3d.RootNode, "Light1", 0, "LIGHTING", "SHADOWMAP")
	horde3d.SetNodeTransform(light, 0, 20, 0, 0, 0, 0, 1, 1, 1)
	horde3d.SetNodeParamF(light, horde3d.Light_RadiusF, 0, 50)

	//add camera
	cam := horde3d.AddCameraNode(horde3d.RootNode, "Camera", pipeRes)
	//Setup Camera Viewport
	horde3d.SetNodeParamI(cam, horde3d.Camera_ViewportXI, 0)
	horde3d.SetNodeParamI(cam, horde3d.Camera_ViewportYI, 0)
	horde3d.SetNodeParamI(cam, horde3d.Camera_ViewportWidthI, appWidth)
	horde3d.SetNodeParamI(cam, horde3d.Camera_ViewportHeightI, appHeight)
	horde3d.SetupCameraView(cam, 45, float32(appWidth)/float32(appHeight), 0.1, 1000)
	horde3d.ResizePipelineBuffers(pipeRes, appWidth, appHeight)

	for running {

		horde3d.Render(cam)
		horde3d.FinalizeFrame()
		horde3d.DumpMessages()
		glfw.SwapBuffers()
		running = glfw.Key(glfw.KeyEsc) == 0 &&
			glfw.WindowParam(glfw.Opened) == 1
	}

	horde3d.Release()
}
