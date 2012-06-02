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
