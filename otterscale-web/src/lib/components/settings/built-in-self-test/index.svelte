<script lang="ts">
	import * as Tabs from '$lib/components/ui/tabs/index';
	import { IOTest } from './io-test';
	import { ObjectStorageTest } from './object-storage-test';
</script>

{#snippet trigger()}
	<Tabs.List>
		<Tabs.Trigger value="io-test">IO Test</Tabs.Trigger>
		<Tabs.Trigger value="object-storage-test">Object Storage Test</Tabs.Trigger>
	</Tabs.List>
{/snippet}
<Tabs.Root value="io-test">
	<Tabs.Content value="io-test">
		<IOTest {trigger} />
	</Tabs.Content>
	<Tabs.Content value="object-storage-test">
		<ObjectStorageTest {trigger} />
	</Tabs.Content>
</Tabs.Root>
