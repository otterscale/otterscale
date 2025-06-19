<script lang="ts">
	import { LayeredTabs } from '$lib/components/custom/tabs/index';
	import Icon from '@iconify/svelte';
	import { BlockImage } from './storage-block-image';
	import { Bucket } from './storage-bucket';
	import { FileSystem } from './storage-filesystem';
	import { ObjectStorageDaemon } from './storage-object-storage-daemon';
	import { Pool } from './storage-pool';
	import { Role } from './storage-role';
	import { User } from './storage-user';
</script>

<LayeredTabs.Root value="block-image">
	<LayeredTabs.Menubar>
		<LayeredTabs.Menu>
			<LayeredTabs.Trigger>
				<Icon icon="ph:cube" />
				CLUSTER
			</LayeredTabs.Trigger>
			<LayeredTabs.MenuContent>
				<LayeredTabs.MenuItem value="pool">
					<Icon icon="ph:database" />
					Pool
				</LayeredTabs.MenuItem>
				<LayeredTabs.MenuItem value="object-storage-daemon">
					<Icon icon="ph:cube" />
					Object Storage Daemon
				</LayeredTabs.MenuItem>
			</LayeredTabs.MenuContent>
		</LayeredTabs.Menu>
		<LayeredTabs.Menu>
			<LayeredTabs.Trigger>
				<Icon icon="ph:hard-drives" />
				BLOCK DEVICE
			</LayeredTabs.Trigger>
			<LayeredTabs.MenuContent>
				<LayeredTabs.MenuItem value="block-image">
					<Icon icon="ph:hard-drives" />
					RADOS Block Device
				</LayeredTabs.MenuItem>
			</LayeredTabs.MenuContent>
		</LayeredTabs.Menu>
		<LayeredTabs.Menu>
			<LayeredTabs.Trigger>
				<Icon icon="ph:folder" />
				FILE SYSTEM
			</LayeredTabs.Trigger>
			<LayeredTabs.MenuContent>
				<LayeredTabs.MenuItem value="file-system">
					<Icon icon="ph:folder" />
					File System
				</LayeredTabs.MenuItem>
			</LayeredTabs.MenuContent>
		</LayeredTabs.Menu>
		<LayeredTabs.Menu>
			<LayeredTabs.Trigger>
				<Icon icon="ph:cube" />
				GATEWAY
			</LayeredTabs.Trigger>
			<LayeredTabs.MenuContent>
				<LayeredTabs.MenuItem value="bucket">
					<Icon icon="ph:cube" />
					Bucket
				</LayeredTabs.MenuItem>
				<LayeredTabs.MenuSeparator />
				<LayeredTabs.MenuItem value="user">
					<Icon icon="ph:user" />
					User
				</LayeredTabs.MenuItem>
				<LayeredTabs.MenuItem value="role">
					<Icon icon="ph:user" />
					Role
				</LayeredTabs.MenuItem>
			</LayeredTabs.MenuContent>
		</LayeredTabs.Menu>
	</LayeredTabs.Menubar>

	<LayeredTabs.Content value="bucket">
		<Bucket />
	</LayeredTabs.Content>
	<LayeredTabs.Content value="block-image">
		<BlockImage />
	</LayeredTabs.Content>
	<LayeredTabs.Content value="file-system">
		<FileSystem />
	</LayeredTabs.Content>
	<LayeredTabs.Content value="object-storage-daemon">
		<ObjectStorageDaemon />
	</LayeredTabs.Content>
	<LayeredTabs.Content value="pool">
		<Pool />
	</LayeredTabs.Content>
	<LayeredTabs.Content value="role">
		<Role />
	</LayeredTabs.Content>
	<LayeredTabs.Content value="user">
		<User />
	</LayeredTabs.Content>
</LayeredTabs.Root>
