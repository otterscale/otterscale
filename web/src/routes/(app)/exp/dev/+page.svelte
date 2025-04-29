<script lang="ts">
	import { Button, buttonVariants } from '$lib/components/ui/button';
	import * as Table from '$lib/components/ui/table';
	import Icon from '@iconify/svelte';
	import * as Select from '$lib/components/ui/select/index.js';
	import { toast } from 'svelte-sonner';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import {
		Nexus,
		type Application_Chart,
		type Application_Chart_Version,
		type CreateReleaseRequest,
		type Facility_Info
	} from '$gen/api/nexus/v1/nexus_pb';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onMount } from 'svelte';
	import * as AlertDialog from '$lib/components/ui/alert-dialog';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import { writable } from 'svelte/store';

	import Badge from '$lib/components/ui/badge/badge.svelte';
	import { ReleaseValuesEdit } from '$lib/components/otterscale/index';
</script>

<div class="flex items-center gap-1">
	<Icon icon="ph:desktop" class="size-5" />
	<Label class="text-base">System</Label>
</div>

<div class="flex items-center gap-1">
	<Icon icon="ph:circuitry" class="size-5" />
	<Label class="text-base">Mainboard</Label>
</div>

<div class="flex items-center gap-1">
	<Icon icon="ph:computer-tower" class="size-5" />
	<Label class="text-base">Chassis</Label>
</div>
