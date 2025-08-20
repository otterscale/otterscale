import { dynamicPaths, type Path } from '$lib/path';

export interface Route {
	path: Path;
	items: Path[];
}

export const platformRoutes = (scope: string): Route[] => [
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
			dynamicPaths.applicationsWorkload(scope),
			dynamicPaths.applicationsService(scope),
			dynamicPaths.applicationsStore(scope)
		]
	},
	{
		path: dynamicPaths.virtualization(scope),
		items: [dynamicPaths.virtualizationVirtualMachine(scope)]
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
	}
];

export const globalRoutes = (scope: string): Route[] => [
	{
		path: dynamicPaths.machines(scope),
		items: [dynamicPaths.machinesMetal(scope)]
	},
	{
		path: dynamicPaths.settings(scope),
		items: [
			dynamicPaths.settingsNetwork(scope),
			// dynamicPaths.settingsSubscription(scope),
			dynamicPaths.settingsBIST(scope)
		]
	}
];
