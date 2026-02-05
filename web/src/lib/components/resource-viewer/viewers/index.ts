import type { Component } from 'svelte';

import DeleteDialog from '$lib/components/form/workspace/delete-dialog.svelte';
import EditSheet from '$lib/components/form/workspace/edit-sheet.svelte';

import Default from './default.svelte';
import Workspaces from './workspaces.svelte';

type ViewerProps = { object: any; schema?: any };
type ViewerType = Component<ViewerProps>;

type EditorProps = { name: string; schema?: any; object?: any; onsuccess?: () => void };
type EditorType = Component<EditorProps> | null;

type DeleterProps = { name: string };
type DeleterType = Component<DeleterProps> | null;

function getResourceViewer(resource: string): ViewerType {
	if (resource === 'workspaces') {
		return Workspaces as ViewerType;
	}
	return Default as ViewerType;
}

function getEditor(resource: string): EditorType {
	if (resource === 'workspaces') {
		return EditSheet as unknown as EditorType;
	}
	return null;
}

function getDeleter(resource: string): DeleterType {
	if (resource === 'workspaces') {
		return DeleteDialog as unknown as DeleterType;
	}
	return null;
}

export { getDeleter, getEditor, getResourceViewer };
export type { DeleterType, EditorType, ViewerType };
