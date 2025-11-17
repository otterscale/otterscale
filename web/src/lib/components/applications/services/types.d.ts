import { type Application_Service } from '$lib/api/application/v1/application_pb';

interface Service extends Application_Service {
	endpoint: string;
}

export type { Service };
