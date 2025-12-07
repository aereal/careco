import { client } from '@/lib/auth0';
import { Result } from '@praha/byethrow';
import { NextRequest, NextResponse } from 'next/server';

const upstreamEndpoint = 'http://localhost:8080/graphql';

class InternalUnauthenticatedError extends Error {
  get message(): string {
    return '[internal] unauthenticated';
  }
}

const getAccessToken = Result.try({
  async try() {
    const session = await client.getSession();
    if (!session) {
      throw new InternalUnauthenticatedError();
    }
    const {
      tokenSet: { accessToken },
    } = session;
    return accessToken;
  },
  catch(err: unknown): InternalUnauthenticatedError | Error {
    if (err instanceof InternalUnauthenticatedError) {
      return err;
    } else if (err instanceof Error) {
      return err;
    } else {
      return new Error(`${err}`);
    }
  },
});

interface RequestBody {
  body: unknown;
}

const getRequestBody = Result.try({
  try: async (req: NextRequest) => ({ body: await req.json() }) as RequestBody,
  catch: (err: unknown) => (err instanceof Error ? err : new Error(`${err}`)),
});

const doRequest = Result.try({
  async try(args: { headers: Headers; downstreamBody: RequestBody }) {
    const init: RequestInit = {
      method: 'POST',
      headers: args.headers,
      body: JSON.stringify(args.downstreamBody.body),
    };
    const upstreamResp = await fetch(upstreamEndpoint, init);
    const upstreamBody = await upstreamResp.json();
    return NextResponse.json<unknown>(upstreamBody, {
      status: upstreamResp.status,
      headers: upstreamResp.headers,
    });
  },
  catch: (err: unknown) => (err instanceof Error ? err : new Error(`${err}`)),
});

export const POST = async (
  req: NextRequest,
): Promise<NextResponse<unknown>> => {
  const ret = await Result.pipe(
    Result.sequence({
      downstreamBody: getRequestBody(req),
      headers: Result.pipe(
        Result.sequence({
          headers: Result.succeed(req.headers),
          accessToken: getAccessToken(),
        }),
        Result.inspect(({ accessToken, headers }) => {
          headers.set('authorization', `Bearer ${accessToken}`);
        }),
        Result.map(({ headers }) => headers),
      ),
    }),
    Result.andThen(doRequest),
    Result.orElse((err) =>
      Result.succeed(
        NextResponse.json<unknown>({ error: `${err}` }, { status: 401 }),
      ),
    ),
  );
  if (Result.isFailure(ret)) {
    console.error(`! unexpected error: ${ret.error}`);
    return NextResponse.json(
      { error: 'unexpected error occurred' },
      { status: 500 },
    );
  }
  return ret.value;
};
