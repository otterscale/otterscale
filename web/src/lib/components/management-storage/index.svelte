<script lang="ts">
	import { LayeredTabs } from '$lib/components/custom/tabs/index';
	import Icon from '@iconify/svelte';
	import { BlockImage } from './storage-block-image';
	import { Bucket } from './storage-bucket';
	import { Volume } from './storage-filesystem-volume';
	import { Subvolume } from './storage-filesystem-subvolume';
	import { SubvolumeGroup } from './storage-filesystem-subvolume-group';
	import { ObjectStorageDaemon } from './storage-object-storage-daemon';
	import { Pool } from './storage-pool';
	import { Monitor } from './storage-monitor';
	import { Role } from './storage-role';
	import { User } from './storage-user';

	let value = $state('user');
</script>

<LayeredTabs.Root bind:value>
	<LayeredTabs.Menubar>
		<LayeredTabs.Menu>
			<LayeredTabs.Trigger>
				<Icon icon="ph:circles-four" />
				CLUSTER
			</LayeredTabs.Trigger>
			<LayeredTabs.MenuContent>
				<LayeredTabs.MenuItem value="monitor">
					<Icon icon="ph:pulse" />
					Monitor
				</LayeredTabs.MenuItem>
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
					<Icon icon="ph:hard-drive" />
					RADOS Block Device
				</LayeredTabs.MenuItem>
			</LayeredTabs.MenuContent>
		</LayeredTabs.Menu>
		<LayeredTabs.Menu>
			<LayeredTabs.Trigger>
				<Icon icon="ph:folder-simple" />
				FILE SYSTEM
			</LayeredTabs.Trigger>
			<LayeredTabs.MenuContent>
				<LayeredTabs.MenuItem value="volume">
					<Icon icon="ph:folder-open" />
					Volume
				</LayeredTabs.MenuItem>
				<LayeredTabs.MenuItem value="subvolume-group">
					<Icon icon="ph:folders" />
					Subvolume Group
				</LayeredTabs.MenuItem>
				<LayeredTabs.MenuItem value="subvolume">
					<Icon icon="ph:folder-simple" />
					Subvolume
				</LayeredTabs.MenuItem>
			</LayeredTabs.MenuContent>
		</LayeredTabs.Menu>
		<LayeredTabs.Menu>
			<LayeredTabs.Trigger>
				<Icon icon="ph:cloud" />
				GATEWAY
			</LayeredTabs.Trigger>
			<LayeredTabs.MenuContent>
				<LayeredTabs.MenuItem value="bucket">
					<Icon icon="ph:cube" />
					Bucket
				</LayeredTabs.MenuItem>
				<LayeredTabs.MenuItem value="user">
					<Icon icon="ph:user-circle" />
					User
				</LayeredTabs.MenuItem>
			</LayeredTabs.MenuContent>
		</LayeredTabs.Menu>
	</LayeredTabs.Menubar>

	<LayeredTabs.Content value="block-image">
		{#if value === 'block-image'}
			<BlockImage />
		{/if}
	</LayeredTabs.Content>
	<LayeredTabs.Content value="bucket">
		{#if value === 'bucket'}
			<Bucket />
		{/if}
	</LayeredTabs.Content>
	<LayeredTabs.Content value="monitor">
		{#if value === 'monitor'}
			<Monitor />
		{/if}
	</LayeredTabs.Content>
	<LayeredTabs.Content value="object-storage-daemon">
		{#if value === 'object-storage-daemon'}
			<ObjectStorageDaemon />
		{/if}
	</LayeredTabs.Content>
	<LayeredTabs.Content value="pool">
		{#if value === 'pool'}
			<Pool />
		{/if}
	</LayeredTabs.Content>
	<LayeredTabs.Content value="role">
		{#if value === 'role'}
			<Role />
		{/if}
	</LayeredTabs.Content>
	<LayeredTabs.Content value="subvolume">
		{#if value === 'subvolume'}
			<Subvolume />
		{/if}
	</LayeredTabs.Content>
	<LayeredTabs.Content value="subvolume-group">
		{#if value === 'subvolume-group'}
			<SubvolumeGroup />
		{/if}
	</LayeredTabs.Content>
	<LayeredTabs.Content value="user">
		{#if value === 'user'}
			<User />
		{/if}
	</LayeredTabs.Content>
	<LayeredTabs.Content value="volume">
		{#if value === 'volume'}
			<Volume />
		{/if}
	</LayeredTabs.Content>
</LayeredTabs.Root>
