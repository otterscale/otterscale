<script lang="ts">
	import * as Tabs from '$lib/components/custom/tabs';
	import { StorageFileSystem } from './storage-filesystem';
	import { StoragePool } from './storage-pool';
	import { StorageBlockImage } from './storage-block-image';
	import { StorageClusterObjectStorageDaemon } from './storage-cluster-object-storage-daemon';
</script>

<Tabs.Root value="filesystem">
	<Tabs.List>
		<Tabs.Trigger value="filesystem">File System</Tabs.Trigger>
		<Tabs.Trigger value="pool">Pool</Tabs.Trigger>
		<Tabs.Trigger value="block-image">Block Image</Tabs.Trigger>
		<Tabs.Trigger value="cluster-osd">Object Storage Daemon</Tabs.Trigger>
	</Tabs.List>
	<Tabs.Content value="filesystem">
		<StorageFileSystem />
	</Tabs.Content>
	<Tabs.Content value="pool">
		<StoragePool />
	</Tabs.Content>
	<Tabs.Content value="block-image">
		<StorageBlockImage />
	</Tabs.Content>
	<Tabs.Content value="cluster-osd">
		<StorageClusterObjectStorageDaemon />
	</Tabs.Content>
</Tabs.Root>
