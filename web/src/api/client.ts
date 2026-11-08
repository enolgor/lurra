import createClient from 'openapi-fetch';
import type { paths } from './schema';

// The one API client. Use it directly (api.GET('/todos'), ...) and add shared behavior,
// such as auth headers, with api.use(middleware).
// Absolute URL: browsers accept '/api/v1', but Node's Request (used in tests) does not.
export const api = createClient<paths>({ baseUrl: `${location.origin}/api/v1` });
