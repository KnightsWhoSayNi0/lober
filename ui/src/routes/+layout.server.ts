import { env } from '$env/dynamic/public';

export async function load() {
    return {
        masterToken: env.PUBLIC_MASTER_TOKEN || 'dev-master-token'
    };
}
