package cl

/*
#include "./opencl.h"
*/
import "C"

import (
	"fmt"
	"unsafe"
)

//////////////// Constants ////////////////
const maxImageFormats = 256

//////////////// Basic Types ////////////////
type ChannelOrder int

const (
	ChannelOrderR         ChannelOrder = C.CL_R
	ChannelOrderA         ChannelOrder = C.CL_A
	ChannelOrderRG        ChannelOrder = C.CL_RG
	ChannelOrderRA        ChannelOrder = C.CL_RA
	ChannelOrderRGB       ChannelOrder = C.CL_RGB
	ChannelOrderRGBA      ChannelOrder = C.CL_RGBA
	ChannelOrderBGRA      ChannelOrder = C.CL_BGRA
	ChannelOrderARGB      ChannelOrder = C.CL_ARGB
	ChannelOrderIntensity ChannelOrder = C.CL_INTENSITY
	ChannelOrderLuminance ChannelOrder = C.CL_LUMINANCE
)

var channelOrderNameMap = map[ChannelOrder]string{
	ChannelOrderR:         "R",
	ChannelOrderA:         "A",
	ChannelOrderRG:        "RG",
	ChannelOrderRA:        "RA",
	ChannelOrderRGB:       "RGB",
	ChannelOrderRGBA:      "RGBA",
	ChannelOrderBGRA:      "BGRA",
	ChannelOrderARGB:      "ARGB",
	ChannelOrderIntensity: "Intensity",
	ChannelOrderLuminance: "Luminance",
}

func (co ChannelOrder) String() string {
	name := channelOrderNameMap[co]
	if name == "" {
		name = fmt.Sprintf("Unknown(%x)", int(co))
	}
	return name
}

type ChannelDataType int

const (
	ChannelDataTypeSNormInt8      ChannelDataType = C.CL_SNORM_INT8
	ChannelDataTypeSNormInt16     ChannelDataType = C.CL_SNORM_INT16
	ChannelDataTypeUNormInt8      ChannelDataType = C.CL_UNORM_INT8
	ChannelDataTypeUNormInt16     ChannelDataType = C.CL_UNORM_INT16
	ChannelDataTypeUNormShort565  ChannelDataType = C.CL_UNORM_SHORT_565
	ChannelDataTypeUNormShort555  ChannelDataType = C.CL_UNORM_SHORT_555
	ChannelDataTypeUNormInt101010 ChannelDataType = C.CL_UNORM_INT_101010
	ChannelDataTypeSignedInt8     ChannelDataType = C.CL_SIGNED_INT8
	ChannelDataTypeSignedInt16    ChannelDataType = C.CL_SIGNED_INT16
	ChannelDataTypeSignedInt32    ChannelDataType = C.CL_SIGNED_INT32
	ChannelDataTypeUnsignedInt8   ChannelDataType = C.CL_UNSIGNED_INT8
	ChannelDataTypeUnsignedInt16  ChannelDataType = C.CL_UNSIGNED_INT16
	ChannelDataTypeUnsignedInt32  ChannelDataType = C.CL_UNSIGNED_INT32
	ChannelDataTypeHalfFloat      ChannelDataType = C.CL_HALF_FLOAT
	ChannelDataTypeFloat          ChannelDataType = C.CL_FLOAT
)

var channelDataTypeNameMap = map[ChannelDataType]string{
	ChannelDataTypeSNormInt8:      "SNormInt8",
	ChannelDataTypeSNormInt16:     "SNormInt16",
	ChannelDataTypeUNormInt8:      "UNormInt8",
	ChannelDataTypeUNormInt16:     "UNormInt16",
	ChannelDataTypeUNormShort565:  "UNormShort565",
	ChannelDataTypeUNormShort555:  "UNormShort555",
	ChannelDataTypeUNormInt101010: "UNormInt101010",
	ChannelDataTypeSignedInt8:     "SignedInt8",
	ChannelDataTypeSignedInt16:    "SignedInt16",
	ChannelDataTypeSignedInt32:    "SignedInt32",
	ChannelDataTypeUnsignedInt8:   "UnsignedInt8",
	ChannelDataTypeUnsignedInt16:  "UnsignedInt16",
	ChannelDataTypeUnsignedInt32:  "UnsignedInt32",
	ChannelDataTypeHalfFloat:      "HalfFloat",
	ChannelDataTypeFloat:          "Float",
}

