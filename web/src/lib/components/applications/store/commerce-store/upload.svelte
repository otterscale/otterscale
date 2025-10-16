<script lang="ts">
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import { ApplicationService } from '$lib/api/application/v1/application_pb';
	import * as AlertDialog from '$lib/components/ui/alert-dialog/index.js';
	import { Button } from '$lib/components/ui/button';
	import { Label } from '$lib/components/ui/label';
	import { cn } from '$lib/utils';
	// FileDropZone component is not available, using basic file input instead
	const MEGABYTE = 1024 * 1024;
	const displaySize = (bytes: number) => {
		if (bytes === 0) return '0 Bytes';
		const k = 1024;
		const sizes = ['Bytes', 'KB', 'MB', 'GB'];
		const i = Math.floor(Math.log(bytes) / Math.log(k));
		return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
	};
	let open = $state(false);
	let files = $state<Array<{ name: string; size: number; url: Promise<string>; uploadedAt: number }>>([]);

	const onUpload = (uploadedFiles: File[]) => {
		const maxFileSize = 20 * MEGABYTE; // 20MB limit for Helm charts

		// Only allow one file at a time
		if (uploadedFiles.length > 1) {
			toast.error('Please select only one file at a time');
			return;
		}

		// Replace existing file if any
		if (files.length > 0) {
			// Clean up existing file URL
			files[0].url.then((url) => URL.revokeObjectURL(url));
			files = [];
		}

		const file = uploadedFiles[0];
		if (!file) return;

		// Check file size
		if (file.size > maxFileSize) {
			toast.error(`${file.name} is too large. Maximum size is ${displaySize(maxFileSize)}`);
			return;
		}

		// Check file type
		if (!file.name.endsWith('.tgz') && !file.name.endsWith('.tar.gz')) {
			toast.error(`${file.name} is not a valid Helm chart file. Only .tgz and .tar.gz files are supported`);
			return;
		}

		// Add the single file
		files = [
			{
				name: file.name,
				size: file.size,
				url: Promise.resolve(URL.createObjectURL(file)),
				uploadedAt: Date.now(),
			},
		];
	};

	// Drag and drop functionality
	let isDragging = $state(false);

	const handleDragOver = (e: DragEvent) => {
		e.preventDefault();
		e.stopPropagation();
		isDragging = true;
	};

	const handleDragEnter = (e: DragEvent) => {
		e.preventDefault();
		e.stopPropagation();
		isDragging = true;
	};

	const handleDragLeave = (e: DragEvent) => {
		e.preventDefault();
		e.stopPropagation();
		// Only set isDragging to false if we're leaving the drop zone entirely
		const currentTarget = e.currentTarget as HTMLElement;
		const relatedTarget = e.relatedTarget as HTMLElement;
		if (!currentTarget?.contains(relatedTarget)) {
			isDragging = false;
		}
	};

	const handleDrop = (e: DragEvent) => {
		e.preventDefault();
		e.stopPropagation();
		isDragging = false;

		const droppedFiles = e.dataTransfer?.files;
		if (droppedFiles) {
			const fileArray = Array.from(droppedFiles);
			onUpload(fileArray); // Let onUpload handle validation
		}
	};

	const handleKeydown = (e: KeyboardEvent) => {
		if (e.key === 'Enter' || e.key === ' ') {
			e.preventDefault();
			document.getElementById('file-upload')?.click();
		}
	};

	const transport: Transport = getContext('transport');
	const client = createClient(ApplicationService, transport);

	const uploadChart = async () => {
		try {
			// Convert files to bytes
			const chartFile = files[0]; // Assuming first file is the chart
			if (!chartFile) {
				toast.error('Please select a chart file to upload');
				return;
			}

			// Validate file type
			if (!chartFile.name.endsWith('.tgz') && !chartFile.name.endsWith('.tar.gz')) {
				toast.error('Please select a valid Helm chart file (.tgz or .tar.gz)');
				return;
			}

			// No need to validate name and version - they will be parsed from chart.yaml

			// Get file content as bytes
			const url = await chartFile.url;
			const response = await fetch(url);
			const arrayBuffer = await response.arrayBuffer();
			const chartContent = new Uint8Array(arrayBuffer);

			const request: UploadChartRequest = {
				chartContent: chartContent,
			};

			await client.uploadChart(request);

			toast.success(`Chart uploaded successfully!`, {
				description: `Saved to local charts directory`,
			});
			close();
		} catch (error) {
			if (error instanceof ConnectError) {
				toast.error(`Upload failed: ${error.message}`);
			} else {
				toast.error('Unknown error occurred during upload');
			}
			console.error('Error uploading charts:', error);
		}
	};

	function close() {
		open = false;
		step = 0;
	}
	let { class: className }: { class?: string } = $props();
	let step = $state(0);
</script>

