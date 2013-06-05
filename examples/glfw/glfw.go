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
	"os"
)

const (
	caption   = "GLFW Sample (Go Implementation)"
	appWidth  = 800
	appHeight = 600
)

var pipeRes horde3d.H3DRes
var cam horde3d.H3DNode

func main() {
	var running bool = true

	if err := glfw.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "[e] %v\n", err)
		return
	}

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

	//horde3d.SetOption(horde3d.Options_DebugViewMode, 1)
	// Add resources
	//pipeline
	pipeRes = horde3d.AddResource(horde3d.ResTypes_Pipeline, "pipelines/hdr.pipeline.xml", 0)

	knightRes := horde3d.AddResource(horde3d.ResTypes_SceneGraph, "models/knight/knight.scene.xml", 0)

	//load resources paths separated by |
	horde3d.LoadResourcesFromDisk("../content")

	model := horde3d.RootNode.AddNodes(knightRes)
	model.SetTransform(0, 0, -30, 0, 0, 0, 0.1, 0.1, 0.1)

	// Add light source
	light := horde3d.RootNode.AddLightNode("Light1", 0, "LIGHTING", "SHADOWMAP")
	light.SetTransform(0, 20, 0, 0, 0, 0, 1, 1, 1)
	light.SetNodeParamF(horde3d.Light_RadiusF, 0, 50)

	//add camera
	cam = horde3d.RootNode.AddCameraNode("Camera", pipeRes)
	glfw.SetWindowSizeCallback(onResize)

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

func onResize(w, h int) {
	if h == 0 {
		h = 1
	}

	cam.SetNodeParamI(horde3d.Camera_ViewportXI, 0)
	cam.SetNodeParamI(horde3d.Camera_ViewportYI, 0)
	cam.SetNodeParamI(horde3d.Camera_ViewportWidthI, w)
	cam.SetNodeParamI(horde3d.Camera_ViewportHeightI, h)

	horde3d.SetupCameraView(cam, 45.0, float32(w)/float32(h), 0.1, 1000.0)
	horde3d.ResizePipelineBuffers(pipeRes, w, h)

}
