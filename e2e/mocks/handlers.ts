import { graphql, HttpResponse } from 'msw';
import { mockGetRootData } from '../fixtures/get-root';
import type { GetRootMockData } from '../types/mock-types';

/**
 * MSWハンドラー: GraphQLリクエストをモック
 * Server Component、Client Component両方のリクエストをインターセプト
 */
export const handlers = [
  graphql.query<GetRootMockData>('GetRoot', () => {
    return HttpResponse.json({
      data: mockGetRootData,
    });
  }),
];
