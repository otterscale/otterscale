<script lang="ts">
	import { Tabs } from '$lib/components/custom/tabs/index';
	import { IoTest } from './io-test';
	import { ObjectStorageTest } from './object-storage-test';
	import Icon from '@iconify/svelte';
</script>

<Tabs.Root value="io-test">
	<Tabs.List class="*:gap-2 [*_&_svg]:size-5">
		<Tabs.Trigger value="io-test">
			<Icon icon="ph:cube" />
			IO Test
		</Tabs.Trigger>
		<Tabs.Trigger value="object-storage-test">
			<Icon icon="ph:cube" />
			Object Storage Test
		</Tabs.Trigger>
	</Tabs.List>
	<Tabs.Content value="io-test">
		<IoTest />
	</Tabs.Content>
	<Tabs.Content value="object-storage-test">
		<ObjectStorageTest />
	</Tabs.Content>
</Tabs.Root>
