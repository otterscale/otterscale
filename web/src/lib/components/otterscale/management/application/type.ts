import type { Facility, Scope } from '$gen/api/nexus/v1/nexus_pb';

export interface Kubernetes {
    scope: Scope;
    facility: Facility;
}