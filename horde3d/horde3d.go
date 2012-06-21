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

package horde3d

/*
#cgo LDFLAGS: -lHorde3D
#include "goHorde3D.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

//typedef int H3DRes;
//typedef int H3DNode;
type H3DRes C.int
type H3DNode C.int

//const H3DNode H3DRootNode = 1;
const RootNode H3DNode = 1

/* Group: Enumerations */

/* Enum: H3DOptions
	The available engine option parameters.

MaxLogLevel         - Defines the maximum log level; only messages which are smaller or equal to this value
                      (hence more important) are published in the message queue. (Default: 4)
MaxNumMessages      - Defines the maximum number of messages that can be stored in the message queue (Default: 512)
TrilinearFiltering  - Enables or disables trilinear filtering for textures. (Values: 0, 1; Default: 1)
MaxAnisotropy       - Sets the maximum quality for anisotropic filtering. (Values: 1, 2, 4, 8; Default: 1)
TexCompression      - Enables or disables texture compression; only affects textures that are
                      loaded after setting the option. (Values: 0, 1; Default: 0)
SRGBLinearization   - Eanbles or disables gamma-to-linear-space conversion of input textures that are tagged as sRGB (Values: 0, 1; Default: 0)
LoadTextures        - Enables or disables loading of textures referenced by materials; this can be useful to reduce
                      loading times for testing. (Values: 0, 1; Default: 1)
FastAnimation       - Disables or enables inter-frame interpolation for animations. (Values: 0, 1; Default: 1)
ShadowMapSize       - Sets the size of the shadow map buffer (Values: 128, 256, 512, 1024, 2048; Default: 1024)
SampleCount         - Maximum number of samples used for multisampled render targets; only affects pipelines
                      that are loaded after setting the option. (Values: 0, 2, 4, 8, 16; Default: 0)
WireframeMode       - Enables or disables wireframe rendering
DebugViewMode       - Enables or disables debug view where geometry is rendered in wireframe without shaders and
                      lights are visualized using their screen space bounding box. (Values: 0, 1; Default: 0)
DumpFailedShaders   - Enables or disables storing of shader code that failed to compile in a text file; this can be
                      useful in combination with the line numbers given back by the shader compiler. (Values: 0, 1; Default: 0)
GatherTimeStats     - Enables or disables gathering of time stats that are useful for profiling (Values: 0, 1; Default: 1)
*/
const (
	_ = iota
	Options_MaxLogLevel
	Options_MaxNumMessages
	ptions_TrilinearFiltering
	Options_MaxAnisotropy
	Options_TexCompression
	Options_SRGBLinearization
	Options_LoadTextures
	Options_FastAnimation
	Options_ShadowMapSize
	Options_SampleCount
	Options_WireframeMode
	Options_DebugViewMode
	Options_DumpFailedShaders
	Options_GatherTimeStats
)

/* Enum: H3DStats
	The available engine statistic parameters.

TriCount          - Number of triangles that were pushed to the renderer
BatchCount        - Number of batches (draw calls)
LightPassCount    - Number of lighting passes
FrameTime         - Time in ms between two h3dFinalizeFrame calls
AnimationTime     - CPU time in ms spent for animation
GeoUpdateTime     - CPU time in ms spent for software skinning and morphing
ParticleSimTime   - CPU time in ms spent for particle simulation and updates
FwdLightsGPUTime  - GPU time in ms spent for forward lighting passes
DefLightsGPUTime  - GPU time in ms spent for drawing deferred light volumes
ShadowsGPUTime    - GPU time in ms spent for generating shadow maps
ParticleGPUTime   - GPU time in ms spent for drawing particles
TextureVMem       - Estimated amount of video memory used by textures (in Mb)
GeometryVMem      - Estimated amount of video memory used by geometry (in Mb)
*/
const (
	_ = iota + 100
	Stats_TriCount
	Stats_BatchCount
	Stats_LightPassCount
	Stats_FrameTime
	Stats_AnimationTime
	Stats_GeoUpdateTime
	Stats_ParticleSimTime
	Stats_FwdLightsGPUTime
	Stats_DefLightsGPUTime
	Stats_ShadowsGPUTime
	Stats_ParticleGPUTime
	Stats_TextureVMem
	Stats_GeometryVMem
)

/* Enum: H3DResTypes
	The available resource types.

Undefined       - An undefined resource, returned by getResourceType in case of error
SceneGraph      - Scene graph subtree stored in XML format
Geometry        - Geometrical data containing bones, vertices and triangles
Animation       - Animation data
Material        - Material script
Code            - Text block containing shader source code
Shader          - Shader program
Texture         - Texture map
ParticleEffect  - Particle configuration
Pipeline        - Rendering pipeline
*/
const (
	ResTypes_Undefined = iota
	ResTypes_SceneGraph
	ResTypes_Geometry
	ResTypes_Animation
	ResTypes_Material
	ResTypes_Code
	ResTypes_Shader
	ResTypes_Texture
	ResTypes_ParticleEffect
	ResTypes_Pipeline
)

