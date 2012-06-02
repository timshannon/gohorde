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
	var cam horde3d.H3DNode
	var width int = 800
	var height int = 600

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

	//Setup Camera Viewport
	horde3d.SetNodeParamI(cam, horde3d.H3DCamera_ViewportXI, 0)
	horde3d.SetNodeParamI(cam, horde3d.H3DCamera_ViewportYI, 0)
	horde3d.SetNodeParamI(cam, horde3d.H3DCamera_ViewportWidthI, width)
	horde3d.SetNodeParamI(cam, horde3d.H3DCamera_ViewportHeightI, height)

	//pipeline
	pipeRes := horde3d.AddResource(horde3d.H3DResTypes_Pipeline, "standard.pipeline.xml", 0)
	modelRes := horde3d.AddResource(horde3d.H3DResTypes_SceneGraph, "character.scene.xml", 0)
	animRes := horde3d.AddResource(horde3d.H3DResTypes_Animation, "walk.anim.xml", 0)

	horde3d.LoadResourcesFromDisk("")

	model := horde3d.AddNodes(horde3d.H3DRootNode, modelRes)
	horde3d.SetupModelAnimStage(model, 0, animRes, 0, "", false)
	for running {
		switch event := sdl.PollEvent(); event.(type) {
		case *sdl.QuitEvent:
			running = false
			break
		}
		horde3d.Render(cam)
		sdl.GL_SwapBuffers()
	}
	horde3d.Release()
	sdl.Quit()

}
