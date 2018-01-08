#ifndef _AL_LISTENER_H_
#define _AL_LISTENER_H_

#include "alMain.h"
#include "alu.h"

#ifdef __cplusplus
extern "C" {
#endif

struct ALlistenerProps {
    ALfloat Position[3];
    ALfloat Velocity[3];
    ALfloat Forward[3];
    ALfloat Up[3];
    ALfloat Gain;
    ALfloat MetersPerUnit;

    ALfloat DopplerFactor;
    ALfloat DopplerVelocity;
    ALfloat SpeedOfSound;
    ALboolean SourceDistanceModel;
    enum DistanceModel DistanceModel;

    ATOMIC(struct ALlistenerProps*) next;
};

typedef struct ALlistener {
    alignas(16) ALfloat Position[3];
    ALfloat Velocity[3];
    ALfloat Forward[3];
    ALfloat Up[3];
    ALfloat Gain;
    ALfloat MetersPerUnit;

    /* Pointer to the most recent property values that are awaiting an update.
     */
    ATOMIC(struct ALlistenerProps*) Update;

    /* A linked list of unused property containers, free to use for future
     * updates.
     */
    ATOMIC(struct ALlistenerProps*) FreeList;

    struct {
        aluMatrixf Matrix;
        aluVector  Velocity;

        ALfloat Gain;
        ALfloat MetersPerUnit;

        ALfloat DopplerFactor;
        ALfloat SpeedOfSound;

        ALboolean SourceDistanceModel;
        enum DistanceModel DistanceModel;
    } Params;
} ALlistener;

void UpdateListenerProps(ALCcontext *context);

#ifdef __cplusplus
}
#endif

#endif
