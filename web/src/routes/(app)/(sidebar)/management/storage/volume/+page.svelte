<script lang="ts">
	import * as StorageManagement from '$lib/components/management-storage';
	import { Tabs } from '$lib/components/custom/tabs/index';
</script>

<Tabs.Root value="volume">
	<Tabs.List>
		<Tabs.Trigger value="volume">Volume</Tabs.Trigger>
		<Tabs.Trigger value="subvolume group">Subvolume Group</Tabs.Trigger>
		<Tabs.Trigger value="subvolume">Subvolume</Tabs.Trigger>
	</Tabs.List>
	<Tabs.Content value="volume">
		<StorageManagement.Volume />
	</Tabs.Content>
	<Tabs.Content value="subvolume group">
		<StorageManagement.SubvolumeGroup />
	</Tabs.Content>
	<Tabs.Content value="subvolume">
		<StorageManagement.Subvolume />
	</Tabs.Content>
</Tabs.Root>
