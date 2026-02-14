import BotIcon from '@lucide/svelte/icons/bot';
import CloudBackupIcon from '@lucide/svelte/icons/cloud-backup';
import CodeIcon from '@lucide/svelte/icons/code';
import CpuIcon from '@lucide/svelte/icons/cpu';
import DatabaseIcon from '@lucide/svelte/icons/database';
import DumbbellIcon from '@lucide/svelte/icons/dumbbell';
import FlagIcon from '@lucide/svelte/icons/flag';
import GaugeIcon from '@lucide/svelte/icons/gauge';
import HardDriveIcon from '@lucide/svelte/icons/hard-drive';
import LayoutGridIcon from '@lucide/svelte/icons/layout-grid';
import MapIcon from '@lucide/svelte/icons/map';
import NetworkIcon from '@lucide/svelte/icons/network';
import PcCaseIcon from '@lucide/svelte/icons/pc-case';
import ScaleIcon from '@lucide/svelte/icons/scale';
import ShieldCheckIcon from '@lucide/svelte/icons/shield-check';
import ShipIcon from '@lucide/svelte/icons/ship';
import TelescopeIcon from '@lucide/svelte/icons/telescope';
import TerminalIcon from '@lucide/svelte/icons/terminal';
import UnplugIcon from '@lucide/svelte/icons/unplug';
import UsersIcon from '@lucide/svelte/icons/users';
import WorkflowIcon from '@lucide/svelte/icons/workflow';

import { m } from '$lib/paraglide/messages';

export const navData = {
	overview: [
		{
			name: m.workspace(),
			url: '#',
			icon: MapIcon,
			edit: true
		}
	],
	aiStudio: [
		{
			title: 'Inference',
			url: '#',
			icon: BotIcon,
			items: [
				{
					title: 'Model',
					url: '#',
					disabled: true
				}
			]
		},
		{
			title: 'Training',
			url: '#',
			icon: DumbbellIcon,
			items: [
				{
					title: 'Finetune Job',
					url: '#',
					disabled: true
				},
				{
					title: 'Dataset',
					url: '#',
					disabled: true
				}
			]
		},
		{
			title: 'Notebooks',
			url: '#',
			icon: TerminalIcon,
			items: [
				{
					title: 'Jupyter',
					url: '#',
					disabled: true
				}
			]
		}
	],
	applications: [
		{
			title: 'Hub',
			url: '#',
			icon: LayoutGridIcon,
			items: [
				{
					title: 'Release',
					url: '#',
					disabled: true
				},
				{
					title: 'Chart',
					url: '#',
					disabled: true
				}
			]
		},
		{
			title: 'Cloud IDE',
			url: '#',
			icon: CodeIcon,
			items: [
				{
					title: 'Coder',
					url: '#',
					disabled: true
				}
			]
		},
		{
			title: 'Database',
			url: '#',
			icon: DatabaseIcon,
			items: [
				{
					title: 'Postgres',
					url: '#',
					disabled: true
				},
				{
					title: 'Redis',
					url: '#',
					disabled: true
				}
			]
		},
		{
			title: 'Workflow',
			url: '#',
			icon: WorkflowIcon,
			items: [
				{
					title: 'Pipeline',
					url: '#',
					disabled: true
				},
				{
					title: 'Task',
					url: '#',
					disabled: true
				}
			]
		}
	],
	resources: [
		{
			title: 'Workloads',
			url: '#',
			icon: FlagIcon,
			items: [
				{
					title: 'Deployment',
					url: '#',
					disabled: true
				},
				{
					title: 'Stateful Set',
					url: '#',
					disabled: true
				},
				{
					title: 'Daemon Set',
					url: '#',
					disabled: true
				},
				{
					title: 'Cron Job',
					url: '#',
					disabled: true
				},
				{
					title: 'Job',
					url: '#',
					disabled: true
				},
				{
					title: 'Pod',
					url: '#',
					disabled: true
				}
			]
		},
		{
			title: 'Compute',
			url: '#',
			icon: CpuIcon,
			items: [
				{
					title: 'Virtual Machine',
					url: '#',
					disabled: true
				}
			]
		},
		{
			title: 'Network',
			url: '#',
			icon: NetworkIcon,
			items: [
				{
					title: 'VPC',
					url: '#',
					disabled: true
				},
				{
					title: 'Load Balancer',
					url: '#',
					disabled: true
				}
			]
		},
		{
			title: 'Storage',
			url: '#',
			icon: HardDriveIcon,
			items: [
				{
					title: 'Block Pool',
					url: '#',
					disabled: true
				},
				{
					title: 'File System',
					url: '#',
					disabled: true
				}
			]
		}
	],
	governance: [
		{
			title: 'Tenant',
			url: '#',
			icon: UsersIcon,
			items: [
				{
					title: 'Workspace',
					url: '#',
					disabled: true
				},
				{
					title: 'User',
					url: '#',
					disabled: true
				}
			]
		},
		{
			title: 'Policy',
			url: '#',
			icon: ScaleIcon,
			items: [
				{
					title: 'Policy',
					url: '#',
					disabled: true
				},
				{
					title: 'Compliance',
					url: '#',
					disabled: true
				}
			]
		},
		{
			title: 'Metering',
			url: '#',
			icon: GaugeIcon,
			items: [
				{
					title: 'Budget',
					url: '#',
					disabled: true
				},
				{
					title: 'Usage',
					url: '#',
					disabled: true
				}
			]
		},
		{
			title: 'Audit',
			url: '#',
			icon: ShieldCheckIcon,
			items: [
				{
					title: 'Trail',
					url: '#',
					disabled: true
				},
				{
					title: 'Log',
					url: '#',
					disabled: true
				}
			]
		}
	],
	reliability: [
		{
			title: 'Telemetry',
			url: '#',
			icon: TelescopeIcon,
			items: [
				{
					title: 'Collector',
					url: '#',
					disabled: true
				},
				{
					title: 'Rule',
					url: '#',
					disabled: true
				}
			]
		},
		{
			title: 'Recovery',
			url: '#',
			icon: CloudBackupIcon,
			items: [
				{
					title: 'Backup',
					url: '#',
					disabled: true
				},
				{
					title: 'Restore',
					url: '#',
					disabled: true
				}
			]
		}
	],
	system: [
		{
			title: 'Fleet',
			url: '#',
			icon: ShipIcon,
			items: [
				{
					title: 'Cluster',
					url: '#',
					disabled: true
				},
				{
					title: 'Config',
					url: '#',
					disabled: true
				},
				{
					title: 'Image',
					url: '#',
					disabled: true
				}
			]
		},
		{
			title: 'Metal',
			url: '#',
			icon: PcCaseIcon,
			items: [
				{
					title: 'Server',
					url: '#',
					disabled: true
				}
			]
		},
		{
			title: 'Tunnels',
			url: '#',
			icon: UnplugIcon,
			items: [
				{
					title: 'Server',
					url: '#',
					disabled: true
				},
				{
					title: 'Client',
					url: '#',
					disabled: true
				}
			]
		}
	]
};
