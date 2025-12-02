import { client } from '@/lib/auth0';
import { NextProxy } from 'next/server';

export const proxy: NextProxy = async (req) => await client.middleware(req);
