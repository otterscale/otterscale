<script lang="ts">
	import { getContext } from 'svelte';
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import * as AlertDialog from '$lib/components/ui/alert-dialog/index.js';
	import { Button } from '$lib/components/ui/button';
	import {
		Nexus,
		type Configuration,
		type Configuration_BootImage,
		type SetDefaultBootImageRequest
	} from '$gen/api/nexus/v1/nexus_pb';
	import { toast } from 'svelte-sonner';

	let {
		bootImage,
		configuration = $bindable()
	}: {
		bootImage: Configuration_BootImage;
		configuration: Configuration;
	} = $props();

	const transport: Transport = getContext('transportNEW');
	const client = createClient(Nexus, transport);

	const DEFAULT_REQUEST = {
		distroSeries: bootImage.distroSeries
	} as SetDefaultBootImageRequest;
	let setDefaultBootImageRequest = $state(DEFAULT_REQUEST);

	function reset() {
		setDefaultBootImageRequest = DEFAULT_REQUEST;
	}

	let open = $state(false);
	function close() {
		open = false;
	}
</script>

<AlertDialog.Root bind:open>
	<AlertDialog.Trigger>
		<Button variant="ghost" class="flex items-center gap-1">
			<Icon icon="ph:star" class="h-4 w-4" />
			Default
		</Button>
	</AlertDialog.Trigger>
	<AlertDialog.Content>
		<AlertDialog.Header>
			<AlertDialog.Title>Set Default Boot Image</AlertDialog.Title>
			<AlertDialog.Description class="rounded-lg bg-muted/50 p-4">
				Are you sure you want to set this as the default boot image? This action will change the
				system's default boot configuration.
			</AlertDialog.Description>
		</AlertDialog.Header>
		<AlertDialog.Footer>
			<AlertDialog.Cancel onclick={reset} class="mr-auto">Cancel</AlertDialog.Cancel>
			<AlertDialog.Action
				onclick={() => {
					client
						.setDefaultBootImage(setDefaultBootImageRequest)
						.then((r) => {
							toast.info(`Set ${setDefaultBootImageRequest.distroSeries} as default.`);
							client.getConfiguration({}).then((r) => {
								configuration = r;
							});
						})
						.catch((e) => {
							toast.info(`Fail to set ${setDefaultBootImageRequest.distroSeries} as default.`);
						});
					console.log(setDefaultBootImageRequest);
					reset();
					close();
				}}
			>
				Continue
			</AlertDialog.Action>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>
