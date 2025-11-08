<script lang="ts" module>
	import { Bucket } from './bucket';
	import { User } from './user';

	import * as Tabs from '$lib/components/ui/tabs/index.js';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let {
		selectedScope = $bindable(),
		selectedFacility = $bindable()
	}: {
		selectedScope: string;
		selectedFacility: string;
	} = $props();
</script>

{#snippet trigger()}
	<Tabs.List>
		<Tabs.Trigger value="bucket">{m.bucket()}</Tabs.Trigger>
		<Tabs.Trigger value="user">{m.user()}</Tabs.Trigger>
	</Tabs.List>
{/snippet}
<Tabs.Root value="bucket">
	<Tabs.Content value="bucket">
		<Bucket bind:selectedScope bind:selectedFacility {trigger} />
	</Tabs.Content>
	<Tabs.Content value="user">
		<User bind:selectedScope bind:selectedFacility {trigger} />
	</Tabs.Content>
</Tabs.Root>
