<script lang="ts">
	import * as StorageManagement from '$lib/components/management-storage';
	import { Tabs } from '$lib/components/custom/tabs/index';
</script>

<Tabs.Root value="subvolume">
	<Tabs.List>
		<!-- <Tabs.Trigger value="volume">Volume</Tabs.Trigger> -->
		<Tabs.Trigger value="subvolume">NFS</Tabs.Trigger>
		<Tabs.Trigger value="subvolume group">Group</Tabs.Trigger>
	</Tabs.List>
	<!-- <Tabs.Content value="volume">
		<StorageManagement.Volume />
	</Tabs.Content> -->
	<Tabs.Content value="subvolume">
		<StorageManagement.Subvolume />
	</Tabs.Content>
	<Tabs.Content value="subvolume group">
		<StorageManagement.SubvolumeGroup />
	</Tabs.Content>
</Tabs.Root>
