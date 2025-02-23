import { page } from '$app/state';
import PocketBase from 'pocketbase';

const pb = new PocketBase('http://127.0.0.1:8090');
pb.autoCancellation(false);

if (pb.authStore.isValid && pb.authStore.record) {
    pb.collection(pb.authStore.record.collectionName)
        .authRefresh()
        .catch((err) => {
            console.warn('Failed to refresh the existing auth token:', err);
            // clear the store only on invalidated/expired token
            const status = err?.status << 0;
            if (status == 401 || status == 403) {
                pb.authStore.clear();
            }
        });
}

export default pb;

export function fetchMessages() {
    return pb.collection('messages').getFullList();
}

export class Helper {
    static isObject(value: any) {
        return value !== null && typeof value === "object" && value.constructor === Object;
    }
    static isEmpty(value: any) {
        return (
            (value === "") ||
            (value === null) ||
            (typeof value === "undefined") ||
            (Array.isArray(value) && value.length === 0) ||
            (Helper.isObject(value) && Object.keys(value).length === 0)
        );
    }
}

export interface pbMessage {
    id: string;
    userId: string;
    senderId: string;
    title: string;
    content: string;
    isRead: boolean;
    isArchived: boolean;
    isDeleted: boolean;
    created: Date;
    updated: Date;
}

export async function listMessages(): Promise<pbMessage[]> {
    var msgs: pbMessage[] = [];
    await pb.collection('messages').getFullList()
        .then((res) => {
            res.forEach((msg: any) => {
                msgs.push(msg as pbMessage);
            });
        })
        .catch((err) => {
            console.error('Failed to list messages:', err)
        });
    return msgs;
}

export async function readMessage(id: string) {
    await pb.collection('messages').update(id, { is_read: true })
        .catch((err) => {
            console.error('Failed to read message:', err)
        });
}

export async function unreadMessage(id: string) {
    await pb.collection('messages').update(id, { is_read: false })
        .catch((err) => {
            console.error('Failed to unread message:', err)
        });
}

export async function archiveMessage(id: string) {
    await pb.collection('messages').update(id, { is_archived: true })
        .catch((err) => {
            console.error('Failed to archive message:', err)
        });
}

export async function unarchiveMessage(id: string) {
    await pb.collection('messages').update(id, { is_archived: false })
        .catch((err) => {
            console.error('Failed to unarchive message:', err)
        });
}

export async function deleteMessage(id: string) {
    await pb.collection('messages').update(id, { is_deleted: true })
        .catch((err) => {
            console.error('Failed to delete message:', err)
        });
}

export interface pbFavorite {
    id: string;
    userId: string;
    path: string;
    name: string;
    created: Date;
    updated: Date;
}

export async function listFavorites(): Promise<pbFavorite[]> {
    var favs: pbFavorite[] = [];
    await pb.collection('favorites').getFullList()
        .then((res) => {
            res.forEach((fav: any) => {
                favs.push(fav as pbFavorite);
            });
        })
        .catch((err) => {
            console.error('Failed to list messages:', err)
        });
    return favs;
}

export async function isFavorite(): Promise<boolean> {
    if (pb.authStore.record) {
        return (await pb.collection('favorites').getFullList())
            .filter((fav: any) => fav.path == page.url.pathname).length > 0;
    }
    return false
}

export async function addFavorite() {
    if (pb.authStore.record) {
        await pb.collection('favorites').create({ user_id: pb.authStore.record.id, path: page.url.pathname, name: document.title || 'Untitled' })
            .catch((err) => {
                console.error('Failed to add favorite:', err)
            });
    }
}

export async function deleteFavorite() {
    if (pb.authStore.record) {
        (await pb.collection('favorites').getFullList())
            .filter((fav: any) => fav.path == page.url.pathname)
            .forEach(async (fav: any) => {
                await pb.collection('favorites').delete(fav.id)
                    .catch((err) => {
                        console.error('Failed to delete favorite:', err)
                    });
            });
    }
}

export interface pbVisit {
    id: string;
    userId: string;
    path: string;
    name: string;
    visited: Date;
    created: Date;
    updated: Date;
}

async function addVisit() {
    if (pb.authStore.record) {
        await pb.collection('visits').create({ user_id: pb.authStore.record.id, path: page.url.pathname, visited: new Date().toUTCString() })
            .catch((err) => {
                console.error('Failed to insert visit:', err)
            });
    }
}

async function updateVisit(id: string) {
    await pb.collection('visits').update(id, { visited: new Date().toUTCString() })
        .catch((err) => {
            console.error('Failed to update visit:', err)
        });
}

export async function upsertVisit() {
    var exists = false;
    (await pb.collection('visits').getFullList())
        .filter((v: any) => v.path == page.url.pathname)
        .forEach(async (v: any) => {
            exists = true;
            await updateVisit(v.id);
        })
    if (!exists) {
        await addVisit();
    }
}