/* Enum: H3DResFlags
	The available flags used when adding a resource.

NoQuery           - Excludes resource from being listed by queryUnloadedResource function.
NoTexCompression  - Disables texture compression for Texture resource.
NoTexMipmaps      - Disables generation of mipmaps for Texture resource.
TexCubemap        - Sets Texture resource to be a cubemap.
TexDynamic        - Enables more efficient updates of Texture resource streams.
TexRenderable     - Makes Texture resource usable as render target.
TexSRGB           - Indicates that Texture resource is in sRGB color space and should be converted
                    to linear space when being sampled.
*/
const (
	ResFlags_NoQuery          = 1
	ResFlags_NoTexCompression = 2
	ResFlags_NoTexMipmaps     = 4
	ResFlags_TexCubemap       = 8
	ResFlags_TexDynamic       = 16
	ResFlags_TexRenderable    = 32
	ResFlags_TexSRGB          = 64
)

/* Enum: H3DFormats
	The available resource stream formats.

Unknown      - Unknown format
TEX_BGRA8    - 8-bit BGRA texture
TEX_DXT1     - DXT1 compressed texture
TEX_DXT3     - DXT3 compressed texture
TEX_DXT5     - DXT5 compressed texture
TEX_RGBA16F  - Half float RGBA texture
TEX_RGBA32F  - Float RGBA texture
*/
const (
	Formats_Unknown = iota
	Formats_TEX_BGRA8
	Formats_TEX_DXT1
	Formats_TEX_DXT3
	Formats_TEX_DXT5
	Formats_TEX_RGBA16F
	Formats_TEX_RGBA32F
)

/* Enum: H3DGeoRes
	The available Geometry resource accessors.

GeometryElem         - Base element
GeoIndexCountI       - Number of indices [read-only]
GeoVertexCountI      - Number of vertices [read-only]
GeoIndices16I        - Flag indicating whether index data is 16 or 32 bit [read-only]
GeoIndexStream       - Triangle index data (uint16 or uint32, depending on flag)
GeoVertPosStream     - Vertex position data (float x, y, z)
GeoVertTanStream     - Vertex tangent frame data (float nx, ny, nz, tx, ty, tz, tw)
GeoVertStaticStream  - Vertex static attribute data (float u0, v0,
                         float4 jointIndices, float4 jointWeights, float u1, v1)
*/
const (
	_ = iota + 200
	GeoRes_GeometryElem
	GeoRes_GeoIndexCountI
	GeoRes_GeoVertexCountI
	GeoRes_GeoIndices16I
	GeoRes_GeoIndexStream
	GeoRes_GeoVertPosStream
	GeoRes_GeoVertTanStream
	GeoRes_GeoVertStaticStream
)

/* Enum: AnimRes
	The available Animation resource accessors.	  

EntityElem      - Stored animation entities (joints and meshes)
EntFrameCountI  - Number of frames stored for a specific entity [read-only]
*/
const (
	_ = iota + 300
	AnimRes_EntityElem
	AnimRes_EntFrameCountI
)

/* Enum: MatRes
	The available Material resource accessors.

MaterialElem  - Base element
SamplerElem   - Sampler element
UniformElem   - Uniform element
MatClassStr   - Material class
MatLinkI      - Material resource that is linked to this material
MatShaderI    - Shader resource
SampNameStr   - Name of sampler [read-only]
SampTexResI   - Texture resource bound to sampler
UnifNameStr   - Name of uniform [read-only]
UnifValueF4   - Value of uniform (a, b, c, d)
*/
const (
	_ = iota + 400
	MatRes_MaterialElem
	MatRes_SamplerElem
	MatRes_UniformElem
	MatRes_MatClassStr
	MatRes_MatLinkI
	MatRes_MatShaderI
	MatRes_SampNameStr
	MatRes_SampTexResI
	MatRes_UnifNameStr
	MatRes_UnifValueF4
)

/* Enum: ShaderRes
	The available Shader resource accessors.

ContextElem     - Context element 
SamplerElem     - Sampler element
UniformElem     - Uniform element
ContNameStr     - Name of context [read-only]
SampNameStr     - Name of sampler [read-only]
UnifNameStr     - Name of uniform [read-only]
UnifSizeI       - Size (number of components) of uniform [read-only]
UnifDefValueF4  - Default value of uniform (a, b, c, d)
*/
const (
	_ = iota + 600
	ShaderRes_ContextElem
	ShaderRes_SamplerElem
	ShaderRes_UniformElem
	ShaderRes_ContNameStr
	ShaderRes_SampNameStr
	ShaderRes_UnifNameStr
	ShaderRes_UnifSizeI
	ShaderRes_UnifDefValueF4
)

