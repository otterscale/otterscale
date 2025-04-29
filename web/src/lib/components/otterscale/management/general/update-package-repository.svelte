<script lang="ts">
	import { getContext } from 'svelte';
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import * as AlertDialog from '$lib/components/ui/alert-dialog/index.js';
	import { Button } from '$lib/components/ui/button';
	import { Label } from '$lib/components/ui/label';
	import { Input } from '$lib/components/ui/input';
	import {
		Nexus,
		type Configuration_PackageRepository,
		type UpdatePackageRepositoryRequest
	} from '$gen/api/nexus/v1/nexus_pb';
	import { toast } from 'svelte-sonner';

	let {
		configuration = $bindable(),
		packageRepository
	}: {
		configuration: Configuration;
		packageRepository: Configuration_PackageRepository;
	} = $props();

	const transport: Transport = getContext('transportNEW');
	const client = createClient(Nexus, transport);

	const DEFAULT_REQUEST = {
		id: packageRepository.id,
		url: packageRepository.url,
		skipJuju: false
	} as UpdatePackageRepositoryRequest;

	let updatePackageRepositoryRequest = $state(DEFAULT_REQUEST);

	function reset() {
		updatePackageRepositoryRequest = DEFAULT_REQUEST;
	}

	let open = $state(false);
	function close() {
		open = false;
	}
</script>

<AlertDialog.Root bind:open>
	<AlertDialog.Trigger>
		<Button variant="ghost">
			<Icon icon="ph:pencil" /> Edit
		</Button>
	</AlertDialog.Trigger>
	<AlertDialog.Content>
		<AlertDialog.Header>
			<AlertDialog.Title>Edit {packageRepository.name}</AlertDialog.Title>
			<AlertDialog.Description class="p-2">
				<span class="grid gap-2">
					<Label>URL</Label>
					<Input bind:value={updatePackageRepositoryRequest.url} />
				</span>
			</AlertDialog.Description>
		</AlertDialog.Header>
		<AlertDialog.Footer>
			<AlertDialog.Cancel onclick={reset} class="mr-auto">Cancel</AlertDialog.Cancel>
			<AlertDialog.Action
				onclick={() => {
					client.updatePackageRepository(updatePackageRepositoryRequest).then((r) => {
						toast.info(`Update package repositories`);
						client.getConfiguration({}).then((r) => {
							configuration = r;
						});
					});
					// toast.info(`Update package repositories`);
					console.log(updatePackageRepositoryRequest);
					reset();
					close();
				}}
			>
				Update
			</AlertDialog.Action>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>
