/*
  This file is used to point the compiler to the actual opencl.h of the system.
  It is also used to check the version of opencl installed
*/
#include <stdlib.h>

#ifdef __APPLE__
	#include <OpenCL/OpenCL.h>
#else
	#include <CL/opencl.h>
#ifdef __WIN32
	#include <CL/cl_d3d10.h>
#endif
#endif

#ifndef CL_VERSION_1_0
	#error "This package requires OpenCL 1.0"
#endif

