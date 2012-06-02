package horde3d

/*
#cgo LDFLAGS: -lHorde3DUtils
#include "goHorde3DUtils.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

const H3DUTMaxStatMode int = 2

func H3dutFreeMem(ptr unsafe.Pointer) {
	//TODO: Review
	C.h3dutFreeMem((**C.char)(ptr))
}

func H3dutDumpMessages() bool {
	return Bool[int(C.h3dutDumpMessages())]
}

//Not implementing
/*
h3dutInitOpenGL
h3dutReleaseOpenGL
h3dutSwapBuffers
h3dutGetResourcePath
h3dutSetResourcePath
*/

func H3dutLoadResourcesFromDisk(contentDir string) bool {
	cContentDir := C.CString(contentDir)
	defer C.free(unsafe.Pointer(cContentDir))

	return Bool[int(C.h3dutLoadResourcesFromDisk(cContentDir))]
}

func H3dutCreateGeometryRes(name string, numVertices int, numTriangleIndices int, posData []float32,
	indexData []uint, normalData []int16, tangentData []int16, bitangentData []int16,
	textData1 []float32, textData2 []float32) H3DRes {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	return H3DRes(C.h3dutCreateGeometryRes(cName, C.int(numVertices), C.int(numTriangleIndices),
		(*C.float)(unsafe.Pointer(&posData[0])), (*C.uint)(unsafe.Pointer(&indexData[0])),
		(*C.short)(unsafe.Pointer(&normalData[0])), (*C.short)(unsafe.Pointer(&tangentData[0])),
		(*C.short)(unsafe.Pointer(&bitangentData[0])), (*C.float)(unsafe.Pointer(&textData1[0])),
		(*C.float)(unsafe.Pointer(&textData2[0]))))

}

//TODO: Implement
//func H3dutCreateTGAImage(

func H3dutScreenshot(filename string) bool {
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	return Bool[int(C.h3dutScreenshot(cFilename))]
}
func H3dutPickRay(cameraNode H3DNode, nwx float32, nwy float32, ox *float32, oy *float32, oz *float32,
	dx *float32, dy *float32, dz *float32) {
	C.h3dutPickRay(C.H3DNode(cameraNode), C.float(nwx), C.float(nwy), (*C.float)(unsafe.Pointer(ox)),
		(*C.float)(unsafe.Pointer(oy)), (*C.float)(unsafe.Pointer(oz)),
		(*C.float)(unsafe.Pointer(dx)), (*C.float)(unsafe.Pointer(dy)),
		(*C.float)(unsafe.Pointer(dz)))
}

func H3dutPickNode(cameraNode H3DNode, nwx float32, nwy float32) H3DNode {
	return H3DNode(C.h3dutPickNode(C.H3DNode(cameraNode), C.float(nwx), C.float(nwy)))
}

func H3dutShowText(text string, x float32, y float32, size float32, colR float32, colG float32, colB float32,
	fontMaterialRes H3DRes) {
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))
	C.h3dutShowText(cText, C.float(x), C.float(y), C.float(size), C.float(colR), C.float(colG),
		C.float(colB), C.H3DRes(fontMaterialRes))
}

func H3dutShowFrameStats(fontMaterialRes H3DRes, panelMaterialRes H3DRes, mode int) {
	C.h3dutShowFrameStats(C.H3DRes(fontMaterialRes), C.H3DRes(panelMaterialRes), C.int(mode))
}
