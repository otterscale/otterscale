import { page } from '$app/state';
import PocketBase, { type RecordAuthResponse, type RecordModel } from 'pocketbase';
import { i18n } from './i18n';

const pb = new PocketBase("http://192.168.197.171:8090/");
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

export function avatarFallback(name: string): string {
    if (name) {
        const names = name.split(' ');
        if (names.length >= 2) {
            return (names[0][0] + names[1][0]).toUpperCase();
        }
    }
    return 'NA';
}

export function avatarURL(avatar: string): string {
    if (pb.authStore.record) {
        return pb.files.getURL(pb.authStore.record, avatar);
    }
    return '';
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
    from: string;
    title: string;
    content: string;
    read: boolean;
    archived: boolean;
    deleted: boolean;
    created: Date;
    updated: Date;
}

export async function listMessages(): Promise<pbMessage[]> {
    var msgs: pbMessage[] = [];
    await pb.collection('messages').getFullList({ expand: "from" })
        .then((res) => {
            res.forEach((msg: any) => {
                msgs.push({
                    id: msg.id,
                    from: msg.expand?.from?.name,
                    title: msg.title,
                    content: msg.content,
                    read: msg.read,
                    archived: msg.archived,
                    deleted: msg.deleted,
                    created: msg.created,
                    updated: msg.updated,
                } as pbMessage);
            });
        })
        .catch((err) => {
            console.error('Failed to list messages:', err)
        });
    return msgs;
}

export async function readMessage(id: string) {
    await pb.collection('messages').update(id, { read: true })
        .catch((err) => {
            console.error('Failed to read message:', err)
        });
}

export async function unreadMessage(id: string) {
    await pb.collection('messages').update(id, { read: false })
        .catch((err) => {
            console.error('Failed to unread message:', err)
        });
}

export async function archiveMessage(id: string) {
    await pb.collection('messages').update(id, { archived: true })
        .catch((err) => {
            console.error('Failed to archive message:', err)
        });
}

export async function unarchiveMessage(id: string) {
    await pb.collection('messages').update(id, { archived: false })
        .catch((err) => {
            console.error('Failed to unarchive message:', err)
        });
}

export async function deleteMessage(id: string) {
    await pb.collection('messages').update(id, { deleted: true })
        .catch((err) => {
            console.error('Failed to delete message:', err)
        });
}

export async function welcomeMessage(userId: string) {
    if (pb.authStore.record) {
        await pb.collection('messages').create({
            user: userId,
            title: 'Welcome to our platform!',
            content: 'Your account has been created successfully. Enjoy your stay!',
        })
            .catch((err) => {
                console.error('Failed to add message:', err)
            });
    }

}

