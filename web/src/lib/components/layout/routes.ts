import { resolve } from '$app/paths';
import { m } from '$lib/paraglide/messages';
import type { Path } from '$lib/path';

interface Route {
	path: Path;
	items: Path[];
}

const globalRoutes = (): Route[] => [
	{
		path: {
			title: m.networking(),
			url: resolve('/(auth)/networking')
		},
		items: [{ title: m.subnets(), url: resolve('/(auth)/networking/subnets') }]
	},
	{
		path: {
			title: m.machines(),
			url: resolve('/(auth)/machines')
		},
		items: [{ title: m.metal(), url: resolve('/(auth)/machines/metal') }]
	},
	{
		path: {
			title: m.configuration(),
			url: resolve('/(auth)/configuration')
		},
		items: []
	}
];

const platformRoutes = (scope: string): Route[] => [
	{
		path: { title: m.models(), url: resolve('/(auth)/scope/[scope]/models', { scope: scope }) },
		items: [{ title: m.llm(), url: resolve('/(auth)/scope/[scope]/models/llm', { scope: scope }) }]
	},
	{
		path: {
			title: m.applications(),
			url: resolve('/(auth)/scope/[scope]/applications', { scope: scope })
		},
		items: [
			{
				title: m.workloads(),
				url: resolve('/(auth)/scope/[scope]/applications/workloads', { scope: scope })
			},
			{
				title: m.services(),
				url: resolve('/(auth)/scope/[scope]/applications/services', { scope: scope })
			},
			{
				title: m.store(),
				url: resolve('/(auth)/scope/[scope]/applications/store', { scope: scope })
			}
		]
	},
	{
		path: { title: m.compute(), url: resolve('/(auth)/scope/[scope]/compute', { scope: scope }) },
		items: [
			{
				title: m.virtual_machine(),
				url: resolve('/(auth)/scope/[scope]/compute/virtual-machine', { scope: scope })
			}
		]
	},
	{
		path: { title: m.storage(), url: resolve('/(auth)/scope/[scope]/storage', { scope: scope }) },
		items: [
			{ title: m.osd(), url: resolve('/(auth)/scope/[scope]/storage/osd', { scope: scope }) },
			{ title: m.pool(), url: resolve('/(auth)/scope/[scope]/storage/pool', { scope: scope }) },
			{
				title: m.block_device(),
				url: resolve('/(auth)/scope/[scope]/storage/block-device', { scope: scope })
			},
			{
				title: m.file_system(),
				url: resolve('/(auth)/scope/[scope]/storage/file-system', { scope: scope })
			},
			{
				title: m.smb(),
				url: resolve('/(auth)/scope/[scope]/storage/smb', { scope: scope })
			},
			{
				title: m.object_gateway(),
				url: resolve('/(auth)/scope/[scope]/storage/object-gateway', { scope: scope })
			}
		]
	},

	{
		path: { title: m.settings(), url: resolve('/(auth)/scope/[scope]/settings', { scope: scope }) },
		items: []
	}
];

export type { Route };
export { globalRoutes, platformRoutes };