<AlertDialog.Root bind:open>
	<AlertDialog.Trigger class={cn('flex items-center gap-1', className)}>
		<Button>
			<Icon icon="ph:upload" />
			Upload
		</Button>
	</AlertDialog.Trigger>
	<AlertDialog.Content>
		{#if step == 0}
			<AlertDialog.Header>
				<AlertDialog.Title>
					Basic Information
					<div class="text-muted-foreground text-sm font-normal">Upload your Helm chart to OCI registry</div>
				</AlertDialog.Title>

				<div class="grid gap-2 space-y-2 py-4">
					<div class="mb-4">
						<p class="text-sm text-gray-600">Chart will be uploaded to the default OCI registry</p>
					</div>
					<div class="grid gap-4">
						<div class="flex w-full flex-col gap-2 p-6">
							<div
								class="rounded-lg border-2 border-dashed p-6 text-center transition-colors {isDragging
									? 'border-blue-400 bg-blue-50'
									: 'border-gray-300 hover:border-gray-400'}"
								role="button"
								tabindex="0"
								ondragover={handleDragOver}
								ondragenter={handleDragEnter}
								ondragleave={handleDragLeave}
								ondrop={handleDrop}
								onkeydown={handleKeydown}
								aria-label="Upload Helm chart file"
							>
								<input
									type="file"
									accept=".tgz,.tar.gz,application/gzip,application/x-gzip"
									class="hidden"
									id="file-upload"
									onchange={(e) => {
										const target = e.target as HTMLInputElement;
										if (target.files) {
											const fileArray = Array.from(target.files);
											onUpload(fileArray);
										}
									}}
								/>
								<Label for="file-upload" class="flex cursor-pointer flex-col items-center">
									<Icon icon="ph:upload" class="mb-2 text-4xl text-gray-400" />
									<p class="text-sm text-gray-600">
										{isDragging
											? 'Drop file here'
											: 'Click to upload or drag and drop Helm chart file'}
									</p>
									<p class="text-xs text-gray-400">Supports .tgz and .tar.gz files (one file only)</p>
								</Label>
							</div>

							<div class="flex flex-col gap-2">
								{#each files as file, i (file.name)}
									<div class="flex place-items-center justify-between gap-2">
										<div class="flex place-items-center gap-2">
											<div class="flex h-9 w-9 items-center justify-center rounded bg-gray-100">
												<Icon icon="ph:file-archive" class="h-5 w-5 text-gray-600" />
											</div>
											<div class="flex flex-col">
												<span class="text-sm font-medium">{file.name}</span>
												<span class="text-muted-foreground text-xs"
													>{displaySize(file.size)}</span
												>
											</div>
										</div>
										{#await file.url then url}
											<Button
												variant="outline"
												size="icon"
												onclick={() => {
													URL.revokeObjectURL(url);
													files = [...files.slice(0, i), ...files.slice(i + 1)];
												}}
											>
												<Icon icon="ph:x" />
											</Button>
										{/await}
									</div>
								{/each}
							</div>
						</div>
					</div>
				</div></AlertDialog.Header
			>
		{:else}
			<AlertDialog.Header>
				<AlertDialog.Title>
					Confirm Upload
					<div class="text-muted-foreground text-sm font-normal">
						Confirm Helm chart information and start upload to OCI registry
					</div>
				</AlertDialog.Title>

				<div class="space-y-4">
					<div class="bg-muted rounded-lg p-4">
						<h4 class="text-foreground font-medium">Chart Information</h4>
						<dl class="mt-2 space-y-1 text-sm">
							<div class="flex justify-between">
								<dt class="text-muted-foreground">Target Registry:</dt>
								<dd class="text-foreground text-xs font-medium">Default</dd>
							</div>
							{#if files.length > 0}
								<div class="flex justify-between">
									<dt class="text-muted-foreground">File:</dt>
									<dd class="text-foreground font-medium">
										{files[0].name} ({displaySize(files[0].size)})
									</dd>
								</div>
							{/if}
						</dl>
					</div>
				</div>
			</AlertDialog.Header>
		{/if}

		<AlertDialog.Footer>
			<AlertDialog.Cancel
				onclick={() => {
					close();
				}}
				class="mr-auto">Cancel</AlertDialog.Cancel
			>
			{#if step > 0}
				<Button
					variant="outline"
					onclick={() => {
						step--;
					}}
				>
					<Icon icon="ph:arrow-left" />
					Previous
				</Button>
			{/if}
			<Button
				disabled={step == 0 && files.length === 0}
				onclick={async () => {
					if (step == 1) {
						// Final step - upload charts
						await uploadChart();
					} else {
						// Move to next step
						if (files.length === 0) {
							toast.error('Please select a chart file before proceeding');
							return;
						}
						step++;
						if (step >= 2) {
							close();
						}
					}
				}}
			>
				{#if step == 1}
					Start Upload
				{:else}
					Next
				{/if}
			</Button>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>
