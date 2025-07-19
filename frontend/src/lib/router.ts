import { createRouter } from '@tanstack/react-router';
import { queryClient } from '~/lib/query-client';
import { routeTree } from '~/routeTree.gen';

export const router = createRouter({
  routeTree,
  context: {
    queryClient,
  },
  defaultPreload: 'intent',
  scrollRestoration: true,
  defaultPreloadStaleTime: 0,
});

declare module '@tanstack/react-router' {
  interface Register {
    router: typeof router;
  }
}
