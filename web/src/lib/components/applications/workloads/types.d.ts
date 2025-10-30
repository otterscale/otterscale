import { type Application as Base } from '$lib/api/application/v1/application_pb';

interface Application extends Base {
	publicAddress: string;
}

export type { Application };
