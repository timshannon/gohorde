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
	// Set options
	horde3d.SetOption(horde3d.H3DOptions_LoadTextures, 1)
	horde3d.SetOption(horde3d.H3DOptions_TexCompression, 0)
	horde3d.SetOption(horde3d.H3DOptions_FastAnimation, 0)
	horde3d.SetOption(horde3d.H3DOptions_MaxAnisotropy, 8)
	horde3d.SetOption(horde3d.H3DOptions_ShadowMapSize, 2048)

	// Add resources
	//pipeline
	hdrPipeline := horde3d.AddResource(horde3d.H3DResTypes_Pipeline, "hdr.pipeline.xml", 0)

	// Font
	//fontMatRes := horde3d.AddResource(horde3d.H3DResTypes_Material, "font.material.xml", 0)
	// Logo
	//logoMatRes := horde3d.AddResource(horde3d.H3DResTypes_Material, "logo.material.xml", 0)
	// Environment
	envRes := horde3d.AddResource(horde3d.H3DResTypes_SceneGraph, "sphere.scene.xml", 0)
	// Knight
	knightRes := horde3d.AddResource(horde3d.H3DResTypes_SceneGraph, "knight.scene.xml", 0)
	knightAnim1Res := horde3d.AddResource(horde3d.H3DResTypes_Animation, "knight_order.anim", 0)
	knightAnim2Res := horde3d.AddResource(horde3d.H3DResTypes_Animation, "knight_attack.anim", 0)
	// Particle system
	particleSysRes := horde3d.AddResource(horde3d.H3DResTypes_SceneGraph, "particleSys1.scene.xml", 0)

	//Add Scene Nodes
	//add camera
	cam := horde3d.AddCameraNode(horde3d.RootNode, "Camera", hdrPipeline)
	//Setup Camera Viewport
	horde3d.SetNodeParamI(cam, horde3d.H3DCamera_ViewportXI, 0)
	horde3d.SetNodeParamI(cam, horde3d.H3DCamera_ViewportYI, 0)
	horde3d.SetNodeParamI(cam, horde3d.H3DCamera_ViewportWidthI, appWidth)
	horde3d.SetNodeParamI(cam, horde3d.H3DCamera_ViewportHeightI, appHeight)
	//env
	env := horde3d.AddNodes(horde3d.RootNode, envRes)
	horde3d.SetNodeTransform(env, 0, -20, 0, 0, 0, 0, 20, 20, 20)
	// Add knight
	knight := horde3d.AddNodes(horde3d.RootNode, knightRes)
	horde3d.SetNodeTransform(knight, 0, 0, 0, 0, 180, 0, 0.1, 0.1, 0.1)
	horde3d.SetupModelAnimStage(knight, 0, knightAnim1Res, 0, "", false)
	horde3d.SetupModelAnimStage(knight, 1, knightAnim2Res, 0, "", false)
	// Attach particle system to hand joint
	horde3d.FindNodes(knight, "Bip01_R_Hand", horde3d.H3DNodeTypes_Joint)
	hand := horde3d.GetNodeFindResult(0)
	particleSys := horde3d.AddNodes(hand, particleSysRes)
	horde3d.SetNodeTransform(particleSys, 0, 40, 0, 90, 0, 0, 1, 1, 1)

	// Add light source
	light := horde3d.AddLightNode(horde3d.RootNode, "Light1", 0, "LIGHTING", "SHADOWMAP")
	horde3d.SetNodeTransform(light, 0, 15, 10, -60, 0, 0, 1, 1, 1)
	horde3d.SetNodeParamF(light, horde3d.H3DLight_RadiusF, 0, 30)
	horde3d.SetNodeParamF(light, horde3d.H3DLight_FovF, 0, 90)
	horde3d.SetNodeParamI(light, horde3d.H3DLight_ShadowMapCountI, 1)
	horde3d.SetNodeParamF(light, horde3d.H3DLight_ShadowMapBiasF, 0, 0.01)
	horde3d.SetNodeParamF(light, horde3d.H3DLight_ColorF3, 0, 1.0)
	horde3d.SetNodeParamF(light, horde3d.H3DLight_ColorF3, 1, 0.8)
	horde3d.SetNodeParamF(light, horde3d.H3DLight_ColorF3, 2, 0.7)

	// Customize post processing effects
	matRes := horde3d.FindResource(horde3d.H3DResTypes_Material, "postHDR.material.xml")
	// hdrParams: exposure, brightpass threshold, brightpass offset (see shader for description)
	horde3d.SetMaterialUniform(matRes, "hdrParams", 2.5, 0.5, 0.08, 0)

	//enable vertical sync if the card supports it
	glfw.SetSwapInterval(1)

	glfw.SetWindowTitle("Horde3d Knight demo implemented in Go")

	//load resources paths separated by |
	horde3d.LoadResourcesFromDisk("../content|" +
		"../content/pipelines|" +
		"../content/models|" +
		"../content/materials|" +
		"../content/shaders|" +
		"../content/textures|" +
		"../content/effects")

	for running {
		horde3d.Render(cam)
		glfw.SwapBuffers()
		running = glfw.Key(glfw.KeyEsc) == 0 &&
			glfw.WindowParam(glfw.Opened) == 1
	}

	horde3d.DumpMessages()
	horde3d.Release()
}
