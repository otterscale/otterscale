<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import type { Row } from '@tanstack/table-core';

	import type { SMBShare } from '$lib/api/storage/v1/storage_pb';
	import {
		SMBShare_CommonConfig_MapToGuest,
		SMBShare_SecurityConfig_Mode
	} from '$lib/api/storage/v1/storage_pb';
	import { Cells } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { ReloadManager } from '$lib/components/custom/reloader';
	import { TagGroup } from '$lib/components/tag-group';
	import { Badge } from '$lib/components/ui/badge';
	import { formatCapacity } from '$lib/formatter';

	import Actions from './cell-actions.svelte';

	export const cells = {
		row_picker,
		name,
		namespace,
		status,
		size,
		browsable,
		read_only,
		guest_ok,
		map_to_guest,
		security_mode,
		valid_users,
		actions
	};

	function getMapToGuestLabel(mapToGuest: SMBShare_CommonConfig_MapToGuest): string {
		switch (mapToGuest) {
			case SMBShare_CommonConfig_MapToGuest.NEVER:
				return 'Never';
			case SMBShare_CommonConfig_MapToGuest.BAD_USER:
				return 'Bad User';
			case SMBShare_CommonConfig_MapToGuest.BAD_PASSWORD:
				return 'Bad Password';
			default:
				return 'Unknown';
		}
	}

	const getSecurityModeLabel = (securityMode: SMBShare_SecurityConfig_Mode) => {
		switch (securityMode) {
			case SMBShare_SecurityConfig_Mode.USER:
				return 'User';
			case SMBShare_SecurityConfig_Mode.ACTIVE_DIRECTORY:
				return 'Active Directory';
			default:
				return 'Unknown';
		}
	};
</script>

{#snippet row_picker(row: Row<SMBShare>)}
	<Layout.Cell class="items-center">
		<Cells.RowPicker {row} />
	</Layout.Cell>
{/snippet}

{#snippet name(row: Row<SMBShare>)}
	<Layout.Cell class="items-start">
		{row.original.name}
	</Layout.Cell>
{/snippet}

{#snippet namespace(row: Row<SMBShare>)}
	<Layout.Cell class="items-start">
		{row.original.namespace}
	</Layout.Cell>
{/snippet}

{#snippet status(row: Row<SMBShare>)}
	<Layout.Cell class="items-start">
		<Badge variant="outline">
			{row.original.healthies}
		</Badge>
	</Layout.Cell>
{/snippet}

{#snippet size(row: Row<SMBShare>)}
	{@const { value, unit } = formatCapacity(Number(row.original.sizeBytes))}
	<Layout.Cell class="items-end">
		{value}
		{unit}
	</Layout.Cell>
{/snippet}

{#snippet browsable(row: Row<SMBShare>)}
	<Layout.Cell class="items-end">
		{@const value = row.original.browsable}
		<Icon icon={value ? 'ph:check' : 'ph:x'} class={value ? 'text-green-500' : 'text-red-500'} />
	</Layout.Cell>
{/snippet}

{#snippet read_only(row: Row<SMBShare>)}
	<Layout.Cell class="items-end">
		{@const value = row.original.readOnly}
		<Icon icon={value ? 'ph:check' : 'ph:x'} class={value ? 'text-green-500' : 'text-red-500'} />
	</Layout.Cell>
{/snippet}

{#snippet guest_ok(row: Row<SMBShare>)}
	<Layout.Cell class="items-end">
		{@const value = row.original.guestOk}
		<Icon icon={value ? 'ph:check' : 'ph:x'} class={value ? 'text-green-500' : 'text-red-500'} />
	</Layout.Cell>
{/snippet}

{#snippet map_to_guest(row: Row<SMBShare>)}
	{#if row.original.commonConfig}
		<Layout.Cell class="items-start">
			<Badge variant="outline">{getMapToGuestLabel(row.original.commonConfig.mapToGuest)}</Badge>
		</Layout.Cell>
	{/if}
{/snippet}

{#snippet security_mode(row: Row<SMBShare>)}
	{#if row.original.securityConfig}
		<Layout.Cell class="items-start">
			<Badge variant="outline">{getSecurityModeLabel(row.original.securityConfig.mode)}</Badge>
		</Layout.Cell>
	{/if}
{/snippet}

{#snippet valid_users(row: Row<SMBShare>)}
	<Layout.Cell class="items-start">
		{#if row.original.validUsers && row.original.validUsers.length > 0}
			<TagGroup
				items={row.original.validUsers.map((validUser) => ({ title: validUser, icon: 'ph:user' }))}
			/>
		{/if}
	</Layout.Cell>
{/snippet}

{#snippet actions(data: { row: Row<SMBShare>; scope: string; reloadManager: ReloadManager })}
	<Layout.Cell class="items-end">
		<Actions smbShare={data.row.original} scope={data.scope} reloadManager={data.reloadManager} />
	</Layout.Cell>
{/snippet}