func (ct ChannelDataType) String() string {
	name := channelDataTypeNameMap[ct]
	if name == "" {
		name = fmt.Sprintf("Unknown(%x)", int(ct))
	}
	return name
}

type ImageInfoParam int

const (
	ImageInfoFormat 	ImageInfoParam = C.CL_IMAGE_FORMAT
	ImageInfoElementSize 	ImageInfoParam = C.CL_IMAGE_ELEMENT_SIZE
	ImageInfoRowPitch 	ImageInfoParam = C.CL_IMAGE_ROW_PITCH
	ImageInfoSlicePitch 	ImageInfoParam = C.CL_IMAGE_SLICE_PITCH
	ImageInfoHeight	 	ImageInfoParam = C.CL_IMAGE_HEIGHT
	ImageInfoWidth 		ImageInfoParam = C.CL_IMAGE_WIDTH
	ImageInfoDepth 		ImageInfoParam = C.CL_IMAGE_DEPTH
)

const (
        MemObjectTypeImage2D MemObjectType = C.CL_MEM_OBJECT_IMAGE2D
        MemObjectTypeImage3D MemObjectType = C.CL_MEM_OBJECT_IMAGE3D
)

//////////////// Abstract Types ////////////////
type ImageFormat struct {
	ChannelOrder    ChannelOrder
	ChannelDataType ChannelDataType
}

func (f ImageFormat) toCl() C.cl_image_format {
	var format C.cl_image_format
	format.image_channel_order = C.cl_channel_order(f.ChannelOrder)
	format.image_channel_data_type = C.cl_channel_type(f.ChannelDataType)
	return format
}

func (ip ImageInfoParam) toCl() C.cl_image_info {
	return C.cl_image_info(ip)
}

type ImageDescription struct {
        Type                            MemObjectType
        Width, Height, Depth            int
        ArraySize, RowPitch, SlicePitch int
        NumMipLevels, NumSamples        int
        Buffer                          *MemObject
}

//////////////// Basic Functions ////////////////
func getImageInfoInt(memobj *MemObject, param_name ImageInfoParam) (int, error) {
	var val C.size_t
	err := C.clGetImageInfo(memobj.clMem, param_name.toCl(), C.size_t(unsafe.Sizeof(val)), unsafe.Pointer(&val), nil)
	if toError(err) != nil {
		return -1, toError(err)
	}
	return int(val), toError(err)
}

//////////////// Abstract Functions ////////////////
func (ctx *Context) CreateImage2D(flags MemFlag, imageFormat ImageFormat, imageDesc ImageDescription, data []byte) (*MemObject, error) {
	format := imageFormat.toCl()
	var dataPtr unsafe.Pointer
	if data != nil {
		dataPtr = unsafe.Pointer(&data[0])
	}
	var err C.cl_int
	clBuffer := C.clCreateImage2D(ctx.clContext, C.cl_mem_flags(flags), &format, C.size_t(imageDesc.Width), C.size_t(imageDesc.Height), C.size_t(imageDesc.RowPitch), dataPtr, &err)
	if err != C.CL_SUCCESS {
		return nil, toError(err)
	}
	if clBuffer == nil {
		return nil, ErrUnknown
	}
	return newMemObject(clBuffer, len(data)), nil
}

func (ctx *Context) CreateImage3D(flags MemFlag, imageFormat ImageFormat, imageDesc ImageDescription, data []byte) (*MemObject, error) {
	format := imageFormat.toCl()
	var dataPtr unsafe.Pointer
	if data != nil {
		dataPtr = unsafe.Pointer(&data[0])
	}
	var err C.cl_int
	clBuffer := C.clCreateImage3D(ctx.clContext, C.cl_mem_flags(flags), &format, C.size_t(imageDesc.Width), C.size_t(imageDesc.Height), C.size_t(imageDesc.Depth),
				      C.size_t(imageDesc.RowPitch), C.size_t(imageDesc.SlicePitch), dataPtr, &err)
	if err != C.CL_SUCCESS {
		return nil, toError(err)
	}
	if clBuffer == nil {
		return nil, ErrUnknown
	}
	return newMemObject(clBuffer, len(data)), nil
}

