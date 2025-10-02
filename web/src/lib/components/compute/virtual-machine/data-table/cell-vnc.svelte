<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';

	import { env } from '$env/dynamic/public';
	import { VirtualMachineService, type VirtualMachine } from '$lib/api/virtual_machine/v1/virtual_machine_pb';
	import { buttonVariants } from '$lib/components/ui/button';
	import * as Sheet from '$lib/components/ui/sheet';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { currentKubernetes } from '$lib/stores';
</script>

<script lang="ts">
	let { virtualMachine }: { virtualMachine: VirtualMachine } = $props();

	const transport: Transport = getContext('transport');
	const virtualMachineClient = createClient(VirtualMachineService, transport);

	const url = new URL(env.PUBLIC_API_URL ?? '');
	const [host, port] = url.host.split(':');
	const defaultPort = url.protocol === 'https:' ? '443' : '80';

	let isMounted = $state(false);
	let vncUrl = $state('');

	async function getVncUrl() {
		try {
			const response = await virtualMachineClient.vNCInstance({
				scopeUuid: $currentKubernetes?.scopeUuid,
				facilityName: $currentKubernetes?.name,
				name: virtualMachine.name,
				namespace: virtualMachine.namespace,
			});
			isMounted = true;
			vncUrl = `/vnc/vnc.html?autoconnect=true&host=${host}&port=${port || defaultPort}&path=vnc/${response.sessionId}`;
		} catch (error) {
			console.error('Error fetching VNC URL:', error);
		}
	}
</script>

<Sheet.Root
	onOpenChange={async (open) => {
		if (open) await getVncUrl();
	}}
>
	<Sheet.Trigger>
		<Tooltip.Provider>
			<Tooltip.Root>
				<Tooltip.Trigger class={buttonVariants({ variant: 'outline' })}>
					<Icon icon="ph:monitor" />
				</Tooltip.Trigger>
				<Tooltip.Content>
					VNC to {virtualMachine.name}
				</Tooltip.Content>
			</Tooltip.Root>
		</Tooltip.Provider>
	</Sheet.Trigger>
	<Sheet.Content class="rounded-l-lg border-none bg-transparent sm:max-w-9/10">
		{#if isMounted}
			<iframe class="size-full rounded-l-lg" title="Remote VNC Console" src={vncUrl}>
				Your browser does not support iframes.
			</iframe>
		{/if}
	</Sheet.Content>
</Sheet.Root>
