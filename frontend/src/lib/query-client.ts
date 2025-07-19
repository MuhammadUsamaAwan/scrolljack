import { MutationCache, QueryClient, type QueryKey } from '@tanstack/react-query';
import { toast } from 'sonner';

declare module '@tanstack/react-query' {
  interface Register {
    mutationMeta: {
      invalidateQueries?: QueryKey;
      successMessage?: string;
      errorMessage?: string;
    };
  }
}

export const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: false,
      refetchOnWindowFocus: false,
    },
    mutations: {
      retry: false,
    },
  },
  mutationCache: new MutationCache({
    onSuccess(_data, _variables, _context, mutation) {
      if (mutation.meta?.successMessage) {
        toast.success(mutation.meta.successMessage);
      }
      if (mutation.meta?.invalidateQueries) {
        queryClient.invalidateQueries({
          queryKey: mutation.meta.invalidateQueries,
        });
      }
    },
    onError(_error, _variables, _context, mutation) {
      if (mutation.meta?.errorMessage) {
        toast.error(mutation.meta.errorMessage);
      }
    },
  }),
});
