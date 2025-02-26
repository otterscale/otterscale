import { page } from '$app/state';
import PocketBase, { type RecordAuthResponse, type RecordModel } from 'pocketbase';
import { siteConfig } from './config/site';
import { i18n } from './i18n';

const pb = new PocketBase('http://192.168.43.102:8090');
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

export async function listAuthMethods(): Promise<string[]> {
    var providers: string[] = [];
    await pb.collection('users').listAuthMethods().then((res) => {
        if (res.oauth2.enabled) {
            res.oauth2.providers.forEach((provider: any) => {
                providers.push(provider.name);
            });
        }
    }).catch((err) => {
        console.error('Failed to list auth methods:', err)
    });
    return providers
}
export async function passwordAuth(email: string, password: string): Promise<RecordAuthResponse<RecordModel>> {
    return await pb.collection('users').authWithPassword(email, password)
}

export async function oauth2Auth(provider: string): Promise<RecordAuthResponse<RecordModel>> {
    return await pb.collection('users').authWithOAuth2({ provider: provider });
}

export async function setEmailVisible(userId: string) {
    if (pb.authStore.record) {
        await pb.collection('users').update(userId, { emailVisibility: true })
            .catch((err) => {
                console.error('Failed to set email visible:', err)
            });;
    }
}

export async function createUser(email: string, password: string, passwordConfirm: string, name: string) {
    await pb.collection('users').create({
        email,
        password,
        passwordConfirm,
        name,
    }).catch((err) => {
        console.error('Failed to create user:', err)
    });
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
                msgs.push({
                    id: msg.id,
                    userId: msg.user_id,
                    senderId: msg.sender_id,
                    title: msg.title,
                    content: msg.content,
                    isRead: msg.is_read,
                    isArchived: msg.is_archived,
                    isDeleted: msg.is_deleted,
                    created: msg.created,
                    updated: msg.updated
                } as pbMessage);
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

export async function welcomeMessage(userId: string) {
    await pb.collection('messages').create({
        user_id: userId,
        sender_id: siteConfig.name,
        title: 'Welcome to our platform!',
        content: 'Your account has been created successfully. Enjoy your stay!',
    })
        .catch((err) => {
            console.error('Failed to add message:', err)
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
                favs.push({
                    id: fav.id,
                    userId: fav.user_id,
                    path: fav.path,
                    name: fav.name,
                    created: fav.created,
                    updated: fav.updated,
                } as pbFavorite);
            });
        })
        .catch((err) => {
            console.error('Failed to list favorites:', err)
        });
    return favs;
}

export async function isFavorite(): Promise<boolean> {
    if (pb.authStore.record) {
        return (await pb.collection('favorites').getFullList())
            .filter((fav: any) => fav.path == i18n.route(page.url.pathname)).length > 0;
    }
    return false
}

export async function addFavorite() {
    if (pb.authStore.record) {
        await pb.collection('favorites').create({ user_id: pb.authStore.record.id, path: i18n.route(page.url.pathname), name: document.title || 'Untitled' })
            .catch((err) => {
                console.error('Failed to add favorite:', err)
            });
    }
}

export async function deleteFavorite() {
    if (pb.authStore.record) {
        (await pb.collection('favorites').getFullList())
            .filter((fav: any) => fav.path == i18n.route(page.url.pathname))
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
        await pb.collection('visits').create({ user_id: pb.authStore.record.id, path: i18n.route(page.url.pathname), name: document.title || 'Untitled', visited: new Date().toUTCString() })
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
        .filter((v: any) => v.path == i18n.route(page.url.pathname))
        .forEach(async (v: any) => {
            exists = true;
            await updateVisit(v.id);
        })
    if (!exists) {
        await addVisit();
    }
}

export async function listVisits(): Promise<pbVisit[]> {
    var visits: pbVisit[] = [];
    await pb.collection('visits').getFullList()
        .then((res) => {
            res.forEach((visit: any) => {
                visits.push({
                    id: visit.id,
                    userId: visit.userId,
                    path: visit.path,
                    name: visit.name,
                    visited: visit.visited,
                    created: visit.created,
                    updated: visit.updated,
                } as pbVisit);
            });
        })
        .catch((err) => {
            console.error('Failed to list visits:', err)
        });
    return visits;
}
