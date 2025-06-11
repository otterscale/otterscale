<script lang="ts">
	import Unstruct from '$lib/components/custom/unstruct.svelte';
	const data = {
		level1: {
			temperature: 25.5,
			sensors: [
				{ type: 'temperature', id: 'temp_1', values: [45, 48, 50] },
				{ type: 'temperature', id: 'temp_2', values: [45, 48, 50] },
				{ type: 'humidity', id: 'humidity_1', values: [45, 48, 50] },
				{ type: 'pressure', id: 'pressure_1', values: [45, 48, 50] }
			],
			systemStatus: {
				rack1: {
					powerStatus: true,
					units: [
						{ name: 'storage_1', status: { online: true, mode: 'normal' } },
						{ name: 'server_1', status: { online: true, mode: 'normal' } },
						{ name: 'network_1', status: { online: true, mode: 'normal' } }
					],
					monitoring: {
						lastCheck: 'operational',
						details: {
							timestamp: new Date('2000-11-8'),
							operator: { id: 'OP123', clearance: 'admin' }
						}
					}
				},
				rack2: [
					{ id: 'SRV01', status: 'running', specs: { cores: 32, services: ['web', 'db'] } },
					{ id: 'SRV02', status: 'maintenance', specs: { cores: 64, services: ['cache', 'queue'] } }
				]
			}
		},
		level2: {
			operational: true,
			devices: [
				{ type: 'switch', id: 'switch_1', status: 'active' },
				{ type: 'firewall', id: 'firewall_1', status: 'active' },
				{ type: 'router', id: 'router_1', status: 'active' }
			],
			network: {
				connected: true,
				configuration: {
					topology: 'mesh',
					protocols: ['BGP', 'OSPF', 'VLAN'],
					security: {
						encryption: 'AES256',
						maxAttempts: 3,
						policies: {
							enabled: true,
							rules: ['IDS', 'IPS']
						}
					}
				}
			},
			metrics: {
				performance: [
					{ period: '2023-Q1', stats: { uptime: 99.9, incidents: 2 } },
					{ period: '2023-Q2', stats: { uptime: 99.95, incidents: 1 } }
				]
			}
		}
	};
</script>

<Unstruct {data} orientation="horizontal" />