/* Enum: TexRes
	The available Texture resource accessors.

TextureElem     - Base element
ImageElem       - Subresources of the texture. A texture consists, depending on the type,
                  of a number of equally sized slices which again can have a fixed number
                  of mipmaps. Each image element represents the base image of a slice or
                  a single mipmap level of the corresponding slice.
TexFormatI      - Texture format [read-only]
TexSliceCountI  - Number of slices (1 for 2D texture and 6 for cubemap) [read-only]
ImgWidthI       - Image width [read-only]
ImgHeightI      - Image height [read-only]
ImgPixelStream  - Pixel data of an image. The data layout matches the layout specified
                  by the texture format with the exception that half-float is converted
                  to float. The first element in the data array corresponds to the lower
                  left corner.
*/
const (
	_ = iota + 700
	TexRes_TextureElem
	TexRes_ImageElem
	TexRes_TexFormatI
	TexRes_TexSliceCountI
	TexRes_ImgWidthI
	TexRes_ImgHeightI
	TexRes_ImgPixelStream
)

/* Enum: PartEffRes
	The available ParticleEffect resource accessors.

ParticleElem     - General particle configuration
ChanMoveVelElem  - Velocity channel
ChanRotVelElem   - Angular velocity channel
ChanSizeElem     - Size channel
ChanColRElem     - Red color component channel
ChanColGElem     - Green color component channel
ChanColBElem     - Blue color component channel
ChanColAElem     - Alpha channel
PartLifeMinF     - Minimum value of random life time (in seconds)
PartLifeMaxF     - Maximum value of random life time (in seconds)
ChanStartMinF    - Minimum for selecting initial random value of channel
ChanStartMaxF    - Maximum for selecting initial random value of channel
ChanEndRateF     - Remaining percentage of initial value when particle is dying
*/
const (
	_ = iota + 800
	PartEffRes_ParticleElem
	PartEffRes_ChanMoveVelElem
	PartEffRes_ChanRotVelElem
	PartEffRes_ChanSizeElem
	PartEffRes_ChanColRElem
	PartEffRes_ChanColGElem
	PartEffRes_ChanColBElem
	PartEffRes_ChanColAElem
	PartEffRes_PartLifeMinF
	PartEffRes_PartLifeMaxF
	PartEffRes_ChanStartMinF
	PartEffRes_ChanStartMaxF
	PartEffRes_ChanEndRateF
	PartEffRes_ChanDragElem
)

/* Enum: PipeRes
	The available Pipeline resource accessors.

StageElem         - Pipeline stage
StageNameStr      - Name of stage [read-only]
StageActivationI  - Flag indicating whether stage is active
*/
const (
	_ = iota + 900
	PipeRes_StageElem
	PipeRes_StageNameStr
	PipeRes_StageActivationI
)

/*	Enum: NodeTypes
		The available scene node types.

	Undefined  - An undefined node type, returned by getNodeType in case of error
	Group      - Group of different scene nodes
	Model      - 3D model with optional skeleton
	Mesh       - Subgroup of a model with triangles of one material
	Joint      - Joint for skeletal animation
	Light      - Light source
	Camera     - Camera giving view on scene
	Emitter    - Particle system emitter
*/
const (
	NodeTypes_Undefined = iota
	NodeTypes_Group
	NodeTypes_Model
	NodeTypes_Mesh
	NodeTypes_Joint
	NodeTypes_Light
	NodeTypes_Camera
	NodeTypes_Emitter
)

/*	Enum: NodeFlags
		The available scene node flags.

	NoDraw         - Excludes scene node from all rendering
	NoCastShadow   - Excludes scene node from list of shadow casters
	NoRayQuery     - Excludes scene node from ray intersection queries
	Inactive       - Deactivates scene node so that it is completely ignored
	                 (combination of all flags above)
*/
const (
	NodeFlags_NoDraw       = 1
	NodeFlags_NoCastShadow = 2
	NodeFlags_NoRayQuery   = 4
	NodeFlags_Inactive     = 7 // NoDraw | NoCastShadow | NoRayQuery
)

/*	Enum: NodeParams
		The available scene node parameters.

	NameStr        - Name of the scene node
	AttachmentStr  - Optional application-specific meta data for a node encapsulated
	                 in an 'Attachment' XML string
*/
const (
	NodeParams_NameStr = 1
	NodeParams_AttachmentStr
)