func (ctx *Context) GetSupportedImageFormats(flags MemFlag, imageType MemObjectType) ([]ImageFormat, error) {
        var formats [maxImageFormats]C.cl_image_format
        var nFormats C.cl_uint
        if err := C.clGetSupportedImageFormats(ctx.clContext, C.cl_mem_flags(flags), C.cl_mem_object_type(imageType), maxImageFormats, &formats[0], &nFormats); err != C.CL_SUCCESS {
                return nil, toError(err)
        }
        fmts := make([]ImageFormat, nFormats)
        for i, f := range formats[:nFormats] {
                fmts[i] = ImageFormat{
                        ChannelOrder:    ChannelOrder(f.image_channel_order),
                        ChannelDataType: ChannelDataType(f.image_channel_data_type),
                }
        }
        return fmts, nil
}

// Enqueues a command to map a region of an image object into the host address space and returns a pointer to this mapped region.
func (q *CommandQueue) EnqueueMapImage(buffer *MemObject, blocking bool, flags MapFlag, origin, region [3]int, eventWaitList []*Event) (*MappedMemObject, *Event, error) {
	cOrigin := sizeT3(origin)
	cRegion := sizeT3(region)
	var event C.cl_event
	var err C.cl_int
	var rowPitch, slicePitch C.size_t
	ptr := C.clEnqueueMapImage(q.clQueue, buffer.clMem, clBool(blocking), flags.toCl(), &cOrigin[0], &cRegion[0], &rowPitch, &slicePitch, C.cl_uint(len(eventWaitList)), eventListPtr(eventWaitList), &event, &err)
	if err != C.CL_SUCCESS {
		return nil, nil, toError(err)
	}
	ev := newEvent(event)
	if ptr == nil {
		return nil, ev, ErrUnknown
	}
	size := 0 // TODO: could calculate this
	return &MappedMemObject{ptr: ptr, size: size, rowPitch: int(rowPitch), slicePitch: int(slicePitch)}, ev, nil
}

// Enqueues a command to read from a 2D or 3D image object to host memory.
func (q *CommandQueue) EnqueueReadImage(image *MemObject, blocking bool, origin, region [3]int, rowPitch, slicePitch int, data []byte, eventWaitList []*Event) (*Event, error) {
	cOrigin := sizeT3(origin)
	cRegion := sizeT3(region)
	var event C.cl_event
	err := toError(C.clEnqueueReadImage(q.clQueue, image.clMem, clBool(blocking), &cOrigin[0], &cRegion[0], C.size_t(rowPitch), C.size_t(slicePitch), unsafe.Pointer(&data[0]), C.cl_uint(len(eventWaitList)), eventListPtr(eventWaitList), &event))
	return newEvent(event), err
}

// Enqueues a command to write from a 2D or 3D image object to host memory.
func (q *CommandQueue) EnqueueWriteImage(image *MemObject, blocking bool, origin, region [3]int, rowPitch, slicePitch int, data []byte, eventWaitList []*Event) (*Event, error) {
	cOrigin := sizeT3(origin)
	cRegion := sizeT3(region)
	var event C.cl_event
	err := toError(C.clEnqueueWriteImage(q.clQueue, image.clMem, clBool(blocking), &cOrigin[0], &cRegion[0], C.size_t(rowPitch), C.size_t(slicePitch), unsafe.Pointer(&data[0]), C.cl_uint(len(eventWaitList)), eventListPtr(eventWaitList), &event))
	return newEvent(event), err
}

