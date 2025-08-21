import { dynamicPaths, type Path } from '$lib/path';

export interface Route {
	path: Path;
	items: Path[];
}

export const applicationRoutes = (scope: string | undefined): Route[] => [
	{
		path: dynamicPaths.models(scope),
		items: [dynamicPaths.modelsLLM(scope)]
	},
	{
		path: dynamicPaths.databases(scope),
		items: [dynamicPaths.databasesRelational(scope), dynamicPaths.databasesNoSQL(scope)]
	},
	{
		path: dynamicPaths.applications(scope),
		items: [
			dynamicPaths.applicationsDeployments(scope),
			dynamicPaths.applicationsServices(scope),
			dynamicPaths.applicationsStore(scope)
		]
	}
];

export const platformRoutes = (scope: string | undefined): Route[] => [
	{
		path: dynamicPaths.compute(scope),
		items: [dynamicPaths.computeVirtualMachine(scope)]
	},
	{
		path: dynamicPaths.storage(scope),
		items: [
			dynamicPaths.storageOSD(scope),
			dynamicPaths.storagePool(scope),
			dynamicPaths.storageBlockDevice(scope),
			dynamicPaths.storageFileSystem(scope),
			dynamicPaths.storageObjectGateway(scope)
		]
	},
	{
		path: dynamicPaths.networking(scope),
		items: [dynamicPaths.networkingSubnets(scope)]
	},
	{
		path: dynamicPaths.machines(scope),
		items: [dynamicPaths.machinesMetal(scope)]
	},
	{
		path: dynamicPaths.settings(scope),
		items: []
	}
];
