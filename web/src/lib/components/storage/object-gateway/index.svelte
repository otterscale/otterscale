<script lang="ts" module>
	import * as Tabs from '$lib/components/ui/tabs/index.js';
	import { m } from '$lib/paraglide/messages';
	import { Bucket } from './bucket';
	import { User } from './user';
</script>

<script lang="ts">
	let {
		selectedScopeUuid = $bindable(),
		selectedFacility = $bindable(),
	}: {
		selectedScopeUuid: string;
		selectedFacility: string;
	} = $props();
</script>

{#snippet trigger()}
	<Tabs.List>
		<Tabs.Trigger value="user">{m.user()}</Tabs.Trigger>
		<Tabs.Trigger value="bucket">{m.bucket()}</Tabs.Trigger>
	</Tabs.List>
{/snippet}
<Tabs.Root value="user">
	<Tabs.Content value="user">
		<User bind:selectedScopeUuid bind:selectedFacility {trigger} />
	</Tabs.Content>
	<Tabs.Content value="bucket">
		<Bucket bind:selectedScopeUuid bind:selectedFacility {trigger} />
	</Tabs.Content>
</Tabs.Root>
