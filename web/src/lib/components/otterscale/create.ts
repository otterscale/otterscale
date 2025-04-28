import type { Configuration } from '$lib/components/otterscale/interfaces';

export let machines: Configuration[] = [
    {
        key: 'bear_metal',
        name: 'Bear Metal',
        icon: 'ph:file-ini',
        steps: [
            {
                description: 'Network',
                parameters: [
                    {
                        key: 'dhcp_on',
                        name: 'DHCP On',
                        value: '',
                        description: ''
                    },
                    {
                        key: 'cidr',
                        name: 'CIDR',
                        value: '',
                        description: ''
                    },
                    {
                        key: 'gateway_ip',
                        name: 'Gateway Ip',
                        value: '',
                        description: ''
                    },
                    {
                        key: 'dns_servers',
                        name: 'DNS Servers',
                        values: [],
                        description: ''
                    },
                    {
                        key: 'start_ip',
                        name: 'Start IP',
                        value: '',
                        description: ''
                    },
                    {
                        key: 'end_ip',
                        name: 'End IP',
                        value: '',
                        description: ''
                    }
                ]
            },
            {
                description: 'Commission Machine',
                parameters: [
                    {
                        key: 'system_id',
                        name: 'System ID ',
                        value: '',
                        description: ''
                    }
                ],
                advancedParameters: [
                    {
                        key: 'enable_ssh',
                        name: 'Enable SSH',
                        value: '',
                        description: ''
                    },
                    {
                        key: 'skip_bmc_config',
                        name: 'Skip BMC Config',
                        value: '',
                        description: ''
                    },
                    {
                        key: 'skip_networking',
                        name: 'Skip Networking',
                        value: '',
                        description: ''
                    },
                    {
                        key: 'skip_storage',
                        name: 'Skip Storage',
                        value: '',
                        description: ''
                    }
                ]
            }
        ],
        templates: [
            {
                name: 'Basic',
                parameters: [
                    {
                        key: 'dhcp_on',
                        value: 'On'
                    }
                ]
            },
            {
                name: 'Advance',
                parameters: [
                    {
                        key: 'enable_ssh',
                        value: 'True'
                    }
                ]
            }
        ]
    }
];