/*	Enum: Model
		The available Model node parameters

	GeoResI      - Geometry resource used for the model
	SWSkinningI  - Enables or disables software skinning (default: 0)
	LodDist1F    - Distance to camera from which on LOD1 is used (default: infinite)
	               (must be a positive value larger than 0.0)
	LodDist2F    - Distance to camera from which on LOD2 is used
	               (may not be smaller than LodDist1) (default: infinite)
	LodDist3F    - Distance to camera from which on LOD3 is used
	               (may not be smaller than LodDist2) (default: infinite)
	LodDist4F    - Distance to camera from which on LOD4 is used
	               (may not be smaller than LodDist3) (default: infinite)
*/
const (
	_             = iota + 200
	Model_GeoResI = 200
	Model_SWSkinningI
	Model_LodDist1F
	Model_LodDist2F
	Model_LodDist3F
	Model_LodDist4F
)

/*	Enum: Mesh
		The available Mesh node parameters.

	MatResI      - Material resource used for the mesh
	BatchStartI  - First triangle index of mesh in Geometry resource of parent Model node [read-only]
	BatchCountI  - Number of triangle indices used for drawing mesh [read-only]
	VertRStartI  - First vertex in Geometry resource of parent Model node [read-only]
	VertREndI    - Last vertex in Geometry resource of parent Model node [read-only]
	LodLevelI    - LOD level of Mesh; the mesh is only rendered if its LOD level corresponds to
	               the model's current LOD level which is calculated based on the LOD distances (default: 0)
*/
const (
	_ = iota + 300
	Mesh_MatResI
	Mesh_BatchStartI
	Mesh_BatchCountI
	Mesh_VertRStartI
	Mesh_VertREndI
	Mesh_LodLevelI
)

/*	Enum: Joint
		The available Joint node parameters.

	JointIndexI  - Index of joint in Geometry resource of parent Model node [read-only]
*/
const (
	Joint_JointIndexI = 400
)

/*	Enum: Light
		The available Light node parameters.

	MatResI             - Material resource used for the light
	RadiusF             - Radius of influence (default: 100.0)
	FovF                - Field of view (FOV) angle (default: 90.0)
	ColorF3             - Diffuse color RGB (default: 1.0, 1.0, 1.0)
	ColorMultiplierF    - Diffuse color multiplier for altering intensity, mainly useful for HDR (default: 1.0)
	ShadowMapCountI     - Number of shadow maps used for light source (values: 0, 1, 2, 3, 4; default: 0)]
	ShadowSplitLambdaF  - Constant determining segmentation of view frustum for Parallel Split Shadow Maps (default: 0.5)
	ShadowMapBiasF      - Bias value for shadow mapping to reduce shadow acne (default: 0.005)
	LightingContextStr  - Name of shader context used for computing lighting
	ShadowContextStr    - Name of shader context used for generating shadow map
*/
const (
	_ = iota + 500
	Light_MatResI
	Light_RadiusF
	Light_FovF
	Light_ColorF3
	Light_ColorMultiplierF
	Light_ShadowMapCountI
	Light_ShadowSplitLambdaF
	Light_ShadowMapBiasF
	Light_LightingContextStr
	Light_ShadowContextStr
)

/*	Enum: Camera
		The available Camera node parameters.

	PipeResI         - Pipeline resource used for rendering
	OutTexResI       - 2D Texture resource used as output buffer (can be 0 to use main framebuffer) (default: 0)
	OutBufIndexI     - Index of the output buffer for stereo rendering (values: 0 for left eye, 1 for right eye) (default: 0)
	LeftPlaneF       - Coordinate of left plane relative to near plane center (default: -0.055228457)
	RightPlaneF      - Coordinate of right plane relative to near plane center (default: 0.055228457)
	BottomPlaneF     - Coordinate of bottom plane relative to near plane center (default: -0.041421354f)
	TopPlaneF        - Coordinate of top plane relative to near plane center (default: 0.041421354f)
	NearPlaneF       - Distance of near clipping plane (default: 0.1)
	FarPlaneF        - Distance of far clipping plane (default: 1000)
	ViewportXI       - Position x-coordinate of the lower left corner of the viewport rectangle (default: 0)
	ViewportYI       - Position y-coordinate of the lower left corner of the viewport rectangle (default: 0)
	ViewportWidthI   - Width of the viewport rectangle (default: 320)
	ViewportHeightI  - Height of the viewport rectangle (default: 240)
	OrthoI           - Flag for setting up an orthographic frustum instead of a perspective one (default: 0)
	OccCullingI      - Flag for enabling occlusion culling (default: 0)
*/

const (
	_ = iota + 600
	Camera_PipeResI
	Camera_OutTexResI
	Camera_OutBufIndexI
	Camera_LeftPlaneF
	Camera_RightPlaneF
	Camera_BottomPlaneF
	Camera_TopPlaneF
	Camera_NearPlaneF
	Camera_FarPlaneF
	Camera_ViewportXI
	Camera_ViewportYI
	Camera_ViewportWidthI
	Camera_ViewportHeightI
	Camera_OrthoI
	Camera_OccCullingI
)

