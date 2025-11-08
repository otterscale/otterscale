<script lang="ts">
	import * as Tabs from '$lib/components/ui/tabs/index';
	import { m } from '$lib/paraglide/messages';

	import { IOTest } from './io-test';
	import { ObjectStorageTest } from './object-storage-test';

	let selectedTab = $state('io-test');
</script>

{#snippet trigger()}
	<Tabs.List>
		<Tabs.Trigger value="io-test">{m.io_test()}</Tabs.Trigger>
		<Tabs.Trigger value="object-storage-test">{m.object_storage_test()}</Tabs.Trigger>
	</Tabs.List>
{/snippet}
<Tabs.Root bind:value={selectedTab}>
	<Tabs.Content value="io-test" class="max-w-[70vw]">
		<IOTest {selectedTab} {trigger} />
	</Tabs.Content>
	<Tabs.Content value="object-storage-test" class="max-w-[70vw]">
		<ObjectStorageTest {selectedTab} {trigger} />
	</Tabs.Content>
</Tabs.Root>
