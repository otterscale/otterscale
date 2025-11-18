<script lang="ts" module>
	import * as Tabs from '$lib/components/ui/tabs/index.js';
	import { m } from '$lib/paraglide/messages';

	import { Group } from './group';
	import { NFS } from './nfs';
</script>

<script lang="ts">
	let {
		scope,
		volume,
		selectedSubvolumeGroupName = $bindable()
	}: {
		scope: string;
		volume: string;
		selectedSubvolumeGroupName: string;
	} = $props();
</script>

{#snippet tabsTrigger()}
	<Tabs.List>
		<Tabs.Trigger value="nfs">{m.network_file_system()}</Tabs.Trigger>
		<Tabs.Trigger value="group">{m.group()}</Tabs.Trigger>
	</Tabs.List>
{/snippet}
<Tabs.Root value="nfs">
	<Tabs.Content value="nfs">
		<NFS {scope} {volume} bind:selectedSubvolumeGroupName trigger={tabsTrigger} />
	</Tabs.Content>
	<Tabs.Content value="group">
		<Group {scope} {volume} trigger={tabsTrigger} />
	</Tabs.Content>
</Tabs.Root>
