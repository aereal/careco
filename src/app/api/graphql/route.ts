import { client } from '@/lib/auth0';
import { NextRequest, NextResponse } from 'next/server';

const upstreamEndpoint = 'http://localhost:8080/graphql';

class InternalUnauthenticatedError extends Error {
  get message(): string {
    return '[internal] unauthenticated';
  }
}

export const POST = async (
  req: NextRequest,
): Promise<NextResponse<unknown>> => {
  const downstreamBody = await req.json();
  const headers = new Headers(req.headers);
  headers.set('content-type', 'application/json');
  try {
    const session = await client.getSession();
    if (!session) {
      throw new InternalUnauthenticatedError();
    }
    const {
      tokenSet: { accessToken },
    } = session;
    headers.set('authorization', `Bearer ${accessToken}`);
  } catch {}
  const init: RequestInit = {
    method: 'POST',
    headers,
    body: JSON.stringify(downstreamBody),
  };
  const upstreamResp = await fetch(upstreamEndpoint, init);
  const upstreamBody = await upstreamResp.json();
  return NextResponse.json(upstreamBody, {
    status: upstreamResp.status,
    headers: upstreamResp.headers,
  });
};
