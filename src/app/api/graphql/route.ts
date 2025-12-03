import { NextRequest, NextResponse } from 'next/server';

const upstreamEndpoint = 'http://localhost:8080/graphql';

export const POST = async (
  req: NextRequest,
): Promise<NextResponse<unknown>> => {
  const downstreamBody = await req.json();
  const headers = new Headers(req.headers);
  headers.set('content-type', 'application/json');
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
