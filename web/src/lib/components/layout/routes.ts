import { dynamicPaths, type Path } from '$lib/path';

interface Route {
	path: Path;
	items: Path[];
}

const globalRoutes = (scope: string | undefined): Route[] => [
	{
		path: dynamicPaths.networking(scope),
		items: [dynamicPaths.networkingSubnets(scope)],
	},
	{
		path: dynamicPaths.machines(scope),
		items: [dynamicPaths.machinesMetal(scope)],
	},
	{
		path: dynamicPaths.globalSettings(scope),
		items: [],
	},
];

const platformRoutes = (scope: string | undefined): Route[] => [
	{
		path: dynamicPaths.models(scope),
		items: [dynamicPaths.modelsLLM(scope)],
	},
	{
		path: dynamicPaths.applications(scope),
		items: [
			dynamicPaths.applicationsWorkloads(scope),
			dynamicPaths.applicationsServices(scope),
			dynamicPaths.applicationsStore(scope),
		],
	},
	{
		path: dynamicPaths.compute(scope),
		items: [dynamicPaths.computeVirtualMachine(scope)],
	},
	{
		path: dynamicPaths.storage(scope),
		items: [
			dynamicPaths.storageOSD(scope),
			dynamicPaths.storagePool(scope),
			dynamicPaths.storageBlockDevice(scope),
			dynamicPaths.storageFileSystem(scope),
			dynamicPaths.storageObjectGateway(scope),
		],
	},

	{
		path: dynamicPaths.scopeBasedSettings(scope),
		items: [],
	},
];

export type { Route };
export { globalRoutes, platformRoutes };