/*	Enum: Emitter
		The available Emitter node parameters.

	MatResI        - Material resource used for rendering
	PartEffResI    - ParticleEffect resource which configures particle properties
	MaxCountI      - Maximal number of particles living at the same time
	RespawnCountI  - Number of times a single particle is recreated after dying (-1 for infinite)
	DelayF         - Time in seconds before emitter begins creating particles (default: 0.0)
	EmissionRateF  - Maximal number of particles to be created per second (default: 0.0)
	SpreadAngleF   - Angle of cone for random emission direction (default: 0.0)
	ForceF3        - Force vector XYZ applied to particles (default: 0.0, 0.0, 0.0)
*/
const (
	_ = iota + 700
	Emitter_MatResI
	Emitter_PartEffResI
	Emitter_MaxCountI
	Emitter_RespawnCountI
	Emitter_DelayF
	Emitter_EmissionRateF
	Emitter_SpreadAngleF
	Emitter_ForceF3
)

//used for bools from c interfaces
// I'm sure there's a better way to do this, but this works for now
var Bool = map[int]bool{
	0: false,
	1: true,
}
var Int = map[bool]C.int{
	false: C.int(1),
	true:  C.int(0),
}

func GetVersionString() string {
	verPointer := C.h3dGetVersionString()
	//defer C.free(unsafe.Pointer(verPointer))

	return C.GoString(verPointer)
}

func CheckExtension(extensionName string) bool {
	cExtName := C.CString(extensionName)
	defer C.free(unsafe.Pointer(cExtName))
	return Bool[int(C.h3dCheckExtension(cExtName))]
}

func GetError() bool {
	return Bool[int(C.h3dGetError())]
}

func Init() bool {
	return Bool[int(C.h3dInit())]
}

func Release() {
	C.h3dRelease()
}

func Render(cameraNode H3DNode) {
	C.h3dRender(C.H3DNode(cameraNode))
}

func FinalizeFrame() {
	C.h3dFinalizeFrame()
}

func Clear() {
	C.h3dClear()
}

func GetMessage(level *int, time *float32) string {
	message := C.h3dGetMessage((*C.int)(unsafe.Pointer(level)),
		(*C.float)(unsafe.Pointer(time)))
	defer C.free(unsafe.Pointer(message))
	return C.GoString(message)
}

func GetOption(param int) float32 {
	return float32(C.h3dGetOption(C.int(param)))
}

func SetOption(param int, value float32) bool {
	return Bool[int(C.h3dSetOption(C.int(param), C.float(value)))]
}

func GetStat(param int, reset bool) float32 {
	return float32(C.h3dGetStat(C.int(param), Int[reset]))
}

func ShowOverlays(verts []float32,
	vertCount int,
	colR float32,
	colG float32,
	colB float32,
	colA float32,
	materialRes H3DRes,
	flags int) {

	C.h3dShowOverlays((*C.float)(unsafe.Pointer(&verts[0])),
		C.int(vertCount),
		C.float(colR),
		C.float(colG),
		C.float(colB),
		C.float(colA),
		C.H3DRes(materialRes),
		C.int(flags))
}

func ClearOverlays() {
	C.h3dClearOverlays()
}

func GetResType(res H3DRes) int {
	return int(C.h3dGetResType(C.H3DRes(res)))
}

func GetResName(res H3DRes) string {
	return C.GoString(C.h3dGetResName(C.H3DRes(res)))
}

func GetNextResource(resType int, start H3DRes) H3DRes {
	return H3DRes(C.h3dGetNextResource(C.int(resType), C.H3DRes(start)))
}

func FindResource(resType int, name string) H3DRes {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	return H3DRes(C.h3dFindResource(C.int(resType), cName))
}

func AddResource(resType int, name string, flags int) H3DRes {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	return H3DRes(C.h3dAddResource(C.int(resType), cName, C.int(flags)))
}

func CloneResource(sourceRes H3DRes, name string) H3DRes {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	return H3DRes(C.h3dCloneResource(C.H3DRes(sourceRes), cName))
}

func RemoveResource(res H3DRes) int {
	return int(C.h3dRemoveResource(C.H3DRes(res)))
}

func IsResLoaded(res H3DRes) bool {
	return Bool[int(C.h3dIsResLoaded(C.H3DRes(res)))]
}

func LoadResource(res H3DRes, data []byte, size int) bool {
	return Bool[int(C.h3dLoadResource(C.H3DRes(res), (*C.char)(unsafe.Pointer(&data[0])), C.int(size)))]
}

func UnloadResource(res H3DRes) {
	C.h3dUnloadResource(C.H3DRes(res))
}

func GetResElemCount(res H3DRes, elem int) int {
	return int(C.h3dGetResElemCount(C.H3DRes(res), C.int(elem)))
}

func FindResElem(res H3DRes, elem int, param int, value string) int {
	cValue := C.CString(value)
	defer C.free(unsafe.Pointer(cValue))
	return int(C.h3dFindResElem(C.H3DRes(res),
		C.int(elem), C.int(param), cValue))
}

