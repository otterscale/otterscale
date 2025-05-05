<script lang="ts">
	import { toast } from 'svelte-sonner';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import { Label } from '$lib/components/ui/label';
	import {
		Nexus,
		type Application_Release,
		type RollbackReleaseRequest
	} from '$gen/api/nexus/v1/nexus_pb';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext } from 'svelte';
	import * as AlertDialog from '$lib/components/ui/alert-dialog';

	let {
		releases = $bindable(),
		release
	}: {
		releases: Application_Release[];
		release: Application_Release;
	} = $props();

	const transport: Transport = getContext('transportNEW');
	const client = createClient(Nexus, transport);

	const DEFAULT_REQUEST = {
		dryRun: false,
		scopeUuid: release.name,
		facilityName: release.name,
		namespace: release.namespace,
		name: release.name
	} as RollbackReleaseRequest;

	let rollbackReleaseRequest = $state(DEFAULT_REQUEST);

	function reset() {
		rollbackReleaseRequest = { dryRun: false } as RollbackReleaseRequest;
	}

	let open = $state(false);
	function close() {
		open = false;
	}
</script>

<AlertDialog.Root bind:open>
	<AlertDialog.Trigger>Rollback</AlertDialog.Trigger>
	<AlertDialog.Content>
		<AlertDialog.Header>
			<AlertDialog.Title>Rollback Release {release.name}</AlertDialog.Title>
			<AlertDialog.Description class="rounded-lg bg-muted/50">
				<div class="grid gap-4 p-4">
					<p class="text-sm text-muted-foreground">
						Are you sure you want to rollback release "{release.name}" in namespace "{release.namespace}"?
					</p>
					<div class="flex items-center justify-end space-x-4">
						<Label for="dry-run">Dry Run</Label>
						<Switch bind:checked={rollbackReleaseRequest.dryRun} />
					</div>
				</div>
			</AlertDialog.Description>
		</AlertDialog.Header>
		<AlertDialog.Footer>
			<AlertDialog.Cancel onclick={reset} class="mr-auto">Cancel</AlertDialog.Cancel>
			<AlertDialog.Action
				onclick={() => {
					client
						.rollbackRelease(rollbackReleaseRequest)
						.then((r) => {
							toast.info(`Rollback ${rollbackReleaseRequest.name}`);
							client.listReleases({}).then((r) => {
								releases = r.releases;
							});
						})
						.catch((e) => {
							toast.error(`Fail to rollback ${rollbackReleaseRequest.name}: ${e.toString()}`);
						});
					reset();
					close();
				}}>Confirm</AlertDialog.Action
			>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>
