#ifdef MORPHTARGETS
	vPosition += (MorphPosition{i} - VertexPosition) * morphTargetInfluences[{i}];
//    vPosition = MorphPosition1;
  #ifdef MORPHTARGETS_NORMAL
	vNormal += (MorphNormal{i} - VertexNormal) * morphTargetInfluences[{i}];
  #endif
#endif