func GetResParamI(res H3DRes, elem int, elemIdx int, param int) int {
	return int(C.h3dGetResParamI(C.H3DRes(res), C.int(elem), C.int(elemIdx), C.int(param)))
}

func SetResParamI(res H3DRes, elem int, elemIdx int, param int, value int) {
	C.h3dSetResParamI(C.H3DRes(res), C.int(elem), C.int(elemIdx), C.int(param), C.int(value))
}

func GetResParamF(res H3DRes, elem int, elemIdx int, param int, compIdx int) float32 {
	return float32(C.h3dGetResParamF(C.H3DRes(res),
		C.int(elem), C.int(elemIdx), C.int(param), C.int(compIdx)))
}

func SetResParamF(res H3DRes, elem int, elemIdx int, param int, compIdx int, value float32) {
	C.h3dSetResParamF(C.H3DRes(res),
		C.int(elem), C.int(elemIdx), C.int(param), C.int(compIdx), C.float(value))
}

func GetResParamStr(res H3DRes, elem int, elemIdx int, param int) string {
	value := C.h3dGetResParamStr(C.H3DRes(res), C.int(elem), C.int(elemIdx), C.int(param))
	defer C.free(unsafe.Pointer(value))
	return C.GoString(value)
}

func SetResParamStr(res H3DRes, elem int, elemIdx int, param int, value string) {
	cValue := C.CString(value)
	defer C.free(unsafe.Pointer(cValue))
	C.h3dSetResParamStr(C.H3DRes(res), C.int(elem), C.int(elemIdx), C.int(param), cValue)
}

func MapResStream(res H3DRes, elem int, elemIdx int, stream int, read bool, write bool, size int) []byte {

	cStream := C.h3dMapResStream(C.H3DRes(res), C.int(elem), C.int(elemIdx), C.int(stream),
		Int[read], Int[write])

	C.free(cStream)

	return C.GoBytes(unsafe.Pointer(cStream), C.int(size))
}

func UnmapResStream(res H3DRes) {
	C.h3dUnmapResStream(C.H3DRes(res))
}

func QueryUnloadedResource(index int) H3DRes {
	return H3DRes(C.h3dQueryUnloadedResource(C.int(index)))
}

func ReleaseUnusedResources() {
	C.h3dReleaseUnusedResources()
}

func CreateTexture(name string, width int, height int, fmt int, flags int) H3DRes {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	return H3DRes(C.h3dCreateTexture(cName, C.int(width), C.int(height), C.int(fmt), C.int(flags)))
}

func SetShaderPreambles(vertPreamble string, fragPreamble string) {
	cVertPreamble := C.CString(vertPreamble)
	cFragPreamble := C.CString(fragPreamble)
	defer C.free(unsafe.Pointer(cVertPreamble))
	defer C.free(unsafe.Pointer(cFragPreamble))

	C.h3dSetShaderPreambles(cVertPreamble, cFragPreamble)
}

func SetMaterialUniform(materialRes H3DRes, name string, a float32, b float32,
	c float32, d float32) bool {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	return Bool[int(C.h3dSetMaterialUniform(C.H3DRes(materialRes), cName,
		C.float(a), C.float(b), C.float(c), C.float(d)))]
}

func ResizePipelineBuffers(pipeRes H3DRes, width int, height int) {
	C.h3dResizePipelineBuffers(C.H3DRes(pipeRes), C.int(width), C.int(height))
}

func GetRenderTargetData(pipelineRes H3DRes, targetName string, bufIndex int, width *int,
	height *int, compCount *int, dataBuffer []byte, bufferSize int) bool {
	cTargetName := C.CString(targetName)
	defer C.free(unsafe.Pointer(cTargetName))

	var cDataBuffer unsafe.Pointer
	defer C.free(cDataBuffer)
	var targetFound bool
	targetFound = Bool[int(C.h3dGetRenderTargetData(C.H3DRes(pipelineRes), cTargetName,
		C.int(bufIndex), (*C.int)(unsafe.Pointer(width)), (*C.int)(unsafe.Pointer(height)),
		(*C.int)(unsafe.Pointer(compCount)), cDataBuffer, C.int(bufferSize)))]
	dataBuffer = C.GoBytes(cDataBuffer, C.int(bufferSize))
	return targetFound
}

func GetNodeType(node H3DNode) int {
	return int(C.h3dGetNodeType(C.H3DNode(node)))
}

func GetNodeParent(node H3DNode) H3DNode {
	return H3DNode(C.h3dGetNodeParent(C.H3DNode(node)))
}

func SetNodeParent(node H3DNode, parent H3DNode) bool {
	return Bool[int(C.h3dSetNodeParent(C.H3DNode(node), C.H3DNode(parent)))]
}

