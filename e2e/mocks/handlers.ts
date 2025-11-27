import { graphql, HttpResponse } from 'msw';
import { mockGetRootData } from '../fixtures/get-root';
import { type GetRootQuery } from '../types';

/**
 * MSWハンドラー: GraphQLリクエストをモック
 * Server Component、Client Component両方のリクエストをインターセプト
 */
export const handlers = [
  graphql.query<GetRootQuery>('GetRoot', () => {
    return HttpResponse.json({
      data: mockGetRootData,
    });
  }),
];
