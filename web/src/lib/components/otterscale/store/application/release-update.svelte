<script lang="ts">
	import { toast } from 'svelte-sonner';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import {
		Nexus,
		type Application_Release,
		type UpdateReleaseRequest
	} from '$gen/api/nexus/v1/nexus_pb';
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import { getContext } from 'svelte';
	import * as AlertDialog from '$lib/components/ui/alert-dialog';
	import { ReleaseValuesEdit } from '$lib/components/otterscale/index';

	let {
		releases = $bindable(),
		release,
		valuesYaml
	}: {
		releases: Application_Release[];
		release: Application_Release;
		valuesYaml: string;
	} = $props();

	const transport: Transport = getContext('transport');
	const client = createClient(Nexus, transport);

	const DEFAULT_REQUEST = {
		dryRun: false,
		scopeUuid: release.scopeUuid,
		facilityName: release.name,
		namespace: release.namespace,
		name: release.name,
		chartRef: release.version?.chartRef,
		valuesYaml: valuesYaml
	} as UpdateReleaseRequest;

	let updateReleaseRequest = $state(DEFAULT_REQUEST);

	function reset() {
		updateReleaseRequest = DEFAULT_REQUEST;
	}

	let open = $state(false);
	function close() {
		open = false;
	}
</script>

<AlertDialog.Root bind:open>
	<AlertDialog.Trigger>Update</AlertDialog.Trigger>
	<AlertDialog.Content>
		<AlertDialog.Header>
			<AlertDialog.Description>
				<fieldset class="items-center rounded-lg border p-3">
					<legend class="text-base">Basic</legend>
					<div class="grid gap-3">
						<span>
							<Label>Scope</Label>
							<Input bind:value={updateReleaseRequest.scopeUuid} />
						</span>
						<span>
							<Label>Facility</Label>
							<Input bind:value={updateReleaseRequest.facilityName} />
						</span>
						<span>
							<Label>Namespace</Label>
							<Input bind:value={updateReleaseRequest.namespace} />
						</span>
						<span>
							<Label>Name</Label>
							<Input bind:value={updateReleaseRequest.name} />
						</span>
						<span>
							<Label>Reference</Label>
							<Input bind:value={updateReleaseRequest.chartRef} />
						</span>
						<span class="flex items-center justify-between">
							<Label>Dry Run</Label>
							<Switch id="enable_ssh" bind:checked={updateReleaseRequest.dryRun} />
						</span>
					</div>
				</fieldset>
				<fieldset class="items-center rounded-lg border p-3">
					<legend class="text-base">Advance</legend>

					<span class="flex items-center justify-between">
						<Label>Configuration</Label>
						<ReleaseValuesEdit
							chartRef={updateReleaseRequest.chartRef}
							bind:valuesYaml={updateReleaseRequest.valuesYaml}
						/>
					</span>
				</fieldset>
			</AlertDialog.Description>
		</AlertDialog.Header>
		<AlertDialog.Footer>
			<AlertDialog.Cancel onclick={reset} class="mr-auto">Cancel</AlertDialog.Cancel>
			<AlertDialog.Action
				onclick={() => {
					toast.promise(() => client.updateRelease(updateReleaseRequest), {
						loading: 'Loading...',
						success: (r) => {
							client.listReleases({}).then((r) => {
								releases = r.releases;
							});
							return `Update ${r.name} success`;
						},
						error: (e) => {
							let msg = `Fail to update ${updateReleaseRequest.name}`;
							toast.error(msg, {
								description: (e as ConnectError).message.toString(),
								duration: Number.POSITIVE_INFINITY
							});
							return msg;
						}
					});

					reset();
					close();
				}}>Confirm</AlertDialog.Action
			>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>
