	vPosition += MorphPosition{i} * morphTargetInfluences[{i}];
  #ifdef MORPHTARGETS_NORMAL
	vNormal += MorphNormal{i} * morphTargetInfluences[{i}];
  #endif