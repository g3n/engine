#ifdef BONE_INFLUENCERS
    #if BONE_INFLUENCERS > 0
	uniform mat4 mBones[TOTAL_BONES];
    in vec4 matricesIndices;
    in vec4 matricesWeights;
//    #if BONE_INFLUENCERS > 4
//        in vec4 matricesIndicesExtra;
//        in vec4 matricesWeightsExtra;
//    #endif
    #endif
#endif
