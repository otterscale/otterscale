<script lang="ts" module>
	import { Group } from './group';
	import { NFS } from './nfs';

	import * as Tabs from '$lib/components/ui/tabs/index.js';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let {
		selectedScope = $bindable(),
		selectedFacility = $bindable(),
		selectedVolume = $bindable(),
		selectedSubvolumeGroupName = $bindable(),
	}: {
		selectedScope: string;
		selectedFacility: string;
		selectedVolume: string;
		selectedSubvolumeGroupName: string;
	} = $props();
</script>

{#snippet trigger()}
	<Tabs.List>
		<Tabs.Trigger value="nfs">{m.network_file_system()}</Tabs.Trigger>
		<Tabs.Trigger value="group">{m.group()}</Tabs.Trigger>
	</Tabs.List>
{/snippet}

<Tabs.Root value="nfs">
	<Tabs.Content value="nfs">
		<NFS bind:selectedScope bind:selectedFacility bind:selectedVolume bind:selectedSubvolumeGroupName {trigger} />
	</Tabs.Content>
	<Tabs.Content value="group">
		<Group bind:selectedScope bind:selectedFacility bind:selectedVolume {trigger} />
	</Tabs.Content>
</Tabs.Root>
