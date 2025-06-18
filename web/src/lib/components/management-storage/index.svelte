<script lang="ts">
	import * as Tabs from '$lib/components/custom/tabs';
	import * as Menubar from '$lib/components/ui/menubar/index';
	import Icon from '@iconify/svelte';
	import { BlockImage } from './storage-block-image';
	import { Bucket } from './storage-bucket';
	import { ObjectStorageDaemon } from './storage-object-storage-daemon';
	import { FileSystem } from './storage-filesystem';
	import { Pool } from './storage-pool';
	import { Role } from './storage-role';
	import { User } from './storage-user';

	let value = $state('block-image');
</script>

<Tabs.Root {value}>
	<Menubar.Root class="border-none shadow-none">
		<Menubar.Menu>
			<Menubar.Trigger class="flex items-center gap-1">
				<Icon icon="ph:cube" />
				CLUSTER</Menubar.Trigger
			>
			<Menubar.Content class="*:text-xs">
				<Menubar.Item
					class="flex items-center gap-1"
					onclick={() => {
						value = 'pool';
					}}
				>
					<Icon icon="ph:database" />
					Pool
				</Menubar.Item>
				<Menubar.Item
					class="flex items-center gap-1"
					onclick={() => {
						value = 'object-storage-daemon';
					}}
				>
					<Icon icon="ph:cube" />
					Object Storage Daemon
				</Menubar.Item>
			</Menubar.Content>
		</Menubar.Menu>
		<Menubar.Menu>
			<Menubar.Trigger class="flex items-center gap-1">
				<Icon icon="ph:hard-drives" />
				RADOS BLOCK DEVICE</Menubar.Trigger
			>
			<Menubar.Content class="*:text-xs">
				<Menubar.Item
					class="flex items-center gap-1"
					onclick={() => {
						value = 'block-image';
					}}
				>
					<Icon icon="ph:hard-drives" />
					RADOS Block Device
				</Menubar.Item>
			</Menubar.Content>
		</Menubar.Menu>
		<Menubar.Menu>
			<Menubar.Trigger class="flex items-center gap-1">
				<Icon icon="ph:folder" />
				FILE SYSTEM</Menubar.Trigger
			>
			<Menubar.Content class="*:text-xs">
				<Menubar.Item
					class="flex items-center gap-1"
					onclick={() => {
						value = 'file-system';
					}}
				>
					<Icon icon="ph:folder" />
					File System
				</Menubar.Item>
			</Menubar.Content>
		</Menubar.Menu>
		<Menubar.Menu>
			<Menubar.Trigger class="flex items-center gap-1">
				<Icon icon="ph:cube" />
				RADOS GATEWAY</Menubar.Trigger
			>
			<Menubar.Content class="*:text-xs">
				<Menubar.Item
					class="flex items-center gap-1"
					onclick={() => {
						value = 'bucket';
					}}
				>
					<Icon icon="ph:cube" />
					Bucket
				</Menubar.Item>
				<Menubar.Separator />
				<Menubar.Item
					class="flex items-center gap-1"
					onclick={() => {
						value = 'user';
					}}
				>
					<Icon icon="ph:user" />
					User
				</Menubar.Item>
				<Menubar.Item
					class="flex items-center gap-1"
					onclick={() => {
						value = 'role';
					}}
				>
					<Icon icon="ph:user" />
					Role
				</Menubar.Item>
			</Menubar.Content>
		</Menubar.Menu>
	</Menubar.Root>

	<Tabs.Content value="block-image">
		<BlockImage />
	</Tabs.Content>
	<Tabs.Content value="file-system">
		<FileSystem />
	</Tabs.Content>
	<Tabs.Content value="pool">
		<Pool />
	</Tabs.Content>
	<Tabs.Content value="object-storage-daemon">
		<ObjectStorageDaemon />
	</Tabs.Content>
	<Tabs.Content value="bucket">
		<Bucket />
	</Tabs.Content>
	<Tabs.Content value="user">
		<User />
	</Tabs.Content>
	<Tabs.Content value="role">
		<Role />
	</Tabs.Content>
</Tabs.Root>
