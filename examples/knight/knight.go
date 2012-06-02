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
	caption   = "Knight - Horde3D Sample (Go Implementation)"
	appWidth  = 800
	appHeight = 600
)

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

	horde3d.Init()
	//pipeline
	hdrPipeline := horde3d.AddResource(horde3d.H3DResTypes_Pipeline, "hdr.pipeline.xml", 0)

	//add camera
	cam := horde3d.AddCameraNode(horde3d.RootNode, "Camera", hdrPipeline)
	//Setup Camera Viewport
	horde3d.SetNodeParamI(cam, horde3d.H3DCamera_ViewportXI, 0)
	horde3d.SetNodeParamI(cam, horde3d.H3DCamera_ViewportYI, 0)
	horde3d.SetNodeParamI(cam, horde3d.H3DCamera_ViewportWidthI, appWidth)
	horde3d.SetNodeParamI(cam, horde3d.H3DCamera_ViewportHeightI, appHeight)

	//enable vertical sync if the card supports it
	glfw.SetSwapInterval(1)

	glfw.SetWindowTitle("Horde3d Knight demo implemented in Go")

	//load resources paths separated by |
	horde3d.LoadResourcesFromDisk("../content|" +
		"../content/pipelines")
	for running {
		horde3d.Render(cam)
		glfw.SwapBuffers()
		running = glfw.Key(glfw.KeyEsc) == 0 &&
			glfw.WindowParam(glfw.Opened) == 1
	}

	horde3d.DumpMessages()
	horde3d.Release()
}
