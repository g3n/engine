#ifdef MORPHTARGETS
	vPosition += (MorphPosition{i} - VertexPosition) * morphTargetInfluences[{i}];
  #ifdef MORPHTARGETS_NORMAL
	vNormal += (MorphNormal{i} - VertexNormal) * morphTargetInfluences[{i}];
  #endif
#endif
