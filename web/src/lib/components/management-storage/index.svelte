<script lang="ts">
	import * as Tabs from '$lib/components/custom/tabs';
	import { StorageFileSystem } from './storage-filesystem';
</script>

<Tabs.Root value="filesystem">
	<Tabs.List>
		<Tabs.Trigger value="filesystem">File System</Tabs.Trigger>
	</Tabs.List>
	<Tabs.Content value="filesystem">
		<StorageFileSystem />
	</Tabs.Content>
</Tabs.Root>
