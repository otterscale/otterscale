import { type Scope } from '$gen/api/scope/v1/scope_pb';
import { type Facility } from '$gen/api/facility/v1/facility_pb';

export interface Kubernetes {
    scope: Scope;
    facility: Facility;
}