// Enqueues a command to copy from a 2D or 3D image object to device memory as image.
func (q *CommandQueue) EnqueueCopyImage(dst, src *MemObject, dst_origin, src_origin, region [3]int, eventWaitList []*Event) (*Event, error) {
	dOrigin := sizeT3(dst_origin)
	sOrigin := sizeT3(src_origin)
	cRegion := sizeT3(region)
	var event C.cl_event
	err := toError(C.clEnqueueCopyImage(q.clQueue, src.clMem, dst.clMem, &sOrigin[0], &dOrigin[0], &cRegion[0], C.cl_uint(len(eventWaitList)), eventListPtr(eventWaitList), &event))
	return newEvent(event), err
}

// Enqueues a command to copy from a 2D or 3D image object to buffer memory.
func (q *CommandQueue) EnqueueCopyImageToBuffer(dst, src *MemObject, src_origin, region [3]int, dst_offset int, eventWaitList []*Event) (*Event, error) {
	sOrigin := sizeT3(src_origin)
	cRegion := sizeT3(region)
	var event C.cl_event
	err := toError(C.clEnqueueCopyImageToBuffer(q.clQueue, src.clMem, dst.clMem, &sOrigin[0], &cRegion[0], C.size_t(dst_offset), C.cl_uint(len(eventWaitList)), eventListPtr(eventWaitList), &event))
	return newEvent(event), err
}

// Enqueues a command to copy from a 2D or 3D image object to buffer memory.
func (q *CommandQueue) EnqueueCopyBufferToImage(dst, src *MemObject, src_offset int, region, dst_origin [3]int, eventWaitList []*Event) (*Event, error) {
	dOrigin := sizeT3(dst_origin)
	cRegion := sizeT3(region)
	var event C.cl_event
	err := toError(C.clEnqueueCopyBufferToImage(q.clQueue, src.clMem, dst.clMem, (C.size_t)(src_offset), &dOrigin[0], &cRegion[0], C.cl_uint(len(eventWaitList)), eventListPtr(eventWaitList), &event))
	return newEvent(event), err
}

func (image_desc *ImageDescription) GetFormat() (*ImageFormat, error) {
	var tmpFormat C.cl_image_format
	err := C.clGetImageInfo(image_desc.Buffer.clMem, ImageInfoFormat.toCl(), C.size_t(unsafe.Sizeof(tmpFormat)), (unsafe.Pointer)(&tmpFormat), nil)
	if toError(err) != nil {
		return nil, toError(err)
	}
	return &ImageFormat{ChannelOrder: (ChannelOrder)(tmpFormat.image_channel_order), ChannelDataType: (ChannelDataType)(tmpFormat.image_channel_data_type)}, toError(err)
}

func (image_desc *ImageDescription) GetElementSize() (int, error) {
        return getImageInfoInt(image_desc.Buffer, ImageInfoElementSize)
}

func (image_desc *ImageDescription) GetRowPitch() (int, error) {
        return getImageInfoInt(image_desc.Buffer, ImageInfoRowPitch)
}

func (image_desc *ImageDescription) GetSlicePitch() (int, error) {
        return getImageInfoInt(image_desc.Buffer, ImageInfoSlicePitch)
}

func (image_desc *ImageDescription) GetHeight() (int, error) {
        return getImageInfoInt(image_desc.Buffer, ImageInfoHeight)
}

func (image_desc *ImageDescription) GetWidth() (int, error) {
        return getImageInfoInt(image_desc.Buffer, ImageInfoWidth)
}

func (image_desc *ImageDescription) GetDepth() (int, error) {
        return getImageInfoInt(image_desc.Buffer, ImageInfoDepth)
}

func (image_desc *ImageDescription) GetContext() (*Context, error) {
	if image_desc.Buffer != nil {
		return image_desc.Buffer.GetContext()
	}
        return nil, toError(C.CL_INVALID_MEM_OBJECT)
}

func (image_desc *ImageDescription) GetMemOffset() (int, error) {
        if image_desc.Buffer != nil {
                return image_desc.Buffer.GetOffset()
        }
        return 0, toError(C.CL_INVALID_MEM_OBJECT)
}

func (image_desc *ImageDescription) GetAssociatedMemObject() (*MemObject, error) {
        if image_desc.Buffer != nil {
                return image_desc.Buffer.GetAssociatedMemObject()
        }
        return nil, toError(C.CL_INVALID_MEM_OBJECT)
}

