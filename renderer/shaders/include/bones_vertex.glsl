#ifdef BONE_INFLUENCERS
    #if BONE_INFLUENCERS > 0

        mat4 influence = mBones[int(matricesIndices[0])] * matricesWeights[0];
        #if BONE_INFLUENCERS > 1
            influence += mBones[int(matricesIndices[1])] * matricesWeights[1];
            #if BONE_INFLUENCERS > 2
                influence += mBones[int(matricesIndices[2])] * matricesWeights[2];
                #if BONE_INFLUENCERS > 3
                    influence += mBones[int(matricesIndices[3])] * matricesWeights[3];
    //                #if BONE_INFLUENCERS > 4
    //                    influence += mBones[int(matricesIndicesExtra[0])] * matricesWeightsExtra[0];
    //                    #if BONE_INFLUENCERS > 5
    //                        influence += mBones[int(matricesIndicesExtra[1])] * matricesWeightsExtra[1];
    //                        #if BONE_INFLUENCERS > 6
    //                            influence += mBones[int(matricesIndicesExtra[2])] * matricesWeightsExtra[2];
    //                            #if BONE_INFLUENCERS > 7
    //                                influence += mBones[int(matricesIndicesExtra[3])] * matricesWeightsExtra[3];
    //                            #endif
    //                        #endif
    //                    #endif
    //                #endif
                #endif
            #endif
        #endif

        finalWorld = finalWorld * influence;

    #endif
#endif