func GetNodeChild(node H3DNode, index int) bool {
	return Bool[int(C.h3dGetNodeChild(C.H3DNode(node), C.int(index)))]
}

func AddNodes(parent H3DNode, sceneGraphRes H3DRes) H3DNode {
	return H3DNode(C.h3dAddNodes(C.H3DNode(parent), C.H3DRes(sceneGraphRes)))
}

func RemoveNode(node H3DNode) {
	C.h3dRemoveNode(C.H3DNode(node))
}

func CheckNodeTransFlag(node H3DNode, reset bool) bool {
	return Bool[int(C.h3dCheckNodeTransFlag(C.H3DNode(node), Int[reset]))]
}

func GetNodeTransform(node H3DNode, tx *float32, ty *float32, tz *float32,
	rx *float32, ry *float32, rz *float32, sx *float32, sy *float32, sz *float32) {
	C.h3dGetNodeTransform(C.H3DNode(node), (*C.float)(unsafe.Pointer(tx)), (*C.float)(unsafe.Pointer(ty)),
		(*C.float)(unsafe.Pointer(tz)), (*C.float)(unsafe.Pointer(rx)), (*C.float)(unsafe.Pointer(ry)),
		(*C.float)(unsafe.Pointer(rz)), (*C.float)(unsafe.Pointer(sx)), (*C.float)(unsafe.Pointer(sy)),
		(*C.float)(unsafe.Pointer(sz)))
}

func SetNodeTransform(node H3DNode, tx float32, ty float32, tz float32,
	rx float32, ry float32, rz float32, sx float32, sy float32, sz float32) {
	C.h3dSetNodeTransform(C.H3DNode(node), C.float(tx), C.float(ty), C.float(tz),
		C.float(rx), C.float(ry), C.float(rz), C.float(sx), C.float(sy), C.float(sz))
}

func GetNodeTransMats(node H3DNode, relMat []float32, absMat []float32) {
	//TODO: Handle nil pointers, possibly check for nil

	C.h3dGetNodeTransMats(C.H3DNode(node), (**C.float)(unsafe.Pointer(&relMat[0])),
		(**C.float)(unsafe.Pointer(&absMat[0])))
}

func SetNodeTransMat(node H3DNode, mat4x4 []float32) {
	C.h3dSetNodeTransMat(C.H3DNode(node), (*C.float)(unsafe.Pointer(&mat4x4[0])))
}

func GetNodeParamI(node H3DNode, param int) int {
	return int(C.h3dGetNodeParamI(C.H3DNode(node), C.int(param)))
}

func SetNodeParamI(node H3DNode, param int, value int) {
	C.h3dSetNodeParamI(C.H3DNode(node), C.int(param), C.int(value))
}

func GetNodeParamF(node H3DNode, param int, compIdx int) float32 {
	return float32(C.h3dGetNodeParamF(C.H3DNode(node), C.int(param), C.int(compIdx)))
}

func SetNodeParamF(node H3DNode, param int, compIdx int, value float32) {
	C.h3dSetNodeParamF(C.H3DNode(node), C.int(param), C.int(compIdx), C.float(value))
}

func GetNodeParamStr(node H3DNode, param int) string {
	value := C.h3dGetNodeParamStr(C.H3DNode(node), C.int(param))
	defer C.free(unsafe.Pointer(value))
	return C.GoString(value)
}

func SetNodeParamStr(node H3DNode, param int, value string) {
	cValue := C.CString(value)
	C.free(unsafe.Pointer(cValue))
	C.h3dSetNodeParamStr(C.H3DNode(node), C.int(param), cValue)
}

func GetNodeFlags(node H3DNode) int {
	return int(C.h3dGetNodeFlags(C.H3DNode(node)))
}

func SetNodeFlags(node H3DNode, flags int, recursive bool) {
	C.h3dSetNodeFlags(C.H3DNode(node), C.int(flags), Int[recursive])
}

func GetNodeAABB(node H3DNode, minX *float32, minY *float32, minZ *float32,
	maxX *float32, maxY *float32, maxZ *float32) {
	C.h3dGetNodeAABB(C.H3DNode(node), (*C.float)(unsafe.Pointer(minX)),
		(*C.float)(unsafe.Pointer(minY)), (*C.float)(unsafe.Pointer(minZ)),
		(*C.float)(unsafe.Pointer(maxX)), (*C.float)(unsafe.Pointer(maxY)),
		(*C.float)(unsafe.Pointer(maxZ)))
}

func FindNodes(node H3DNode, name string, nodeType int) int {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	return int(C.h3dFindNodes(C.H3DNode(node), cName, C.int(nodeType)))
}

func GetNodeFindResult(index int) H3DNode {
	return H3DNode(C.h3dGetNodeFindResult(C.int(index)))
}

