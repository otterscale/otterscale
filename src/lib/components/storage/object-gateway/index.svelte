<script lang="ts" module>
	import * as Tabs from '$lib/components/ui/tabs/index.js';
	import { m } from '$lib/paraglide/messages';

	import { Bucket } from './bucket';
	import { User } from './user';
</script>

<script lang="ts">
	let {
		scope
	}: {
		scope: string;
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
		<Bucket {scope} {trigger} />
	</Tabs.Content>
	<Tabs.Content value="user">
		<User {scope} {trigger} />
	</Tabs.Content>
</Tabs.Root>
