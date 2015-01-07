// Copyright 2014 Benjamin Savs.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "../../cmd/ld/textflag.h"

// func Vector3Lerp(a, b Vector3, t float32) Vector3 
TEXT ·Vector3Lerp(SB),NOSPLIT,$0
//	MOVSS  t+24(FP), X0
//	MOVSS  $(1.0), X1
//	SUBSS  X0, X1
//	MOVSS  X0, ret+32(FP)
//	MOVSS  X1, ret+36(FP)
	MOVSS  aX+12(FP), X0
	MOVSS  aY+16(FP), X1
	MOVSS  aZ+20(FP), X2
	MOVSS  X0, ret+32(FP)
	MOVSS  X1, ret+36(FP)
	MOVSS  X2, ret+40(FP)
	RET