export interface pbFavorite {
    id: string;
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
        await pb.collection('favorites').create({ user: pb.authStore.record.id, path: i18n.route(page.url.pathname), name: document.title || 'Untitled' })
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

export interface pbRecent {
    id: string;
    path: string;
    name: string;
    visited: Date;
    created: Date;
    updated: Date;
}

async function addRecent() {
    if (pb.authStore.record) {
        await pb.collection('recents').create({ user: pb.authStore.record.id, path: i18n.route(page.url.pathname), name: document.title || 'Untitled', visited: new Date().toUTCString() })
            .catch((err) => {
                console.error('Failed to insert recent:', err)
            });
    }
}

async function updateRecent(id: string) {
    await pb.collection('recents').update(id, { visited: new Date().toUTCString() })
        .catch((err) => {
            console.error('Failed to update recent:', err)
        });
}

export async function upsertRecent() {
    var exists = false;
    (await pb.collection('recents').getFullList())
        .filter((v: any) => v.path == i18n.route(page.url.pathname))
        .forEach(async (v: any) => {
            exists = true;
            await updateRecent(v.id);
        })
    if (!exists) {
        await addRecent();
    }
}

export async function listRecents(): Promise<pbRecent[]> {
    var recents: pbRecent[] = [];
    await pb.collection('recents').getFullList()
        .then((res) => {
            res.forEach((rec: any) => {
                recents.push({
                    id: rec.id,
                    path: rec.path,
                    name: rec.name,
                    visited: rec.visited,
                    created: rec.created,
                    updated: rec.updated,
                } as pbRecent);
            });
        })
        .catch((err) => {
            console.error('Failed to list recents:', err)
        });
    return recents;
}

export interface pbWorkload {
    user: string;
    avatar: string;
    json: object;
    created: Date;
}

export interface pbConnector {
    id: string;
    kind: string;
    type: string;
    name: string;
    image: boolean;
    enabled: boolean;
    user: string;
    avatar: string;
    created: Date;
    updated: Date;
    workload: pbWorkload;
}

function toPbConnector(rec: any): pbConnector {
    var workload = rec.expand?.workloads_via_connector?.reduce((prev: any, curr: any) => (!prev || new Date(curr.created) > new Date(prev.created)) ? curr : prev, null);
    return {
        id: rec.id,
        kind: rec.kind,
        type: rec.type,
        name: rec.name,
        image: rec.image,
        enabled: rec.enabled,
        user: rec.expand?.user?.name,
        avatar: rec.expand?.user?.avatar,
        created: rec.created,
        updated: rec.updated,
        workload: {
            user: workload?.expand?.user?.name,
            avatar: workload?.expand?.user?.avatar,
            json: workload?.json,
            created: workload?.created,
        }
    } as pbConnector
}

export async function listConnectors(filter: string): Promise<pbConnector[]> {
    var connectors: pbConnector[] = [];
    await pb.collection('connectors').getFullList({ filter: filter, expand: "workloads_via_connector.user,user" })
        .then((res) => {
            res.forEach((rec: any) => {
                connectors.push(toPbConnector(rec));
            });
        })
        .catch((err) => {
            console.error('Failed to list connectors:', err)
        });
    return connectors;
}


export interface pbPipeline {
    id: string;
    source: pbConnector;
    destination: pbConnector;
    schedule: string;
    user: string;
    avatar: string;
    created: Date;
    updated: Date;
}

export async function listPipelines(filter: string): Promise<pbPipeline[]> {
    var pipelines: pbPipeline[] = [];
    await pb.collection('pipelines').getFullList({ filter: filter, expand: "user,source,destination,source.user,destination.user,source.workloads_via_connector.user,destination.workloads_via_connector.user" })
        .then((res) => {
            res.forEach((rec: any) => {
                if (rec.deleted) {
                    return
                }
                pipelines.push({
                    id: rec.id,
                    source: toPbConnector(rec.expand?.source),
                    destination: toPbConnector(rec.expand?.destination),
                    schedule: rec.schedule,
                    user: rec.expand?.user?.name,
                    avatar: rec.expand?.user?.avatar,
                    created: rec.created,
                    updated: rec.updated,
                } as pbPipeline)
            });
        })
        .catch((err) => {
            console.error('Failed to list pipelines:', err)
        });
    return pipelines;
}

export interface pbInstance {
    id: string;
    type: string;
    name: string;
    tags: string[]
    neighbors: string[];
    information: {},
    created: Date;
    updated: Date;
}
export async function listInstances(filter: string): Promise<pbInstance[]> {
    var instances: pbInstance[] = [];
    await pb.collection('nodes').getFullList({ filter: filter })
        .then((res) => {
            res.forEach((rec: any) => {
                instances.push({
                    id: rec.id,
                    type: rec.type,
                    name: rec.name,
                    tags: rec.tags,
                    neighbors: rec.neighbors,
                    information: rec.information,
                    created: rec.created,
                    updated: rec.updated,
                } as pbInstance)
            });
        })
        .catch((err) => {
            console.error('Failed to list nodes:', err)
        });
    return instances;
}
