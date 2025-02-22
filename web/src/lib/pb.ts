import PocketBase from 'pocketbase';

const pb = new PocketBase('http://127.0.0.1:8090');

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
    await pb.collection('messages').getList()
        .then((res) => {
            res.items.forEach((msg: any) => {
                msgs.push({
                    id: msg.id,
                    userId: msg.user_id,
                    senderId: msg.sender_id,
                    title: msg.title,
                    content: msg.content,
                    isRead: msg.is_read,
                    isArchived: msg.is_archived,
                    isDeleted: msg.is_deleted,
                    created: new Date(msg.created),
                    updated: new Date(msg.updated)
                });
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
            console.error('Failed to read message:', err)
        });
}

export async function archiveMessage(id: string) {
    await pb.collection('messages').update(id, { is_archived: true })
        .catch((err) => {
            console.error('Failed to read message:', err)
        });
}

export async function unarchiveMessage(id: string) {
    await pb.collection('messages').update(id, { is_archived: false })
        .catch((err) => {
            console.error('Failed to read message:', err)
        });
}

export async function deleteMessage(id: string) {
    await pb.collection('messages').update(id, { is_deleted: true })
        .catch((err) => {
            console.error('Failed to read message:', err)
        });
}