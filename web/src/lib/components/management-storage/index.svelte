<script lang="ts">
	import * as Tabs from '$lib/components/custom/tabs';
	import { StorageFileSystem } from './storage-filesystem';
	import { StoragePool } from './storage-pool';
	import { StorageBlockImage } from './storage-block-image';
	import { StorageClusterOSD } from './storage-cluster-osd';
	import Icon from '@iconify/svelte';
</script>

<Tabs.Root value="block-image">
	<Tabs.List class="*:gap-2 [*_&_svg]:size-5">
		<Tabs.Trigger value="block-image" class="flex items-center">
			<Icon icon="ph:hard-drives" />
			BLOCK IMAGE
		</Tabs.Trigger>
		<Tabs.Trigger value="filesystem" class="flex items-center">
			<Icon icon="ph:folder" />
			FILE SYSTEM
		</Tabs.Trigger>
		<Tabs.Trigger value="pool" class="flex items-center">
			<Icon icon="ph:database" />
			POOL
		</Tabs.Trigger>
		<Tabs.Trigger value="cluster-osd">
			<Icon icon="ph:cube" />
			OBJECT STORAGE DAEMON
		</Tabs.Trigger>
	</Tabs.List>
	<Tabs.Content value="block-image">
		<StorageBlockImage />
	</Tabs.Content>
	<Tabs.Content value="filesystem">
		<StorageFileSystem />
	</Tabs.Content>
	<Tabs.Content value="pool">
		<StoragePool />
	</Tabs.Content>
	<Tabs.Content value="cluster-osd">
		<StorageClusterOSD />
	</Tabs.Content>
</Tabs.Root>
