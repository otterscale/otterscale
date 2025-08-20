<script lang="ts">
	import * as Tabs from '$lib/components/ui/tabs/index.js';
	import { Group } from './group';
	import { NFS } from './nfs';

	let {
		selectedScopeUuid = $bindable(),
		selectedFacility = $bindable(),
		selectedVolume = $bindable(),
		selectedSubvolumeGroupName = $bindable()
	}: {
		selectedScopeUuid: string;
		selectedFacility: string;
		selectedVolume: string;
		selectedSubvolumeGroupName: string;
	} = $props();
</script>

{#snippet trigger()}
	<Tabs.List>
		<Tabs.Trigger value="nfs">NFS</Tabs.Trigger>
		<Tabs.Trigger value="group">Group</Tabs.Trigger>
	</Tabs.List>
{/snippet}
<Tabs.Root value="nfs">
	<Tabs.Content value="nfs">
		<NFS
			bind:selectedScopeUuid
			bind:selectedFacility
			bind:selectedVolume
			bind:selectedSubvolumeGroupName
			{trigger}
		/>
	</Tabs.Content>
	<Tabs.Content value="group">
		<Group bind:selectedScopeUuid bind:selectedFacility bind:selectedVolume {trigger} />
	</Tabs.Content>
</Tabs.Root>