func CastRay(node H3DNode, ox float32, oy float32, oz float32,
	dx float32, dy float32, dz float32, numNearest int) int {
	return int(C.h3dCastRay(C.H3DNode(node), C.float(ox), C.float(oy), C.float(oz),
		C.float(dx), C.float(dy), C.float(dz), C.int(numNearest)))
}

func GetCastRayResult(index int, node *H3DNode, distance *float32, intersection []float32) bool {
	return Bool[int(C.h3dGetCastRayResult(C.int(index), (*C.H3DNode)(unsafe.Pointer(node)),
		(*C.float)(unsafe.Pointer(distance)), (*C.float)(unsafe.Pointer(&intersection[0]))))]
}

func CheckNodeVisibility(node H3DNode, cameraNode H3DNode, checkOcclusion bool, calcLod bool) int {
	return int(C.h3dCheckNodeVisibility(C.H3DNode(node), C.H3DNode(cameraNode),
		Int[checkOcclusion], Int[calcLod]))
}

func AddGroupNode(parent H3DNode, name string) H3DNode {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	return H3DNode(C.h3dAddGroupNode(C.H3DNode(parent), cName))
}

func AddModelNode(parent H3DNode, name string, geometryRes H3DRes) H3DNode {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	return H3DNode(C.h3dAddModelNode(C.H3DNode(parent), cName, C.H3DRes(geometryRes)))
}

func SetupModelAnimStage(modelNode H3DNode, stage int, animationRes H3DRes, layer int,
	startNode string, additive bool) {
	cStartNode := C.CString(startNode)
	defer C.free(unsafe.Pointer(cStartNode))
	C.h3dSetupModelAnimStage(C.H3DNode(modelNode), C.int(stage), C.H3DRes(animationRes),
		C.int(layer), cStartNode, Int[additive])
}

func SetModelAnimParams(modelNode H3DNode, stage int, time float32, weight float32) {
	C.h3dSetModelAnimParams(C.H3DNode(modelNode), C.int(stage), C.float(time), C.float(weight))
}

func SetModelMorpher(modelNode H3DNode, target string, weight float32) bool {
	cTarget := C.CString(target)
	defer C.free(unsafe.Pointer(cTarget))
	return Bool[int(C.h3dSetModelMorpher(C.H3DNode(modelNode), cTarget, C.float(weight)))]
}

func AddMeshNode(parent H3DNode, name string, materialRes H3DRes, batchStart int, batchCount int,
	vertRStart int, vertEnd int) H3DNode {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	return H3DNode(C.h3dAddMeshNode(C.H3DNode(parent), cName, C.H3DRes(materialRes), C.int(batchStart),
		C.int(batchCount), C.int(vertRStart), C.int(vertEnd)))
}

func AddJointNode(parent H3DNode, name string, jointIndex int) H3DNode {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	return H3DNode(C.h3dAddJointNode(C.H3DNode(parent), cName, C.int(jointIndex)))
}

func AddLightNode(parent H3DNode, name string, materialRes H3DRes, lightingContext string,
	shadowContext string) H3DNode {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	cLightingContext := C.CString(lightingContext)
	defer C.free(unsafe.Pointer(cLightingContext))
	cShadowContext := C.CString(shadowContext)
	defer C.free(unsafe.Pointer(cShadowContext))

	return H3DNode(C.h3dAddLightNode(C.H3DNode(parent), cName, C.H3DRes(materialRes), cLightingContext,
		cShadowContext))
}

func AddCameraNode(parent H3DNode, name string, pipelineRes H3DRes) H3DNode {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	return H3DNode(C.h3dAddCameraNode(C.H3DNode(parent), cName, C.H3DRes(pipelineRes)))
}

func SetupCameraView(cameraNode H3DNode, fov float32, aspect float32,
	nearDist float32, farDist float32) {
	C.h3dSetupCameraView(C.H3DNode(cameraNode), C.float(fov), C.float(aspect),
		C.float(nearDist), C.float(farDist))
}

func GetCameraProjMat(cameraNode H3DNode, projMat []float32) {

	C.h3dGetCameraProjMat(C.H3DNode(cameraNode), (*C.float)(unsafe.Pointer(&projMat[0])))
}

func AddEmitterNode(parent H3DNode, name string, materialRes H3DRes, particleEffectRes H3DRes,
	maxParticleCount int, respawnCount int) H3DNode {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	return H3DNode(C.h3dAddEmitterNode(C.H3DNode(parent), cName, C.H3DRes(materialRes),
		C.H3DRes(particleEffectRes), C.int(maxParticleCount), C.int(respawnCount)))
}

func AdvanceEmitterTime(emitterNode H3DNode, timeDelta float32) {
	C.h3dAdvanceEmitterTime(C.H3DNode(emitterNode), C.float(timeDelta))
}

func HasEmitterFinished(emitterNode H3DNode) bool {
	return Bool[int(C.h3dHasEmitterFinished(C.H3DNode(emitterNode)))]
}
