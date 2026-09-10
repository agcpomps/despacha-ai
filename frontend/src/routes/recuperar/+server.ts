import type { RequestHandler } from './$types';

const retired = () => new Response(null, { status: 410 });

/**
 * This legacy recovery path is deliberately unavailable while a secure,
 * first-party password-reset flow is designed and implemented.
 */
export const GET: RequestHandler = retired;
export const POST: RequestHandler = retired;
export const PUT: RequestHandler = retired;
export const PATCH: RequestHandler = retired;
export const DELETE: RequestHandler = retired;
