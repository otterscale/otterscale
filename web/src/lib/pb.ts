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

export default pb;